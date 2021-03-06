// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package build

import (
	"github.com/bmatcuk/doublestar"
	copypkg "github.com/otiai10/copy"
	"path/filepath"
	"strings"
)

func init() {
	AddTarget("install-resources", "Installs Resources", InstallResources)
}

func InstallResources(args []string) error {
	files, err := collectIncludedResources(BuildConfig.Resources)
	if err != nil {
		return err
	}

	for _, inputFilePath := range files {
		mappedInputFilePath := getResourcePathMapping(BuildConfig.Resources, inputFilePath)
		outputFilePath := filepath.Join(BuildConfig.OutputResourcesPath(), mappedInputFilePath)
		logger.Infof("Copying %s to %s", inputFilePath, outputFilePath)
		if err = copypkg.Copy(inputFilePath, outputFilePath); err != nil {
			return err
		}
	}

	return nil
}

func collectIncludedResources(resources Resources) ([]string, error) {
	var results []string
	for _, inc := range resources.Includes {
		if strings.HasPrefix(inc, "/") {
			inc = inc[1:]
		}

		incFiles, err := doublestar.Glob(inc)
		if err != nil {
			return nil, err
		}

		for _, incFile := range incFiles {
			excluded, err := excludeFilteredResource(resources, incFile)
			if err != nil {
				return nil, err
			}

			if !excluded {
				results = append(results, incFile)
			}
		}
	}

	return results, nil
}

func excludeFilteredResource(resources Resources, included string) (bool, error) {
	excludes := resources.Excludes
	excludes = append(excludes, "/dist/**", "/test/**", "/local/**", "/vendor/**")
	for _, exc := range excludes {
		if strings.HasPrefix(exc, "/") {
			exc = exc[1:]
		}

		matches, err := doublestar.Match(exc, included)
		if err != nil {
			return false, err
		} else if matches {
			return true, err
		}
	}
	return false, nil
}

func getResourcePathMapping(resources Resources, resourcePath string) string {
	if !strings.HasPrefix(resourcePath, "/") {
		resourcePath = "/" + resourcePath
	}
	for _, pathMapping := range resources.Mappings {
		pathFrom, pathTo := pathMapping.From, pathMapping.To
		if !strings.HasSuffix(pathFrom, "/") {
			pathFrom += "/"
		}
		if !strings.HasSuffix(pathTo, "/") {
			pathTo += "/"
		}
		if strings.HasPrefix(resourcePath, pathFrom) {
			resourcePath = strings.TrimPrefix(resourcePath, pathFrom)
			resourcePath = pathTo + resourcePath
			return resourcePath
		}
	}
	return resourcePath
}
