// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	context "context"

	models "github.com/aetheris-lab/aetheris-id/api/internal/models"
	mock "github.com/stretchr/testify/mock"
)

// OAuthServiceMock is an autogenerated mock type for the OAuthService type
type OAuthServiceMock struct {
	mock.Mock
}

type OAuthServiceMock_Expecter struct {
	mock *mock.Mock
}

func (_m *OAuthServiceMock) EXPECT() *OAuthServiceMock_Expecter {
	return &OAuthServiceMock_Expecter{mock: &_m.Mock}
}

// Authorize provides a mock function with given fields: ctx, input
func (_m *OAuthServiceMock) Authorize(ctx context.Context, input models.AuthorizeInput) (*models.AuthorizeResponse, error) {
	ret := _m.Called(ctx, input)

	if len(ret) == 0 {
		panic("no return value specified for Authorize")
	}

	var r0 *models.AuthorizeResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, models.AuthorizeInput) (*models.AuthorizeResponse, error)); ok {
		return rf(ctx, input)
	}
	if rf, ok := ret.Get(0).(func(context.Context, models.AuthorizeInput) *models.AuthorizeResponse); ok {
		r0 = rf(ctx, input)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.AuthorizeResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, models.AuthorizeInput) error); ok {
		r1 = rf(ctx, input)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// OAuthServiceMock_Authorize_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Authorize'
type OAuthServiceMock_Authorize_Call struct {
	*mock.Call
}

// Authorize is a helper method to define mock.On call
//   - ctx context.Context
//   - input models.AuthorizeInput
func (_e *OAuthServiceMock_Expecter) Authorize(ctx interface{}, input interface{}) *OAuthServiceMock_Authorize_Call {
	return &OAuthServiceMock_Authorize_Call{Call: _e.mock.On("Authorize", ctx, input)}
}

func (_c *OAuthServiceMock_Authorize_Call) Run(run func(ctx context.Context, input models.AuthorizeInput)) *OAuthServiceMock_Authorize_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(models.AuthorizeInput))
	})
	return _c
}

func (_c *OAuthServiceMock_Authorize_Call) Return(_a0 *models.AuthorizeResponse, _a1 error) *OAuthServiceMock_Authorize_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *OAuthServiceMock_Authorize_Call) RunAndReturn(run func(context.Context, models.AuthorizeInput) (*models.AuthorizeResponse, error)) *OAuthServiceMock_Authorize_Call {
	_c.Call.Return(run)
	return _c
}

// ExchangeCodeForToken provides a mock function with given fields: ctx, input
func (_m *OAuthServiceMock) ExchangeCodeForToken(ctx context.Context, input models.ExchangeAuthorizationCodeInput) (*models.TokenResponse, error) {
	ret := _m.Called(ctx, input)

	if len(ret) == 0 {
		panic("no return value specified for ExchangeCodeForToken")
	}

	var r0 *models.TokenResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, models.ExchangeAuthorizationCodeInput) (*models.TokenResponse, error)); ok {
		return rf(ctx, input)
	}
	if rf, ok := ret.Get(0).(func(context.Context, models.ExchangeAuthorizationCodeInput) *models.TokenResponse); ok {
		r0 = rf(ctx, input)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.TokenResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, models.ExchangeAuthorizationCodeInput) error); ok {
		r1 = rf(ctx, input)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// OAuthServiceMock_ExchangeCodeForToken_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ExchangeCodeForToken'
type OAuthServiceMock_ExchangeCodeForToken_Call struct {
	*mock.Call
}

// ExchangeCodeForToken is a helper method to define mock.On call
//   - ctx context.Context
//   - input models.ExchangeAuthorizationCodeInput
func (_e *OAuthServiceMock_Expecter) ExchangeCodeForToken(ctx interface{}, input interface{}) *OAuthServiceMock_ExchangeCodeForToken_Call {
	return &OAuthServiceMock_ExchangeCodeForToken_Call{Call: _e.mock.On("ExchangeCodeForToken", ctx, input)}
}

func (_c *OAuthServiceMock_ExchangeCodeForToken_Call) Run(run func(ctx context.Context, input models.ExchangeAuthorizationCodeInput)) *OAuthServiceMock_ExchangeCodeForToken_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(models.ExchangeAuthorizationCodeInput))
	})
	return _c
}

func (_c *OAuthServiceMock_ExchangeCodeForToken_Call) Return(_a0 *models.TokenResponse, _a1 error) *OAuthServiceMock_ExchangeCodeForToken_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *OAuthServiceMock_ExchangeCodeForToken_Call) RunAndReturn(run func(context.Context, models.ExchangeAuthorizationCodeInput) (*models.TokenResponse, error)) *OAuthServiceMock_ExchangeCodeForToken_Call {
	_c.Call.Return(run)
	return _c
}

// NewOAuthServiceMock creates a new instance of OAuthServiceMock. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewOAuthServiceMock(t interface {
	mock.TestingT
	Cleanup(func())
}) *OAuthServiceMock {
	mock := &OAuthServiceMock{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
