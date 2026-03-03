package shared

import "errors"

type WebpayIntegrationType string

const (
	Live               WebpayIntegrationType = "LIVE"
	Test               WebpayIntegrationType = "TEST"
	BaseUrlProduction  string                = "https://webpay3g.transbank.cl"
	BaseUrlIntegration string                = "https://webpay3gint.transbank.cl"
)

type Options struct {
	ApiKey       string
	CommerceCode string
	Environment  WebpayIntegrationType
}

func (o *Options) Validate() error {
	if o == nil || (o.ApiKey == "" && o.CommerceCode == "") {
		return errors.New("No credentials")
	} else if o.ApiKey == "" {
		return errors.New("ApiKey is required")
	} else if o.CommerceCode == "" {
		return errors.New("CommerceCode is required")
	}

	return nil
}

func (o *Options) GetBaseUrl() string {
	switch o.Environment {
	case Live:
		return BaseUrlProduction
	default:
		return BaseUrlIntegration
	}
}
