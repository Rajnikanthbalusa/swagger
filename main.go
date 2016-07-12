package main

import (
	"bytes"
	"fmt"
	"io/ioutil"

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

type SwaggerPath struct {
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
	Definitions map[string]SwaggerDefinition      `yaml:"definitions"`
	Paths       map[string]map[string]SwaggerPath `yaml:"paths"`
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

	buildTypes(swagger.Definitions)
	buildRouter()
	// TODO: Is this really the way I want to do this???
	for _, pathObj := range swagger.Paths {
		for _, responseObj := range pathObj {
			buildContexts(responseObj.Responses)
		}
	}
	buildHandlers()
}

type Generator struct {
	buf bytes.Buffer
}

func (g *Generator) Printf(format string, args ...interface{}) {
	fmt.Fprintf(&g.buf, format, args...)
}

// TODO: Add a nice comment!
// TODO: Make this write out to a file...
func buildTypes(definitions map[string]SwaggerDefinition) {

	// TODO: Verify that the types are correct. In particular make sure they have the right references...

	var g Generator
	for name, definition := range definitions {
		definition.Printf(&g, name)
	}

	fmt.Printf(g.buf.String())
}

func buildRouter() {
	var g Generator
	g.Printf("package main\n\n")
	g.Printf("import net/http\n\n")
	g.Printf("func withRoutes(r *mux.Router) {\n")

	// TODO: Note that this is coupled with the handler names...
	g.Printf("r.Methods(\"%s\").Path(\"%s\").HandlerFunc(%s))\n", "GET", "/book/{id}", "getBookHandler")

	g.Printf("}\n")

}

func buildContexts(responses map[string]SwaggerResponse) {
	// This includes the interfaces...
	var g Generator

	// TODO: Create the Input type
	// TODO: If params validate them (either from query, path or body) - this should probably be from the input type???

	// TODO: How should I name these things???
	g.Printf("type BooksIDController interface {\n")

	// TODO: Can we return something better than an error?
	g.Printf("\tGetBook(input BookIDInput) error\n")
	g.Printf("}\n")

	fmt.Printf(g.buf.String())
}

func buildHandlers() {

	var g Generator
	// TODO: Don't hard-code getBookHandler
	g.Printf("func %s {\n", "getBookHandler")

	// Build a context object
	// Check for an error

	// For each parameter
	// Validate (possibly recursively...)... maybe this should be built into a context handler...

	// Build the context and

	g.Printf("}")
}
