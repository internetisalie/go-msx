// Code generated by mockery v2.14.0. DO NOT EDIT.

package idm

import (
	integration "cto-github.cisco.com/NFV-BU/go-msx/integration"
	mock "github.com/stretchr/testify/mock"

	paging "cto-github.cisco.com/NFV-BU/go-msx/paging"
)

// MockIdm is an autogenerated mock type for the Api type
type MockIdm struct {
	mock.Mock
}

// BatchCreateCapabilities provides a mock function with given fields: populator, owner, capabilities
func (_m *MockIdm) BatchCreateCapabilities(populator bool, owner string, capabilities []CapabilityCreateRequest) (*integration.MsxResponse, error) {
	ret := _m.Called(populator, owner, capabilities)

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func(bool, string, []CapabilityCreateRequest) *integration.MsxResponse); ok {
		r0 = rf(populator, owner, capabilities)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(bool, string, []CapabilityCreateRequest) error); ok {
		r1 = rf(populator, owner, capabilities)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// BatchUpdateCapabilities provides a mock function with given fields: populator, owner, capabilities
func (_m *MockIdm) BatchUpdateCapabilities(populator bool, owner string, capabilities []CapabilityUpdateRequest) (*integration.MsxResponse, error) {
	ret := _m.Called(populator, owner, capabilities)

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func(bool, string, []CapabilityUpdateRequest) *integration.MsxResponse); ok {
		r0 = rf(populator, owner, capabilities)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(bool, string, []CapabilityUpdateRequest) error); ok {
		r1 = rf(populator, owner, capabilities)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateRole provides a mock function with given fields: dbinstaller, body
func (_m *MockIdm) CreateRole(dbinstaller bool, body RoleCreateRequest) (*integration.MsxResponse, error) {
	ret := _m.Called(dbinstaller, body)

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func(bool, RoleCreateRequest) *integration.MsxResponse); ok {
		r0 = rf(dbinstaller, body)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(bool, RoleCreateRequest) error); ok {
		r1 = rf(dbinstaller, body)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteCapability provides a mock function with given fields: populator, owner, name
func (_m *MockIdm) DeleteCapability(populator bool, owner string, name string) (*integration.MsxResponse, error) {
	ret := _m.Called(populator, owner, name)

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func(bool, string, string) *integration.MsxResponse); ok {
		r0 = rf(populator, owner, name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(bool, string, string) error); ok {
		r1 = rf(populator, owner, name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteRole provides a mock function with given fields: roleName
func (_m *MockIdm) DeleteRole(roleName string) (*integration.MsxResponse, error) {
	ret := _m.Called(roleName)

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func(string) *integration.MsxResponse); ok {
		r0 = rf(roleName)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(roleName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetCapabilities provides a mock function with given fields: p
func (_m *MockIdm) GetCapabilities(p paging.Request) (*integration.MsxResponse, error) {
	ret := _m.Called(p)

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func(paging.Request) *integration.MsxResponse); ok {
		r0 = rf(p)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(paging.Request) error); ok {
		r1 = rf(p)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetMyProvider provides a mock function with given fields:
func (_m *MockIdm) GetMyProvider() (*integration.MsxResponse, error) {
	ret := _m.Called()

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func() *integration.MsxResponse); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
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

// GetProviderByName provides a mock function with given fields: providerName
func (_m *MockIdm) GetProviderByName(providerName string) (*integration.MsxResponse, error) {
	ret := _m.Called(providerName)

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func(string) *integration.MsxResponse); ok {
		r0 = rf(providerName)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(providerName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetProviderExtensionByName provides a mock function with given fields: name
func (_m *MockIdm) GetProviderExtensionByName(name string) (*integration.MsxResponse, error) {
	ret := _m.Called(name)

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func(string) *integration.MsxResponse); ok {
		r0 = rf(name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetRoles provides a mock function with given fields: resolvePermissionNames, p
func (_m *MockIdm) GetRoles(resolvePermissionNames bool, p paging.Request) (*integration.MsxResponse, error) {
	ret := _m.Called(resolvePermissionNames, p)

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func(bool, paging.Request) *integration.MsxResponse); ok {
		r0 = rf(resolvePermissionNames, p)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(bool, paging.Request) error); ok {
		r1 = rf(resolvePermissionNames, p)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetTenantById provides a mock function with given fields: tenantId
func (_m *MockIdm) GetTenantById(tenantId string) (*integration.MsxResponse, error) {
	ret := _m.Called(tenantId)

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func(string) *integration.MsxResponse); ok {
		r0 = rf(tenantId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(tenantId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetTenantByIdV8 provides a mock function with given fields: tenantId
func (_m *MockIdm) GetTenantByIdV8(tenantId string) (*integration.MsxResponse, error) {
	ret := _m.Called(tenantId)

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func(string) *integration.MsxResponse); ok {
		r0 = rf(tenantId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(tenantId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetTenantByName provides a mock function with given fields: tenantName
func (_m *MockIdm) GetTenantByName(tenantName string) (*integration.MsxResponse, error) {
	ret := _m.Called(tenantName)

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func(string) *integration.MsxResponse); ok {
		r0 = rf(tenantName)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(tenantName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUserById provides a mock function with given fields: userId
func (_m *MockIdm) GetUserById(userId string) (*integration.MsxResponse, error) {
	ret := _m.Called(userId)

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func(string) *integration.MsxResponse); ok {
		r0 = rf(userId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(userId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUserByIdV8 provides a mock function with given fields: userId
func (_m *MockIdm) GetUserByIdV8(userId string) (*integration.MsxResponse, error) {
	ret := _m.Called(userId)

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func(string) *integration.MsxResponse); ok {
		r0 = rf(userId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(userId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IsTokenActive provides a mock function with given fields:
func (_m *MockIdm) IsTokenActive() (*integration.MsxResponse, error) {
	ret := _m.Called()

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func() *integration.MsxResponse); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
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

// UpdateRole provides a mock function with given fields: dbinstaller, body
func (_m *MockIdm) UpdateRole(dbinstaller bool, body RoleUpdateRequest) (*integration.MsxResponse, error) {
	ret := _m.Called(dbinstaller, body)

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func(bool, RoleUpdateRequest) *integration.MsxResponse); ok {
		r0 = rf(dbinstaller, body)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(bool, RoleUpdateRequest) error); ok {
		r1 = rf(dbinstaller, body)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewMockIdm interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockIdm creates a new instance of MockIdm. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockIdm(t mockConstructorTestingTNewMockIdm) *MockIdm {
	mock := &MockIdm{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
