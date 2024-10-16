package internal

import (
	"bufio"
	// "bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	// "path/filepath"
	// "strings"

	"github.com/chzyer/readline"
)

type Vendor string
type Model string

const (
	Anthropic Vendor = "anthropic"
	OpenAI    Vendor = "OpenAI"

	gpt4o    Model = "gpt-4o"
	sonnet35 Model = "claude-3-5-sonnet-20240620"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatCompletionChunk struct {
	ID                string `json:"id"`
	Object            string `json:"object"`
	Created           int64  `json:"created"`
	Model             string `json:"model"`
	SystemFingerprint string `json:"system_fingerprint"`
	Choices           []struct {
		Index        int     `json:"index"`
		Delta        Message `json:"delta"`
		FinishReason string  `json:"finish_reason"`
	} `json:"choices"`
}

type RequestBodyOpenAI struct {
	Model    Model     `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
}

type RequestBodyAnthropic struct {
	Model         Model     `json:"model"`
	Messages      []Message `json:"messages"`
	Max_token     int       `json:"max_token,omitempty"`
	Stop_sequence []string  `json:"stop_sequence,omitempty"`
	Stream        bool      `json:"stream"`
	System        string    `json:"system"`
	Temperature   float32   `json:"temperature"`
}

type responseAnthropic struct {
	Id      string `json:"id"`
	Type    string `json:"type"`
	Role    string `json:"role"`
	Content struct {
		Type string `json:"type"`
		Text string `json:"text"`
	}
	Model string `json:"model"`
}

type RequestBundle struct {
	vendor        Vendor
	model         Model
	systemMessage string
	prompt        string
	body          []byte
}

func isValidVendor(vendor Vendor) bool {
	return vendor == Anthropic || vendor == OpenAI
}

func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func isValidModel(model string) bool {
	validModel := false
	models := []string{string(sonnet35), "gpt-4o"}
	for _, v := range models {
		if v == model {
			validModel = true
			break
		}
	}
	if !validModel {
		return false
	}
	return true
}

func getFileContent(filePath string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}
	return string(data), nil
}

func decodeJSON(filePath string, body interface{}) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("Error opening JSON file %s", err)
	}
	defer file.Close()

	return json.NewDecoder(file).Decode(body)
}

func encodeJSON(filePath string, body interface{}) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("Error creating JSON file %s", err)
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(body)
}

func HandleInput(prompt bool) []string {
	qlines := make([]string, 0)

	if prompt {
		qlines = readFromPrompt()
	} else {
		qlines = readFromStdin()
	}

	return qlines
}

func readFromStdin() []string {
	scanner := bufio.NewScanner(os.Stdin)
	lines := make([]string, 0)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Error reading from standard input:", err)
	}

	return lines
}

func readFromPrompt() []string {
	lines := make([]string, 0)
	rl, err := readline.New("> ")
	if err != nil {
		panic(err)
	}
	defer rl.Close()

	fmt.Println("Enter/paste your text. Press Ctrl+D on a new line to finish.")

	for {
		line, err := rl.Readline()
		if err != nil { // io.EOF will be returned if Ctrl+D is pressed
			if err == readline.ErrInterrupt {
				// Handle interrupt (Ctrl+C)
				continue
			} else if err == io.EOF {
				break
			}
			panic(err)
		}

		lines = append(lines, line)
	}
	return lines
}

func FormatOut(filePath string) {
	var body RequestBodyOpenAI
	err := decodeJSON(filePath, &body)
	if err != nil {
		fmt.Println("Error reading file.\n[ERROR] -", err)
		return
	}
	for _, v := range body.Messages {
		if v.Role == "user" || v.Role == "assistant" {
			fmt.Println("*" + v.Role + "*" + ":")
			fmt.Println("\t" + v.Content)
			fmt.Println()
		}
	}
}

// func getVendorFromFilePath(filePath string) (vendor string) {
// 	base := filepath.Base(filePath)
// 	parts := strings.Split(base, "_")
// 	vendor := parts[len(parts)-1]
// 	return
// }

func getRequest(filePath string, newSession bool, bundle RequestBundle) (requestBody []byte) {
	if newSession || !fileExists(filePath) {
		requestBody = createNewRequestBody(bundle)
		// Open the file for writing
		if err := encodeJSON(filePath, requestBody); err != nil {
			fmt.Println("Error writing new session file:", err)
		}
	} else {
		if bundle.vendor == Anthropic {

		} else if bundle.vendor == OpenAI {

		}
		var request RequestBodyOpenAI
		err := decodeJSON(filePath, &request)
		if err != nil {
			fmt.Println("Error decoding JSON.\n[ERROR] -", err)
			os.Exit(1)
		}

		// Append data to slice
		msg := Message{
			Role:    "user",
			Content: bundle.prompt,
		}
		request.Messages = append(request.Messages, msg)

		// cfg.model = request.Model

		requestBody, _ = json.Marshal(request)

		err2 := encodeJSON(filePath, request)
		if err2 != nil {
			fmt.Println("Error encoding JSON.\n[ERROR] -", err2)
			os.Exit(1)
		}
	}
	return
}

func createNewRequestBody(bundle RequestBundle) (requestBody []byte) {
	// body := map[string]interface{}{}
	if bundle.vendor == Anthropic {
		body := RequestBodyAnthropic{
			Model:  bundle.model,
			System: bundle.systemMessage,
			Messages: []Message{{
				Role:    "user",
				Content: bundle.prompt,
			}},
			Stream: true,
		}
		requestBody, _ = json.Marshal(body)

	} else if bundle.vendor == OpenAI {
		body := RequestBodyOpenAI{
			Model: bundle.model,
			Messages: []Message{
				{
					Role:    "system",
					Content: bundle.systemMessage,
				},
				{
					Role:    "user",
					Content: bundle.prompt,
				},
			},
			Stream: true,
		}
		requestBody, _ = json.Marshal(body)
	}
	return
}

func setHeader(req *http.Request, key string, vendor Vendor) {
	if vendor == Anthropic {
		req.Header.Set("x-api-key", key)

	} else if vendor == OpenAI {
		req.Header.Set("Authorization", "Bearer "+key)
	}
	req.Header.Set("Content-Type", "application/json")
}

// func performRequest(url string, requestBody []byte, key string) string {
//
// 	// Create a new request
// 	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
// 	setHeader(req, key, vendor)
//
// 	// Send the request
// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		fmt.Println("Error on response.\n[ERROR] -", err)
// 		os.Exit(1)
// 	}
// 	defer resp.Body.Close()
//
// 	var responses []string
// 	// Read the streamed response
// 	reader := bufio.NewReader(resp.Body)
// 	for {
// 		line, err := reader.ReadBytes('\n')
// 		if err != nil {
// 			if err == io.EOF {
// 				break
// 			}
// 			fmt.Println("Error on reading response.\n[ERROR] -", err)
// 			os.Exit(1)
// 		}
// 		if len(line) <= 1 {
// 			continue
// 		}
//
// 		jsonData := bytes.TrimPrefix(line, []byte("data: "))
//
// 		// Unmarshal JSON into the struct
// 		var chunk ChatCompletionChunk
// 		err2 := json.Unmarshal(jsonData, &chunk)
// 		if err2 != nil {
// 			fmt.Println("Error unmarshalling JSON.\n[ERROR] -", err2, string(jsonData))
// 			os.Exit(1)
// 		}
//
// 		if chunk.Choices[0].FinishReason == "stop" {
// 			break
// 		}
//
// 		words := chunk.Choices[0].Delta.Content
// 		fmt.Print(words)
// 		responses = append(responses, words)
// 	}
// 	response := strings.Join(responses, "")
// 	return response
// }

func persistConversation(filePath string, role string, response string) {
	if !fileExists(filePath) {
		return
	}
	var request RequestBodyOpenAI
	err := decodeJSON(filePath, &request)
	if err != nil {
		fmt.Println("Error decoding JSON.\n[ERROR] -", err)
		return
	}

	msg := Message{
		Role:    role,
		Content: response,
	}
	request.Messages = append(request.Messages, msg)

	// Write back to the JSON file
	err2 := encodeJSON(filePath, request)
	if err2 != nil {
		fmt.Println("Error encoding JSON.\n[ERROR] -", err2)
		return
	}
}
