// Code generated by mockery v2.52.2. DO NOT EDIT.

package mocks

import (
	context "context"

	entity "github.com/ardfard/sb-test/internal/domain/entity"
	mock "github.com/stretchr/testify/mock"
)

// MockAudioRepository is an autogenerated mock type for the AudioRepository type
type MockAudioRepository struct {
	mock.Mock
}

type MockAudioRepository_Expecter struct {
	mock *mock.Mock
}

func (_m *MockAudioRepository) EXPECT() *MockAudioRepository_Expecter {
	return &MockAudioRepository_Expecter{mock: &_m.Mock}
}

// GetByID provides a mock function with given fields: ctx, id
func (_m *MockAudioRepository) GetByID(ctx context.Context, id uint) (*entity.Audio, error) {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for GetByID")
	}

	var r0 *entity.Audio
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint) (*entity.Audio, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint) *entity.Audio); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.Audio)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockAudioRepository_GetByID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetByID'
type MockAudioRepository_GetByID_Call struct {
	*mock.Call
}

// GetByID is a helper method to define mock.On call
//   - ctx context.Context
//   - id uint
func (_e *MockAudioRepository_Expecter) GetByID(ctx interface{}, id interface{}) *MockAudioRepository_GetByID_Call {
	return &MockAudioRepository_GetByID_Call{Call: _e.mock.On("GetByID", ctx, id)}
}

func (_c *MockAudioRepository_GetByID_Call) Run(run func(ctx context.Context, id uint)) *MockAudioRepository_GetByID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uint))
	})
	return _c
}

func (_c *MockAudioRepository_GetByID_Call) Return(_a0 *entity.Audio, _a1 error) *MockAudioRepository_GetByID_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockAudioRepository_GetByID_Call) RunAndReturn(run func(context.Context, uint) (*entity.Audio, error)) *MockAudioRepository_GetByID_Call {
	_c.Call.Return(run)
	return _c
}

// GetByUserIDAndPhraseID provides a mock function with given fields: ctx, userID, phraseID
func (_m *MockAudioRepository) GetByUserIDAndPhraseID(ctx context.Context, userID uint, phraseID uint) (*entity.Audio, error) {
	ret := _m.Called(ctx, userID, phraseID)

	if len(ret) == 0 {
		panic("no return value specified for GetByUserIDAndPhraseID")
	}

	var r0 *entity.Audio
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint, uint) (*entity.Audio, error)); ok {
		return rf(ctx, userID, phraseID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint, uint) *entity.Audio); ok {
		r0 = rf(ctx, userID, phraseID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.Audio)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint, uint) error); ok {
		r1 = rf(ctx, userID, phraseID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockAudioRepository_GetByUserIDAndPhraseID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetByUserIDAndPhraseID'
type MockAudioRepository_GetByUserIDAndPhraseID_Call struct {
	*mock.Call
}

// GetByUserIDAndPhraseID is a helper method to define mock.On call
//   - ctx context.Context
//   - userID uint
//   - phraseID uint
func (_e *MockAudioRepository_Expecter) GetByUserIDAndPhraseID(ctx interface{}, userID interface{}, phraseID interface{}) *MockAudioRepository_GetByUserIDAndPhraseID_Call {
	return &MockAudioRepository_GetByUserIDAndPhraseID_Call{Call: _e.mock.On("GetByUserIDAndPhraseID", ctx, userID, phraseID)}
}

func (_c *MockAudioRepository_GetByUserIDAndPhraseID_Call) Run(run func(ctx context.Context, userID uint, phraseID uint)) *MockAudioRepository_GetByUserIDAndPhraseID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uint), args[2].(uint))
	})
	return _c
}

func (_c *MockAudioRepository_GetByUserIDAndPhraseID_Call) Return(_a0 *entity.Audio, _a1 error) *MockAudioRepository_GetByUserIDAndPhraseID_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockAudioRepository_GetByUserIDAndPhraseID_Call) RunAndReturn(run func(context.Context, uint, uint) (*entity.Audio, error)) *MockAudioRepository_GetByUserIDAndPhraseID_Call {
	_c.Call.Return(run)
	return _c
}

// Store provides a mock function with given fields: ctx, audio
func (_m *MockAudioRepository) Store(ctx context.Context, audio *entity.Audio) (*entity.Audio, error) {
	ret := _m.Called(ctx, audio)

	if len(ret) == 0 {
		panic("no return value specified for Store")
	}

	var r0 *entity.Audio
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *entity.Audio) (*entity.Audio, error)); ok {
		return rf(ctx, audio)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *entity.Audio) *entity.Audio); ok {
		r0 = rf(ctx, audio)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.Audio)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *entity.Audio) error); ok {
		r1 = rf(ctx, audio)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockAudioRepository_Store_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Store'
type MockAudioRepository_Store_Call struct {
	*mock.Call
}

// Store is a helper method to define mock.On call
//   - ctx context.Context
//   - audio *entity.Audio
func (_e *MockAudioRepository_Expecter) Store(ctx interface{}, audio interface{}) *MockAudioRepository_Store_Call {
	return &MockAudioRepository_Store_Call{Call: _e.mock.On("Store", ctx, audio)}
}

func (_c *MockAudioRepository_Store_Call) Run(run func(ctx context.Context, audio *entity.Audio)) *MockAudioRepository_Store_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*entity.Audio))
	})
	return _c
}

func (_c *MockAudioRepository_Store_Call) Return(_a0 *entity.Audio, _a1 error) *MockAudioRepository_Store_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockAudioRepository_Store_Call) RunAndReturn(run func(context.Context, *entity.Audio) (*entity.Audio, error)) *MockAudioRepository_Store_Call {
	_c.Call.Return(run)
	return _c
}

// Update provides a mock function with given fields: ctx, audio
func (_m *MockAudioRepository) Update(ctx context.Context, audio *entity.Audio) error {
	ret := _m.Called(ctx, audio)

	if len(ret) == 0 {
		panic("no return value specified for Update")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *entity.Audio) error); ok {
		r0 = rf(ctx, audio)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockAudioRepository_Update_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Update'
type MockAudioRepository_Update_Call struct {
	*mock.Call
}

// Update is a helper method to define mock.On call
//   - ctx context.Context
//   - audio *entity.Audio
func (_e *MockAudioRepository_Expecter) Update(ctx interface{}, audio interface{}) *MockAudioRepository_Update_Call {
	return &MockAudioRepository_Update_Call{Call: _e.mock.On("Update", ctx, audio)}
}

func (_c *MockAudioRepository_Update_Call) Run(run func(ctx context.Context, audio *entity.Audio)) *MockAudioRepository_Update_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*entity.Audio))
	})
	return _c
}

func (_c *MockAudioRepository_Update_Call) Return(_a0 error) *MockAudioRepository_Update_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockAudioRepository_Update_Call) RunAndReturn(run func(context.Context, *entity.Audio) error) *MockAudioRepository_Update_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockAudioRepository creates a new instance of MockAudioRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockAudioRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockAudioRepository {
	mock := &MockAudioRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
