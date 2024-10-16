package openai

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"jervis/models"
	"net/http"
	"os"
	"strings"
)

type MessageOpenAI struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type RequestBodyOpenAI struct {
	Model    string           `json:"model"`
	Messages []models.Message `json:"messages"`
	Stream   bool             `json:"stream"`
}

type ChatCompletionChunk struct {
	ID                string `json:"id"`
	Object            string `json:"object"`
	Created           int64  `json:"created"`
	Model             string `json:"model"`
	SystemFingerprint string `json:"system_fingerprint"`
	Choices           []struct {
		Index        int            `json:"index"`
		Delta        models.Message `json:"delta"`
		FinishReason string         `json:"finish_reason"`
	} `json:"choices"`
}

type Client struct {
	model         string
	systemMessage string
	apiKey        string
}

func NewClient(model, systemMessage, apiKey string) *Client {
	return &Client{model: model, systemMessage: systemMessage, apiKey: apiKey}
}

func (c *Client) SendMessage(conversation []models.Message) (models.Response, error) {
	url := "https://api.openai.com/v1/chat/completions"
	if len(conversation) == 1 {
		conversation = append([]models.Message{{
			Role:    "system",
			Content: c.systemMessage,
		}}, conversation...)
	}

	body := RequestBodyOpenAI{
		Model:    c.model,
		Messages: conversation,
		Stream:   true,
	}
	requestBody, err := json.Marshal(body)
	if err != nil {
		return models.Response{}, fmt.Errorf("failed to marshal request body: %w", err)
	}
	// Create a new request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return models.Response{}, fmt.Errorf("failed to create HTTP request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	fmt.Println(req)
	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error on response.\n[ERROR] -", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	var responses []string
	// Read the streamed response
	reader := bufio.NewReader(resp.Body)
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("Error on reading response.\n[ERROR] -", err)
			os.Exit(1)
		}
		if len(line) <= 1 {
			continue
		}

		jsonData := bytes.TrimPrefix(line, []byte("data: "))

		// Unmarshal JSON into the struct
		var chunk ChatCompletionChunk
		err2 := json.Unmarshal(jsonData, &chunk)
		if err2 != nil {
			return models.Response{}, fmt.Errorf("Error unmarshalling JSON.\n[ERROR] - %w, %s", err2, string(jsonData))
		}

		if chunk.Choices[0].FinishReason == "stop" {
			break
		}

		words := chunk.Choices[0].Delta.Content
		fmt.Print(words)
		responses = append(responses, words)
	}
	response := strings.Join(responses, "")
	fmt.Println(response)
	return models.Response{Content: response}, nil
}
