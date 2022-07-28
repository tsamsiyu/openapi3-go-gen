package generator

import (
	spec3 "github.com/getkin/kin-openapi/openapi3"
)

type Flattener struct {
	doc *spec3.T
}

func NewFlattener(doc *spec3.T) *Flattener {
	return &Flattener{doc: doc}
}

func (i *Flattener) Flatten() map[string]*spec3.SchemaRef {
	flatSchemaRefs := make(map[string]*spec3.SchemaRef)

	for schemaName, schema := range i.doc.Components.Schemas {
		i.collectSchema("", schemaName, schema, flatSchemaRefs)
	}

	for schemaName, schema := range i.doc.Components.Schemas {
		i.deepFlatSchemaRef(schemaName, schema, flatSchemaRefs)
	}

	return flatSchemaRefs
}

func (i *Flattener) deepFlatSchemaRef(schemaName string, schemaRef *spec3.SchemaRef, flatSchemas map[string]*spec3.SchemaRef) {
	for propName, propSchema := range schemaRef.Value.Properties {
		propSchemaName := i.collectSchema(schemaName, propName, propSchema, flatSchemas)
		if propSchemaName != "" {
			i.deepFlatSchemaRef(propSchemaName, propSchema, flatSchemas)
		}
	}
}

func (i *Flattener) collectSchema(
	parentName string,
	name string,
	schema *spec3.SchemaRef,
	flatSchemaRefs map[string]*spec3.SchemaRef,
) string {
	custom := getCustomTypeSchemaRef(schema)
	if custom == nil {
		return ""
	}

	var modelName string

	if custom.Ref != "" {
		modelName = refToModelName(custom.Ref)
	} else {
		if parentName != "" {
			modelName = embeddedObjectToModelName(parentName, name)
		} else {
			modelName = propToModelName(name)
		}
	}

	flatSchemaRefs[modelName] = custom

	return modelName
}

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
