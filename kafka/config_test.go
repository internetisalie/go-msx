// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package kafka

import (
	"cto-github.cisco.com/NFV-BU/go-msx/certificate"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers/configtest"
	"reflect"
	"testing"
)

func TestNewConnectionConfig(t *testing.T) {
	type args struct {
		cfg *config.Config
	}
	tests := []struct {
		name    string
		args    args
		want    *ConnectionConfig
		wantErr bool
	}{
		{
			name: "Defaults",
			args: args{
				cfg: configtest.NewInMemoryConfig(map[string]string{
					"spring.application.name":     "test",
					"spring.application.instance": "XYZABC",
				}),
			},
			want: &ConnectionConfig{
				Brokers:                []string{"localhost"},
				DefaultBrokerPort:      9092,
				ZkNodes:                []string{"localhost"},
				DefaultZkPort:          2181,
				OffsetUpdateTimeWindow: 10000,
				OffsetUpdateCount:      0,
				RequiredAcks:           1,
				MinPartitionCount:      1,
				ReplicationFactor:      1,
				AutoCreateTopics:       true,
				DefaultPartitions:      1,
				Version:                "2.0.0",
				ClientId:               "test",
				ClientIdSuffix:         "XYZABC",
				Enabled:                false,
				Partitioner:            "hash",
				Tls: certificate.TLSConfig{
					MinVersion:        "tls12",
					CertificateSource: "server",
					CipherSuites: []string{
						"TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305",
						"TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256",
						"TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384",
						"TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA",
						"TLS_RSA_WITH_AES_256_GCM_SHA384",
						"TLS_RSA_WITH_AES_256_CBC_SHA",
					},
				},
			},
		},
		{
			name: "Custom",
			args: args{
				cfg: configtest.NewInMemoryConfig(map[string]string{
					"spring.cloud.stream.kafka.binder.brokers":                   "remote-vm,remote-vm2",
					"spring.cloud.stream.kafka.binder.broker-cert-name":          "kafka",
					"spring.cloud.stream.kafka.binder.default-broker-port":       "9999",
					"spring.cloud.stream.kafka.binder.zk-nodes":                  "remote-vm,remote-vm2",
					"spring.cloud.stream.kafka.binder.default-zk-port":           "9998",
					"spring.cloud.stream.kafka.binder.offset-update-time-window": "9997",
					"spring.cloud.stream.kafka.binder.offset-update-count":       "1",
					"spring.cloud.stream.kafka.binder.required-acks":             "2",
					"spring.cloud.stream.kafka.binder.min-partition-count":       "3",
					"spring.cloud.stream.kafka.binder.replication-factor":        "4",
					"spring.cloud.stream.kafka.binder.auto-create-topics":        "false",
					"spring.cloud.stream.kafka.binder.default-partitions":        "5",
					"spring.cloud.stream.kafka.binder.version":                   "2.2.0",
					"spring.cloud.stream.kafka.binder.client-id":                 "test",
					"spring.cloud.stream.kafka.binder.client-id-suffix":          "XYZABC",
					"spring.cloud.stream.kafka.binder.enabled":                   "true",
					"spring.cloud.stream.kafka.binder.partitioner":               "random",
				}),
			},
			want: &ConnectionConfig{
				Brokers:                []string{"remote-vm", "remote-vm2"},
				DefaultBrokerPort:      9999,
				ZkNodes:                []string{"remote-vm", "remote-vm2"},
				DefaultZkPort:          9998,
				OffsetUpdateTimeWindow: 9997,
				OffsetUpdateCount:      1,
				RequiredAcks:           2,
				MinPartitionCount:      3,
				ReplicationFactor:      4,
				AutoCreateTopics:       false,
				DefaultPartitions:      5,
				Version:                "2.2.0",
				ClientId:               "test",
				ClientIdSuffix:         "XYZABC",
				Enabled:                true,
				Partitioner:            "random",
				Tls: certificate.TLSConfig{
					MinVersion:        "tls12",
					CertificateSource: "server",
					CipherSuites: []string{
						"TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305",
						"TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256",
						"TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384",
						"TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA",
						"TLS_RSA_WITH_AES_256_GCM_SHA384",
						"TLS_RSA_WITH_AES_256_CBC_SHA",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewConnectionConfig(tt.args.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewConnectionConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf(testhelpers.Diff(tt.want, got))
			}
		})
	}
}
