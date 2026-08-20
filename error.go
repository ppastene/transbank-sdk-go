package transbank

import (
	"fmt"
	"strings"
)

// ValidationError indicates invalid SDK inputs or credentials. The API is not
// called.
type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

// TransportError indicates the request or response could not be completed
// (network, encoding, parsing). Err is the root cause.
type TransportError struct {
	Message string
	Err     error
}

func (e *TransportError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *TransportError) Unwrap() error {
	return e.Err
}

// HTTPError indicates the Transbank API responded with a non-2xx status.
type HTTPError struct {
	StatusCode int
	Body       string
}

func (e *HTTPError) Error() string {
	if e.Body == "" {
		return fmt.Sprintf("transbank: status %d", e.StatusCode)
	}
	return fmt.Sprintf("transbank: status %d: %s", e.StatusCode, strings.TrimSpace(e.Body))
}
