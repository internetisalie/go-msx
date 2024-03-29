// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package secrets

import "testing"

func TestSecretsResponse_Value(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name string
		s    SecretsResponse
		args args
		want string
	}{
		{
			name: "Nil",
			s:    nil,
			args: args{key: "any"},
			want: "",
		},
		{
			name: "NotNil",
			s: SecretsResponse{
				"secret-key-1": "secret-value-1",
			},
			args: args{key: "secret-key-1"},
			want: "secret-value-1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Value(tt.args.key); got != tt.want {
				t.Errorf("Value() = %v, want %v", got, tt.want)
			}
		})
	}
}
