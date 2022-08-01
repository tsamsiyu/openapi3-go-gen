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

	_ = os.RemoveAll(genPath)
	_ = os.Remove(specPath)

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
	l.IsExternalRefsAllowed = true

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

func TestSimplest(t *testing.T) {
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

func TestSimplestRef(t *testing.T) {
	beforeTest(t)

	schemasYaml := `
Foo:
  type: object
  properties:
    bar:
      $ref: "#/components/schemas/Bar"
Bar:
  type: object
  properties:
    name:
      type: string
`

	expectedFoo := strings.TrimPrefix(`
package openapi

type Foo struct {
	Bar Bar
}

func (instance *Foo) Validate() error {
	return nil
}
`, "\n")

	expectedBar := strings.TrimPrefix(`
package openapi

type Bar struct {
	Name string
}

func (instance *Bar) Validate() error {
	return nil
}
`, "\n")

	err := generate(schemasYaml)
	require.NoError(t, err)

	foo, err := readGoFile("foo.go")
	require.NoError(t, err)

	bar, err := readGoFile("bar.go")
	require.NoError(t, err)

	require.Equal(t, expectedFoo, foo)
	require.Equal(t, expectedBar, bar)
}

func TestSimplestNested(t *testing.T) {
	beforeTest(t)

	schemasYaml := `
Foo:
  type: object
  properties:
    bar:
      type: object
      properties:
        name:
          type: string
`

	expectedFoo := strings.TrimPrefix(`
package openapi

type Foo struct {
	Bar FooBar
}

func (instance *Foo) Validate() error {
	return nil
}
`, "\n")

	expectedBar := strings.TrimPrefix(`
package openapi

type FooBar struct {
	Name string
}

func (instance *FooBar) Validate() error {
	return nil
}
`, "\n")

	err := generate(schemasYaml)
	require.NoError(t, err)

	foo, err := readGoFile("foo.go")
	require.NoError(t, err)

	bar, err := readGoFile("foo_bar.go")
	require.NoError(t, err)

	require.Equal(t, expectedFoo, foo)
	require.Equal(t, expectedBar, bar)
}
