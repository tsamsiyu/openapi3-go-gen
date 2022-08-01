package app

import (
	"context"
	"os"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
	"openapi3-go-gen/pkg/generator"

	spec3 "github.com/getkin/kin-openapi/openapi3"
)

func Run(input string, output string) error {
	rootCtx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if err := readTemplates(); err != nil {
		return errors.WithStack(err)
	}

	l := spec3.NewLoader()
	l.IsExternalRefsAllowed = true

	doc, err := l.LoadFromFile(input)
	if err != nil {
		return errors.Wrapf(err, "failed while loading openapi spec")
	}

	if err := doc.Validate(rootCtx); err != nil {
		return errors.Wrapf(err, "failed while validating openapi spec")
	}

	flattener := generator.NewFlattener(doc)

	flatSchemaRefs := flattener.Flatten()

	schemaResolver := generator.NewSchemaResolver(flatSchemaRefs)

	models := schemaResolver.Resolve()

	gen := generator.NewGenerator()

	if err := gen.GenerateToFile(models, output); err != nil {
		return err
	}

	return nil
}

func readTemplates() error {
	var structTmplPath string

	templatesFolder := os.Getenv("CODEGEN_TEMPLATES_FOLDER")
	if templatesFolder == "" {
		structTmplPath = "pkg/generator/templates/struct.tmpl"
	} else {
		structTmplPath = filepath.Join(templatesFolder, "struct.tmpl")
	}

	return generator.ReadTemplates(structTmplPath)
}
