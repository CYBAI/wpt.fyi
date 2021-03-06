// Code generated by MockGen. DO NOT EDIT.
// Source: api.go

package checks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockSuitesAPI is a mock of SuitesAPI interface
type MockSuitesAPI struct {
	ctrl     *gomock.Controller
	recorder *MockSuitesAPIMockRecorder
}

// MockSuitesAPIMockRecorder is the mock recorder for MockSuitesAPI
type MockSuitesAPIMockRecorder struct {
	mock *MockSuitesAPI
}

// NewMockSuitesAPI creates a new mock instance
func NewMockSuitesAPI(ctrl *gomock.Controller) *MockSuitesAPI {
	mock := &MockSuitesAPI{ctrl: ctrl}
	mock.recorder = &MockSuitesAPIMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockSuitesAPI) EXPECT() *MockSuitesAPIMockRecorder {
	return m.recorder
}

// Context mocks base method
func (m *MockSuitesAPI) Context() context.Context {
	ret := m.ctrl.Call(m, "Context")
	ret0, _ := ret[0].(context.Context)
	return ret0
}

// Context indicates an expected call of Context
func (mr *MockSuitesAPIMockRecorder) Context() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Context", reflect.TypeOf((*MockSuitesAPI)(nil).Context))
}

// PendingCheckRun mocks base method
func (m *MockSuitesAPI) PendingCheckRun(sha, browser string) (bool, error) {
	ret := m.ctrl.Call(m, "PendingCheckRun", sha, browser)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PendingCheckRun indicates an expected call of PendingCheckRun
func (mr *MockSuitesAPIMockRecorder) PendingCheckRun(sha, browser interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PendingCheckRun", reflect.TypeOf((*MockSuitesAPI)(nil).PendingCheckRun), sha, browser)
}

// CompleteCheckRun mocks base method
func (m *MockSuitesAPI) CompleteCheckRun(sha, browser string) (bool, error) {
	ret := m.ctrl.Call(m, "CompleteCheckRun", sha, browser)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CompleteCheckRun indicates an expected call of CompleteCheckRun
func (mr *MockSuitesAPIMockRecorder) CompleteCheckRun(sha, browser interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CompleteCheckRun", reflect.TypeOf((*MockSuitesAPI)(nil).CompleteCheckRun), sha, browser)
}
