// Code generated by mockery v2.3.0. DO NOT EDIT.

package transit

import mock "github.com/stretchr/testify/mock"

// MockEncrypter is an autogenerated mock type for the Encrypter type
type MockEncrypter struct {
	mock.Mock
}

// CreateKey provides a mock function with given fields:
func (_m *MockEncrypter) CreateKey() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Decode provides a mock function with given fields: insecureValue
func (_m *MockEncrypter) Decode(insecureValue string) (map[string]*string, error) {
	ret := _m.Called(insecureValue)

	var r0 map[string]*string
	if rf, ok := ret.Get(0).(func(string) map[string]*string); ok {
		r0 = rf(insecureValue)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]*string)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(insecureValue)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Decrypt provides a mock function with given fields: secureValue
func (_m *MockEncrypter) Decrypt(secureValue string) (map[string]*string, error) {
	ret := _m.Called(secureValue)

	var r0 map[string]*string
	if rf, ok := ret.Get(0).(func(string) map[string]*string); ok {
		r0 = rf(secureValue)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]*string)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(secureValue)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DecryptSet provides a mock function with given fields: set
func (_m *MockEncrypter) DecryptSet(set BulkSet) error {
	ret := _m.Called(set)

	var r0 error
	if rf, ok := ret.Get(0).(func(BulkSet) error); ok {
		r0 = rf(set)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DecryptSets provides a mock function with given fields: sets
func (_m *MockEncrypter) DecryptSets(sets BulkSets) error {
	ret := _m.Called(sets)

	var r0 error
	if rf, ok := ret.Get(0).(func(BulkSets) error); ok {
		r0 = rf(sets)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Encrypt provides a mock function with given fields: secureValue
func (_m *MockEncrypter) Encrypt(value map[string]*string) (string, bool, error) {
	ret := _m.Called(value)

	var r0 string
	if rf, ok := ret.Get(0).(func(map[string]*string) string); ok {
		r0 = rf(value)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 bool
	if rf, ok := ret.Get(1).(func(map[string]*string) bool); ok {
		r1 = rf(value)
	} else {
		r1 = ret.Get(1).(bool)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(map[string]*string) error); ok {
		r2 = rf(value)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}
