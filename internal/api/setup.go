package api

import (
	"github.com/Sugar-pack/test-task/internal/handler"
	"github.com/Sugar-pack/test-task/internal/logging"
	"github.com/go-chi/chi/v5"
)

func SetupRouter(logger logging.Logger, handler *handler.CompanyHandler, whiteList []string) *chi.Mux {
	router := chi.NewRouter()
	qualifier := &IPAPICountryQualifier{}
	router.Use(LoggingMiddleware(logger), WithLogRequestBoundaries(), CountryAccessMiddleware(qualifier, whiteList))
	router.Post("/companies", handler.CreateCompany)
	router.Post("/companies/search", handler.GetCompany)
	router.Post("/companies/delete", handler.DeleteCompanies)
	router.Post("/companies/update", handler.UpdateCompanies)

	return router
}
