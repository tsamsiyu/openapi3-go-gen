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

func TestSimplestObject(t *testing.T) {
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

func TestRef(t *testing.T) {
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

func TestNested(t *testing.T) {
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

func TestOneOf(t *testing.T) {
	beforeTest(t)

	schemasYaml := `
Foo:
  type: object
  properties:
    bar:
      type: string
    baz:
      oneOf:
        - $ref: "#/components/schemas/Baz"
        - type: string
Baz:
  type: object
  properties:
    name:
      type: string
`

	expectedFoo := strings.TrimPrefix(`
package openapi

type Foo struct {
	Baz interface{}
	Bar string
}

func (instance *Foo) Validate() error {
	return nil
}
`, "\n")

	expectedBaz := strings.TrimPrefix(`
package openapi

type Baz struct {
	Name string
}

func (instance *Baz) Validate() error {
	return nil
}
`, "\n")

	err := generate(schemasYaml)
	require.NoError(t, err)

	foo, err := readGoFile("foo.go")
	require.NoError(t, err)

	bar, err := readGoFile("baz.go")
	require.NoError(t, err)

	require.Equal(t, expectedFoo, foo)
	require.Equal(t, expectedBaz, bar)
}

func TestOneOfWithJustRef(t *testing.T) {
	beforeTest(t)

	schemasYaml := `
Foo:
  type: object
  properties:
    bar:
      type: string
    baz:
      oneOf:
        - $ref: "#/components/schemas/Baz"
Baz:
  type: object
  properties:
    name:
      type: string
`

	expectedFoo := strings.TrimPrefix(`
package openapi

type Foo struct {
	Baz interface{}
	Bar string
}

func (instance *Foo) Validate() error {
	return nil
}
`, "\n")

	expectedBaz := strings.TrimPrefix(`
package openapi

type Baz struct {
	Name string
}

func (instance *Baz) Validate() error {
	return nil
}
`, "\n")

	err := generate(schemasYaml)
	require.NoError(t, err)

	foo, err := readGoFile("foo.go")
	require.NoError(t, err)

	bar, err := readGoFile("baz.go")
	require.NoError(t, err)

	require.Equal(t, expectedFoo, foo)
	require.Equal(t, expectedBaz, bar)
}

func TestAnyOf(t *testing.T) {
	beforeTest(t)

	schemasYaml := `
Foo:
  type: object
  properties:
    bar:
      type: string
    baz:
      anyOf:
        - $ref: "#/components/schemas/Baz"
        - type: string
Baz:
  type: object
  properties:
    name:
      type: string
`

	expectedFoo := strings.TrimPrefix(`
package openapi

type Foo struct {
	Baz interface{}
	Bar string
}

func (instance *Foo) Validate() error {
	return nil
}
`, "\n")

	expectedBaz := strings.TrimPrefix(`
package openapi

type Baz struct {
	Name string
}

func (instance *Baz) Validate() error {
	return nil
}
`, "\n")

	err := generate(schemasYaml)
	require.NoError(t, err)

	foo, err := readGoFile("foo.go")
	require.NoError(t, err)

	bar, err := readGoFile("baz.go")
	require.NoError(t, err)

	require.Equal(t, expectedFoo, foo)
	require.Equal(t, expectedBaz, bar)
}

func TestNullable(t *testing.T) {
	beforeTest(t)

	schemasYaml := `
Foo:
  type: object
  properties:
    bar:
      type: string
    baz:
      type: number
      nullable: true
`

	expectedFoo := strings.TrimPrefix(`
package openapi

type Foo struct {
	Baz *float64
	Bar string
}

func (instance *Foo) Validate() error {
	return nil
}
`, "\n")

	err := generate(schemasYaml)
	require.NoError(t, err)

	foo, err := readGoFile("foo.go")
	require.NoError(t, err)

	require.Equal(t, expectedFoo, foo)
}

func TestArrayWithRef(t *testing.T) {
	beforeTest(t)

	schemasYaml := `
Foo:
  type: object
  properties:
    name:
      type: string
    bars:
      type: array
      items:
        $ref: "#/components/schemas/Bar"
Bar:
  type: object
  properties:
    age:
      type: integer
`

	expectedFoo := strings.TrimPrefix(`
package openapi

type Foo struct {
	Name string
	Bars []Bar
}

func (instance *Foo) Validate() error {
	return nil
}
`, "\n")

	expectedBar := strings.TrimPrefix(`
package openapi

type Bar struct {
	Age int
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

func TestArrayWithInt(t *testing.T) {
	beforeTest(t)

	schemasYaml := `
Foo:
  type: object
  properties:
    name:
      type: string
    bars:
      type: array
      items:
        type: integer
`

	expectedFoo := strings.TrimPrefix(`
package openapi

type Foo struct {
	Name string
	Bars []int
}

func (instance *Foo) Validate() error {
	return nil
}
`, "\n")

	err := generate(schemasYaml)
	require.NoError(t, err)

	foo, err := readGoFile("foo.go")
	require.NoError(t, err)

	require.Equal(t, expectedFoo, foo)
}

func TestArrayWithNestedObject(t *testing.T) {
	beforeTest(t)

	schemasYaml := `
Foo:
  type: object
  properties:
    name:
      type: string
    bars:
      type: array
      items:
        type: object
        properties:
          zoo:
            type: string
`

	expectedFoo := strings.TrimPrefix(`
package openapi

type Foo struct {
	Name string
	Bars []FooBar
}

func (instance *Foo) Validate() error {
	return nil
}
`, "\n")

	expectedBar := strings.TrimPrefix(`
package openapi

type FooBar struct {
	Zoo string
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

func TestArrayWithNullableType(t *testing.T) {
	beforeTest(t)

	schemasYaml := `
Foo:
  type: object
  properties:
    name:
      type: string
    bars:
      type: array
      items:
        type: string
        nullable: true
`

	expectedFoo := strings.TrimPrefix(`
package openapi

type Foo struct {
	Name string
	Bars []string
}

func (instance *Foo) Validate() error {
	return nil
}
`, "\n")

	err := generate(schemasYaml)
	require.NoError(t, err)

	foo, err := readGoFile("foo.go")
	require.NoError(t, err)

	require.Equal(t, expectedFoo, foo)
}

func TestNullableArray(t *testing.T) {
	beforeTest(t)

	schemasYaml := `
Foo:
  type: object
  properties:
    name:
      type: string
    bars:
      type: array
      nullable: true
      items:
        type: string
`

	expectedFoo := strings.TrimPrefix(`
package openapi

type Foo struct {
	Name string
	Bars []string
}

func (instance *Foo) Validate() error {
	return nil
}
`, "\n")

	err := generate(schemasYaml)
	require.NoError(t, err)

	foo, err := readGoFile("foo.go")
	require.NoError(t, err)

	require.Equal(t, expectedFoo, foo)
}

func TestAllOfWithRefAndObject(t *testing.T) {
	beforeTest(t)

	schemasYaml := `
Foo:
  type: object
  properties:
    name:
      type: string
    plum:
      allOf:
        - $ref: "#/components/schemas/Bar"
        - type: object
          properties:
            kek:
              type: string
              nullable: true
Bar:
  type: object
  properties:
    bazzer:
      type: string
`

	expectedFoo := strings.TrimPrefix(`
package openapi

type Foo struct {
	Plum FooPlum
	Name string
}

func (instance *Foo) Validate() error {
	return nil
}
`, "\n")

	expectedBar := strings.TrimPrefix(`
package openapi

type Bar struct {
	Bazzer string
}

func (instance *Bar) Validate() error {
	return nil
}
`, "\n")

	expectedPlum := strings.TrimPrefix(`
package openapi

type FooPlum struct {
	Kek    *string
	Bazzer string
}

func (instance *FooPlum) Validate() error {
	return nil
}
`, "\n")

	err := generate(schemasYaml)
	require.NoError(t, err)

	foo, err := readGoFile("foo.go")
	require.NoError(t, err)

	bar, err := readGoFile("bar.go")
	require.NoError(t, err)

	plum, err := readGoFile("foo_plum.go")
	require.NoError(t, err)

	require.Equal(t, expectedFoo, foo)
	require.Equal(t, expectedBar, bar)
	require.Equal(t, expectedPlum, plum)
}

func TestAllOfWithJustOneRef(t *testing.T) {
	beforeTest(t)

	schemasYaml := `
Foo:
  type: object
  properties:
    name:
      type: string
    plum:
      allOf:
        - $ref: "#/components/schemas/Bar"
Bar:
  type: object
  properties:
    bazzer:
      type: string
`

	expectedFoo := strings.TrimPrefix(`
package openapi

type Foo struct {
	Plum Bar
	Name string
}

func (instance *Foo) Validate() error {
	return nil
}
`, "\n")

	expectedBar := strings.TrimPrefix(`
package openapi

type Bar struct {
	Bazzer string
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

func TestAllOfWithJustOneNestedObject(t *testing.T) {
	beforeTest(t)

	schemasYaml := `
Foo:
  type: object
  properties:
    name:
      type: string
    plum:
      allOf:
        - type: object
          properties:
            is_agree:
              type: boolean
`

	expectedFoo := strings.TrimPrefix(`
package openapi

type Foo struct {
	Plum FooPlum
	Name string
}

func (instance *Foo) Validate() error {
	return nil
}
`, "\n")

	expectedPlum := strings.TrimPrefix(`
package openapi

type FooPlum struct {
	IsAgree bool
}

func (instance *FooPlum) Validate() error {
	return nil
}
`, "\n")

	err := generate(schemasYaml)
	require.NoError(t, err)

	foo, err := readGoFile("foo.go")
	require.NoError(t, err)

	plum, err := readGoFile("foo_plum.go")
	require.NoError(t, err)

	require.Equal(t, expectedFoo, foo)
	require.Equal(t, expectedPlum, plum)
}

func TestRequiredStringAndNullableString(t *testing.T) {
	beforeTest(t)

	schemasYaml := `
Foo:
  type: object
  required: [name, last_name]
  properties:
    name:
      type: string
    last_name:
      type: string
      nullable: true
`

	expectedFoo := strings.TrimPrefix(`
package openapi

import (
	"errors"
)

type Foo struct {
	Name     string
	LastName *string
}

func (instance *Foo) Validate() error {
	if instance.Name == "" {
		return errors.New("Value for field Name must be not empty")
	}
	if instance.LastName == nil {
		return errors.New("Value for field LastName must be present")
	}
	return nil
}
`, "\n")

	err := generate(schemasYaml)
	require.NoError(t, err)

	foo, err := readGoFile("foo.go")
	require.NoError(t, err)

	require.Equal(t, expectedFoo, foo)
}

func TestMinMaxLength(t *testing.T) {
	beforeTest(t)

	schemasYaml := `
Foo:
  type: object
  properties:
    name:
      type: string
      maxLength: 10
      minLength: 3
`

	expectedFoo := strings.TrimPrefix(`
package openapi

import (
	"errors"
)

type Foo struct {
	Name string
}

func (instance *Foo) Validate() error {
	if len(instance.Name) > 10 {
		return errors.New("Field Name size should not be greater than 10")
	}
	if len(instance.Name) < 3 {
		return errors.New("Field Name size should not be less than 3")
	}
	return nil
}
`, "\n")

	err := generate(schemasYaml)
	require.NoError(t, err)

	foo, err := readGoFile("foo.go")
	require.NoError(t, err)

	require.Equal(t, expectedFoo, foo)
}

func TestMinMax(t *testing.T) {
	beforeTest(t)

	schemasYaml := `
Foo:
  type: object
  properties:
    name:
      type: integer
      maximum: 10
      minimum: 3
`

	expectedFoo := strings.TrimPrefix(`
package openapi

import (
	"errors"
)

type Foo struct {
	Name int
}

func (instance *Foo) Validate() error {
	if instance.Name > 10 {
		return errors.New("Field Name should not be greater than 10")
	}
	if instance.Name < 3 {
		return errors.New("Field Name should not be less than 3")
	}
	return nil
}
`, "\n")

	err := generate(schemasYaml)
	require.NoError(t, err)

	foo, err := readGoFile("foo.go")
	require.NoError(t, err)

	require.Equal(t, expectedFoo, foo)
}

func TestExclusiveMinMax(t *testing.T) {
	beforeTest(t)

	schemasYaml := `
Foo:
  type: object
  properties:
    name:
      type: integer
      maximum: 10
      minimum: 3
      exclusiveMinimum: true
      exclusiveMaximum: true
`

	expectedFoo := strings.TrimPrefix(`
package openapi

import (
	"errors"
)

type Foo struct {
	Name int
}

func (instance *Foo) Validate() error {
	if instance.Name >= 10 {
		return errors.New("Field Name should not be greater or equal than 10")
	}
	if instance.Name <= 3 {
		return errors.New("Field Name should not be less or equal than 3")
	}
	return nil
}
`, "\n")

	err := generate(schemasYaml)
	require.NoError(t, err)

	foo, err := readGoFile("foo.go")
	require.NoError(t, err)

	require.Equal(t, expectedFoo, foo)
}

func TestPattern(t *testing.T) {
	beforeTest(t)

	schemasYaml := `
Foo:
  type: object
  properties:
    name:
      type: string
      pattern: ^\d{3}-\d{2}-\d{4}$
`

	expectedFoo := strings.TrimPrefix(`
package openapi

import (
	"errors"
	"regexp"
)

type Foo struct {
	Name string
}

func (instance *Foo) Validate() error {
	if match, _ := regexp.MatchString(`+"`"+`^\d{3}-\d{2}-\d{4}$`+"`"+`, instance.Name); !match {
		return errors.New("Field Name is not formatted correctly")
	}
	return nil
}
`, "\n")

	err := generate(schemasYaml)
	require.NoError(t, err)

	foo, err := readGoFile("foo.go")
	require.NoError(t, err)

	require.Equal(t, expectedFoo, foo)
}

func TestEnum(t *testing.T) {
	beforeTest(t)

	schemasYaml := `
Foo:
  type: object
  properties:
    name:
      type: string
      enum: [Katty, Petty]
    level:
      type: number
      enum: [1.1, 2.2]
`

	expectedFoo := strings.TrimPrefix(`
package openapi

import (
	"errors"
)

type Foo struct {
	Name  string
	Level float64
}

func (instance *Foo) Validate() error {
	containsName := false
	enumName := []string{"Katty", "Petty"}
	for _, v := range enumName {
		if v == instance.Name {
			containsName = true
			break
		}
	}

	if !containsName {
		return errors.New("Value for field Name is not allowed")
	}
	containsLevel := false
	enumLevel := []float64{"1.1", "2.2"}
	for _, v := range enumLevel {
		if v == instance.Level {
			containsLevel = true
			break
		}
	}

	if !containsLevel {
		return errors.New("Value for field Level is not allowed")
	}
	return nil
}
`, "\n")

	err := generate(schemasYaml)
	require.NoError(t, err)

	foo, err := readGoFile("foo.go")
	require.NoError(t, err)

	require.Equal(t, expectedFoo, foo)
}
