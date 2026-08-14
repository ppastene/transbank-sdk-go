package internal

import (
	"net/http"
	"time"

	"github.com/ppastene/transbank-sdk-go"
)

type Credentials struct {
	ApiKey       string
	CommerceCode string
}

type Config struct {
	Credentials Credentials
	HTTP        transbank.HTTPClient
	BaseURL     string
	Headers     map[string]string
}

func NewConfig(commerceCode, apiKey, defaultBaseURL string) Config {
	return Config{
		Credentials: Credentials{
			ApiKey:       apiKey,
			CommerceCode: commerceCode,
		},
		HTTP:    &http.Client{Timeout: 30 * time.Second},
		BaseURL: defaultBaseURL,
	}
}

func (cfg *Config) SetBaseURL(url string) {
	cfg.BaseURL = url
}
