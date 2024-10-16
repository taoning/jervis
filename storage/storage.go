package storage

import (
	"encoding/json"
	"fmt"
	"jervis/models"
	"os"
)

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

func LoadConversation(filePath string, newSession bool) ([]models.Message, error) {
	var message []models.Message
	if newSession {
		return message, nil
	}
	decodeJSON(filePath, &message)
	return message, nil
}

func SaveConversation(filePath string, conversation []models.Message) error {
	err := encodeJSON(filePath, conversation)
	return err
}
