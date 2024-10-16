package anthropic

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

type requestBody struct {
	Model         string           `json:"model"`
	MaxToken      int              `json:"max_tokens"`
	Messages      []models.Message `json:"messages"`
	SystemMessage string           `json:"system"`
	Stream        bool             `json:"stream"`
}

type Client struct {
	model         string
	systemMessage string
	apiKey        string
}

type EventData struct {
	Type    string `json:"type"`
	Message struct {
		Usage struct {
			InputTokens  int `json:"input_tokens"`
			OutputTokens int `json:"output_tokens"`
		} `json:"usage"`
	} `json:"message,omitempty"`
	Delta struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"delta,omitempty"`
	Index      int    `json:"index,omitempty"`
	StopReason string `json:"stop_reason,omitempty"`
	Usage      struct {
		OutputTokens int `json:"output_tokens"`
	} `json:"usage,omitempty"`
	Error struct {
		Type    string `json:"type"`
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

func NewClient(model string, systemMessage string, apiKey string) *Client {
	return &Client{model: model, systemMessage: systemMessage, apiKey: apiKey}
}

func (c *Client) SendMessage(conversation []models.Message) (models.Response, error) {
	url := "https://api.anthropic.com/v1/messages"
	body := requestBody{
		Model:         c.model,
		Messages:      conversation,
		MaxToken:      8192,
		SystemMessage: c.systemMessage,
		Stream:        true,
	}
	requestBody, _ := json.Marshal(body)
	// Create a new request
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	req.Header.Set("x-api-key", c.apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")
	req.Header.Set("content-type", "application/json")

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
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		// Ignore empty or ping lines
		if len(strings.TrimSpace(line)) == 0 || strings.HasPrefix(line, "event: ping") {
			continue
		}

		// Parse event data
		if strings.HasPrefix(line, "data: ") {
			jsonData := strings.TrimPrefix(line, "data: ")
			var eventData EventData
			if err := json.Unmarshal([]byte(jsonData), &eventData); err != nil {
				fmt.Println("Error parsing JSON:", err)
				continue
			}

			// Handle different event types
			switch eventData.Type {
			case "message_start":
				// get input token count
				// fmt.Println("Input tokens:", string(eventData.Message.Usage.InputTokens))
				// inputTokens := eventData.Message.Usage.InputTokens
				// if inputTokens > 150000 {
				// 	fmt.Println("Token count > 150000")
				// }
				continue
			case "content_block_delta":
				words := eventData.Delta.Text
				fmt.Print(words)
				responses = append(responses, words)
			case "content_block_start", "ping", "content_block_stop":
				continue
			case "message_delta":
				// get output token count
				// outputTokens := eventData.Usage.OutputTokens
				continue
			case "message_stop":
				// fmt.Println("Message Stop")
				break
			case "error":
				fmt.Println(eventData.Error.Message)
			default:
				fmt.Println("Unhandled event type:", eventData.Type)
			}
		}
	}
	response := strings.Join(responses, "")
	return models.Response{Content: response}, nil
}
