// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package app

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/app/version"
	"cto-github.cisco.com/NFV-BU/go-msx/cli"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/config/cobraprovider"
	"cto-github.cisco.com/NFV-BU/go-msx/populate"
	"cto-github.cisco.com/NFV-BU/go-msx/repository/migrate"
	_ "cto-github.cisco.com/NFV-BU/go-msx/stats/logstats"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

const (
	configKeyRedisEnable           = "spring.redis.enable"
	configKeyKafkaEnable           = "spring.cloud.stream.kafka.binder.enabled"
	configKeyConsulDiscoveryEnable = "spring.cloud.consul.discovery.enabled"
	configKeyConsulEnable          = "spring.cloud.consul.enabled"
	configKeyVaultEnable           = "spring.cloud.vault.enabled"
	configKeyServerEnable          = "server.enabled"
	configKeyLeaderEnable          = "consul.leader.election.enabled"
	configKeyCassandraEnable       = "spring.data.cassandra.enable"
	configKeySqlDbEnable           = "spring.datasource.enabled"

	CommandRoot     = ""
	CommandMigrate  = "migrate"
	CommandPopulate = "populate"
	CommandVersion  = "version"

	configValueFalse = "false"
)

func init() {
	// Configure the root command
	cmd := cli.RootCmd()

	cmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		RegisterProviderFactory(SourceCommandLine, func(name string, cfg *config.Config) ([]config.Provider, error) {

			appInstance := types.RandString(5)

			return []config.Provider{
				config.NewCacheProvider(
					cobraprovider.NewProvider(name, cmd, "cli.flag."),
				),
				config.NewCacheProvider(config.NewInMemoryProvider("Built-In", map[string]string{
					"spring.application.name":     strings.Fields(cli.RootCmd().Use)[0],
					"spring.application.instance": appInstance,
					"info.app.name":               strings.Fields(cli.RootCmd().Use)[0],
				})),
			}, nil
		})
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		return application.Run(CommandRoot)
	}

	if migrateCommand, err := AddCommand(CommandMigrate, "Migrate database schema", migrate.Migrate, commandMigrateInit); err != nil {
		cli.Fatal(err)
	} else {
		migrate.CustomizeCommand(migrateCommand)
	}

	if populateCommand, err := AddCommand(CommandPopulate, "Populate data", populate.Populate, commandPopulateInit); err != nil {
		cli.Fatal(err)
	} else {
		populate.CustomizeCommand(populateCommand)
	}

	if _, err := AddCommand(CommandVersion, "Show version", version.Version, commandVersionInit); err != nil {
		cli.Fatal(err)
	}
}

func AddCommand(path, brief string, command CommandObserver, init Observer) (cmd *cobra.Command, err error) {
	return cli.AddCommand(path, brief, func(args []string) error {
		OnEvent(EventCommand, path, init)
		OnEvent(EventReady, PhaseAfter, func(ctx context.Context) error {
			logger.Infof("Executing command: %s", strings.Fields(cmd.Use)[0])
			if err := command(ctx, args); err != nil {
				logger.Errorf("Command %s returned error: %v", strings.Fields(cmd.Use)[0], err)
				cli.SetExitCode(1)
			}
			return application.Stop()
		})

		return application.Run(path)
	})
}

func Run(appName string) {
	// Convert environment variable POPULATE into migrate command
	if strings.ToLower(os.Getenv("POPULATE")) == "database" {
		os.Args = append(os.Args, "migrate")
	}

	cli.Run(appName)
}

func Noop(context.Context) error {
	return nil
}

func commandMigrateInit(context.Context) error {
	OverrideConfig(map[string]string{
		configKeyRedisEnable:           configValueFalse,
		configKeyKafkaEnable:           configValueFalse,
		configKeyConsulDiscoveryEnable: configValueFalse,
		configKeyServerEnable:          configValueFalse,
		configKeyLeaderEnable:          configValueFalse,
	})

	return nil
}

func commandPopulateInit(context.Context) error {
	OverrideConfig(map[string]string{
		configKeyConsulDiscoveryEnable: configValueFalse,
		configKeyServerEnable:          configValueFalse,
		configKeyLeaderEnable:          configValueFalse,
	})

	return nil
}

func commandVersionInit(context.Context) error {
	OverrideConfig(map[string]string{
		configKeyRedisEnable:           configValueFalse,
		configKeyKafkaEnable:           configValueFalse,
		configKeyConsulDiscoveryEnable: configValueFalse,
		configKeyServerEnable:          configValueFalse,
		configKeyLeaderEnable:          configValueFalse,
		configKeyCassandraEnable:       configValueFalse,
		configKeySqlDbEnable:           configValueFalse,
		configKeyConsulEnable:          configValueFalse,
		configKeyVaultEnable:           configValueFalse,
	})

	return nil
}
