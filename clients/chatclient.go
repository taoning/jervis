package clients

import "jervis/models"

type ChatClient interface {
	SendMessage(conversation []models.Message) (models.Response, error)
}
