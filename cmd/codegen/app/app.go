package app

import (
	"context"
	"os/exec"
	"path/filepath"
	"time"

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
		return err
	}

	if err := doc.Validate(rootCtx); err != nil {
		return err
	}

	documentInspector := generator.NewFlattener(doc)

	flatSchemaRefs := documentInspector.Flatten()

	schemaResolver := generator.NewSchemaResolver(flatSchemaRefs)

	models := schemaResolver.Resolve()

	gen := generator.NewGenerator()

	err = gen.Generate(models, output)
	if err != nil {
		return err
	}

	generatedPath, err := filepath.Abs(output)
	if err != nil {
		return err
	}

	cmd := exec.Command("goimports", "-w", generatedPath)
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
