package test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"openapi3-go-gen/pkg/generator"

	"github.com/kr/text"
	"github.com/stretchr/testify/require"

	spec3 "github.com/getkin/kin-openapi/openapi3"
)

const oasLayout = `
openapi: "3.0.0"
info:
  title: "Test"
  version: "1.0.0"
paths: {}
components:
  schemas: %s
`

func beforeTest(t *testing.T) {
	tmplPath, err := filepath.Abs("../pkg/generator/templates/struct.tmpl")
	require.NoError(t, err)

	genPath, err := filepath.Abs("gen")
	require.NoError(t, err)

	specPath, err := filepath.Abs("oas.yml")
	require.NoError(t, err)

	err = generator.ReadTemplates(tmplPath)
	require.NoError(t, err)

	err = os.RemoveAll(genPath)
	require.NoError(t, err)

	err = os.Remove(specPath)
	require.NoError(t, err)

	err = os.Mkdir("gen", 0777)
	require.NoError(t, err)
}

func generate(yml string) error {
	oasStr := fmt.Sprintf(oasLayout, text.Indent(yml, strings.Repeat("  ", 2)))

	err := os.WriteFile("oas.yml", []byte(oasStr), 0777)
	if err != nil {
		return err
	}

	l := spec3.NewLoader()

	doc, err := l.LoadFromData([]byte(oasStr))
	if err != nil {
		return err
	}

	flattener := generator.NewFlattener(doc)

	flatSchemaRefs := flattener.Flatten()

	schemaResolver := generator.NewSchemaResolver(flatSchemaRefs)

	models := schemaResolver.Resolve()

	gen := generator.NewGenerator()

	err = gen.GenerateToFile(models, "gen")
	if err != nil {
		return err
	}

	return nil
}

func readGoFile(filename string) (string, error) {
	file, err := os.ReadFile(fmt.Sprintf("gen/%s", filename))
	return string(file), err
}

func TestOne(t *testing.T) {
	beforeTest(t)

	schemasYaml := `
Foo:
  type: object
  properties:
    str:
      type: string
    num:
      type: number
    int:
      type: integer
`

	expected := strings.TrimPrefix(`
package openapi

type Foo struct {
	Str string
	Num float64
	Int int
}

func (instance *Foo) Validate() error {
	return nil
}
`, "\n")

	err := generate(schemasYaml)
	require.NoError(t, err)

	f, err := readGoFile("foo.go")
	require.NoError(t, err)

	require.Equal(t, expected, f)
}
