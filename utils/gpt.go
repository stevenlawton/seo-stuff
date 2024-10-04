package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

const openAIAPIURL = "https://api.openai.com/v1/chat/completions"

// GPT4Request represents the request payload for the OpenAI API
type GPT4Request struct {
	Model       string       `json:"model"`
	Messages    []GPTMessage `json:"messages"`
	MaxTokens   int          `json:"max_tokens,omitempty"`
	Temperature float64      `json:"temperature,omitempty"`
}

// GPTMessage represents a message in the conversation
type GPTMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// GPT4Response represents the response from the OpenAI API
type GPT4Response struct {
	Choices []struct {
		Message GPTMessage `json:"message"`
	} `json:"choices"`
}

// CallGPT4 makes a request to the GPT-4 API
func CallGPT4(prompt string) (string, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("OPENAI_API_KEY environment variable is not set")
	}

	requestBody := GPT4Request{
		Model: "gpt-4",
		Messages: []GPTMessage{
			{Role: "system", Content: "You are an assistant that provides brand-aligned suggestions for web content improvements."},
			{Role: "user", Content: prompt},
		},
		MaxTokens:   150,
		Temperature: 0.7,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", openAIAPIURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		return "", fmt.Errorf("OpenAI API error: %s", string(bodyBytes))
	}

	var gptResponse GPT4Response
	if err := json.NewDecoder(resp.Body).Decode(&gptResponse); err != nil {
		return "", err
	}

	if len(gptResponse.Choices) == 0 {
		return "", fmt.Errorf("No response from GPT-4")
	}

	return gptResponse.Choices[0].Message.Content, nil
}
