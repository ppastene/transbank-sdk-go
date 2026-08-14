package patpass

import (
	"github.com/ppastene/transbank-sdk-go"
)

// Environment identifies a PatPass environment.
type Environment int

const (
	// Integration is the PatPass integration (sandbox) environment.
	Integration Environment = iota
	// Production is the PatPass production environment.
	Production
)

// Options configures a PatPass service. It requires an Environment, a numeric
// CommerceCode and an Authorization key, and optionally an HTTPClient.
type Options struct {
	Environment   Environment          // integration or production
	CommerceCode  string               // numeric commerce code
	Authorization string               // authorization key
	HTTPClient    transbank.HTTPClient // optional custom HTTP client
}

// Validate returns an error of type *transbank.ValidationError if the options
// are invalid: unknown environment, non-numeric commerce code or empty
// authorization.
func (o Options) Validate() error {
	if o.Environment != Integration && o.Environment != Production {
		return &transbank.ValidationError{Message: "invalid environment"}
	}
	if len(o.CommerceCode) == 0 || !isDigits(o.CommerceCode) {
		return &transbank.ValidationError{Message: "commerce code must be numeric"}
	}
	if o.Authorization == "" {
		return &transbank.ValidationError{Message: "authorization must not be empty"}
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
