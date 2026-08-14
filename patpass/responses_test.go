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

func TestUnmarshalInscriptionStatusResponseAuthorizedTrue(t *testing.T) {
	resp := patpass.InscriptionStatusResponse{}
	err := json.Unmarshal([]byte(`{"authorized":true,"voucherUrl":"`+testVoucherURL+`"}`), &resp)
	if err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}
	if !resp.Authorized {
		t.Error("Authorized = false, want true")
	}
	if resp.VoucherUrl != testVoucherURL {
		t.Errorf("VoucherUrl = %q, want %q", resp.VoucherUrl, testVoucherURL)
	}
}

func TestUnmarshalInscriptionStatusResponseAuthorizedFalse(t *testing.T) {
	resp := patpass.InscriptionStatusResponse{}
	err := json.Unmarshal([]byte(`{"authorized":false,"voucherUrl":"`+testVoucherURL+`"}`), &resp)
	if err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}
	if resp.Authorized {
		t.Error("Authorized = true, want false")
	}
}

func TestUnmarshalInscriptionStatusResponseRejectsStringAuthorized(t *testing.T) {
	resp := patpass.InscriptionStatusResponse{}
	err := json.Unmarshal([]byte(`{"authorized":"true","voucherUrl":"`+testVoucherURL+`"}`), &resp)
	if err == nil {
		t.Fatal("expected error for string authorized, got nil")
	}
}
