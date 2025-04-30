package internal

import (
	"log"
	"net/http"
)

/*
  Define any http route handlers here
*/

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
