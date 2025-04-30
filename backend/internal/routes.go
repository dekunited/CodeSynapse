package internal

import (
	"net/http"
)

func AddRoutes(
	mux *http.ServeMux,
	// add dependencies here
) {
	// add routes here
	mux.Handle("/hello", HandleHelloHandler("hello"))
	mux.Handle("/health", HandleHealthHandler())
}
