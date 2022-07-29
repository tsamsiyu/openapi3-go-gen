package main

import (
	"fmt"
	"os"

	"openapi3-go-gen/cmd/codegen/app"
)

const (
	src  = "example/testdata/openapi.yaml"
	dest = "example/generated"
)

func main() {
	if err := os.RemoveAll(dest); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := os.Mkdir(dest, 0777); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err := app.Run(src, dest)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
