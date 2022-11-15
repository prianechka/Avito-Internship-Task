// Code generated by MockGen. DO NOT EDIT.
// Source: interface.go

// Package mocks is a generated GoMock package.
package mocks

import (
	http "net/http"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockReportHandlerInterface is a mock of ReportHandlerInterface interface.
type MockReportHandlerInterface struct {
	ctrl     *gomock.Controller
	recorder *MockReportHandlerInterfaceMockRecorder
}

// MockReportHandlerInterfaceMockRecorder is the mock recorder for MockReportHandlerInterface.
type MockReportHandlerInterfaceMockRecorder struct {
	mock *MockReportHandlerInterface
}

// NewMockReportHandlerInterface creates a new mock instance.
func NewMockReportHandlerInterface(ctrl *gomock.Controller) *MockReportHandlerInterface {
	mock := &MockReportHandlerInterface{ctrl: ctrl}
	mock.recorder = &MockReportHandlerInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockReportHandlerInterface) EXPECT() *MockReportHandlerInterfaceMockRecorder {
	return m.recorder
}

// GetFinanceReport mocks base method.
func (m *MockReportHandlerInterface) GetFinanceReport(w http.ResponseWriter, r *http.Request) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "GetFinanceReport", w, r)
}

// GetFinanceReport indicates an expected call of GetFinanceReport.
func (mr *MockReportHandlerInterfaceMockRecorder) GetFinanceReport(w, r interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFinanceReport", reflect.TypeOf((*MockReportHandlerInterface)(nil).GetFinanceReport), w, r)
}