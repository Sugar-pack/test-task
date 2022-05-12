package api

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type CountryQualifierTestSuite struct {
	suite.Suite
	qualifier *IPAPICountryQualifier
}

func TestCountryQualifierTestSuite(t *testing.T) {
	s := new(CountryQualifierTestSuite)
	suite.Run(t, s)
}

func (s *CountryQualifierTestSuite) SetupTest() {
	WhiteList := []string{"Russia", "United States"}
	s.qualifier = NewIPAPI(WhiteList)
}

func (s *CountryQualifierTestSuite) TestIPAPICountryQualifier_QualifyCountry() {
	t := s.T()
	ctx := context.Background()
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	ip := "some-ip"
	httpmock.RegisterResponder(http.MethodGet, fmt.Sprintf(urlPattern, ip),
		httpmock.NewStringResponder(http.StatusNotFound, ""))
	isAllow := s.qualifier.QualifyCountry(ctx, ip)
	assert.False(t, isAllow)
}

func (s *CountryQualifierTestSuite) TestIPAPICountryQualifier_NotInList() {
	t := s.T()
	ctx := context.Background()
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	ip := "some-ip"
	httpmock.RegisterResponder(http.MethodGet, fmt.Sprintf(urlPattern, ip),
		httpmock.NewStringResponder(http.StatusNotFound, "United Kingdom"))
	isAllow := s.qualifier.QualifyCountry(ctx, ip)
	assert.False(t, isAllow)
}

func (s *CountryQualifierTestSuite) TestIPAPICountryQualifier_Allow() {
	t := s.T()
	ctx := context.Background()
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	ip := "some-ip"
	httpmock.RegisterResponder(http.MethodGet, fmt.Sprintf(urlPattern, ip),
		httpmock.NewStringResponder(http.StatusNotFound, s.qualifier.WhiteList[0]))
	isAllow := s.qualifier.QualifyCountry(ctx, ip)
	assert.True(t, isAllow)
}
