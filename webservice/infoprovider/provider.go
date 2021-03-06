// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package infoprovider

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"cto-github.cisco.com/NFV-BU/go-msx/webservice"
	"cto-github.cisco.com/NFV-BU/go-msx/webservice/adminprovider"
	"github.com/emicklei/go-restful"
)

const (
	configKeyInfo = "info"
	endpointName  = "info"
)

type InfoProvider struct{}

func (h InfoProvider) infoReport(req *restful.Request) (interface{}, error) {
	type Info struct {
		App struct {
			Name        string `json:"name"`
			Version     string `json:"version" config:"default="`
			Description string `json:"description"`
			Attributes  struct {
				DisplayName string `json:"displayName"`
				Parent      string `json:"parent"`
				Type        string `json:"type"`
			} `json:"attributes"`
		} `json:"app"`
		Build struct {
			Version       string       `json:"version"`
			BuildNumber   string       `json:"number" config:"buildNumber"`
			BuildDateTime string       `json:"-"`
			Artifact      string       `json:"artifact"`
			Name          string       `json:"name"`
			Time          epochSeconds `json:"time" config:"default=0"`
			Group         string       `json:"group"`
			CommitHash    string       `json:"commitHash" config:"default="`
			DiffHash      string       `json:"diffHash" config:"default="`
		} `json:"build"`
	}

	i := Info{}
	if err := config.MustFromContext(req.Request.Context()).Populate(&i, configKeyInfo); err != nil {
		return nil, webservice.NewStatusError(err, 500)
	}

	i.App.Version = i.Build.Version
	buildTime, err := types.ParseTime(i.Build.BuildDateTime)
	if err == nil {
		i.Build.Time = newEpochSeconds(buildTime.ToTimeTime())
	}

	return i, nil
}

func (h InfoProvider) EndpointName() string {
	return endpointName
}

func (h InfoProvider) Actuate(infoService *restful.WebService) error {
	infoService.Consumes(restful.MIME_JSON, restful.MIME_XML)
	infoService.Produces(restful.MIME_JSON, restful.MIME_XML)

	infoService.Path(infoService.RootPath() + "/admin/info")

	// Unsecured routes for info
	infoService.Route(infoService.GET("").
		Operation("admin.info").
		To(adminprovider.RawAdminController(h.infoReport)).
		Doc("Get System info").
		Do(webservice.Returns200))

	return nil
}

func RegisterProvider(ctx context.Context) error {
	server := webservice.WebServerFromContext(ctx)
	if server != nil {
		server.RegisterActuator(new(InfoProvider))
		adminprovider.RegisterLink("info", "info", false)
	}
	return nil
}
