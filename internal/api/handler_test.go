package api

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gavv/httpexpect/v2"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/Sugar-pack/test-task/internal/logging"
	"github.com/Sugar-pack/test-task/internal/mocks/qualifier"
	"github.com/Sugar-pack/test-task/internal/mocks/repository"
	"github.com/Sugar-pack/test-task/internal/model"
	repo "github.com/Sugar-pack/test-task/internal/repository"
)

const localhost = "127.0.0.1"

type CompanyTestSuite struct {
	suite.Suite
	repo      *repository.CompanyRepository
	qualifier *qualifier.CountryQualifier
	server    *httptest.Server
}

func TestOrderStatusSuite(t *testing.T) {
	s := new(CompanyTestSuite)
	suite.Run(t, s)
}

func (s *CompanyTestSuite) SetupTest() {
	logger := logging.GetLogger()
	repo := &repository.CompanyRepository{}
	q := &qualifier.CountryQualifier{}
	s.repo = repo
	s.qualifier = q
	handler := NewCompanyHandler(repo)
	router := SetupRouter(logger, handler, q)
	s.server = httptest.NewServer(router)
}

func (s *CompanyTestSuite) TestCompanyHandler_CreateCompany_DecodeErr() {
	t := s.T()

	httpExpect := httpexpect.New(t, s.server.URL)
	badBody := []byte("bad body")
	s.qualifier.On("QualifyCountry", mock.AnythingOfType("*context.valueCtx"), localhost).Return(true)
	httpExpect.POST("/companies/create").WithJSON(badBody).Expect().Status(http.StatusBadRequest)
	s.qualifier.AssertExpectations(t)
}

func (s *CompanyTestSuite) TestCompanyHandler_CreateCompany_RepoErr() {
	t := s.T()

	httpExpect := httpexpect.New(t, s.server.URL)
	company := model.Company{
		Name:    "test name",
		Code:    "code",
		Country: "country",
		Website: "website",
		Phone:   "phone",
	}
	err := errors.New("error")
	s.qualifier.On("QualifyCountry", mock.AnythingOfType("*context.valueCtx"), localhost).
		Return(true)
	s.repo.On("CreateCompany", mock.AnythingOfType("*context.valueCtx"), MapJSONCompanyToDB(&company)).
		Return(err)
	httpExpect.POST("/companies/create").WithJSON(company).Expect().Status(http.StatusInternalServerError)
}

func (s *CompanyTestSuite) TestCompanyHandler_CreateCompany_OK() {
	t := s.T()

	httpExpect := httpexpect.New(t, s.server.URL)
	company := model.Company{
		Name:    "test name",
		Code:    "code",
		Country: "country",
		Website: "website",
		Phone:   "phone",
	}
	s.qualifier.On("QualifyCountry", mock.AnythingOfType("*context.valueCtx"), localhost).
		Return(true)
	s.repo.On("CreateCompany", mock.AnythingOfType("*context.valueCtx"), MapJSONCompanyToDB(&company)).
		Return(nil)
	httpExpect.POST("/companies/create").WithJSON(company).Expect().Status(http.StatusOK)
}

func (s *CompanyTestSuite) TestCompanyHandler_UpdateCompany_DecodeErr() {
	t := s.T()

	httpExpect := httpexpect.New(t, s.server.URL)
	badBody := []byte("bad body")
	httpExpect.PATCH("/companies/update").WithJSON(badBody).Expect().Status(http.StatusBadRequest)
}

func (s *CompanyTestSuite) TestCompanyHandler_UpdateCompany_RepoErr() {
	t := s.T()

	httpExpect := httpexpect.New(t, s.server.URL)
	companyUpdate := model.CompanyForUpdate{
		FilterFields: model.Company{
			Name:    "name",
			Code:    "code",
			Country: "country",
			Website: "website",
			Phone:   "phone",
		},
		FieldsForUpdate: model.Company{
			Name:    "new name",
			Code:    "new code",
			Country: "new country",
			Website: "new website",
			Phone:   "new phone",
		},
	}
	err := errors.New("error")
	s.repo.On("UpdateCompany", mock.AnythingOfType("*context.valueCtx"), MapJSONUpdateToDB(&companyUpdate)).
		Return(int64(0), err)
	httpExpect.PATCH("/companies/update").WithJSON(companyUpdate).Expect().Status(http.StatusInternalServerError)
	s.repo.AssertExpectations(t)
}

func (s *CompanyTestSuite) TestCompanyHandler_UpdateCompany_OK() {
	t := s.T()

	httpExpect := httpexpect.New(t, s.server.URL)
	companyUpdate := model.CompanyForUpdate{
		FilterFields: model.Company{
			Name:    "name",
			Code:    "code",
			Country: "country",
			Website: "website",
			Phone:   "phone",
		},
		FieldsForUpdate: model.Company{
			Name:    "new name",
			Code:    "new code",
			Country: "new country",
			Website: "new website",
			Phone:   "new phone",
		},
	}
	s.repo.On("UpdateCompany", mock.AnythingOfType("*context.valueCtx"), MapJSONUpdateToDB(&companyUpdate)).
		Return(int64(1), nil)
	httpExpect.PATCH("/companies/update").WithJSON(companyUpdate).Expect().Status(http.StatusOK)
	s.repo.AssertExpectations(t)
}

func (s *CompanyTestSuite) TestCompanyHandler_DeleteCompany_NoAccess() {
	t := s.T()

	httpExpect := httpexpect.New(t, s.server.URL)
	company := repo.CompanyForFilter{
		Name:    "name",
		Code:    "code",
		Country: "country",
		Website: "website",
		Phone:   "phone",
	}

	s.qualifier.On("QualifyCountry", mock.AnythingOfType("*context.valueCtx"), localhost).
		Return(false)
	path := fmt.Sprintf("/companies/name=%s&code=%s&country=%s&website=%s&phone=%s/",
		company.Name, company.Code, company.Country, company.Website, company.Phone)
	httpExpect.DELETE(path).Expect().Status(http.StatusForbidden)
	s.qualifier.AssertExpectations(t)
}

func (s *CompanyTestSuite) TestCompanyHandler_DeleteCompany_RepoErr() {
	t := s.T()

	httpExpect := httpexpect.New(t, s.server.URL)
	company := repo.CompanyForFilter{
		Name:    "name",
		Code:    "code",
		Country: "country",
		Website: "website",
		Phone:   "phone",
	}

	s.qualifier.On("QualifyCountry", mock.AnythingOfType("*context.valueCtx"), localhost).
		Return(true)
	path := fmt.Sprintf("/companies/name=%s&code=%s&country=%s&website=%s&phone=%s/",
		company.Name, company.Code, company.Country, company.Website, company.Phone)
	err := errors.New("error")
	s.repo.On("DeleteCompany", mock.AnythingOfType("*context.valueCtx"), &company).Return(int64(0), err)
	httpExpect.DELETE(path).Expect().Status(http.StatusInternalServerError)
	s.repo.AssertExpectations(t)
	s.qualifier.AssertExpectations(t)
}

func (s *CompanyTestSuite) TestCompanyHandler_DeleteCompany_OK() {
	t := s.T()

	httpExpect := httpexpect.New(t, s.server.URL)
	company := repo.CompanyForFilter{
		Name:    "name",
		Code:    "code",
		Country: "country",
		Website: "website",
		Phone:   "phone",
	}

	s.qualifier.On("QualifyCountry", mock.AnythingOfType("*context.valueCtx"), localhost).
		Return(true)
	path := fmt.Sprintf("/companies/name=%s&code=%s&country=%s&website=%s&phone=%s/",
		company.Name, company.Code, company.Country, company.Website, company.Phone)
	s.repo.On("DeleteCompany", mock.AnythingOfType("*context.valueCtx"), &company).Return(int64(1), nil)
	httpExpect.DELETE(path).Expect().Status(http.StatusOK)
	s.repo.AssertExpectations(t)
	s.qualifier.AssertExpectations(t)
}

func (s *CompanyTestSuite) TestCompanyHandler_GetCompany_RepoErr() {
	t := s.T()

	httpExpect := httpexpect.New(t, s.server.URL)
	company := repo.CompanyForFilter{
		Name:    "name",
		Code:    "code",
		Country: "country",
		Website: "website",
		Phone:   "phone",
	}

	path := fmt.Sprintf("/companies/name=%s&code=%s&country=%s&website=%s&phone=%s/",
		company.Name, company.Code, company.Country, company.Website, company.Phone)
	err := errors.New("error")
	s.repo.On("GetCompany", mock.AnythingOfType("*context.valueCtx"), &company).Return(nil, err)
	httpExpect.GET(path).Expect().Status(http.StatusInternalServerError)
	s.repo.AssertExpectations(t)
}

func (s *CompanyTestSuite) TestCompanyHandler_GetCompany_NoResult() {
	t := s.T()

	httpExpect := httpexpect.New(t, s.server.URL)
	company := repo.CompanyForFilter{
		Name:    "name",
		Code:    "code",
		Country: "country",
		Website: "website",
		Phone:   "phone",
	}

	path := fmt.Sprintf("/companies/name=%s&code=%s&country=%s&website=%s&phone=%s/",
		company.Name, company.Code, company.Country, company.Website, company.Phone)
	s.repo.On("GetCompany", mock.AnythingOfType("*context.valueCtx"), &company).Return(nil, nil)
	httpExpect.GET(path).Expect().Status(http.StatusNotFound)
	s.repo.AssertExpectations(t)
}

func (s *CompanyTestSuite) TestCompanyHandler_GetCompany_OK() {
	t := s.T()

	httpExpect := httpexpect.New(t, s.server.URL)
	company := repo.CompanyForFilter{
		Name:    "name",
		Code:    "code",
		Country: "country",
		Website: "website",
		Phone:   "phone",
	}

	path := fmt.Sprintf("/companies/name=%s&code=%s&country=%s&website=%s&phone=%s/",
		company.Name, company.Code, company.Country, company.Website, company.Phone)
	returnCompanies := []repo.Company{
		{
			Name: "name",
		},
	}
	s.repo.On("GetCompany", mock.AnythingOfType("*context.valueCtx"), &company).Return(returnCompanies, nil)
	httpExpect.GET(path).Expect().Status(http.StatusOK).JSON().Equal(returnCompanies)
	s.repo.AssertExpectations(t)
}
