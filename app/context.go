package app

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/cassandra"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/consul"
	"cto-github.cisco.com/NFV-BU/go-msx/health"
	"cto-github.cisco.com/NFV-BU/go-msx/health/cassandracheck"
	"cto-github.cisco.com/NFV-BU/go-msx/health/consulcheck"
	"cto-github.cisco.com/NFV-BU/go-msx/health/kafkacheck"
	"cto-github.cisco.com/NFV-BU/go-msx/health/redischeck"
	"cto-github.cisco.com/NFV-BU/go-msx/health/vaultcheck"
	"cto-github.cisco.com/NFV-BU/go-msx/kafka"
	"cto-github.cisco.com/NFV-BU/go-msx/redis"
	"cto-github.cisco.com/NFV-BU/go-msx/vault"
	"cto-github.cisco.com/NFV-BU/go-msx/webservice"
	"github.com/pkg/errors"
)

func init() {
	OnEvent(EventConfigure, PhaseAfter, withConfig(configureConsulPool))
	OnEvent(EventConfigure, PhaseAfter, withConfig(configureVaultPool))
	OnEvent(EventConfigure, PhaseAfter, withConfig(configureCassandraPool))
	OnEvent(EventConfigure, PhaseAfter, withConfig(configureRedisPool))
	OnEvent(EventConfigure, PhaseAfter, withConfig(configureKafkaPool))
	OnEvent(EventConfigure, PhaseAfter, configureWebService)
}

type configHandler func(cfg *config.Config) error

func withConfig(handler configHandler) Observer {
	return func(ctx context.Context) error {
		var cfg *config.Config
		if cfg = config.FromContext(ctx); cfg == nil {
			return errors.New("Failed to retrieve config from context")
		}

		return handler(cfg)
	}
}

func configureConsulPool(cfg *config.Config) error {
	if err := consul.ConfigurePool(cfg); err != nil && err != consul.ErrDisabled {
		return err
	} else if err != consul.ErrDisabled {
		RegisterInjector(consul.ContextWithPool)
		health.RegisterCheck("consul", consulcheck.Check)
	}

	return nil
}

func configureVaultPool(cfg *config.Config) error {
	if err := vault.ConfigurePool(cfg); err != nil && err != vault.ErrDisabled {
		return err
	} else if err != vault.ErrDisabled {
		RegisterInjector(vault.ContextWithPool)
		health.RegisterCheck("vault", vaultcheck.Check)
	}

	return nil
}

func configureCassandraPool(cfg *config.Config) error {
	if err := cassandra.ConfigurePool(cfg); err != nil && err != cassandra.ErrDisabled {
		return err
	} else if err != cassandra.ErrDisabled {
		RegisterInjector(cassandra.ContextWithPool)
		health.RegisterCheck("cassandra", cassandracheck.Check)
	}

	return nil
}

func configureRedisPool(cfg *config.Config) error {
	if err := redis.ConfigurePool(cfg); err != nil && err != redis.ErrDisabled {
		return err
	} else if err != redis.ErrDisabled {
		RegisterInjector(redis.ContextWithPool)
		health.RegisterCheck("redis", redischeck.Check)
	}

	return nil
}

func configureKafkaPool(cfg *config.Config) error {
	if err := kafka.ConfigurePool(cfg); err != nil && err != kafka.ErrDisabled {
		return err
	} else if err != kafka.ErrDisabled {
		RegisterInjector(kafka.ContextWithPool)
		health.RegisterCheck("kafka", kafkacheck.Check)
	}

	return nil
}

func configureWebService(ctx context.Context) error {
	return withConfig(func(cfg *config.Config) error {
		if err := webservice.ConfigureWebServer(cfg, ctx); err != nil && err != webservice.ErrDisabled {
			return err
		} else if err != webservice.ErrDisabled {
			RegisterInjector(webservice.ContextWithWebServer)
			// TODO: health check?
		}

		return nil
	})(ctx)
}

type ContextInjector func(ctx context.Context) context.Context

var contextInjectors []ContextInjector

func RegisterInjector(injector ContextInjector) {
	contextInjectors = append(contextInjectors, injector)
}

func injectContextValues(ctx context.Context) context.Context {
	for _, contextInjector := range contextInjectors {
		ctx = contextInjector(ctx)
	}
	return ctx
}
