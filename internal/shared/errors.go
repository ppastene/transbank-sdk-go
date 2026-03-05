package shared

import (
	"encoding/json"
	"fmt"
)

type WebpayError struct {
	Code           int    `json:"code"`
	ServiceMessage string `json:"service_message"`
	Cause          error  `json:"cause"`
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

func (e *WebpayError) MarshalJSON() ([]byte, error) {
	type Alias WebpayError

	var causeMsg string
	if e.Cause != nil {
		causeMsg = e.Cause.Error()
	}

	return json.Marshal(&struct {
		*Alias
		Cause string `json:"cause"`
	}{
		Alias: (*Alias)(e),
		Cause: causeMsg,
	})
}
