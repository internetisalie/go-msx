package kafka

import (
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/pkg/errors"
	"strings"
	"time"
)

const (
	configRootKafka       = "spring.cloud.stream.kafka.binder"
	configKeyAppName      = "spring.application.name"
	configDefaultClientId = "sarama"
)

type ConnectionConfig struct {
	Brokers                []string `config:"default=localhost"`
	DefaultBrokerPort      int      `config:"default=9092"`
	ZkNodes                []string `config:"default=localhost"`
	DefaultZkPort          int      `config:"default=2181"`
	OffsetUpdateTimeWindow int      `config:"default=10000"`
	OffsetUpdateCount      int      `config:"default=0"`
	RequiredAcks           int      `config:"default=1"`
	MinPartitionCount      int      `config:"default=1"`
	ReplicationFactor      int      `config:"default=1"`
	AutoCreateTopics       bool     `config:"default=true"`
	DefaultPartitions      int      `config:"default=12"`
	Version                string   `config:"default=2.0.1"`
	ClientId               string   `config:"default=sarama"`
	Enabled                bool     `config:"default=false"`
}

func (c *ConnectionConfig) SaramaConfig() (*sarama.Config, error) {
	saramaConfig := sarama.NewConfig()
	saramaConfig.ClientID = c.ClientId
	saramaConfig.Consumer.Fetch.Default = 1024 * 1024
	saramaConfig.Consumer.MaxWaitTime = 500 * time.Millisecond
	saramaConfig.Consumer.MaxProcessingTime = 15 * time.Second
	saramaConfig.Consumer.Return.Errors = true
	saramaConfig.Producer.Retry.Max = 5
	saramaConfig.Producer.Return.Successes = true
	saramaConfig.Net.DialTimeout = 15 * time.Second
	saramaConfig.Net.ReadTimeout = 15 * time.Second
	saramaConfig.Net.WriteTimeout = 15 * time.Second
	saramaConfig.Metadata.Retry.Backoff = time.Second * 2

	if c.RequiredAcks == 1 {
		saramaConfig.Producer.RequiredAcks = sarama.WaitForLocal
	} else if c.RequiredAcks == 0 {
		saramaConfig.Producer.RequiredAcks = sarama.NoResponse
	} else {
		saramaConfig.Producer.RequiredAcks = sarama.WaitForAll
	}

	if kafkaVersion, err := c.findSupportedVersion(c.Version); err == nil {
		saramaConfig.Version = *kafkaVersion
	} else {
		return nil, errors.Wrap(err, "Failed to find supported kafka version")
	}

	return saramaConfig, nil
}

func (c *ConnectionConfig) BrokerAddresses() []string {
	var results []string

	for _, broker := range c.Brokers {
		if !strings.Contains(broker, ":") {
			results = append(results, fmt.Sprintf("%s:%d", broker, c.DefaultBrokerPort))
		} else {
			results = append(results, broker)
		}
	}

	return results
}

func (c *ConnectionConfig) findSupportedVersion(version string) (*sarama.KafkaVersion, error) {
	for _, v := range sarama.SupportedVersions {
		if v.String() == version {
			return &v, nil
		}
	}
	return nil, errors.New("Unsupported version: " + version)
}

func NewConnectionConfig(cfg *config.Config) (*ConnectionConfig, error) {
	connectionConfig := new(ConnectionConfig)
	if err := cfg.Populate(connectionConfig, configRootKafka); err != nil {
		return nil, err
	}

	// Override the ClientID
	if connectionConfig.ClientId == configDefaultClientId {
		if appName, err := cfg.String(configKeyAppName); err != nil {
			return nil, err
		} else {
			connectionConfig.ClientId = appName
		}
	}

	return connectionConfig, nil
}