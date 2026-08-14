package webpayplus_test

import (
	"testing"

	"github.com/ppastene/transbank-sdk-go/webpayplus"
)

func TestTransactionStatusResponseIsApproved(t *testing.T) {
	tests := []struct {
		name       string
		response   webpayplus.TransactionStatusResponse
		wantResult bool
	}{
		{name: "authorized", response: webpayplus.TransactionStatusResponse{ResponseCode: 0, Status: "AUTHORIZED"}, wantResult: true},
		{name: "captured", response: webpayplus.TransactionStatusResponse{ResponseCode: 0, Status: "CAPTURED"}, wantResult: true},
		{name: "reversed", response: webpayplus.TransactionStatusResponse{ResponseCode: 0, Status: "REVERSED"}, wantResult: false},
		{name: "nullified", response: webpayplus.TransactionStatusResponse{ResponseCode: 0, Status: "NULLIFIED"}, wantResult: false},
		{name: "initialized", response: webpayplus.TransactionStatusResponse{ResponseCode: 0, Status: "INITIALIZED"}, wantResult: false},
		{name: "unknown status", response: webpayplus.TransactionStatusResponse{ResponseCode: 0, Status: "WHATEVER"}, wantResult: false},
		{name: "non zero response code", response: webpayplus.TransactionStatusResponse{ResponseCode: -1, Status: "AUTHORIZED"}, wantResult: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.response.IsApproved(); got != tt.wantResult {
				t.Errorf("IsApproved() = %v, want %v", got, tt.wantResult)
			}
		})
	}
}

func TestMallTransactionDetailsResponseIsApproved(t *testing.T) {
	tests := []struct {
		name       string
		response   webpayplus.MallTransactionDetailsResponse
		wantResult bool
	}{
		{name: "authorized", response: webpayplus.MallTransactionDetailsResponse{ResponseCode: 0, Status: "AUTHORIZED"}, wantResult: true},
		{name: "captured", response: webpayplus.MallTransactionDetailsResponse{ResponseCode: 0, Status: "CAPTURED"}, wantResult: true},
		{name: "reversed", response: webpayplus.MallTransactionDetailsResponse{ResponseCode: 0, Status: "REVERSED"}, wantResult: false},
		{name: "nullified", response: webpayplus.MallTransactionDetailsResponse{ResponseCode: 0, Status: "NULLIFIED"}, wantResult: false},
		{name: "failed", response: webpayplus.MallTransactionDetailsResponse{ResponseCode: -1, Status: "AUTHORIZED"}, wantResult: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.response.IsApproved(); got != tt.wantResult {
				t.Errorf("IsApproved() = %v, want %v", got, tt.wantResult)
			}
		})
	}
}

func TestMallTransactionStatusResponseIsApproved(t *testing.T) {
	approved := webpayplus.MallTransactionDetailsResponse{ResponseCode: 0, Status: "AUTHORIZED"}
	declined := webpayplus.MallTransactionDetailsResponse{ResponseCode: -1, Status: "REJECTED"}

	tests := []struct {
		name       string
		response   webpayplus.MallTransactionStatusResponse
		wantResult bool
	}{
		{name: "all approved", response: webpayplus.MallTransactionStatusResponse{Details: []webpayplus.MallTransactionDetailsResponse{approved, approved}}, wantResult: true},
		{name: "any declined", response: webpayplus.MallTransactionStatusResponse{Details: []webpayplus.MallTransactionDetailsResponse{approved, declined}}, wantResult: false},
		{name: "no details", response: webpayplus.MallTransactionStatusResponse{Details: nil}, wantResult: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.response.IsApproved(); got != tt.wantResult {
				t.Errorf("IsApproved() = %v, want %v", got, tt.wantResult)
			}
		})
	}
}
