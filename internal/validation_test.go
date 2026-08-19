package internal

import (
	"errors"
	"strings"
	"math"
	"testing"

	"github.com/ppastene/transbank-sdk-go"
)

func validationErr(t *testing.T, err error) *transbank.ValidationError {
	t.Helper()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var tbErr *transbank.ValidationError
	if !errors.As(err, &tbErr) {
		t.Fatalf("error type = %T, want *transbank.ValidationError", err)
	}
	return tbErr
}

func TestValidateCommerceCode(t *testing.T) {
	tests := []struct {
		name         string
		commerceCode string
		wantErr      bool
	}{
		{name: "valid", commerceCode: "597055555532", wantErr: false},
		{name: "empty", commerceCode: "", wantErr: true},
		{name: "11 digits", commerceCode: "59705555553", wantErr: true},
		{name: "13 digits", commerceCode: "5970555555321", wantErr: true},
		{name: "non numeric", commerceCode: "59705555553a", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCommerceCode(tt.commerceCode)
			if tt.wantErr {
				if err := validationErr(t, err); err.Message == "" {
					t.Error("expected non-empty validation message")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestValidateToken(t *testing.T) {
	valid := strings.Repeat("a", 64)
	tests := []struct {
		name    string
		token   string
		wantErr bool
	}{
		{name: "valid", token: valid, wantErr: false},
		{name: "empty", token: "", wantErr: true},
		{name: "63 chars", token: valid[:63], wantErr: true},
		{name: "65 chars", token: valid + "a", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateToken(tt.token)
			if tt.wantErr {
				if err := validationErr(t, err); err.Message == "" {
					t.Error("expected non-empty validation message")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestValidateBuyOrder(t *testing.T) {
	tests := []struct {
		name     string
		buyOrder string
		wantErr  bool
	}{
		{name: "valid", buyOrder: "orden-de-compra-123", wantErr: false},
		{name: "empty", buyOrder: "", wantErr: true},
		{name: "26 chars", buyOrder: strings.Repeat("a", 26), wantErr: false},
		{name: "27 chars", buyOrder: strings.Repeat("a", 27), wantErr: true},
		{name: "invalid characters", buyOrder: "orden con espacios", wantErr: true},
		{name: "accented", buyOrder: "órden-1", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateBuyOrder(tt.buyOrder)
			if tt.wantErr {
				if err := validationErr(t, err); err.Message == "" {
					t.Error("expected non-empty validation message")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestValidateSessionID(t *testing.T) {
	tests := []struct {
		name      string
		sessionID string
		wantErr   bool
	}{
		{name: "valid", sessionID: "sesion-123", wantErr: false},
		{name: "empty", sessionID: "", wantErr: true},
		{name: "61 chars", sessionID: strings.Repeat("a", 61), wantErr: false},
		{name: "62 chars", sessionID: strings.Repeat("a", 62), wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateSessionID(tt.sessionID)
			if tt.wantErr {
				if err := validationErr(t, err); err.Message == "" {
					t.Error("expected non-empty validation message")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestValidateAmount(t *testing.T) {
	tests := []struct {
		name    string
		amount  float64
		wantErr bool
	}{
		{name: "valid", amount: 10000, wantErr: false},
		{name: "valid decimals", amount: 10000.99, wantErr: false},
		{name: "zero", amount: 0, wantErr: true},
		{name: "negative", amount: -100, wantErr: true},
		{name: "three decimals", amount: 10000.123, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateAmount(tt.amount)
			if tt.wantErr {
				if err := validationErr(t, err); err.Message == "" {
					t.Error("expected non-empty validation message")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestValidateUsername(t *testing.T) {
	tests := []struct {
		name     string
		username string
		wantErr  bool
	}{
		{name: "valid", username: "juanperez", wantErr: false},
		{name: "empty", username: "", wantErr: true},
		{name: "40 chars", username: strings.Repeat("a", 40), wantErr: false},
		{name: "41 chars", username: strings.Repeat("a", 41), wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateUsername(tt.username)
			if tt.wantErr {
				if err := validationErr(t, err); err.Message == "" {
					t.Error("expected non-empty validation message")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		wantErr bool
	}{
		{name: "valid", email: "juan.perez@gmail.com", wantErr: false},
		{name: "empty", email: "", wantErr: true},
		{name: "no at", email: "correo-sin-arroba", wantErr: true},
		{name: "100 chars", email: strings.Repeat("a", 92) + "@mail.cl", wantErr: false},
		{name: "101 chars", email: strings.Repeat("a", 93) + "@mail.cl", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateEmail(tt.email)
			if tt.wantErr {
				if err := validationErr(t, err); err.Message == "" {
					t.Error("expected non-empty validation message")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestValidateTbkUser(t *testing.T) {
	tests := []struct {
		name    string
		tbkUser string
		wantErr bool
	}{
		{name: "valid", tbkUser: "b6bd6ba3-e718-4107-9386-d2b099a8dd42", wantErr: false},
		{name: "empty", tbkUser: "", wantErr: true},
		{name: "40 chars", tbkUser: strings.Repeat("a", 40), wantErr: false},
		{name: "41 chars", tbkUser: strings.Repeat("a", 41), wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTbkUser(tt.tbkUser)
			if tt.wantErr {
				if err := validationErr(t, err); err.Message == "" {
					t.Error("expected non-empty validation message")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestValidateReturnURL(t *testing.T) {
	tests := []struct {
		name      string
		returnURL string
		wantErr   bool
	}{
		{name: "valid https", returnURL: "https://www.mi-tienda.cl/retorno", wantErr: false},
		{name: "valid http", returnURL: "http://www.mi-tienda.cl/retorno", wantErr: false},
		{name: "empty", returnURL: "", wantErr: true},
		{name: "relative", returnURL: "/retorno", wantErr: true},
		{name: "no host", returnURL: "https:///retorno", wantErr: true},
		{name: "ftp scheme", returnURL: "ftp://www.mi-tienda.cl/retorno", wantErr: true},
		{name: "255 chars", returnURL: "https://" + strings.Repeat("a", 244) + ".cl", wantErr: false},
		{name: "256 chars", returnURL: "https://" + strings.Repeat("a", 245) + ".cl", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateReturnURL(tt.returnURL)
			if tt.wantErr {
				if err := validationErr(t, err); err.Message == "" {
					t.Error("expected non-empty validation message")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestValidateCardNumber(t *testing.T) {
	tests := []struct {
		name       string
		cardNumber string
		wantErr    bool
	}{
		{name: "valid", cardNumber: "4239000000000000", wantErr: false},
		{name: "empty", cardNumber: "", wantErr: true},
		{name: "16 chars", cardNumber: strings.Repeat("9", 16), wantErr: false},
		{name: "17 chars", cardNumber: strings.Repeat("9", 17), wantErr: true},
		{name: "non numeric", cardNumber: "4239-0000-0000-0000", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCardNumber(tt.cardNumber)
			if tt.wantErr {
				if err := validationErr(t, err); err.Message == "" {
					t.Error("expected non-empty validation message")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestValidateCardExpirationDate(t *testing.T) {
	tests := []struct {
		name               string
		cardExpirationDate string
		wantErr            bool
	}{
		{name: "valid", cardExpirationDate: "22/10", wantErr: false},
		{name: "valid two digits", cardExpirationDate: "01/25", wantErr: false},
		{name: "empty", cardExpirationDate: "", wantErr: true},
		{name: "wrong format", cardExpirationDate: "2210", wantErr: true},
		{name: "bad separator", cardExpirationDate: "22-10", wantErr: true},
		{name: "non numeric month", cardExpirationDate: "ab/10", wantErr: true},
		{name: "too long", cardExpirationDate: "1/2025", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCardExpirationDate(tt.cardExpirationDate)
			if tt.wantErr {
				if err := validationErr(t, err); err.Message == "" {
					t.Error("expected non-empty validation message")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestValidateCVV(t *testing.T) {
	tests := []struct {
		name    string
		cvv     string
		wantErr bool
	}{
		{name: "valid", cvv: "123", wantErr: false},
		{name: "empty is valid", cvv: "", wantErr: false},
		{name: "4 chars", cvv: "1234", wantErr: false},
		{name: "5 chars", cvv: "12345", wantErr: true},
		{name: "non numeric", cvv: "12a", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCVV(tt.cvv)
			if tt.wantErr {
				if err := validationErr(t, err); err.Message == "" {
					t.Error("expected non-empty validation message")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestValidateInstallmentsNumber(t *testing.T) {
	tests := []struct {
		name               string
		installmentsNumber int
		wantErr            bool
	}{
		{name: "valid", installmentsNumber: 10, wantErr: false},
		{name: "one", installmentsNumber: 1, wantErr: false},
		{name: "99", installmentsNumber: 99, wantErr: false},
		{name: "zero", installmentsNumber: 0, wantErr: true},
		{name: "negative", installmentsNumber: -1, wantErr: true},
		{name: "100", installmentsNumber: 100, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateInstallmentsNumber(tt.installmentsNumber)
			if tt.wantErr {
				if err := validationErr(t, err); err.Message == "" {
					t.Error("expected non-empty validation message")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestValidateAmountEdgeCases(t *testing.T) {
	tests := []struct {
		name    string
		amount  float64
		wantErr bool
	}{
		{name: "NaN", amount: math.NaN(), wantErr: true},
		{name: "+Inf", amount: math.Inf(1), wantErr: true},
		{name: "-Inf", amount: math.Inf(-1), wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateAmount(tt.amount)
			if tt.wantErr {
				if err := validationErr(t, err); err.Message == "" {
					t.Error("expected non-empty validation message")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestValidateCardNumberShort(t *testing.T) {
	if err := ValidateCardNumber("4"); err != nil {
		t.Errorf("short numeric card number should be accepted: %v", err)
	}
}
