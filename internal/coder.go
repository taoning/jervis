package internal

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

func DoCode() {

	filePath := viper.GetString("coder.fileName") + ".json"

	if viper.GetBool("format") {
		printFormattedOutput(filePath)
		return
	}

	qlines := handleInput(viper.GetBool("readlines"))
	model := viper.GetString("model")
	validModel := checkModel(model)
	if !validModel {
		fmt.Println("Invalid model: ", model)
		return
	}

	question := strings.Join(qlines, "\n")
    if question == "" {
        return
    }

	requestBody := createRequest(filePath, viper.GetBool("newSession"), model, viper.GetString("coder.systemMessage"), question)

	if len(requestBody) == 0 {
		return
	}

	response := performAPIRequest(viper.GetString("url"), requestBody, viper.GetString("api_key"))

	if len(response) > 0 {
		persistConversation(filePath, "assistant", response)

	}
}
