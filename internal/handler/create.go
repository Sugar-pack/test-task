package handler

import (
	"encoding/json"
	"github.com/Sugar-pack/test-task/internal/logging"
	"github.com/Sugar-pack/test-task/internal/model"
	"github.com/Sugar-pack/test-task/internal/repository"
	"net/http"
)

type CompanyHandler struct {
	CompanyRepository repository.CompanyRepository
}

func NewCompanyHandler(companyRepository repository.CompanyRepository) *CompanyHandler {
	return &CompanyHandler{
		CompanyRepository: companyRepository,
	}
}

func (h *CompanyHandler) CreateCompany(writer http.ResponseWriter, request *http.Request) {
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
	err = h.CompanyRepository.CreateCompany(ctx, company)
	if err != nil {
		logger.WithError(err).Error("CreateCompany repository error")
		InternalError(ctx, writer, "Cant create company")
		return
	}
	StatusOk(ctx, writer, "Company created")
}