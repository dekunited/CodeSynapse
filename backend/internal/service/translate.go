package service

// Store the translation types here
var bestTranslationMap = map[string]string {
}

type TranslationRequest struct {
	SourceCode     string `json:"sourceCode"`
	SourceLanguage string `json:"sourceLanguage"`
	TargetLanguage string `json:"targetLanguage"`
}

type TranslationResponse struct {
	TranslatedCode string    `json:"translatedCode"`
	ModelUsed      string    `json:"modelUsed"`
}

// TODO
func PromptBuilder() {
}

// TODO
func TranslateCode() {
}
