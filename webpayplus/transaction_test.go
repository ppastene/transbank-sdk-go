package webpayplus_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/ppastene/transbank-sdk-go"
	"github.com/ppastene/transbank-sdk-go/webpayplus"
)

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
			name: "empty commerce code",
			opts: transbank.Options{
				ApiKey:      testAPIKey,
				Environment: transbank.Integration,
			},
			want: &transbank.ValidationError{},
		},
		{
			name: "short commerce code",
			opts: transbank.Options{
				CommerceCode: "5970",
				ApiKey:       testAPIKey,
				Environment:  transbank.Integration,
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
			tx, err := webpayplus.NewTransaction(tt.opts)
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

func TestTransactionCreate(t *testing.T) {
	server := newMockServer(t, http.StatusOK, `{"token":"tok_123","url":"https://webpay.cl/pagar/tok_123"}`)
	tx, err := webpayplus.NewTransaction(testOptions(server))
	if err != nil {
		t.Fatalf("NewTransaction: %v", err)
	}

	create, err := tx.Create(testBuyOrder, testSessionID, 10000.50, testReturnURL)
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	if create.Token != "tok_123" {
		t.Errorf("Token = %q, want tok_123", create.Token)
	}
	if create.Url != "https://webpay.cl/pagar/tok_123" {
		t.Errorf("Url = %q, want https://webpay.cl/pagar/tok_123", create.Url)
	}

	req := server.LastRequest()
	assertRequest(t, req, http.MethodPost, testTransactionsPath)

	var body map[string]any
	if err := json.Unmarshal([]byte(req.Body), &body); err != nil {
		t.Fatalf("decoding request body: %v", err)
	}
	if body["buy_order"] != testBuyOrder {
		t.Errorf("body buy_order = %v, want %q", body["buy_order"], testBuyOrder)
	}
	if body["session_id"] != testSessionID {
		t.Errorf("body session_id = %v, want %q", body["session_id"], testSessionID)
	}
	if body["amount"] != 10000.5 {
		t.Errorf("body amount = %v, want 10000.5", body["amount"])
	}
	if body["return_url"] != testReturnURL {
		t.Errorf("body return_url = %v, want %q", body["return_url"], testReturnURL)
	}
}

func TestTransactionCommit(t *testing.T) {
	server := newMockServer(t, http.StatusOK, `{"response_code":0,"status":"AUTHORIZED","buy_order":"`+testBuyOrder+`"}`)
	tx, err := webpayplus.NewTransaction(testOptions(server))
	if err != nil {
		t.Fatalf("NewTransaction: %v", err)
	}

	commit, err := tx.Commit(testToken)
	if err != nil {
		t.Fatalf("Commit returned error: %v", err)
	}
	if commit.ResponseCode != 0 {
		t.Errorf("ResponseCode = %d, want 0", commit.ResponseCode)
	}
	if !commit.IsApproved() {
		t.Error("IsApproved = false, want true")
	}

	assertRequest(t, server.LastRequest(), http.MethodPut, testTransactionsPath+"/"+testToken)
}

func TestTransactionStatus(t *testing.T) {
	server := newMockServer(t, http.StatusOK, `{"vci":"TSY","amount":10000.5,"status":"AUTHORIZED","buy_order":"`+testBuyOrder+`","session_id":"`+testSessionID+`","card_detail":{"card_number":"6623"},"response_code":0}`)
	tx, err := webpayplus.NewTransaction(testOptions(server))
	if err != nil {
		t.Fatalf("NewTransaction: %v", err)
	}

	status, err := tx.Status(testToken)
	if err != nil {
		t.Fatalf("Status returned error: %v", err)
	}
	if status.Status != "AUTHORIZED" {
		t.Errorf("Status = %q, want AUTHORIZED", status.Status)
	}
	if status.Amount != 10000.5 {
		t.Errorf("Amount = %v, want 10000.5", status.Amount)
	}
	if status.CardDetail.CardNumber != "6623" {
		t.Errorf("CardNumber = %q, want 6623", status.CardDetail.CardNumber)
	}

	assertRequest(t, server.LastRequest(), http.MethodGet, testTransactionsPath+"/"+testToken)
}

func TestTransactionRefund(t *testing.T) {
	server := newMockServer(t, http.StatusOK, `{"type":"REFUNDED","authorization_code":"`+testAuthCode+`","authorization_date":"2026-08-13","nullified_amount":10000.5,"balance":0,"response_code":0}`)
	tx, err := webpayplus.NewTransaction(testOptions(server))
	if err != nil {
		t.Fatalf("NewTransaction: %v", err)
	}

	refund, err := tx.Refund(testToken, 10000.50)
	if err != nil {
		t.Fatalf("Refund returned error: %v", err)
	}
	if refund.Type != "REFUNDED" {
		t.Errorf("Type = %q, want REFUNDED", refund.Type)
	}
	if refund.NullifiedAmount != 10000.5 {
		t.Errorf("NullifiedAmount = %v, want 10000.5", refund.NullifiedAmount)
	}

	req := server.LastRequest()
	assertRequest(t, req, http.MethodPost, testTransactionsPath+"/"+testToken+"/refunds")

	var body map[string]any
	if err := json.Unmarshal([]byte(req.Body), &body); err != nil {
		t.Fatalf("decoding request body: %v", err)
	}
	if body["amount"] != 10000.5 {
		t.Errorf("body amount = %v, want 10000.5", body["amount"])
	}
}

func TestTransactionCapture(t *testing.T) {
	server := newMockServer(t, http.StatusOK, `{"authorization_code":"`+testAuthCode+`","authorization_date":"2026-08-13","captured_amount":10000.5,"response_code":0}`)
	tx, err := webpayplus.NewTransaction(testOptions(server))
	if err != nil {
		t.Fatalf("NewTransaction: %v", err)
	}

	capture, err := tx.Capture(testToken, testBuyOrder, testAuthCode, 10000.50)
	if err != nil {
		t.Fatalf("Capture returned error: %v", err)
	}
	if capture.CapturedAmount != 10000.5 {
		t.Errorf("CapturedAmount = %v, want 10000.5", capture.CapturedAmount)
	}

	req := server.LastRequest()
	assertRequest(t, req, http.MethodPut, testTransactionsPath+"/"+testToken+"/capture")

	var body map[string]any
	if err := json.Unmarshal([]byte(req.Body), &body); err != nil {
		t.Fatalf("decoding request body: %v", err)
	}
	if body["buy_order"] != testBuyOrder {
		t.Errorf("body buy_order = %v, want %q", body["buy_order"], testBuyOrder)
	}
	if body["authorization_code"] != testAuthCode {
		t.Errorf("body authorization_code = %v, want %q", body["authorization_code"], testAuthCode)
	}
	if body["capture_amount"] != 10000.5 {
		t.Errorf("body capture_amount = %v, want 10000.5", body["capture_amount"])
	}
}

func TestTransactionHTTPError(t *testing.T) {
	server := newMockServer(t, http.StatusUnauthorized, `{"error_message":"unauthorized"}`)
	tx, err := webpayplus.NewTransaction(testOptions(server))
	if err != nil {
		t.Fatalf("NewTransaction: %v", err)
	}

	_, err = tx.Status(testToken)
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
}

func TestTransactionValidationSkipsRequest(t *testing.T) {
	server := newMockServer(t, http.StatusOK, `{}`)
	tx, err := webpayplus.NewTransaction(testOptions(server))
	if err != nil {
		t.Fatalf("NewTransaction: %v", err)
	}

	if _, err := tx.Create("", testSessionID, 10000, testReturnURL); err == nil {
		t.Error("expected error for empty buy_order")
	}
	if _, err := tx.Create(testBuyOrder, testSessionID, 0, testReturnURL); err == nil {
		t.Error("expected error for zero amount")
	}
	if _, err := tx.Create(testBuyOrder, testSessionID, 10000.123, testReturnURL); err == nil {
		t.Error("expected error for amount with more than 2 decimals")
	}
	if _, err := tx.Create(testBuyOrder, testSessionID, 10000, "not-a-url"); err == nil {
		t.Error("expected error for relative return_url")
	}
	if _, err := tx.Status("short-token"); err == nil {
		t.Error("expected error for invalid token")
	}
	if _, err := tx.Refund(testToken, 0); err == nil {
		t.Error("expected error for zero refund amount")
	}

	if got := server.RequestCount(); got != 0 {
		t.Errorf("request count = %d, want 0 (validation must not hit the API)", got)
	}
}
