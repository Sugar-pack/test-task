package api

import (
	"fmt"
	"github.com/Sugar-pack/test-task/internal/handler"
	"github.com/Sugar-pack/test-task/internal/logging"
	"net"
	"net/http"
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
		handlerFn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			logger := logging.FromContext(ctx)
			requestURI := r.RequestURI
			requestMethod := r.Method
			logRequest := fmt.Sprintf("%s %s", requestMethod, requestURI)
			logger.WithField("request", logRequest).Trace("REQUEST_STARTED")
			next.ServeHTTP(w, r)
			logger.WithField("request", logRequest).Trace("REQUEST_COMPLETED")
		}
		return http.HandlerFunc(handlerFn)
	}
	return httpMw
}

func CountryAccessMiddleware(c CountryQualifier, whiteList []string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			logger := logging.FromContext(ctx)
			adr := r.RemoteAddr
			logger.Info("request from ", adr)
			userIP, err := IPbyRequest(r)
			if err != nil {
				logger.WithError(err).Error("failed to get user ip")
				handler.InternalError(ctx, w, "cant get ip")
				return
			}
			logger.Info("user ip ", userIP)
			logger.Info("ip")
			country, err := c.QualifyCountry(ctx, userIP)
			if err != nil {
				logger.WithError(err).Error("cant get country")
				handler.InternalError(ctx, w, "cant get country")
				return
			}
			for _, countryFromWL := range whiteList {
				if country == countryFromWL {
					next.ServeHTTP(w, r)
					return
				}
			}
			handler.Forbidden(ctx, w, "country not allowed")
			return
		})
	}
}

func IPbyRequest(req *http.Request) (string, error) {
	ip, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		return "", err
	}
	userIP := net.ParseIP(ip)
	if userIP == nil {
		return "", fmt.Errorf("invalid ip: %s", ip)
	}
	return userIP.String(), nil
}
