# Openapi3 Go Gen

    Tool for generating Go models for schemas defined in components section of openapi specification.

### Details:
- generates models
- generate validations
- meet oneOf and anyOf as interface type
- correctly handles allOf
- all files are generated into single folder

Feel free to check `test/testdata` along with `test/generated` folders as an example.

### Usage

Be sure you have `openapi.yaml` and `generated` folder in your active directory, then type:

> docker run --rm -v "$PWD:/usr/run" tsamsiyu/openapi3-go-gen --input=/usr/run/openapi.yaml --output=/usr/run/generated 
