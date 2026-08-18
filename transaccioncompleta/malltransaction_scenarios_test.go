package transaccioncompleta_test

import (
	"net/http"
	"testing"

	"github.com/ppastene/transbank-sdk-go/transaccioncompleta"
)

func TestMallTransactionCreateNoValidation(t *testing.T) {
	server := newMockServer(t, http.StatusOK, `{"token":"`+testToken+`"}`)
	tx, err := transaccioncompleta.NewMallTransaction(testMallOptionsNoValidation(server))
	if err != nil {
		t.Fatalf("NewMallTransaction: %v", err)
	}
	resp, err := tx.Create(testBuyOrder, testSessionID, testCardNumber, testCardExpiry, testMallDetails(), testCVV)
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

func TestMallTransactionCreateNoValidationAPIRejection(t *testing.T) {
	server := newMockServer(t, http.StatusBadRequest, `{"error_message":"invalid parameter"}`)
	tx, err := transaccioncompleta.NewMallTransaction(testMallOptionsNoValidation(server))
	if err != nil {
		t.Fatalf("NewMallTransaction: %v", err)
	}
	_, err = tx.Create("", testSessionID, testCardNumber, testCardExpiry, testMallDetails(), testCVV)
	wantHTTPError(t, err, http.StatusBadRequest, `{"error_message":"invalid parameter"}`)
	if server.RequestCount() != 1 {
		t.Errorf("request count = %d, want 1 (API must be called when validation is off)", server.RequestCount())
	}
}

func TestMallTransactionInstallmentsNoValidation(t *testing.T) {
	server := newMockServer(t, http.StatusOK, `[{"installments_amount":3334,"id_query_installments":11,"deferred_periods":[{"amount":1000,"period":1}]},{"installments_amount":6667,"id_query_installments":12,"deferred_periods":[]}]`)
	tx, err := transaccioncompleta.NewMallTransaction(testMallOptionsNoValidation(server))
	if err != nil {
		t.Fatalf("NewMallTransaction: %v", err)
	}
	resp, err := tx.Installments(testToken, testMallInstallmentsDetails())
	if err != nil {
		t.Fatalf("Installments returned error: %v", err)
	}
	if len(*resp) != 2 {
		t.Fatalf("len(resp) = %d, want 2", len(*resp))
	}
	if server.RequestCount() != 1 {
		t.Errorf("request count = %d, want 1", server.RequestCount())
	}
}

func TestMallTransactionInstallmentsNoValidationAPIRejection(t *testing.T) {
	server := newMockServer(t, http.StatusBadRequest, `{"error_message":"invalid parameter"}`)
	tx, err := transaccioncompleta.NewMallTransaction(testMallOptionsNoValidation(server))
	if err != nil {
		t.Fatalf("NewMallTransaction: %v", err)
	}
	_, err = tx.Installments("", testMallInstallmentsDetails())
	wantHTTPError(t, err, http.StatusBadRequest, `{"error_message":"invalid parameter"}`)
	if server.RequestCount() != 1 {
		t.Errorf("request count = %d, want 1 (API must be called when validation is off)", server.RequestCount())
	}
}

func TestMallTransactionCommitNoValidation(t *testing.T) {
	server := newMockServer(t, http.StatusOK, `{"buy_order":"orden-compra-123","card_detail":{"card_number":"4239"},"accounting_date":"0321","transaction_date":"2019-03-21T15:43:48.523Z","details":[{"amount":10000,"status":"AUTHORIZED","authorization_code":"123456","payment_type_code":"VN","response_code":0,"installments_number":0,"commerce_code":"597055555552","buy_order":"orden-hija-1"}]}`)
	tx, err := transaccioncompleta.NewMallTransaction(testMallOptionsNoValidation(server))
	if err != nil {
		t.Fatalf("NewMallTransaction: %v", err)
	}
	resp, err := tx.Commit(testToken, testMallCommitDetails())
	if err != nil {
		t.Fatalf("Commit returned error: %v", err)
	}
	if !resp.IsApproved() {
		t.Error("IsApproved = false, want true")
	}
	if server.RequestCount() != 1 {
		t.Errorf("request count = %d, want 1", server.RequestCount())
	}
}

func TestMallTransactionCommitNoValidationAPIRejection(t *testing.T) {
	server := newMockServer(t, http.StatusBadRequest, `{"error_message":"invalid parameter"}`)
	tx, err := transaccioncompleta.NewMallTransaction(testMallOptionsNoValidation(server))
	if err != nil {
		t.Fatalf("NewMallTransaction: %v", err)
	}
	_, err = tx.Commit("", testMallCommitDetails())
	wantHTTPError(t, err, http.StatusBadRequest, `{"error_message":"invalid parameter"}`)
	if server.RequestCount() != 1 {
		t.Errorf("request count = %d, want 1 (API must be called when validation is off)", server.RequestCount())
	}
}

func TestMallTransactionStatusNoValidation(t *testing.T) {
	server := newMockServer(t, http.StatusOK, `{"buy_order":"orden-compra-123","card_detail":{"card_number":"4239"},"accounting_date":"0321","transaction_date":"2019-03-21T15:43:48.523Z","details":[{"amount":10000,"status":"AUTHORIZED","authorization_code":"123456","payment_type_code":"VN","response_code":0,"installments_number":0,"commerce_code":"597055555552","buy_order":"orden-hija-1"}]}`)
	tx, err := transaccioncompleta.NewMallTransaction(testMallOptionsNoValidation(server))
	if err != nil {
		t.Fatalf("NewMallTransaction: %v", err)
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

func TestMallTransactionStatusNoValidationAPIRejection(t *testing.T) {
	server := newMockServer(t, http.StatusBadRequest, `{"error_message":"invalid parameter"}`)
	tx, err := transaccioncompleta.NewMallTransaction(testMallOptionsNoValidation(server))
	if err != nil {
		t.Fatalf("NewMallTransaction: %v", err)
	}
	_, err = tx.Status("")
	wantHTTPError(t, err, http.StatusBadRequest, `{"error_message":"invalid parameter"}`)
	if server.RequestCount() != 1 {
		t.Errorf("request count = %d, want 1 (API must be called when validation is off)", server.RequestCount())
	}
}

func TestMallTransactionRefundNoValidation(t *testing.T) {
	server := newMockServer(t, http.StatusOK, `{"type":"NULLIFY","authorization_code":"123456","authorization_date":"2019-03-20T20:18:20Z","nullified_amount":1000.00,"balance":0.00,"response_code":0}`)
	tx, err := transaccioncompleta.NewMallTransaction(testMallOptionsNoValidation(server))
	if err != nil {
		t.Fatalf("NewMallTransaction: %v", err)
	}
	resp, err := tx.Refund(testToken, testChildBuyOrder1, testChildCode1, 1000)
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

func TestMallTransactionRefundNoValidationAPIRejection(t *testing.T) {
	server := newMockServer(t, http.StatusBadRequest, `{"error_message":"invalid parameter"}`)
	tx, err := transaccioncompleta.NewMallTransaction(testMallOptionsNoValidation(server))
	if err != nil {
		t.Fatalf("NewMallTransaction: %v", err)
	}
	_, err = tx.Refund("", testChildBuyOrder1, testChildCode1, 1000)
	wantHTTPError(t, err, http.StatusBadRequest, `{"error_message":"invalid parameter"}`)
	if server.RequestCount() != 1 {
		t.Errorf("request count = %d, want 1 (API must be called when validation is off)", server.RequestCount())
	}
}

func TestMallTransactionCaptureNoValidation(t *testing.T) {
	server := newMockServer(t, http.StatusOK, `{"authorization_code":"123456","authorization_date":"2019-03-20T20:18:20Z","captured_amount":1000,"response_code":0}`)
	tx, err := transaccioncompleta.NewMallTransaction(testMallOptionsNoValidation(server))
	if err != nil {
		t.Fatalf("NewMallTransaction: %v", err)
	}
	resp, err := tx.Capture(testToken, testChildCode1, testBuyOrder, testAuthCode, 1000)
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

func TestMallTransactionCaptureNoValidationAPIRejection(t *testing.T) {
	server := newMockServer(t, http.StatusBadRequest, `{"error_message":"invalid parameter"}`)
	tx, err := transaccioncompleta.NewMallTransaction(testMallOptionsNoValidation(server))
	if err != nil {
		t.Fatalf("NewMallTransaction: %v", err)
	}
	_, err = tx.Capture("", testChildCode1, testBuyOrder, testAuthCode, 1000)
	wantHTTPError(t, err, http.StatusBadRequest, `{"error_message":"invalid parameter"}`)
	if server.RequestCount() != 1 {
		t.Errorf("request count = %d, want 1 (API must be called when validation is off)", server.RequestCount())
	}
}
