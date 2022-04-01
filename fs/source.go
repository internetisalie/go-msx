// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package fs

import (
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"path/filepath"
	"runtime"
)

// For unit testing
func SetSources() error {
	var err error
	if fsConfig == nil {
		fsConfig = new(FileSystemConfig)
	}
	_, file, _, _ := runtime.Caller(1)
	thence := types.FindSourceDirFromFile(file)
	fsConfig.Sources, err = filepath.Abs(thence)
	return err
}
