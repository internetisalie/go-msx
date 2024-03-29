// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	context "context"

	exp "github.com/doug-martin/goqu/v9/exp"
	mock "github.com/stretchr/testify/mock"

	paging "cto-github.cisco.com/NFV-BU/go-msx/paging"
)

// CrudRepositoryApi is an autogenerated mock type for the CrudRepositoryApi type
type CrudRepositoryApi struct {
	mock.Mock
}

// CountAll provides a mock function with given fields: ctx, dest
func (_m *CrudRepositoryApi) CountAll(ctx context.Context, dest *int64) error {
	ret := _m.Called(ctx, dest)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *int64) error); ok {
		r0 = rf(ctx, dest)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CountAllBy provides a mock function with given fields: ctx, where, dest
func (_m *CrudRepositoryApi) CountAllBy(ctx context.Context, where map[string]interface{}, dest *int64) error {
	ret := _m.Called(ctx, where, dest)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, map[string]interface{}, *int64) error); ok {
		r0 = rf(ctx, where, dest)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CountAllByExpression provides a mock function with given fields: ctx, where, dest
func (_m *CrudRepositoryApi) CountAllByExpression(ctx context.Context, where exp.Expression, dest *int64) error {
	ret := _m.Called(ctx, where, dest)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, exp.Expression, *int64) error); ok {
		r0 = rf(ctx, where, dest)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteBy provides a mock function with given fields: ctx, where
func (_m *CrudRepositoryApi) DeleteBy(ctx context.Context, where map[string]interface{}) error {
	ret := _m.Called(ctx, where)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, map[string]interface{}) error); ok {
		r0 = rf(ctx, where)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FindAll provides a mock function with given fields: ctx, dest
func (_m *CrudRepositoryApi) FindAll(ctx context.Context, dest interface{}) error {
	ret := _m.Called(ctx, dest)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, interface{}) error); ok {
		r0 = rf(ctx, dest)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FindAllBy provides a mock function with given fields: ctx, where, dest
func (_m *CrudRepositoryApi) FindAllBy(ctx context.Context, where map[string]interface{}, dest interface{}) error {
	ret := _m.Called(ctx, where, dest)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, map[string]interface{}, interface{}) error); ok {
		r0 = rf(ctx, where, dest)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FindAllByExpression provides a mock function with given fields: ctx, where, dest
func (_m *CrudRepositoryApi) FindAllByExpression(ctx context.Context, where exp.Expression, dest interface{}) error {
	ret := _m.Called(ctx, where, dest)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, exp.Expression, interface{}) error); ok {
		r0 = rf(ctx, where, dest)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FindAllDistinctBy provides a mock function with given fields: ctx, distinct, where, dest
func (_m *CrudRepositoryApi) FindAllDistinctBy(ctx context.Context, distinct []string, where map[string]interface{}, dest interface{}) error {
	ret := _m.Called(ctx, distinct, where, dest)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, []string, map[string]interface{}, interface{}) error); ok {
		r0 = rf(ctx, distinct, where, dest)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FindAllPagedBy provides a mock function with given fields: ctx, where, preq, dest
func (_m *CrudRepositoryApi) FindAllPagedBy(ctx context.Context, where map[string]interface{}, preq paging.Request, dest interface{}) (paging.Response, error) {
	ret := _m.Called(ctx, where, preq, dest)

	var r0 paging.Response
	if rf, ok := ret.Get(0).(func(context.Context, map[string]interface{}, paging.Request, interface{}) paging.Response); ok {
		r0 = rf(ctx, where, preq, dest)
	} else {
		r0 = ret.Get(0).(paging.Response)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, map[string]interface{}, paging.Request, interface{}) error); ok {
		r1 = rf(ctx, where, preq, dest)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindAllPagedByExpression provides a mock function with given fields: ctx, where, preq, dest
func (_m *CrudRepositoryApi) FindAllPagedByExpression(ctx context.Context, where exp.Expression, preq paging.Request, dest interface{}) (paging.Response, error) {
	ret := _m.Called(ctx, where, preq, dest)

	var r0 paging.Response
	if rf, ok := ret.Get(0).(func(context.Context, exp.Expression, paging.Request, interface{}) paging.Response); ok {
		r0 = rf(ctx, where, preq, dest)
	} else {
		r0 = ret.Get(0).(paging.Response)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, exp.Expression, paging.Request, interface{}) error); ok {
		r1 = rf(ctx, where, preq, dest)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindAllSortedBy provides a mock function with given fields: ctx, where, sortOrder, dest
func (_m *CrudRepositoryApi) FindAllSortedBy(ctx context.Context, where map[string]interface{}, sortOrder paging.SortOrder, dest interface{}) error {
	ret := _m.Called(ctx, where, sortOrder, dest)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, map[string]interface{}, paging.SortOrder, interface{}) error); ok {
		r0 = rf(ctx, where, sortOrder, dest)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FindAllSortedByExpression provides a mock function with given fields: ctx, where, sortOrder, dest
func (_m *CrudRepositoryApi) FindAllSortedByExpression(ctx context.Context, where exp.Expression, sortOrder paging.SortOrder, dest interface{}) error {
	ret := _m.Called(ctx, where, sortOrder, dest)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, exp.Expression, paging.SortOrder, interface{}) error); ok {
		r0 = rf(ctx, where, sortOrder, dest)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FindOneBy provides a mock function with given fields: ctx, where, dest
func (_m *CrudRepositoryApi) FindOneBy(ctx context.Context, where map[string]interface{}, dest interface{}) error {
	ret := _m.Called(ctx, where, dest)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, map[string]interface{}, interface{}) error); ok {
		r0 = rf(ctx, where, dest)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FindOneSortedBy provides a mock function with given fields: ctx, where, sortOrder, dest
func (_m *CrudRepositoryApi) FindOneSortedBy(ctx context.Context, where map[string]interface{}, sortOrder paging.SortOrder, dest interface{}) error {
	ret := _m.Called(ctx, where, sortOrder, dest)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, map[string]interface{}, paging.SortOrder, interface{}) error); ok {
		r0 = rf(ctx, where, sortOrder, dest)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Insert provides a mock function with given fields: ctx, value
func (_m *CrudRepositoryApi) Insert(ctx context.Context, value interface{}) error {
	ret := _m.Called(ctx, value)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, interface{}) error); ok {
		r0 = rf(ctx, value)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Save provides a mock function with given fields: ctx, value
func (_m *CrudRepositoryApi) Save(ctx context.Context, value interface{}) error {
	ret := _m.Called(ctx, value)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, interface{}) error); ok {
		r0 = rf(ctx, value)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SaveAll provides a mock function with given fields: ctx, values
func (_m *CrudRepositoryApi) SaveAll(ctx context.Context, values []interface{}) error {
	ret := _m.Called(ctx, values)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, []interface{}) error); ok {
		r0 = rf(ctx, values)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Truncate provides a mock function with given fields: ctx
func (_m *CrudRepositoryApi) Truncate(ctx context.Context) error {
	ret := _m.Called(ctx)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Update provides a mock function with given fields: ctx, where, value
func (_m *CrudRepositoryApi) Update(ctx context.Context, where map[string]interface{}, value interface{}) error {
	ret := _m.Called(ctx, where, value)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, map[string]interface{}, interface{}) error); ok {
		r0 = rf(ctx, where, value)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewCrudRepositoryApi interface {
	mock.TestingT
	Cleanup(func())
}

// NewCrudRepositoryApi creates a new instance of CrudRepositoryApi. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewCrudRepositoryApi(t mockConstructorTestingTNewCrudRepositoryApi) *CrudRepositoryApi {
	mock := &CrudRepositoryApi{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
