// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

//go:generate staticfiles -o templates.go _templates/

package skel

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"cto-github.cisco.com/NFV-BU/go-msx/exec"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type FileFormat int

const (
	FileFormatGo FileFormat = iota
	FileFormatMakefile
	FileFormatJson
	FileFormatSql
	FileFormatYaml
	FileFormatXml
	FileFormatGroovy
	FileFormatProperties
	FileFormatMarkdown
	FileFormatGoMod
	FileFormatDocker
	FileFormatBash
)

type RenderOptions struct {
	Variables  map[string]string
	Conditions map[string]bool
	Strings    map[string]string
}

func (r RenderOptions) AddString(source, dest string) {
	r.Strings[source] = dest
}

func (r RenderOptions) AddStrings(strings map[string]string) {
	for k, v := range strings {
		r.Strings[k] = v
	}
}

func (r RenderOptions) AddCondition(condition string, value bool) {
	r.Conditions[condition] = value
}

func (r RenderOptions) AddConditions(conditions map[string]bool) {
	for k, v := range conditions {
		r.Conditions[k] = v
	}
}

func NewRenderOptions() RenderOptions {
	return RenderOptions{
		Variables: map[string]string{
			"app.name":                     skeletonConfig.AppName,
			"app.shortname":                strings.TrimSuffix(skeletonConfig.AppName, "service"),
			"app.description":              skeletonConfig.AppDescription,
			"app.displayname":              skeletonConfig.AppDisplayName,
			"app.version":                  skeletonConfig.AppVersion,
			"app.migrateversion":           skeletonConfig.AppMigrateVersion(),
			"app.packageurl":               skeletonConfig.AppPackageUrl(),
			"deployment.group":             skeletonConfig.DeploymentGroup,
			"server.port":                  strconv.Itoa(skeletonConfig.ServerPort),
			"server.contextpath":           path.Clean("/" + skeletonConfig.ServerContextPath),
			"kubernetes.group":             skeletonConfig.KubernetesGroup,
			"target.dir":                   skeletonConfig.TargetDirectory(),
			"repository.cassandra.enabled": strconv.FormatBool(skeletonConfig.Repository == "cassandra"),
			"repository.cockroach.enabled": strconv.FormatBool(skeletonConfig.Repository == "cockroach"),
			"jenkins.publish.trunk":        strconv.FormatBool(skeletonConfig.KubernetesGroup != "platformms"),
			"generator":                    skeletonConfig.Archetype,
			"beat.protocol":                skeletonConfig.BeatProtocol,
			"service.type":                 skeletonConfig.ServiceType,
			"slack.channel":                skeletonConfig.SlackChannel,
			"trunk":                        skeletonConfig.Trunk,
		},
		Conditions: map[string]bool{
			"REPOSITORY_COCKROACH": skeletonConfig.Repository == "cockroach",
			"REPOSITORY_CASSANDRA": skeletonConfig.Repository == "cassandra",
			"GENERATOR_APP":        skeletonConfig.Archetype == archetypeKeyApp,
			"GENERATOR_BEAT":       skeletonConfig.Archetype == archetypeKeyBeat,
			"GENERATOR_SP":         skeletonConfig.Archetype == archetypeKeyServicePack,
			"UI":                   hasUI(),
		},
		Strings: make(map[string]string),
	}
}

type Template struct {
	Name       string
	DestFile   string
	SourceFile string
	SourceData []byte
	Format     FileFormat
}

func (t Template) source(options RenderOptions) (string, error) {
	if t.SourceData != nil {
		return string(t.SourceData), nil
	}

	sourceFile := substituteVariables(t.SourceFile, options.Variables)

	f, ok := staticFiles[sourceFile]
	if !ok {
		return "", errors.Errorf("Template file not found: %s", sourceFile)
	}

	var reader io.Reader
	if f.size != 0 {
		var err error
		reader, err = gzip.NewReader(strings.NewReader(f.data))
		if err != nil {
			return "", err
		}
	} else {
		reader = strings.NewReader(f.data)
	}

	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func (t Template) destinationFile(options RenderOptions) (string, error) {
	if t.DestFile == "" {
		if t.SourceFile == "" {
			return "", errors.New("Missing destination filename")
		}
		return substituteVariables(t.SourceFile, options.Variables), nil
	}
	return substituteVariables(t.DestFile, options.Variables), nil
}

func (t Template) Render(options RenderOptions) error {
	// Load the source
	contents, err := t.source(options)
	if err != nil {
		return err
	}

	// Find the destination
	destFile, err := t.destinationFile(options)
	if err != nil {
		return err
	}

	logger.Infof("- %s (%s)", t.Name, destFile)

	// Replace strings
	for sourceString, destString := range options.Strings {
		contents = strings.ReplaceAll(contents, sourceString, destString)
	}

	// Substitute variables
	contents = substituteVariables(contents, options.Variables)

	// Execute conditions
	for condition, value := range options.Conditions {
		contents, err = processConditionalBlocks(contents, t.Format, condition, value)
	}
	if err != nil {
		return err
	}

	// Ensure the target parent directory exists
	targetFileName := path.Join(skeletonConfig.TargetDirectory(), destFile)
	targetDirectory := path.Dir(targetFileName)
	err = os.MkdirAll(targetDirectory, 0755)
	if err != nil {
		return err
	}

	// Write the rendered contents to the destination file
	err = ioutil.WriteFile(targetFileName, []byte(contents), 0644)
	if err != nil {
		return err
	}

	if t.Format == FileFormatGo {
		err = exec.ExecutePipes(
			exec.Info("  - Reformatting %q", path.Base(destFile)),
			exec.ExecSimple("go", "fmt", targetFileName))
	}

	return err
}

type TemplateSet []Template

func (t TemplateSet) Render(options RenderOptions) error {
	for _, template := range t {
		if err := template.Render(options); err != nil {
			return err
		}
	}
	return nil
}

func substituteVariables(source string, variableValues map[string]string) string {
	rendered := source
	variableInstanceRegex := regexp.MustCompile(`\${([^}]+)}`)
	for _, variableInstance := range variableInstanceRegex.FindAllStringSubmatch(rendered, -1) {
		variableName := variableInstance[1]
		variableValue, ok := variableValues[strings.ToLower(variableName)]
		if ok {
			rendered = strings.ReplaceAll(rendered, "${"+variableName+"}", variableValue)
		}
	}
	return rendered
}

func conditionalMarkers(format FileFormat) (string, string) {
	switch format {
	case FileFormatMakefile, FileFormatYaml, FileFormatProperties, FileFormatDocker, FileFormatBash:
		return "#", ""
	case FileFormatSql:
		return "--#", ""
	case FileFormatXml, FileFormatMarkdown:
		return "<--#", "-->"
	default:
		return "//#", ""
	}
}

func processConditionalBlocks(data string, format FileFormat, condition string, output bool) (result string, err error) {
	type parserState int
	const outside parserState = 0
	const insideIf parserState = 1
	const insideElse parserState = 2

	sb := strings.Builder{}
	write := func(out bool, line string) {
		if !out {
			return
		}
		sb.WriteString(line)
		sb.WriteRune('\n')
	}
	insideCondition := outside
	prefix, suffix := conditionalMarkers(format)
	startMarker := prefix + "if " + condition + suffix
	middleMarker := prefix + "else " + condition + suffix
	endMarker := prefix + "endif " + condition + suffix

	scanner := bufio.NewScanner(strings.NewReader(data))
	for scanner.Scan() {
		line := scanner.Text()
		lineTrimmed := strings.TrimSpace(line)
		switch insideCondition {
		case outside:
			switch lineTrimmed {
			case startMarker:
				insideCondition = insideIf
			default:
				write(true, line)
			}

		case insideIf:
			switch lineTrimmed {
			case endMarker:
				insideCondition = outside
			case middleMarker:
				insideCondition = insideElse
			default:
				write(output, line)
			}

		case insideElse:
			switch lineTrimmed {
			case endMarker:
				insideCondition = outside
			default:
				write(!output, line)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return "", errors.Wrap(err, "Failed to process conditional blocks")
	}

	return sb.String(), nil
}

func initializePackageFromFile(fileName, packageUrl string) error {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, fileName, nil, 0)
	if err != nil {
		return err
	}

	// Add the imports
	for i := 0; i < len(f.Decls); i++ {
		d := f.Decls[i]

		switch d.(type) {
		case *ast.GenDecl:
			dd := d.(*ast.GenDecl)

			// IMPORT Declarations
			if dd.Tok == token.IMPORT {
				// Add the new import
				iSpec := &ast.ImportSpec{
					Path: &ast.BasicLit{Value: strconv.Quote(packageUrl)},
					Name: ast.NewIdent("_"),
				}

				dd.Specs = append(dd.Specs, iSpec)
			}
		}
	}

	// Sort the imports
	ast.SortImports(fset, f)

	var output []byte
	buffer := bytes.NewBuffer(output)
	if err := printer.Fprint(buffer, fset, f); err != nil {
		return err
	}

	return ioutil.WriteFile(fileName, buffer.Bytes(), 0644)
}

func hasUI() bool {
	uiPath := filepath.Join(skeletonConfig.TargetDirectory(), "ui", "package.json")
	if st, err := os.Stat(uiPath); err != nil {
		return false
	} else {
		return !st.IsDir() && st.Size() > 0
	}
}

func iff(cond bool, truth, falsehood string) string {
	if cond {
		return truth
	}
	return falsehood
}

// addYamlConf attempts to add conf to the file at filePath.
// If confKey exists in the file, the existing configs that match regEx are replaced
// and the result is written back to the file.
// If the confKey does not exist in the file, we append conf to the end of the file.
func addYamlConf(filePath, confKey, conf string, regEx *regexp.Regexp) error {
	logger.Infof("Adding configuration for %s to %s", confKey, filePath)
	config, err := getYamlConf(filePath, confKey)
	if err != nil {
		return err
	}

	if config == nil {
		if err = appendYaml(filePath, []byte("\n"+conf)); err != nil {
			return err
		}
	} else if err = replaceYaml(filePath, []byte(conf), regEx); err != nil {
		return err
	}

	return nil
}

// getYamlConf retrieves an interface mapped to a given conf
// within the file at the given filePath.
func getYamlConf(filePath, conf string) (interface{}, error) {
	sourceData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	var result map[string]interface{}
	if err = yaml.Unmarshal(sourceData, &result); err != nil {
		return "", err
	}

	return result[conf], nil
}

// appendYaml appends the yaml to the end of the file at
// the given filePath. The resulting data is written back to the file.
func appendYaml(filePath string, yaml []byte) error {
	sourceData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	sourceData = append(sourceData, yaml...)
	if err = ioutil.WriteFile(filePath, sourceData, 0644); err != nil {
		return err
	}

	return nil
}

// replaceYaml replaces all strings that match regEx in
// the file at filePath with yaml, and logs a warning
// if there are no matches. The resulting data is written back to the file.
func replaceYaml(filePath string, yaml []byte, regEx *regexp.Regexp) error {
	sourceData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	if ok := regEx.Match(sourceData); ok {
		sourceData = regEx.ReplaceAll(sourceData, yaml)
	} else {
		logger.Warnf("Failed to add the following configuation to %s:\n%sAs it already exists with a different configuration in %s", filePath, yaml, filePath)
	}

	if err = ioutil.WriteFile(filePath, sourceData, 0644); err != nil {
		return err
	}

	return nil
}
