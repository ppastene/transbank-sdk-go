package patpass_test

import (
	"encoding/json"
	"testing"

	"github.com/ppastene/transbank-sdk-go/patpass"
)

func TestUnmarshalInscriptionStartResponse(t *testing.T) {
	resp := patpass.InscriptionStartResponse{}
	err := json.Unmarshal([]byte(`{"token":"`+testToken+`","url":"`+testFormURL+`"}`), &resp)
	if err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}
	if resp.Token != testToken {
		t.Errorf("Token = %q, want %q", resp.Token, testToken)
	}
	if resp.Url != testFormURL {
		t.Errorf("Url = %q, want %q", resp.Url, testFormURL)
	}
}

func TestUnmarshalInscriptionStatusResponse(t *testing.T) {
	tests := []struct {
		name       string
		json       string
		authorized bool
		voucher    string
	}{
		{
			name:       "authorized true",
			json:       `{"authorized":true,"voucherUrl":"` + testVoucherURL + `"}`,
			authorized: true,
			voucher:    testVoucherURL,
		},
		{
			name:       "authorized false",
			json:       `{"authorized":false,"voucherUrl":"` + testVoucherURL + `"}`,
			authorized: false,
			voucher:    testVoucherURL,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := patpass.InscriptionStatusResponse{}
			err := json.Unmarshal([]byte(tt.json), &resp)
			if err != nil {
				t.Fatalf("Unmarshal error: %v", err)
			}
			if resp.Authorized != tt.authorized {
				t.Errorf("Authorized = %v, want %v", resp.Authorized, tt.authorized)
			}
			if resp.VoucherUrl != tt.voucher {
				t.Errorf("VoucherUrl = %q, want %q", resp.VoucherUrl, tt.voucher)
			}
		})
	}
}

func TestUnmarshalInscriptionStatusResponseRejectsStringAuthorized(t *testing.T) {
	resp := patpass.InscriptionStatusResponse{}
	err := json.Unmarshal([]byte(`{"authorized":"true","voucherUrl":"`+testVoucherURL+`"}`), &resp)
	if err == nil {
		t.Fatal("expected error for string authorized, got nil")
	}
}
