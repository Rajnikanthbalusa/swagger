package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"

	"gopkg.in/yaml.v2"
)

type SwaggerDefinition struct {
	Type string `yaml:"type"`
	// TODO: Properties should probably be more tightly typed
	Properties map[string]map[string]string `yaml:"properties"`
}

// TODO: I don't like that we have to pass a Generator in here...
// TODO: Should this even be object oriented?
func (d SwaggerDefinition) Printf(g *Generator, name string) {
	// TODO: Handle different types...

	g.Printf("type %s struct {\n", name)
	for key, value := range d.Properties {
		// TODO: Error handling? Or verify that earlier??? Also have switch?
		propType := value["type"]
		structType := ""
		if propType == "string" {
			structType = "string"

		} else if propType == "integer" {
			// TODO: Distinguish between int32 and other possibilities...
			structType = "int"
		} else {
			panic(fmt.Sprintf("Type %s not supported", propType))
		}
		// TODO: Upper case
		g.Printf("\t%s %s `json:\"%s\"`\n", key, structType, key)
	}
	g.Printf("}\n\n")
}

type SwaggerOperation struct {
	OperationID string                     `yaml:"operationId"`
	Description string                     `yaml:"description"`
	Responses   map[string]SwaggerResponse `yaml:"responses"`
	// TODO: Will need to add parameters...
}

type SwaggerResponse struct {
	Description string `yaml:"description"`
	// TODO: Add more types to schema???
	Schema map[string]string `yaml:"schema"`
}

type Swagger struct {
	Definitions map[string]SwaggerDefinition           `yaml:"definitions"`
	Paths       map[string]map[string]SwaggerOperation `yaml:"paths"`
}

func capitalize(input string) string {
	return strings.ToUpper(input[0:1]) + input[1:]
}

func main() {

	// TODO: Make this configurable
	bytes, err := ioutil.ReadFile("test.yml")
	if err != nil {
		panic(err)
	}

	var swagger Swagger
	if err := yaml.Unmarshal(bytes, &swagger); err != nil {
		panic(err)
	}

	fmt.Printf("Swagger: %+v\n", swagger)

	if err := buildTypes(swagger.Definitions); err != nil {
		panic(err)
	}
	if err := buildRouter(swagger.Paths); err != nil {
		panic(err)
	}
	// TODO: Is this really the way I want to do this???
	if err := buildContexts(swagger.Paths); err != nil {
		panic(err)
	}
	if err := buildHandlers(swagger.Paths); err != nil {
		panic(err)
	}
}

type Generator struct {
	buf bytes.Buffer
}

func (g *Generator) Printf(format string, args ...interface{}) {
	fmt.Fprintf(&g.buf, format, args...)
}

// TODO: Add a nice comment!
// TODO: Make this write out to a file...
func buildTypes(definitions map[string]SwaggerDefinition) error {

	// TODO: Verify that the types are correct. In particular make sure they have the right references...

	var g Generator
	for name, definition := range definitions {
		definition.Printf(&g, name)
	}

	fmt.Printf(g.buf.String())

	return ioutil.WriteFile("generated/types.go", g.buf.Bytes(), 0644)
}

func buildRouter(paths map[string]map[string]SwaggerOperation) error {
	var g Generator

	// TODO: Add something to all these about being auto-generated

	g.Printf("package main\n\n")
	g.Printf("import \"github.com/gorilla/mux\"\n\n")
	g.Printf("func withRoutes(r *mux.Router) *mux.Router {\n")

	for path, pathObj := range paths {
		for method, op := range pathObj {
			// TODO: Validate the method
			// TODO: Note the coupling for the handler name here and in the handler function. Does that mean these should be
			// together? Probably...
			g.Printf("\tr.Methods(\"%s\").Path(\"%s\").HandlerFunc(%sHandler)\n",
				strings.ToUpper(method), path, capitalize(op.OperationID))
		}
	}
	// TODO: It's a bit weird that this returns a pointer that it modifies...
	g.Printf("\treturn r\n")
	g.Printf("}\n")

	fmt.Printf(g.buf.String())
	return ioutil.WriteFile("generated/router.go", g.buf.Bytes(), 0644)
}

func buildContexts(paths map[string]map[string]SwaggerOperation) error {
	// This includes the interfaces...
	var g Generator

	g.Printf("package main\n\n")

	g.Printf("import (\n")
	g.Printf("\t\"net/http\"\n")
	g.Printf(")\n\n")

	for _, path := range paths {
		for _, op := range path {
			// TODO: Should this be more functional??

			g.Printf("type %sInput struct {\n", capitalize(op.OperationID))
			// TODO: Add in input params...
			g.Printf("}\n")

			g.Printf("func New%sInput(r *http.Request) (*%sInput, error) {\n", capitalize(op.OperationID), capitalize(op.OperationID))
			// TODO: This should take in the http.Request object probably so that it can get the request parameters from that
			g.Printf("\treturn &%sInput{}, nil\n", capitalize(op.OperationID))
			g.Printf("}\n")

			g.Printf("func (i %sInput) Validate() error{\n", capitalize(op.OperationID))
			// TODO: Add in any validation...
			g.Printf("\treturn nil\n")
			g.Printf("}\n")
		}
	}

	// TODO: How should I name these things??? Should they be on a per-tag basis???
	g.Printf("\ntype Controller interface {\n")

	for _, path := range paths {
		for _, op := range path {
			g.Printf("\t%s(input %sInput) error\n", capitalize(op.OperationID), capitalize(op.OperationID))
		}
	}
	g.Printf("}\n")

	fmt.Printf(g.buf.String())
	return ioutil.WriteFile("generated/contexts.go", g.buf.Bytes(), 0644)
}

func buildHandlers(paths map[string]map[string]SwaggerOperation) error {

	var g Generator

	g.Printf("package main\n\n")

	g.Printf("import (\n")
	g.Printf("\t\"net/http\"\n")
	g.Printf(")\n\n")

	for _, path := range paths {
		for _, op := range path {
			g.Printf("func %sHandler(w http.ResponseWriter, r *http.Request) {\n", capitalize(op.OperationID))
			g.Printf("\tinput, err := New%sInput(r)\n", capitalize(op.OperationID))
			g.Printf("\tif err != nil {\n")
			// TODO: Handle these errors better... by returning something... the default presumably?
			g.Printf("\t\treturn\n")
			g.Printf("\t}\n")
			g.Printf("\terr = input.Validate()\n")
			g.Printf("\tif err != nil {\n")
			g.Printf("\t\treturn\n")
			g.Printf("\t}\n")

			// TODO: Call the actual interface... Not sure quite how I want to instantiate this...
			// TODO: Do something to handle the actual response

			g.Printf("}\n")
		}
	}

	fmt.Printf(g.buf.String())

	return ioutil.WriteFile("generated/handlers.go", g.buf.Bytes(), 0644)
}
