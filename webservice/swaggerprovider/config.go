package swaggerprovider

import (
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"strings"
)

const (
	configRootDocumentation = "swagger"
)

type DocumentationConfig struct {
	Enabled     bool   `config:"default=false"`
	ApiPath     string `config:"default=/apidocs.json"`
	SwaggerPath string `config:"default=/swagger-resources"`
	Security    struct {
		Sso struct {
			BaseUrl       string `config:"default=http://localhost:9103/idm"`
			TokenPath     string `config:"default=/v2/token"`
			AuthorizePath string `config:"default=/v2/authorize"`
			ClientId      string `config:"default=nfv-client"`
			ClientSecret  string `config:"default=nfv-secret"`
		}
	}
	Ui struct {
		Enabled  bool   `config:"default=true"`
		Endpoint string `config:"default=/swagger"`
		View     string `config:"default=/swagger-ui.html"`
	}
}

func DocumentationConfigFromConfig(cfg *config.Config) (*DocumentationConfig, error) {
	var documentationConfig DocumentationConfig
	if err := cfg.Populate(&documentationConfig, configRootDocumentation); err != nil {
		return nil, err
	}

	if !strings.HasPrefix(documentationConfig.SwaggerPath, "/") {
		documentationConfig.SwaggerPath = "/" + documentationConfig.SwaggerPath
	}

	if strings.HasSuffix(documentationConfig.SwaggerPath, "/") {
		documentationConfig.SwaggerPath = strings.TrimSuffix(documentationConfig.SwaggerPath, "/")
	}

	return &documentationConfig, nil
}
