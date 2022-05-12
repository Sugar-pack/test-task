package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Sugar-pack/test-task/internal/logging"
	"github.com/Sugar-pack/test-task/internal/model"
)

func (h *CompanyHandler) DeleteCompanies(writer http.ResponseWriter, request *http.Request) {
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
	deletedRows, err := h.CompanyRepository.DeleteCompany(ctx, &companyForFilter)
	if err != nil {
		logger.WithError(err).Error("DeleteCompany repository error")
		InternalError(ctx, writer, "Cant delete company")

		return
	}

	StatusOk(ctx, writer, fmt.Sprintf("%d rows deleted", deletedRows))
}
