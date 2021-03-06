package api

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/Sugar-pack/test-task/internal/logging"
)

type CountryQualifier interface {
	QualifyCountry(ctx context.Context, ip string) bool
}

type IPAPICountryQualifier struct {
	WhiteList []string
}

func NewIPAPI(whiteList []string) *IPAPICountryQualifier {
	return &IPAPICountryQualifier{WhiteList: whiteList}
}

const urlPattern = "https://ipapi.co/%s/country_name/"

func (i *IPAPICountryQualifier) QualifyCountry(ctx context.Context, ip string) bool {
	logger := logging.FromContext(ctx)
	ipapiClient := http.Client{}
	url := fmt.Sprintf(urlPattern, ip)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logger.WithError(err).Error("cant create request")

		return false
	}
	req.Header.Set("User-Agent", "ipapi.co/#go-v1.5")
	resp, err := ipapiClient.Do(req)
	if err != nil {
		logger.WithError(err).Error("cant get response")

		return false
	}
	defer func(Body io.ReadCloser) {
		errClose := Body.Close()
		if errClose != nil {
			logger.WithError(errClose).Error("cant close body")
		}
	}(resp.Body)
	country, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.WithError(err).Error("cant get country")

		return false
	}
	logger.Info("country ", string(country))

	for _, countryName := range i.WhiteList {
		if string(country) == countryName {
			return true
		}
	}

	return false
}
