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
	QualifyCountry(ctx context.Context, ip string) (string, error)
}

type IPAPICountryQualifier struct{}

func (i *IPAPICountryQualifier) QualifyCountry(ctx context.Context, ip string) (string, error) {
	logger := logging.FromContext(ctx)
	ipapiClient := http.Client{}
	url := fmt.Sprintf("https://ipapi.co/%s/country_name/", ip)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logger.WithError(err).Error("cant create request")
		return "", err
	}
	req.Header.Set("User-Agent", "ipapi.co/#go-v1.5")
	resp, err := ipapiClient.Do(req)
	if err != nil {
		logger.WithError(err).Error("cant get response")
		return "", err
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
		return "", err
	}
	logger.Info("country ", string(country))
	return string(country), nil
}
