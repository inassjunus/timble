// Code generated by mockery v2.52.2. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// APIResponser is an autogenerated mock type for the APIResponser type
type APIResponser struct {
	mock.Mock
}

// ToBytes provides a mock function with no fields
func (_m *APIResponser) ToBytes() []byte {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for ToBytes")
	}

	var r0 []byte
	if rf, ok := ret.Get(0).(func() []byte); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	return r0
}

// NewAPIResponser creates a new instance of APIResponser. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewAPIResponser(t interface {
	mock.TestingT
	Cleanup(func())
}) *APIResponser {
	mock := &APIResponser{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
