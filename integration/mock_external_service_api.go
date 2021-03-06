// Code generated by mockery v2.3.0. DO NOT EDIT.

package integration

import (
	http "net/http"

	httpclient "cto-github.cisco.com/NFV-BU/go-msx/httpclient"
	mock "github.com/stretchr/testify/mock"

	url "net/url"
)

// MockExternalServiceApi is an autogenerated mock type for the ExternalServiceApi type
type MockExternalServiceApi struct {
	mock.Mock
}

// AddInterceptor provides a mock function with given fields: interceptor
func (_m *MockExternalServiceApi) AddInterceptor(interceptor httpclient.RequestInterceptor) {
	_m.Called(interceptor)
}

// Do provides a mock function with given fields: req, responseBody
func (_m *MockExternalServiceApi) Do(req *http.Request, responseBody interface{}) (*http.Response, []byte, error) {
	ret := _m.Called(req, responseBody)

	var r0 *http.Response
	if rf, ok := ret.Get(0).(func(*http.Request, interface{}) *http.Response); ok {
		r0 = rf(req, responseBody)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*http.Response)
		}
	}

	var r1 []byte
	if rf, ok := ret.Get(1).(func(*http.Request, interface{}) []byte); ok {
		r1 = rf(req, responseBody)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).([]byte)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(*http.Request, interface{}) error); ok {
		r2 = rf(req, responseBody)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// Request provides a mock function with given fields: endpoint, uriVariables, queryVariables, headers, body
func (_m *MockExternalServiceApi) Request(endpoint Endpoint, uriVariables map[string]string, queryVariables url.Values, headers http.Header, body []byte) (*http.Request, error) {
	ret := _m.Called(endpoint, uriVariables, queryVariables, headers, body)

	var r0 *http.Request
	if rf, ok := ret.Get(0).(func(Endpoint, map[string]string, url.Values, http.Header, []byte) *http.Request); ok {
		r0 = rf(endpoint, uriVariables, queryVariables, headers, body)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*http.Request)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(Endpoint, map[string]string, url.Values, http.Header, []byte) error); ok {
		r1 = rf(endpoint, uriVariables, queryVariables, headers, body)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
