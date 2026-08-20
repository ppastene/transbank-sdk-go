package transaccioncompleta_test

import (
	"errors"
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

func TestTransactionHTTPError(t *testing.T) {
	tests := []struct {
		name   string
		method func(tx *transaccioncompleta.Transaction) error
	}{
		{
			name: "create",
			method: func(tx *transaccioncompleta.Transaction) error {
				_, err := tx.Create(testBuyOrder, testSessionID, testAmount, testCardNumber, testCardExpiry, testCVV)
				return err
			},
		},
		{
			name: "commit",
			method: func(tx *transaccioncompleta.Transaction) error {
				_, err := tx.Commit(testToken, nil, nil, nil)
				return err
			},
		},
		{
			name: "status",
			method: func(tx *transaccioncompleta.Transaction) error {
				_, err := tx.Status(testToken)
				return err
			},
		},
		{
			name: "refund",
			method: func(tx *transaccioncompleta.Transaction) error {
				_, err := tx.Refund(testToken, 1000)
				return err
			},
		},
		{
			name: "capture",
			method: func(tx *transaccioncompleta.Transaction) error {
				_, err := tx.Capture(testToken, testBuyOrder, testAuthCode, 1000)
				return err
			},
		},
		{
			name: "installments",
			method: func(tx *transaccioncompleta.Transaction) error {
				_, err := tx.Installments(testToken, 10)
				return err
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := newMockServer(t, http.StatusUnauthorized, `{"error_message":"unauthorized"}`)
			tx, err := transaccioncompleta.NewTransaction(testOptions(server))
			if err != nil {
				t.Fatalf("NewTransaction: %v", err)
			}
			err = tt.method(tx)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			var tbErr *transbank.HTTPError
			if !errors.As(err, &tbErr) {
				t.Fatalf("error type = %T, want *transbank.HTTPError", err)
			}
			if tbErr.StatusCode != http.StatusUnauthorized {
				t.Errorf("StatusCode = %d, want %d", tbErr.StatusCode, http.StatusUnauthorized)
			}
			if tbErr.Body != `{"error_message":"unauthorized"}` {
				t.Errorf("Body = %q, want raw response", tbErr.Body)
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
