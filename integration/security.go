// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package integration

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"encoding/base64"
)

const (
	configRootIntegrationSecurityClient = "integration.security.client"
)

type SecurityClientSettings struct {
	ClientId     string `config:"default=nfv-service"`
	ClientSecret string `config:"default=nfv-service-secret"`
}

func (c *SecurityClientSettings) Authorization() string {
	clientCredentials := []byte(c.ClientId + ":" + c.ClientSecret)
	return "Basic " + base64.StdEncoding.EncodeToString(clientCredentials)
}

func NewSecurityClientSettings(ctx context.Context) (cfg *SecurityClientSettings, err error) {
	cfg = &SecurityClientSettings{}
	err = config.FromContext(ctx).LatestValues().Populate(cfg, configRootIntegrationSecurityClient)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
