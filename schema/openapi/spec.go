// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package openapi

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/swaggest/openapi-go/openapi3"
	"net/http"
	"regexp"
	"sort"
	"strings"
)

func Spec() *openapi3.Spec {
	return documentationReflector.SpecEns()
}

var findPathParameterRegex = regexp.MustCompile(`{([^}:]+)(:[^/]+)?(?:})`)

func SetSpecOperation(method, path string, operation *openapi3.Operation) error {
	pathParametersSubmatches := findPathParameterRegex.FindAllStringSubmatch(path, -1)

	switch method {
	case http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete,
		http.MethodPatch, http.MethodOptions, http.MethodHead, http.MethodTrace:
		break
	default:
		return fmt.Errorf("unexpected http method: %s", method)
	}

	s := Spec()

	pathItem := s.Paths.MapOfPathItemValues[path]
	pathParams := map[string]bool{}

	if len(pathParametersSubmatches) > 0 {
		for _, submatch := range pathParametersSubmatches {
			pathParams[submatch[1]] = true

			if submatch[2] != "" { // Remove gorilla.Mux-style regexp in path
				path = strings.Replace(path, submatch[0], "{"+submatch[1]+"}", 1)
			}
		}
	}

	paramIndex := make(map[string]bool, len(operation.Parameters))
	var errs []string
	for _, p := range operation.Parameters {
		p = ResolveParameter(p)
		if p.Parameter == nil {
			continue
		}

		if found := paramIndex[p.Parameter.Name+string(p.Parameter.In)]; found {
			errs = append(errs, "duplicate parameter in "+string(p.Parameter.In)+": "+p.Parameter.Name)

			continue
		}

		if found := pathParams[p.Parameter.Name]; !found && p.Parameter.In == openapi3.ParameterInPath {
			errs = append(errs, "missing path parameter placeholder in url: "+p.Parameter.Name)

			continue
		}

		paramIndex[p.Parameter.Name+string(p.Parameter.In)] = true
	}

	for pathParam := range pathParams {
		if !paramIndex[pathParam+string(openapi3.ParameterInPath)] {
			errs = append(errs, "undefined path parameter: "+pathParam)
		}
	}

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, ", "))
	}

	pathItem.WithMapOfOperationValuesItem(strings.ToLower(method), *operation)

	s.Paths.WithMapOfPathItemValuesItem(path, pathItem)

	return nil
}

func ResolveParameter(p openapi3.ParameterOrRef) openapi3.ParameterOrRef {
	if p.Parameter != nil {
		return p
	}
	if p.ParameterReference == nil {
		return p
	}

	refName := ParameterRefName(p.ParameterReference)
	if target, ok := documentationReflector.SpecEns().ComponentsEns().ParametersEns().MapOfParameterOrRefValues[refName]; ok {
		return ResolveParameter(target)
	}

	return p
}

func LookupSchema(schemaName string) (*openapi3.Schema, bool) {
	result, ok := documentationReflector.SpecEns().ComponentsEns().SchemasEns().MapOfSchemaOrRefValues[schemaName]
	if !ok {
		return nil, false
	}
	return result.Schema, ok
}

// Tags is a sortable, unique list of openapi3 tags
type Tags []openapi3.Tag

func (t Tags) Len() int {
	return len(t)
}

func (t Tags) Less(i, j int) bool {
	return strings.Compare(t[i].Name, t[j].Name) < 0
}

func (t Tags) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

func (t Tags) Index(name string) int {
	for i, tag := range t {
		if tag.Name == name {
			return i
		}
	}

	return -1
}

func (t Tags) Add(tag openapi3.Tag) Tags {
	result := append(Tags(nil), t...)

	i := t.Index(tag.Name)
	if i == -1 {
		result = append(result, tag)
	} else {
		result[i] = tag
	}

	sort.Sort(result)

	return result
}

func AddTag(name string, description string) {
	tag := openapi3.Tag{}
	tag.WithName(name)
	tag.WithDescription(description)

	tags := Tags(Spec().Tags).Add(tag)
	Spec().WithTags(tags...)
}
