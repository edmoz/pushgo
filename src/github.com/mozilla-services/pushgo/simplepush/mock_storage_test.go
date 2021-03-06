// Automatically generated by MockGen. DO NOT EDIT!
// Source: src/github.com/mozilla-services/pushgo/simplepush/storage.go

package simplepush

import (
	gomock "github.com/rafrombrc/gomock/gomock"
	time "time"
)

// Mock of Store interface
type MockStore struct {
	ctrl     *gomock.Controller
	recorder *_MockStoreRecorder
}

// Recorder for MockStore (not exported)
type _MockStoreRecorder struct {
	mock *MockStore
}

func NewMockStore(ctrl *gomock.Controller) *MockStore {
	mock := &MockStore{ctrl: ctrl}
	mock.recorder = &_MockStoreRecorder{mock}
	return mock
}

func (_m *MockStore) EXPECT() *_MockStoreRecorder {
	return _m.recorder
}

func (_m *MockStore) CanStore(channels int) bool {
	ret := _m.ctrl.Call(_m, "CanStore", channels)
	ret0, _ := ret[0].(bool)
	return ret0
}

func (_mr *_MockStoreRecorder) CanStore(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "CanStore", arg0)
}

func (_m *MockStore) Close() error {
	ret := _m.ctrl.Call(_m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockStoreRecorder) Close() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Close")
}

func (_m *MockStore) KeyToIDs(key string) (string, string, error) {
	ret := _m.ctrl.Call(_m, "KeyToIDs", key)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

func (_mr *_MockStoreRecorder) KeyToIDs(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "KeyToIDs", arg0)
}

func (_m *MockStore) IDsToKey(suaid string, schid string) (string, error) {
	ret := _m.ctrl.Call(_m, "IDsToKey", suaid, schid)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockStoreRecorder) IDsToKey(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "IDsToKey", arg0, arg1)
}

func (_m *MockStore) Status() (bool, error) {
	ret := _m.ctrl.Call(_m, "Status")
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockStoreRecorder) Status() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Status")
}

func (_m *MockStore) Exists(suaid string) bool {
	ret := _m.ctrl.Call(_m, "Exists", suaid)
	ret0, _ := ret[0].(bool)
	return ret0
}

func (_mr *_MockStoreRecorder) Exists(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Exists", arg0)
}

func (_m *MockStore) Register(suaid string, schid string, version int64) error {
	ret := _m.ctrl.Call(_m, "Register", suaid, schid, version)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockStoreRecorder) Register(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Register", arg0, arg1, arg2)
}

func (_m *MockStore) Update(suaid string, schid string, version int64) error {
	ret := _m.ctrl.Call(_m, "Update", suaid, schid, version)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockStoreRecorder) Update(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Update", arg0, arg1, arg2)
}

func (_m *MockStore) Unregister(suaid string, schid string) error {
	ret := _m.ctrl.Call(_m, "Unregister", suaid, schid)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockStoreRecorder) Unregister(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Unregister", arg0, arg1)
}

func (_m *MockStore) Drop(suaid string, schid string) error {
	ret := _m.ctrl.Call(_m, "Drop", suaid, schid)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockStoreRecorder) Drop(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Drop", arg0, arg1)
}

func (_m *MockStore) FetchAll(suaid string, since time.Time) ([]Update, []string, error) {
	ret := _m.ctrl.Call(_m, "FetchAll", suaid, since)
	ret0, _ := ret[0].([]Update)
	ret1, _ := ret[1].([]string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

func (_mr *_MockStoreRecorder) FetchAll(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "FetchAll", arg0, arg1)
}

func (_m *MockStore) DropAll(suaid string) error {
	ret := _m.ctrl.Call(_m, "DropAll", suaid)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockStoreRecorder) DropAll(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "DropAll", arg0)
}

func (_m *MockStore) FetchPing(suaid string) ([]byte, error) {
	ret := _m.ctrl.Call(_m, "FetchPing", suaid)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockStoreRecorder) FetchPing(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "FetchPing", arg0)
}

func (_m *MockStore) PutPing(suaid string, pingData []byte) error {
	ret := _m.ctrl.Call(_m, "PutPing", suaid, pingData)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockStoreRecorder) PutPing(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "PutPing", arg0, arg1)
}

func (_m *MockStore) DropPing(suaid string) error {
	ret := _m.ctrl.Call(_m, "DropPing", suaid)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockStoreRecorder) DropPing(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "DropPing", arg0)
}
