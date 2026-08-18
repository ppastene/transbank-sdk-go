package transbank

import (
	"net/http"
)

type Environment int

const (
	Integration Environment = iota
	Production
)

// HTTPClient executes an HTTP request. Both *http.Client and any custom
// implementation with a Do method satisfy it, so users can inject their own
// client (for example the underlying client of a Resty client) through
// Options.HTTPClient.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Options struct {
	Environment    Environment
	CommerceCode   string
	ApiKey         string
	HTTPClient     HTTPClient
	ValidateInputs bool
}

func (o Options) Validate() error {
	if o.Environment != Integration && o.Environment != Production {
		return &ValidationError{Message: "invalid environment"}
	}
	if len(o.CommerceCode) == 0 || !isDigits(o.CommerceCode) {
		return &ValidationError{Message: "commerce code must be numeric"}
	}
	if o.ApiKey == "" {
		return &ValidationError{Message: "api key must not be empty"}
	}
	return nil
}

func isDigits(s string) bool {
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}
