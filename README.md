# swagger
Swagger test

# Running
```
go run main.go
cd generated
go run main.go middleware.go router.go handlers.go contexts.go controller.go
```

## After Generating

### Main
generated/main.go is a hard-coded file that implement the initial main functionality.

### Files to Change
- In controller.go implement the logic of your handlers.
- In middleware.go add any middleware specific to your service


## Swagger Spec

Currently, this repo doesn't implement the entire Swagger Spec. This is a non-exhaustive list of the parts of the Swagger specification that haven't been implemented yet.

### Planning to Implement
All Swagger Data Types (long, float, double, etc...)
Schema
  - parameters
  - tags
$ref in multiple places
All HTTP operations (this may already just work)
Operation
  - tags
Most Parameter Types
  - Path
  - Query
  - Header
  - Body
Required Fields
Response
  - Headers

### Not Planning to Implement (at least for now)
All Mime Types
Patterned Fields (these are vendor specific extensions anyway)
Multi-File Swagger Definitions
Everything in JSON-Schema
Schema
  - host
  - basePath
  - scheme (just http / https)
  - consumes
  - produces
  - securityDefinitions
  - security
Consumes
  - produces
  - consumes
  - schemes
  - security
Form parameter type (though if it's easy maybe we should just add...)
Parameter
  - items
  - collectionFormat
  - all the json schema requirements? (uniqueItems, multipleOf, etc...)
Schema object (for now going to try to get this from somewhere else...)
Discriminators
XML Modeling
Security Objects

