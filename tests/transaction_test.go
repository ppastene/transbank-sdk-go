package tests

import (
	"errors"
	"strings"
	"testing"

	webpay "github.com/ppastene/transbank-sdk-go"
	"github.com/ppastene/transbank-sdk-go/internal/shared"
)

var transaction *webpay.Transaction

var options = &webpay.Options{
	ApiKey:       "api-key",
	CommerceCode: "commerce-code",
}

func TestTransactionCreate_InputError(t *testing.T) {
	tests := []struct {
		name          string
		buyOrder      string
		sessionId     string
		amount        float64
		returnUrl     string
		expectedError string
	}{
		// BuyOrder errors
		{
			name:          "BuyOrder exceedes maximum length",
			buyOrder:      strings.Repeat("a", 27),
			sessionId:     "S1",
			amount:        10000,
			returnUrl:     "http://test.com",
			expectedError: "SDK Validation Error: 'buyOrder' is too long, the maximum length is 26",
		},
		{
			name:          "BuyOrder is empty",
			buyOrder:      "",
			sessionId:     "S1",
			amount:        10000,
			returnUrl:     "http://test.com",
			expectedError: "SDK Validation Error: 'buyOrder' cannot be empty",
		},
		// SessionId errors
		{
			name:          "SessionId exceedes maximum length",
			buyOrder:      "OC123",
			sessionId:     strings.Repeat("a", 62),
			amount:        10000,
			returnUrl:     "http://test.com",
			expectedError: "SDK Validation Error: 'sessionId' is too long, the maximum length is 61",
		},
		{
			name:          "SessionId is empty",
			buyOrder:      "OC123",
			sessionId:     "",
			amount:        10000,
			returnUrl:     "http://test.com",
			expectedError: "SDK Validation Error: 'sessionId' cannot be empty",
		},
		// ReturnUrl errors
		{
			name:          "ReturnUrl is empty",
			buyOrder:      "OC123",
			sessionId:     "S1",
			amount:        10000,
			returnUrl:     "",
			expectedError: "SDK Validation Error: 'returnUrl' cannot be empty",
		},
		{
			name:          "ReturnUrl exceedes maximum length",
			buyOrder:      "OC123",
			sessionId:     "S1",
			amount:        10000,
			returnUrl:     "https://example.com/?" + strings.Repeat("a", 236),
			expectedError: "SDK Validation Error: 'returnUrl' is too long, the maximum length is 256",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := transaction.Create(tt.buyOrder, tt.sessionId, tt.amount, tt.returnUrl)
			if err == nil {
				t.Errorf("Se esperaba error pero se obtuvo nil")
				return
			}
			werr, ok := err.(*shared.WebpayError)
			if !ok {
				t.Fatalf("El error devuelto no es de tipo *shared.WebpayError, se obtuvo: %T", err)
			}

			if werr.Error() != tt.expectedError {
				t.Errorf("Mensaje incorrecto.\nEsperado: %q\nObtenido: %q", tt.expectedError, werr.Error())
			}

			if werr.Code != -1 {
				t.Errorf("Código de error incorrecto. Esperado: -1, Obtenido: %d", werr.Code)
			}
		})
	}
}

func TestTransactionCreate_ServerError(t *testing.T) {
	ms := NewMockServer()
	defer ms.Close()
	mockClient := &mockClient{ms.Server.URL}
	tx := webpay.NewTransactionWithClient(mockClient, options)
	tests := []struct {
		name           string
		amount         int
		mockStatusCode int
		mockResponse   map[string]string
		expectedError  string
	}{
		{
			name:           "Amount rejected by Transbank",
			amount:         -1,
			mockStatusCode: 400,
			mockResponse:   map[string]string{"error_message": "Invalid value for parameter: amount"},
			expectedError:  "Invalid value for parameter: amount",
		},
		{
			name:           "Credenciales Inválidas",
			amount:         1000,
			mockStatusCode: 401,
			mockResponse:   map[string]string{"error_message": "not authorized"},
			expectedError:  "not authorized",
		},
		{
			name:           "Error Interno de Transbank",
			amount:         5000,
			mockStatusCode: 500,
			mockResponse:   map[string]string{"error_message": "internal server error"},
			expectedError:  "internal server error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms.StatusCode = tt.mockStatusCode
			ms.Response = tt.mockResponse

			_, err := tx.Create("order", "session", float64(tt.amount), "http://return.cl")

			if err == nil {
				t.Fatal("Se esperaba un error del servidor, pero err fue nil")
			}

			cause := errors.Unwrap(err)

			if cause == nil {
				t.Fatalf("El error no tiene una causa envuelta. Error obtenido: %v", err)
			}

			if cause.Error() != tt.expectedError {
				t.Errorf("Causa del error incorrecta.\nEsperado: %q\nObtenido: %q\n(Mensaje completo: %v)",
					tt.expectedError, cause.Error(), err)
			}

			var werr *shared.WebpayError
			if errors.As(err, &werr) {
				if werr.Code != tt.mockStatusCode {
					t.Errorf("Código HTTP incorrecto. Esperado %d, obtuve %d", tt.mockStatusCode, werr.Code)
				}
			}
		})
	}
}

func TestTransactionCreate_Success(t *testing.T) {
	ms := NewMockServer()
	defer ms.Server.Close()

	ms.Response = map[string]string{
		"token": "webpay_token_123456",
		"url":   "https://webpay.cl/formulario-pago",
	}
	ms.StatusCode = 200

	mockClient := &mockClient{ms.Server.URL}
	tx := webpay.NewTransactionWithClient(mockClient, options)

	res, err := tx.Create("orden123", "sesion456", 15000, "https://mi-sitio.cl/return")

	if err != nil {
		t.Fatalf("No se esperaba error, se obtuvo: %v", err)
	}

	if res.Token != "webpay_token_123456" {
		t.Errorf("Token incorrecto. Esperaba %q, obtuve %q", "webpay_token_123456", res.Token)
	}
}

func TestTransactionStatus_InputError(t *testing.T) {
	tests := []struct {
		name          string
		token         string
		expectedError string
	}{
		{
			name:          "Token exceedes maximum length",
			token:         strings.Repeat("a", 65),
			expectedError: "SDK Validation Error: 'token' is too long, the maximum length is 64",
		},
		{
			name:          "Token is empty",
			token:         "",
			expectedError: "SDK Validation Error: 'token' cannot be empty",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := transaction.Status(tt.token)
			if err == nil {
				t.Errorf("Se esperaba error pero se obtuvo nil")
				return
			}
			werr, ok := err.(*shared.WebpayError)
			if !ok {
				t.Fatalf("El error devuelto no es de tipo *shared.WebpayError, se obtuvo: %T", err)
			}

			if werr.Error() != tt.expectedError {
				t.Errorf("Mensaje incorrecto.\nEsperado: %q\nObtenido: %q", tt.expectedError, werr.Error())
			}

			if werr.Code != -1 {
				t.Errorf("Código de error incorrecto. Esperado: -1, Obtenido: %d", werr.Code)
			}
		})
	}
}

func TestTransactionStatus_ServerError(t *testing.T) {
	ms := NewMockServer()
	defer ms.Close()
	mockClient := &mockClient{ms.Server.URL}
	tx := webpay.NewTransactionWithClient(mockClient, options)

	tests := []struct {
		name           string
		token          string
		mockStatusCode int
		mockResponse   map[string]string
		expectedError  string
	}{
		{
			name:           "Token doesn't exists on the system",
			token:          strings.Repeat("a", 10),
			mockStatusCode: 422,
			mockResponse:   map[string]string{"error_message": "Invalid value for parameter: token"},
			expectedError:  "Invalid value for parameter: token",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms.StatusCode = tt.mockStatusCode
			ms.Response = tt.mockResponse

			_, err := tx.Status(tt.token)

			if err == nil {
				t.Fatal("Se esperaba un error del servidor, pero err fue nil")
			}

			cause := errors.Unwrap(err)

			if cause == nil {
				t.Fatalf("El error no tiene una causa envuelta. Error obtenido: %v", err)
			}

			if cause.Error() != tt.expectedError {
				t.Errorf("Causa del error incorrecta.\nEsperado: %q\nObtenido: %q\n(Mensaje completo: %v)",
					tt.expectedError, cause.Error(), err)
			}

			var werr *shared.WebpayError
			if errors.As(err, &werr) {
				if werr.Code != tt.mockStatusCode {
					t.Errorf("Código HTTP incorrecto. Esperado %d, obtuve %d", tt.mockStatusCode, werr.Code)
				}
			}
		})
	}
}
func TestTransactionStatus_Init(t *testing.T) {
	ms := NewMockServer()
	defer ms.Server.Close()

	ms.Response = map[string]any{
		"amount":              10000,
		"status":              "INITIALIZED",
		"buy_order":           "123456",
		"session_id":          "session123456",
		"accounting_date":     "0223",
		"transaction_date":    "2026-02-23T23:10:58.559Z",
		"installments_number": 0,
	}
	ms.StatusCode = 200

	mockClient := &mockClient{ms.Server.URL}
	tx := webpay.NewTransactionWithClient(mockClient, options)
	res, err := tx.Status(strings.Repeat("a", 64))
	if err != nil {
		t.Fatalf("Error al obtener status: %v", err)
	}
	if res.Status != "INITIALIZED" {
		t.Errorf("Status incorrecto. Esperaba INITIALIZED, obtuve %s", res.Status)
	}
	if res.Amount != 10000 {
		t.Errorf("Monto incorrecto. Esperaba 10000, obtuve %f", res.Amount)
	}

	if res.Vci != "" {
		t.Errorf("VCI debería estar vacío en estado INITIALIZED, obtuve %s", res.Vci)
	}
	if res.AuthorizationCode != "" {
		t.Errorf("AuthorizationCode debería estar vacío, obtuve %s", res.AuthorizationCode)
	}

	if res.CardDetail.CardNumber != "" {
		t.Errorf("CardNumber debería estar vacío, obtuve %s", res.CardDetail.CardNumber)
	}
}

func TestTransactionStatus_Failed(t *testing.T) {
	ms := NewMockServer()
	defer ms.Server.Close()

	ms.Response = map[string]any{
		"vci":        "TSN",
		"amount":     10000,
		"status":     "FAILED",
		"buy_order":  "123456",
		"session_id": "session123456",
		"card_detail": map[string]string{
			"card_number": "6623",
		},
		"accounting_date":     "0223",
		"transaction_date":    "2026-02-24T00:32:34.673Z",
		"authorization_code":  "000000",
		"payment_type_code":   "VN",
		"response_code":       -1,
		"installments_number": 0,
	}
	ms.StatusCode = 200

	mockClient := &mockClient{ms.Server.URL}
	tx := webpay.NewTransactionWithClient(mockClient, options)

	res, err := tx.Status(strings.Repeat("a", 64))

	if err != nil {
		t.Fatalf("No se esperaba error de red, pero se obtuvo: %v", err)
	}

	if res.Status != "FAILED" {
		t.Errorf("Status incorrecto. Esperaba FAILED, obtuve %s", res.Status)
	}

	if res.ResponseCode != -1 {
		t.Errorf("ResponseCode incorrecto. Esperaba -1, obtuve %d", res.ResponseCode)
	}

	if res.Vci != "TSN" {
		t.Errorf("VCI incorrecto. Esperaba TSN, obtuve %s", res.Vci)
	}

	if res.CardDetail.CardNumber != "6623" {
		t.Errorf("CardNumber incorrecto. Esperaba 6623, obtuve %s", res.CardDetail.CardNumber)
	}

	if res.AuthorizationCode != "000000" {
		t.Errorf("AuthorizationCode incorrecto. Esperaba 000000, obtuve %s", res.AuthorizationCode)
	}
}

func TestTransactionStatus_Reversed(t *testing.T) {
	ms := NewMockServer()
	defer ms.Server.Close()

	ms.Response = map[string]any{
		"vci":        "TSY",
		"amount":     10000,
		"status":     "REVERSED",
		"buy_order":  "123456",
		"session_id": "session123456",
		"card_detail": map[string]string{
			"card_number": "6623",
		},
		"accounting_date":     "0223",
		"transaction_date":    "2026-02-24T01:25:48.474Z",
		"authorization_code":  "1213",
		"payment_type_code":   "VN",
		"response_code":       0,
		"installments_number": 0,
	}
	ms.StatusCode = 200

	mockClient := &mockClient{ms.Server.URL}
	tx := webpay.NewTransactionWithClient(mockClient, options)

	res, err := tx.Status(strings.Repeat("a", 64))

	if err != nil {
		t.Fatalf("No se esperaba error de red, pero se obtuvo: %v", err)
	}

	if res.Status != "REVERSED" {
		t.Errorf("Status incorrecto. Esperaba REVERSED, obtuve %s", res.Status)
	}

	if res.ResponseCode != 0 {
		t.Errorf("ResponseCode incorrecto. Esperaba 0, obtuve %d", res.ResponseCode)
	}

	if res.Vci != "TSY" {
		t.Errorf("VCI incorrecto. Esperaba TSY, obtuve %s", res.Vci)
	}

	if res.CardDetail.CardNumber != "6623" {
		t.Errorf("CardNumber incorrecto. Esperaba 6623, obtuve %s", res.CardDetail.CardNumber)
	}

	if res.AuthorizationCode != "1213" {
		t.Errorf("AuthorizationCode incorrecto. Esperaba 1213, obtuve %s", res.AuthorizationCode)
	}
}

func TestTransactionStatus_Success(t *testing.T) {
	ms := NewMockServer()
	defer ms.Server.Close()

	ms.Response = map[string]any{
		"vci":        "TSY",
		"amount":     10000,
		"status":     "AUTHORIZED",
		"buy_order":  "123456",
		"session_id": "session123456",
		"card_detail": map[string]string{
			"card_number": "6623",
		},
		"accounting_date":     "0223",
		"transaction_date":    "2026-02-23T23:10:58.559Z",
		"authorization_code":  "1213",
		"payment_type_code":   "VN",
		"response_code":       0,
		"installments_number": 0,
	}
	ms.StatusCode = 200

	mockClient := &mockClient{ms.Server.URL}
	tx := webpay.NewTransactionWithClient(mockClient, options)

	res, err := tx.Status(strings.Repeat("a", 64))

	if err != nil {
		t.Fatalf("No se esperaba error de red, pero se obtuvo: %v", err)
	}

	if res.Status != "AUTHORIZED" {
		t.Errorf("Status incorrecto. Esperaba AUTHORIZED, obtuve %s", res.Status)
	}

	if res.ResponseCode != 0 {
		t.Errorf("ResponseCode incorrecto. Esperaba 0, obtuve %d", res.ResponseCode)
	}

	if res.Vci != "TSY" {
		t.Errorf("VCI incorrecto. Esperaba TSY, obtuve %s", res.Vci)
	}

	if res.CardDetail.CardNumber != "6623" {
		t.Errorf("CardNumber incorrecto. Esperaba 6623, obtuve %s", res.CardDetail.CardNumber)
	}

	if res.AuthorizationCode != "1213" {
		t.Errorf("AuthorizationCode incorrecto. Esperaba 1213, obtuve %s", res.AuthorizationCode)
	}
}

func TestTRansactionRefund_InputError(t *testing.T) {
	tests := []struct {
		name          string
		token         string
		amount        float64
		expectedError string
	}{
		{
			name:          "Token is empty",
			token:         "",
			amount:        10000,
			expectedError: "SDK Validation Error: 'token' cannot be empty",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := transaction.Refund(tt.token, tt.amount)
			if err == nil {
				t.Errorf("Se esperaba error pero se obtuvo nil")
				return
			}
			werr, ok := err.(*shared.WebpayError)
			if !ok {
				t.Fatalf("El error devuelto no es de tipo *shared.WebpayError, se obtuvo: %T", err)
			}

			if werr.Error() != tt.expectedError {
				t.Errorf("Mensaje incorrecto.\nEsperado: %q\nObtenido: %q", tt.expectedError, werr.Error())
			}

			if werr.Code != -1 {
				t.Errorf("Código de error incorrecto. Esperado: -1, Obtenido: %d", werr.Code)
			}
		})
	}
}

func TestTRansactionRefund_ServerError(t *testing.T) {
	ms := NewMockServer()
	defer ms.Close()
	mockClient := &mockClient{ms.Server.URL}
	tx := webpay.NewTransactionWithClient(mockClient, options)
	tests := []struct {
		name           string
		token          string
		amount         float64
		mockStatusCode int
		mockResponse   map[string]string
		expectedError  string
	}{
		{
			name:           "Monto no aceptado por Transbank",
			token:          strings.Repeat("a", 64),
			amount:         0,
			mockStatusCode: 422,
			mockResponse:   map[string]string{"error_message": "Invalid value for parameter: amount"},
			expectedError:  "Invalid value for parameter: amount",
		},
		{
			name:           "Token es invalido",
			token:          strings.Repeat("a", 10),
			amount:         10000,
			mockStatusCode: 422,
			mockResponse:   map[string]string{"error_message": "Invalid value for parameter: token"},
			expectedError:  "Invalid value for parameter: token",
		},
		{
			name:           "Monto es mayor al pagado",
			token:          strings.Repeat("a", 64),
			amount:         100000,
			mockStatusCode: 422,
			mockResponse:   map[string]string{"error_message": "Amount to refund is bigger than authorized"},
			expectedError:  "Amount to refund is bigger than authorized",
		},
		{
			name:           "Se intenta hacer reembolso a una transaccion rechazada",
			token:          strings.Repeat("a", 64),
			amount:         10000,
			mockStatusCode: 422,
			mockResponse:   map[string]string{"error_message": "Transaction has an invalid state (8) to reverse"},
			expectedError:  "Transaction has an invalid state (8) to reverse",
		},
		{
			name:           "Monto a reembolsar no es exacto",
			token:          strings.Repeat("a", 64),
			amount:         9999,
			mockStatusCode: 400,
			mockResponse:   map[string]string{"error_message": "Transaction not found"},
			expectedError:  "Transaction not found",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms.StatusCode = tt.mockStatusCode
			ms.Response = tt.mockResponse

			_, err := tx.Refund(tt.token, tt.amount)

			if err == nil {
				t.Fatal("Se esperaba un error del servidor, pero err fue nil")
			}

			cause := errors.Unwrap(err)

			if cause == nil {
				t.Fatalf("El error no tiene una causa envuelta. Error obtenido: %v", err)
			}

			if cause.Error() != tt.expectedError {
				t.Errorf("Causa del error incorrecta.\nEsperado: %q\nObtenido: %q\n(Mensaje completo: %v)",
					tt.expectedError, cause.Error(), err)
			}

			var werr *shared.WebpayError
			if errors.As(err, &werr) {
				if werr.Code != tt.mockStatusCode {
					t.Errorf("Código HTTP incorrecto. Esperado %d, obtuve %d", tt.mockStatusCode, werr.Code)
				}
			}
		})
	}
}

func TestTRansactionRefund_ReverseSuccess(t *testing.T) {
	ms := NewMockServer()
	defer ms.Close()
	mockClient := &mockClient{ms.Server.URL}
	tx := webpay.NewTransactionWithClient(mockClient, options)

	ms.Response = map[string]any{
		"type": "REVERSED",
	}
	ms.StatusCode = 200

	res, err := tx.Refund(strings.Repeat("a", 64), 10000)

	if err != nil {
		t.Fatalf("No se esperaba error de red, pero se obtuvo: %v", err)
	}

	if res.Type != "REVERSED" {
		t.Errorf("Tipo incorrecto. Esperaba REVERSED, obtuve %s", res.Type)
	}
}

func TestTRansactionRefund_NullifiedSuccess(t *testing.T) {
	ms := NewMockServer()
	defer ms.Close()
	mockClient := &mockClient{ms.Server.URL}
	tx := webpay.NewTransactionWithClient(mockClient, options)

	ms.Response = map[string]any{
		"type":               "NULLIFIED",
		"authorization_code": "123456",
		"authorization_date": "2019-03-20T20:18:20Z",
		"nullified_amount":   1000.00,
		"balance":            0.00,
		"response_code":      0,
	}
	ms.StatusCode = 200

	res, err := tx.Refund(strings.Repeat("a", 64), 1000)

	if err != nil {
		t.Fatalf("No se esperaba error de red, pero se obtuvo: %v", err)
	}

	if res.Type != "NULLIFIED" {
		t.Errorf("Tipo incorrecto. Esperaba NULLIFIED, obtuve %s", res.Type)
	}

	if res.AuthorizationCode == "" {
		t.Errorf("Authorization Code no puede estar vacio. Esperaba 123456, obtuve %s", res.Type)
	}

	if res.AuthorizationDate == "" {
		t.Errorf("Authorization Date no puede estar vacio. Esperaba 2019-03-20T20:18:20Z, obtuve %s", res.Type)
	}

	if res.ResponseCode != 0 {
		t.Errorf("Response Code Incorecto. Esperaba 0, obtuve %s", res.Type)
	}
}
