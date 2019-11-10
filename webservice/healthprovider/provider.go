package healthprovider

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/health"
	"cto-github.cisco.com/NFV-BU/go-msx/security"
	"cto-github.cisco.com/NFV-BU/go-msx/webservice"
	"cto-github.cisco.com/NFV-BU/go-msx/webservice/adminprovider"
	"github.com/emicklei/go-restful"
	"github.com/pkg/errors"
	"net/http"
)

type HealthProvider struct{}


func (h HealthProvider) healthReport(ctx context.Context) (interface{}, error) {
	userContext := security.UserContextFromContext(ctx)
	if userContext != nil {
		return health.GenerateReport(ctx), nil
	} else {
		return health.GenerateSummary(ctx), nil
	}
}

func (h HealthProvider) healthComponentReport(req *restful.Request) (interface{}, error) {
	component := req.PathParameter("component")
	report := health.GenerateReport(req.Request.Context())
	if details, ok := report.Details[component]; !ok {
		return nil, webservice.NewStatusError(
			errors.New("Component not found"),
			http.StatusNotFound)
	} else {
		return details, nil
	}
}

func (h HealthProvider) Actuate(healthService *restful.WebService) error {
	healthService.Consumes(restful.MIME_JSON, restful.MIME_XML)
	healthService.Produces(restful.MIME_JSON, restful.MIME_XML)

	healthService.Path(healthService.RootPath() + "/admin/health")

	healthService.Route(healthService.GET("").
		To(webservice.RawContextController(h.healthReport)).
		Doc("Get System health").
		Do(webservice.Returns200))

	healthService.Route(healthService.GET("/{component}").
		To(webservice.RawController(h.healthComponentReport)).
		Param(healthService.PathParameter("component", "Name of component to probe")).
		Doc("Get component health").
		Do(webservice.Returns(200, 404)))

	return nil
}

func RegisterProvider(ctx context.Context) error {
	server := webservice.WebServerFromContext(ctx)
	if server != nil {
		server.RegisterActuator(new(HealthProvider))
		adminprovider.RegisterLink("health", "health", false)
		adminprovider.RegisterLink("health-component", "health/{component}", true)
	}
	return nil
}
