package internal

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

/*
 * Map for storing best models for each translation
 */
var bestTranslationMap = map[string]string{
	"python-java": "Transcoder",
	"go-python":   "ChatGPT",
	"c++-python":  "Transcoder",
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
	case "ChatGPT":
		log.Printf("    Translation supported -- using %s for %s\n", model, req.Translation)
		prompt, err := BuildPrompt(ctx, req.Translation, req.Code)
		if err != nil {
			return nil, err
		}

		resp, err := OpenAIRequest(ctx, prompt)
		if err != nil {
			return nil, err
		}

		return resp, nil

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

  // TODO: Tune this, specify the exact format we want the response in to simplify parsing it easily.
	prompt := fmt.Sprintf(
		"Please translate the following %s code to the same equivalent in %s\nCODE: %s\n",
		from,
		to,
		code,
	)

	return prompt, nil
}

// TODO: Parse the response and get just the code
func ParseModelResponse() { } 

// Function to send request to OpenAI API
func OpenAIRequest(ctx context.Context, prompt string) (*TranslationResponse, error) {
	key := os.Getenv("OPEN_AI_KEY")

	client := openai.NewClient(
		option.WithAPIKey(key),
	)

	log.Println("OpenAI Request: building and sending chat")

	chatCompletion, err := client.Chat.Completions.New(
		ctx,
		openai.ChatCompletionNewParams{
			Messages: []openai.ChatCompletionMessageParamUnion{
				openai.UserMessage(prompt),
			},
			Model: openai.ChatModelGPT4o,
		},
	)
	if err != nil {
		log.Println("   Open AI (ERROR) - Issue with completion: ", err)
		return nil, err
	}

	log.Println("     Open AI: ", chatCompletion.Choices[0].Message.Content)
  // TODO: Parse response

	return nil, nil
}

// Function to send request to GCP service

// Function to send request to Hugging Face? / Whatever else
