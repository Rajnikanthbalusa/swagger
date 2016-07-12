package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// TODO: Some initialization of config???

func main() {
	router := mux.NewRouter()
	router = withMiddleware(withRoutes(router))
	server := &http.Server{
		// TODO: This should be configurable???
		Addr:    fmt.Sprintf(":80"),
		Handler: router,
	}

	log.Fatal(server.ListenAndServe())
}
