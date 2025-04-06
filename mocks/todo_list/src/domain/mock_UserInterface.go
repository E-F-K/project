// Code generated by mockery. DO NOT EDIT.

package domain

import (
	context "context"
	domain "todo_list/internal/domain"

	mock "github.com/stretchr/testify/mock"
)

// MockUserInterface is an autogenerated mock type for the UserInterface type
type MockUserInterface struct {
	mock.Mock
}

type MockUserInterface_Expecter struct {
	mock *mock.Mock
}

func (_m *MockUserInterface) EXPECT() *MockUserInterface_Expecter {
	return &MockUserInterface_Expecter{mock: &_m.Mock}
}

// Authenticate provides a mock function with given fields: ctx, token
func (_m *MockUserInterface) Authenticate(ctx context.Context, token string) (domain.User, error) {
	ret := _m.Called(ctx, token)

	if len(ret) == 0 {
		panic("no return value specified for Authenticate")
	}

	var r0 domain.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (domain.User, error)); ok {
		return rf(ctx, token)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) domain.User); ok {
		r0 = rf(ctx, token)
	} else {
		r0 = ret.Get(0).(domain.User)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, token)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockUserInterface_Authenticate_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Authenticate'
type MockUserInterface_Authenticate_Call struct {
	*mock.Call
}

// Authenticate is a helper method to define mock.On call
//   - ctx context.Context
//   - token string
func (_e *MockUserInterface_Expecter) Authenticate(ctx interface{}, token interface{}) *MockUserInterface_Authenticate_Call {
	return &MockUserInterface_Authenticate_Call{Call: _e.mock.On("Authenticate", ctx, token)}
}

func (_c *MockUserInterface_Authenticate_Call) Run(run func(ctx context.Context, token string)) *MockUserInterface_Authenticate_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *MockUserInterface_Authenticate_Call) Return(_a0 domain.User, _a1 error) *MockUserInterface_Authenticate_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockUserInterface_Authenticate_Call) RunAndReturn(run func(context.Context, string) (domain.User, error)) *MockUserInterface_Authenticate_Call {
	_c.Call.Return(run)
	return _c
}

// Close provides a mock function with no fields
func (_m *MockUserInterface) Close() error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Close")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockUserInterface_Close_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Close'
type MockUserInterface_Close_Call struct {
	*mock.Call
}

// Close is a helper method to define mock.On call
func (_e *MockUserInterface_Expecter) Close() *MockUserInterface_Close_Call {
	return &MockUserInterface_Close_Call{Call: _e.mock.On("Close")}
}

func (_c *MockUserInterface_Close_Call) Run(run func()) *MockUserInterface_Close_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockUserInterface_Close_Call) Return(_a0 error) *MockUserInterface_Close_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockUserInterface_Close_Call) RunAndReturn(run func() error) *MockUserInterface_Close_Call {
	_c.Call.Return(run)
	return _c
}

// Login provides a mock function with given fields: ctx, email, password
func (_m *MockUserInterface) Login(ctx context.Context, email string, password string) error {
	ret := _m.Called(ctx, email, password)

	if len(ret) == 0 {
		panic("no return value specified for Login")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) error); ok {
		r0 = rf(ctx, email, password)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockUserInterface_Login_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Login'
type MockUserInterface_Login_Call struct {
	*mock.Call
}

// Login is a helper method to define mock.On call
//   - ctx context.Context
//   - email string
//   - password string
func (_e *MockUserInterface_Expecter) Login(ctx interface{}, email interface{}, password interface{}) *MockUserInterface_Login_Call {
	return &MockUserInterface_Login_Call{Call: _e.mock.On("Login", ctx, email, password)}
}

func (_c *MockUserInterface_Login_Call) Run(run func(ctx context.Context, email string, password string)) *MockUserInterface_Login_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *MockUserInterface_Login_Call) Return(_a0 error) *MockUserInterface_Login_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockUserInterface_Login_Call) RunAndReturn(run func(context.Context, string, string) error) *MockUserInterface_Login_Call {
	_c.Call.Return(run)
	return _c
}

// RegisterUser provides a mock function with given fields: ctx, name, email, passwordHash, token
func (_m *MockUserInterface) RegisterUser(ctx context.Context, name string, email string, passwordHash string, token string) error {
	ret := _m.Called(ctx, name, email, passwordHash, token)

	if len(ret) == 0 {
		panic("no return value specified for RegisterUser")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string, string) error); ok {
		r0 = rf(ctx, name, email, passwordHash, token)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockUserInterface_RegisterUser_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RegisterUser'
type MockUserInterface_RegisterUser_Call struct {
	*mock.Call
}

// RegisterUser is a helper method to define mock.On call
//   - ctx context.Context
//   - name string
//   - email string
//   - passwordHash string
//   - token string
func (_e *MockUserInterface_Expecter) RegisterUser(ctx interface{}, name interface{}, email interface{}, passwordHash interface{}, token interface{}) *MockUserInterface_RegisterUser_Call {
	return &MockUserInterface_RegisterUser_Call{Call: _e.mock.On("RegisterUser", ctx, name, email, passwordHash, token)}
}

func (_c *MockUserInterface_RegisterUser_Call) Run(run func(ctx context.Context, name string, email string, passwordHash string, token string)) *MockUserInterface_RegisterUser_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string), args[3].(string), args[4].(string))
	})
	return _c
}

func (_c *MockUserInterface_RegisterUser_Call) Return(_a0 error) *MockUserInterface_RegisterUser_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockUserInterface_RegisterUser_Call) RunAndReturn(run func(context.Context, string, string, string, string) error) *MockUserInterface_RegisterUser_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateToken provides a mock function with given fields: ctx, email, token
func (_m *MockUserInterface) UpdateToken(ctx context.Context, email string, token string) error {
	ret := _m.Called(ctx, email, token)

	if len(ret) == 0 {
		panic("no return value specified for UpdateToken")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) error); ok {
		r0 = rf(ctx, email, token)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockUserInterface_UpdateToken_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateToken'
type MockUserInterface_UpdateToken_Call struct {
	*mock.Call
}

// UpdateToken is a helper method to define mock.On call
//   - ctx context.Context
//   - email string
//   - token string
func (_e *MockUserInterface_Expecter) UpdateToken(ctx interface{}, email interface{}, token interface{}) *MockUserInterface_UpdateToken_Call {
	return &MockUserInterface_UpdateToken_Call{Call: _e.mock.On("UpdateToken", ctx, email, token)}
}

func (_c *MockUserInterface_UpdateToken_Call) Run(run func(ctx context.Context, email string, token string)) *MockUserInterface_UpdateToken_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *MockUserInterface_UpdateToken_Call) Return(_a0 error) *MockUserInterface_UpdateToken_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockUserInterface_UpdateToken_Call) RunAndReturn(run func(context.Context, string, string) error) *MockUserInterface_UpdateToken_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockUserInterface creates a new instance of MockUserInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockUserInterface(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockUserInterface {
	mock := &MockUserInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
