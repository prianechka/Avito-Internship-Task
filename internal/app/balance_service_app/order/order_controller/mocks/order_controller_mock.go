// Code generated by MockGen. DO NOT EDIT.
// Source: interface.go

// Package mocks is a generated GoMock package.
package mocks

import (
	order "Avito-Internship-Task/internal/app/balance_service_app/order"
	report "Avito-Internship-Task/internal/app/balance_service_app/report"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockOrderControllerInterface is a mock of OrderControllerInterface interface.
type MockOrderControllerInterface struct {
	ctrl     *gomock.Controller
	recorder *MockOrderControllerInterfaceMockRecorder
}

// MockOrderControllerInterfaceMockRecorder is the mock recorder for MockOrderControllerInterface.
type MockOrderControllerInterfaceMockRecorder struct {
	mock *MockOrderControllerInterface
}

// NewMockOrderControllerInterface creates a new mock instance.
func NewMockOrderControllerInterface(ctrl *gomock.Controller) *MockOrderControllerInterface {
	mock := &MockOrderControllerInterface{ctrl: ctrl}
	mock.recorder = &MockOrderControllerInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockOrderControllerInterface) EXPECT() *MockOrderControllerInterfaceMockRecorder {
	return m.recorder
}

// CheckOrderIsExist mocks base method.
func (m *MockOrderControllerInterface) CheckOrderIsExist(orderID, userID, serviceID int) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckOrderIsExist", orderID, userID, serviceID)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckOrderIsExist indicates an expected call of CheckOrderIsExist.
func (mr *MockOrderControllerInterfaceMockRecorder) CheckOrderIsExist(orderID, userID, serviceID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckOrderIsExist", reflect.TypeOf((*MockOrderControllerInterface)(nil).CheckOrderIsExist), orderID, userID, serviceID)
}

// CreateNewOrder mocks base method.
func (m *MockOrderControllerInterface) CreateNewOrder(orderID, userID, serviceID int, sum float64, comment string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateNewOrder", orderID, userID, serviceID, sum, comment)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateNewOrder indicates an expected call of CreateNewOrder.
func (mr *MockOrderControllerInterfaceMockRecorder) CreateNewOrder(orderID, userID, serviceID, sum, comment interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateNewOrder", reflect.TypeOf((*MockOrderControllerInterface)(nil).CreateNewOrder), orderID, userID, serviceID, sum, comment)
}

// FinishOrder mocks base method.
func (m *MockOrderControllerInterface) FinishOrder(orderID, userID, serviceID int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FinishOrder", orderID, userID, serviceID)
	ret0, _ := ret[0].(error)
	return ret0
}

// FinishOrder indicates an expected call of FinishOrder.
func (mr *MockOrderControllerInterfaceMockRecorder) FinishOrder(orderID, userID, serviceID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FinishOrder", reflect.TypeOf((*MockOrderControllerInterface)(nil).FinishOrder), orderID, userID, serviceID)
}

// GetFinanceReports mocks base method.
func (m *MockOrderControllerInterface) GetFinanceReports(month, year int) ([]report.FinanceReport, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFinanceReports", month, year)
	ret0, _ := ret[0].([]report.FinanceReport)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFinanceReports indicates an expected call of GetFinanceReports.
func (mr *MockOrderControllerInterfaceMockRecorder) GetFinanceReports(month, year interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFinanceReports", reflect.TypeOf((*MockOrderControllerInterface)(nil).GetFinanceReports), month, year)
}

// GetOrder mocks base method.
func (m *MockOrderControllerInterface) GetOrder(orderID, userID, serviceID int) (order.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOrder", orderID, userID, serviceID)
	ret0, _ := ret[0].(order.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOrder indicates an expected call of GetOrder.
func (mr *MockOrderControllerInterfaceMockRecorder) GetOrder(orderID, userID, serviceID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrder", reflect.TypeOf((*MockOrderControllerInterface)(nil).GetOrder), orderID, userID, serviceID)
}

// ReserveOrder mocks base method.
func (m *MockOrderControllerInterface) ReserveOrder(orderID, userID, serviceID int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReserveOrder", orderID, userID, serviceID)
	ret0, _ := ret[0].(error)
	return ret0
}

// ReserveOrder indicates an expected call of ReserveOrder.
func (mr *MockOrderControllerInterfaceMockRecorder) ReserveOrder(orderID, userID, serviceID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReserveOrder", reflect.TypeOf((*MockOrderControllerInterface)(nil).ReserveOrder), orderID, userID, serviceID)
}

// ReturnOrder mocks base method.
func (m *MockOrderControllerInterface) ReturnOrder(orderID, userID, serviceID int) (float64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReturnOrder", orderID, userID, serviceID)
	ret0, _ := ret[0].(float64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReturnOrder indicates an expected call of ReturnOrder.
func (mr *MockOrderControllerInterfaceMockRecorder) ReturnOrder(orderID, userID, serviceID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReturnOrder", reflect.TypeOf((*MockOrderControllerInterface)(nil).ReturnOrder), orderID, userID, serviceID)
}
