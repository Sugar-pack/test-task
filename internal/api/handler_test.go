package api

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Sugar-pack/test-task/internal/model"

	"github.com/stretchr/testify/mock"

	"github.com/gavv/httpexpect/v2"
	"github.com/stretchr/testify/suite"

	"github.com/Sugar-pack/test-task/internal/logging"
	"github.com/Sugar-pack/test-task/internal/mocks/qualifier"
	"github.com/Sugar-pack/test-task/internal/mocks/repository"
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
