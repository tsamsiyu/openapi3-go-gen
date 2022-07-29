# Openapi3 Go Gen

    Tool for generating Go models for schemas defined in components section of openapi specification.

### Details:
- generates models
- generates validations
- meets oneOf and anyOf as interface type
- correctly handles allOf
- all files are generated into a single folder

Feel free to check `example` folder to see a generated result

### Usage

Be sure there is `openapi.yaml` file and `generated` folder in your current directory, then type:

> docker run --rm -v "$PWD:/usr/run" tsamsiyu/openapi3-go-gen --input=/usr/run/openapi.yaml --output=/usr/run/generated 
