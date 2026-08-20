package internal

import (
	"regexp"
	"strconv"

	"github.com/ppastene/transbank-sdk-go"
)

var urlParamPattern = regexp.MustCompile(`^[A-Za-z0-9_-]+$`)

func isNumeric(s string) bool {
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

func ValidateURLParam(name, value string, maxLen int) error {
	if value == "" {
		return &transbank.ValidationError{Message: name + " must not be empty"}
	}
	if len(value) > maxLen {
		return &transbank.ValidationError{Message: name + " must be at most " + strconv.Itoa(maxLen) + " characters"}
	}
	if !urlParamPattern.MatchString(value) {
		return &transbank.ValidationError{Message: name + " contains invalid characters"}
	}
	return nil
}
