package model

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMessage_Success(t *testing.T) {
	msg := "test message"
	message := NewMessage(http.StatusOK, msg)
	assert.Equal(t, SuccessType, message.Type)
	assert.Equal(t, msg, message.Data)
}

func TestNewMessage_Error(t *testing.T) {
	msg := "test message"
	message := NewMessage(http.StatusInternalServerError, msg)
	assert.Equal(t, ErrorType, message.Type)
	assert.Equal(t, msg, message.Data)
}
