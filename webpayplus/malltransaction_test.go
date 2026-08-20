package webpayplus_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/ppastene/transbank-sdk-go"
	"github.com/ppastene/transbank-sdk-go/webpayplus"
)

func TestMallTransactionCreate(t *testing.T) {
	server := newMockServer(t, http.StatusOK, `{"token":"tok_mall","url":"https://webpay.cl/pagar/tok_mall"}`)
	mt, err := webpayplus.NewMallTransaction(testOptions(server))
	if err != nil {
		t.Fatalf("NewMallTransaction: %v", err)
	}

	details := []webpayplus.MallDetails{
		{Amount: 10000, CommerceCode: testChildCode1, BuyOrder: "orden-detalle-1"},
		{Amount: 2500.50, CommerceCode: testChildCode2, BuyOrder: "orden-detalle-2"},
	}

	create, err := mt.Create(testBuyOrder, testSessionID, testReturnURL, details)
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	if create.Token != "tok_mall" {
		t.Errorf("Token = %q, want tok_mall", create.Token)
	}

	req := server.LastRequest()
	assertRequest(t, req, http.MethodPost, testTransactionsPath)

	var body struct {
		BuyOrder  string                   `json:"buy_order"`
		SessionID string                   `json:"session_id"`
		ReturnURL string                   `json:"return_url"`
		Details   []webpayplus.MallDetails `json:"details"`
	}
	if err := json.Unmarshal([]byte(req.Body), &body); err != nil {
		t.Fatalf("decoding request body: %v", err)
	}
	if body.BuyOrder != testBuyOrder || body.SessionID != testSessionID || body.ReturnURL != testReturnURL {
		t.Errorf("body = %+v, want buy_order=%q session_id=%q return_url=%q", body, testBuyOrder, testSessionID, testReturnURL)
	}
	if len(body.Details) != 2 {
		t.Fatalf("details count = %d, want 2", len(body.Details))
	}
	if body.Details[0].Amount != 10000 || body.Details[0].CommerceCode != testChildCode1 || body.Details[0].BuyOrder != "orden-detalle-1" {
		t.Errorf("details[0] = %+v, unexpected", body.Details[0])
	}
	if body.Details[1].Amount != 2500.5 || body.Details[1].CommerceCode != testChildCode2 {
		t.Errorf("details[1] = %+v, unexpected", body.Details[1])
	}
}

func TestMallTransactionCommit(t *testing.T) {
	server := newMockServer(t, http.StatusOK, `{"details":[{"response_code":0,"status":"AUTHORIZED","commerce_code":"`+testChildCode1+`","buy_order":"orden-detalle-1"}],"buy_order":"`+testBuyOrder+`"}`)
	mt, err := webpayplus.NewMallTransaction(testOptions(server))
	if err != nil {
		t.Fatalf("NewMallTransaction: %v", err)
	}

	commit, err := mt.Commit(testToken)
	if err != nil {
		t.Fatalf("Commit returned error: %v", err)
	}
	if len(commit.Details) != 1 {
		t.Fatalf("details count = %d, want 1", len(commit.Details))
	}
	if !commit.IsApproved() {
		t.Error("IsApproved = false, want true")
	}

	assertRequest(t, server.LastRequest(), http.MethodPut, testTransactionsPath+"/"+testToken)
}

func TestMallTransactionStatus(t *testing.T) {
	server := newMockServer(t, http.StatusOK, `{"details":[{"response_code":0,"status":"AUTHORIZED","amount":10000,"commerce_code":"`+testChildCode1+`","buy_order":"orden-detalle-1"}],"buy_order":"`+testBuyOrder+`"}`)
	mt, err := webpayplus.NewMallTransaction(testOptions(server))
	if err != nil {
		t.Fatalf("NewMallTransaction: %v", err)
	}

	status, err := mt.Status(testToken)
	if err != nil {
		t.Fatalf("Status returned error: %v", err)
	}
	if len(status.Details) != 1 {
		t.Fatalf("details count = %d, want 1", len(status.Details))
	}
	if status.Details[0].Status != "AUTHORIZED" {
		t.Errorf("details[0].Status = %q, want AUTHORIZED", status.Details[0].Status)
	}

	assertRequest(t, server.LastRequest(), http.MethodGet, testTransactionsPath+"/"+testToken)
}

func TestMallTransactionRefund(t *testing.T) {
	server := newMockServer(t, http.StatusOK, `{"type":"REFUNDED","authorization_code":"`+testAuthCode+`","authorization_date":"2026-08-13","nullified_amount":10000,"balance":0,"response_code":0}`)
	mt, err := webpayplus.NewMallTransaction(testOptions(server))
	if err != nil {
		t.Fatalf("NewMallTransaction: %v", err)
	}

	refund, err := mt.Refund(testToken, "orden-detalle-1", testChildCode1, 10000)
	if err != nil {
		t.Fatalf("Refund returned error: %v", err)
	}
	if refund.Type != "REFUNDED" {
		t.Errorf("Type = %q, want REFUNDED", refund.Type)
	}

	req := server.LastRequest()
	assertRequest(t, req, http.MethodPost, testTransactionsPath+"/"+testToken+"/refunds")

	var body map[string]any
	if err := json.Unmarshal([]byte(req.Body), &body); err != nil {
		t.Fatalf("decoding request body: %v", err)
	}
	if body["buy_order"] != "orden-detalle-1" {
		t.Errorf("body buy_order = %v, want orden-detalle-1", body["buy_order"])
	}
	if body["commerce_code"] != testChildCode1 {
		t.Errorf("body commerce_code = %v, want %q", body["commerce_code"], testChildCode1)
	}
	if body["amount"] != 10000.0 {
		t.Errorf("body amount = %v, want 10000", body["amount"])
	}
}

func TestMallTransactionCapture(t *testing.T) {
	server := newMockServer(t, http.StatusOK, `{"authorization_code":"`+testAuthCode+`","authorization_date":"2026-08-13","captured_amount":10000,"response_code":0}`)
	mt, err := webpayplus.NewMallTransaction(testOptions(server))
	if err != nil {
		t.Fatalf("NewMallTransaction: %v", err)
	}

	capture, err := mt.Capture(testToken, testChildCode1, "orden-detalle-1", testAuthCode, 10000)
	if err != nil {
		t.Fatalf("Capture returned error: %v", err)
	}
	if capture.CapturedAmount != 10000 {
		t.Errorf("CapturedAmount = %v, want 10000", capture.CapturedAmount)
	}

	req := server.LastRequest()
	assertRequest(t, req, http.MethodPut, testTransactionsPath+"/"+testToken+"/capture")

	var body map[string]any
	if err := json.Unmarshal([]byte(req.Body), &body); err != nil {
		t.Fatalf("decoding request body: %v", err)
	}
	if body["commerce_code"] != testChildCode1 {
		t.Errorf("body commerce_code = %v, want %q", body["commerce_code"], testChildCode1)
	}
	if body["capture_amount"] != 10000.0 {
		t.Errorf("body capture_amount = %v, want 10000", body["capture_amount"])
	}
}

func TestMallTransactionHTTPError(t *testing.T) {
	tests := []struct {
		name   string
		method func(mt *webpayplus.MallTransaction) error
	}{
		{
			name: "commit",
			method: func(mt *webpayplus.MallTransaction) error {
				_, err := mt.Commit(testToken)
				return err
			},
		},
		{
			name: "status",
			method: func(mt *webpayplus.MallTransaction) error {
				_, err := mt.Status(testToken)
				return err
			},
		},
		{
			name: "refund",
			method: func(mt *webpayplus.MallTransaction) error {
				_, err := mt.Refund(testToken, "orden-detalle-1", testChildCode1, 10000)
				return err
			},
		},
		{
			name: "capture",
			method: func(mt *webpayplus.MallTransaction) error {
				_, err := mt.Capture(testToken, testChildCode1, "orden-detalle-1", testAuthCode, 10000)
				return err
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := newMockServer(t, http.StatusUnauthorized, `{"error_message":"unauthorized"}`)
			mt, err := webpayplus.NewMallTransaction(testOptions(server))
			if err != nil {
				t.Fatalf("NewMallTransaction: %v", err)
			}
			err = tt.method(mt)
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
