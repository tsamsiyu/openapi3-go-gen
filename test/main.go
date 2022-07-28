package main

import (
	"fmt"
	"os"

	"openapi3-go-gen/cmd/codegen/app"
)

func main() {
	err := app.Run("test/testdata/openapi.yaml", "test/generated")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
