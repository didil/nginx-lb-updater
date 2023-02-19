// Code generated by MockGen. DO NOT EDIT.
// Source: services/lb_updater.go

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	services "github.com/didil/nginx-lb-updater/services"
	gomock "github.com/golang/mock/gomock"
)

// MockLBUpdater is a mock of LBUpdater interface.
type MockLBUpdater struct {
	ctrl     *gomock.Controller
	recorder *MockLBUpdaterMockRecorder
}

// MockLBUpdaterMockRecorder is the mock recorder for MockLBUpdater.
type MockLBUpdaterMockRecorder struct {
	mock *MockLBUpdater
}

// NewMockLBUpdater creates a new mock instance.
func NewMockLBUpdater(ctrl *gomock.Controller) *MockLBUpdater {
	mock := &MockLBUpdater{ctrl: ctrl}
	mock.recorder = &MockLBUpdaterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLBUpdater) EXPECT() *MockLBUpdaterMockRecorder {
	return m.recorder
}

// UpdateStream mocks base method.
func (m *MockLBUpdater) UpdateStream(backendName string, port int, protocol string, servers []services.Server, proxyTimeoutSeconds, proxyConnectTimeoutSeconds int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateStream", backendName, port, protocol, servers, proxyTimeoutSeconds, proxyConnectTimeoutSeconds)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateStream indicates an expected call of UpdateStream.
func (mr *MockLBUpdaterMockRecorder) UpdateStream(backendName, port, protocol, servers, proxyTimeoutSeconds, proxyConnectTimeoutSeconds interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateStream", reflect.TypeOf((*MockLBUpdater)(nil).UpdateStream), backendName, port, protocol, servers, proxyTimeoutSeconds, proxyConnectTimeoutSeconds)
}