// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package service

import (
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers/configtest"
	"reflect"
	"testing"
)

func TestNewSecurityAccountsDefaultSettings(t *testing.T) {
	type args struct {
		cfg *config.Config
	}
	tests := []struct {
		name    string
		args    args
		want    *SecurityAccountsDefaultSettings
		wantErr bool
	}{
		{
			name: "Defaults",
			args: args{
				cfg: configtest.NewInMemoryConfig(map[string]string{}),
			},
			want: &SecurityAccountsDefaultSettings{
				Username: "system",
				Password: "system",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewSecurityAccountsDefaultSettings(tt.args.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewSecurityAccountsDefaultSettings() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSecurityAccountsDefaultSettings() gotCfg = %v, want %v", got, tt.want)
			}
		})
	}
}
