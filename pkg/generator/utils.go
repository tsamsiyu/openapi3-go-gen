package generator

import (
	spec3 "github.com/getkin/kin-openapi/openapi3"
	"path/filepath"
	"strings"

	"github.com/iancoleman/strcase"

	plur "github.com/gertd/go-pluralize"
)

var (
	inflector = plur.NewClient()
)

func isInterface(schema *spec3.Schema) bool {
	return schema.OneOf != nil || schema.AnyOf != nil
}

func isScalar(tp string) bool {
	return tp == "string" || tp == "integer" || tp == "boolean" || tp == "float"
}

func isArray(tp string) bool {
	return tp == "array"
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
