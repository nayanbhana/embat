// Code generated by MockGen. DO NOT EDIT.
// Source: embat.go
//
// Generated by this command:
//
//	mockgen -source=embat.go -destination=./mock/embat.go -package=mock
//

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	embat "github.com/nayanbhana/embat"
	gomock "go.uber.org/mock/gomock"
)

// MockBatchProcessor is a mock of BatchProcessor interface.
type MockBatchProcessor[J any, R any] struct {
	ctrl     *gomock.Controller
	recorder *MockBatchProcessorMockRecorder[J, R]
}

// MockBatchProcessorMockRecorder is the mock recorder for MockBatchProcessor.
type MockBatchProcessorMockRecorder[J any, R any] struct {
	mock *MockBatchProcessor[J, R]
}

// NewMockBatchProcessor creates a new mock instance.
func NewMockBatchProcessor[J any, R any](ctrl *gomock.Controller) *MockBatchProcessor[J, R] {
	mock := &MockBatchProcessor[J, R]{ctrl: ctrl}
	mock.recorder = &MockBatchProcessorMockRecorder[J, R]{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBatchProcessor[J, R]) EXPECT() *MockBatchProcessorMockRecorder[J, R] {
	return m.recorder
}

// Process mocks base method.
func (m *MockBatchProcessor[J, R]) Process(batch []embat.Job[J]) []embat.Result[R] {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Process", batch)
	ret0, _ := ret[0].([]embat.Result[R])
	return ret0
}

// Process indicates an expected call of Process.
func (mr *MockBatchProcessorMockRecorder[J, R]) Process(batch any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Process", reflect.TypeOf((*MockBatchProcessor[J, R])(nil).Process), batch)
}
