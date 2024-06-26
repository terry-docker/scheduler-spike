// Code generated by MockGen. DO NOT EDIT.
// Source: cronjobs/cron.go

// Package mock_cronjobs is a generated GoMock package.
package mock_cronjobs

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	cron "github.com/robfig/cron/v3"
)

// MockICronScheduler is a mock of ICronScheduler interface.
type MockICronScheduler struct {
	ctrl     *gomock.Controller
	recorder *MockICronSchedulerMockRecorder
}

// MockICronSchedulerMockRecorder is the mock recorder for MockICronScheduler.
type MockICronSchedulerMockRecorder struct {
	mock *MockICronScheduler
}

// NewMockICronScheduler creates a new mock instance.
func NewMockICronScheduler(ctrl *gomock.Controller) *MockICronScheduler {
	mock := &MockICronScheduler{ctrl: ctrl}
	mock.recorder = &MockICronSchedulerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockICronScheduler) EXPECT() *MockICronSchedulerMockRecorder {
	return m.recorder
}

// AddTask mocks base method.
func (m *MockICronScheduler) AddTask(spec string, task func()) (cron.EntryID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddTask", spec, task)
	ret0, _ := ret[0].(cron.EntryID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddTask indicates an expected call of AddTask.
func (mr *MockICronSchedulerMockRecorder) AddTask(spec, task interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddTask", reflect.TypeOf((*MockICronScheduler)(nil).AddTask), spec, task)
}

// List mocks base method.
func (m *MockICronScheduler) List() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "List")
}

// List indicates an expected call of List.
func (mr *MockICronSchedulerMockRecorder) List() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockICronScheduler)(nil).List))
}

// RemoveTask mocks base method.
func (m *MockICronScheduler) RemoveTask(id cron.EntryID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveTask", id)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveTask indicates an expected call of RemoveTask.
func (mr *MockICronSchedulerMockRecorder) RemoveTask(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveTask", reflect.TypeOf((*MockICronScheduler)(nil).RemoveTask), id)
}

// Start mocks base method.
func (m *MockICronScheduler) Start() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Start")
}

// Start indicates an expected call of Start.
func (mr *MockICronSchedulerMockRecorder) Start() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Start", reflect.TypeOf((*MockICronScheduler)(nil).Start))
}

// Stop mocks base method.
func (m *MockICronScheduler) Stop() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Stop")
}

// Stop indicates an expected call of Stop.
func (mr *MockICronSchedulerMockRecorder) Stop() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stop", reflect.TypeOf((*MockICronScheduler)(nil).Stop))
}
