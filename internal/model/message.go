package model

import "net/http"

type Message struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

func NewMessage(code int, msg string) *Message {
	message := Message{}
	if code >= http.StatusOK && code < http.StatusMultipleChoices {
		message.Type = "Success"
	} else {
		message.Type = "Error"
	}
	message.Data = msg

	return &message
}
