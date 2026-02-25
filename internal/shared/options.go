package shared

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

func (o *Options) GetBaseUrl() string {
	switch o.Environment {
	case Live:
		return BaseUrlProduction
	default:
		return BaseUrlIntegration
	}
}
