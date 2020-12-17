package webservice

import (
	"cto-github.cisco.com/NFV-BU/go-msx/fs"
	"cto-github.cisco.com/NFV-BU/go-msx/resource"
	"net/http"
	"path/filepath"
)

// noIndexFileSystem prevents directory listings
type noIndexFileSystem struct {
	fs http.FileSystem
}

func (nfs noIndexFileSystem) Open(path string) (http.File, error) {
	f, err := nfs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	s, err := f.Stat()
	if s.IsDir() {
		index := filepath.Join(path, "index.html")
		if _, err := nfs.fs.Open(index); err != nil {
			closeErr := f.Close()
			if closeErr != nil {
				return nil, closeErr
			}

			return nil, err
		}
	}

	return f, nil
}

// NewWebRoot creates an http.FileSystem at the specified resource path
func NewWebRoot(webRootPath string) (http.FileSystem, error) {
	vfs, err := resource.FileSystem()
	if err != nil {
		return nil, err
	}

	return fs.NewPrefixFileSystem(vfs, webRootPath)
}
