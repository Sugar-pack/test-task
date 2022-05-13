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

func (h *CompanyHandler) UpdateCompanies(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := logging.FromContext(ctx)
	logger.Info("CreateCompany begin")
	updateCompany := &model.CompanyForUpdate{}
	err := json.NewDecoder(request.Body).Decode(updateCompany)
	if err != nil {
		logger.WithError(err).Error("Decode error")
		helper.BadRequest(ctx, writer, "Cant decode request body")

		return
	}
	updateCompanyForFilter := MapJSONUpdateToDB(updateCompany)
	updatedRows, err := h.CompanyRepository.UpdateCompany(ctx, updateCompanyForFilter)
	if err != nil {
		logger.WithError(err).Error("UpdateCompany repository error")
		helper.InternalError(ctx, writer, "Cant update company")

		return
	}
	err = h.Producer.PublishMessage(ctx, sender.JSONType,
		model.NewMessage(http.StatusOK, fmt.Sprintf("%d rows updated", updatedRows)))
	if err != nil {
		logger.WithError(err).Error("Cant publish message")
	}
	helper.StatusOk(ctx, writer, fmt.Sprintf("%d rows updated", updatedRows))
}

func MapJSONUpdateToDB(companyUpdate *model.CompanyForUpdate) *repository.CompanyForUpdate {
	return &repository.CompanyForUpdate{
		FilterFields: repository.CompanyForFilter{
			Name:    companyUpdate.FilterFields.Name,
			Code:    companyUpdate.FilterFields.Code,
			Country: companyUpdate.FilterFields.Country,
			Website: companyUpdate.FilterFields.Website,
			Phone:   companyUpdate.FilterFields.Phone,
		},
		FieldsForUpdate: repository.CompanyUpdatable{
			Code:    companyUpdate.FieldsForUpdate.Code,
			Country: companyUpdate.FieldsForUpdate.Country,
			Website: companyUpdate.FieldsForUpdate.Website,
			Phone:   companyUpdate.FieldsForUpdate.Phone,
		},
	}
}
