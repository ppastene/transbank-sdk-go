package oneclick_test

import (
	"testing"

	"github.com/ppastene/transbank-sdk-go/oneclick"
)

func TestOneclickMallTransactionDetailsIsApproved(t *testing.T) {
	tests := []struct {
		name       string
		response   oneclick.OneclickMallTransactionDetails
		wantResult bool
	}{
		{name: "authorized", response: oneclick.OneclickMallTransactionDetails{ResponseCode: 0, Status: "AUTHORIZED"}, wantResult: true},
		{name: "captured", response: oneclick.OneclickMallTransactionDetails{ResponseCode: 0, Status: "CAPTURED"}, wantResult: true},
		{name: "rejected", response: oneclick.OneclickMallTransactionDetails{ResponseCode: 0, Status: "REJECTED"}, wantResult: false},
		{name: "non zero response code", response: oneclick.OneclickMallTransactionDetails{ResponseCode: -1, Status: "AUTHORIZED"}, wantResult: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.response.IsApproved(); got != tt.wantResult {
				t.Errorf("IsApproved() = %v, want %v", got, tt.wantResult)
			}
		})
	}
}

func TestOneclickMallTransactionStatusResponseIsApproved(t *testing.T) {
	approved := oneclick.OneclickMallTransactionDetails{ResponseCode: 0, Status: "AUTHORIZED"}
	declined := oneclick.OneclickMallTransactionDetails{ResponseCode: -1, Status: "REJECTED"}

	tests := []struct {
		name       string
		response   oneclick.OneclickMallTransactionStatusResponse
		wantResult bool
	}{
		{name: "all approved", response: oneclick.OneclickMallTransactionStatusResponse{Details: []oneclick.OneclickMallTransactionDetails{approved, approved}}, wantResult: true},
		{name: "any declined", response: oneclick.OneclickMallTransactionStatusResponse{Details: []oneclick.OneclickMallTransactionDetails{approved, declined}}, wantResult: false},
		{name: "no details", response: oneclick.OneclickMallTransactionStatusResponse{Details: nil}, wantResult: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.response.IsApproved(); got != tt.wantResult {
				t.Errorf("IsApproved() = %v, want %v", got, tt.wantResult)
			}
		})
	}
}
