package api

import (
	"errors"
	"fmt"
	"net"
	"net/http"

	"github.com/Sugar-pack/test-task/internal/handler"
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

func CountryAccessMiddleware(qualifier CountryQualifier, whiteList []string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			ctx := request.Context()
			logger := logging.FromContext(ctx)
			adr := request.RemoteAddr
			logger.Info("request from ", adr)
			userIP, err := IPbyRequest(request)
			if err != nil {
				logger.WithError(err).Error("failed to get user ip")
				handler.InternalError(ctx, writer, "cant get ip")

				return
			}
			logger.Info("user ip ", userIP)
			logger.Info("ip")
			country, err := qualifier.QualifyCountry(ctx, userIP)
			if err != nil {
				logger.WithError(err).Error("cant get country")
				handler.InternalError(ctx, writer, "cant get country")

				return
			}
			for _, countryFromWL := range whiteList {
				if country == countryFromWL {
					next.ServeHTTP(writer, request)

					return
				}
			}
			handler.Forbidden(ctx, writer, "country not allowed")

			return
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
