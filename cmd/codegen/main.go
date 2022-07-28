package main

import (
	"flag"
	"fmt"
	"log"
	"openapi3-go-gen/cmd/codegen/app"
	"os"
)

func main() {
	input := flag.String("input", "", "Path to openapi.yaml or openapi.json")
	output := flag.String("output", "", "Path to where generated files will be located")
	flag.Parse()

	if *input == "" {
		log.Fatalln("'input' flag must be provided")
	}

	if *output == "" {
		log.Fatalln("'output' flag must be provided")
	}

	if _, err := os.Stat(*input); os.IsNotExist(err) {
		log.Fatalf("File %s does not exist\n", *input)
	}

	if _, err := os.Stat(*output); os.IsNotExist(err) {
		log.Fatalf("Directory %s does not exist\n", *output)
	}

	err := app.Run(*input, *output)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
