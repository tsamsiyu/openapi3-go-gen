package generator

import (
	"errors"
	"fmt"

	spec3 "github.com/getkin/kin-openapi/openapi3"
)

const (
	GeneratedFilesPkgName = "openapi"
)

type GoType struct {
	Name       string
	IsNullable bool
	IsPtr      bool
}

type Prop struct {
	*spec3.Schema

	GoType     *GoType
	Name       string
	IsRequired bool
}

type Model struct {
	PkgName string
	Name    string
	Props   []Prop
}

type SchemaResolver struct {
	data map[string]*spec3.SchemaRef
}

func NewSchemaResolver(data map[string]*spec3.SchemaRef) *SchemaResolver {
	return &SchemaResolver{
		data: data,
	}
}

func (r *SchemaResolver) Resolve() map[string]*Model {
	models := make(map[string]*Model)

	for name, schemaRef := range r.data {
		models[name] = &Model{
			PkgName: GeneratedFilesPkgName,
			Name:    name,
			Props:   r.buildProps(name, schemaRef),
		}
	}

	return models
}

func (r *SchemaResolver) buildProps(name string, schemaRef *spec3.SchemaRef) []Prop {
	props := make([]Prop, 0)

	if schemaRef.Value.AllOf != nil {
		for _, elementSchemaRef := range schemaRef.Value.AllOf {
			elementProps := r.buildProps(name, elementSchemaRef)
			props = append(props, elementProps...)
		}
	} else {
		for propName, propSchemaRef := range schemaRef.Value.Properties {
			prop := r.mapSchemaRefToProp(name, schemaRef.Value, propName, propSchemaRef)
			props = append(props, *prop)
		}
	}

	return props
}

func (r *SchemaResolver) findSchema(name string) *spec3.SchemaRef {
	var res *spec3.SchemaRef
	for n, v := range r.data {
		if n == name {
			res = v
			break
		}
	}

	return res
}

func (r *SchemaResolver) mapSchemaRefToProp(parentName string, parentSchema *spec3.Schema, name string, schemaRef *spec3.SchemaRef) *Prop {
	var prop *Prop

	custom := getCustomTypeSchemaRef(schemaRef)

	if custom == nil {
		prop = &Prop{
			Schema:     schemaRef.Value,
			Name:       propName(name),
			GoType:     mapSimpleSchema2GoType(schemaRef.Value),
			IsRequired: isPropRequired(parentSchema.Required, name),
		}
	} else {
		var modelName string

		if custom.Ref != "" {
			modelName = refToModelName(custom.Ref)
		} else {
			modelName = embeddedObjectToModelName(parentName, name)
		}

		referenced := r.findSchema(modelName)
		if referenced == nil {
			msg := fmt.Sprintf("There is no component [%s] found by ref %s", modelName, custom.Ref)
			panic(errors.New(msg).(any))
		}

		prop = &Prop{
			Schema:     custom.Value,
			Name:       propName(name),
			GoType:     mapCustomSchemaToGoType(modelName, schemaRef.Value),
			IsRequired: isPropRequired(parentSchema.Required, name),
		}
	}

	return prop
}

func mapCustomSchemaToGoType(typeName string, schema *spec3.Schema) *GoType {
	if schema.Type == "array" {
		return &GoType{
			Name:       fmt.Sprintf("[]%s", typeName),
			IsNullable: true,
			IsPtr:      false,
		}
	}

	if schema.Nullable {
		return &GoType{
			Name:       "*" + typeName,
			IsNullable: true,
			IsPtr:      true,
		}
	}

	return &GoType{
		Name:       typeName,
		IsNullable: false,
		IsPtr:      false,
	}
}

func mapScalarType2GoType(schema *spec3.Schema) *GoType {
	var goTypeStr string

	switch schema.Type {
	case "integer":
		goTypeStr = "int"
	case "float":
		goTypeStr = "float64"
	case "boolean":
		goTypeStr = "bool"
	case "string":
		goTypeStr = "string"
	}

	if goTypeStr == "" {
		return nil
	}

	if schema.Nullable {
		return &GoType{
			Name:       "*" + goTypeStr,
			IsNullable: schema.Nullable,
			IsPtr:      true,
		}
	}

	return &GoType{
		Name:       goTypeStr,
		IsNullable: false,
		IsPtr:      false,
	}
}

func mapSimpleSchema2GoType(schema *spec3.Schema) *GoType {
	scalarGoType := mapScalarType2GoType(schema)
	if scalarGoType != nil {
		return scalarGoType
	}

	if schema.OneOf != nil || schema.AnyOf != nil || schema.Type == "object" {
		return &GoType{
			Name:       "interface{}",
			IsNullable: true,
			IsPtr:      false,
		}
	}

	if schema.AllOf != nil && len(schema.AllOf) == 1 {
		return mapSimpleSchema2GoType(schema.AllOf[0].Value)
	}

	if schema.Type == "array" {
		scalarGoType := mapScalarType2GoType(schema.Items.Value)
		if scalarGoType != nil {
			return scalarGoType
		}

		if schema.Items.Value.OneOf != nil || schema.Items.Value.AnyOf != nil || schema.Items.Value.Type == "object" {
			return &GoType{
				Name:       "[]interface{}",
				IsNullable: true,
				IsPtr:      false,
			}
		}

		if schema.Items.Value.AllOf != nil && len(schema.Items.Value.AllOf) == 1 {
			return mapSimpleSchema2GoType(schema.Items.Value.AllOf[0].Value)
		}

		panic(errors.New(fmt.Sprintf("Not a simple array type provided: %s", schema.Items.Value.Type)).(any))
	}

	panic(errors.New(fmt.Sprintf("Not a simple type provided: %s", schema.Type)).(any))
}

func isPropRequired(objRequired []string, propName string) bool {
	res := false

	for _, v := range objRequired {
		if v == propName {
			res = true
			break
		}
	}

	return res
}
