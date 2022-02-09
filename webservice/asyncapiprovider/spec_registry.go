package asyncapiprovider

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/schema/asyncapi"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"cto-github.cisco.com/NFV-BU/go-msx/webservice"
	"cto-github.cisco.com/NFV-BU/go-msx/webservice/swaggerprovider"
	"encoding/json"
	"fmt"
	"github.com/swaggest/jsonschema-go"
	yaml2 "gopkg.in/yaml.v2"
)

type RegistrySpecProvider struct {
	appInfo *swaggerprovider.AppInfo
}

func (p RegistrySpecProvider) Spec() ([]byte, error) {
	spec := p.RenderSpec()

	specJsonBytes, err := json.Marshal(spec)
	if err != nil {
		return nil, err
	}

	var specYaml = yaml2.MapSlice{}
	err = yaml2.Unmarshal(specJsonBytes, &specYaml)
	if err != nil {
		return nil, err
	}

	specYamlBytes, err := yaml2.Marshal(specYaml)
	if err != nil {
		return nil, err
	}

	return specYamlBytes, nil
}

func (p RegistrySpecProvider) RenderSpec() asyncapi.Spec {
	spec := &asyncapi.Spec{}
	spec.ID = types.NewStringPtr(fmt.Sprintf("uri:%s.cpx.plus.cisco.com", p.appInfo.Name))
	spec.DefaultContentType = types.NewStringPtr(webservice.MIME_JSON)
	spec.WithInfo(p.Info())
	spec.WithServersItem("cpx", p.Server())
	spec.WithChannels(p.Channels())
	spec.ComponentsEns().WithMessages(p.Messages())
	spec.ComponentsEns().WithSchemas(p.Schemas())
	spec.ComponentsEns().WithOperationBindingsItem("cpx", p.Binding())
	spec.ComponentsEns().WithSecuritySchemes(p.SecuritySchemes())
	return *spec
}

func (p RegistrySpecProvider) Info() asyncapi.Info {
	return asyncapi.Info{
		Title: p.appInfo.DisplayName,
		Description: types.NewStringPtr("Kafka Stream documentation for " + p.appInfo.DisplayName + "\n" +
			" \n> " + p.appInfo.Description),
		TermsOfService: types.NewStringPtr("http://www.cisco.com"),
		Contact: &asyncapi.Contact{
			Name:  types.NewStringPtr("Cisco Systems Inc."),
			URL:   types.NewStringPtr("http://www.cisco.com"),
			Email: types.NewStringPtr("somecontact@cisco.com"),
		},
		License: &asyncapi.License{
			Name: "Apache License Version 2.0",
			URL:  types.NewStringPtr("http://www.apache.org/licenses/LICENSE-2.0.html"),
		},
		Version: p.appInfo.Version,
	}
}

func (p RegistrySpecProvider) Server() asyncapi.Server {
	return asyncapi.Server{
		URL:             types.NewStringPtr("kafka"),
		Description:     types.NewStringPtr("CPX Internal Kafka"),
		Protocol:        types.NewStringPtr("kafka-secure"),
		ProtocolVersion: types.NewStringPtr("2.2.0"),
		Security: []map[string][]string{
			{"cpx": {}},
		},
	}
}

func (p RegistrySpecProvider) Channels() map[string]asyncapi.ChannelItem {
	return asyncapi.Reflector.SpecEns().Channels
}

func (p RegistrySpecProvider) Messages() map[string]asyncapi.MessageChoices {
	return asyncapi.Reflector.SpecEns().ComponentsEns().Messages
}

func (p RegistrySpecProvider) Schemas() map[string]jsonschema.Schema {
	return asyncapi.Reflector.SpecEns().ComponentsEns().Schemas
}

func (p RegistrySpecProvider) Binding() asyncapi.BindingsObject {
	var kafkaProps interface{} = types.Pojo{
		"groupId":  "{APP_NAME}-{TOPIC_NAME}-GP",
		"clientId": "{APP_NAME}-{TOPIC_NAME}-{APP_INSTANCE_ID}",
	}

	return asyncapi.BindingsObject{
		Kafka: &kafkaProps,
	}
}

func (p RegistrySpecProvider) SecuritySchemes() asyncapi.ComponentsSecuritySchemes {
	saslPlainSecuritySchema := asyncapi.SaslPlainSecurityScheme{}
	saslSecurityScheme := (&asyncapi.SaslSecurityScheme{}).WithSaslPlainSecurityScheme(saslPlainSecuritySchema)
	securityScheme := (&asyncapi.SecurityScheme{}).WithSaslSecurityScheme(*saslSecurityScheme)
	componentsSecuritySchemesWd := (&asyncapi.ComponentsSecuritySchemesWD{}).WithSecurityScheme(*securityScheme)
	schemes := (&asyncapi.ComponentsSecuritySchemes{}).WithMapOfComponentsSecuritySchemesWDValuesItem("cpx", *componentsSecuritySchemesWd)
	return *schemes
}

func NewRegistrySpecProvider(ctx context.Context) (*RegistrySpecProvider, error) {
	appInfo, err := swaggerprovider.AppInfoFromConfig(config.FromContext(ctx))
	if err != nil {
		return nil, err
	}

	return &RegistrySpecProvider{
		appInfo: appInfo,
	}, nil
}