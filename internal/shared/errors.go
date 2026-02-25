package shared

import "fmt"

type WebpayError struct {
	Code           int    `json:"code"`
	ServiceMessage string `json:"service_message"`
	Cause          error
}

func (e *WebpayError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.ServiceMessage, e.Cause)
	}
	return e.ServiceMessage
}

func (e *WebpayError) Unwrap() error {
	return e.Cause
}
