package api

import (
	"fmt"
	"net/http"

	"github.com/Sugar-pack/test-task/internal/helper"

	"github.com/Sugar-pack/test-task/internal/logging"
)

func (h *CompanyHandler) DeleteCompanies(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := logging.FromContext(ctx)
	logger.Info("CreateCompany begin")
	companyForFilter := CompanyFilterFromRequest(request)
	deletedRows, err := h.CompanyRepository.DeleteCompany(ctx, &companyForFilter)
	if err != nil {
		logger.WithError(err).Error("DeleteCompany repository error")
		helper.InternalError(ctx, writer, "Cant delete company")

		return
	}

	helper.StatusOk(ctx, writer, fmt.Sprintf("%d rows deleted", deletedRows))
}
