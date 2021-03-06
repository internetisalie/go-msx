// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package fs

import (
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/bmatcuk/doublestar"
	"github.com/shurcooL/httpfs/filter"
	"github.com/shurcooL/httpfs/vfsutil"
	"io"
	"net/http"
	"os"
	pathpkg "path"
	"strings"
)

func NewGlobFileSystem(source http.FileSystem, includes []string, excludes []string) (http.FileSystem, error) {
	var keepFiles = make(types.StringSet)
	var keepDirs = make(types.StringSet)
	err := vfsutil.WalkFiles(source, "/", func(path string, info os.FileInfo, rs io.ReadSeeker, err error) (err2 error) {
		included := false
		for _, inc := range includes {
			if included, err2 = doublestar.Match(inc, path); err2 != nil {
				return
			} else if included {
				break
			}
		}

		if !included {
			return nil
		}

		var excluded = false
		for _, exc := range excludes {
			if !strings.HasPrefix(exc, "/") && !strings.HasPrefix(exc, "**/") && exc != "**" {
				exc = "/" + exc
			}
			if excluded, err2 = doublestar.Match(exc, path); err != nil {
				return
			} else if excluded {
				break
			}
		}

		if !excluded {
			keepFiles.Add(path)
			dir := pathpkg.Dir(path)
			for dir != "/" && !keepDirs.Contains(dir) {
				keepDirs.Add(dir)
				dir = pathpkg.Dir(dir)
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return filter.Keep(source, func(p string, fi os.FileInfo) bool {
		if keepFiles.Contains(p) {
			return true
		}
		if !fi.IsDir() {
			return false
		}
		if p == "/" {
			return true
		}
		if keepDirs.Contains(p) {
			return true
		}
		return false
	}), nil
}

func ListFiles(source http.FileSystem) ([]string, error) {
	var results []string
	err := vfsutil.WalkFiles(source, "/", func(path string, info os.FileInfo, rs io.ReadSeeker, e error) (err error) {
		if !info.IsDir() {
			results = append(results, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return results, nil
}
