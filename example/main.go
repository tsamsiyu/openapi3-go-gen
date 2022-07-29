package main

import (
	"fmt"
	"os"

	"openapi3-go-gen/cmd/codegen/app"
)

func main() {
	if err := os.RemoveAll("example/generated"); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := os.Mkdir("example/generated", 0777); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err := app.Run("example/testdata/openapi.yaml", "example/generated")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
