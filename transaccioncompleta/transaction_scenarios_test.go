package transaccioncompleta_test

import (
	"net/http"
	"testing"

	"github.com/ppastene/transbank-sdk-go/transaccioncompleta"
)

func TestTransactionCreateNoValidation(t *testing.T) {
	server := newMockServer(t, http.StatusOK, `{"token":"`+testToken+`"}`)
	tx, err := transaccioncompleta.NewTransaction(testOptionsNoValidation(server))
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
	if server.RequestCount() != 1 {
		t.Errorf("request count = %d, want 1", server.RequestCount())
	}
}

func TestTransactionCreateNoValidationAPIRejection(t *testing.T) {
	server := newMockServer(t, http.StatusBadRequest, `{"error_message":"invalid parameter"}`)
	tx, err := transaccioncompleta.NewTransaction(testOptionsNoValidation(server))
	if err != nil {
		t.Fatalf("NewTransaction: %v", err)
	}
	_, err = tx.Create("", testSessionID, testAmount, testCardNumber, testCardExpiry, testCVV)
	wantHTTPError(t, err, http.StatusBadRequest, `{"error_message":"invalid parameter"}`)
	if server.RequestCount() != 1 {
		t.Errorf("request count = %d, want 1 (API must be called when validation is off)", server.RequestCount())
	}
}

func TestTransactionInstallmentsNoValidation(t *testing.T) {
	server := newMockServer(t, http.StatusOK, `{"installments_amount":3334,"id_query_installments":11,"deferred_periods":[{"amount":1000,"period":1}]}`)
	tx, err := transaccioncompleta.NewTransaction(testOptionsNoValidation(server))
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
	if server.RequestCount() != 1 {
		t.Errorf("request count = %d, want 1", server.RequestCount())
	}
}

func TestTransactionInstallmentsNoValidationAPIRejection(t *testing.T) {
	server := newMockServer(t, http.StatusBadRequest, `{"error_message":"invalid parameter"}`)
	tx, err := transaccioncompleta.NewTransaction(testOptionsNoValidation(server))
	if err != nil {
		t.Fatalf("NewTransaction: %v", err)
	}
	_, err = tx.Installments("", 10)
	wantHTTPError(t, err, http.StatusBadRequest, `{"error_message":"invalid parameter"}`)
	if server.RequestCount() != 1 {
		t.Errorf("request count = %d, want 1 (API must be called when validation is off)", server.RequestCount())
	}
}

func TestTransactionCommitNoValidation(t *testing.T) {
	server := newMockServer(t, http.StatusOK, `{"amount":10000,"status":"AUTHORIZED","buy_order":"orden-compra-123","response_code":0}`)
	tx, err := transaccioncompleta.NewTransaction(testOptionsNoValidation(server))
	if err != nil {
		t.Fatalf("NewTransaction: %v", err)
	}
	resp, err := tx.Commit(testToken, nil, nil, nil)
	if err != nil {
		t.Fatalf("Commit returned error: %v", err)
	}
	if resp.Status != "AUTHORIZED" {
		t.Errorf("Status = %q, want %q", resp.Status, "AUTHORIZED")
	}
	if server.RequestCount() != 1 {
		t.Errorf("request count = %d, want 1", server.RequestCount())
	}
}

func TestTransactionCommitNoValidationAPIRejection(t *testing.T) {
	server := newMockServer(t, http.StatusBadRequest, `{"error_message":"invalid parameter"}`)
	tx, err := transaccioncompleta.NewTransaction(testOptionsNoValidation(server))
	if err != nil {
		t.Fatalf("NewTransaction: %v", err)
	}
	_, err = tx.Commit("", nil, nil, nil)
	wantHTTPError(t, err, http.StatusBadRequest, `{"error_message":"invalid parameter"}`)
	if server.RequestCount() != 1 {
		t.Errorf("request count = %d, want 1 (API must be called when validation is off)", server.RequestCount())
	}
}

func TestTransactionStatusNoValidation(t *testing.T) {
	server := newMockServer(t, http.StatusOK, `{"amount":10000,"status":"AUTHORIZED","buy_order":"orden-compra-123","card_detail":{"card_number":"4239"},"response_code":0}`)
	tx, err := transaccioncompleta.NewTransaction(testOptionsNoValidation(server))
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
	if server.RequestCount() != 1 {
		t.Errorf("request count = %d, want 1", server.RequestCount())
	}
}

func TestTransactionStatusNoValidationAPIRejection(t *testing.T) {
	server := newMockServer(t, http.StatusBadRequest, `{"error_message":"invalid parameter"}`)
	tx, err := transaccioncompleta.NewTransaction(testOptionsNoValidation(server))
	if err != nil {
		t.Fatalf("NewTransaction: %v", err)
	}
	_, err = tx.Status("")
	wantHTTPError(t, err, http.StatusBadRequest, `{"error_message":"invalid parameter"}`)
	if server.RequestCount() != 1 {
		t.Errorf("request count = %d, want 1 (API must be called when validation is off)", server.RequestCount())
	}
}

func TestTransactionRefundNoValidation(t *testing.T) {
	server := newMockServer(t, http.StatusOK, `{"type":"NULLIFY","authorization_code":"123456","authorization_date":"2019-03-20T20:18:20Z","nullified_amount":1000.00,"balance":0.00,"response_code":0}`)
	tx, err := transaccioncompleta.NewTransaction(testOptionsNoValidation(server))
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
	if server.RequestCount() != 1 {
		t.Errorf("request count = %d, want 1", server.RequestCount())
	}
}

func TestTransactionRefundNoValidationAPIRejection(t *testing.T) {
	server := newMockServer(t, http.StatusBadRequest, `{"error_message":"invalid parameter"}`)
	tx, err := transaccioncompleta.NewTransaction(testOptionsNoValidation(server))
	if err != nil {
		t.Fatalf("NewTransaction: %v", err)
	}
	_, err = tx.Refund("", 1000)
	wantHTTPError(t, err, http.StatusBadRequest, `{"error_message":"invalid parameter"}`)
	if server.RequestCount() != 1 {
		t.Errorf("request count = %d, want 1 (API must be called when validation is off)", server.RequestCount())
	}
}

func TestTransactionCaptureNoValidation(t *testing.T) {
	server := newMockServer(t, http.StatusOK, `{"authorization_code":"123456","authorization_date":"2019-03-20T20:18:20Z","captured_amount":1000,"response_code":0}`)
	tx, err := transaccioncompleta.NewTransaction(testOptionsNoValidation(server))
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
	if server.RequestCount() != 1 {
		t.Errorf("request count = %d, want 1", server.RequestCount())
	}
}

func TestTransactionCaptureNoValidationAPIRejection(t *testing.T) {
	server := newMockServer(t, http.StatusBadRequest, `{"error_message":"invalid parameter"}`)
	tx, err := transaccioncompleta.NewTransaction(testOptionsNoValidation(server))
	if err != nil {
		t.Fatalf("NewTransaction: %v", err)
	}
	_, err = tx.Capture("", testBuyOrder, testAuthCode, 1000)
	wantHTTPError(t, err, http.StatusBadRequest, `{"error_message":"invalid parameter"}`)
	if server.RequestCount() != 1 {
		t.Errorf("request count = %d, want 1 (API must be called when validation is off)", server.RequestCount())
	}
}
