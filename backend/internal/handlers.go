package internal

import (
	"encoding/json"
	"log"
	"net/http"
)

type TranslationRequest struct {
	Translation string  `json:"translation"`
	Code        string  `json:"code"`
	Prompt      *string `json:"prompt,omitempty"`
}

type TranslationResponse struct {
	TranslatedCode string `json:"translatedCode"`
	ModelUsed      string `json:"modelUsed"`
}

/*
 * Define the actual routes here
 */

func AddRoutes(mux *http.ServeMux) {
	mux.Handle("/hello", HandleHelloHandler("hello"))
	mux.Handle("/health", HandleHealthHandler())
	mux.Handle("/api/translate", TranslateHandler())
}

/*
 * Define any route handlers here
 */

func TranslateHandler() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				return
			}

			// decode
			var req TranslationRequest
			err := json.NewDecoder(r.Body).Decode(&req)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			log.Printf("Translation request received: %s\n", req.Translation)

			// Translate the Code
			_, err = TranslateCode(r.Context(), req)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				log.Printf("ERROR DURING TRANSLATION")
				return
			}

			// encode response, return TODO
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
