package metricsprovider

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/webservice"
	"cto-github.cisco.com/NFV-BU/go-msx/webservice/adminprovider"
	"github.com/emicklei/go-restful"
)

const (
	endpointName = "metrics"
)

var metricNames = []string{}

type Report struct {
	Names []string `json:"names"`
}

type Provider struct{}

func (h Provider) Report(req *restful.Request) (interface{}, error) {
	return Report{Names: metricNames}, nil
}

func (h Provider) Actuate(webService *restful.WebService) error {
	webService.Consumes(restful.MIME_JSON, restful.MIME_XML)
	webService.Produces(restful.MIME_JSON, restful.MIME_XML)

	webService.Path(webService.RootPath() + "/admin/" + endpointName)

	// Unsecured routes for info
	webService.Route(webService.GET("").
		Operation("admin.metrics").
		To(webservice.RawController(h.Report)).
		Do(webservice.Returns200))

	return nil
}

func RegisterProvider(ctx context.Context) error {
	server := webservice.WebServerFromContext(ctx)
	if server != nil {
		server.RegisterActuator(new(Provider))
		adminprovider.RegisterLink(endpointName, endpointName, false)
	}
	return nil
}
