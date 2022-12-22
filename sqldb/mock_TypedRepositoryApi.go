// Code generated by mockery v2.14.0. DO NOT EDIT.

package sqldb

import (
	context "context"

	goqu "github.com/doug-martin/goqu/v9"
	exp "github.com/doug-martin/goqu/v9/exp"

	mock "github.com/stretchr/testify/mock"

	paging "cto-github.cisco.com/NFV-BU/go-msx/paging"
)

// MockTypedRepositoryApi is an autogenerated mock type for the TypedRepositoryApi type
type MockTypedRepositoryApi[I interface{}] struct {
	mock.Mock
}

// CountAll provides a mock function with given fields: ctx, dest, where
func (_m *MockTypedRepositoryApi[I]) CountAll(ctx context.Context, dest *int64, where WhereOption) error {
	ret := _m.Called(ctx, dest, where)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *int64, WhereOption) error); ok {
		r0 = rf(ctx, dest, where)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteAll provides a mock function with given fields: ctx, where
func (_m *MockTypedRepositoryApi[I]) DeleteAll(ctx context.Context, where WhereOption) error {
	ret := _m.Called(ctx, where)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, WhereOption) error); ok {
		r0 = rf(ctx, where)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteOne provides a mock function with given fields: ctx, keys
func (_m *MockTypedRepositoryApi[I]) DeleteOne(ctx context.Context, keys exp.Ex) error {
	ret := _m.Called(ctx, keys)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, exp.Ex) error); ok {
		r0 = rf(ctx, keys)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FindAll provides a mock function with given fields: ctx, dest, options
func (_m *MockTypedRepositoryApi[I]) FindAll(ctx context.Context, dest *[]I, options ...func(*goqu.SelectDataset, *paging.Request) (*goqu.SelectDataset, *paging.Request)) (paging.Response, error) {
	_va := make([]interface{}, len(options))
	for _i := range options {
		_va[_i] = options[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, dest)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 paging.Response
	if rf, ok := ret.Get(0).(func(context.Context, *[]I, ...func(*goqu.SelectDataset, *paging.Request) (*goqu.SelectDataset, *paging.Request)) paging.Response); ok {
		r0 = rf(ctx, dest, options...)
	} else {
		r0 = ret.Get(0).(paging.Response)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *[]I, ...func(*goqu.SelectDataset, *paging.Request) (*goqu.SelectDataset, *paging.Request)) error); ok {
		r1 = rf(ctx, dest, options...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindOne provides a mock function with given fields: ctx, dest, where
func (_m *MockTypedRepositoryApi[I]) FindOne(ctx context.Context, dest *I, where WhereOption) error {
	ret := _m.Called(ctx, dest, where)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *I, WhereOption) error); ok {
		r0 = rf(ctx, dest, where)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Insert provides a mock function with given fields: ctx, value
func (_m *MockTypedRepositoryApi[I]) Insert(ctx context.Context, value I) error {
	ret := _m.Called(ctx, value)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, I) error); ok {
		r0 = rf(ctx, value)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Truncate provides a mock function with given fields: ctx
func (_m *MockTypedRepositoryApi[I]) Truncate(ctx context.Context) error {
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
func (_m *MockTypedRepositoryApi[I]) Update(ctx context.Context, where WhereOption, value I) error {
	ret := _m.Called(ctx, where, value)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, WhereOption, I) error); ok {
		r0 = rf(ctx, where, value)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Upsert provides a mock function with given fields: ctx, value
func (_m *MockTypedRepositoryApi[I]) Upsert(ctx context.Context, value I) error {
	ret := _m.Called(ctx, value)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, I) error); ok {
		r0 = rf(ctx, value)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewMockTypedRepositoryApi interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockTypedRepositoryApi creates a new instance of MockTypedRepositoryApi. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockTypedRepositoryApi[I interface{}](t mockConstructorTestingTNewMockTypedRepositoryApi) *MockTypedRepositoryApi[I] {
	mock := &MockTypedRepositoryApi[I]{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}