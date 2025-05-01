package internal

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/hupe1980/go-huggingface"
)

/*
 * Map for storing best models for each translation
 */
var bestTranslationMap = map[string]string{
	"python-java": "CodeBERT",
	"go-python":   "CodeBERT",
}

func TranslateCode(
	ctx context.Context,
	translation string,
	code string,
) (*TranslationResponse, error) {
	model, ok := bestTranslationMap[translation]
	if !ok {
		// TODO: Error here, not supported, etc.
	}

	switch model {
	case "CodeBERT":
		log.Printf("    Translation supported -- using %s for %s\n", model, translation)
		TranslateWithCodeBERT(ctx, translation, code)
	case "Codex":
		log.Printf("    Translation supported -- using %s for %s\n", model, translation)
	case "Transcoder":
		log.Printf("    Translation supported -- using %s for %s\n", model, translation)
	default:
		log.Printf("    Translation not supported")
	}

	return nil, nil
}

func TranslateWithCodeBERT(ctx context.Context, translation string, code string) {
	token := os.Getenv("HUGGING_FACE_TOKEN")
	if token == "" {
		// todo - error here
	}

	client := huggingface.NewInferenceClient(token)
	log.Println("Translating... ", client)

	prompt, err := BuildPrompt(ctx, translation, code)
	if err != nil {
		// todo - erro here
	}
	log.Println("     Prompt... ", prompt)

	maxNewTokens := 1024
	temp := 0.2
	topP := 0.95
	returnFullText := false

	result, err := client.TextGeneration(ctx, &huggingface.TextGenerationRequest{
		Model:  "microsoft/codebert-base", 
		Inputs: prompt,
		Parameters: huggingface.TextGenerationParameters{
			MaxNewTokens:   &maxNewTokens,   // Adjust based on expected response size
			Temperature:    &temp,           // Lower for more deterministic results
			TopP:           &topP,           // Nucleus sampling
			ReturnFullText: &returnFullText, // Only return the generated text, not the prompt
		},
	})

  if err != nil {
    log.Println("BERT -- ", err)
  }

  log.Println("IDKIDKIDIKD: ", result)
}

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
