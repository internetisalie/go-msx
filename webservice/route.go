package webservice

import (
	"cto-github.cisco.com/NFV-BU/go-msx/integration"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"cto-github.cisco.com/NFV-BU/go-msx/rbac"
	"cto-github.cisco.com/NFV-BU/go-msx/security"
	"cto-github.cisco.com/NFV-BU/go-msx/security/httprequest"
	"fmt"
	"github.com/emicklei/go-restful"
	restfulspec "github.com/emicklei/go-restful-openapi"
	"github.com/go-openapi/spec"
	"github.com/pkg/errors"
	"net/http"
	"path"
	"reflect"
	"strings"
)

const (
	HeaderNameAuthorization     = "Authorization"
	MetadataKeyResponseEnvelope = "MSX_RESPONSE_ENVELOPE"
	MetadataKeyResponsePayload  = "MSX_RESPONSE_PAYLOAD"
	MetadataTagDefinition       = "TagDefinition"
	requestAttributeParams      = "params"
)

var (
	HeaderAuthorization *restful.Parameter
	logger              = log.NewLogger("msx.webservice")
	responseTypes       = make(map[reflect.Type]string)
)

func init() {
	HeaderAuthorization = restful.
		HeaderParameter(HeaderNameAuthorization, "Authentication token in form 'Bearer {token}'").
		Required(false)
}

func StandardList(b *restful.RouteBuilder) {
	b.Do(StandardReturns, ProducesJson)
}

func StandardRetrieve(b *restful.RouteBuilder) {
	b.Do(StandardReturns, ProducesJson)
}

func StandardCreate(b *restful.RouteBuilder) {
	b.Do(CreateReturns, ProducesJson, ConsumesJson)
}

func StandardUpdate(b *restful.RouteBuilder) {
	b.Do(StandardReturns, ProducesJson, ConsumesJson)
}

func StandardDelete(b *restful.RouteBuilder) {
	b.Do(StandardReturns, ProducesJson)
}

func ResponseTypeName(t reflect.Type) (string, bool) {
	typeName, ok := responseTypes[t]
	return typeName, ok
}

func newResponse(payload interface{}) interface{} {
	structType := reflect.TypeOf(integration.MsxEnvelope{})
	var structFields []reflect.StructField
	for i := 0; i < structType.NumField(); i++ {
		structField := structType.Field(i)
		if structField.Name == "Payload" {
			structField.Type = reflect.TypeOf(payload)
		}
		structFields = append(structFields, structField)
	}

	payloadType := reflect.TypeOf(payload)
	if payloadType.Kind() == reflect.Ptr {
		payloadType = payloadType.Elem()
	}
	payloadPackageName := path.Base(payloadType.PkgPath())
	if payloadPackageName != "" {
		payloadPackageName += "."
	}
	payloadTypeName := payloadPackageName + payloadType.Name()
	responseTypeName := fmt.Sprintf("integration.MsxEnvelope«%s»", payloadTypeName)
	responseType := reflect.StructOf(structFields)
	responseTypes[responseType] = responseTypeName
	return reflect.New(responseType).Interface()
}

func ResponsePayload(payload interface{}) func(*restful.RouteBuilder) {
	return func(b *restful.RouteBuilder) {
		example := newResponse(payload)
		b.DefaultReturns("Success", example)
		b.Writes(example)
	}
}

func ResponseRawPayload(payload interface{}) func(*restful.RouteBuilder) {
	return func(b *restful.RouteBuilder) {
		b.DefaultReturns("Success", payload)
		b.Writes(payload)
	}
}

func StandardReturns(b *restful.RouteBuilder) {
	b.Do(Returns(200, 400, 401, 403))
}

func CreateReturns(b *restful.RouteBuilder) {
	b.Do(Returns(200, 201, 400, 401, 403))
}

func ProducesJson(b *restful.RouteBuilder) {
	b.Produces(MIME_JSON)
}

func ConsumesJson(b *restful.RouteBuilder) {
	b.Consumes(MIME_JSON)
}

func securityContextFilter(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	token, err := httprequest.ExtractToken(req.Request)
	if err != nil && err != httprequest.ErrNotFound {
		WriteErrorEnvelope(req, resp, http.StatusBadRequest, err)
		return
	}

	if err == nil {
		userContext, err := security.NewUserContextFromToken(req.Request.Context(), token)
		if err != nil {
			WriteErrorEnvelope(req, resp, http.StatusBadRequest, err)
		}

		ctx := security.ContextWithUserContext(req.Request.Context(), userContext)
		req.Request = req.Request.WithContext(ctx)
	}

	chain.ProcessFilter(req, resp)
}

func authenticationFilter(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	authenticationProvider := AuthenticationProviderFromContext(req.Request.Context())
	if authenticationProvider != nil {
		err := authenticationProvider.Authenticate(req)
		if err != nil {
			WriteErrorEnvelope(req, resp, http.StatusUnauthorized, err)
			return
		}
	}

	chain.ProcessFilter(req, resp)
}

func PermissionsFilter(anyOf ...string) restful.FilterFunction {
	return func(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
		var ctx = req.Request.Context()
		if err := rbac.HasPermission(ctx, anyOf); err != nil {
			logger.WithError(err).WithField("perms", anyOf).Error("Permission denied")
			WriteErrorEnvelope(req, resp, http.StatusForbidden, err)
			return
		}

		chain.ProcessFilter(req, resp)
	}
}

func getParameter(parameter *restful.Parameter, req *restful.Request) (string, error) {
	switch parameter.Kind() {
	case restful.PathParameterKind:
		return req.PathParameter(parameter.Data().Name), nil

	case restful.BodyParameterKind:
		return req.BodyParameter(parameter.Data().Name)

	case restful.QueryParameterKind:
		return req.QueryParameter(parameter.Data().Name), nil

	case restful.HeaderParameterKind:
		return req.HeaderParameter(parameter.Data().Name), nil

	default:
		return "", errors.Errorf("Unsupported parameter type: %v", parameter.Kind())
	}
}

func TenantFilter(parameter *restful.Parameter) restful.FilterFunction {
	return func(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
		tenantId, err := getParameter(parameter, req)
		if err != nil {
			WriteErrorEnvelope(req, resp, http.StatusBadRequest, err)
			return
		}

		ctx := req.Request.Context()
		if err := rbac.HasTenant(ctx, tenantId); err != nil {
			ctx = log.ExtendContext(ctx, log.LogContext{
				"tenant": tenantId,
			})
			req.Request = req.Request.WithContext(ctx)
			WriteErrorEnvelope(req, resp, http.StatusForbidden, err)
			return
		}

		chain.ProcessFilter(req, resp)
	}
}

func optionsFilter(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	if "OPTIONS" != req.Request.Method {
		chain.ProcessFilter(req, resp)
		return
	}

	var container = ContainerFromContext(req.Request.Context())
	var router = RouterFromContext(req.Request.Context())
	var newHttpRequest = *req.Request
	var allowedMethods = make(map[string]struct{})
	for _, method := range []string{"PATCH", "POST", "GET", "PUT", "DELETE", "HEAD", "TRACE"} {
		newHttpRequest.Method = method
		_, route, err := router.SelectRoute(container.RegisteredWebServices(), &newHttpRequest)
		if err != nil || route == nil {
			continue
		}
		allowedMethods[route.Method] = struct{}{}
	}

	if len(allowedMethods) == 0 {
		http.NotFound(resp, req.Request)
		return
	}

	allowedMethods["OPTIONS"] = struct{}{}
	var allowMethods []string
	for k := range allowedMethods {
		allowMethods = append(allowMethods, k)
	}

	resp.AddHeader("Allow", strings.Join(allowMethods, ","))
}

var DefaultSuccessEnvelope = integration.MsxEnvelope{}

type RouteBuilderFunc func(*restful.RouteBuilder)

func Returns200(b *restful.RouteBuilder) {
	b.Returns(http.StatusOK, "OK", nil)
}
func Returns201(b *restful.RouteBuilder) {
	b.Returns(http.StatusCreated, "Created", nil)
}
func Returns204(b *restful.RouteBuilder) {
	b.Returns(http.StatusNoContent, "No Content", nil)
}
func Returns400(b *restful.RouteBuilder) {
	b.Returns(http.StatusBadRequest, "Bad Request", nil)
}
func Returns401(b *restful.RouteBuilder) {
	b.Returns(http.StatusUnauthorized, "Not Authorized", nil)
}
func Returns403(b *restful.RouteBuilder) {
	b.Returns(http.StatusForbidden, "Forbidden", nil)
}
func Returns404(b *restful.RouteBuilder) {
	b.Returns(http.StatusNotFound, "Not Found", nil)
}
func Returns409(b *restful.RouteBuilder) {
	b.Returns(http.StatusConflict, "Conflict", nil)
}
func Returns424(b *restful.RouteBuilder) {
	b.Returns(http.StatusFailedDependency, "Failed Dependency", nil)
}
func Returns500(b *restful.RouteBuilder) {
	b.Returns(http.StatusInternalServerError, "Internal Server Error", nil)
}
func Returns503(b *restful.RouteBuilder) {
	b.Returns(http.StatusInternalServerError, "Bad Gateway", nil)
}

func Returns(statuses ...int) RouteBuilderFunc {
	var statusFuncs []RouteBuilderFunc
	for _, status := range statuses {
		switch status {
		case 200:
			statusFuncs = append(statusFuncs, Returns200)
		case 201:
			statusFuncs = append(statusFuncs, Returns201)
		case 204:
			statusFuncs = append(statusFuncs, Returns204)
		case 400:
			statusFuncs = append(statusFuncs, Returns400)
		case 401:
			statusFuncs = append(statusFuncs, Returns401)
		case 404:
			statusFuncs = append(statusFuncs, Returns404)
		case 409:
			statusFuncs = append(statusFuncs, Returns409)
		case 424:
			statusFuncs = append(statusFuncs, Returns424)
		case 500:
			statusFuncs = append(statusFuncs, Returns500)
		case 503:
			statusFuncs = append(statusFuncs, Returns503)
		}
	}

	return func(b *restful.RouteBuilder) {
		for _, statusFunc := range statusFuncs {
			statusFunc(b)
		}
	}
}

func TagDefinition(name, description string) RouteBuilderFunc {
	return func(b *restful.RouteBuilder) {
		b.Metadata(restfulspec.KeyOpenAPITags, []string{name})
		b.Metadata(MetadataTagDefinition, spec.TagProps{
			Name: name,
			Description:description,
		})
	}
}

type RouteFunction func(svc *restful.WebService) *restful.RouteBuilder

func PopulateParams(template interface{}) RouteBuilderFunc {
	templateType := reflect.TypeOf(template)
	if templateType.Kind() == reflect.Ptr {
		templateType = templateType.Elem()
	}

	return func(builder *restful.RouteBuilder) {
		builder.Filter(func(req *restful.Request, response *restful.Response, chain *restful.FilterChain) {
			// Instantiate a new object of the same type as template
			target := reflect.New(templateType).Interface()

			// Populate the target
			if err := Populate(req, target); err != nil {
				WriteErrorEnvelope(req, response, 400, err)
				return
			}

			req.SetAttribute(requestAttributeParams, target)

			chain.ProcessFilter(req, response)
		})
	}
}

func Params(req *restful.Request) interface{} {
	return req.Attribute(requestAttributeParams)
}
