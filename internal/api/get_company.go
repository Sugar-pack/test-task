package api

import (
	"encoding/json"
	"net/http"

	"github.com/Sugar-pack/test-task/internal/helper"

	"github.com/go-chi/chi/v5"

	"github.com/Sugar-pack/test-task/internal/logging"
	"github.com/Sugar-pack/test-task/internal/model"
	"github.com/Sugar-pack/test-task/internal/repository"
)

func (h *CompanyHandler) GetCompany(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := logging.FromContext(ctx)
	logger.Info("CreateCompany begin")
	companyForFilter := CompanyFilterFromRequest(request)
	companies, err := h.CompanyRepository.GetCompany(ctx, &companyForFilter)
	if err != nil {
		logger.WithError(err).Error("GetCompany repository error")
		helper.InternalError(ctx, writer, "Cant get company")

		return
	}
	if len(companies) == 0 {
		helper.NotFound(ctx, writer, "Companies not found")

		return
	}
	var companiesJSON []model.Company
	for _, comp := range companies {
		companiesJSON = append(companiesJSON, MapDBCompanyToJSON(&comp))
	}
	err = json.NewEncoder(writer).Encode(companies)
	if err != nil {
		logger.WithError(err).Error("Encode error")
		helper.InternalError(ctx, writer, "Cant encode response")

		return
	}
	helper.StatusOk(ctx, writer, "")
}

func MapDBCompanyToJSON(company *repository.Company) model.Company {
	return model.Company{
		Name:    company.Name,
		Code:    company.Code,
		Country: company.Country,
		Website: company.Website,
		Phone:   company.Phone,
	}
}

func CompanyFilterFromRequest(request *http.Request) repository.CompanyForFilter {
	companyName := chi.URLParam(request, "name")
	companyCode := chi.URLParam(request, "code")
	companyCountry := chi.URLParam(request, "country")
	companyWebsite := chi.URLParam(request, "website")
	companyPhone := chi.URLParam(request, "phone")

	return repository.CompanyForFilter{
		Name:    companyName,
		Code:    companyCode,
		Country: companyCountry,
		Website: companyWebsite,
		Phone:   companyPhone,
	}
}
