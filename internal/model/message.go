package model

import "net/http"

type Message struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

const (
	SuccessType = "Success"
	ErrorType   = "Error"
)

func NewMessage(code int, msg string) *Message {
	message := Message{}
	if code >= http.StatusOK && code < http.StatusMultipleChoices {
		message.Type = SuccessType
	} else {
		message.Type = ErrorType
	}
	message.Data = msg

	return &message
}
