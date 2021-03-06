// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package populate

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	api "cto-github.cisco.com/NFV-BU/go-msx/integration/serviceconfigmanager"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"cto-github.cisco.com/NFV-BU/go-msx/populate"
	"cto-github.cisco.com/NFV-BU/go-msx/resource"
	"cto-github.cisco.com/NFV-BU/go-msx/security/service"
	"github.com/pkg/errors"
	"path"
)

const (
	manifestFile = "manifest.json"

	artifactKeyServiceConfigs = "serviceconfigs"

	serviceConfigPopulatorConfigRoot = "populate.serviceconfig"
)

var logger = log.NewLogger("msx.integration.serviceconfigmanager.populate")

type ServiceConfigPopulatorConfig struct {
	Enabled bool   `config:"default=false"`
	Root    string `config:"default=${populate.root}/serviceconfig"`
}

func NewServiceConfigPopulatorConfigFromConfig(cfg *config.Config) (*ServiceConfigPopulatorConfig, error) {
	var populatorConfig ServiceConfigPopulatorConfig
	if err := cfg.Populate(&populatorConfig, serviceConfigPopulatorConfigRoot); err != nil {
		return nil, err
	}

	return &populatorConfig, nil
}

type ServiceConfigPopulator struct {
	cfg ServiceConfigPopulatorConfig
}

type manifest map[string][]artifact

type artifact struct {
	api.ServiceConfigurationRequest
	populate.Artifact
}

func (p ServiceConfigPopulator) Populate(ctx context.Context) error {
	if !p.cfg.Enabled {
		logger.WithContext(ctx).Warn("Service Config populator disabled.")
		return nil
	}

	return service.WithDefaultServiceAccount(ctx, func(ctx context.Context) error {
		logger.WithContext(ctx).Info("Populating service configs")

		var manifest manifest
		err := resource.
			Reference(path.Join(p.cfg.Root, manifestFile)).
			Unmarshal(&manifest)
		if err != nil {
			return err
		}

		scm, err := api.NewIntegration(ctx)
		if err != nil {
			return err
		}

		for _, serviceDefinitionArtifact := range manifest[artifactKeyServiceConfigs] {
			err = p.populateServiceConfig(ctx, scm, serviceDefinitionArtifact)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func (p ServiceConfigPopulator) populateServiceConfig(ctx context.Context, scm api.Api, artifact artifact) (err error) {
	logger.WithContext(ctx).Infof("Populating service config %q", artifact.TemplateFileName)

	var request = artifact.ServiceConfigurationRequest
	if artifact.TemplateFileName != "" {
		configuration, err := resource.Reference(
			path.Join(p.cfg.Root, artifact.TemplateFileName)).
			ReadAll()
		if err != nil {
			return errors.Wrapf(err, "Failed to load service config %q", artifact.TemplateFileName)
		}
		request.Configuration = string(configuration)
	}

	resp, err := scm.CreateServiceConfiguration(request)

	if resp == nil || (err != nil && resp.StatusCode != 409) {
		return errors.Wrapf(err, "Failed to populate service config %q", artifact.TemplateFileName)
	}

	if resp.StatusCode == 409 {
		return nil
	}

	logger.WithContext(ctx).Infof("Successfully populated service config %q", artifact.TemplateFileName)
	return nil
}

func init() {
	populate.RegisterPopulationTask(
		populate.NewPopulatorTask(
			"Populate service configurations",
			1000,
			[]string{"all", "serviceConfig", "serviceMetadata"},
			func(ctx context.Context) (populate.Populator, error) {
				populatorConfig, err := NewServiceConfigPopulatorConfigFromConfig(config.MustFromContext(ctx))
				if err != nil {
					return nil, err
				}

				return &ServiceConfigPopulator{
					cfg: *populatorConfig,
				}, nil
			}))
}
