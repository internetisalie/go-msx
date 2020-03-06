// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	context "context"
	time "time"

	mock "github.com/stretchr/testify/mock"
)

// CrudRepositoryApi is an autogenerated mock type for the CrudRepositoryApi type
type CrudRepositoryApi struct {
	mock.Mock
}

// CountAll provides a mock function with given fields: ctx, dest
func (_m *CrudRepositoryApi) CountAll(ctx context.Context, dest interface{}) error {
	ret := _m.Called(ctx, dest)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, interface{}) error); ok {
		r0 = rf(ctx, dest)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CountAllBy provides a mock function with given fields: ctx, where, dest
func (_m *CrudRepositoryApi) CountAllBy(ctx context.Context, where map[string]interface{}, dest interface{}) error {
	ret := _m.Called(ctx, where, dest)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, map[string]interface{}, interface{}) error); ok {
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

// FindAllByLuceneSearch provides a mock function with given fields: ctx, index, search, dest
func (_m *CrudRepositoryApi) FindAllByLuceneSearch(ctx context.Context, index string, search string, dest interface{}) error {
	ret := _m.Called(ctx, index, search, dest)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, interface{}) error); ok {
		r0 = rf(ctx, index, search, dest)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FindAllCql provides a mock function with given fields: ctx, stmt, names, where, dest
func (_m *CrudRepositoryApi) FindAllCql(ctx context.Context, stmt string, names []string, where map[string]interface{}, dest interface{}) error {
	ret := _m.Called(ctx, stmt, names, where, dest)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, []string, map[string]interface{}, interface{}) error); ok {
		r0 = rf(ctx, stmt, names, where, dest)
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

// FindPartitionKeys provides a mock function with given fields: ctx, dest
func (_m *CrudRepositoryApi) FindPartitionKeys(ctx context.Context, dest interface{}) error {
	ret := _m.Called(ctx, dest)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, interface{}) error); ok {
		r0 = rf(ctx, dest)
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

// SaveWithTtl provides a mock function with given fields: ctx, value, ttl
func (_m *CrudRepositoryApi) SaveWithTtl(ctx context.Context, value interface{}, ttl time.Duration) error {
	ret := _m.Called(ctx, value, ttl)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, interface{}, time.Duration) error); ok {
		r0 = rf(ctx, value, ttl)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateBy provides a mock function with given fields: ctx, where, values
func (_m *CrudRepositoryApi) UpdateBy(ctx context.Context, where map[string]interface{}, values map[string]interface{}) error {
	ret := _m.Called(ctx, where, values)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, map[string]interface{}, map[string]interface{}) error); ok {
		r0 = rf(ctx, where, values)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
