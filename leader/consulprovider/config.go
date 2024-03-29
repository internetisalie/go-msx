// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package consulprovider

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"github.com/pkg/errors"
)

const configRootConsulLeaderElection = "consul.leader.election"

var logger = log.NewLogger("msx.leader.consulprovider")

type ConsulLeaderElectionConfig struct {
	Enabled          bool   `config:"default=false"`
	DefaultMasterKey string `config:"default=service/${spring.application.name}/leader"`
	Disconnected     bool   `config:"default=${cli.flag.disconnected:false}"`
	LeaderProperties []LeaderProperties
}

type LeaderProperties struct {
	Key             string `config:"default=${consul.leader.election.defaultMasterKey}"`
	HeartBeatMillis int    `config:"default=2000"`
	BusyWaitMillis  int    `config:"default=5000"`
	Disconnected    bool   `config:"default=${cli.flag.disconnected:false}"`
}

func NewConsulLeaderElectionConfigFromConfig(cfg *config.Config) (*ConsulLeaderElectionConfig, error) {
	var leaderElectionConfig ConsulLeaderElectionConfig
	if err := cfg.Populate(&leaderElectionConfig, configRootConsulLeaderElection); err != nil {
		return nil, err
	}
	return &leaderElectionConfig, nil
}

func NewConsulLeaderElectionConfig(ctx context.Context) (*ConsulLeaderElectionConfig, error) {
	cfg := config.FromContext(ctx)
	if cfg == nil {
		return nil, errors.New("Configuration not found in context")
	}
	return NewConsulLeaderElectionConfigFromConfig(cfg)
}
