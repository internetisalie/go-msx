// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

//go:generate mockery --name RestController --structname MockRestController --filename mock_RestController.go --inpackage

package webservice

import (
	"cto-github.cisco.com/NFV-BU/go-msx/audit/auditlog"
	"cto-github.cisco.com/NFV-BU/go-msx/integration"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"cto-github.cisco.com/NFV-BU/go-msx/paging"
	"cto-github.cisco.com/NFV-BU/go-msx/rbac"
	"cto-github.cisco.com/NFV-BU/go-msx/schema"
	"cto-github.cisco.com/NFV-BU/go-msx/security"
	"cto-github.cisco.com/NFV-BU/go-msx/security/certdetailsprovider"
	"cto-github.cisco.com/NFV-BU/go-msx/security/httprequest"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"cto-github.cisco.com/NFV-BU/go-msx/validate"
	"cto-github.cisco.com/NFV-BU/go-msx/webservice/restfulcontext"
	"github.com/emicklei/go-restful"
	restfulspec "github.com/emicklei/go-restful-openapi"
	"github.com/go-openapi/spec"
	"github.com/iancoleman/strcase"
	"github.com/pkg/errors"
	"github.com/swaggest/refl"
	"net/http"
	"reflect"
	"strings"
)

const (
	HeaderNameAuthorization   = "Authorization"
	HeaderNameContentEncoding = "Content-Encoding"
	HeaderNameContentType     = "Content-Type"

	MetadataTagDefinition     = "TagDefinition"
	MetadataPermissions       = "Permissions"
	MetadataSuccessResponse   = "SuccessResponse"
	MetadataErrorPayload      = "ErrorPayload"
	MetadataRequest           = "Request"
	MetadataDefaultReturnCode = "DefaultReturnCode"
	MetadataEnvelope          = "Envelope"
	MetadataParams            = "InputParamsStruct"
	MetadataApiIgnore         = "ApiIgnore"

	AttributeDefaultReturnCode = "DefaultReturnCode"
	AttributeSuccessResponse   = "SuccessResponse"
	AttributeErrorPayload      = "ErrorPayload"
	AttributeError             = "Error"
	AttributeParams            = "Params"
	AttributeSilenceLog        = "ManagementSilenceLog"

	StructTagRequest  = "req"
	StructTagResponse = "resp"
)

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

func StandardAccept(b *restful.RouteBuilder) {
	b.Do(AcceptReturns, ProducesJson, ConsumesJson)
}

func StandardNoContent(b *restful.RouteBuilder) {
	b.Do(NoContentReturns, ProducesJson, ConsumesJson)
}

func NoBodyNoContent(b *restful.RouteBuilder) {
	b.Do(NoContentReturns, ProducesJson, ConsumesAny)
}

func NewEnvelopedResponse(envelopeType reflect.Type, payloadInstance interface{}) interface{} {
	responseType := schema.NewParameterizedStruct(
		envelopeType,
		payloadInstance)

	return reflect.New(responseType).Interface()
}

func ResponsePayload(payload interface{}) func(*restful.RouteBuilder) {
	errorPayloadFn := ErrorPayload(new(integration.MsxEnvelope))
	return func(b *restful.RouteBuilder) {
		example := NewEnvelopedResponse(
			reflect.TypeOf(integration.MsxEnvelope{}),
			payload)
		RouteBuilderWithSuccessResponse(b, example)
		RouteBuilderWithEnvelopedPayload(b, payload)
		b.DefaultReturns("Success", example)
		b.Writes(example)
		b.Do(errorPayloadFn)
	}
}

func PaginatedResponsePayload(payload interface{}) func(*restful.RouteBuilder) {
	errorPayloadFn := ErrorPayload(new(integration.MsxEnvelope))
	return func(b *restful.RouteBuilder) {
		paginatedPayload := NewEnvelopedResponse(
			reflect.TypeOf(paging.PaginatedResponse{}),
			payload)
		envelopedPayload := NewEnvelopedResponse(
			reflect.TypeOf(integration.MsxEnvelope{}),
			paginatedPayload)
		RouteBuilderWithSuccessResponse(b, envelopedPayload)
		RouteBuilderWithEnvelopedPayload(b, payload)
		b.DefaultReturns("Success", envelopedPayload)
		b.Writes(envelopedPayload)
		b.Do(errorPayloadFn)
	}
}

func PaginatedV8ResponsePayload(payload interface{}) func(*restful.RouteBuilder) {
	errorPayloadFn := ErrorPayload(new(integration.MsxEnvelope))
	return func(b *restful.RouteBuilder) {
		paginatedPayload := NewEnvelopedResponse(
			reflect.TypeOf(paging.PaginatedResponseV8{}),
			payload)
		envelopedPayload := NewEnvelopedResponse(
			reflect.TypeOf(integration.MsxEnvelope{}),
			paginatedPayload)
		RouteBuilderWithSuccessResponse(b, envelopedPayload)
		RouteBuilderWithEnvelopedPayload(b, paginatedPayload)
		b.DefaultReturns("Success", envelopedPayload)
		b.Writes(envelopedPayload)
		b.Do(errorPayloadFn)
	}
}

func ResponseRawPayload(payload interface{}) func(*restful.RouteBuilder) {
	errorPayloadFn := ErrorPayload(new(integration.ErrorDTO))
	return func(b *restful.RouteBuilder) {
		RouteBuilderWithSuccessResponse(b, payload)
		b.DefaultReturns("Success", payload)
		if payload != nil {
			b.Writes(payload)
		}
		b.Do(errorPayloadFn)
	}
}

func ErrorPayload(payload interface{}) func(*restful.RouteBuilder) {
	return func(builder *restful.RouteBuilder) {
		RouteBuilderWithErrorPayload(builder, payload)
		builder.Filter(func(request *restful.Request, response *restful.Response, chain *restful.FilterChain) {
			RequestWithErrorPayload(request, payload)
			chain.ProcessFilter(request, response)
		})
	}
}

func SuccessResponse(template interface{}) restfulcontext.RouteBuilderFunc {
	templateType := reflect.TypeOf(template)
	for templateType.Kind() == reflect.Ptr {
		templateType = templateType.Elem()
	}

	// Create an instance for the metadata (swagger)
	template = reflect.Zero(templateType).Interface()

	// Extract the payload
	var payload = samplePayload(template, StructTagResponse)

	return func(builder *restful.RouteBuilder) {
		RouteBuilderWithSuccessResponse(builder, template)
		builder.DefaultReturns("Success", payload)
		if payload != nil {
			builder.Writes(payload)
		}

		builder.Filter(func(req *restful.Request, response *restful.Response, chain *restful.FilterChain) {
			// Instantiate a new object of the same type as template for the request
			target := reflect.New(templateType).Interface()

			// Store it in the request attributes
			req.SetAttribute(AttributeSuccessResponse, target)

			// Continue processing
			chain.ProcessFilter(req, response)
		})
	}
}

func samplePayload(template interface{}, structTagName string) interface{} {
	hasReq := false
	templateType := reflect.TypeOf(template)

	if templateType.Kind() == reflect.Struct {
		// search for `req:"body"` in struct tags
		for i := 0; i < templateType.NumField(); i++ {
			templateField := templateType.Field(i)
			req := templateField.Tag.Get(structTagName)
			if req != "" {
				hasReq = true
			}
			if req != "body" {
				continue
			}

			payloadType := templateField.Type
			for payloadType.Kind() == reflect.Ptr {
				payloadType = payloadType.Elem()
			}

			return types.Instantiate(payloadType)
		}
	}

	if !hasReq {
		// Assume the template is the payload
		return types.Instantiate(templateType)
	}

	return nil
}

func StandardReturns(b *restful.RouteBuilder) {
	b.Do(Returns(200, 400, 401, 403), DefaultReturns(200))
}

func CreateReturns(b *restful.RouteBuilder) {
	b.Do(Returns(200, 201, 400, 401, 403), DefaultReturns(201))
}

func AcceptReturns(b *restful.RouteBuilder) {
	b.Do(Returns(202, 400, 401, 403), DefaultReturns(202))
}

func NoContentReturns(b *restful.RouteBuilder) {
	b.Do(Returns(204, 400, 401, 403), DefaultReturns(204))
}

func ProducesJson(b *restful.RouteBuilder) {
	b.Produces(MIME_JSON)
}

func ConsumesJson(b *restful.RouteBuilder) {
	b.Consumes(MIME_JSON)
}

func ConsumesAny(b *restful.RouteBuilder) {
	b.Consumes("*/*")
}

func ConsumesNone(b *restful.RouteBuilder) {
	b.Consumes("")
}

func ProducesTextPlain(b *restful.RouteBuilder) {
	b.Produces(MIME_TEXT_PLAIN)
}

func ConsumesTextPlain(b *restful.RouteBuilder) {
	b.Consumes(MIME_TEXT_PLAIN)
}

func DefaultReturns(code int) restfulcontext.RouteBuilderFunc {
	return func(b *restful.RouteBuilder) {
		RouteBuilderWithDefaultReturnCode(b, code)
		b.Filter(func(request *restful.Request, response *restful.Response, chain *restful.FilterChain) {
			request.SetAttribute(AttributeDefaultReturnCode, code)
			chain.ProcessFilter(request, response)
		})
	}
}

func logSilenceFilter(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	RequestWithSilenceLog(req)
	chain.ProcessFilter(req, resp)
}

func tokenUserContextFilter(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	token, err := httprequest.ExtractToken(req.Request)
	if err != nil && err != httprequest.ErrNotFound {
		WriteError(req, resp, http.StatusUnauthorized, err)
		return
	}

	if err == nil {
		userContext, err := security.NewUserContextFromToken(req.Request.Context(), token)
		if err != nil {
			WriteError(req, resp, http.StatusUnauthorized, err)
			return
		}

		ctx := security.ContextWithUserContext(req.Request.Context(), userContext)
		req.Request = req.Request.WithContext(ctx)
	}

	chain.ProcessFilter(req, resp)
}

func certificateUserContextFilter(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	cert, err := httprequest.ExtractCertificate(req.Request)
	if err != nil && err != httprequest.ErrNotFound {
		WriteError(req, resp, http.StatusUnauthorized, err)
		return
	}

	if err == nil {
		// Found a certificate
		userContext, err := security.NewUserContextFromCertificate(req.Request.Context(), cert)
		if err != nil {
			WriteError(req, resp, http.StatusUnauthorized, err)
			return
		}

		// Inject the derived UserContext
		ctx := security.ContextWithUserContext(req.Request.Context(), userContext)

		// Make sure we answer UserContextDetails queries from the certificate
		ctx = security.ContextWithTokenDetailsProvider(ctx, new(certdetailsprovider.TokenDetailsProvider))

		// Apply the updated context back to the request
		req.Request = req.Request.WithContext(ctx)
	}

	chain.ProcessFilter(req, resp)
}

func filterFilter(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	filters := FiltersFromContext(req.Request.Context())
	if filters != nil {
		chain.Filters = append(chain.Filters, filters...)
	}
	chain.ProcessFilter(req, resp)
}

func authenticationFilter(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	authenticationProvider := AuthenticationProviderFromContext(req.Request.Context())
	if authenticationProvider != nil {
		err := authenticationProvider.Authenticate(req)
		if err != nil {
			WriteError(req, resp, http.StatusUnauthorized, err)
			return
		}
	}

	chain.ProcessFilter(req, resp)
}

func auditContextFilter(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	server := WebServerFromContext(req.Request.Context())
	auditDetails := auditlog.ExtractRequestDetails(req, server.cfg.Host, server.cfg.Port)
	ctx := auditlog.ContextWithRequestDetails(req.Request.Context(), auditDetails)
	req.Request = req.Request.WithContext(ctx)
	chain.ProcessFilter(req, resp)
}

func Permissions(anyOf ...string) restfulcontext.RouteBuilderFunc {
	return func(b *restful.RouteBuilder) {
		RouteBuilderWithPermissions(b, anyOf)
		b.Filter(PermissionsFilter(anyOf...))
	}
}

// Deprecated
func PermissionsFilter(anyOf ...string) restful.FilterFunction {
	return func(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
		var ctx = req.Request.Context()
		// Temporarily allow system user
		var userContext = security.UserContextFromContext(ctx)
		if userContext.UserName != "system" {
			if err := rbac.HasPermission(ctx, anyOf); err != nil {
				logger.WithError(err).WithField("perms", anyOf).Error("Permission denied")
				WriteError(req, resp, http.StatusForbidden, err)
				return
			}
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
			WriteError(req, resp, http.StatusBadRequest, err)
			return
		}

		tenantUuid, err := types.ParseUUID(tenantId)
		if err != nil {
			WriteError(req, resp, http.StatusBadRequest, err)
		}

		ctx := req.Request.Context()
		if err := rbac.HasTenant(ctx, tenantUuid); err != nil {
			ctx = log.ExtendContext(ctx, log.LogContext{
				"tenant": tenantId,
			})
			req.Request = req.Request.WithContext(ctx)
			WriteError(req, resp, http.StatusForbidden, err)
			return
		}

		chain.ProcessFilter(req, resp)
	}
}

func Returns200(b *restful.RouteBuilder) {
	b.Returns(http.StatusOK, "OK", nil)
}
func Returns201(b *restful.RouteBuilder) {
	b.Returns(http.StatusCreated, "Created", nil)
}
func Returns202(b *restful.RouteBuilder) {
	b.Returns(http.StatusAccepted, "Accepted", nil)
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
func Returns502(b *restful.RouteBuilder) {
	b.Returns(http.StatusBadGateway, "Bad Gateway", nil)
}
func Returns503(b *restful.RouteBuilder) {
	b.Returns(http.StatusServiceUnavailable, "Service Unavailable", nil)
}

func Returns(statuses ...int) restfulcontext.RouteBuilderFunc {
	var statusFuncs []restfulcontext.RouteBuilderFunc
	for _, status := range statuses {
		switch status {
		case 200:
			statusFuncs = append(statusFuncs, Returns200)
		case 201:
			statusFuncs = append(statusFuncs, Returns201)
		case 202:
			statusFuncs = append(statusFuncs, Returns202)
		case 204:
			statusFuncs = append(statusFuncs, Returns204)
		case 400:
			statusFuncs = append(statusFuncs, Returns400)
		case 401:
			statusFuncs = append(statusFuncs, Returns401)
		case 403:
			statusFuncs = append(statusFuncs, Returns403)
		case 404:
			statusFuncs = append(statusFuncs, Returns404)
		case 409:
			statusFuncs = append(statusFuncs, Returns409)
		case 424:
			statusFuncs = append(statusFuncs, Returns424)
		case 500:
			statusFuncs = append(statusFuncs, Returns500)
		case 502:
			statusFuncs = append(statusFuncs, Returns502)
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

func TagDefinition(name, description string) restfulcontext.RouteBuilderFunc {
	return func(b *restful.RouteBuilder) {
		b.Metadata(restfulspec.KeyOpenAPITags, []string{name})
		RouteBuilderWithTagDefinition(b,
			spec.TagProps{
				Name:        name,
				Description: description,
			})
	}
}

func TagsFromRoute(route restful.Route) []string {
	tags := route.Metadata[restfulspec.KeyOpenAPITags]
	if tags == nil {
		return nil
	}
	return tags.([]string)
}

type RestController interface {
	Routes(svc *restful.WebService)
}

type RouteFunction func(svc *restful.WebService) *restful.RouteBuilder

func PopulateParams(inputPort interface{}) restfulcontext.RouteBuilderFunc {
	inputPortType := refl.DeepIndirect(reflect.TypeOf(inputPort))

	return func(builder *restful.RouteBuilder) {
		builder.
			Do(RouteInputParamsInjector(reflect.New(inputPortType).Interface())).
			Filter(func(req *restful.Request, response *restful.Response, chain *restful.FilterChain) {
				// Instantiate a new object of the same type as template
				target := reflect.New(inputPortType).Interface()

				// Populate the target
				if err := Populate(req, target); err != nil {
					WriteError(req, response, 400, err)
					return
				}

				RequestWithParams(req, target)

				chain.ProcessFilter(req, response)
			})
	}
}

func Request(template interface{}) restfulcontext.RouteBuilderFunc {
	populateParamsFunc := PopulateParams(template)

	templateType := reflect.TypeOf(template)
	for templateType.Kind() == reflect.Ptr {
		templateType = templateType.Elem()
	}

	// Create an instance for the metadata (swagger)
	template = reflect.Zero(templateType).Interface()

	// Extract the parameters
	var params = extractRequestParameters(template)

	// Extract the payload
	var payload = samplePayload(template, StructTagRequest)

	return func(builder *restful.RouteBuilder) {
		RouteBuilderWithRequest(builder, template)
		if payload != nil {
			builder.Reads(payload)
		}

		for _, param := range params {
			paramName := param.Data().Name
			existingParam := builder.ParameterNamed(paramName)
			if existingParam == nil {
				builder.Param(param)
			}
		}

		populateParamsFunc(builder)
	}
}

func extractRequestParameters(template interface{}) []*restful.Parameter {
	templateType := reflect.TypeOf(template)
	for templateType.Kind() == reflect.Ptr {
		templateType = templateType.Elem()
	}

	var results []*restful.Parameter

	if templateType.Kind() == reflect.Struct {
		for i := 0; i < templateType.NumField(); i++ {
			if p := extractRequestParameter(templateType.Field(i)); p != nil {
				results = append(results, p)
			}
		}
	}

	return results
}

type ParameterTag struct {
	TagName string
	Source  string
	Name    string
}

func NewParameterTag(field reflect.StructField, tagName string) ParameterTag {
	source := field.Tag.Get(tagName)
	name := ""

	if strings.Contains(source, "=") {
		reqParts := strings.SplitN(source, "=", 2)
		source, name = reqParts[0], reqParts[1]
	}

	if name == "" {
		name = strcase.ToLowerCamel(field.Name)
	}

	return ParameterTag{
		TagName: tagName,
		Source:  source,
		Name:    name,
	}
}

var typesTimeType = reflect.TypeOf(types.Time{})
var typesUuidType = reflect.TypeOf(types.UUID{})
var textUnmarshalerInstance types.TextUnmarshaler
var textUnmarshalerType = reflect.TypeOf(&textUnmarshalerInstance).Elem()

func parameterType(fieldType reflect.Type) (dataType, format string) {
	fieldType = refl.DeepIndirect(fieldType)

	if fieldType == typesTimeType {
		return "string", "date-time"
	} else if fieldType == typesUuidType {
		return "string", "uuid"
	} else if fieldType.Implements(textUnmarshalerType) {
		return "string", ""
	}

	switch fieldType.Kind() {
	case reflect.Struct, reflect.Map:
		return "object", ""
	case reflect.Array, reflect.Slice:
		return "array", ""
	case reflect.Uint, reflect.Uint64, reflect.Int, reflect.Int64:
		return "integer", "int64"
	case reflect.Uint32, reflect.Uint16, reflect.Uint8,
		reflect.Int32, reflect.Int16, reflect.Int8:
		return "integer", "int32"
	case reflect.Float64:
		return "number", "double"
	case reflect.Float32:
		return "number", "float"
	case reflect.Bool:
		return "boolean", ""
	case reflect.String:
		return "string", ""
	default:
		return "string", ""
	}
}

func extractRequestParameter(templateField reflect.StructField) *restful.Parameter {
	parameter := NewParameterTag(templateField, StructTagRequest)
	desc := templateField.Tag.Get("description")
	required := templateField.Tag.Get("required") == "true"
	dataType, format := parameterType(templateField.Type)

	switch parameter.Source {
	case "path":
		return restful.PathParameter(parameter.Name, desc).Required(true).DataType(dataType).DataFormat(format).Description(desc)
	case "query":
		return restful.QueryParameter(parameter.Name, desc).Required(required).DataType(dataType).DataFormat(format).Description(desc)
	case "header":
		return restful.HeaderParameter(parameter.Name, desc).Required(required).DataType(dataType).DataFormat(format).Description(desc)
	case "form":
		return restful.FormParameter(parameter.Name, desc).Required(required).DataType(dataType).DataFormat(format).Description(desc)
	}

	return nil
}

func ValidateParams(fn ValidatorFunction) restfulcontext.RouteBuilderFunc {
	return func(builder *restful.RouteBuilder) {
		builder.Filter(func(req *restful.Request, response *restful.Response, chain *restful.FilterChain) {
			err := validate.Validate(requestValidator{fn: fn, req: req})
			if err != nil {
				WriteError(req, response, 400, err)
				return
			}

			chain.ProcessFilter(req, response)
		})
	}
}

type ValidatorFunction func(req *restful.Request) (err error)

type requestValidator struct {
	req *restful.Request
	fn  ValidatorFunction
}

func (r requestValidator) Validate() error {
	return r.fn(r.req)
}

// Routes adds the routes to the specified webservice after optionally tagging
func Routes(svc *restful.WebService, tag restfulcontext.RouteBuilderFunc, routeFunctions ...RouteFunction) {
	for _, routeFunction := range routeFunctions {
		routeBuilder := routeFunction(svc)
		if tag != nil {
			routeBuilder.Do(tag)
		}
		svc.Route(routeBuilder)

		routes := svc.Routes()
		route := routes[len(routes)-1]

		logger.Infof("Registered route: %s %s", route.Method, route.Path)
	}
}
