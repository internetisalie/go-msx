// Code generated by mockery v2.3.0. DO NOT EDIT.

package fs

import (
	os "os"

	mock "github.com/stretchr/testify/mock"
)

// MockFile is an autogenerated mock type for the File type
type MockFile struct {
	mock.Mock
}

// Close provides a mock function with given fields:
func (_m *MockFile) Close() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Read provides a mock function with given fields: p
func (_m *MockFile) Read(p []byte) (int, error) {
	ret := _m.Called(p)

	var r0 int
	if rf, ok := ret.Get(0).(func([]byte) int); ok {
		r0 = rf(p)
	} else {
		r0 = ret.Get(0).(int)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func([]byte) error); ok {
		r1 = rf(p)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Readdir provides a mock function with given fields: count
func (_m *MockFile) Readdir(count int) ([]os.FileInfo, error) {
	ret := _m.Called(count)

	var r0 []os.FileInfo
	if rf, ok := ret.Get(0).(func(int) []os.FileInfo); ok {
		r0 = rf(count)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]os.FileInfo)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(count)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Seek provides a mock function with given fields: offset, whence
func (_m *MockFile) Seek(offset int64, whence int) (int64, error) {
	ret := _m.Called(offset, whence)

	var r0 int64
	if rf, ok := ret.Get(0).(func(int64, int) int64); ok {
		r0 = rf(offset, whence)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int64, int) error); ok {
		r1 = rf(offset, whence)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Stat provides a mock function with given fields:
func (_m *MockFile) Stat() (os.FileInfo, error) {
	ret := _m.Called()

	var r0 os.FileInfo
	if rf, ok := ret.Get(0).(func() os.FileInfo); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(os.FileInfo)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}