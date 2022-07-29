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

func (f *Flattener) Flatten() map[string]*spec3.SchemaRef {
	flatSchemaRefs := make(map[string]*spec3.SchemaRef)

	for schemaName, schema := range f.doc.Components.Schemas {
		f.collectCustomSchemaRef("", schemaName, schema, flatSchemaRefs)
	}

	for schemaName, schema := range f.doc.Components.Schemas {
		f.collectDeepCustomPropsSchemaRef(schemaName, schema, flatSchemaRefs)
	}

	return flatSchemaRefs
}

func (f *Flattener) collectDeepCustomPropsSchemaRef(schemaName string, schemaRef *spec3.SchemaRef, flatSchemaRefs map[string]*spec3.SchemaRef) {
	custom := getCustomTypeSchemaRef(schemaRef)
	if custom == nil {
		return
	}

	var manyRefs []*spec3.SchemaRef

	if schemaRef.Value.AllOf != nil {
		manyRefs = schemaRef.Value.AllOf
	}

	if schemaRef.Value.OneOf != nil {
		manyRefs = schemaRef.Value.OneOf
	}

	if schemaRef.Value.AnyOf != nil {
		manyRefs = schemaRef.Value.AnyOf
	}

	for propName, propSchema := range custom.Value.Properties {
		propSchemaName := f.collectCustomSchemaRef(schemaName, propName, propSchema, flatSchemaRefs)
		if propSchemaName != "" {
			f.collectDeepCustomPropsSchemaRef(propSchemaName, propSchema, flatSchemaRefs)
		}
	}

	if manyRefs != nil && len(manyRefs) > 1 {
		for _, elementSchema := range manyRefs {
			f.collectDeepCustomPropsSchemaRef(schemaName, elementSchema, flatSchemaRefs)
		}
	}
}

func (f *Flattener) collectCustomSchemaRef(
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
