package chat

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
	"jervis/clients"
	"jervis/clients/anthropic"
	"jervis/clients/openai"
	"jervis/internal"
	"jervis/models"
	"jervis/storage"
)

func DoChat() error {

	vendor := viper.GetString("vendor")
	filePath := viper.GetString("chat.fileName") + "_" + vendor + ".json"
	model := viper.GetString("model")
	systemMessage := viper.GetString("chat.systemMessage")
	newSession := viper.GetBool("newSession")

	if viper.GetBool("format") {
		internal.FormatOut(filePath)
		return nil
	}

	var client clients.ChatClient
	switch vendor {
	case "openai":
		client = openai.NewClient(model, systemMessage, viper.GetString("openai_api_key"))
	case "anthropic":
		client = anthropic.NewClient(model, systemMessage, viper.GetString("anthropic_api_key"))
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

	// bundle := RequestBundle{
	// 	vendor:        Vendor(vendor),
	// 	model:         Model(viper.GetString("model")),
	// 	systemMessage: viper.GetString("chat.systemMessage"),
	// 	prompt:        prompt,
	// }

	// requestBody := getRequest(
	// 	filePath,
	// 	viper.GetBool("newSession"),
	// 	bundle)
	//
	// if len(requestBody) == 0 {
	// 	return
	// }
	//
	// response := performRequest(viper.GetString("url"), requestBody, viper.GetString("api_key"))
	//
	// if len(response) > 0 {
	// 	persistConversation(filePath, "assistant", response)
	//
	// }
	return nil
}
