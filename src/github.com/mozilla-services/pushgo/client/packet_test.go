/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package client

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/mozilla-services/pushgo/id"
)

var ErrInvalidCaseTest = &ClientError{"Invalid case test type."}

// CaseTestType is used by `Case{ClientPing, ACK}.MarshalJSON()` to generate
// different JSON representations of the underlying ping or ACK packet. If
// a packet doesn't support a particular representation (e.g., `CaseACK`
// doesn't support any of the `FieldType*` test types), its `MarshalJSON()`
// method will return ``ErrInvalidCaseTest`.
type CaseTestType int

const (
	FieldTypeLower CaseTestType = iota + 1
	FieldTypeSpace
	FieldChansCap
	FieldIdCap
	ValueTypeUpper
	ValueTypeCap
	ValueTypeEmpty
)

func (t CaseTestType) String() string {
	switch t {
	case FieldTypeLower:
		return "lowercase message type field name"
	case FieldTypeSpace:
		return "leading and trailing whitespace in message type field name"
	case FieldChansCap:
		return "mixed-case channel IDs field name"
	case FieldIdCap:
		return "mixed-case device ID field name"
	case ValueTypeUpper:
		return "uppercase message type value"
	case ValueTypeCap:
		return "mixed-case message type value"
	case ValueTypeEmpty:
		return "empty message type value"
	}
	return "unknown case test type"
}

func decodePing(c *Conn, fields Fields, statusCode int, errorText string) (HasType, error) {
	if len(errorText) > 0 {
		return nil, &ServerError{"ping", c.Origin(), errorText, statusCode}
	}
	return &ServerPing{statusCode}, nil
}

func decodeCaseACK(c *Conn, fields Fields, statusCode int, errorText string) (HasType, error) {
	if len(errorText) > 0 {
		return nil, &ServerError{"ack", c.Origin(), errorText, statusCode}
	}
	return nil, nil
}

type caseTest struct {
	CaseTestType
	statusCode int
	forceReset bool
}

func (t caseTest) TestPing() error {
	conn, err := DialOrigin(Origin)
	if err != nil {
		return fmt.Errorf("On test %v, error dialing origin: %#v", t.CaseTestType, err)
	}
	defer conn.Close()
	defer conn.Purge()
	conn.RegisterDecoder("ping", DecoderFunc(decodePing))
	request := &CaseClientPing{
		CaseTestType: t.CaseTestType,
		replies:      make(chan Reply),
		ClientPing:   make(ClientPing),
	}
	reply, err := conn.WriteRequest(request)
	if err != nil {
		return fmt.Errorf("On test %v, error writing ping packet: %#v", t.CaseTestType, err)
	}
	if reply.Status() != t.statusCode {
		return fmt.Errorf("On test %v, unexpected status code: got %#v; want %#v", t.CaseTestType, reply.Status(), t.statusCode)
	}
	return nil
}

func (t caseTest) TestACK() error {
	conn, err := DialOrigin(Origin)
	if err != nil {
		return fmt.Errorf("On test %v, error dialing origin: %#v", t.CaseTestType, err)
	}
	defer conn.Close()
	defer conn.Purge()
	conn.RegisterDecoder("ack", DecoderFunc(decodeCaseACK))
	request := &CaseACK{
		CaseTestType: t.CaseTestType,
		replies:      make(chan Reply),
		ClientACK: ClientACK{
			errors: make(chan error),
		},
	}
	_, err = conn.WriteRequest(request)
	if t.statusCode >= 200 && t.statusCode < 300 {
		if err != nil {
			return fmt.Errorf("On test %v, error writing acknowledgement: %#v", t.CaseTestType, err)
		}
		return nil
	}
	if err != io.EOF {
		return fmt.Errorf("On test %v, error writing acknowledgement: got %#v; want io.EOF", t.CaseTestType, err)
	}
	err = conn.Close()
	clientErr, ok := err.(Error)
	if !ok {
		return fmt.Errorf("On test %v, type assertion failed for close error: %#v", t.CaseTestType, err)
	}
	if clientErr.Status() != t.statusCode {
		return fmt.Errorf("On test %v, unexpected close error status: got %#v; want %#v", t.CaseTestType, clientErr.Status(), t.statusCode)
	}
	return nil
}

func (t caseTest) TestHelo() error {
	deviceId, err := id.Generate()
	if err != nil {
		return fmt.Errorf("On test %v, error generating device ID: %#v", t.CaseTestType, err)
	}
	channelId, err := id.Generate()
	if err != nil {
		return fmt.Errorf("On test %v, error generating channel ID: %#v", t.CaseTestType, err)
	}
	conn, err := DialOrigin(Origin)
	if err != nil {
		return fmt.Errorf("On test %v, error dialing origin: %#v", t.CaseTestType, err)
	}
	defer conn.Close()
	defer conn.Purge()
	request := &CaseHelo{
		t.CaseTestType,
		ClientHelo{
			DeviceId:   deviceId,
			ChannelIds: []string{channelId},
			replies:    make(chan Reply),
			errors:     make(chan error),
		},
	}
	reply, err := conn.WriteRequest(request)
	if t.statusCode >= 200 && t.statusCode < 300 {
		if err != nil {
			return fmt.Errorf("On test %v, error writing handshake request: %#v", t.CaseTestType, err)
		}
		helo, ok := reply.(*ServerHelo)
		if !ok {
			return fmt.Errorf("On test %v, type assertion failed for handshake reply: %#v", t.CaseTestType, reply)
		}
		if helo.StatusCode != t.statusCode {
			return fmt.Errorf("On test %v, unexpected reply status: got %#v; want %#v", t.CaseTestType, helo.StatusCode, t.statusCode)
		}
		if t.forceReset {
			if helo.DeviceId == deviceId {
				return fmt.Errorf("On test %v, want new device ID; got %#v", t.CaseTestType, deviceId)
			}
			return nil
		}
		if helo.DeviceId != deviceId {
			return fmt.Errorf("On test %v, mismatched device ID: got %#v; want %#v", t.CaseTestType, helo.DeviceId, deviceId)
		}
		return nil
	}
	if err != io.EOF {
		return fmt.Errorf("On test %v, error writing handshake: got %#v; want io.EOF", t.CaseTestType, err)
	}
	err = conn.Close()
	clientErr, ok := err.(Error)
	if !ok {
		return fmt.Errorf("On test %v, type assertion failed for close error: %#v", t.CaseTestType, err)
	}
	if clientErr.Status() != t.statusCode {
		return fmt.Errorf("On test %v, unexpected close error status: got %#v; want %#v", t.CaseTestType, clientErr.Status(), t.statusCode)
	}
	return nil
}

type CaseACK struct {
	CaseTestType
	replies chan Reply
	ClientACK
}

func (*CaseACK) CanReply() bool      { return true }
func (*CaseACK) Sync() bool          { return true }
func (a *CaseACK) Reply(reply Reply) { a.replies <- reply }

func (a *CaseACK) Close() {
	a.ClientACK.Close()
	close(a.replies)
}

func (a *CaseACK) Do() (reply Reply, err error) {
	select {
	case reply = <-a.replies:
	case err = <-a.getErrors():
	}
	return
}

func (a *CaseACK) MarshalJSON() ([]byte, error) {
	var results interface{}
	messageType := a.Type().String()
	switch a.CaseTestType {
	case ValueTypeUpper:
		results = struct {
			MessageType string   `json:"messageType"`
			Updates     []Update `json:"updates"`
		}{strings.ToUpper(messageType), a.Updates}

	case ValueTypeCap:
		results = struct {
			MessageType string   `json:"messageType"`
			Updates     []Update `json:"updates"`
		}{strings.ToUpper(messageType[:1]) + strings.ToLower(messageType[1:]), a.Updates}

	case ValueTypeEmpty:
		results = struct {
			MessageType string   `json:"messageType"`
			Updates     []Update `json:"updates"`
		}{"", a.Updates}

	default:
		return nil, ErrInvalidCaseTest
	}
	return json.Marshal(results)
}

type CaseClientPing struct {
	CaseTestType
	replies chan Reply
	ClientPing
}

func (p *CaseClientPing) CanReply() bool    { return true }
func (p *CaseClientPing) Sync() bool        { return true }
func (p *CaseClientPing) Reply(reply Reply) { p.replies <- reply }

func (p *CaseClientPing) Close() {
	p.ClientPing.Close()
	close(p.replies)
}

func (p *CaseClientPing) Do() (reply Reply, err error) {
	select {
	case reply = <-p.replies:
	case err = <-p.getErrors():
	}
	return
}

func (p *CaseClientPing) MarshalJSON() ([]byte, error) {
	switch p.CaseTestType {
	case ValueTypeUpper:
		return []byte(`{"messageType":"PING"}`), nil

	case ValueTypeCap:
		return []byte(`{"messageType":"Ping"}`), nil

	case ValueTypeEmpty:
		return []byte{'{', '}'}, nil
	}
	return nil, ErrInvalidCaseTest
}

type ServerPing struct {
	StatusCode int
}

func (*ServerPing) Type() PacketType { return Ping }
func (*ServerPing) Id() interface{}  { return PingId }
func (*ServerPing) HasRequest() bool { return true }
func (*ServerPing) Sync() bool       { return true }
func (p *ServerPing) Status() int    { return p.StatusCode }

type CaseHelo struct {
	CaseTestType
	ClientHelo
}

func (h *CaseHelo) MarshalJSON() ([]byte, error) {
	var results interface{}
	messageType := h.Type().String()
	switch h.CaseTestType {
	case FieldTypeLower:
		results = struct {
			MessageType string   `json:"messagetype"`
			DeviceId    string   `json:"uaid"`
			ChannelIds  []string `json:"channelIDs"`
		}{messageType, h.DeviceId, h.ChannelIds}

	case FieldTypeSpace:
		results = struct {
			MessageType string   `json:" messageType "`
			DeviceId    string   `json:"uaid"`
			ChannelIds  []string `json:"channelIDs"`
		}{messageType, h.DeviceId, h.ChannelIds}

	case FieldChansCap:
		results = struct {
			MessageType string   `json:"messageType"`
			DeviceId    string   `json:"uaid"`
			ChannelIds  []string `json:"ChannelIDs"`
		}{messageType, h.DeviceId, h.ChannelIds}

	case FieldIdCap:
		results = struct {
			MessageType string   `json:"messageType"`
			DeviceId    string   `json:"uaiD"`
			ChannelIds  []string `json:"channelIDs"`
		}{messageType, h.DeviceId, h.ChannelIds}

	case ValueTypeUpper:
		results = struct {
			MessageType string   `json:"messageType"`
			DeviceId    string   `json:"uaid"`
			ChannelIds  []string `json:"channelIDs"`
		}{strings.ToUpper(messageType), h.DeviceId, h.ChannelIds}

	case ValueTypeCap:
		results = struct {
			MessageType string   `json:"messageType"`
			DeviceId    string   `json:"uaid"`
			ChannelIds  []string `json:"channelIDs"`
		}{strings.ToUpper(messageType[:1]) + strings.ToLower(messageType[1:]), h.DeviceId, h.ChannelIds}

	case ValueTypeEmpty:
		results = struct {
			MessageType string   `json:"messageType"`
			DeviceId    string   `json:"uaid"`
			ChannelIds  []string `json:"channelIDs"`
		}{"", h.DeviceId, h.ChannelIds}

	default:
		return nil, ErrInvalidCaseTest
	}
	return json.Marshal(results)
}

var ackTests = []caseTest{
	caseTest{ValueTypeUpper, 401, false},
	caseTest{ValueTypeCap, 401, false},
	caseTest{ValueTypeEmpty, 401, false},
}

func TestACKCase(t *testing.T) {
	for _, test := range ackTests {
		if err := test.TestACK(); err != nil {
			t.Error(err)
		}
	}
}

var pingTests = []caseTest{
	caseTest{ValueTypeUpper, 200, false},
	caseTest{ValueTypeCap, 200, false},
	caseTest{ValueTypeEmpty, 200, false},
}

func TestPingCase(t *testing.T) {
	for _, test := range pingTests {
		if err := test.TestPing(); err != nil {
			t.Error(err)
		}
	}
}

var heloTests = []caseTest{
	caseTest{FieldTypeLower, 401, false},
	caseTest{FieldTypeSpace, 401, false},
	caseTest{FieldChansCap, 401, false},
	caseTest{FieldIdCap, 200, true},
	caseTest{ValueTypeUpper, 200, true},
	caseTest{ValueTypeCap, 200, true},
	caseTest{ValueTypeEmpty, 401, false},
}

func TestHeloCase(t *testing.T) {
	for _, test := range heloTests {
		if err := test.TestHelo(); err != nil {
			t.Error(err)
		}
	}
}