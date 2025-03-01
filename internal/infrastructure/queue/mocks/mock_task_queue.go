// Code generated by mockery v2.52.2. DO NOT EDIT.

package mocks

import (
	context "context"

	queue "github.com/ardfard/sb-test/internal/domain/queue"
	mock "github.com/stretchr/testify/mock"
)

// MockTaskQueue is an autogenerated mock type for the TaskQueue type
type MockTaskQueue struct {
	mock.Mock
}

type MockTaskQueue_Expecter struct {
	mock *mock.Mock
}

func (_m *MockTaskQueue) EXPECT() *MockTaskQueue_Expecter {
	return &MockTaskQueue_Expecter{mock: &_m.Mock}
}

// Complete provides a mock function with given fields: ctx, taskID
func (_m *MockTaskQueue) Complete(ctx context.Context, taskID string) error {
	ret := _m.Called(ctx, taskID)

	if len(ret) == 0 {
		panic("no return value specified for Complete")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, taskID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockTaskQueue_Complete_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Complete'
type MockTaskQueue_Complete_Call struct {
	*mock.Call
}

// Complete is a helper method to define mock.On call
//   - ctx context.Context
//   - taskID string
func (_e *MockTaskQueue_Expecter) Complete(ctx interface{}, taskID interface{}) *MockTaskQueue_Complete_Call {
	return &MockTaskQueue_Complete_Call{Call: _e.mock.On("Complete", ctx, taskID)}
}

func (_c *MockTaskQueue_Complete_Call) Run(run func(ctx context.Context, taskID string)) *MockTaskQueue_Complete_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *MockTaskQueue_Complete_Call) Return(_a0 error) *MockTaskQueue_Complete_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockTaskQueue_Complete_Call) RunAndReturn(run func(context.Context, string) error) *MockTaskQueue_Complete_Call {
	_c.Call.Return(run)
	return _c
}

// Dequeue provides a mock function with given fields: ctx
func (_m *MockTaskQueue) Dequeue(ctx context.Context) (*queue.Task, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for Dequeue")
	}

	var r0 *queue.Task
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (*queue.Task, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) *queue.Task); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*queue.Task)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockTaskQueue_Dequeue_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Dequeue'
type MockTaskQueue_Dequeue_Call struct {
	*mock.Call
}

// Dequeue is a helper method to define mock.On call
//   - ctx context.Context
func (_e *MockTaskQueue_Expecter) Dequeue(ctx interface{}) *MockTaskQueue_Dequeue_Call {
	return &MockTaskQueue_Dequeue_Call{Call: _e.mock.On("Dequeue", ctx)}
}

func (_c *MockTaskQueue_Dequeue_Call) Run(run func(ctx context.Context)) *MockTaskQueue_Dequeue_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *MockTaskQueue_Dequeue_Call) Return(_a0 *queue.Task, _a1 error) *MockTaskQueue_Dequeue_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockTaskQueue_Dequeue_Call) RunAndReturn(run func(context.Context) (*queue.Task, error)) *MockTaskQueue_Dequeue_Call {
	_c.Call.Return(run)
	return _c
}

// Enqueue provides a mock function with given fields: ctx, payload
func (_m *MockTaskQueue) Enqueue(ctx context.Context, payload uint) error {
	ret := _m.Called(ctx, payload)

	if len(ret) == 0 {
		panic("no return value specified for Enqueue")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uint) error); ok {
		r0 = rf(ctx, payload)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockTaskQueue_Enqueue_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Enqueue'
type MockTaskQueue_Enqueue_Call struct {
	*mock.Call
}

// Enqueue is a helper method to define mock.On call
//   - ctx context.Context
//   - payload uint
func (_e *MockTaskQueue_Expecter) Enqueue(ctx interface{}, payload interface{}) *MockTaskQueue_Enqueue_Call {
	return &MockTaskQueue_Enqueue_Call{Call: _e.mock.On("Enqueue", ctx, payload)}
}

func (_c *MockTaskQueue_Enqueue_Call) Run(run func(ctx context.Context, payload uint)) *MockTaskQueue_Enqueue_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uint))
	})
	return _c
}

func (_c *MockTaskQueue_Enqueue_Call) Return(_a0 error) *MockTaskQueue_Enqueue_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockTaskQueue_Enqueue_Call) RunAndReturn(run func(context.Context, uint) error) *MockTaskQueue_Enqueue_Call {
	_c.Call.Return(run)
	return _c
}

// Fail provides a mock function with given fields: ctx, taskID, errMsg
func (_m *MockTaskQueue) Fail(ctx context.Context, taskID string, errMsg string) error {
	ret := _m.Called(ctx, taskID, errMsg)

	if len(ret) == 0 {
		panic("no return value specified for Fail")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) error); ok {
		r0 = rf(ctx, taskID, errMsg)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockTaskQueue_Fail_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Fail'
type MockTaskQueue_Fail_Call struct {
	*mock.Call
}

// Fail is a helper method to define mock.On call
//   - ctx context.Context
//   - taskID string
//   - errMsg string
func (_e *MockTaskQueue_Expecter) Fail(ctx interface{}, taskID interface{}, errMsg interface{}) *MockTaskQueue_Fail_Call {
	return &MockTaskQueue_Fail_Call{Call: _e.mock.On("Fail", ctx, taskID, errMsg)}
}

func (_c *MockTaskQueue_Fail_Call) Run(run func(ctx context.Context, taskID string, errMsg string)) *MockTaskQueue_Fail_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *MockTaskQueue_Fail_Call) Return(_a0 error) *MockTaskQueue_Fail_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockTaskQueue_Fail_Call) RunAndReturn(run func(context.Context, string, string) error) *MockTaskQueue_Fail_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockTaskQueue creates a new instance of MockTaskQueue. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockTaskQueue(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockTaskQueue {
	mock := &MockTaskQueue{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
