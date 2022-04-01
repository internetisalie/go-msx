// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package vault

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
)

const (
	configRootTokenSourceAppRole = "spring.cloud.vault.token-source.approle"
)

type AppRoleConfig struct {
	RoleId   string
	SecretId string
}

type AppRoleSource struct {
	cfg  *AppRoleConfig
	conn ConnectionApi
}

func (c *AppRoleSource) GetToken(ctx context.Context) (token string, err error) {
	return c.conn.LoginWithAppRole(ctx, c.cfg.RoleId, c.cfg.SecretId)
}

func (c *AppRoleSource) Renewable() bool {
	return true
}

func NewAppRoleConfig(cfg *config.Config) (*AppRoleConfig, error) {
	var appRoleConfig AppRoleConfig
	if err := cfg.Populate(&appRoleConfig, configRootTokenSourceAppRole); err != nil {
		return nil, err
	}
	return &appRoleConfig, nil
}

func NewAppRoleSource(cfg *config.Config, conn ConnectionApi) (*AppRoleSource, error) {
	appRoleConfig, err := NewAppRoleConfig(cfg)
	if err != nil {
		return nil, err
	}

	return &AppRoleSource{
		cfg:  appRoleConfig,
		conn: conn,
	}, nil
}
