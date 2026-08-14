package transaccioncompleta_test

import (
	"encoding/json"
	"testing"

	"github.com/ppastene/transbank-sdk-go/transaccioncompleta"
)

func TestUnmarshalInstallmentsResponse(t *testing.T) {
	body := `{"installments_amount":3334,"id_query_installments":11,"deferred_periods":[{"amount":1000,"period":1},{"amount":2000,"period":2}]}`
	var resp transaccioncompleta.TransactionInstallmentsResponse
	if err := json.Unmarshal([]byte(body), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if resp.InstallmentsAmount != 3334 {
		t.Errorf("InstallmentsAmount = %d, want 3334", resp.InstallmentsAmount)
	}
	if resp.IdQueryInstallments != 11 {
		t.Errorf("IdQueryInstallments = %d, want 11", resp.IdQueryInstallments)
	}
	if len(resp.DeferredPeriods) != 2 {
		t.Fatalf("len(DeferredPeriods) = %d, want 2", len(resp.DeferredPeriods))
	}
	if resp.DeferredPeriods[0].Amount != 1000 || resp.DeferredPeriods[0].Period != 1 {
		t.Errorf("DeferredPeriods[0] = %+v, want {Amount:1000 Period:1}", resp.DeferredPeriods[0])
	}
	if resp.DeferredPeriods[1].Amount != 2000 || resp.DeferredPeriods[1].Period != 2 {
		t.Errorf("DeferredPeriods[1] = %+v, want {Amount:2000 Period:2}", resp.DeferredPeriods[1])
	}
}

func TestUnmarshalInstallmentsResponseEmptyDeferredPeriods(t *testing.T) {
	body := `{"installments_amount":3334,"id_query_installments":11,"deferred_periods":[]}`
	var resp transaccioncompleta.TransactionInstallmentsResponse
	if err := json.Unmarshal([]byte(body), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if resp.InstallmentsAmount != 3334 {
		t.Errorf("InstallmentsAmount = %d, want 3334", resp.InstallmentsAmount)
	}
	if len(resp.DeferredPeriods) != 0 {
		t.Errorf("len(DeferredPeriods) = %d, want 0", len(resp.DeferredPeriods))
	}
}

func TestUnmarshalInstallmentsResponseRejectsStringTypes(t *testing.T) {
	body := `{"installments_amount":"3334","id_query_installments":"11","deferred_periods":[{"amount":"1000","period":"1"}]}`
	var resp transaccioncompleta.TransactionInstallmentsResponse
	if err := json.Unmarshal([]byte(body), &resp); err == nil {
		t.Fatal("expected unmarshal error when the API sends strings, got nil")
	}
}

func TestUnmarshalRefundResponseType(t *testing.T) {
	body := `{"type":"NULLIFY","authorization_code":"123456","authorization_date":"2019-03-20T20:18:20Z","nullified_amount":1000.00,"balance":0.00,"response_code":0}`
	var resp transaccioncompleta.TransactionRefundResponse
	if err := json.Unmarshal([]byte(body), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if resp.Type != "NULLIFY" {
		t.Errorf("Type = %q, want NULLIFY", resp.Type)
	}
	if resp.ResponseCode != 0 {
		t.Errorf("ResponseCode = %d, want 0", resp.ResponseCode)
	}
	if resp.NullifiedAmount != 1000 {
		t.Errorf("NullifiedAmount = %f, want 1000", resp.NullifiedAmount)
	}
}

func TestUnmarshalRefundResponseReversed(t *testing.T) {
	body := `{"type":"REVERSED"}`
	var resp transaccioncompleta.TransactionRefundResponse
	if err := json.Unmarshal([]byte(body), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if resp.Type != "REVERSED" {
		t.Errorf("Type = %q, want REVERSED", resp.Type)
	}
	if resp.AuthorizationCode != "" || resp.NullifiedAmount != 0 {
		t.Errorf("expected zero values when only type is sent, got %+v", resp)
	}
}

func TestUnmarshalCommitAndStatusResponse(t *testing.T) {
	body := `{"vci":"TSY","amount":10000,"status":"AUTHORIZED","buy_order":"orden-compra-123","session_id":"sesion-456","card_detail":{"card_number":"4239"},"accounting_date":"2026-08-13","transaction_date":"2026-08-13T10:00:00.000Z","authorization_code":"123456","payment_type_code":"VN","response_code":0,"installments_number":3,"installments_amount":3334,"balance":0}`

	var commit transaccioncompleta.TransactionCommitResponse
	if err := json.Unmarshal([]byte(body), &commit); err != nil {
		t.Fatalf("unmarshal commit: %v", err)
	}
	if commit.Status != "AUTHORIZED" {
		t.Errorf("Status = %q, want AUTHORIZED", commit.Status)
	}
	if commit.Amount != 10000 {
		t.Errorf("Amount = %f, want 10000", commit.Amount)
	}
	if commit.CardDetail.CardNumber != "4239" {
		t.Errorf("CardDetail.CardNumber = %q, want 4239", commit.CardDetail.CardNumber)
	}
	if commit.InstallmentsNumber != 3 {
		t.Errorf("InstallmentsNumber = %d, want 3", commit.InstallmentsNumber)
	}

	var status transaccioncompleta.TransactionStatusResponse
	if err := json.Unmarshal([]byte(body), &status); err != nil {
		t.Fatalf("unmarshal status: %v", err)
	}
	if status.BuyOrder != "orden-compra-123" {
		t.Errorf("BuyOrder = %q, want orden-compra-123", status.BuyOrder)
	}
	if !status.IsApproved() {
		t.Error("IsApproved = false, want true")
	}
}

func TestTransactionStatusResponseIsApproved(t *testing.T) {
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
			resp := transaccioncompleta.TransactionStatusResponse{Status: tt.status, ResponseCode: tt.response}
			if got := resp.IsApproved(); got != tt.want {
				t.Errorf("IsApproved = %v, want %v", got, tt.want)
			}
		})
	}
}
