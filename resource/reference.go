// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package resource

import (
	"cto-github.cisco.com/NFV-BU/go-msx/fs"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

var logger = log.NewLogger("msx.resource")

type Ref string

func References(resourceGlob string) (refs []Ref) {
	absGlob := abs(resourceGlob)
	resourceFs, err := FileSystem()
	if err != nil {
		logger.WithError(err).Error("Failed to open resource FileSystem")
		return nil
	}

	globFs, err := fs.NewGlobFileSystem(resourceFs, []string{absGlob}, nil)
	if err != nil {
		logger.WithError(err).Error("Failed to create glob filesystem")
		return nil
	}

	files, err := fs.ListFiles(globFs)
	if err != nil {
		logger.WithError(err).Error("Failed to list files from glob filesystem")
		return nil
	}

	for _, file := range files {
		refs = append(refs, Ref(file))
	}
	return refs
}

func Reference(resourceName string) (ref Ref) {
	return Ref(abs(resourceName))
}

func (r Ref) String() string {
	return string(r)
}

func (r Ref) ReadAll() (data []byte, err error) {
	return load(string(r))
}

func (r Ref) Open() (http.File, error) {
	return open(string(r))
}

func (r Ref) Unmarshal(target interface{}) (err error) {
	bytes, err := load(string(r))
	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, target)
}

func ReadAll(resourceName string) (data []byte, err error) {
	return Ref(abs(resourceName)).ReadAll()
}

// Deprecated.  Use resource.ReadAll()
func Load(resourceName string) (data []byte, err error) {
	return Ref(abs(resourceName)).ReadAll()
}

func Unmarshal(resourceName string, target interface{}) (err error) {
	return Ref(abs(resourceName)).Unmarshal(target)
}

func sourcePath(p string) (string, error) {
	if fs.Config() == nil || fs.Config().Sources == "" {
		logger.Error("SourcePath called with nil FS configuration")
		return "/nil", ErrFilesystemUnavailable
	}

	fp := filepath.Clean(p)
	sp := fs.Config().Sources
	p = strings.TrimPrefix(fp, sp)
	return filepath.ToSlash(p), nil
}

func open(resourcePath string) (http.File, error) {
	fileSystem, err := FileSystem()
	if err != nil {
		return nil, err
	}

	reader, err := fileSystem.Open(resourcePath)
	if err != nil {
		return nil, err
	}

	return reader, err
}

func load(resourcePath string) ([]byte, error) {
	reader, err := open(resourcePath)
	if err != nil {
		return nil, err
	}

	defer reader.Close()

	return ioutil.ReadAll(reader)
}

func abs(filename string) string {
	if strings.HasPrefix(filename, "/") {
		return filename
	}

	_, file, _, ok := runtime.Caller(2)
	if !ok {
		logger.Errorf("Failed to identify caller: %s", filename)
		return filename
	}

	base := path.Dir(file)
	full := path.Join(base, filename)

	absPath, err := sourcePath(full)
	if err != nil {
		logger.WithError(err).Errorf("Failed to resolve source path: %s", full)
	}

	return absPath
}
