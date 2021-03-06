// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package sqldb

import (
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers/configtest"
	"reflect"
	"testing"
)

func TestNewSqlConfigFromConfig(t *testing.T) {
	type args struct {
		cfg *config.Config
	}
	tests := []struct {
		name    string
		args    args
		want    *Config
		wantErr bool
	}{
		{
			name: "EmbeddedDefaults",
			args: args{
				cfg: configtest.NewInMemoryConfig(map[string]string{
					"db.cockroach.host":                  "localhost",
					"db.cockroach.port":                  "26257",
					"spring.datasource.driver":           "postgres",
					"spring.datasource.name":             "",
					"spring.datasource.username":         "root",
					"spring.datasource.password":         "",
					"spring.datasource.data-source-name": "postgresql://${spring.datasource.username}:${spring.datasource.password}@${db.cockroach.host}:${db.cockroach.port}/${spring.datasource.name}?sslmode=disable",
				}),
			},
			want: &Config{
				Driver:         "postgres",
				DataSourceName: "postgresql://root:@localhost:26257/?sslmode=disable",
				Enabled:        false,
			},
			wantErr: false,
		},
		{
			name: "Custom",
			args: args{
				cfg: configtest.NewInMemoryConfig(map[string]string{
					"spring.datasource.driver":           "sqlite3",
					"spring.datasource.enabled":          "true",
					"spring.datasource.name":             "TestNewSqlConfigFromConfig",
					"spring.datasource.data-source-name": "file:${spring.datasource.name}?cache=shared&mode=memory",
				}),
			},
			want: &Config{
				Driver:         "sqlite3",
				DataSourceName: "file:TestNewSqlConfigFromConfig?cache=shared&mode=memory",
				Enabled:        true,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewSqlConfigFromConfig(tt.args.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewSqlConfigFromConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSqlConfigFromConfig() got = %v, want %v", got, tt.want)
			}
		})
	}
}
