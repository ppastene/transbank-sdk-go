package tests

import (
	"errors"
	"strings"
	"testing"

	webpay "github.com/ppastene/transbank-sdk-go"
	"github.com/ppastene/transbank-sdk-go/internal/shared"
)

var mallTransactionOptions = &webpay.Options{
	ApiKey:       "api-key",
	CommerceCode: "commerce-code",
}

var validMallDetails = webpay.WebpayPlusMallDetails{
	Amount:       10000,
	CommerceCode: "597055555536",
	BuyOrder:     "m1-123456",
}

func TestMallTransaction_InitWithoutCredentials(t *testing.T) {
	mockClient := &mockClient{}

	tests := []struct {
		name          string
		options       *webpay.Options
		expectedError string
	}{
		{
			name:          "Credentials missing",
			options:       &webpay.Options{},
			expectedError: "No credentials",
		},
		{
			name:          "No Api Key",
			options:       &webpay.Options{CommerceCode: "579B532A7440BB0C9079DED94D31EA1615BACEB56610332264630D42D0A36B1C"},
			expectedError: "ApiKey is required",
		},
		{
			name:          "No Commerce Code",
			options:       &webpay.Options{ApiKey: "597055555540"},
			expectedError: "CommerceCode is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := webpay.NewMallTransactionWithClient(mockClient, tt.options)

			_, err := tx.Status("token")
			if err == nil {
				t.Fatalf("Expected error, got nil")
			}

			var werr *shared.WebpayError
			if !errors.As(err, &werr) {
				t.Errorf("Error is not a WebpayError")
			}

			cause := errors.Unwrap(err)
			if cause.Error() != tt.expectedError {
				t.Errorf("Wrong error cause.\nExpected: %q\nGot: %q", tt.expectedError, cause.Error())
			}
		})
	}
}

func TestMallTransactionCreate_ServerError(t *testing.T) {
	ms := NewMockServer()
	defer ms.Close()
	mockClient := &mockClient{ms.Server.URL}
	tx := webpay.NewMallTransactionWithClient(mockClient, mallTransactionOptions)
	tests := []struct {
		name           string
		buyOrder       string
		sessionId      string
		returnUrl      string
		details        []webpay.WebpayPlusMallDetails
		mockStatusCode int
		mockResponse   map[string]string
		expectedError  string
	}{
		{
			name:           "buy_order not found",
			buyOrder:       "",
			sessionId:      "session123456",
			returnUrl:      "https://webpay.cl/formulario-pago",
			details:        []webpay.WebpayPlusMallDetails{validMallDetails},
			mockResponse:   map[string]string{"error_message": "buy_order is required!"},
			mockStatusCode: 422,
			expectedError:  "buy_order is required!",
		},
		{
			name:           "buy_order not valid",
			buyOrder:       "p-123456àèìòù",
			sessionId:      "session123456",
			returnUrl:      "https://webpay.cl/formulario-pago",
			details:        []webpay.WebpayPlusMallDetails{validMallDetails},
			mockResponse:   map[string]string{"error_message": "\"Buy order\" rejected with value p-123456àèìòù"},
			mockStatusCode: 422,
			expectedError:  "\"Buy order\" rejected with value p-123456àèìòù",
		},
		{
			name:           "buy_order large exceeded",
			buyOrder:       strings.Repeat("a", 27),
			sessionId:      "session123456",
			returnUrl:      "https://webpay.cl/formulario-pago",
			details:        []webpay.WebpayPlusMallDetails{validMallDetails},
			mockResponse:   map[string]string{"error_message": "Invalid value for parameter: buy_order"},
			mockStatusCode: 422,
			expectedError:  "Invalid value for parameter: buy_order",
		},
		{
			name:           "session_id not found",
			buyOrder:       "p-123456",
			sessionId:      "",
			returnUrl:      "https://webpay.cl/formulario-pago",
			details:        []webpay.WebpayPlusMallDetails{validMallDetails},
			mockResponse:   map[string]string{"error_message": "session_id is required!"},
			mockStatusCode: 422,
			expectedError:  "session_id is required!",
		},
		{
			name:           "session_id large exceeded",
			buyOrder:       "p-123456",
			sessionId:      strings.Repeat("a", 62),
			returnUrl:      "https://webpay.cl/formulario-pago",
			details:        []webpay.WebpayPlusMallDetails{validMallDetails},
			mockResponse:   map[string]string{"error_message": "Invalid value for parameter: session_id"},
			mockStatusCode: 422,
			expectedError:  "Invalid value for parameter: session_id",
		},
		{
			name:           "return_url not found",
			buyOrder:       "p-123456",
			sessionId:      "session123456",
			returnUrl:      "",
			details:        []webpay.WebpayPlusMallDetails{validMallDetails},
			mockResponse:   map[string]string{"error_message": "return_url is required!"},
			mockStatusCode: 422,
			expectedError:  "return_url is required!",
		},
		{
			name:           "return_url not valid",
			buyOrder:       "p-123456",
			sessionId:      "session123456",
			returnUrl:      strings.Repeat("a", 256),
			details:        []webpay.WebpayPlusMallDetails{validMallDetails},
			mockResponse:   map[string]string{"error_message": "Invalid value for parameter: return_url"},
			mockStatusCode: 422,
			expectedError:  "Invalid value for parameter: return_url",
		},
		{
			name:           "details is empty",
			buyOrder:       "p-123456",
			sessionId:      "session123456",
			returnUrl:      "https://webpay.cl/formulario-pago",
			details:        []webpay.WebpayPlusMallDetails{},
			mockResponse:   map[string]string{"error_message": "at least one detail is required"},
			mockStatusCode: 422,
			expectedError:  "at least one detail is required",
		},
		{
			name:      "details[0].amount not valid",
			buyOrder:  "p-123456",
			sessionId: "session123456",
			returnUrl: "https://webpay.cl/formulario-pago",
			details: []webpay.WebpayPlusMallDetails{
				{
					Amount:       0,
					CommerceCode: "597055555536",
					BuyOrder:     "m1-123456",
				},
			},
			mockResponse:   map[string]string{"error_message": "Invalid value for parameter: details[0].amount"},
			mockStatusCode: 422,
			expectedError:  "Invalid value for parameter: details[0].amount",
		},
		{
			name:      "details[0].commerce_code not found",
			buyOrder:  "p-123456",
			sessionId: "session123456",
			returnUrl: "https://webpay.cl/formulario-pago",
			details: []webpay.WebpayPlusMallDetails{
				{
					Amount:       10000,
					CommerceCode: "",
					BuyOrder:     "m1-123456",
				},
			},
			mockResponse:   map[string]string{"error_message": "details[0].commerce_code is required!"},
			mockStatusCode: 422,
			expectedError:  "details[0].commerce_code is required!",
		},
		{
			name:      "details[0].commerce_code not valid",
			buyOrder:  "p-123456",
			sessionId: "session123456",
			returnUrl: "https://webpay.cl/formulario-pago",
			details: []webpay.WebpayPlusMallDetails{
				{
					Amount:       10000,
					CommerceCode: strings.Repeat("a", 13),
					BuyOrder:     "m1-123456",
				},
			},
			mockResponse:   map[string]string{"error_message": "Invalid value for parameter: details[0].commerce_code"},
			mockStatusCode: 422,
			expectedError:  "Invalid value for parameter: details[0].commerce_code",
		},
		{
			name:      "details[0].commerce_code not registered",
			buyOrder:  "p-123456",
			sessionId: "session123456",
			returnUrl: "https://webpay.cl/formulario-pago",
			details: []webpay.WebpayPlusMallDetails{
				{
					Amount:       10000,
					CommerceCode: strings.Repeat("1", 12),
					BuyOrder:     "m1-123456",
				},
			},
			mockResponse:   map[string]string{"error_message": "Unexpected error"},
			mockStatusCode: 500,
			expectedError:  "Unexpected error",
		},
		{
			name:      "details[0].buy_order not found",
			buyOrder:  "p-123456",
			sessionId: "session123456",
			returnUrl: "https://webpay.cl/formulario-pago",
			details: []webpay.WebpayPlusMallDetails{
				{
					Amount:       10000,
					CommerceCode: "597055555536",
					BuyOrder:     "",
				},
			},
			mockResponse:   map[string]string{"error_message": "details[0].buy_order is required!"},
			mockStatusCode: 422,
			expectedError:  "details[0].buy_order is required!",
		},
		{
			name:      "details[0].buy_order not valid",
			buyOrder:  "p-123456",
			sessionId: "session123456",
			returnUrl: "https://webpay.cl/formulario-pago",
			details: []webpay.WebpayPlusMallDetails{
				{
					Amount:       10000,
					CommerceCode: "597055555536",
					BuyOrder:     "m1-123456àèìoù",
				},
			},
			mockResponse:   map[string]string{"error_message": "Parameter \"Detail buy order\" rejected with value m1-123456àèìoù"},
			mockStatusCode: 422,
			expectedError:  "Parameter \"Detail buy order\" rejected with value m1-123456àèìoù",
		},
		{
			name:      "details[0].buy_order large exceeded",
			buyOrder:  "p-123456",
			sessionId: "session123456",
			returnUrl: "https://webpay.cl/formulario-pago",
			details: []webpay.WebpayPlusMallDetails{
				{
					Amount:       10000,
					CommerceCode: "597055555536",
					BuyOrder:     strings.Repeat("a", 27),
				},
			},
			mockResponse:   map[string]string{"error_message": "Invalid value for parameter: details[0].buy_order"},
			mockStatusCode: 422,
			expectedError:  "Invalid value for parameter: details[0].buy_order",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms.StatusCode = tt.mockStatusCode
			ms.Response = tt.mockResponse

			_, err := tx.Create(tt.buyOrder, tt.sessionId, tt.returnUrl, tt.details)

			if err == nil {
				t.Fatalf("Expected error, got nil")
			}

			cause := errors.Unwrap(err)

			if cause.Error() != tt.expectedError {
				t.Errorf("Wrong error cause.\nExpected: %q\nGot: %q)", tt.expectedError, cause.Error())
			}

			var werr *shared.WebpayError
			if errors.As(err, &werr) {
				if werr.Code != tt.mockStatusCode {
					t.Errorf("Wrong error code. Expected %d, Got %d", tt.mockStatusCode, werr.Code)
				}
			}
		})
	}
}

func TestMallTransactionCreate_Success(t *testing.T) {
	ms := NewMockServer()
	defer ms.Server.Close()
	ms.Response = map[string]string{
		"token": "webpay_token_123456",
		"url":   "https://webpay.cl/formulario-pago",
	}
	ms.StatusCode = 200
	mockClient := &mockClient{ms.Server.URL}
	tx := webpay.NewMallTransactionWithClient(mockClient, mallTransactionOptions)
	res, err := tx.Create("p-123456", "session123456", "https://webpay.cl/formulario-pago", []webpay.WebpayPlusMallDetails{validMallDetails})
	if err != nil {
		t.Fatalf("Expected result, got: %v", err)
	}
	if res.Token != "webpay_token_123456" {
		t.Errorf("Wrong token. Expected %q, got %q", "webpay_token_123456", res.Token)
	}
}

func TestMallTransactionStatus_InputError(t *testing.T) {
	mockClient := &mockClient{}
	tx := webpay.NewMallTransactionWithClient(mockClient, mallTransactionOptions)
	tests := []struct {
		name          string
		token         string
		expectedError string
	}{
		{
			name:          "Token exceedes maximum length",
			token:         strings.Repeat("a", 65),
			expectedError: "token is too long, the maximum length is 64",
		},
		{
			name:          "Token is empty",
			token:         "",
			expectedError: "token cannot be empty",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tx.Status(tt.token)
			if err == nil {
				t.Fatalf("Expected error, got nil")
			}

			cause := errors.Unwrap(err)

			if cause.Error() != tt.expectedError {
				t.Errorf("Wrong error cause.\nExpected: %q\nGot: %q)", tt.expectedError, cause.Error())
			}

			var werr *shared.WebpayError
			if errors.As(err, &werr) {
				if werr.Code != -1 {
					t.Errorf("Wrong error code. Expected: -1, Got: %d", werr.Code)
				}
			}
		})
	}
}

func TestMallTransactionStatus_ServerError(t *testing.T) {
	ms := NewMockServer()
	defer ms.Close()
	mockClient := &mockClient{ms.Server.URL}
	tx := webpay.NewTransactionWithClient(mockClient, mallTransactionOptions)

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
				t.Fatalf("Expected error, got nil")
			}

			cause := errors.Unwrap(err)

			if cause.Error() != tt.expectedError {
				t.Errorf("Wrong error cause.\nExpected: %q\nGot: %q)", tt.expectedError, cause.Error())
			}

			var werr *shared.WebpayError
			if errors.As(err, &werr) {
				if werr.Code != tt.mockStatusCode {
					t.Errorf("Wrong error code. Expected %d, Got %d", tt.mockStatusCode, werr.Code)
				}
			}
		})
	}
}

func TestMallTransactionStatus_Success(t *testing.T) {
	ms := NewMockServer()
	defer ms.Server.Close()
	ms.Response = map[string]any{
		"vci": "TSY",
		"details": []any{
			map[string]any{
				"amount":              10000,
				"status":              "AUTHORIZED",
				"authorization_code":  "1213",
				"payment_type_code":   "VN",
				"response_code":       0,
				"installments_number": 0,
				"commerce_code":       "597055555536",
				"buy_order":           "m1-123456",
			},
			map[string]any{
				"amount":              20000,
				"status":              "AUTHORIZED",
				"authorization_code":  "1213",
				"payment_type_code":   "VN",
				"response_code":       0,
				"installments_number": 0,
				"commerce_code":       "597055555537",
				"buy_order":           "m2-123456",
			},
		},
		"buy_order":  "p-123456",
		"session_id": "session123456",
		"card_detail": map[string]any{
			"card_number": "6623",
		},
		"accounting_date":  "0302",
		"transaction_date": "2026-03-03T00:14:17.848Z",
	}

	ms.StatusCode = 200

	mockClient := &mockClient{ms.Server.URL}
	tx := webpay.NewMallTransactionWithClient(mockClient, mallTransactionOptions)
	res, err := tx.Status(strings.Repeat("a", 64))
	if err != nil {
		t.Fatalf("Expected result, got: %v", err)
	}

	if res.Vci != "TSY" {
		t.Errorf("Wrong vci. Expected TSY, got %s", res.Vci)
	}

	if res.CardDetail.CardNumber != "6623" {
		t.Errorf("Wrong card_number. Expected 6623, got %s", res.CardDetail.CardNumber)
	}
	for _, detail := range res.Details {
		if detail.ResponseCode != 0 {
			t.Errorf("Wrong response_code. Expected 0, got %d", detail.ResponseCode)
		}
		if detail.Status != "AUTHORIZED" {
			t.Errorf("Wrong status. Expected AUTHORIZED, got %s", detail.Status)
		}
		if detail.AuthorizationCode != "1213" {
			t.Errorf("Wrong authorization_code. Expected 1213, got %s", detail.AuthorizationCode)
		}
	}
}

func TestMallTransactionRefund_InputError(t *testing.T) {
	mockClient := &mockClient{}
	tx := webpay.NewMallTransactionWithClient(mockClient, mallTransactionOptions)
	tests := []struct {
		name          string
		token         string
		buyOrder      string
		amount        float64
		commerceCode  string
		expectedError string
	}{
		{
			name:          "Token exceedes maximum length",
			token:         strings.Repeat("a", 65),
			buyOrder:      "ordenCompra12345678",
			amount:        10000,
			commerceCode:  "597055555536",
			expectedError: "token is too long, the maximum length is 64",
		},
		{
			name:          "Token is empty",
			token:         "",
			buyOrder:      "ordenCompra12345678",
			amount:        10000,
			commerceCode:  "597055555536",
			expectedError: "token cannot be empty",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tx.Refund(tt.token, tt.buyOrder, tt.commerceCode, tt.amount)
			if err == nil {
				t.Fatalf("Expected error, got nil")
			}

			cause := errors.Unwrap(err)

			if cause.Error() != tt.expectedError {
				t.Errorf("Wrong error cause.\nExpected: %q\nGot: %q)", tt.expectedError, cause.Error())
			}

			var werr *shared.WebpayError
			if errors.As(err, &werr) {
				if werr.Code != -1 {
					t.Errorf("Wrong error code. Expected: -1, Got: %d", werr.Code)
				}
			}
		})
	}
}

func TestMallTransactionRefund_ServerError(t *testing.T) {
	ms := NewMockServer()
	defer ms.Server.Close()
	ms.StatusCode = 200

	mockClient := &mockClient{ms.Server.URL}
	tx := webpay.NewMallTransactionWithClient(mockClient, mallTransactionOptions)
	tests := []struct {
		name           string
		token          string
		buyOrder       string
		amount         float64
		commerceCode   string
		mockStatusCode int
		mockResponse   map[string]string
		expectedError  string
	}{
		{
			name:           "commerce_code is empty",
			token:          strings.Repeat("a", 64),
			buyOrder:       "ordenCompra12345678",
			amount:         10000,
			commerceCode:   "",
			mockStatusCode: 422,
			mockResponse:   map[string]string{"error_message": "commerce_code is required!"},
			expectedError:  "commerce_code is required!",
		},
		{
			name:           "commerce_code is not valid",
			token:          strings.Repeat("a", 64),
			buyOrder:       "ordenCompra12345678",
			amount:         10000,
			commerceCode:   strings.Repeat("a", 13),
			mockStatusCode: 422,
			mockResponse:   map[string]string{"error_message": "Invalid value for parameter: commerce_code"},
			expectedError:  "Invalid value for parameter: commerce_code",
		},
		{
			name:           "buy_order is empty",
			token:          strings.Repeat("a", 64),
			buyOrder:       "",
			amount:         10000,
			commerceCode:   "597055555536",
			mockStatusCode: 422,
			mockResponse:   map[string]string{"error_message": "buy_order is required!"},
			expectedError:  "buy_order is required!",
		},
		{
			name:           "buy_order does no exist",
			token:          strings.Repeat("a", 64),
			buyOrder:       "buyOrderThatDoesNotExist",
			amount:         10000,
			commerceCode:   "597055555536",
			mockStatusCode: 422,
			mockResponse:   map[string]string{"error_message": "Transaction detail was not found"},
			expectedError:  "Transaction detail was not found",
		},
		{
			name:           "buy_order not valid",
			token:          strings.Repeat("a", 64),
			buyOrder:       strings.Repeat("a", 27),
			amount:         10000,
			commerceCode:   "597055555536",
			mockStatusCode: 422,
			mockResponse:   map[string]string{"error_message": "Invalid value for parameter: buy_order"},
			expectedError:  "Invalid value for parameter: buy_order",
		},
		{
			name:           "amount is 0",
			token:          strings.Repeat("a", 64),
			buyOrder:       "ordenCompra12345678",
			amount:         0,
			commerceCode:   "597055555536",
			mockStatusCode: 422,
			mockResponse:   map[string]string{"error_message": "Invalid value for parameter: amount"},
			expectedError:  "Invalid value for parameter: amount",
		},
		{
			name:           "amount is bigger than authorized",
			token:          strings.Repeat("a", 64),
			buyOrder:       "ordenCompra12345678",
			amount:         10000000000,
			commerceCode:   "597055555536",
			mockStatusCode: 422,
			mockResponse:   map[string]string{"error_message": "Amount to refund is bigger than authorized"},
			expectedError:  "Amount to refund is bigger than authorized",
		},
		{
			name:           "transaction not found",
			token:          strings.Repeat("a", 64),
			buyOrder:       "ordenCompra12345678",
			amount:         500,
			commerceCode:   "597055555536",
			mockStatusCode: 400,
			mockResponse:   map[string]string{"error_message": "Transaction not found"},
			expectedError:  "Transaction not found",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms.StatusCode = tt.mockStatusCode
			ms.Response = tt.mockResponse

			_, err := tx.Refund(tt.token, tt.buyOrder, tt.commerceCode, tt.amount)

			if err == nil {
				t.Fatalf("Expected error, got nil")
			}

			cause := errors.Unwrap(err)

			if cause.Error() != tt.expectedError {
				t.Errorf("Wrong error cause.\nExpected: %q\nGot: %q)", tt.expectedError, cause.Error())
			}

			var werr *shared.WebpayError
			if errors.As(err, &werr) {
				if werr.Code != tt.mockStatusCode {
					t.Errorf("Wrong error code. Expected %d, Got %d", tt.mockStatusCode, werr.Code)
				}
			}
		})
	}
}

func TestMallTransactionRefund_Success(t *testing.T) {
	ms := NewMockServer()
	defer ms.Server.Close()
	ms.Response = map[string]string{
		"type": "REVERSED",
	}
	ms.StatusCode = 200
	mockClient := &mockClient{ms.Server.URL}
	tx := webpay.NewMallTransactionWithClient(mockClient, mallTransactionOptions)
	res, err := tx.Refund("token", "buyOrder12345678", "597055555536", 10000)
	if err != nil {
		t.Fatalf("Expected result, got: %v", err)
	}
	if res.Type != "REVERSED" {
		t.Errorf("Response error. Expected REVERSED, got %s", res.Type)
	}
}

func TestMallTransactionRefund_NullifiedSuccess(t *testing.T) {
	ms := NewMockServer()
	defer ms.Close()
	mockClient := &mockClient{ms.Server.URL}
	ms.Response = map[string]any{
		"type":               "NULLIFIED",
		"authorization_code": "123456",
		"authorization_date": "2019-03-20T20:18:20Z",
		"nullified_amount":   500.00,
		"balance":            9500.00,
		"response_code":      0,
	}
	ms.StatusCode = 200
	tx := webpay.NewMallTransactionWithClient(mockClient, transactionOptions)
	res, err := tx.Refund("token", "buyOrder12345678", "597055555536", 500)
	if err != nil {
		t.Fatalf("Expected result, got: %v", err)
	}
	if res.Type != "NULLIFIED" {
		t.Errorf("Response error. Expected NULLIFIED, got %s", res.Type)
	}
	if res.AuthorizationCode == "" {
		t.Errorf("authorization_code must not be empty. Expected 123456, got %s", res.Type)
	}
	if res.ResponseCode != 0 {
		t.Errorf("Wrong response_code. Expected 0, got %s", res.Type)
	}
}
