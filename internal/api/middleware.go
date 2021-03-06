package api

import (
	"errors"
	"fmt"
	"net"
	"net/http"

	"github.com/Sugar-pack/test-task/internal/helper"
	"github.com/Sugar-pack/test-task/internal/logging"
)

func LoggingMiddleware(logger logging.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctxWithLogger := logging.WithContext(ctx, logger)
			r = r.WithContext(ctxWithLogger)
			next.ServeHTTP(w, r)
		})
	}
}

func WithLogRequestBoundaries() func(next http.Handler) http.Handler {
	httpMw := func(next http.Handler) http.Handler {
		handlerFn := func(writer http.ResponseWriter, request *http.Request) {
			ctx := request.Context()
			logger := logging.FromContext(ctx)
			requestURI := request.RequestURI
			requestMethod := request.Method
			logRequest := fmt.Sprintf("%s %s", requestMethod, requestURI)
			logger.WithField("request", logRequest).Trace("REQUEST_STARTED")
			next.ServeHTTP(writer, request)
			logger.WithField("request", logRequest).Trace("REQUEST_COMPLETED")
		}

		return http.HandlerFunc(handlerFn)
	}

	return httpMw
}

func CountryAccessMiddleware(qualifier CountryQualifier) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			ctx := request.Context()
			logger := logging.FromContext(ctx)
			adr := request.RemoteAddr
			method := request.Method
			if method != http.MethodPost && method != http.MethodDelete {
				next.ServeHTTP(writer, request)

				return
			}

			logger.Info("request from ", adr)
			userIP, err := IPbyRequest(request)
			if err != nil {
				logger.WithError(err).Error("failed to get user ip")
				helper.InternalError(ctx, writer, "cant get ip")

				return
			}
			logger.Info("user ip ", userIP)
			logger.Info("ip")
			isAllow := qualifier.QualifyCountry(ctx, userIP)
			if !isAllow {
				helper.Forbidden(ctx, writer, "country not allowed")

				return
			}
			next.ServeHTTP(writer, request)
		})
	}
}

func IPbyRequest(req *http.Request) (string, error) {
	ip, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		return "", fmt.Errorf("cant split ip: %w", err)
	}
	userIP := net.ParseIP(ip)
	if userIP == nil {
		return "", errors.New(fmt.Sprintf("invalid ip: %s", ip)) //nolint:goerr113 // we want to return error with ip
	}

	return userIP.String(), nil
}
