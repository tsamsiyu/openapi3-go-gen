package generator

import (
	"path/filepath"
	"strings"

	"github.com/iancoleman/strcase"

	plur "github.com/gertd/go-pluralize"
	spec3 "github.com/getkin/kin-openapi/openapi3"
)

var (
	inflector = plur.NewClient()
)

func getCustomTypeSchemaRef(schemaRef *spec3.SchemaRef) *spec3.SchemaRef {
	var targetSchemaRef *spec3.SchemaRef

	if isArray(schemaRef.Value.Type) {
		targetSchemaRef = schemaRef.Value.Items
	} else {
		targetSchemaRef = schemaRef
	}

	if isScalar(targetSchemaRef.Value.Type) || isInterface(targetSchemaRef.Value) {
		return nil
	}

	return targetSchemaRef
}

func isInterface(schema *spec3.Schema) bool {
	if schema.OneOf != nil || schema.AnyOf != nil {
		return true
	}

	if schema.Type == "object" && (schema.Properties == nil) || len(schema.Properties) < 1 {
		return true
	}

	return false
}

func isScalar(tp string) bool {
	return tp == "string" || tp == "integer" || tp == "boolean" || tp == "float"
}

func isArray(tp string) bool {
	return tp == "array"
}

func modelToFilename(modelName string) string {
	return strcase.ToLowerCamel(modelName)
}

func getBaseFilename(path string) string {
	fileName := filepath.Base(path)
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}

func refToModelName(ref string) string {
	if strings.HasSuffix(ref, "yaml") || strings.HasSuffix(ref, "yml") {
		return strcase.ToCamel(getBaseFilename(ref))
	}

	parts := strings.Split(ref, "/")
	return parts[len(parts)-1]
}

func propToModelName(prop string) string {
	return strcase.ToCamel(inflector.Singular(prop))
}

func propName(prop string) string {
	return strcase.ToCamel(prop)
}

func embeddedObjectToModelName(schemaName string, prop string) string {
	return strcase.ToCamel(schemaName + "_" + inflector.Singular(prop))
}
