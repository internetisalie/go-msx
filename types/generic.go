// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package types

import (
	"fmt"
	"path"
	"reflect"
)

var parameterizedTypeNames = make(map[reflect.Type]string)
var customTypeNames = map[string]struct{}{
	"types.UUID": {},
	"types.Time": {},
}

func NewParameterizedStruct(structType reflect.Type, payloadField string, payload interface{}) reflect.Type {
	var structFields []reflect.StructField
	for i := 0; i < structType.NumField(); i++ {
		structField := structType.Field(i)
		if structField.Name == payloadField {
			if payload == nil {
				continue
			} else {
				structField.Type = reflect.TypeOf(payload)
			}
		}
		structFields = append(structFields, structField)
	}

	structName := GetTypeName(structType, false)
	payloadTypeName := GetInstanceTypeName(payload)

	parameterizedStructType := reflect.StructOf(structFields)
	parameterizedTypeNames[parameterizedStructType] = fmt.Sprintf("%s«%s»", structName, payloadTypeName)
	return parameterizedStructType
}

func GetInstanceTypeName(instance interface{}) string {
	if instance == nil {
		return "Void"
	}

	instanceType := reflect.TypeOf(instance)
	return GetTypeName(instanceType, true)
}

func GetTypeName(instanceType reflect.Type, root bool) string {
	if instanceType.Kind() == reflect.Ptr {
		instanceType = instanceType.Elem()
	}

	typeNamePrefix, typeNameSuffix, typeName := "", "", ""
	ok := false

	if typeName, ok = parameterizedTypeNames[instanceType]; ok {
		return typeName
	}

	if typeName, ok = GetNamedTypeName(instanceType); ok {
		if _, ok = customTypeNames[typeName]; ok {
			return typeName
		}
	}

	switch instanceType.Kind() {
	case reflect.Array, reflect.Slice:
		if root {
			// marked as an array in parent
			return GetTypeName(instanceType.Elem(), false)
		} else {
			typeNamePrefix = "List«"
			typeNameSuffix = "»"
			typeName = GetTypeName(instanceType.Elem(), false)
			return typeNamePrefix + typeName + typeNameSuffix
		}

	case reflect.Map:
		typeNamePrefix = "Map«" + GetTypeName(instanceType.Key(), false) + ","
		typeNameSuffix = "»"
		typeName = GetTypeName(instanceType.Elem(), false)
		return typeNamePrefix + typeName + typeNameSuffix
	}

	switch instanceType.Kind() {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Bool, reflect.String:
		return instanceType.Name()

	default:
		return typeName
	}
}

func GetNamedTypeName(instanceType reflect.Type) (string, bool) {
	instanceTypePackagePath := instanceType.PkgPath()
	instanceTypePackageName := ""
	if instanceTypePackagePath != "" {
		instanceTypePackageName = path.Base(instanceTypePackagePath)
		if instanceTypePackageName != "" {
			instanceTypePackageName += "."
		}
	} else {
		return "", false
	}
	return instanceTypePackageName + instanceType.Name(), true
}
