package internal

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/chzyer/readline"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatCompletionRequestBody struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
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

func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func checkModel(model string) bool {
	validModel := false
	models := []string{"gpt-3.5-turbo-1106", "gpt-4-1106-preview"}
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

func getFileContents(filePath string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}
	return string(data), nil
}

func decodeJSON(filePath string, body *ChatCompletionRequestBody) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("Error opening JSON file %s", err)
	}
	defer file.Close()

	if err := json.NewDecoder(file).Decode(body); err != nil {
		return fmt.Errorf("Error decoding JSON file %s", err)
	}
	return nil
}

func encodeJSON(filePath string, body ChatCompletionRequestBody) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("Error creating JSON file %s", err)
	}
	defer file.Close()

	if err := json.NewEncoder(file).Encode(body); err != nil {
		return fmt.Errorf("Error encoding JSON file %s", err)
	}
	return nil
}

func handleInput(prompt bool) []string {
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

func printFormattedOutput(filePath string) {
	var body ChatCompletionRequestBody
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

func createRequest(filePath string, newSession bool, model string, systemMessage string, question string) []byte {
	var requestBody []byte
	if newSession {
		if fileExists(filePath) {
			// delete the JSON file
			err := os.Remove(filePath)
			if err != nil {
				panic(err)
			}
		}
		// Open the file for writing
		file, _ := os.Create(filePath)
		defer file.Close()
		// Write the marshaled data to the file
		requestBody = getRequestBody(model, systemMessage, question)
		file.Write(requestBody)
	} else {
		if fileExists(filePath) {
			var request ChatCompletionRequestBody
			err := decodeJSON(filePath, &request)
			if err != nil {
				fmt.Println("Error decoding JSON.\n[ERROR] -", err)
				os.Exit(1)
			}

			// Append data to slice
			msg := Message{
				Role:    "user",
				Content: question,
			}
			request.Messages = append(request.Messages, msg)

			// cfg.model = request.Model

			requestBody, _ = json.Marshal(request)

			err2 := encodeJSON(filePath, request)
			if err2 != nil {
				fmt.Println("Error encoding JSON.\n[ERROR] -", err2)
				os.Exit(1)
			}
		} else {
			requestBody = getRequestBody(model, systemMessage, question)
		}
	}
	return requestBody
}

func getRequestBody(model string, systemMessge string, question string) []byte {
	body := map[string]interface{}{
		"model": model,
		"messages": []map[string]string{
			{
				"role":    "system",
				"content": systemMessge,
			},
			{
				"role":    "user",
				"content": question,
			},
		},
		"stream": true,
	}
	requestBody, _ := json.Marshal(body)
	return requestBody
}

func performAPIRequest(url string, requestBody []byte, key string) string {

	// Create a new request
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+key)

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
			fmt.Println("Error unmarshalling JSON.\n[ERROR] -", err2, string(jsonData))
			os.Exit(1)
		}

		if chunk.Choices[0].FinishReason == "stop" {
			break
		}

		words := chunk.Choices[0].Delta.Content
		fmt.Print(words)
		responses = append(responses, words)
	}
	response := strings.Join(responses, "")
	return response
}

func persistConversation(filePath string, role string, response string) {
	if !fileExists(filePath) {
		return
	}
	var request ChatCompletionRequestBody
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
