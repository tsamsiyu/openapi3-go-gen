package generator

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"sort"
	"text/template"

	"github.com/pkg/errors"
)

var (
	structTemplate *template.Template
)

func ReadTemplates(structTmplPath string) error {
	tplBytes, err := os.ReadFile(structTmplPath)
	if err != nil {
		return err
	}

	structTemplate, err = template.New("struct").Funcs(template.FuncMap{
		"NotNil": func(v interface{}) bool {
			reflval := reflect.ValueOf(v)
			return !reflval.IsNil()
		},
		"Deref": func(v interface{}) interface{} {
			reflval := reflect.ValueOf(v)

			if !reflval.IsValid() || reflval.IsNil() {
				return nil
			}

			if reflval.Kind() == reflect.Ptr {
				elem := reflval.Elem()
				return elem.Interface()
			}

			return v
		},
	}).Parse(string(tplBytes))
	if err != nil {
		return err
	}

	return nil
}

type sortingProp struct {
	props *[]Prop
}

func (p *sortingProp) Len() int {
	return len(*p.props)
}

func (p *sortingProp) Less(i, j int) bool {
	props := *p.props
	return props[i].Name > props[j].Name
}

func (p *sortingProp) Swap(i, j int) {
	props := *p.props
	props[i], props[j] = props[j], props[i]
}

type Generator struct {
}

func NewGenerator() *Generator {
	return &Generator{}
}

func (g *Generator) GenerateForModel(writer io.Writer, model *Model) error {
	if err := structTemplate.Execute(writer, model); err != nil {
		return err
	}

	return nil
}

func (g *Generator) GenerateToFile(models map[string]*Model, path string) error {
	for name, model := range models {
		sorted := &sortingProp{
			props: &model.Props,
		}

		sort.Sort(sorted)

		filename := modelToFilename(name) + ".go"

		fmt.Printf("Generating: %s\n", filename)

		file, err := os.Create(filepath.Join(path, filename))
		if err != nil {
			return err
		}
		defer file.Close()

		fileWriter := bufio.NewWriter(file)

		if err := g.GenerateForModel(fileWriter, model); err != nil {
			return err
		}

		if err := file.Sync(); err != nil {
			return err
		}

		if err := fileWriter.Flush(); err != nil {
			return err
		}

	}

	cmd := exec.Command("goimports", "-w", path)
	if err := cmd.Run(); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
