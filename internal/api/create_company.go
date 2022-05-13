package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Sugar-pack/test-task/internal/sender"

	"github.com/Sugar-pack/test-task/internal/helper"

	"github.com/Sugar-pack/test-task/internal/logging"
	"github.com/Sugar-pack/test-task/internal/model"
	"github.com/Sugar-pack/test-task/internal/repository"
)

type CompanyHandler struct {
	CompanyRepository repository.CompanyRepository
	Producer          sender.Producer
}

func NewCompanyHandler(companyRepository repository.CompanyRepository, producer sender.Producer) *CompanyHandler {
	return &CompanyHandler{
		CompanyRepository: companyRepository,
		Producer:          producer,
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
		helper.BadRequest(ctx, writer, "Cant decode request body")

		return
	}
	err = h.CompanyRepository.CreateCompany(ctx, MapJSONCompanyToDB(company))
	if err != nil {
		logger.WithError(err).Error("CreateCompany repository error")
		helper.InternalError(ctx, writer, "Cant create company")

		return
	}
	err = h.Producer.PublishMessage(ctx, sender.JSONType, model.NewMessage(http.StatusOK,
		fmt.Sprintf("Company %s created", company.Name)))
	if err != nil {
		logger.WithError(err).Error("Cant publish message")
	}
	helper.StatusOk(ctx, writer, "Company created")
}

func MapJSONCompanyToDB(company *model.Company) *repository.Company {
	return &repository.Company{
		Name:    company.Name,
		Code:    company.Code,
		Country: company.Country,
		Website: company.Website,
		Phone:   company.Phone,
	}
}
