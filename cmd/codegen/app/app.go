package app

import (
	"context"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
	"openapi3-go-gen/pkg/generator"

	spec3 "github.com/getkin/kin-openapi/openapi3"
)

func Run(input string, output string) error {
	rootCtx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	l := spec3.NewLoader()
	l.IsExternalRefsAllowed = true
	doc, err := l.LoadFromFile(input)
	if err != nil {
		return errors.WithStack(err)
	}

	if err := doc.Validate(rootCtx); err != nil {
		return errors.WithStack(err)
	}

	documentInspector := generator.NewFlattener(doc)

	flatSchemaRefs := documentInspector.Flatten()

	schemaResolver := generator.NewSchemaResolver(flatSchemaRefs)

	models := schemaResolver.Resolve()

	gen := generator.NewGenerator()

	err = gen.Generate(models, output)
	if err != nil {
		return errors.WithStack(err)
	}

	generatedPath, err := filepath.Abs(output)
	if err != nil {
		return errors.WithStack(err)
	}

	cmd := exec.Command("goimports", "-w", generatedPath)
	if err := cmd.Run(); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
