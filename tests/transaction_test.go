package tests

import (
	"errors"
	"strings"
	"testing"

	"github.com/ppastene/transbank-sdk-go"
	"github.com/ppastene/transbank-sdk-go/internal/shared"
)

var transactionOptions = &transbank.Options{
	ApiKey:       "api-key",
	CommerceCode: "commerce-code",
}

func TestTransactionCreate_ServerError(t *testing.T) {
	ms := NewMockServer()
	defer ms.Close()
	mockClient := &mockClient{ms.Server.URL}
	tx := transbank.NewTransactionWithClient(mockClient, transactionOptions)
	tests := []struct {
		name           string
		buyOrder       string
		sessionId      string
		amount         float64
		returnUrl      string
		mockStatusCode int
		mockResponse   map[string]string
		expectedError  string
	}{
		{
			name:           "'amount' cannot be a negative number",
			buyOrder:       "11223344",
			sessionId:      "S1",
			amount:         -1,
			returnUrl:      "http://www.return-url.com",
			mockStatusCode: 422,
			mockResponse:   map[string]string{"error_message": "Invalid value for parameter: amount"},
			expectedError:  "Invalid value for parameter: amount",
		},
		{
			name:           "'amount' cannot be zero",
			buyOrder:       "11223344",
			sessionId:      "S1",
			amount:         0,
			returnUrl:      "http://www.return-url.com",
			mockStatusCode: 422,
			mockResponse:   map[string]string{"error_message": "Invalid value for parameter: amount"},
			expectedError:  "Invalid value for parameter: amount",
		},
		{
			name:           "'buy_order' is empty",
			buyOrder:       "",
			sessionId:      "S1",
			amount:         1000,
			returnUrl:      "http://www.return-url.com",
			mockStatusCode: 422,
			mockResponse:   map[string]string{"error_message": "buy_order is required!"},
			expectedError:  "buy_order is required!",
		},
		{
			name:           "'buy_order' has invalid characters",
			buyOrder:       "buy-order;",
			sessionId:      "S1",
			amount:         1000,
			returnUrl:      "http://www.return-url.com",
			mockStatusCode: 422,
			mockResponse:   map[string]string{"error_message": "Parameter TBK_ORDEN_COMPRA rejected with value buy-order;"},
			expectedError:  "Parameter TBK_ORDEN_COMPRA rejected with value buy-order;",
		},
		{
			name:           "'buy_order' exceedes max length",
			buyOrder:       strings.Repeat("a", 27),
			sessionId:      "S1",
			amount:         1000,
			returnUrl:      "http://www.return-url.com",
			mockStatusCode: 422,
			mockResponse:   map[string]string{"error_message": "Invalid value for parameter: buy_order"},
			expectedError:  "Invalid value for parameter: buy_order",
		},
		{
			name:           "'session_id' is empty",
			buyOrder:       "11223344",
			sessionId:      "",
			amount:         1000,
			returnUrl:      "http://www.return-url.com",
			mockStatusCode: 422,
			mockResponse:   map[string]string{"error_message": "session_id is required!"},
			expectedError:  "session_id is required!",
		},
		{
			name:           "'session_id' exceedes max length",
			buyOrder:       "11223344",
			sessionId:      strings.Repeat("a", 62),
			amount:         1000,
			returnUrl:      "http://www.return-url.com",
			mockStatusCode: 422,
			mockResponse:   map[string]string{"error_message": "Invalid value for parameter: session_id"},
			expectedError:  "Invalid value for parameter: session_id",
		},
		{
			name:           "'return_url' is empty",
			buyOrder:       "11223344",
			sessionId:      "S1",
			amount:         1000,
			returnUrl:      "",
			mockStatusCode: 400,
			mockResponse:   map[string]string{"error_message": "return_url is required"},
			expectedError:  "return_url is required",
		},
		{
			name:           "'return_url' is not an url",
			buyOrder:       "11223344",
			sessionId:      "S1",
			amount:         1000,
			returnUrl:      "some-random-url",
			mockStatusCode: 422,
			mockResponse:   map[string]string{"error_message": "Invalid value for parameter: return_url"},
			expectedError:  "Invalid value for parameter: return_url",
		},
		{
			name:           "'return_url' is not valid",
			buyOrder:       "11223344",
			sessionId:      "S1",
			amount:         1000,
			returnUrl:      "www.return-url.com",
			mockStatusCode: 422,
			mockResponse:   map[string]string{"error_message": "Invalid value for parameter: return_url"},
			expectedError:  "Invalid value for parameter: return_url",
		},
		{
			name:           "Invalid credentials",
			buyOrder:       "11223344",
			sessionId:      "S1",
			amount:         1000,
			returnUrl:      "http://www.return-url.com",
			mockStatusCode: 401,
			mockResponse:   map[string]string{"error_message": "not authorized"},
			expectedError:  "not authorized",
		},
		{
			name:           "Transbank error",
			buyOrder:       "11223344",
			sessionId:      "S1",
			amount:         5000,
			returnUrl:      "http://www.return-url.com",
			mockStatusCode: 500,
			mockResponse:   map[string]string{"error_message": "internal server error"},
			expectedError:  "internal server error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms.StatusCode = tt.mockStatusCode
			ms.Response = tt.mockResponse

			_, err := tx.Create(tt.buyOrder, tt.sessionId, tt.amount, tt.returnUrl)

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
	tx := transbank.NewTransactionWithClient(mockClient, transactionOptions)

	res, err := tx.Create("orden123", "session456", 15000, "https://mi-sitio.cl/return")

	if err != nil {
		t.Fatalf("No se esperaba error, se obtuvo: %v", err)
	}

	if res.Token != "webpay_token_123456" {
		t.Errorf("Token incorrecto. Esperaba %q, obtuve %q", "webpay_token_123456", res.Token)
	}
}

func TestTransactionStatus_InputError(t *testing.T) {
	ms := NewMockServer()
	defer ms.Server.Close()
	mockClient := &mockClient{ms.Server.URL}
	tx := transbank.NewTransactionWithClient(mockClient, transactionOptions)

	tests := []struct {
		name          string
		token         string
		expectedError string
	}{
		{
			name:          "Token exceedes maximum length",
			token:         strings.Repeat("a", 65),
			expectedError: "SDK Error: token is too long, the maximum length is 64",
		},
		{
			name:          "Token is empty",
			token:         "",
			expectedError: "SDK Error: token cannot be empty",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tx.Status(tt.token)
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
	tx := transbank.NewTransactionWithClient(mockClient, transactionOptions)

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
	tx := transbank.NewTransactionWithClient(mockClient, transactionOptions)
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
	tx := transbank.NewTransactionWithClient(mockClient, transactionOptions)

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
	tx := transbank.NewTransactionWithClient(mockClient, transactionOptions)

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
	tx := transbank.NewTransactionWithClient(mockClient, transactionOptions)

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

func TestTransactionRefund_InputError(t *testing.T) {
	ms := NewMockServer()
	defer ms.Server.Close()
	mockClient := &mockClient{ms.Server.URL}
	tx := transbank.NewTransactionWithClient(mockClient, transactionOptions)

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
			expectedError: "SDK Error: token cannot be empty",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tx.Refund(tt.token, tt.amount)
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

func TestTransactionRefund_ServerError(t *testing.T) {
	ms := NewMockServer()
	defer ms.Close()
	mockClient := &mockClient{ms.Server.URL}
	tx := transbank.NewTransactionWithClient(mockClient, transactionOptions)
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

func TestTransactionRefund_ReverseSuccess(t *testing.T) {
	ms := NewMockServer()
	defer ms.Close()
	mockClient := &mockClient{ms.Server.URL}
	tx := transbank.NewTransactionWithClient(mockClient, transactionOptions)

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

func TestTransactionRefund_NullifiedSuccess(t *testing.T) {
	ms := NewMockServer()
	defer ms.Close()
	mockClient := &mockClient{ms.Server.URL}
	tx := transbank.NewTransactionWithClient(mockClient, transactionOptions)

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

func TestTransactionCapture_InputError(t *testing.T) {
	ms := NewMockServer()
	defer ms.Server.Close()
	mockClient := &mockClient{ms.Server.URL}
	tx := transbank.NewTransactionWithClient(mockClient, transactionOptions)

	tests := []struct {
		name              string
		token             string
		buyOrder          string
		authorizationCode string
		amount            float64
		expectedError     string
	}{
		{
			name:              "Token is empty",
			token:             "",
			buyOrder:          "123456",
			authorizationCode: "1213",
			amount:            10000,
			expectedError:     "SDK Error: token cannot be empty",
		},
		{
			name:              "Token exceedes maximum length",
			token:             strings.Repeat("a", 65),
			buyOrder:          "123456",
			authorizationCode: "1213",
			amount:            10000,
			expectedError:     "SDK Error: token is too long, the maximum length is 64",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tx.Capture(tt.token, tt.buyOrder, tt.authorizationCode, tt.amount)
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

func TestTransactionCapture_ServerError(t *testing.T) {
	ms := NewMockServer()
	defer ms.Close()
	mockClient := &mockClient{ms.Server.URL}
	tx := transbank.NewTransactionWithClient(mockClient, transactionOptions)
	tests := []struct {
		name              string
		token             string
		buyOrder          string
		authorizationCode string
		amount            float64
		mockStatusCode    int
		mockResponse      map[string]string
		expectedError     string
	}{
		{
			name:              "'authorization_code' is empty",
			token:             strings.Repeat("a", 64),
			buyOrder:          "123456",
			authorizationCode: "",
			amount:            1000,
			mockStatusCode:    422,
			mockResponse:      map[string]string{"error_message": "authorization_code is required!"},
			expectedError:     "authorization_code is required!",
		},
		{
			name:              "Wrong 'authorization_code'",
			token:             strings.Repeat("a", 64),
			buyOrder:          "123456",
			authorizationCode: "1212",
			amount:            1000,
			mockStatusCode:    400,
			mockResponse:      map[string]string{"error_message": "Transaction not found"},
			expectedError:     "Transaction not found",
		},
		{
			name:              "'buy_order' is empty",
			token:             strings.Repeat("a", 64),
			buyOrder:          "",
			authorizationCode: "1213",
			amount:            1000,
			mockStatusCode:    422,
			mockResponse:      map[string]string{"error_message": "buy_order is required!"},
			expectedError:     "buy_order is required!",
		},
		{
			name:              "Wrong 'buy_order'",
			token:             strings.Repeat("a", 64),
			buyOrder:          "123457",
			authorizationCode: "1213",
			amount:            1000,
			mockStatusCode:    400,
			mockResponse:      map[string]string{"error_message": "Transaction not found"},
			expectedError:     "Transaction not found",
		},
		{
			name:              "'capture_amount' cannot be zero",
			token:             strings.Repeat("a", 64),
			buyOrder:          "123456",
			authorizationCode: "1213",
			amount:            0,
			mockStatusCode:    422,
			mockResponse:      map[string]string{"error_message": "Invalid value for parameter: capture_amount"},
			expectedError:     "Invalid value for parameter: capture_amount",
		},
		{
			name:              "'capture_amount' cannot be a negative number",
			token:             strings.Repeat("a", 64),
			buyOrder:          "123456",
			authorizationCode: "1213",
			amount:            -1000,
			mockStatusCode:    422,
			mockResponse:      map[string]string{"error_message": "Invalid value for parameter: capture_amount"},
			expectedError:     "Invalid value for parameter: capture_amount",
		},
		{
			name:              "Token doesn't exist",
			token:             strings.Repeat("b", 64),
			buyOrder:          "123456",
			authorizationCode: "1213",
			amount:            1000,
			mockStatusCode:    422,
			mockResponse:      map[string]string{"error_message": "Invalid value for parameter: token"},
			expectedError:     "Invalid value for parameter: token",
		},
		{
			name:              "Token is less than 64 characters",
			token:             strings.Repeat("b", 20),
			buyOrder:          "123456",
			authorizationCode: "1213",
			amount:            1000,
			mockStatusCode:    422,
			mockResponse:      map[string]string{"error_message": "Invalid value for parameter: token"},
			expectedError:     "Invalid value for parameter: token",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms.StatusCode = tt.mockStatusCode
			ms.Response = tt.mockResponse

			_, err := tx.Capture(tt.token, tt.buyOrder, tt.authorizationCode, tt.amount)

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

func TestTransactionCapture_Success(t *testing.T) {
	ms := NewMockServer()
	defer ms.Close()
	mockClient := &mockClient{ms.Server.URL}
	ms.Response = map[string]any{
		"authorization_code": "120050",
		"authorization_date": "2026-04-07T01:28:04.145Z",
		"captured_amount":    10000,
		"response_code":      0,
	}
	ms.StatusCode = 200
	tx := transbank.NewTransactionWithClient(mockClient, transactionOptions)
	res, err := tx.Capture(strings.Repeat("a", 64), "buyOrder12345678", "1213", 10000)
	if err != nil {
		t.Fatalf("Expected result, got: %v", err)
	}
	if res.AuthorizationCode == "" {
		t.Errorf("authorization_code must not be empty. Expected 1213, got %s", res.AuthorizationCode)
	}
	if res.ResponseCode != 0 {
		t.Errorf("Wrong response_code. Expected 0, got %d", res.ResponseCode)
	}
}
