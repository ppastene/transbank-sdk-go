package webpayplus_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/ppastene/transbank-sdk-go"
	"github.com/ppastene/transbank-sdk-go/webpayplus"
)

func TestMallTransactionCreateNoValidation(t *testing.T) {
	server := newMockServer(t, http.StatusOK, `{"token":"tok_mall","url":"https://webpay.cl/pagar/tok_mall"}`)
	tx, err := webpayplus.NewMallTransaction(testOptionsNoValidation(server))
	if err != nil {
		t.Fatalf("NewMallTransaction: %v", err)
	}
	details := []webpayplus.MallDetails{{Amount: 10000, CommerceCode: testChildCode1, BuyOrder: "orden-detalle-1"}}
	resp, err := tx.Create(testBuyOrder, testSessionID, testReturnURL, details)
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	if resp.Token != "tok_mall" {
		t.Errorf("Token = %q, want %q", resp.Token, "tok_mall")
	}
	if server.RequestCount() != 1 {
		t.Errorf("request count = %d, want 1", server.RequestCount())
	}
}

func TestMallTransactionCreateNoValidationAPIRejection(t *testing.T) {
	server := newMockServer(t, http.StatusBadRequest, `{"error_message":"invalid parameter"}`)
	tx, err := webpayplus.NewMallTransaction(testOptionsNoValidation(server))
	if err != nil {
		t.Fatalf("NewMallTransaction: %v", err)
	}
	_, err = tx.Create(testBuyOrder, testSessionID, testReturnURL, nil)
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

func TestMallTransactionCommitNoValidation(t *testing.T) {
	server := newMockServer(t, http.StatusOK, `{"details":[{"response_code":0,"status":"AUTHORIZED","commerce_code":"`+testChildCode1+`","buy_order":"orden-detalle-1"}],"buy_order":"`+testBuyOrder+`"}`)
	tx, err := webpayplus.NewMallTransaction(testOptionsNoValidation(server))
	if err != nil {
		t.Fatalf("NewMallTransaction: %v", err)
	}
	resp, err := tx.Commit(testToken)
	if err != nil {
		t.Fatalf("Commit returned error: %v", err)
	}
	if !resp.IsApproved() {
		t.Errorf("IsApproved() = false, want true")
	}
	if server.RequestCount() != 1 {
		t.Errorf("request count = %d, want 1", server.RequestCount())
	}
}

func TestMallTransactionCommitNoValidationAPIRejection(t *testing.T) {
	server := newMockServer(t, http.StatusBadRequest, `{"error_message":"invalid parameter"}`)
	tx, err := webpayplus.NewMallTransaction(testOptionsNoValidation(server))
	if err != nil {
		t.Fatalf("NewMallTransaction: %v", err)
	}
	_, err = tx.Commit("")
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

func TestMallTransactionStatusNoValidation(t *testing.T) {
	server := newMockServer(t, http.StatusOK, `{"details":[{"response_code":0,"status":"AUTHORIZED","amount":10000,"commerce_code":"`+testChildCode1+`","buy_order":"orden-detalle-1"}],"buy_order":"`+testBuyOrder+`"}`)
	tx, err := webpayplus.NewMallTransaction(testOptionsNoValidation(server))
	if err != nil {
		t.Fatalf("NewMallTransaction: %v", err)
	}
	resp, err := tx.Status(testToken)
	if err != nil {
		t.Fatalf("Status returned error: %v", err)
	}
	if len(resp.Details) != 1 {
		t.Fatalf("len(Details) = %d, want 1", len(resp.Details))
	}
	if resp.Details[0].Status != "AUTHORIZED" {
		t.Errorf("Details[0].Status = %q, want %q", resp.Details[0].Status, "AUTHORIZED")
	}
	if server.RequestCount() != 1 {
		t.Errorf("request count = %d, want 1", server.RequestCount())
	}
}

func TestMallTransactionStatusNoValidationAPIRejection(t *testing.T) {
	server := newMockServer(t, http.StatusBadRequest, `{"error_message":"invalid parameter"}`)
	tx, err := webpayplus.NewMallTransaction(testOptionsNoValidation(server))
	if err != nil {
		t.Fatalf("NewMallTransaction: %v", err)
	}
	_, err = tx.Status("")
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

func TestMallTransactionRefundNoValidation(t *testing.T) {
	server := newMockServer(t, http.StatusOK, `{"type":"REFUNDED","authorization_code":"`+testAuthCode+`","authorization_date":"2026-08-13","nullified_amount":10000,"balance":0,"response_code":0}`)
	tx, err := webpayplus.NewMallTransaction(testOptionsNoValidation(server))
	if err != nil {
		t.Fatalf("NewMallTransaction: %v", err)
	}
	resp, err := tx.Refund(testToken, "orden-detalle-1", testChildCode1, 10000)
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

func TestMallTransactionRefundNoValidationAPIRejection(t *testing.T) {
	server := newMockServer(t, http.StatusBadRequest, `{"error_message":"invalid parameter"}`)
	tx, err := webpayplus.NewMallTransaction(testOptionsNoValidation(server))
	if err != nil {
		t.Fatalf("NewMallTransaction: %v", err)
	}
	_, err = tx.Refund(testToken, "orden-detalle-1", testChildCode1, 0)
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

func TestMallTransactionCaptureNoValidation(t *testing.T) {
	server := newMockServer(t, http.StatusOK, `{"authorization_code":"`+testAuthCode+`","authorization_date":"2026-08-13","captured_amount":10000,"response_code":0}`)
	tx, err := webpayplus.NewMallTransaction(testOptionsNoValidation(server))
	if err != nil {
		t.Fatalf("NewMallTransaction: %v", err)
	}
	resp, err := tx.Capture(testToken, testChildCode1, "orden-detalle-1", testAuthCode, 10000)
	if err != nil {
		t.Fatalf("Capture returned error: %v", err)
	}
	if resp.CapturedAmount != 10000 {
		t.Errorf("CapturedAmount = %f, want 10000", resp.CapturedAmount)
	}
	if server.RequestCount() != 1 {
		t.Errorf("request count = %d, want 1", server.RequestCount())
	}
}

func TestMallTransactionCaptureNoValidationAPIRejection(t *testing.T) {
	server := newMockServer(t, http.StatusBadRequest, `{"error_message":"invalid parameter"}`)
	tx, err := webpayplus.NewMallTransaction(testOptionsNoValidation(server))
	if err != nil {
		t.Fatalf("NewMallTransaction: %v", err)
	}
	_, err = tx.Capture(testToken, testChildCode1, "", testAuthCode, 10000)
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
