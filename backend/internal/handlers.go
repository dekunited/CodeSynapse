package internal

import (
	"encoding/json"
	"log"
	"net/http"
)

type TranslationRequest struct {
	Translation string  `json:"translation"`
	Code        string  `json:"code"`
	Prompt      *string `json:"prompt,omitempty"` // might not need depending on setup for other models
}

type TranslationResponse struct {
	TranslatedCode string `json:"translatedCode"`
	ModelUsed      string `json:"modelUsed"`
}

/*
 * Define the actual routes here
 */

func AddRoutes(mux *http.ServeMux) {
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
        log.Printf("[TranslateHandler]: Error decoding request: %v", err)
				return
			}

      log.Printf("[TranslateHandler]: Translation request received")

			// Translate the Code
			_, err = TranslateCode(r.Context(), req)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
        log.Printf("[TranslateHandler]: Error translating code:\n%v", err)
				return
			}

			// encode response, return TODO
		},
	)
}

