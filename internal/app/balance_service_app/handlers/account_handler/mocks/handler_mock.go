// Code generated by MockGen. DO NOT EDIT.
// Source: interface.go

// Package mocks is a generated GoMock package.
package mocks

import (
	http "net/http"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockAccountHandlerInterface is a mock of AccountHandlerInterface interface.
type MockAccountHandlerInterface struct {
	ctrl     *gomock.Controller
	recorder *MockAccountHandlerInterfaceMockRecorder
}

// MockAccountHandlerInterfaceMockRecorder is the mock recorder for MockAccountHandlerInterface.
type MockAccountHandlerInterfaceMockRecorder struct {
	mock *MockAccountHandlerInterface
}

// NewMockAccountHandlerInterface creates a new mock instance.
func NewMockAccountHandlerInterface(ctrl *gomock.Controller) *MockAccountHandlerInterface {
	mock := &MockAccountHandlerInterface{ctrl: ctrl}
	mock.recorder = &MockAccountHandlerInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAccountHandlerInterface) EXPECT() *MockAccountHandlerInterfaceMockRecorder {
	return m.recorder
}

// GetBalance mocks base method.
func (m *MockAccountHandlerInterface) GetBalance(w http.ResponseWriter, r *http.Request) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "GetBalance", w, r)
}

// GetBalance indicates an expected call of GetBalance.
func (mr *MockAccountHandlerInterfaceMockRecorder) GetBalance(w, r interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBalance", reflect.TypeOf((*MockAccountHandlerInterface)(nil).GetBalance), w, r)
}

// RefillBalance mocks base method.
func (m *MockAccountHandlerInterface) RefillBalance(w http.ResponseWriter, r *http.Request) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RefillBalance", w, r)
}

// RefillBalance indicates an expected call of RefillBalance.
func (mr *MockAccountHandlerInterfaceMockRecorder) RefillBalance(w, r interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RefillBalance", reflect.TypeOf((*MockAccountHandlerInterface)(nil).RefillBalance), w, r)
}

// Transfer mocks base method.
func (m *MockAccountHandlerInterface) Transfer(w http.ResponseWriter, r *http.Request) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Transfer", w, r)
}

// Transfer indicates an expected call of Transfer.
func (mr *MockAccountHandlerInterfaceMockRecorder) Transfer(w, r interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Transfer", reflect.TypeOf((*MockAccountHandlerInterface)(nil).Transfer), w, r)
}