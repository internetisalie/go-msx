// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package fs

import "testing"

func TestSetSources(t *testing.T) {
	tests := []struct {
		name     string
		fsConfig *FileSystemConfig
		wantErr  bool
	}{
		{
			name:     "Simple",
			fsConfig: nil,
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fsConfig = tt.fsConfig
			if err := SetSources(); (err != nil) != tt.wantErr {
				t.Errorf("SetSources() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
