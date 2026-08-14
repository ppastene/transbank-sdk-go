package transaccioncompleta_test

import (
	"net/http"
	"testing"

	"github.com/ppastene/transbank-sdk-go"
	"github.com/ppastene/transbank-sdk-go/transaccioncompleta"
)

func intPtr(i int) *int {
	return &i
}

func boolPtr(b bool) *bool {
	return &b
}

func TestNewTransaction(t *testing.T) {
	tests := []struct {
		name string
		opts transbank.Options
		want *transbank.ValidationError
	}{
		{
			name: "valid integration",
			opts: transbank.Options{
				CommerceCode: testCommerceCode,
				ApiKey:       testAPIKey,
				Environment:  transbank.Integration,
			},
		},
		{
			name: "valid production",
			opts: transbank.Options{
				CommerceCode: testCommerceCode,
				ApiKey:       testAPIKey,
				Environment:  transbank.Production,
			},
		},
		{
			name: "invalid environment",
			opts: transbank.Options{
				CommerceCode: testCommerceCode,
				ApiKey:       testAPIKey,
				Environment:  transbank.Environment(99),
			},
			want: &transbank.ValidationError{},
		},
		{
			name: "non numeric commerce code",
			opts: transbank.Options{
				CommerceCode: "59705555553a",
				ApiKey:       testAPIKey,
				Environment:  transbank.Integration,
			},
			want: &transbank.ValidationError{},
		},
		{
			name: "empty api key",
			opts: transbank.Options{
				CommerceCode: testCommerceCode,
				Environment:  transbank.Integration,
			},
			want: &transbank.ValidationError{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx, err := transaccioncompleta.NewTransaction(tt.opts)
			if tt.want != nil {
				wantValidationError(t, err)
				if tx != nil {
					t.Error("expected nil Transaction on error")
				}
				return
			}
			if err != nil {
				t.Fatalf("NewTransaction returned error: %v", err)
			}
			if tx == nil {
				t.Error("expected non-nil Transaction")
			}
		})
	}
}

func TestTransactionCreateWithCVV(t *testing.T) {
	server := newMockServer(t, http.StatusOK, `{"token":"`+testToken+`"}`)
	tx, err := transaccioncompleta.NewTransaction(testOptions(server))
	if err != nil {
		t.Fatalf("NewTransaction: %v", err)
	}

	resp, err := tx.Create(testBuyOrder, testSessionID, testAmount, testCardNumber, testCardExpiry, testCVV)
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	if resp.Token != testToken {
		t.Errorf("Token = %q, want %q", resp.Token, testToken)
	}

	req := server.LastRequest()
	assertRequest(t, req, http.MethodPost, testTransactionsPath)
	assertBody(t, req, `{"amount":10000,"buy_order":"orden-compra-123","card_expiration_date":"22/10","card_number":"4051885600446623","cvv":"123","session_id":"sesion-456"}`)
}

func TestTransactionCreateWithoutCVV(t *testing.T) {
	server := newMockServer(t, http.StatusOK, `{"token":"`+testToken+`"}`)
	tx, err := transaccioncompleta.NewTransaction(testOptions(server))
	if err != nil {
		t.Fatalf("NewTransaction: %v", err)
	}

	resp, err := tx.Create(testBuyOrder, testSessionID, testAmount, testCardNumber, testCardExpiry, "")
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	if resp.Token != testToken {
		t.Errorf("Token = %q, want %q", resp.Token, testToken)
	}

	req := server.LastRequest()
	assertRequest(t, req, http.MethodPost, testTransactionsPath)
	assertBody(t, req, `{"amount":10000,"buy_order":"orden-compra-123","card_expiration_date":"22/10","card_number":"4051885600446623","session_id":"sesion-456"}`)
	assertBodyNotContains(t, req, "cvv")
}

func TestTransactionCreateHTTPError(t *testing.T) {
	server := newMockServer(t, http.StatusBadRequest, `{"error_message":"bad request"}`)
	tx, err := transaccioncompleta.NewTransaction(testOptions(server))
	if err != nil {
		t.Fatalf("NewTransaction: %v", err)
	}

	_, err = tx.Create(testBuyOrder, testSessionID, testAmount, testCardNumber, testCardExpiry, testCVV)
	wantHTTPError(t, err, http.StatusBadRequest, `{"error_message":"bad request"}`)
}

func TestTransactionCreateValidation(t *testing.T) {
	server := newMockServer(t, http.StatusOK, `{"token":"`+testToken+`"}`)
	tx, err := transaccioncompleta.NewTransaction(testOptions(server))
	if err != nil {
		t.Fatalf("NewTransaction: %v", err)
	}

	tests := []struct {
		name           string
		buyOrder       string
		sessionID      string
		amount         float64
		cardNumber     string
		cardExpiration string
		cvv            string
	}{
		{name: "empty buy order", sessionID: testSessionID, amount: testAmount, cardNumber: testCardNumber, cardExpiration: testCardExpiry, cvv: testCVV},
		{name: "empty session id", buyOrder: testBuyOrder, amount: testAmount, cardNumber: testCardNumber, cardExpiration: testCardExpiry, cvv: testCVV},
		{name: "zero amount", buyOrder: testBuyOrder, sessionID: testSessionID, amount: 0, cardNumber: testCardNumber, cardExpiration: testCardExpiry, cvv: testCVV},
		{name: "empty card number", buyOrder: testBuyOrder, sessionID: testSessionID, amount: testAmount, cardExpiration: testCardExpiry, cvv: testCVV},
		{name: "non numeric card number", buyOrder: testBuyOrder, sessionID: testSessionID, amount: testAmount, cardNumber: "4239-0000", cardExpiration: testCardExpiry, cvv: testCVV},
		{name: "empty card expiration", buyOrder: testBuyOrder, sessionID: testSessionID, amount: testAmount, cardNumber: testCardNumber, cvv: testCVV},
		{name: "bad card expiration", buyOrder: testBuyOrder, sessionID: testSessionID, amount: testAmount, cardNumber: testCardNumber, cardExpiration: "2210", cvv: testCVV},
		{name: "cvv too long", buyOrder: testBuyOrder, sessionID: testSessionID, amount: testAmount, cardNumber: testCardNumber, cardExpiration: testCardExpiry, cvv: "12345"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tx.Create(tt.buyOrder, tt.sessionID, tt.amount, tt.cardNumber, tt.cardExpiration, tt.cvv)
			wantValidationError(t, err)
			if server.RequestCount() != 0 {
				t.Errorf("request count = %d, want 0 (validation must happen before the API call)", server.RequestCount())
			}
		})
	}
}

func TestTransactionInstallments(t *testing.T) {
	server := newMockServer(t, http.StatusOK, `{"installments_amount":3334,"id_query_installments":11,"deferred_periods":[{"amount":1000,"period":1}]}`)
	tx, err := transaccioncompleta.NewTransaction(testOptions(server))
	if err != nil {
		t.Fatalf("NewTransaction: %v", err)
	}

	resp, err := tx.Installments(testToken, 10)
	if err != nil {
		t.Fatalf("Installments returned error: %v", err)
	}
	if resp.InstallmentsAmount != 3334 {
		t.Errorf("InstallmentsAmount = %d, want 3334", resp.InstallmentsAmount)
	}
	if resp.IdQueryInstallments != 11 {
		t.Errorf("IdQueryInstallments = %d, want 11", resp.IdQueryInstallments)
	}
	if len(resp.DeferredPeriods) != 1 {
		t.Fatalf("len(DeferredPeriods) = %d, want 1", len(resp.DeferredPeriods))
	}
	if resp.DeferredPeriods[0].Amount != 1000 {
		t.Errorf("DeferredPeriods[0].Amount = %d, want 1000", resp.DeferredPeriods[0].Amount)
	}
	if resp.DeferredPeriods[0].Period != 1 {
		t.Errorf("DeferredPeriods[0].Period = %d, want 1", resp.DeferredPeriods[0].Period)
	}

	req := server.LastRequest()
	assertRequest(t, req, http.MethodPost, testTransactionsPath+"/"+testToken+"/installments")
	assertBody(t, req, `{"installments_number":10}`)
}

func TestTransactionInstallmentsValidation(t *testing.T) {
	server := newMockServer(t, http.StatusOK, `{}`)
	tx, err := transaccioncompleta.NewTransaction(testOptions(server))
	if err != nil {
		t.Fatalf("NewTransaction: %v", err)
	}

	tests := []struct {
		name               string
		token              string
		installmentsNumber int
	}{
		{name: "empty token", installmentsNumber: 10},
		{name: "zero installments", token: testToken, installmentsNumber: 0},
		{name: "negative installments", token: testToken, installmentsNumber: -1},
		{name: "100 installments", token: testToken, installmentsNumber: 100},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tx.Installments(tt.token, tt.installmentsNumber)
			wantValidationError(t, err)
			if server.RequestCount() != 0 {
				t.Errorf("request count = %d, want 0", server.RequestCount())
			}
		})
	}
}

func TestTransactionCommitNoInstallments(t *testing.T) {
	server := newMockServer(t, http.StatusOK, `{"amount":10000,"status":"AUTHORIZED","buy_order":"orden-compra-123","response_code":0}`)
	tx, err := transaccioncompleta.NewTransaction(testOptions(server))
	if err != nil {
		t.Fatalf("NewTransaction: %v", err)
	}

	resp, err := tx.Commit(testToken, nil, nil, nil)
	if err != nil {
		t.Fatalf("Commit returned error: %v", err)
	}
	if resp.Amount != 10000 {
		t.Errorf("Amount = %f, want 10000", resp.Amount)
	}
	if resp.Status != "AUTHORIZED" {
		t.Errorf("Status = %q, want AUTHORIZED", resp.Status)
	}

	req := server.LastRequest()
	assertRequest(t, req, http.MethodPut, testTransactionsPath+"/"+testToken)
	assertBody(t, req, `{}`)
	assertBodyNotContains(t, req, "null")
}

func TestTransactionCommitWithInstallments(t *testing.T) {
	server := newMockServer(t, http.StatusOK, `{"amount":10000,"status":"AUTHORIZED","response_code":0}`)
	tx, err := transaccioncompleta.NewTransaction(testOptions(server))
	if err != nil {
		t.Fatalf("NewTransaction: %v", err)
	}

	resp, err := tx.Commit(testToken, intPtr(15), intPtr(1), boolPtr(false))
	if err != nil {
		t.Fatalf("Commit returned error: %v", err)
	}
	if !resp.IsApproved() {
		t.Error("IsApproved = false, want true")
	}

	req := server.LastRequest()
	assertRequest(t, req, http.MethodPut, testTransactionsPath+"/"+testToken)
	assertBody(t, req, `{"deferred_period_index":1,"grace_period":false,"id_query_installments":15}`)
	assertBodyNotContains(t, req, "null")
}

func TestTransactionCommitValidation(t *testing.T) {
	server := newMockServer(t, http.StatusOK, `{}`)
	tx, err := transaccioncompleta.NewTransaction(testOptions(server))
	if err != nil {
		t.Fatalf("NewTransaction: %v", err)
	}

	_, err = tx.Commit("", nil, nil, nil)
	wantValidationError(t, err)
	if server.RequestCount() != 0 {
		t.Errorf("request count = %d, want 0", server.RequestCount())
	}
}

func TestTransactionStatus(t *testing.T) {
	server := newMockServer(t, http.StatusOK, `{"amount":10000,"status":"AUTHORIZED","buy_order":"orden-compra-123","card_detail":{"card_number":"4239"},"response_code":0}`)
	tx, err := transaccioncompleta.NewTransaction(testOptions(server))
	if err != nil {
		t.Fatalf("NewTransaction: %v", err)
	}

	resp, err := tx.Status(testToken)
	if err != nil {
		t.Fatalf("Status returned error: %v", err)
	}
	if resp.BuyOrder != testBuyOrder {
		t.Errorf("BuyOrder = %q, want %q", resp.BuyOrder, testBuyOrder)
	}
	if resp.CardDetail.CardNumber != "4239" {
		t.Errorf("CardDetail.CardNumber = %q, want 4239", resp.CardDetail.CardNumber)
	}
	if !resp.IsApproved() {
		t.Error("IsApproved = false, want true")
	}

	req := server.LastRequest()
	assertRequest(t, req, http.MethodGet, testTransactionsPath+"/"+testToken)
}

func TestTransactionStatusValidation(t *testing.T) {
	server := newMockServer(t, http.StatusOK, `{}`)
	tx, err := transaccioncompleta.NewTransaction(testOptions(server))
	if err != nil {
		t.Fatalf("NewTransaction: %v", err)
	}

	_, err = tx.Status("")
	wantValidationError(t, err)
	if server.RequestCount() != 0 {
		t.Errorf("request count = %d, want 0", server.RequestCount())
	}
}

func TestTransactionRefund(t *testing.T) {
	server := newMockServer(t, http.StatusOK, `{"type":"NULLIFY","authorization_code":"123456","authorization_date":"2019-03-20T20:18:20Z","nullified_amount":1000.00,"balance":0.00,"response_code":0}`)
	tx, err := transaccioncompleta.NewTransaction(testOptions(server))
	if err != nil {
		t.Fatalf("NewTransaction: %v", err)
	}

	resp, err := tx.Refund(testToken, 1000)
	if err != nil {
		t.Fatalf("Refund returned error: %v", err)
	}
	if resp.Type != "NULLIFY" {
		t.Errorf("Type = %q, want NULLIFY", resp.Type)
	}
	if resp.NullifiedAmount != 1000 {
		t.Errorf("NullifiedAmount = %f, want 1000", resp.NullifiedAmount)
	}
	if resp.ResponseCode != 0 {
		t.Errorf("ResponseCode = %d, want 0", resp.ResponseCode)
	}

	req := server.LastRequest()
	assertRequest(t, req, http.MethodPost, testTransactionsPath+"/"+testToken+"/refunds")
	assertBody(t, req, `{"amount":1000}`)
}

func TestTransactionRefundValidation(t *testing.T) {
	server := newMockServer(t, http.StatusOK, `{}`)
	tx, err := transaccioncompleta.NewTransaction(testOptions(server))
	if err != nil {
		t.Fatalf("NewTransaction: %v", err)
	}

	tests := []struct {
		name   string
		token  string
		amount float64
	}{
		{name: "empty token", token: "", amount: 1000},
		{name: "zero amount", token: testToken, amount: 0},
		{name: "negative amount", token: testToken, amount: -1000},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tx.Refund(tt.token, tt.amount)
			wantValidationError(t, err)
			if server.RequestCount() != 0 {
				t.Errorf("request count = %d, want 0", server.RequestCount())
			}
		})
	}
}

func TestTransactionCapture(t *testing.T) {
	server := newMockServer(t, http.StatusOK, `{"authorization_code":"123456","authorization_date":"2019-03-20T20:18:20Z","captured_amount":1000,"response_code":0}`)
	tx, err := transaccioncompleta.NewTransaction(testOptions(server))
	if err != nil {
		t.Fatalf("NewTransaction: %v", err)
	}

	resp, err := tx.Capture(testToken, testBuyOrder, testAuthCode, 1000)
	if err != nil {
		t.Fatalf("Capture returned error: %v", err)
	}
	if resp.CapturedAmount != 1000 {
		t.Errorf("CapturedAmount = %f, want 1000", resp.CapturedAmount)
	}
	if resp.AuthorizationCode != testAuthCode {
		t.Errorf("AuthorizationCode = %q, want %q", resp.AuthorizationCode, testAuthCode)
	}

	req := server.LastRequest()
	assertRequest(t, req, http.MethodPut, testTransactionsPath+"/"+testToken+"/capture")
	assertBody(t, req, `{"authorization_code":"123456","buy_order":"orden-compra-123","capture_amount":1000}`)
}

func TestTransactionCaptureValidation(t *testing.T) {
	server := newMockServer(t, http.StatusOK, `{}`)
	tx, err := transaccioncompleta.NewTransaction(testOptions(server))
	if err != nil {
		t.Fatalf("NewTransaction: %v", err)
	}

	tests := []struct {
		name          string
		token         string
		buyOrder      string
		authorization string
		captureAmount float64
	}{
		{name: "empty token", buyOrder: testBuyOrder, authorization: testAuthCode, captureAmount: 1000},
		{name: "empty buy order", token: testToken, authorization: testAuthCode, captureAmount: 1000},
		{name: "empty authorization code", token: testToken, buyOrder: testBuyOrder, captureAmount: 1000},
		{name: "zero capture amount", token: testToken, buyOrder: testBuyOrder, authorization: testAuthCode, captureAmount: 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tx.Capture(tt.token, tt.buyOrder, tt.authorization, tt.captureAmount)
			wantValidationError(t, err)
			if server.RequestCount() != 0 {
				t.Errorf("request count = %d, want 0", server.RequestCount())
			}
		})
	}
}
