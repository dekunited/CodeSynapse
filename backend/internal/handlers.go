package internal

import (
	"log"
	"net/http"
)

/*
  Define the actual routes here
*/

func AddRoutes(mux *http.ServeMux) {
	mux.Handle("/hello", HandleHelloHandler("hello"))
	mux.Handle("/health", HandleHealthHandler())
	mux.Handle("/api/translate", TranslateHandler())
}

/*
  Define any route handlers here
*/

func TranslateHandler() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
		},
	)
}

func HandleHelloHandler(hello string) http.Handler {
	hello2 := hello + "2"
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"message": "` + hello2 + `"}`))
			log.Printf("HandleHelloHandler(): Request recieved")
		},
	)
}

func HandleHealthHandler() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"status": "ok"}`))
			log.Printf("HandleHealthHandler(): Request recieved")
		},
	)
}
