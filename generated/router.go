package main

func New(version string) http.Handler {
	router := mux.NewRouter()

	return withMiddleware(router)
}
