package transaccioncompleta_test

import (
	"encoding/json"
	"testing"

	"github.com/ppastene/transbank-sdk-go/transaccioncompleta"
)

func TestUnmarshalMallInstallmentsResponse(t *testing.T) {
	body := `[{"installments_amount":3334,"id_query_installments":11,"deferred_periods":[{"amount":1000,"period":1}]},{"installments_amount":6667,"id_query_installments":12,"deferred_periods":[]}]`
	var resp transaccioncompleta.MallTransactionInstallmentsResponse
	if err := json.Unmarshal([]byte(body), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(resp) != 2 {
		t.Fatalf("len(resp) = %d, want 2", len(resp))
	}
	if resp[0].InstallmentsAmount != 3334 {
		t.Errorf("resp[0].InstallmentsAmount = %d, want 3334", resp[0].InstallmentsAmount)
	}
	if resp[0].IdQueryInstallments != 11 {
		t.Errorf("resp[0].IdQueryInstallments = %d, want 11", resp[0].IdQueryInstallments)
	}
	if len(resp[0].DeferredPeriods) != 1 {
		t.Errorf("len(resp[0].DeferredPeriods) = %d, want 1", len(resp[0].DeferredPeriods))
	}
	if resp[0].DeferredPeriods[0].Amount != 1000 || resp[0].DeferredPeriods[0].Period != 1 {
		t.Errorf("resp[0].DeferredPeriods[0] = %+v, want {Amount:1000 Period:1}", resp[0].DeferredPeriods[0])
	}
	if len(resp[1].DeferredPeriods) != 0 {
		t.Errorf("len(resp[1].DeferredPeriods) = %d, want 0", len(resp[1].DeferredPeriods))
	}
}

func TestUnmarshalMallCommitAndStatusResponse(t *testing.T) {
	body := `{"buy_order":"orden-compra-123","card_detail":{"card_number":"4239"},"accounting_date":"0321","transaction_date":"2019-03-21T15:43:48.523Z","details":[{"amount":10000,"status":"AUTHORIZED","authorization_code":"123456","payment_type_code":"VN","response_code":0,"installments_amount":5000,"installments_number":2,"commerce_code":"597055555552","buy_order":"orden-hija-1","balance":0}]}`

	var commit transaccioncompleta.MallTransactionCommitResponse
	if err := json.Unmarshal([]byte(body), &commit); err != nil {
		t.Fatalf("unmarshal commit: %v", err)
	}
	if commit.BuyOrder != "orden-compra-123" {
		t.Errorf("BuyOrder = %q, want orden-compra-123", commit.BuyOrder)
	}
	if commit.CardDetail.CardNumber != "4239" {
		t.Errorf("CardDetail.CardNumber = %q, want 4239", commit.CardDetail.CardNumber)
	}
	if len(commit.Details) != 1 {
		t.Fatalf("len(Details) = %d, want 1", len(commit.Details))
	}
	if commit.Details[0].Amount != 10000 {
		t.Errorf("Details[0].Amount = %f, want 10000", commit.Details[0].Amount)
	}
	if commit.Details[0].InstallmentsNumber != 2 {
		t.Errorf("Details[0].InstallmentsNumber = %d, want 2", commit.Details[0].InstallmentsNumber)
	}

	var status transaccioncompleta.MallTransactionStatusResponse
	if err := json.Unmarshal([]byte(body), &status); err != nil {
		t.Fatalf("unmarshal status: %v", err)
	}
	if status.AccountingDate != "0321" {
		t.Errorf("AccountingDate = %q, want 0321", status.AccountingDate)
	}
	if !status.IsApproved() {
		t.Error("IsApproved = false, want true")
	}
}

func TestMallTransactionStatusResponseIsApproved(t *testing.T) {
	tests := []struct {
		name     string
		statuses []string
		codes    []int
		want     bool
	}{
		{name: "all authorized", statuses: []string{"AUTHORIZED", "AUTHORIZED"}, codes: []int{0, 0}, want: true},
		{name: "all captured", statuses: []string{"CAPTURED", "CAPTURED"}, codes: []int{0, 0}, want: true},
		{name: "one rejected", statuses: []string{"AUTHORIZED", "REJECTED"}, codes: []int{0, 0}, want: false},
		{name: "one failed response code", statuses: []string{"AUTHORIZED", "AUTHORIZED"}, codes: []int{0, -1}, want: false},
		{name: "empty details", statuses: nil, codes: nil, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := transaccioncompleta.MallTransactionStatusResponse{}
			for i := range tt.statuses {
				resp.Details = append(resp.Details, transaccioncompleta.MallTransactionDetails{
					Status:       tt.statuses[i],
					ResponseCode: tt.codes[i],
				})
			}
			if got := resp.IsApproved(); got != tt.want {
				t.Errorf("IsApproved = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMallTransactionDetailsIsApproved(t *testing.T) {
	tests := []struct {
		name     string
		status   string
		response int
		want     bool
	}{
		{name: "authorized", status: "AUTHORIZED", response: 0, want: true},
		{name: "captured", status: "CAPTURED", response: 0, want: true},
		{name: "rejected", status: "REJECTED", response: 0, want: false},
		{name: "failed response code", status: "AUTHORIZED", response: -1, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := transaccioncompleta.MallTransactionDetails{Status: tt.status, ResponseCode: tt.response}
			if got := d.IsApproved(); got != tt.want {
				t.Errorf("IsApproved = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnmarshalMallRefundResponseType(t *testing.T) {
	body := `{"type":"NULLIFY","authorization_code":"123456","authorization_date":"2019-03-20T20:18:20Z","nullified_amount":1000.00,"balance":0.00,"response_code":0}`
	var resp transaccioncompleta.MallTransactionRefundResponse
	if err := json.Unmarshal([]byte(body), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if resp.Type != "NULLIFY" {
		t.Errorf("Type = %q, want NULLIFY", resp.Type)
	}
	if resp.ResponseCode != 0 {
		t.Errorf("ResponseCode = %d, want 0", resp.ResponseCode)
	}
}

func TestUnmarshalMallCaptureResponse(t *testing.T) {
	body := `{"authorization_code":"123456","authorization_date":"2019-03-20T20:18:20Z","captured_amount":1000,"response_code":0}`
	var resp transaccioncompleta.MallTransactionCaptureResponse
	if err := json.Unmarshal([]byte(body), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if resp.CapturedAmount != 1000 {
		t.Errorf("CapturedAmount = %f, want 1000", resp.CapturedAmount)
	}
	if resp.AuthorizationCode != "123456" {
		t.Errorf("AuthorizationCode = %q, want 123456", resp.AuthorizationCode)
	}
}
