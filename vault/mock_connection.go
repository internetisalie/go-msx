// Code generated by mockery v2.9.4. DO NOT EDIT.

package vault

import (
	context "context"

	api "github.com/hashicorp/vault/api"

	mock "github.com/stretchr/testify/mock"

	tls "crypto/tls"

	x509 "crypto/x509"
)

// MockConnection is an autogenerated mock type for the ConnectionApi type
type MockConnection struct {
	mock.Mock
}

// CreateTransitKey provides a mock function with given fields: ctx, keyName, request
func (_m *MockConnection) CreateTransitKey(ctx context.Context, keyName string, request CreateTransitKeyRequest) error {
	ret := _m.Called(ctx, keyName, request)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, CreateTransitKeyRequest) error); ok {
		r0 = rf(ctx, keyName, request)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteVersionedMetadata provides a mock function with given fields: ctx, path
func (_m *MockConnection) DeleteVersionedMetadata(ctx context.Context, path string) error {
	ret := _m.Called(ctx, path)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, path)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteVersionedSecrets provides a mock function with given fields: ctx, path, request
func (_m *MockConnection) DeleteVersionedSecrets(ctx context.Context, path string, request VersionRequest) error {
	ret := _m.Called(ctx, path, request)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, VersionRequest) error); ok {
		r0 = rf(ctx, path, request)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteVersionedSecretsLatest provides a mock function with given fields: ctx, p
func (_m *MockConnection) DeleteVersionedSecretsLatest(ctx context.Context, p string) error {
	ret := _m.Called(ctx, p)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, p)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DestroyVersionedSecrets provides a mock function with given fields: ctx, path, request
func (_m *MockConnection) DestroyVersionedSecrets(ctx context.Context, path string, request VersionRequest) error {
	ret := _m.Called(ctx, path, request)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, VersionRequest) error); ok {
		r0 = rf(ctx, path, request)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GenerateRandomBytes provides a mock function with given fields: ctx, length
func (_m *MockConnection) GenerateRandomBytes(ctx context.Context, length int) ([]byte, error) {
	ret := _m.Called(ctx, length)

	var r0 []byte
	if rf, ok := ret.Get(0).(func(context.Context, int) []byte); ok {
		r0 = rf(ctx, length)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, length)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetVersionedMetadata provides a mock function with given fields: ctx, path
func (_m *MockConnection) GetVersionedMetadata(ctx context.Context, path string) (VersionedMetadata, error) {
	ret := _m.Called(ctx, path)

	var r0 VersionedMetadata
	if rf, ok := ret.Get(0).(func(context.Context, string) VersionedMetadata); ok {
		r0 = rf(ctx, path)
	} else {
		r0 = ret.Get(0).(VersionedMetadata)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, path)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetVersionedSecrets provides a mock function with given fields: ctx, path, version
func (_m *MockConnection) GetVersionedSecrets(ctx context.Context, path string, version *int) (map[string]interface{}, error) {
	ret := _m.Called(ctx, path, version)

	var r0 map[string]interface{}
	if rf, ok := ret.Get(0).(func(context.Context, string, *int) map[string]interface{}); ok {
		r0 = rf(ctx, path, version)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]interface{})
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, *int) error); ok {
		r1 = rf(ctx, path, version)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Health provides a mock function with given fields: ctx
func (_m *MockConnection) Health(ctx context.Context) (*api.HealthResponse, error) {
	ret := _m.Called(ctx)

	var r0 *api.HealthResponse
	if rf, ok := ret.Get(0).(func(context.Context) *api.HealthResponse); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*api.HealthResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IssueCertificate provides a mock function with given fields: ctx, role, request
func (_m *MockConnection) IssueCertificate(ctx context.Context, role string, request IssueCertificateRequest) (*tls.Certificate, error) {
	ret := _m.Called(ctx, role, request)

	var r0 *tls.Certificate
	if rf, ok := ret.Get(0).(func(context.Context, string, IssueCertificateRequest) *tls.Certificate); ok {
		r0 = rf(ctx, role, request)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*tls.Certificate)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, IssueCertificateRequest) error); ok {
		r1 = rf(ctx, role, request)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListSecrets provides a mock function with given fields: ctx, path
func (_m *MockConnection) ListSecrets(ctx context.Context, path string) (map[string]string, error) {
	ret := _m.Called(ctx, path)

	var r0 map[string]string
	if rf, ok := ret.Get(0).(func(context.Context, string) map[string]string); ok {
		r0 = rf(ctx, path)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]string)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, path)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// LoginWithAppRole provides a mock function with given fields: ctx, roleId, secretId
func (_m *MockConnection) LoginWithAppRole(ctx context.Context, roleId string, secretId string) (string, error) {
	ret := _m.Called(ctx, roleId, secretId)

	var r0 string
	if rf, ok := ret.Get(0).(func(context.Context, string, string) string); ok {
		r0 = rf(ctx, roleId, secretId)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, roleId, secretId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// LoginWithKubernetes provides a mock function with given fields: ctx, jwt, role
func (_m *MockConnection) LoginWithKubernetes(ctx context.Context, jwt string, role string) (string, error) {
	ret := _m.Called(ctx, jwt, role)

	var r0 string
	if rf, ok := ret.Get(0).(func(context.Context, string, string) string); ok {
		r0 = rf(ctx, jwt, role)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, jwt, role)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PatchVersionedSecrets provides a mock function with given fields: ctx, path, request
func (_m *MockConnection) PatchVersionedSecrets(ctx context.Context, path string, request VersionedWriteRequest) error {
	ret := _m.Called(ctx, path, request)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, VersionedWriteRequest) error); ok {
		r0 = rf(ctx, path, request)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ReadCaCertificate provides a mock function with given fields: ctx
func (_m *MockConnection) ReadCaCertificate(ctx context.Context) (*x509.Certificate, error) {
	ret := _m.Called(ctx)

	var r0 *x509.Certificate
	if rf, ok := ret.Get(0).(func(context.Context) *x509.Certificate); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*x509.Certificate)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RemoveSecrets provides a mock function with given fields: ctx, path
func (_m *MockConnection) RemoveSecrets(ctx context.Context, path string) error {
	ret := _m.Called(ctx, path)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, path)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// StoreSecrets provides a mock function with given fields: ctx, path, secrets
func (_m *MockConnection) StoreSecrets(ctx context.Context, path string, secrets map[string]string) error {
	ret := _m.Called(ctx, path, secrets)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, map[string]string) error); ok {
		r0 = rf(ctx, path, secrets)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// StoreVersionedMetadata provides a mock function with given fields: ctx, path, request
func (_m *MockConnection) StoreVersionedMetadata(ctx context.Context, path string, request VersionedMetadataRequest) error {
	ret := _m.Called(ctx, path, request)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, VersionedMetadataRequest) error); ok {
		r0 = rf(ctx, path, request)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// StoreVersionedSecrets provides a mock function with given fields: ctx, path, request
func (_m *MockConnection) StoreVersionedSecrets(ctx context.Context, path string, request VersionedWriteRequest) error {
	ret := _m.Called(ctx, path, request)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, VersionedWriteRequest) error); ok {
		r0 = rf(ctx, path, request)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// TransitBulkDecrypt provides a mock function with given fields: ctx, keyName, ciphertext
func (_m *MockConnection) TransitBulkDecrypt(ctx context.Context, keyName string, ciphertext ...string) ([]string, error) {
	_va := make([]interface{}, len(ciphertext))
	for _i := range ciphertext {
		_va[_i] = ciphertext[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, keyName)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 []string
	if rf, ok := ret.Get(0).(func(context.Context, string, ...string) []string); ok {
		r0 = rf(ctx, keyName, ciphertext...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, ...string) error); ok {
		r1 = rf(ctx, keyName, ciphertext...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// TransitDecrypt provides a mock function with given fields: ctx, keyName, ciphertext
func (_m *MockConnection) TransitDecrypt(ctx context.Context, keyName string, ciphertext string) (string, error) {
	ret := _m.Called(ctx, keyName, ciphertext)

	var r0 string
	if rf, ok := ret.Get(0).(func(context.Context, string, string) string); ok {
		r0 = rf(ctx, keyName, ciphertext)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, keyName, ciphertext)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// TransitEncrypt provides a mock function with given fields: ctx, keyName, plaintext
func (_m *MockConnection) TransitEncrypt(ctx context.Context, keyName string, plaintext string) (string, error) {
	ret := _m.Called(ctx, keyName, plaintext)

	var r0 string
	if rf, ok := ret.Get(0).(func(context.Context, string, string) string); ok {
		r0 = rf(ctx, keyName, plaintext)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, keyName, plaintext)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UndeleteVersionedSecrets provides a mock function with given fields: ctx, path, request
func (_m *MockConnection) UndeleteVersionedSecrets(ctx context.Context, path string, request VersionRequest) error {
	ret := _m.Called(ctx, path, request)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, VersionRequest) error); ok {
		r0 = rf(ctx, path, request)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
