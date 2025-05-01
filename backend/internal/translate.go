package internal

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
)

/*
 * Map for storing best models for each translation
 */
var bestTranslationMap = map[string]string{
	"python-java": "CodeBERT",
	"go-python":   "CodeBERT",
}

/*
 * TranslateCode() 
 */
func TranslateCode(ctx context.Context, req TranslationRequest) (*TranslationResponse, error) {
	model, ok := bestTranslationMap[req.Translation]
	if !ok {
		// TODO: Error here, not supported, etc.
	}

	switch model {
	case "CodeBERT":
		log.Printf("    Translation supported -- using %s for %s\n", model, req.Translation)
		//TranslateWithCodeBERT(ctx, req)
	case "Codex":
		log.Printf("    Translation supported -- using %s for %s\n", model, req.Translation)
	case "Transcoder":
		log.Printf("    Translation supported -- using %s for %s\n", model, req.Translation)
	default:
		log.Printf("    Translation not supported")
	}

	return nil, nil
}

/*
 * Build the prompt to send to the LLM
 */
func BuildPrompt(ctx context.Context, translation string, code string) (string, error) {
	translationPair := strings.Split(translation, "-")
	if len(translationPair) != 2 {
		return "", errors.New(
			"BuildPrompt(): Translation string is not the correct format of from-to",
		)
	}

	from := translationPair[0]
	to := translationPair[1]

	prompt := fmt.Sprintf(
		"Please translate the following %s code to the same equivalent in %s\nCODE: %s\n",
		from,
		to,
		code,
	)

	return prompt, nil
}
