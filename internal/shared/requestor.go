package shared

import (
	"encoding/json"
	"errors"
)

type HTTPClientInterface interface {
	Request(method string, url string, headers map[string]string, payload any) ([]byte, int, error)
}

type Requestor struct {
	Client  HTTPClientInterface
	Options *Options
}

func (r *Requestor) Do(method, path string, payload any, result any) (int, error) {
	if r.Options == nil || (r.Options.ApiKey == "" && r.Options.CommerceCode == "") {
		return -1, &WebpayError{Code: -1, ServiceMessage: "Error while requesting to Transbank:", Cause: errors.New("Missing credentials")}
	}

	fullUrl := r.Options.GetBaseUrl() + path

	headers := map[string]string{
		"Tbk-Api-Key-Id":     r.Options.ApiKey,
		"Tbk-Api-Key-Secret": r.Options.CommerceCode,
		"Content-Type":       "application/json",
	}

	body, code, err := r.Client.Request(method, fullUrl, headers, payload)
	if err != nil {
		return 0, &WebpayError{
			Code:           0,
			ServiceMessage: "Connection refused or timeout",
			Cause:          err,
		}
	}

	if code < 200 || code >= 300 {
		var apiErr struct {
			ErrorMessage string `json:"error_message"`
		}
		_ = json.Unmarshal(body, &apiErr)

		msg := apiErr.ErrorMessage
		if msg == "" {
			msg = string(body)
		}

		return code, &WebpayError{
			Code:           code,
			ServiceMessage: "Transbank API error",
			Cause:          errors.New(msg),
		}
	}

	if result != nil && len(body) > 0 {
		if err := json.Unmarshal(body, result); err != nil {
			return code, &WebpayError{
				Code:           code,
				ServiceMessage: "Error unmarshaling success response",
				Cause:          err,
			}
		}
	}

	return code, nil
}
