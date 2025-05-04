package internal

import (
	"bufio"
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
	"go-python": "ChatGPT",
}

/*
 * TranslateCode() will call the corresponding model that is best at the 
 * given code translation pair.
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

	// Add the prompt here. May need more tuning.
	// NOTE: Adding the tokens <CODE START> <CODE END> is used to help parse the returned code
	prompt := fmt.Sprintf(
		"Please translate the following %s code to the same equivalent in %s\n In your response, please add the <CODE START> delimiter before showing the actual resulting code and then <CODE END> after. Additionally, just include the code in plaintext rather than including any markdown. Thank you.\nCODE: %s\n",
		from,
		to,
		code,
	)

	return prompt, nil
}

/*
 * ParseModelResponse() will look for the delimiters/tokens added to the response.
 * The goal is to parse out only the returned code from the model's full response.
 * NOTE: Anything in between <CODE START> and <CODE END> is code to be returned.
 */
func ParseModelResponse(modelResp string) *TranslationResponse {
	resp := TranslationResponse{
		TranslatedCode: "",
	}

	scanner := bufio.NewScanner(strings.NewReader(modelResp))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Compare(line, "<CODE START>") == 0 {
			for scanner.Scan() {
				codeLine := scanner.Text()
				if strings.Compare(codeLine, "<CODE END>") == 0 {
					break
				}

				resp.TranslatedCode += codeLine
				resp.TranslatedCode += "\n"
			}
			break
		}
	}

	log.Println("CODE SNIPPET TRANSLATED: \n", resp.TranslatedCode)

	return &resp
}

/*
 * OpenAIRequest() will send the prompt to the OpenAI model and
 * receive the chat response.
 */
func OpenAIRequest(ctx context.Context, prompt string) (*TranslationResponse, error) {
	key := os.Getenv("OPEN_AI_KEY")

	client := openai.NewClient(
		option.WithAPIKey(key),
	)

	log.Println("OpenAI Request: building and sending chat")

	// Sending prompt to the specified model
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

	// Grab the response, parse, return
	// NOTE: update the ModelUsed field once actual model is decided
	modelResponse := chatCompletion.Choices[0].Message.Content
	resp := ParseModelResponse(modelResponse)
	resp.ModelUsed = "OpenAI Model Placeholder"

	return resp, nil
}

// Function to send request to GCP service

// Function to send request to Hugging Face? / Whatever else
