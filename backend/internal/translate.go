package internal

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

/*
 * TranslateCode() will call the corresponding model that is best at the
 * given code translation pair.
 */
func TranslateCode(ctx context.Context, req TranslationRequest) (*TranslationResponse, error) {
	if req.Model == "" {
		log.Println("[TranslateCode]: Error. No model provided")
		return nil, errors.ErrUnsupported
	}

	model := req.Model

	switch model {
	case "gpt4o":
		log.Printf(
			"[TranslateCode]: Translation supported -- using %s for %s\n",
			model,
			req.Translation,
		)
		prompt, err := BuildPrompt(ctx, req.Translation, req.Code)
		if err != nil {
			return nil, err
		}

		resp, err := OpenAIRequest(ctx, prompt)
		if err != nil {
			return nil, err
		}

		return resp, nil

	case "llama-3.2-3b":
		log.Printf(
			"[TranslateCode]: Translation supported -- using %s for %s\n",
			model,
			req.Translation,
		)
		prompt, err := BuildPrompt(ctx, req.Translation, req.Code)
		if err != nil {
			return nil, err
		}

		resp, err := LlamaRequest(ctx, prompt)
		if err != nil {
			return nil, err
		}

		return resp, nil

	case "deepseek-6.7b":
		log.Printf(
			"[TranslateCode]: Translation supported -- using %s for %s\n",
			model,
			req.Translation,
		)
		prompt, err := BuildPrompt(ctx, req.Translation, req.Code)
		if err != nil {
			return nil, err
		}

		resp, err := DeepSeekRequest(ctx, prompt)
		if err != nil {
			return nil, err
		}

		return resp, nil

	case "phi-2.7b":
		log.Printf(
			"[TranslateCode]: Translation supported -- using %s for %s\n",
			model,
			req.Translation,
		)
		prompt, err := BuildPhiPrompt(ctx, req.Translation, req.Code)
		if err != nil {
			return nil, err
		}

		resp, err := PhiRequest(ctx, prompt, req.Code)
		if err != nil {
			return nil, err
		}

		return resp, nil

	default:
		log.Printf("[TranslateCode]: Given model is not supported")
	}

	return nil, nil
}

/*
 * BuildPrompt() will build the prompt to send to the LLM
 * Used for GPT, Llama, and Deepseek-Coder
 */
func BuildPrompt(ctx context.Context, translation string, code string) (string, error) {
	translationPair := strings.Split(translation, "-")
	if len(translationPair) != 2 {
		return "", errors.New(
			"[BuildPrompt]: Translation string is not the correct format of from-to",
		)
	}

	from := translationPair[0]
	to := translationPair[1]

	// NOTE: Adding the tokens <CODE START> <CODE END> is used to help parse the returned code
	prompt := fmt.Sprintf(
		"You are an AI that only responds with translated code wrapped between the following delimiters:\n<CODE START>\n...translated code here...\n<CODE END>\nOnly return the code block. Do NOT include any extra text, markdown, or commentary.\nPlease translate the following %s code to the same equivalent in %s:\n%s\n",
		from,
		to,
		code,
	)

	return prompt, nil
}

/*
 * BuildPhiPrompt(): Phi2 kind of sucks at following directions. Here, this is used to
 * build a simpler prompt that the model can hopefully comprehend most of the time.
 */
func BuildPhiPrompt(ctx context.Context, translation string, code string) (string, error) {
	translationPair := strings.Split(translation, "-")
	if len(translationPair) != 2 {
		return "", errors.New(
			"[BuildPrompt]: Translation string is not the correct format of from-to",
		)
	}

	from := translationPair[0]
	to := translationPair[1]

	prompt := fmt.Sprintf(
		"Please translate the following %s code to %s. Only return the translated code and no other text. Thank you\n CODE TO TRANSLATE:\n%s\n",
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
	inCodeBlock := false
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "<CODE START>") {
			inCodeBlock = true
			continue
		}

		if strings.Contains(line, "<CODE END>") {
			inCodeBlock = false
			break
		}

		if inCodeBlock {
			resp.TranslatedCode += line + "\n"
		}
	}

	log.Println("[ParseModelResponse]: Code has been successfully parsed and translated...")
	return &resp
}

/*
 * ParsePhiModelResponse() will look for the most often used delimiters / formatters it returns.
 * NOTE: This model literally never listened to the other prompt (using the <CODE> delimiters).
 * In this case, since the model is barely even used, I just look to parse out the most common
 * formatters it returns (usually markdown symbols)
 */
func ParsePhiModelResponse(modelResp string) *TranslationResponse {
	resp := TranslationResponse{
		TranslatedCode: "",
	}

	scanner := bufio.NewScanner(strings.NewReader(modelResp))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Compare(line, "```") == 0 || strings.Compare(line, "```python") == 0 {
			for scanner.Scan() {
				codeLine := scanner.Text()
				if strings.Compare(codeLine, "```") == 0 {
					break
				}

				resp.TranslatedCode += codeLine
				resp.TranslatedCode += "\n"
			}
			break
		}
	}

	log.Println("[ParseModelResponse]: Prased and Translated code:\n", resp.TranslatedCode)
	return &resp
}

/*
 * OpenAIRequest() will send the prompt to the OpenAI model API and
 * receive the chat response.
 */
func OpenAIRequest(ctx context.Context, prompt string) (*TranslationResponse, error) {
	key := os.Getenv("OPEN_AI_KEY")
	if key == "" {
		return nil, fmt.Errorf("[OpenAIRequest]: OPEN_AI_KEY environment variable not set")
	}

	client := openai.NewClient(
		option.WithAPIKey(key),
	)

	log.Println("[OpenAIRequest]: building and sending chat...")

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
		log.Println("[OpenAIRequest]: Issue with completion: ", err)
		return nil, err
	}

	// Grab the response, parse, return
	modelResponse := chatCompletion.Choices[0].Message.Content
	resp := ParseModelResponse(modelResponse)
	resp.ModelUsed = "gpt4o"

	return resp, nil
}

/*
 * LlamaRequest() will send the given prompt to llama-3.2-3b-instruct API
 * and recieve back the response.
 * NOTE: This model is accessed through the Nvidia API.
 */
func LlamaRequest(ctx context.Context, prompt string) (*TranslationResponse, error) {
	// url := "https://integrate.api.nvidia.com/v1/chat/completions"
	url := os.Getenv("NVIDIA_LLAMA_URL")
	if url == "" {
		return nil, fmt.Errorf("[LlamaRequest]: NVIDIA_LLAMA_URL environment variable not set")
	}

	key := os.Getenv("NVIDIA_API_KEY")
	if key == "" {
		return nil, fmt.Errorf("[LlamaRequest]: NVIDIA_API_KEY environment variable not set")
	}

	requestData := map[string]interface{}{
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": prompt,
			},
		},
		"stream":            false,
		"model":             "meta/llama-3.2-3b-instruct",
		"max_tokens":        1024,
		"presence_penalty":  0,
		"frequency_penalty": 0,
		"top_p":             0.7,
		"temperature":       0.2,
	}

	jsonData, err := json.Marshal(requestData)
	if err != nil {
		return nil, fmt.Errorf("[LlamaRequest]: error marshaling JSON: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("[LlamaRequest]: error creating request: %v", err)
	}

	// Add headers
	req.Header.Set("Authorization", "Bearer "+key)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("[LlamaRequest]: error sending request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf(
			"[LlamaRequest]: API error (status %d): %s",
			resp.StatusCode,
			string(bodyBytes),
		)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("[LlamaRequest]: error reading response body: %v", err)
	}

	var responseData map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &responseData); err != nil {
		return nil, fmt.Errorf("[LlamaRequest]: error parsing response JSON: %v", err)
	}

	// Extract the content from the response
	// NOTE: lots of stuff to parse and extract here..
	choices, ok := responseData["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		return nil, fmt.Errorf("[LlamaRequest]: unexpected response format: missing choices")
	}

	choice, ok := choices[0].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("[LlamaRequest]: unexpected response format: invalid choice")
	}

	message, ok := choice["message"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("[LlamaRequest]: unexpected response format: missing message")
	}

	content, ok := message["content"].(string)
	if !ok {
		return nil, fmt.Errorf("[LlamaRequest]: unexpected response format: missing content")
	}

	response := ParseModelResponse(content)
	response.ModelUsed = "llama-3.2-3b"
	return response, nil
}

/*
 * DeepSeekRequest() will send the given prompt to the locally running Deepseek-Coder
 * model to get the code translation.
 * NOTE: This does require the deepseek-coder:6.7b model to be ran locally using Ollama.
 * NOTE: This model requires quite a bit of RAM to actually be ran locally. Do use with caution.
 */
func DeepSeekRequest(ctx context.Context, prompt string) (*TranslationResponse, error) {
	// Calling local Ollama server from a docker container
	// url := "http://host.docker.internal:11434/api/generate"
	url := os.Getenv("OLLAMA_URL")
	if url == "" {
		return nil, fmt.Errorf("[DeepSeekRequest]: OLLAMA_URL environment variable not set")
	}

	payload := map[string]interface{}{
		"model":  "deepseek-coder:6.7b",
		"prompt": prompt,
		"stream": false,
	}

	jsonData, _ := json.Marshal(payload)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("[DeepSeekRequest]: error getting response: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var responseMap map[string]interface{}
	err = json.Unmarshal(body, &responseMap)
	if err != nil {
		return nil, fmt.Errorf("[DeepSeekRequest]: error decoding JSON: %v", err)
	}

	modelResponse := responseMap["response"].(string)
	response := ParseModelResponse(modelResponse)
	response.ModelUsed = "deepseek-6.7b"

	return response, nil
}

/*
 * PhiRequest() will send the given prompt to the locally running Phi model to get the
 * code translation.
 * NOTE: This does require the Phi:2.7b model to be ran locally using Ollama.
 */
func PhiRequest(ctx context.Context, prompt string, code string) (*TranslationResponse, error) {
	// Calling local Ollama server from a docker container
	//	url := "http://host.docker.internal:11434/api/generate"

	url := os.Getenv("OLLAMA_URL")
	if url == "" {
		return nil, fmt.Errorf("[PhiRequest]: OLLAMA_URL environment variable not set")
	}
	payload := map[string]interface{}{
		"model":  "phi:2.7b",
		"prompt": prompt,
		"stream": false,
	}

	jsonData, _ := json.Marshal(payload)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("[PhiRequest]: error getting response: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var responseMap map[string]interface{}
	err = json.Unmarshal(body, &responseMap)
	if err != nil {
		return nil, fmt.Errorf("[PhiRequest]: error getting response: %v", err)
	}

	modelResponse := responseMap["response"].(string)
	response := ParsePhiModelResponse(modelResponse)
	if response.TranslatedCode == "" {
		response.TranslatedCode = modelResponse
	}
	response.ModelUsed = "phi-2.7b"

	return response, nil
}
