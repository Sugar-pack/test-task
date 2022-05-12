package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Sugar-pack/test-task/internal/logging"
	"github.com/Sugar-pack/test-task/internal/model"
	"github.com/Sugar-pack/test-task/internal/repository"
)

func (h *CompanyHandler) GetCompany(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := logging.FromContext(ctx)
	logger.Info("CreateCompany begin")
	company := &model.Company{}
	err := json.NewDecoder(request.Body).Decode(company)
	if err != nil {
		logger.WithError(err).Error("Decode error")
		BadRequest(ctx, writer, "Cant decode request body")

		return
	}
	companyForFilter := MapJSONToFilter(company)
	companies, err := h.CompanyRepository.GetCompany(ctx, &companyForFilter)
	if err != nil {
		logger.WithError(err).Error("GetCompany repository error")
		InternalError(ctx, writer, "Cant get company")

		return
	}
	if len(companies) == 0 {
		NotFound(ctx, writer, "Companies not found")

		return
	}
	logger.Info(companies)
	var companiesJSON []model.Company
	for _, comp := range companies {
		companiesJSON = append(companiesJSON, MapDBCompanyToJSON(&comp))
	}
	err = json.NewEncoder(writer).Encode(companies)
	if err != nil {
		logger.WithError(err).Error("Encode error")
		InternalError(ctx, writer, "Cant encode response")

		return
	}
	StatusOk(ctx, writer, "")
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

func MapJSONToFilter(company *model.Company) repository.CompanyForFilter {
	return repository.CompanyForFilter{
		Name:    company.Name,
		Code:    company.Code,
		Country: company.Country,
		Website: company.Website,
		Phone:   company.Phone,
	}
}
