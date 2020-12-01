// Code generated by MockGen. DO NOT EDIT.
// Source: ./event_handlers.go

// Package mock_controller is a generated GoMock package.
package mock_controller

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	v1alpha1 "github.com/solo-io/gloo-mesh/pkg/api/xds.enterprise.agent.mesh.gloo.solo.io/v1alpha1"
	controller "github.com/solo-io/gloo-mesh/pkg/api/xds.enterprise.agent.mesh.gloo.solo.io/v1alpha1/controller"
	predicate "sigs.k8s.io/controller-runtime/pkg/predicate"
)

// MockXdsConfigEventHandler is a mock of XdsConfigEventHandler interface
type MockXdsConfigEventHandler struct {
	ctrl     *gomock.Controller
	recorder *MockXdsConfigEventHandlerMockRecorder
}

// MockXdsConfigEventHandlerMockRecorder is the mock recorder for MockXdsConfigEventHandler
type MockXdsConfigEventHandlerMockRecorder struct {
	mock *MockXdsConfigEventHandler
}

// NewMockXdsConfigEventHandler creates a new mock instance
func NewMockXdsConfigEventHandler(ctrl *gomock.Controller) *MockXdsConfigEventHandler {
	mock := &MockXdsConfigEventHandler{ctrl: ctrl}
	mock.recorder = &MockXdsConfigEventHandlerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockXdsConfigEventHandler) EXPECT() *MockXdsConfigEventHandlerMockRecorder {
	return m.recorder
}

// CreateXdsConfig mocks base method
func (m *MockXdsConfigEventHandler) CreateXdsConfig(obj *v1alpha1.XdsConfig) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateXdsConfig", obj)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateXdsConfig indicates an expected call of CreateXdsConfig
func (mr *MockXdsConfigEventHandlerMockRecorder) CreateXdsConfig(obj interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateXdsConfig", reflect.TypeOf((*MockXdsConfigEventHandler)(nil).CreateXdsConfig), obj)
}

// UpdateXdsConfig mocks base method
func (m *MockXdsConfigEventHandler) UpdateXdsConfig(old, new *v1alpha1.XdsConfig) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateXdsConfig", old, new)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateXdsConfig indicates an expected call of UpdateXdsConfig
func (mr *MockXdsConfigEventHandlerMockRecorder) UpdateXdsConfig(old, new interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateXdsConfig", reflect.TypeOf((*MockXdsConfigEventHandler)(nil).UpdateXdsConfig), old, new)
}

// DeleteXdsConfig mocks base method
func (m *MockXdsConfigEventHandler) DeleteXdsConfig(obj *v1alpha1.XdsConfig) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteXdsConfig", obj)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteXdsConfig indicates an expected call of DeleteXdsConfig
func (mr *MockXdsConfigEventHandlerMockRecorder) DeleteXdsConfig(obj interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteXdsConfig", reflect.TypeOf((*MockXdsConfigEventHandler)(nil).DeleteXdsConfig), obj)
}

// GenericXdsConfig mocks base method
func (m *MockXdsConfigEventHandler) GenericXdsConfig(obj *v1alpha1.XdsConfig) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenericXdsConfig", obj)
	ret0, _ := ret[0].(error)
	return ret0
}

// GenericXdsConfig indicates an expected call of GenericXdsConfig
func (mr *MockXdsConfigEventHandlerMockRecorder) GenericXdsConfig(obj interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenericXdsConfig", reflect.TypeOf((*MockXdsConfigEventHandler)(nil).GenericXdsConfig), obj)
}

// MockXdsConfigEventWatcher is a mock of XdsConfigEventWatcher interface
type MockXdsConfigEventWatcher struct {
	ctrl     *gomock.Controller
	recorder *MockXdsConfigEventWatcherMockRecorder
}

// MockXdsConfigEventWatcherMockRecorder is the mock recorder for MockXdsConfigEventWatcher
type MockXdsConfigEventWatcherMockRecorder struct {
	mock *MockXdsConfigEventWatcher
}

// NewMockXdsConfigEventWatcher creates a new mock instance
func NewMockXdsConfigEventWatcher(ctrl *gomock.Controller) *MockXdsConfigEventWatcher {
	mock := &MockXdsConfigEventWatcher{ctrl: ctrl}
	mock.recorder = &MockXdsConfigEventWatcherMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockXdsConfigEventWatcher) EXPECT() *MockXdsConfigEventWatcherMockRecorder {
	return m.recorder
}

// AddEventHandler mocks base method
func (m *MockXdsConfigEventWatcher) AddEventHandler(ctx context.Context, h controller.XdsConfigEventHandler, predicates ...predicate.Predicate) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, h}
	for _, a := range predicates {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "AddEventHandler", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddEventHandler indicates an expected call of AddEventHandler
func (mr *MockXdsConfigEventWatcherMockRecorder) AddEventHandler(ctx, h interface{}, predicates ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, h}, predicates...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddEventHandler", reflect.TypeOf((*MockXdsConfigEventWatcher)(nil).AddEventHandler), varargs...)
}
