// Code generated by mockery v2.43.2. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// ChatGptApi is an autogenerated mock type for the ChatGptApi type
type ChatGptApi struct {
	mock.Mock
}

// Call provides a mock function with given fields: request
func (_m *ChatGptApi) Call(request string) (string, error) {
	ret := _m.Called(request)

	if len(ret) == 0 {
		panic("no return value specified for Call")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (string, error)); ok {
		return rf(request)
	}
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(request)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(request)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewChatGptApi creates a new instance of ChatGptApi. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewChatGptApi(t interface {
	mock.TestingT
	Cleanup(func())
}) *ChatGptApi {
	mock := &ChatGptApi{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
