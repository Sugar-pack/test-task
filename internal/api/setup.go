package api

import (
	"fmt"

	"github.com/go-chi/chi/v5"

	"github.com/Sugar-pack/test-task/internal/logging"
)

func SetupRouter(logger logging.Logger, handler *CompanyHandler, qualifier CountryQualifier) *chi.Mux {
	router := chi.NewRouter()
	router.Use(LoggingMiddleware(logger), WithLogRequestBoundaries(), CountryAccessMiddleware(qualifier))
	router.Post("/companies/create", handler.CreateCompany)
	filterString := "name={name}&code={code}&country={country}&website={website}&phone={phone}"
	router.Get(fmt.Sprintf("/companies/%s/", filterString), handler.GetCompany)
	router.Delete(fmt.Sprintf("/companies/%s/", filterString), handler.DeleteCompanies)
	router.Patch("/companies/update", handler.UpdateCompanies)

	return router
}
