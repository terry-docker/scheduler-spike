// Code generated by MockGen. DO NOT EDIT.
// Source: persistence/persistence.go

// Package mock_persistence is a generated GoMock package.
package mock_persistence

import (
	reflect "reflect"

	common "example.com/scheduler/common"
	gomock "github.com/golang/mock/gomock"
)

// MockIPersistenceManager is a mock of IPersistenceManager interface.
type MockIPersistenceManager struct {
	ctrl     *gomock.Controller
	recorder *MockIPersistenceManagerMockRecorder
}

// MockIPersistenceManagerMockRecorder is the mock recorder for MockIPersistenceManager.
type MockIPersistenceManagerMockRecorder struct {
	mock *MockIPersistenceManager
}

// NewMockIPersistenceManager creates a new mock instance.
func NewMockIPersistenceManager(ctrl *gomock.Controller) *MockIPersistenceManager {
	mock := &MockIPersistenceManager{ctrl: ctrl}
	mock.recorder = &MockIPersistenceManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIPersistenceManager) EXPECT() *MockIPersistenceManagerMockRecorder {
	return m.recorder
}

// LoadTasks mocks base method.
func (m *MockIPersistenceManager) LoadTasks() ([]common.TaskConfig, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LoadTasks")
	ret0, _ := ret[0].([]common.TaskConfig)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LoadTasks indicates an expected call of LoadTasks.
func (mr *MockIPersistenceManagerMockRecorder) LoadTasks() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoadTasks", reflect.TypeOf((*MockIPersistenceManager)(nil).LoadTasks))
}

// SaveTasks mocks base method.
func (m *MockIPersistenceManager) SaveTasks(arg0 []common.TaskConfig) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveTasks", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveTasks indicates an expected call of SaveTasks.
func (mr *MockIPersistenceManagerMockRecorder) SaveTasks(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveTasks", reflect.TypeOf((*MockIPersistenceManager)(nil).SaveTasks), arg0)
}