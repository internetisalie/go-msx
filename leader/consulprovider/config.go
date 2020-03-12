package consulprovider

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"errors"
)

const configRootConsulLeaderElection = "consul.leader.election"

var logger = log.NewLogger("msx.leader.consulprovider")

type ConsulLeaderElectionConfig struct {
	Enabled          bool   `config:"default=false"`
	DefaultMasterKey string `config:"default=service/${info.app.name}/leader"`
	LeaderProperties []LeaderProperties
}

type LeaderProperties struct {
	Key             string
	HeartBeatMillis int
	BusyWaitMillis  int
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