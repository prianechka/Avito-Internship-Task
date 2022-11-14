// Code generated by MockGen. DO NOT EDIT.
// Source: interface.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockAccountControllerInterface is a mock of AccountControllerInterface interface.
type MockAccountControllerInterface struct {
	ctrl     *gomock.Controller
	recorder *MockAccountControllerInterfaceMockRecorder
}

// MockAccountControllerInterfaceMockRecorder is the mock recorder for MockAccountControllerInterface.
type MockAccountControllerInterfaceMockRecorder struct {
	mock *MockAccountControllerInterface
}

// NewMockAccountControllerInterface creates a new mock instance.
func NewMockAccountControllerInterface(ctrl *gomock.Controller) *MockAccountControllerInterface {
	mock := &MockAccountControllerInterface{ctrl: ctrl}
	mock.recorder = &MockAccountControllerInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAccountControllerInterface) EXPECT() *MockAccountControllerInterfaceMockRecorder {
	return m.recorder
}

// CheckAbleToBuyService mocks base method.
func (m *MockAccountControllerInterface) CheckAbleToBuyService(userID int64, servicePrice float64) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckAbleToBuyService", userID, servicePrice)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckAbleToBuyService indicates an expected call of CheckAbleToBuyService.
func (mr *MockAccountControllerInterfaceMockRecorder) CheckAbleToBuyService(userID, servicePrice interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckAbleToBuyService", reflect.TypeOf((*MockAccountControllerInterface)(nil).CheckAbleToBuyService), userID, servicePrice)
}

// CheckAccountIsExist mocks base method.
func (m *MockAccountControllerInterface) CheckAccountIsExist(userID int64) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckAccountIsExist", userID)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckAccountIsExist indicates an expected call of CheckAccountIsExist.
func (mr *MockAccountControllerInterfaceMockRecorder) CheckAccountIsExist(userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckAccountIsExist", reflect.TypeOf((*MockAccountControllerInterface)(nil).CheckAccountIsExist), userID)
}

// CheckBalance mocks base method.
func (m *MockAccountControllerInterface) CheckBalance(userID int64) (float64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckBalance", userID)
	ret0, _ := ret[0].(float64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckBalance indicates an expected call of CheckBalance.
func (mr *MockAccountControllerInterfaceMockRecorder) CheckBalance(userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckBalance", reflect.TypeOf((*MockAccountControllerInterface)(nil).CheckBalance), userID)
}

// CreateNewAccount mocks base method.
func (m *MockAccountControllerInterface) CreateNewAccount(userID int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateNewAccount", userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateNewAccount indicates an expected call of CreateNewAccount.
func (mr *MockAccountControllerInterfaceMockRecorder) CreateNewAccount(userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateNewAccount", reflect.TypeOf((*MockAccountControllerInterface)(nil).CreateNewAccount), userID)
}

// DonateMoney mocks base method.
func (m *MockAccountControllerInterface) DonateMoney(userID int64, sum float64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DonateMoney", userID, sum)
	ret0, _ := ret[0].(error)
	return ret0
}

// DonateMoney indicates an expected call of DonateMoney.
func (mr *MockAccountControllerInterfaceMockRecorder) DonateMoney(userID, sum interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DonateMoney", reflect.TypeOf((*MockAccountControllerInterface)(nil).DonateMoney), userID, sum)
}

// SpendMoney mocks base method.
func (m *MockAccountControllerInterface) SpendMoney(userID int64, sum float64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SpendMoney", userID, sum)
	ret0, _ := ret[0].(error)
	return ret0
}

// SpendMoney indicates an expected call of SpendMoney.
func (mr *MockAccountControllerInterfaceMockRecorder) SpendMoney(userID, sum interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SpendMoney", reflect.TypeOf((*MockAccountControllerInterface)(nil).SpendMoney), userID, sum)
}
