package cassandra

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"github.com/gocql/gocql"
	"github.com/pkg/errors"
	"sync"
)

var pool *ConnectionPool
var poolMtx sync.Mutex

type ConnectionPool struct {
	cfg     *ClusterConfig
	cluster *Cluster
}

func (p *ConnectionPool) WithSession(action func(*gocql.Session) error) error {
	session, err := p.cluster.CreateSession()
	if err != nil {
		return err
	}
	defer session.Close()

	return action(session)
}

func (p *ConnectionPool) ClusterConfig() ClusterConfig {
	return *p.cfg
}

func CreateKeyspaceForPool(ctx context.Context) error {
	if pool == nil {
		// Cassandra is disabled
		return nil
	}

	systemClusterConfig := pool.ClusterConfig()
	targetKeyspaceName := systemClusterConfig.KeyspaceName
	systemClusterConfig.KeyspaceName = keyspaceSystem

	logger.WithContext(ctx).Infof("Ensuring keyspace %s exists", targetKeyspaceName)

	cluster, err := NewCluster(&systemClusterConfig)
	if err != nil {
		return err
	}

	return cluster.createKeyspace(ctx, targetKeyspaceName, systemClusterConfig.KeyspaceOptions)
}

func Pool() *ConnectionPool {
	return pool
}

func ConfigurePool(cfg *config.Config) error {
	poolMtx.Lock()
	defer poolMtx.Unlock()

	if pool != nil {
		return nil
	}

	clusterConfig, err := NewClusterConfigFromConfig(cfg)
	if err != nil {
		return err
	}

	cluster, err := NewClusterFromConfig(cfg)
	if err != nil {
		return err
	}

	pool = &ConnectionPool{
		cfg:     clusterConfig,
		cluster: cluster,
	}

	return nil
}

type cassandraContextKey int

const contextKeyCassandraPool cassandraContextKey = iota

func ContextWithPool(ctx context.Context) context.Context {
	return context.WithValue(ctx, contextKeyCassandraPool, pool)
}

func PoolFromContext(ctx context.Context) (*ConnectionPool, error) {
	connectionPoolInterface := ctx.Value(contextKeyCassandraPool)
	if connectionPoolInterface == nil {
		return nil, ErrDisabled
	}
	if connectionPool, ok := connectionPoolInterface.(*ConnectionPool); !ok {
		return nil, errors.New("Context cassandra connection pool value is the wrong type")
	} else {
		return connectionPool, nil
	}
}