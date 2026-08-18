package oneclick_test

import (
	"net/http"
	"testing"

	"github.com/ppastene/transbank-sdk-go/oneclick"
)

func TestMallTransactionAuthorizeNoValidation(t *testing.T) {
	server := newMockServer(t, http.StatusOK, `{"buy_order":"`+testBuyOrder+`","card_detail":{"card_number":"6623"},"accounting_date":"2026-08-13","transaction_date":"2026-08-13T10:00:00.000Z","details":[{"amount":10000,"status":"AUTHORIZED","authorization_code":"`+testAuthCode+`","payment_type_code":"VN","response_code":0,"installments_number":3,"commerce_code":"`+testChildCode1+`","buy_order":"`+testChildBuyOrder+`"}]}`)
	tx, err := oneclick.NewMallTransaction(testOptionsNoValidation(server))
	if err != nil {
		t.Fatalf("NewMallTransaction: %v", err)
	}

	auth, err := tx.Authorize(testUsername, testTbkUser, testBuyOrder, testDetails())
	if err != nil {
		t.Fatalf("Authorize returned error: %v", err)
	}
	if auth.BuyOrder != testBuyOrder {
		t.Errorf("BuyOrder = %q, want %q", auth.BuyOrder, testBuyOrder)
	}
	if server.RequestCount() != 1 {
		t.Errorf("request count = %d, want 1", server.RequestCount())
	}
}

func TestMallTransactionAuthorizeNoValidationAPIRejection(t *testing.T) {
	server := newMockServer(t, http.StatusBadRequest, `{"error_message":"bad request"}`)
	tx, err := oneclick.NewMallTransaction(testOptionsNoValidation(server))
	if err != nil {
		t.Fatalf("NewMallTransaction: %v", err)
	}

	_, err = tx.Authorize("", testTbkUser, testBuyOrder, testDetails())
	wantHTTPError(t, err, http.StatusBadRequest, `{"error_message":"bad request"}`)
	if server.RequestCount() != 1 {
		t.Errorf("request count = %d, want 1 (API must be called when validation is off)", server.RequestCount())
	}
}

func TestMallTransactionStatusNoValidation(t *testing.T) {
	server := newMockServer(t, http.StatusOK, `{"buy_order":"`+testBuyOrder+`","card_detail":{"card_number":"6623"},"accounting_date":"2026-08-13","transaction_date":"2026-08-13T10:00:00.000Z","details":[{"amount":10000,"status":"AUTHORIZED","authorization_code":"`+testAuthCode+`","payment_type_code":"VN","response_code":0,"installments_number":3,"commerce_code":"`+testChildCode1+`","buy_order":"`+testChildBuyOrder+`"}]}`)
	tx, err := oneclick.NewMallTransaction(testOptionsNoValidation(server))
	if err != nil {
		t.Fatalf("NewMallTransaction: %v", err)
	}

	status, err := tx.Status(testBuyOrder)
	if err != nil {
		t.Fatalf("Status returned error: %v", err)
	}
	if status.BuyOrder != testBuyOrder {
		t.Errorf("BuyOrder = %q, want %q", status.BuyOrder, testBuyOrder)
	}
	if server.RequestCount() != 1 {
		t.Errorf("request count = %d, want 1", server.RequestCount())
	}
}

func TestMallTransactionStatusNoValidationAPIRejection(t *testing.T) {
	server := newMockServer(t, http.StatusBadRequest, `{"error_message":"bad request"}`)
	tx, err := oneclick.NewMallTransaction(testOptionsNoValidation(server))
	if err != nil {
		t.Fatalf("NewMallTransaction: %v", err)
	}

	_, err = tx.Status("")
	wantHTTPError(t, err, http.StatusBadRequest, `{"error_message":"bad request"}`)
	if server.RequestCount() != 1 {
		t.Errorf("request count = %d, want 1 (API must be called when validation is off)", server.RequestCount())
	}
}

func TestMallTransactionRefundNoValidation(t *testing.T) {
	server := newMockServer(t, http.StatusOK, `{"type":"REFUNDED","authorization_code":"`+testAuthCode+`","authorization_date":"2026-08-13","nullified_amount":10000,"balance":0,"response_code":0}`)
	tx, err := oneclick.NewMallTransaction(testOptionsNoValidation(server))
	if err != nil {
		t.Fatalf("NewMallTransaction: %v", err)
	}

	refund, err := tx.Refund(testBuyOrder, testChildCode1, testChildBuyOrder, 10000)
	if err != nil {
		t.Fatalf("Refund returned error: %v", err)
	}
	if refund.Type != "REFUNDED" {
		t.Errorf("Type = %q, want REFUNDED", refund.Type)
	}
	if server.RequestCount() != 1 {
		t.Errorf("request count = %d, want 1", server.RequestCount())
	}
}

func TestMallTransactionRefundNoValidationAPIRejection(t *testing.T) {
	server := newMockServer(t, http.StatusBadRequest, `{"error_message":"bad request"}`)
	tx, err := oneclick.NewMallTransaction(testOptionsNoValidation(server))
	if err != nil {
		t.Fatalf("NewMallTransaction: %v", err)
	}

	_, err = tx.Refund(testBuyOrder, testChildCode1, testChildBuyOrder, 0)
	wantHTTPError(t, err, http.StatusBadRequest, `{"error_message":"bad request"}`)
	if server.RequestCount() != 1 {
		t.Errorf("request count = %d, want 1 (API must be called when validation is off)", server.RequestCount())
	}
}

func TestMallTransactionCaptureNoValidation(t *testing.T) {
	server := newMockServer(t, http.StatusOK, `{"authorization_code":"`+testAuthCode+`","authorization_date":"2026-08-13T10:00:00.000Z","captured_amount":50,"response_code":0}`)
	tx, err := oneclick.NewMallTransaction(testOptionsNoValidation(server))
	if err != nil {
		t.Fatalf("NewMallTransaction: %v", err)
	}

	capture, err := tx.Capture(testBuyOrder, testCaptureCode, testAuthCode, 50)
	if err != nil {
		t.Fatalf("Capture returned error: %v", err)
	}
	if capture.CapturedAmount != 50 {
		t.Errorf("CapturedAmount = %v, want 50", capture.CapturedAmount)
	}
	if server.RequestCount() != 1 {
		t.Errorf("request count = %d, want 1", server.RequestCount())
	}
}

func TestMallTransactionCaptureNoValidationAPIRejection(t *testing.T) {
	server := newMockServer(t, http.StatusBadRequest, `{"error_message":"bad request"}`)
	tx, err := oneclick.NewMallTransaction(testOptionsNoValidation(server))
	if err != nil {
		t.Fatalf("NewMallTransaction: %v", err)
	}

	_, err = tx.Capture(testBuyOrder, testCaptureCode, "", 50)
	wantHTTPError(t, err, http.StatusBadRequest, `{"error_message":"bad request"}`)
	if server.RequestCount() != 1 {
		t.Errorf("request count = %d, want 1 (API must be called when validation is off)", server.RequestCount())
	}
}
