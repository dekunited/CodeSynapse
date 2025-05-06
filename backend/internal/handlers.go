package internal

import (
	"encoding/json"
	"log"
	"net/http"
  "strings"
)

type TranslationRequest struct {
	Translation string  `json:"translation"`
	Code        string  `json:"code"`
	Model       *string `json:"model,omitempty"`
}

type TranslationResponse struct {
	TranslatedCode string `json:"translatedCode"`
	ModelUsed      string `json:"modelUsed"`
}

/*
 * Define routes here
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
			translationResp, err := TranslateCode(r.Context(), req)
			if err != nil {
        if strings.Contains(err.Error(), "connection refused") {
          http.Error(w, err.Error(), http.StatusServiceUnavailable)
          log.Printf("[TranslateHandler]: Model or service is unavailable")
          return
        }
				http.Error(w, err.Error(), http.StatusInternalServerError)
				log.Printf("[TranslateHandler]: Error translating code:\n%v", err)
				return
			}

      // Return
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			err = json.NewEncoder(w).Encode(translationResp)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				log.Printf("[TranslateHandler]: Error encoding response:\n%v", err)
				return
			}
		},
	)
}
