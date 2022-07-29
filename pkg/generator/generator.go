package generator

import (
	"bufio"
	"os"
	"path/filepath"
	"reflect"
	"text/template"
)

var (
	structTemplate *template.Template
)

func init() {
	var structTmplPath string

	templatesFolder := os.Getenv("CODEGEN_TEMPLATES_FOLDER")
	if templatesFolder == "" {
		structTmplPath = "pkg/generator/templates/struct.tmpl"
	} else {
		structTmplPath = filepath.Join(templatesFolder, "struct.tmpl")
	}

	tplBytes, err := os.ReadFile(structTmplPath)
	if err != nil {
		panic(err.(any))
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
		panic(err.(any))
	}
}

type Generator struct {
}

func NewGenerator() *Generator {
	return &Generator{}
}

func (g *Generator) Generate(models map[string]*Model, path string) error {
	for name, model := range models {
		file, err := os.Create(filepath.Join(path, modelToFilename(name)+".go"))
		if err != nil {
			return err
		}
		defer file.Close()

		fileWriter := bufio.NewWriter(file)

		if err := structTemplate.Execute(fileWriter, model); err != nil {
			return err
		}

		if err := file.Sync(); err != nil {
			return err
		}

		if err := fileWriter.Flush(); err != nil {
			return err
		}
	}

	return nil
}
