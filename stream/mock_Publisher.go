// Code generated by mockery v2.3.0. DO NOT EDIT.

package stream

import (
	message "github.com/ThreeDotsLabs/watermill/message"
	mock "github.com/stretchr/testify/mock"
)

// MockPublisher is an autogenerated mock type for the Publisher type
type MockPublisher struct {
	mock.Mock
}

// Close provides a mock function with given fields:
func (_m *MockPublisher) Close() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Publish provides a mock function with given fields: _a0
func (_m *MockPublisher) Publish(_a0 *message.Message) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(*message.Message) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
