package oneclick_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/ppastene/transbank-sdk-go"
	"github.com/ppastene/transbank-sdk-go/oneclick"
)

func TestNewMallTransaction(t *testing.T) {
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
			tx, err := oneclick.NewMallTransaction(tt.opts)
			if tt.want != nil {
				wantValidationError(t, err)
				if tx != nil {
					t.Error("expected nil MallTransaction on error")
				}
				return
			}
			if err != nil {
				t.Fatalf("NewMallTransaction returned error: %v", err)
			}
			if tx == nil {
				t.Error("expected non-nil MallTransaction")
			}
		})
	}
}

func testDetails() []oneclick.MallDetails {
	return []oneclick.MallDetails{
		{
			Amount:             10000,
			CommerceCode:       testChildCode1,
			BuyOrder:           testChildBuyOrder,
			InstallmentsNumber: 3,
		},
		{
			Amount:             50000,
			CommerceCode:       testChildCode2,
			BuyOrder:           "orden-hija-789",
			InstallmentsNumber: 3,
		},
	}
}

func TestMallTransactionAuthorize(t *testing.T) {
	server := newMockServer(t, http.StatusOK, `{"buy_order":"`+testBuyOrder+`","card_detail":{"card_number":"6623"},"accounting_date":"2026-08-13","transaction_date":"2026-08-13T10:00:00.000Z","details":[{"amount":10000,"status":"AUTHORIZED","authorization_code":"`+testAuthCode+`","payment_type_code":"VN","response_code":0,"installments_number":3,"commerce_code":"`+testChildCode1+`","buy_order":"`+testChildBuyOrder+`"}]}`)
	tx, err := oneclick.NewMallTransaction(testOptions(server))
	if err != nil {
		t.Fatalf("NewMallTransaction: %v", err)
	}

	details := testDetails()
	auth, err := tx.Authorize(testUsername, testTbkUser, testBuyOrder, details)
	if err != nil {
		t.Fatalf("Authorize returned error: %v", err)
	}
	if auth.BuyOrder != testBuyOrder {
		t.Errorf("BuyOrder = %q, want %q", auth.BuyOrder, testBuyOrder)
	}
	if len(auth.Details) != 1 {
		t.Fatalf("len(Details) = %d, want 1", len(auth.Details))
	}
	if auth.Details[0].Status != "AUTHORIZED" {
		t.Errorf("Details[0].Status = %q, want AUTHORIZED", auth.Details[0].Status)
	}
	if !auth.IsApproved() {
		t.Error("IsApproved = false, want true")
	}

	req := server.LastRequest()
	assertRequest(t, req, http.MethodPost, testTransactionsPath)

	var body map[string]any
	if err := json.Unmarshal([]byte(req.Body), &body); err != nil {
		t.Fatalf("decoding request body: %v", err)
	}
	if body["username"] != testUsername {
		t.Errorf("body username = %v, want %q", body["username"], testUsername)
	}
	if body["tbk_user"] != testTbkUser {
		t.Errorf("body tbk_user = %v, want %q", body["tbk_user"], testTbkUser)
	}
	if body["buy_order"] != testBuyOrder {
		t.Errorf("body buy_order = %v, want %q", body["buy_order"], testBuyOrder)
	}
	rawDetails, ok := body["details"].([]any)
	if !ok {
		t.Fatalf("body details type = %T, want []any", body["details"])
	}
	if len(rawDetails) != 2 {
		t.Fatalf("len(details) = %d, want 2", len(rawDetails))
	}
	first, ok := rawDetails[0].(map[string]any)
	if !ok {
		t.Fatalf("details[0] type = %T, want map[string]any", rawDetails[0])
	}
	if first["commerce_code"] != testChildCode1 {
		t.Errorf("details[0].commerce_code = %v, want %q", first["commerce_code"], testChildCode1)
	}
	if first["amount"] != 10000.0 {
		t.Errorf("details[0].amount = %v, want 10000", first["amount"])
	}
	if first["installments_number"] != 3.0 {
		t.Errorf("details[0].installments_number = %v, want 3", first["installments_number"])
	}
}

func TestMallTransactionStatus(t *testing.T) {
	server := newMockServer(t, http.StatusOK, `{"buy_order":"`+testBuyOrder+`","card_detail":{"card_number":"6623"},"accounting_date":"2026-08-13","transaction_date":"2026-08-13T10:00:00.000Z","details":[{"amount":10000,"status":"AUTHORIZED","authorization_code":"`+testAuthCode+`","payment_type_code":"VN","response_code":0,"installments_number":3,"commerce_code":"`+testChildCode1+`","buy_order":"`+testChildBuyOrder+`"}]}`)
	tx, err := oneclick.NewMallTransaction(testOptions(server))
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
	if status.CardDetail.CardNumber != "6623" {
		t.Errorf("CardNumber = %q, want 6623", status.CardDetail.CardNumber)
	}

	assertRequest(t, server.LastRequest(), http.MethodGet, testTransactionsPath+"/"+testBuyOrder)
}

func TestMallTransactionRefund(t *testing.T) {
	server := newMockServer(t, http.StatusOK, `{"type":"REFUNDED","authorization_code":"`+testAuthCode+`","authorization_date":"2026-08-13","nullified_amount":10000,"balance":0,"response_code":0}`)
	tx, err := oneclick.NewMallTransaction(testOptions(server))
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
	if refund.NullifiedAmount != 10000 {
		t.Errorf("NullifiedAmount = %v, want 10000", refund.NullifiedAmount)
	}

	req := server.LastRequest()
	assertRequest(t, req, http.MethodPost, testTransactionsPath+"/"+testBuyOrder+"/refunds")

	var body map[string]any
	if err := json.Unmarshal([]byte(req.Body), &body); err != nil {
		t.Fatalf("decoding request body: %v", err)
	}
	if body["commerce_code"] != testChildCode1 {
		t.Errorf("body commerce_code = %v, want %q", body["commerce_code"], testChildCode1)
	}
	if body["detail_buy_order"] != testChildBuyOrder {
		t.Errorf("body detail_buy_order = %v, want %q", body["detail_buy_order"], testChildBuyOrder)
	}
	if body["amount"] != 10000.0 {
		t.Errorf("body amount = %v, want 10000", body["amount"])
	}
}

func TestMallTransactionCapture(t *testing.T) {
	server := newMockServer(t, http.StatusOK, `{"authorization_code":"`+testAuthCode+`","authorization_date":"2026-08-13T10:00:00.000Z","captured_amount":50,"response_code":0}`)
	tx, err := oneclick.NewMallTransaction(testOptions(server))
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
	if capture.ResponseCode != 0 {
		t.Errorf("ResponseCode = %d, want 0", capture.ResponseCode)
	}

	req := server.LastRequest()
	assertRequest(t, req, http.MethodPut, testTransactionsPath+"/capture")

	var body map[string]any
	if err := json.Unmarshal([]byte(req.Body), &body); err != nil {
		t.Fatalf("decoding request body: %v", err)
	}
	if body["commerce_code"] != testCaptureCode {
		t.Errorf("body commerce_code = %v, want %q", body["commerce_code"], testCaptureCode)
	}
	if body["buy_order"] != testBuyOrder {
		t.Errorf("body buy_order = %v, want %q", body["buy_order"], testBuyOrder)
	}
	if body["capture_amount"] != 50.0 {
		t.Errorf("body capture_amount = %v, want 50", body["capture_amount"])
	}
	if body["authorization_code"] != testAuthCode {
		t.Errorf("body authorization_code = %v, want %q", body["authorization_code"], testAuthCode)
	}
}

func TestMallTransactionHTTPError(t *testing.T) {
	server := newMockServer(t, http.StatusBadRequest, `{"error_message":"bad request"}`)
	tx, err := oneclick.NewMallTransaction(testOptions(server))
	if err != nil {
		t.Fatalf("NewMallTransaction: %v", err)
	}

	_, err = tx.Status(testBuyOrder)
	wantHTTPError(t, err, http.StatusBadRequest, `{"error_message":"bad request"}`)
}

func TestMallTransactionValidationSkipsRequest(t *testing.T) {
	server := newMockServer(t, http.StatusOK, `{}`)
	tx, err := oneclick.NewMallTransaction(testOptions(server))
	if err != nil {
		t.Fatalf("NewMallTransaction: %v", err)
	}

	if _, err := tx.Authorize("", testTbkUser, testBuyOrder, testDetails()); err == nil {
		t.Error("expected error for empty username")
	}
	if _, err := tx.Authorize(testUsername, "", testBuyOrder, testDetails()); err == nil {
		t.Error("expected error for empty tbk_user")
	}
	if _, err := tx.Authorize(testUsername, testTbkUser, "", testDetails()); err == nil {
		t.Error("expected error for empty buy_order")
	}
	if _, err := tx.Authorize(testUsername, testTbkUser, testBuyOrder, nil); err == nil {
		t.Error("expected error for empty details")
	}
	if _, err := tx.Authorize(testUsername, testTbkUser, testBuyOrder, []oneclick.MallDetails{{Amount: 100, CommerceCode: "bad", BuyOrder: testChildBuyOrder}}); err == nil {
		t.Error("expected error for invalid child commerce code")
	}
	if _, err := tx.Authorize(testUsername, testTbkUser, testBuyOrder, []oneclick.MallDetails{{Amount: 0, CommerceCode: testChildCode1, BuyOrder: testChildBuyOrder}}); err == nil {
		t.Error("expected error for zero amount")
	}
	if _, err := tx.Authorize(testUsername, testTbkUser, testBuyOrder, []oneclick.MallDetails{{Amount: 100, CommerceCode: testChildCode1, BuyOrder: testChildBuyOrder, InstallmentsNumber: 100}}); err == nil {
		t.Error("expected error for invalid installments_number")
	}
	if _, err := tx.Status(""); err == nil {
		t.Error("expected error for empty buy_order")
	}
	if _, err := tx.Refund(testBuyOrder, "bad", testChildBuyOrder, 100); err == nil {
		t.Error("expected error for invalid child commerce code")
	}
	if _, err := tx.Refund(testBuyOrder, testChildCode1, testChildBuyOrder, 0); err == nil {
		t.Error("expected error for zero refund amount")
	}
	if _, err := tx.Capture(testBuyOrder, testCaptureCode, "", 50); err == nil {
		t.Error("expected error for empty authorization_code")
	}

	if got := server.RequestCount(); got != 0 {
		t.Errorf("request count = %d, want 0 (validation must not hit the API)", got)
	}
}
