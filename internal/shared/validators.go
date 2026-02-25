package shared

import (
	"fmt"
	"math"
	"net/url"
	"regexp"
	"strings"
)

var validTbkRegex = regexp.MustCompile(`^[a-zA-Z0-9|_=&%.,~:/?[+!@()>\-. ]*$`)

func IsTextEmpty(value, fieldName string) error {
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("'%s' cannot be empty", fieldName)
	}
	return nil
}

func HasTextWithMaxLength(value string, length int, fieldName string) error {
	if err := IsTextEmpty(value, fieldName); err != nil {
		return err
	}
	if len(value) > length {
		return fmt.Errorf("'%s' is too long, the maximum length is %d", fieldName, length)
	}
	return nil
}

func HasInvalidCharacters(value, fieldName string) error {
	if !validTbkRegex.MatchString(value) {
		return fmt.Errorf("%s contains invalid characters", fieldName)
	}
	return nil
}

func IsValueNumeric(value any, fieldName string) error {
	switch value.(type) {
	case int, int32, int64, float32, float64:
		return nil
	default:
		return fmt.Errorf("%s must be a numeric value (int or float), got %T", fieldName, value)
	}
}

func IsValueGreaterThanZero(value any, fieldName string) error {
	var floatValue float64
	switch v := value.(type) {
	case int:
		floatValue = float64(v)
	case int32:
		floatValue = float64(v)
	case int64:
		floatValue = float64(v)
	case float32:
		floatValue = float64(v)
	case float64:
		floatValue = v
	default:
		return IsValueNumeric(value, fieldName)
	}

	if floatValue <= 0 {
		return fmt.Errorf("%s must be greater than zero", fieldName)
	}

	return nil
}

func HasValidDecimalPlaces(value any, maxDecimal int, fieldName string) error {
	var floatValue float64
	switch v := value.(type) {
	case float64:
		floatValue = v
	case float32:
		floatValue = float64(v)
	case int, int32, int64:
		return nil
	default:
		return IsValueNumeric(value, fieldName)
	}

	shift := math.Pow(10, float64(maxDecimal))
	shiftedValue := floatValue * shift

	if math.Abs(shiftedValue-math.Trunc(shiftedValue)) > 0.0000001 {
		return fmt.Errorf("%s cannot have more than %d decimal places", fieldName, maxDecimal)
	}

	return nil
}

func IsValidURL(rawURL, fieldName string) error {
	u, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return fmt.Errorf("%s is not a valid URL", fieldName)
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("%s must be an absolute URL including protocol (e.g., https://example.com)", fieldName)
	}

	if u.Host == "" {
		return fmt.Errorf("%s must have a host", fieldName)
	}

	return nil
}
