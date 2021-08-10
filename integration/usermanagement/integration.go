package usermanagement

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/httpclient"
	"cto-github.cisco.com/NFV-BU/go-msx/integration"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"cto-github.cisco.com/NFV-BU/go-msx/paging"
	"cto-github.cisco.com/NFV-BU/go-msx/security"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"encoding/json"
	"github.com/pkg/errors"
	"net/http"
	"net/url"
	"strconv"
)

const (
	endpointNameGetAdminHealth = "getAdminHealth"

	endpointNameLogin  = "login"
	endpointNameLogout = "logout"

	endpointNameIsTokenValid    = "isTokenValid"
	endpointNameGetTokenDetails = "getTokenDetails"

	endpointNameGetMyProvider = "getMyProvider"

	endpointNameGetProviderByName = "getProviderByName"

	endpointNameGetProviderExtensionByName = "getProviderExtensionByName"

	endpointNameGetTenantById   = "getTenantById"
	endpointNameGetTenantByIdV8 = "getTenantByIdV8"
	endpointNameGetTenantByName = "getTenantByName"

	endpointNameGetUserById   = "getUserById"
	endpointNameGetUserByIdV8 = "getUserByIdV8"

	endpointNameGetSystemSecrets      = "getSystemSecrets"
	endpointNameEncryptSystemSecrets  = "encryptSystemSecrets"
	endpointNameAddSystemSecrets      = "addSystemSecrets"
	endpointNameReplaceSystemSecrets  = "replaceSystemSecrets"
	endpointNameRemoveSystemSecrets   = "removeSystemSecrets"
	endpointNameGenerateSystemSecrets = "generateSystemSecrets"

	endpointNameGetTenantSecrets      = "getTenantSecrets"
	endpointNameEncryptTenantSecrets  = "encryptTenantSecrets"
	endpointNameAddTenantSecrets      = "addTenantSecrets"
	endpointNameReplaceTenantSecrets  = "replaceTenantSecrets"
	endpointNameRemoveTenantSecrets   = "removeTenantSecrets"
	endpointNameGenerateTenantSecrets = "generateTenantSecrets"

	endpointNameGetRoles   = "getRoles"
	endpointNameCreateRole = "createRole"
	endpointNameUpdateRole = "updateRole"
	endpointNameDeleteRole = "deleteRole"

	endpointNameGetCapabilities         = "getCapabilities"
	endpointNameBatchCreateCapabilities = "batchCreateCapabilities"
	endpointNameBatchUpdateCapabilities = "batchUpdateCapabilities"
	endpointNameDeleteCapability        = "deleteCapability"

	endpointNameGetSecretPolicy   = "getSecretPolicy"
	endpointNameSetSecretPolicy   = "setSecretPolicy"
	endpointNameUnsetSecretPolicy = "unsetSecretPolicy"

	endpointTenantHierarchyRoot      = "getTenantHierarchyRoot"
	endpointTenantHierarchyParent    = "getTenantHierarchyParent"
	endpointTenantHierarchyAncestors = "getTenantHierarchyAncestors"

	serviceName = integration.ServiceNameUserManagement
)

var (
	logger    = log.NewLogger("msx.integration.usermanagement")
	endpoints = map[string]integration.MsxServiceEndpoint{
		endpointNameGetAdminHealth: {Method: "GET", Path: "/admin/health"},

		endpointNameLogin:  {Method: "POST", Path: "/v2/token"},
		endpointNameLogout: {Method: "GET", Path: "/v2/logout"},

		endpointNameIsTokenValid:    {Method: "GET", Path: "/api/v1/isTokenValid"},
		endpointNameGetTokenDetails: {Method: "POST", Path: "/v2/check_token"},

		endpointNameGetMyProvider:     {Method: "GET", Path: "/api/v1/providers"},
		endpointNameGetProviderByName: {Method: "GET", Path: "/api/v1/providers/{{.providerName}}"},

		endpointNameGetProviderExtensionByName: {Method: "GET", Path: "/api/v1/providers/providerextension/parameters/{{.providerExtensionName}}"},

		endpointNameGetTenantById:   {Method: "GET", Path: "/api/v3/tenants/{{.tenantId}}"},
		endpointNameGetTenantByIdV8: {Method: "GET", Path: "/api/v8/tenants/{{.tenantId}}"},
		endpointNameGetTenantByName: {Method: "GET", Path: "/api/v1/tenants/{{.tenantName}}"},

		endpointNameGetUserById:   {Method: "GET", Path: "/api/v2/users/{{.userId}}"},
		endpointNameGetUserByIdV8: {Method: "GET", Path: "/api/v8/users/{{.userId}}"},

		endpointNameGetSystemSecrets:      {Method: "GET", Path: "/api/v2/secrets/scope/{{.scope}}"},
		endpointNameAddSystemSecrets:      {Method: "POST", Path: "/api/v2/secrets/scope/{{.scope}}"},
		endpointNameReplaceSystemSecrets:  {Method: "PUT", Path: "/api/v2/secrets/scope/{{.scope}}"},
		endpointNameRemoveSystemSecrets:   {Method: "DELETE", Path: "/api/v2/secrets/scope/{{.scope}}"},
		endpointNameEncryptSystemSecrets:  {Method: "POST", Path: "/api/v2/secrets/scope/{{.scope}}/encrypt"},
		endpointNameGenerateSystemSecrets: {Method: "POST", Path: "/api/v2/secrets/scope/{{.scope}}/generate"},

		endpointNameGetTenantSecrets:      {Method: "GET", Path: "/api/v2/secrets/tenant/{{.tenantId}}/scope/{{.scope}}"},
		endpointNameAddTenantSecrets:      {Method: "POST", Path: "/api/v2/secrets/tenant/{{.tenantId}}/scope/{{.scope}}"},
		endpointNameReplaceTenantSecrets:  {Method: "PUT", Path: "/api/v2/secrets/tenant/{{.tenantId}}/scope/{{.scope}}"},
		endpointNameRemoveTenantSecrets:   {Method: "DELETE", Path: "/api/v2/secrets/tenant/{{.tenantId}}/scope/{{.scope}}"},
		endpointNameEncryptTenantSecrets:  {Method: "POST", Path: "/api/v2/secrets/tenant/{{.tenantId}}/scope/{{.scope}}/encrypt"},
		endpointNameGenerateTenantSecrets: {Method: "POST", Path: "/api/v2/secrets/tenant/{{.tenantId}}/scope/{{.scope}}/generate"},

		endpointNameGetRoles:   {Method: "GET", Path: "/api/v1/roles"},
		endpointNameCreateRole: {Method: "POST", Path: "/api/v1/roles"},
		endpointNameUpdateRole: {Method: "PUT", Path: "/api/v1/roles/{roleName}"},
		endpointNameDeleteRole: {Method: "DELETE", Path: "/api/v1/roles/{roleName}"},

		endpointNameGetCapabilities:         {Method: "GET", Path: "/api/v1/roles/capabilities"},
		endpointNameBatchCreateCapabilities: {Method: "POST", Path: "/api/v1/roles/capabilities"},
		endpointNameBatchUpdateCapabilities: {Method: "PUT", Path: "/api/v1/roles/capabilities"},
		endpointNameDeleteCapability:        {Method: "DELETE", Path: "/api/v1/roles/capabilities/{capabilityName}"},

		endpointNameGetSecretPolicy:   {Method: "GET", Path: "/api/v2/secrets/policy/{policyName}"},
		endpointNameSetSecretPolicy:   {Method: "PUT", Path: "/api/v2/secrets/policy/{policyName}"},
		endpointNameUnsetSecretPolicy: {Method: "DELETE", Path: "/api/v2/secrets/policy/{policyName}"},

		endpointTenantHierarchyRoot:      {Method: "GET", Path: "/v2/tenant_hierarchy/root"},
		endpointTenantHierarchyParent:    {Method: "GET", Path: "/v2/tenant_hierarchy/parent"},
		endpointTenantHierarchyAncestors: {Method: "GET", Path: "/v2/tenant_hierarchy/ancestors"},
	}
)

func NewIntegration(ctx context.Context) (Api, error) {
	integrationInstance := IntegrationFromContext(ctx)
	if integrationInstance == nil {
		integrationInstance = &Integration{
			MsxContextServiceExecutor: integration.NewMsxService(ctx, serviceName, endpoints),
		}
	}
	return integrationInstance, nil
}

func NewIntegrationWithExecutor(executor integration.MsxContextServiceExecutor) *Integration {
	return &Integration{
		MsxContextServiceExecutor: executor,
	}
}

type Integration struct {
	integration.MsxContextServiceExecutor
}

func (i *Integration) GetAdminHealth() (result *integration.MsxResponse, err error) {
	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameGetAdminHealth,
		Payload:      &integration.HealthDTO{},
		NoToken:      true,
	})
}

func (i *Integration) Login(user, password string) (result *integration.MsxResponse, err error) {
	securityClientSettings, err := integration.NewSecurityClientSettings(i.Context())
	if err != nil {
		return nil, err
	}

	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameLogin,
		Headers: http.Header(map[string][]string{
			"Authorization": {securityClientSettings.Authorization()},
			"Content-Type":  {httpclient.MimeTypeApplicationWwwFormUrlencoded},
		}),
		Body: []byte(url.Values(map[string][]string{
			"grant_type": {"password"},
			"username":   {user},
			"password":   {password},
		}).Encode()),
		Payload:      new(LoginResponse),
		ErrorPayload: new(integration.OAuthErrorDTO),
		NoToken:      true,
	})
}

func (i *Integration) Logout() (result *integration.MsxResponse, err error) {
	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameLogout,
	})
}

func (i *Integration) IsTokenActive() (*integration.MsxResponse, error) {
	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameIsTokenValid,
		ErrorPayload: new(integration.OAuthErrorDTO),
	})
}

func (i *Integration) GetTokenDetails(noDetails bool) (*integration.MsxResponse, error) {
	securityClientSettings, err := integration.NewSecurityClientSettings(i.Context())
	if err != nil {
		return nil, err
	}

	var headers = make(http.Header)
	headers.Set("Authorization", securityClientSettings.Authorization())
	headers.Set("Content-Type", httpclient.MimeTypeApplicationWwwFormUrlencoded)
	headers.Set("Accept", httpclient.MimeTypeApplicationJson)

	var body = make(url.Values)
	userContext := security.UserContextFromContext(i.Context())
	body.Set("token", userContext.Token)
	if noDetails {
		body.Set("no_details", "true")
	}
	var bodyBytes = []byte(body.Encode())

	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName:   endpointNameGetTokenDetails,
		Headers:        headers,
		Body:           bodyBytes,
		ExpectEnvelope: false,
		NoToken:        true,
		Payload:        new(TokenDetails),
		ErrorPayload:   new(integration.ErrorDTO3),
	})
}

func (i *Integration) GetMyProvider() (*integration.MsxResponse, error) {
	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameGetMyProvider,
		Payload:      &ProviderResponse{},
		ErrorPayload: new(integration.ErrorDTO),
	})
}

func (i *Integration) GetProviderByName(name string) (*integration.MsxResponse, error) {
	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameGetProviderByName,
		EndpointParameters: map[string]string{
			"providerName": name,
		},
		Payload:      &ProviderResponse{},
		ErrorPayload: new(integration.ErrorDTO),
	})
}

func (i *Integration) GetProviderExtensionByName(name string) (*integration.MsxResponse, error) {
	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameGetProviderExtensionByName,
		EndpointParameters: map[string]string{
			"providerExtensionName": name,
		},
		Payload:      &ProviderExtensionResponse{},
		ErrorPayload: new(integration.ErrorDTO),
	})
}

// Deprecated: The underlying REST API was deprecated in 3.10.0.  v8 (or newer) API should be used instead.
func (i *Integration) GetUserById(userId string) (*integration.MsxResponse, error) {
	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameGetUserById,
		EndpointParameters: map[string]string{
			"userId": userId,
		},
		ExpectEnvelope: true,
		Payload:        new(UserResponse),
	})
}

func (i *Integration) GetUserByIdV8(userId string) (*integration.MsxResponse, error) {
	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameGetUserByIdV8,
		EndpointParameters: map[string]string{
			"userId": userId,
		},
		ExpectEnvelope: false,
		Payload:        new(UserResponseV8),
	})
}

// Deprecated: The underlying REST API was deprecated in 3.10.0.  v8 (or newer) API should be used instead.
func (i *Integration) GetTenantById(tenantId string) (result *integration.MsxResponse, err error) {
	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameGetTenantById,
		EndpointParameters: map[string]string{
			"tenantId": tenantId,
		},
		ExpectEnvelope: true,
		Payload:        new(TenantResponse),
	})
}

func (i *Integration) GetTenantByIdV8(tenantId string) (result *integration.MsxResponse, err error) {
	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameGetTenantByIdV8,
		EndpointParameters: map[string]string{
			"tenantId": tenantId,
		},
		ExpectEnvelope: false,
		Payload:        new(TenantResponseV8),
	})
}

// Deprecated: Tenants should generally be access by ID, not tenantName.  The REST API is retired
// and due for decomissioning.
func (i *Integration) GetTenantByName(tenantName string) (result *integration.MsxResponse, err error) {
	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameGetTenantByName,
		EndpointParameters: map[string]string{
			"tenantName": tenantName,
		},
		Payload: new(TenantResponse),
	})
}

func (i *Integration) GetSystemSecrets(scope string) (result *integration.MsxResponse, err error) {
	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameGetSystemSecrets,
		EndpointParameters: map[string]string{
			"scope": scope,
		},
		Payload:        new(SecretsResponse),
		ExpectEnvelope: true,
	})
}

func (i *Integration) AddSystemSecrets(scope string, secrets map[string]string) (result *integration.MsxResponse, err error) {
	var bodyBytes []byte
	if bodyBytes, err = json.Marshal(secrets); err != nil {
		return nil, errors.Wrap(err, "Failed to serialize body")
	}
	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameAddSystemSecrets,
		EndpointParameters: map[string]string{
			"scope": scope,
		},
		Body:           bodyBytes,
		ExpectEnvelope: true,
	})
}

func (i *Integration) ReplaceSystemSecrets(scope string, secrets map[string]string) (result *integration.MsxResponse, err error) {
	var bodyBytes []byte
	if bodyBytes, err = json.Marshal(secrets); err != nil {
		return nil, errors.Wrap(err, "Failed to serialize body")
	}
	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameReplaceSystemSecrets,
		EndpointParameters: map[string]string{
			"scope": scope,
		},
		Body:           bodyBytes,
		ExpectEnvelope: true,
	})
}

func (i *Integration) EncryptSystemSecrets(scope string, names []string, encrypt EncryptSecretsDTO) (result *integration.MsxResponse, err error) {
	body := &GetSecretRequestDTO{
		Names:   names,
		Encrypt: encrypt,
	}

	var bodyBytes []byte
	if bodyBytes, err = json.Marshal(body); err != nil {
		return nil, errors.Wrap(err, "Failed to serialize body")
	}

	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameEncryptSystemSecrets,
		EndpointParameters: map[string]string{
			"scope": scope,
		},
		Body:           bodyBytes,
		Payload:        new(Pojo),
		ExpectEnvelope: true,
	})
}

func (i *Integration) RemoveSystemSecrets(scope string) (result *integration.MsxResponse, err error) {
	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameRemoveSystemSecrets,
		EndpointParameters: map[string]string{
			"scope": scope,
		},
		ExpectEnvelope: true,
	})
}

func (i *Integration) GenerateSystemSecrets(scope string, names []string, save bool) (result *integration.MsxResponse, err error) {
	body := &GenerateSecretRequestDTO{
		Names:   names,
		Save:    save,
		Encrypt: nil,
	}

	var bodyBytes []byte
	if bodyBytes, err = json.Marshal(body); err != nil {
		return nil, errors.Wrap(err, "Failed to serialize body")
	}

	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameGenerateSystemSecrets,
		EndpointParameters: map[string]string{
			"scope": scope,
		},
		Body:           bodyBytes,
		Payload:        new(Pojo),
		ExpectEnvelope: true,
	})
}

func (i *Integration) GetTenantSecrets(tenantId, scope string) (result *integration.MsxResponse, err error) {
	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameGetTenantSecrets,
		EndpointParameters: map[string]string{
			"tenantId": tenantId,
			"scope":    scope,
		},
		Payload:        new(SecretsResponse),
		ExpectEnvelope: true,
	})
}

func (i *Integration) AddTenantSecrets(tenantId, scope string, secrets map[string]string) (result *integration.MsxResponse, err error) {
	var bodyBytes []byte
	if bodyBytes, err = json.Marshal(secrets); err != nil {
		return nil, errors.Wrap(err, "Failed to serialize body")
	}
	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameAddTenantSecrets,
		EndpointParameters: map[string]string{
			"tenantId": tenantId,
			"scope":    scope,
		},
		Body:           bodyBytes,
		ExpectEnvelope: true,
	})
}

func (i *Integration) ReplaceTenantSecrets(tenantId, scope string, secrets map[string]string) (result *integration.MsxResponse, err error) {
	var bodyBytes []byte
	if bodyBytes, err = json.Marshal(secrets); err != nil {
		return nil, errors.Wrap(err, "Failed to serialize body")
	}
	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameReplaceTenantSecrets,
		EndpointParameters: map[string]string{
			"tenantId": tenantId,
			"scope":    scope,
		},
		Body:           bodyBytes,
		ExpectEnvelope: true,
	})
}

func (i *Integration) EncryptTenantSecrets(tenantId, scope string, names []string, encrypt EncryptSecretsDTO) (result *integration.MsxResponse, err error) {
	body := &GetSecretRequestDTO{
		Names:   names,
		Encrypt: encrypt,
	}

	var bodyBytes []byte
	if bodyBytes, err = json.Marshal(body); err != nil {
		return nil, errors.Wrap(err, "Failed to serialize body")
	}

	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameEncryptTenantSecrets,
		EndpointParameters: map[string]string{
			"tenantId": tenantId,
			"scope":    scope,
		},
		Body:           bodyBytes,
		Payload:        new(Pojo),
		ExpectEnvelope: true,
	})
}

func (i *Integration) RemoveTenantSecrets(tenantId, scope string) (result *integration.MsxResponse, err error) {
	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameRemoveTenantSecrets,
		EndpointParameters: map[string]string{
			"tenantId": tenantId,
			"scope":    scope,
		},
		ExpectEnvelope: true,
	})
}

func (i *Integration) GenerateTenantSecrets(tenantId, scope string, names []string, save bool) (result *integration.MsxResponse, err error) {
	body := &GenerateSecretRequestDTO{
		Names:   names,
		Save:    save,
		Encrypt: nil,
	}

	var bodyBytes []byte
	if bodyBytes, err = json.Marshal(body); err != nil {
		return nil, errors.Wrap(err, "Failed to serialize body")
	}

	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameGenerateTenantSecrets,
		EndpointParameters: map[string]string{
			"tenantId": tenantId,
			"scope":    scope,
		},
		Body:           bodyBytes,
		Payload:        new(Pojo),
		ExpectEnvelope: true,
	})
}

func (i *Integration) GetRoles(resolvePermissionNames bool, p paging.Request) (result *integration.MsxResponse, err error) {
	qp := p.QueryParameters()
	qp.Set("resolvepermissionname", strconv.FormatBool(resolvePermissionNames))

	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName:    endpointNameGetRoles,
		QueryParameters: qp,
		Payload:         new(RoleListResponse),
		ExpectEnvelope:  false,
	})
}

func (i *Integration) CreateRole(populator bool, body RoleCreateRequest) (result *integration.MsxResponse, err error) {
	var bodyBytes []byte
	if bodyBytes, err = json.Marshal(body); err != nil {
		return nil, errors.Wrap(err, "Failed to serialize body")
	}

	qp := url.Values{}
	if body.Owner != "" {
		qp.Set("owner", body.Owner)
	}
	if populator {
		qp.Set("dbinstaller", strconv.FormatBool(populator))
	}

	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName:    endpointNameCreateRole,
		QueryParameters: qp,
		Body:            bodyBytes,
		Payload:         new(RoleResponse),
		ExpectEnvelope:  false,
	})
}

func (i *Integration) UpdateRole(populator bool, body RoleUpdateRequest) (result *integration.MsxResponse, err error) {
	var bodyBytes []byte
	if bodyBytes, err = json.Marshal(body); err != nil {
		return nil, errors.Wrap(err, "Failed to serialize body")
	}

	qp := url.Values{}
	if body.Owner != "" {
		qp.Set("owner", body.Owner)
	}
	if populator {
		qp.Set("dbinstaller", strconv.FormatBool(populator))
	}

	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameUpdateRole,
		EndpointParameters: map[string]string{
			"roleName": body.RoleName,
		},
		QueryParameters: qp,
		Body:            bodyBytes,
		Payload:         new(RoleResponse),
		ExpectEnvelope:  false,
	})
}

func (i *Integration) DeleteRole(roleName string) (result *integration.MsxResponse, err error) {
	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameDeleteRole,
		EndpointParameters: map[string]string{
			"roleName": roleName,
		},
		Payload:        new(RoleResponse),
		ExpectEnvelope: false,
	})
}

func (i *Integration) GetCapabilities(p paging.Request) (*integration.MsxResponse, error) {
	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName:    endpointNameGetCapabilities,
		QueryParameters: p.QueryParameters(),
		Payload:         new(CapabilityListResponse),
		ExpectEnvelope:  false,
	})
}

func (i *Integration) BatchCreateCapabilities(populator bool, owner string, capabilities []CapabilityCreateRequest) (result *integration.MsxResponse, err error) {
	var bodyBytes []byte
	if bodyBytes, err = json.Marshal(CapabilityBatchCreateRequest{capabilities}); err != nil {
		return nil, errors.Wrap(err, "Failed to serialize body")
	}

	qp := url.Values{}
	if owner != "" {
		qp.Set("owner", owner)
	}
	if populator {
		qp.Set("dbinstaller", strconv.FormatBool(populator))
	}

	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName:    endpointNameBatchCreateCapabilities,
		QueryParameters: qp,
		Body:            bodyBytes,
		Payload:         new(CapabilityListResponse),
		ExpectEnvelope:  false,
	})
}

func (i *Integration) BatchUpdateCapabilities(populator bool, owner string, capabilities []CapabilityUpdateRequest) (result *integration.MsxResponse, err error) {
	var bodyBytes []byte
	if bodyBytes, err = json.Marshal(CapabilityBatchUpdateRequest{capabilities}); err != nil {
		return nil, errors.Wrap(err, "Failed to serialize body")
	}

	qp := url.Values{}
	if owner != "" {
		qp.Set("owner", owner)
	}
	if populator {
		qp.Set("dbinstaller", strconv.FormatBool(populator))
	}

	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName:    endpointNameBatchUpdateCapabilities,
		QueryParameters: qp,
		Body:            bodyBytes,
		Payload:         new(CapabilityListResponse),
		ExpectEnvelope:  false,
	})
}

func (i *Integration) DeleteCapability(populator bool, owner string, name string) (*integration.MsxResponse, error) {
	qp := url.Values{}
	if owner != "" {
		qp.Set("owner", owner)
	}
	if populator {
		qp.Set("dbinstaller", strconv.FormatBool(populator))
	}

	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameDeleteCapability,
		EndpointParameters: map[string]string{
			"capabilityName": name,
		},
		QueryParameters: qp,
		Payload:         new(CapabilityResponse),
		ExpectEnvelope:  false,
	})
}

func (i *Integration) GetSecretPolicy(name string) (*integration.MsxResponse, error) {
	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameGetSecretPolicy,
		EndpointParameters: map[string]string{
			"policyName": name,
		},
		Payload:        new(SecretPolicyResponse),
		ExpectEnvelope: true,
	})
}

func (i *Integration) StoreSecretPolicy(name string, policy SecretPolicySetRequest) (result *integration.MsxResponse, err error) {
	var body struct {
		SecretPolicySetRequest
		Name string `json:"name"`
	}

	body.SecretPolicySetRequest = policy
	body.Name = name

	var bodyBytes []byte
	if bodyBytes, err = json.Marshal(body); err != nil {
		return nil, errors.Wrap(err, "Failed to serialize body")
	}

	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameSetSecretPolicy,
		EndpointParameters: map[string]string{
			"policyName": name,
		},
		Body:           bodyBytes,
		Payload:        new(SecretPolicyResponse),
		ExpectEnvelope: true,
	})
}

func (i *Integration) DeleteSecretPolicy(name string) (*integration.MsxResponse, error) {
	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameUnsetSecretPolicy,
		EndpointParameters: map[string]string{
			"policyName": name,
		},
		ExpectEnvelope: true,
	})
}

func (i *Integration) GetTenantHierarchyRoot() (*integration.MsxResponse, error) {
	msxEndpointRequest, err := i.buildTenantHierarchyMsxEndpointRequest(endpointTenantHierarchyRoot)
	if err != nil {
		return nil, err
	}

	return i.Execute(msxEndpointRequest)
}

func (i *Integration) GetTenantHierarchyParent(tenantId types.UUID) (*integration.MsxResponse, error) {
	msxEndpointRequest, err := i.buildTenantHierarchyMsxEndpointRequest(endpointTenantHierarchyParent)
	if err != nil {
		return nil, err
	}

	qp := url.Values{}
	qp.Set("tenantId", tenantId.String())
	msxEndpointRequest.QueryParameters = qp

	return i.Execute(msxEndpointRequest)
}

func (i *Integration) GetTenantHierarchyAncestors(tenantId types.UUID) (*integration.MsxResponse, []types.UUID, error) {
	request, err := i.buildTenantHierarchyMsxEndpointRequest(endpointTenantHierarchyAncestors)
	if err != nil {
		return nil, nil, err
	}

	request.QueryParameters = url.Values{
		"tenantId": []string{tenantId.String()},
	}

	response, err := i.Execute(request)
	if err != nil {
		return nil, nil, err
	}

	var result []types.UUID
	if err = json.Unmarshal(response.Body, &result); err != nil {
		return nil, nil, err
	}

	return response, result, err
}

func (i *Integration) buildTenantHierarchyMsxEndpointRequest(endpointName string) (*integration.MsxEndpointRequest, error) {

	securityClientSettings, err := integration.NewSecurityClientSettings(i.Context())
	if err != nil {
		return nil, err
	}

	var headers = make(http.Header)
	headers.Set("Authorization", securityClientSettings.Authorization())
	headers.Set("Content-Type", httpclient.MimeTypeApplicationWwwFormUrlencoded)
	headers.Set("Accept", httpclient.MimeTypeApplicationJson)

	return &integration.MsxEndpointRequest{
		EndpointName:   endpointName,
		Headers:        headers,
		ExpectEnvelope: false,
		NoToken:        true,
		Payload:        nil,
		ErrorPayload:   new(integration.ErrorDTO3),
	}, nil
}
