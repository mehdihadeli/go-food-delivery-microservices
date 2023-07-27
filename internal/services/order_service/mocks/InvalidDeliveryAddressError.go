// Code generated by mockery v2.30.16. DO NOT EDIT.

package mocks

import (
	fmt "fmt"

	mock "github.com/stretchr/testify/mock"
)

// InvalidDeliveryAddressError is an autogenerated mock type for the InvalidDeliveryAddressError type
type InvalidDeliveryAddressError struct {
	mock.Mock
}

type InvalidDeliveryAddressError_Expecter struct {
	mock *mock.Mock
}

func (_m *InvalidDeliveryAddressError) EXPECT() *InvalidDeliveryAddressError_Expecter {
	return &InvalidDeliveryAddressError_Expecter{mock: &_m.Mock}
}

// Cause provides a mock function with given fields:
func (_m *InvalidDeliveryAddressError) Cause() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// InvalidDeliveryAddressError_Cause_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Cause'
type InvalidDeliveryAddressError_Cause_Call struct {
	*mock.Call
}

// Cause is a helper method to define mock.On call
func (_e *InvalidDeliveryAddressError_Expecter) Cause() *InvalidDeliveryAddressError_Cause_Call {
	return &InvalidDeliveryAddressError_Cause_Call{Call: _e.mock.On("Cause")}
}

func (_c *InvalidDeliveryAddressError_Cause_Call) Run(run func()) *InvalidDeliveryAddressError_Cause_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *InvalidDeliveryAddressError_Cause_Call) Return(_a0 error) *InvalidDeliveryAddressError_Cause_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *InvalidDeliveryAddressError_Cause_Call) RunAndReturn(run func() error) *InvalidDeliveryAddressError_Cause_Call {
	_c.Call.Return(run)
	return _c
}

// Error provides a mock function with given fields:
func (_m *InvalidDeliveryAddressError) Error() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// InvalidDeliveryAddressError_Error_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Error'
type InvalidDeliveryAddressError_Error_Call struct {
	*mock.Call
}

// Error is a helper method to define mock.On call
func (_e *InvalidDeliveryAddressError_Expecter) Error() *InvalidDeliveryAddressError_Error_Call {
	return &InvalidDeliveryAddressError_Error_Call{Call: _e.mock.On("Error")}
}

func (_c *InvalidDeliveryAddressError_Error_Call) Run(run func()) *InvalidDeliveryAddressError_Error_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *InvalidDeliveryAddressError_Error_Call) Return(_a0 string) *InvalidDeliveryAddressError_Error_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *InvalidDeliveryAddressError_Error_Call) RunAndReturn(run func() string) *InvalidDeliveryAddressError_Error_Call {
	_c.Call.Return(run)
	return _c
}

// Format provides a mock function with given fields: f, verb
func (_m *InvalidDeliveryAddressError) Format(f fmt.State, verb rune) {
	_m.Called(f, verb)
}

// InvalidDeliveryAddressError_Format_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Format'
type InvalidDeliveryAddressError_Format_Call struct {
	*mock.Call
}

// Format is a helper method to define mock.On call
//   - f fmt.State
//   - verb rune
func (_e *InvalidDeliveryAddressError_Expecter) Format(f interface{}, verb interface{}) *InvalidDeliveryAddressError_Format_Call {
	return &InvalidDeliveryAddressError_Format_Call{Call: _e.mock.On("Format", f, verb)}
}

func (_c *InvalidDeliveryAddressError_Format_Call) Run(run func(f fmt.State, verb rune)) *InvalidDeliveryAddressError_Format_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(fmt.State), args[1].(rune))
	})
	return _c
}

func (_c *InvalidDeliveryAddressError_Format_Call) Return() *InvalidDeliveryAddressError_Format_Call {
	_c.Call.Return()
	return _c
}

func (_c *InvalidDeliveryAddressError_Format_Call) RunAndReturn(run func(fmt.State, rune)) *InvalidDeliveryAddressError_Format_Call {
	_c.Call.Return(run)
	return _c
}

// IsBadRequestError provides a mock function with given fields:
func (_m *InvalidDeliveryAddressError) IsBadRequestError() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// InvalidDeliveryAddressError_IsBadRequestError_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'IsBadRequestError'
type InvalidDeliveryAddressError_IsBadRequestError_Call struct {
	*mock.Call
}

// IsBadRequestError is a helper method to define mock.On call
func (_e *InvalidDeliveryAddressError_Expecter) IsBadRequestError() *InvalidDeliveryAddressError_IsBadRequestError_Call {
	return &InvalidDeliveryAddressError_IsBadRequestError_Call{Call: _e.mock.On("IsBadRequestError")}
}

func (_c *InvalidDeliveryAddressError_IsBadRequestError_Call) Run(run func()) *InvalidDeliveryAddressError_IsBadRequestError_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *InvalidDeliveryAddressError_IsBadRequestError_Call) Return(_a0 bool) *InvalidDeliveryAddressError_IsBadRequestError_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *InvalidDeliveryAddressError_IsBadRequestError_Call) RunAndReturn(run func() bool) *InvalidDeliveryAddressError_IsBadRequestError_Call {
	_c.Call.Return(run)
	return _c
}

// IsCustomError provides a mock function with given fields:
func (_m *InvalidDeliveryAddressError) IsCustomError() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// InvalidDeliveryAddressError_IsCustomError_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'IsCustomError'
type InvalidDeliveryAddressError_IsCustomError_Call struct {
	*mock.Call
}

// IsCustomError is a helper method to define mock.On call
func (_e *InvalidDeliveryAddressError_Expecter) IsCustomError() *InvalidDeliveryAddressError_IsCustomError_Call {
	return &InvalidDeliveryAddressError_IsCustomError_Call{Call: _e.mock.On("IsCustomError")}
}

func (_c *InvalidDeliveryAddressError_IsCustomError_Call) Run(run func()) *InvalidDeliveryAddressError_IsCustomError_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *InvalidDeliveryAddressError_IsCustomError_Call) Return(_a0 bool) *InvalidDeliveryAddressError_IsCustomError_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *InvalidDeliveryAddressError_IsCustomError_Call) RunAndReturn(run func() bool) *InvalidDeliveryAddressError_IsCustomError_Call {
	_c.Call.Return(run)
	return _c
}

// IsInvalidDeliveryAddressError provides a mock function with given fields:
func (_m *InvalidDeliveryAddressError) IsInvalidDeliveryAddressError() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// InvalidDeliveryAddressError_IsInvalidDeliveryAddressError_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'IsInvalidDeliveryAddressError'
type InvalidDeliveryAddressError_IsInvalidDeliveryAddressError_Call struct {
	*mock.Call
}

// IsInvalidDeliveryAddressError is a helper method to define mock.On call
func (_e *InvalidDeliveryAddressError_Expecter) IsInvalidDeliveryAddressError() *InvalidDeliveryAddressError_IsInvalidDeliveryAddressError_Call {
	return &InvalidDeliveryAddressError_IsInvalidDeliveryAddressError_Call{Call: _e.mock.On("IsInvalidDeliveryAddressError")}
}

func (_c *InvalidDeliveryAddressError_IsInvalidDeliveryAddressError_Call) Run(run func()) *InvalidDeliveryAddressError_IsInvalidDeliveryAddressError_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *InvalidDeliveryAddressError_IsInvalidDeliveryAddressError_Call) Return(_a0 bool) *InvalidDeliveryAddressError_IsInvalidDeliveryAddressError_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *InvalidDeliveryAddressError_IsInvalidDeliveryAddressError_Call) RunAndReturn(run func() bool) *InvalidDeliveryAddressError_IsInvalidDeliveryAddressError_Call {
	_c.Call.Return(run)
	return _c
}

// Message provides a mock function with given fields:
func (_m *InvalidDeliveryAddressError) Message() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// InvalidDeliveryAddressError_Message_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Message'
type InvalidDeliveryAddressError_Message_Call struct {
	*mock.Call
}

// Message is a helper method to define mock.On call
func (_e *InvalidDeliveryAddressError_Expecter) Message() *InvalidDeliveryAddressError_Message_Call {
	return &InvalidDeliveryAddressError_Message_Call{Call: _e.mock.On("Message")}
}

func (_c *InvalidDeliveryAddressError_Message_Call) Run(run func()) *InvalidDeliveryAddressError_Message_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *InvalidDeliveryAddressError_Message_Call) Return(_a0 string) *InvalidDeliveryAddressError_Message_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *InvalidDeliveryAddressError_Message_Call) RunAndReturn(run func() string) *InvalidDeliveryAddressError_Message_Call {
	_c.Call.Return(run)
	return _c
}

// Status provides a mock function with given fields:
func (_m *InvalidDeliveryAddressError) Status() int {
	ret := _m.Called()

	var r0 int
	if rf, ok := ret.Get(0).(func() int); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int)
	}

	return r0
}

// InvalidDeliveryAddressError_Status_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Status'
type InvalidDeliveryAddressError_Status_Call struct {
	*mock.Call
}

// Status is a helper method to define mock.On call
func (_e *InvalidDeliveryAddressError_Expecter) Status() *InvalidDeliveryAddressError_Status_Call {
	return &InvalidDeliveryAddressError_Status_Call{Call: _e.mock.On("Status")}
}

func (_c *InvalidDeliveryAddressError_Status_Call) Run(run func()) *InvalidDeliveryAddressError_Status_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *InvalidDeliveryAddressError_Status_Call) Return(_a0 int) *InvalidDeliveryAddressError_Status_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *InvalidDeliveryAddressError_Status_Call) RunAndReturn(run func() int) *InvalidDeliveryAddressError_Status_Call {
	_c.Call.Return(run)
	return _c
}

// Unwrap provides a mock function with given fields:
func (_m *InvalidDeliveryAddressError) Unwrap() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// InvalidDeliveryAddressError_Unwrap_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Unwrap'
type InvalidDeliveryAddressError_Unwrap_Call struct {
	*mock.Call
}

// Unwrap is a helper method to define mock.On call
func (_e *InvalidDeliveryAddressError_Expecter) Unwrap() *InvalidDeliveryAddressError_Unwrap_Call {
	return &InvalidDeliveryAddressError_Unwrap_Call{Call: _e.mock.On("Unwrap")}
}

func (_c *InvalidDeliveryAddressError_Unwrap_Call) Run(run func()) *InvalidDeliveryAddressError_Unwrap_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *InvalidDeliveryAddressError_Unwrap_Call) Return(_a0 error) *InvalidDeliveryAddressError_Unwrap_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *InvalidDeliveryAddressError_Unwrap_Call) RunAndReturn(run func() error) *InvalidDeliveryAddressError_Unwrap_Call {
	_c.Call.Return(run)
	return _c
}

// NewInvalidDeliveryAddressError creates a new instance of InvalidDeliveryAddressError. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewInvalidDeliveryAddressError(t interface {
	mock.TestingT
	Cleanup(func())
}) *InvalidDeliveryAddressError {
	mock := &InvalidDeliveryAddressError{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}