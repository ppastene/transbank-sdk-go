package internal

import (
	"math"
	"net/url"
	"regexp"
	"strings"

	"github.com/ppastene/transbank-sdk-go"
)

var buyOrderPattern = regexp.MustCompile(`^[A-Za-z0-9_.\-%~.,+!@()=]+$`)

var cardExpirationPattern = regexp.MustCompile(`^\d{2}/\d{2}$`)

func ValidateCommerceCode(commerceCode string) error {
	if len(commerceCode) != 12 || !isNumeric(commerceCode) {
		return &transbank.ValidationError{Message: "commerce code must be 12 digits"}
	}
	return nil
}

func isNumeric(s string) bool {
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

func ValidateToken(token string) error {
	if len(token) != 64 {
		return &transbank.ValidationError{Message: "token must be 64 characters"}
	}
	return nil
}

func ValidateBuyOrder(buyOrder string) error {
	if buyOrder == "" {
		return &transbank.ValidationError{Message: "buy_order must not be empty"}
	}
	if len(buyOrder) > 26 {
		return &transbank.ValidationError{Message: "buy_order must be at most 26 characters"}
	}
	if !buyOrderPattern.MatchString(buyOrder) {
		return &transbank.ValidationError{Message: "buy_order contains invalid characters"}
	}
	return nil
}

func ValidateSessionID(sessionID string) error {
	if sessionID == "" {
		return &transbank.ValidationError{Message: "session_id must not be empty"}
	}
	if len(sessionID) > 61 {
		return &transbank.ValidationError{Message: "session_id must be at most 61 characters"}
	}
	return nil
}

func ValidateAmount(amount float64) error {
	if amount <= 0 {
		return &transbank.ValidationError{Message: "amount must be greater than zero"}
	}
	if math.Abs(amount*100-math.Round(amount*100)) > 1e-9 {
		return &transbank.ValidationError{Message: "amount must have at most 2 decimal places"}
	}
	return nil
}

func ValidateUsername(username string) error {
	if username == "" {
		return &transbank.ValidationError{Message: "username must not be empty"}
	}
	if len(username) > 40 {
		return &transbank.ValidationError{Message: "username must be at most 40 characters"}
	}
	return nil
}

func ValidateEmail(email string) error {
	if email == "" {
		return &transbank.ValidationError{Message: "email must not be empty"}
	}
	if len(email) > 100 {
		return &transbank.ValidationError{Message: "email must be at most 100 characters"}
	}
	if !strings.Contains(email, "@") {
		return &transbank.ValidationError{Message: "email must contain @"}
	}
	return nil
}

func ValidateTbkUser(tbkUser string) error {
	if tbkUser == "" {
		return &transbank.ValidationError{Message: "tbk_user must not be empty"}
	}
	if len(tbkUser) > 40 {
		return &transbank.ValidationError{Message: "tbk_user must be at most 40 characters"}
	}
	return nil
}

func ValidateReturnURL(returnURL string) error {
	if len(returnURL) > 255 {
		return &transbank.ValidationError{Message: "return_url must be at most 255 characters"}
	}
	u, err := url.Parse(returnURL)
	if err != nil || !u.IsAbs() || u.Host == "" {
		return &transbank.ValidationError{Message: "return_url must be an absolute URL"}
	}
	scheme := strings.ToLower(u.Scheme)
	if scheme != "http" && scheme != "https" {
		return &transbank.ValidationError{Message: "return_url must use http or https"}
	}
	return nil
}

func ValidateCardNumber(cardNumber string) error {
	if cardNumber == "" {
		return &transbank.ValidationError{Message: "card_number must not be empty"}
	}
	if len(cardNumber) > 16 {
		return &transbank.ValidationError{Message: "card_number must be at most 16 digits"}
	}
	if !isNumeric(cardNumber) {
		return &transbank.ValidationError{Message: "card_number must be numeric"}
	}
	return nil
}

func ValidateCardExpirationDate(cardExpirationDate string) error {
	if cardExpirationDate == "" {
		return &transbank.ValidationError{Message: "card_expiration_date must not be empty"}
	}
	if len(cardExpirationDate) > 5 {
		return &transbank.ValidationError{Message: "card_expiration_date must be at most 5 characters"}
	}
	if !cardExpirationPattern.MatchString(cardExpirationDate) {
		return &transbank.ValidationError{Message: "card_expiration_date must be in MM/YY format"}
	}
	return nil
}

func ValidateCVV(cvv string) error {
	if cvv == "" {
		return nil
	}
	if len(cvv) > 4 {
		return &transbank.ValidationError{Message: "cvv must be at most 4 digits"}
	}
	if !isNumeric(cvv) {
		return &transbank.ValidationError{Message: "cvv must be numeric"}
	}
	return nil
}

func ValidateInstallmentsNumber(installmentsNumber int) error {
	if installmentsNumber < 1 || installmentsNumber > 99 {
		return &transbank.ValidationError{Message: "installments_number must be between 1 and 99"}
	}
	return nil
}
