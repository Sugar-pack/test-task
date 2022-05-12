package handler

import (
	"bytes"
	"context"
	"net/http"
	"strings"

	"github.com/Sugar-pack/test-task/internal/logging"
)

const ErrMsgWritingResponse = "Error while writing response"

func BadRequest(ctx context.Context, writer http.ResponseWriter, msg string) {
	logger := logging.FromContext(ctx)
	writer.WriteHeader(http.StatusBadRequest)
	_, wErr := writer.Write([]byte(msg))
	if wErr != nil {
		logger.WithError(wErr).Error(ErrMsgWritingResponse)
	}
}

func InternalError(ctx context.Context, writer http.ResponseWriter, s string) {
	logger := logging.FromContext(ctx)
	writer.WriteHeader(http.StatusInternalServerError)
	_, wErr := writer.Write([]byte(s))
	if wErr != nil {
		logger.WithError(wErr).Error(ErrMsgWritingResponse)
	}
}

func StatusOk(ctx context.Context, writer http.ResponseWriter, s string) {
	logger := logging.FromContext(ctx)
	writer.WriteHeader(http.StatusOK)
	_, wErr := writer.Write([]byte(s))
	if wErr != nil {
		logger.WithError(wErr).Error(ErrMsgWritingResponse)
	}
}

func Forbidden(ctx context.Context, writer http.ResponseWriter, msg string) {
	logger := logging.FromContext(ctx)
	writer.WriteHeader(http.StatusForbidden)
	_, wErr := writer.Write([]byte(msg))
	if wErr != nil {
		logger.WithError(wErr).Error(ErrMsgWritingResponse)
	}
}

func NotFound(ctx context.Context, writer http.ResponseWriter, msg string) {
	logger := logging.FromContext(ctx)
	body := strings.NewReader(msg)
	buff := new(bytes.Buffer)
	if _, err := buff.ReadFrom(body); err != nil {
		logger.WithError(err).Error(ErrMsgWritingResponse)

		return
	}
	rawResponse(ctx, writer, http.StatusNotFound, nil, buff.Bytes())
}

func rawResponse(ctx context.Context, writer http.ResponseWriter, httpCode int, httpHeaders http.Header, body []byte) {
	logger := logging.FromContext(ctx)
	for k, vs := range httpHeaders {
		for _, v := range vs {
			writer.Header().Add(k, v)
		}
	}
	writer.WriteHeader(httpCode)
	if _, wErr := writer.Write(body); wErr != nil {
		logger.WithError(wErr).Error(ErrMsgWritingResponse)
	}
}
