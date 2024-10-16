package edit

import (
	"fmt"
	"jervis/clients"
	"jervis/clients/anthropic"
	"jervis/clients/openai"
	"jervis/internal"
	"jervis/models"
	"jervis/storage"
	"strings"

	"github.com/spf13/viper"
)

func DoEdit() error {

	filePath := viper.GetString("editor.fileName") + ".json"
	vendor := viper.GetString("vendor")
	apiKey := viper.GetString("api_key")
	model := viper.GetString("model")
	systemMessage := viper.GetString("edit.systemMessage")
	newSession := viper.GetBool("newSession")

	if viper.GetBool("format") {
		internal.FormatOut(filePath)
		return nil
	}

	var client clients.ChatClient
	switch vendor {
	case "openai":
		client = openai.NewClient(model, systemMessage, apiKey)
	case "anthropic":
		client = anthropic.NewClient(model, systemMessage, apiKey)
	default:
		return fmt.Errorf("unsupported vendor type: %s", vendor)
	}

	conversation, err := storage.LoadConversation(filePath, newSession)
	if err != nil {
		return err
	}

	prompt := strings.Join(internal.HandleInput(viper.GetBool("readlines")), "\n")
	if prompt == "" {
		return nil
	}
	userMessage := models.Message{
		Role:    "user",
		Content: prompt,
	}

	conversation = append(conversation, userMessage)

	response, err := client.SendMessage(conversation)
	if err != nil {
		return err
	}

	assistantMessage := models.Message{
		Role:    "assistant",
		Content: response.Content,
	}

	conversation = append(conversation, assistantMessage)

	if err := storage.SaveConversation(filePath, conversation); err != nil {
		return err
	}
	return nil
}
