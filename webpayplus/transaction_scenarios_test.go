package webpayplus_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/ppastene/transbank-sdk-go"
	"github.com/ppastene/transbank-sdk-go/webpayplus"
)

func TestTransactionCreateNoValidation(t *testing.T) {
	server := newMockServer(t, http.StatusOK, `{"token":"tok_123","url":"https://webpay.cl/pagar/tok_123"}`)
	tx, err := webpayplus.NewTransaction(testOptionsNoValidation(server))
	if err != nil {
		t.Fatalf("NewTransaction: %v", err)
	}
	resp, err := tx.Create(testBuyOrder, testSessionID, 10000.50, testReturnURL)
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	if resp.Token != "tok_123" {
		t.Errorf("Token = %q, want %q", resp.Token, "tok_123")
	}
	if server.RequestCount() != 1 {
		t.Errorf("request count = %d, want 1", server.RequestCount())
	}
}

func TestTransactionCreateNoValidationAPIRejection(t *testing.T) {
	server := newMockServer(t, http.StatusBadRequest, `{"error_message":"invalid parameter"}`)
	tx, err := webpayplus.NewTransaction(testOptionsNoValidation(server))
	if err != nil {
		t.Fatalf("NewTransaction: %v", err)
	}
	_, err = tx.Create("", testSessionID, 10000.50, testReturnURL)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var tbErr *transbank.HTTPError
	if !errors.As(err, &tbErr) {
		t.Fatalf("error type = %T, want *transbank.HTTPError", err)
	}
	if tbErr.StatusCode != http.StatusBadRequest {
		t.Errorf("StatusCode = %d, want %d", tbErr.StatusCode, http.StatusBadRequest)
	}
	if server.RequestCount() != 1 {
		t.Errorf("request count = %d, want 1 (API must be called when validation is off)", server.RequestCount())
	}
}

func TestTransactionCommitNoValidation(t *testing.T) {
	server := newMockServer(t, http.StatusOK, `{"response_code":0,"status":"AUTHORIZED","buy_order":"`+testBuyOrder+`"}`)
	tx, err := webpayplus.NewTransaction(testOptionsNoValidation(server))
	if err != nil {
		t.Fatalf("NewTransaction: %v", err)
	}
	resp, err := tx.Commit(testToken)
	if err != nil {
		t.Fatalf("Commit returned error: %v", err)
	}
	if resp.ResponseCode != 0 {
		t.Errorf("ResponseCode = %d, want 0", resp.ResponseCode)
	}
	if server.RequestCount() != 1 {
		t.Errorf("request count = %d, want 1", server.RequestCount())
	}
}

func TestTransactionCommitTokenAlwaysValidated(t *testing.T) {
	server := newMockServer(t, http.StatusOK, `{}`)
	tx, err := webpayplus.NewTransaction(testOptionsNoValidation(server))
	if err != nil {
		t.Fatalf("NewTransaction: %v", err)
	}
	_, err = tx.Commit("")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	wantValidationError(t, err)
	if server.RequestCount() != 0 {
		t.Errorf("request count = %d, want 0 (token validation must not hit API)", server.RequestCount())
	}
}

func TestTransactionStatusNoValidation(t *testing.T) {
	server := newMockServer(t, http.StatusOK, `{"vci":"TSY","amount":10000.5,"status":"AUTHORIZED","buy_order":"`+testBuyOrder+`","session_id":"`+testSessionID+`","card_detail":{"card_number":"6623"},"response_code":0}`)
	tx, err := webpayplus.NewTransaction(testOptionsNoValidation(server))
	if err != nil {
		t.Fatalf("NewTransaction: %v", err)
	}
	resp, err := tx.Status(testToken)
	if err != nil {
		t.Fatalf("Status returned error: %v", err)
	}
	if resp.Status != "AUTHORIZED" {
		t.Errorf("Status = %q, want %q", resp.Status, "AUTHORIZED")
	}
	if server.RequestCount() != 1 {
		t.Errorf("request count = %d, want 1", server.RequestCount())
	}
}

func TestTransactionStatusTokenAlwaysValidated(t *testing.T) {
	server := newMockServer(t, http.StatusOK, `{}`)
	tx, err := webpayplus.NewTransaction(testOptionsNoValidation(server))
	if err != nil {
		t.Fatalf("NewTransaction: %v", err)
	}
	_, err = tx.Status("short-token")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	wantValidationError(t, err)
	if server.RequestCount() != 0 {
		t.Errorf("request count = %d, want 0 (token validation must not hit API)", server.RequestCount())
	}
}

func TestTransactionRefundNoValidation(t *testing.T) {
	server := newMockServer(t, http.StatusOK, `{"type":"REFUNDED","authorization_code":"`+testAuthCode+`","authorization_date":"2026-08-13","nullified_amount":10000.5,"balance":0,"response_code":0}`)
	tx, err := webpayplus.NewTransaction(testOptionsNoValidation(server))
	if err != nil {
		t.Fatalf("NewTransaction: %v", err)
	}
	resp, err := tx.Refund(testToken, 10000.50)
	if err != nil {
		t.Fatalf("Refund returned error: %v", err)
	}
	if resp.Type != "REFUNDED" {
		t.Errorf("Type = %q, want %q", resp.Type, "REFUNDED")
	}
	if server.RequestCount() != 1 {
		t.Errorf("request count = %d, want 1", server.RequestCount())
	}
}

func TestTransactionRefundNoValidationAPIRejection(t *testing.T) {
	server := newMockServer(t, http.StatusBadRequest, `{"error_message":"invalid parameter"}`)
	tx, err := webpayplus.NewTransaction(testOptionsNoValidation(server))
	if err != nil {
		t.Fatalf("NewTransaction: %v", err)
	}
	_, err = tx.Refund(testToken, 0)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var tbErr *transbank.HTTPError
	if !errors.As(err, &tbErr) {
		t.Fatalf("error type = %T, want *transbank.HTTPError", err)
	}
	if tbErr.StatusCode != http.StatusBadRequest {
		t.Errorf("StatusCode = %d, want %d", tbErr.StatusCode, http.StatusBadRequest)
	}
	if server.RequestCount() != 1 {
		t.Errorf("request count = %d, want 1 (API must be called when validation is off)", server.RequestCount())
	}
}

func TestTransactionCaptureNoValidation(t *testing.T) {
	server := newMockServer(t, http.StatusOK, `{"authorization_code":"`+testAuthCode+`","authorization_date":"2026-08-13","captured_amount":10000.5,"response_code":0}`)
	tx, err := webpayplus.NewTransaction(testOptionsNoValidation(server))
	if err != nil {
		t.Fatalf("NewTransaction: %v", err)
	}
	resp, err := tx.Capture(testToken, testBuyOrder, testAuthCode, 10000.50)
	if err != nil {
		t.Fatalf("Capture returned error: %v", err)
	}
	if resp.CapturedAmount != 10000.5 {
		t.Errorf("CapturedAmount = %f, want 10000.5", resp.CapturedAmount)
	}
	if server.RequestCount() != 1 {
		t.Errorf("request count = %d, want 1", server.RequestCount())
	}
}

func TestTransactionCaptureNoValidationAPIRejection(t *testing.T) {
	server := newMockServer(t, http.StatusBadRequest, `{"error_message":"invalid parameter"}`)
	tx, err := webpayplus.NewTransaction(testOptionsNoValidation(server))
	if err != nil {
		t.Fatalf("NewTransaction: %v", err)
	}
	_, err = tx.Capture(testToken, testBuyOrder, "", 10000.50)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var tbErr *transbank.HTTPError
	if !errors.As(err, &tbErr) {
		t.Fatalf("error type = %T, want *transbank.HTTPError", err)
	}
	if tbErr.StatusCode != http.StatusBadRequest {
		t.Errorf("StatusCode = %d, want %d", tbErr.StatusCode, http.StatusBadRequest)
	}
	if server.RequestCount() != 1 {
		t.Errorf("request count = %d, want 1 (API must be called when validation is off)", server.RequestCount())
	}
}
