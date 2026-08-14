package oneclick_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/ppastene/transbank-sdk-go"
	"github.com/ppastene/transbank-sdk-go/oneclick"
)

func TestNewMallInscription(t *testing.T) {
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
			ins, err := oneclick.NewMallInscription(tt.opts)
			if tt.want != nil {
				wantValidationError(t, err)
				if ins != nil {
					t.Error("expected nil MallInscription on error")
				}
				return
			}
			if err != nil {
				t.Fatalf("NewMallInscription returned error: %v", err)
			}
			if ins == nil {
				t.Error("expected non-nil MallInscription")
			}
		})
	}
}

func TestMallInscriptionStart(t *testing.T) {
	server := newMockServer(t, http.StatusOK, `{"token":"`+testToken+`","url_webpay":"https://webpay.cl/form_inscription"}`)
	ins, err := oneclick.NewMallInscription(testOptions(server))
	if err != nil {
		t.Fatalf("NewMallInscription: %v", err)
	}

	start, err := ins.Start(testUsername, testEmail, testResponseURL)
	if err != nil {
		t.Fatalf("Start returned error: %v", err)
	}
	if start.Token != testToken {
		t.Errorf("Token = %q, want %q", start.Token, testToken)
	}
	if start.UrlWebpay != "https://webpay.cl/form_inscription" {
		t.Errorf("UrlWebpay = %q, want https://webpay.cl/form_inscription", start.UrlWebpay)
	}

	req := server.LastRequest()
	assertRequest(t, req, http.MethodPost, testInscriptionsPath)

	var body map[string]any
	if err := json.Unmarshal([]byte(req.Body), &body); err != nil {
		t.Fatalf("decoding request body: %v", err)
	}
	if body["username"] != testUsername {
		t.Errorf("body username = %v, want %q", body["username"], testUsername)
	}
	if body["email"] != testEmail {
		t.Errorf("body email = %v, want %q", body["email"], testEmail)
	}
	if body["response_url"] != testResponseURL {
		t.Errorf("body response_url = %v, want %q", body["response_url"], testResponseURL)
	}
}

func TestMallInscriptionFinish(t *testing.T) {
	server := newMockServer(t, http.StatusOK, `{"response_code":0,"tbk_user":"`+testTbkUser+`","authorization_code":"`+testAuthCode+`","card_type":"Visa","card_number":"XXXXXXXXXXXX6623"}`)
	ins, err := oneclick.NewMallInscription(testOptions(server))
	if err != nil {
		t.Fatalf("NewMallInscription: %v", err)
	}

	finish, err := ins.Finish(testToken)
	if err != nil {
		t.Fatalf("Finish returned error: %v", err)
	}
	if finish.ResponseCode != 0 {
		t.Errorf("ResponseCode = %d, want 0", finish.ResponseCode)
	}
	if finish.TbkUser != testTbkUser {
		t.Errorf("TbkUser = %q, want %q", finish.TbkUser, testTbkUser)
	}
	if finish.CardNumber != "XXXXXXXXXXXX6623" {
		t.Errorf("CardNumber = %q, want XXXXXXXXXXXX6623", finish.CardNumber)
	}

	req := server.LastRequest()
	assertRequest(t, req, http.MethodPut, testInscriptionsPath+"/"+testToken)
	if req.Body != "{}" {
		t.Errorf("body = %q, want {}", req.Body)
	}
}

func TestMallInscriptionDelete(t *testing.T) {
	server := newMockServer(t, http.StatusNoContent, ``)
	ins, err := oneclick.NewMallInscription(testOptions(server))
	if err != nil {
		t.Fatalf("NewMallInscription: %v", err)
	}

	if err := ins.Delete(testTbkUser, testUsername); err != nil {
		t.Fatalf("Delete returned error: %v", err)
	}

	req := server.LastRequest()
	assertRequest(t, req, http.MethodDelete, testInscriptionsPath)

	var body map[string]any
	if err := json.Unmarshal([]byte(req.Body), &body); err != nil {
		t.Fatalf("decoding request body: %v", err)
	}
	if body["tbk_user"] != testTbkUser {
		t.Errorf("body tbk_user = %v, want %q", body["tbk_user"], testTbkUser)
	}
	if body["username"] != testUsername {
		t.Errorf("body username = %v, want %q", body["username"], testUsername)
	}
}

func TestMallInscriptionHTTPError(t *testing.T) {
	server := newMockServer(t, http.StatusUnauthorized, `{"error_message":"unauthorized"}`)
	ins, err := oneclick.NewMallInscription(testOptions(server))
	if err != nil {
		t.Fatalf("NewMallInscription: %v", err)
	}

	_, err = ins.Finish(testToken)
	wantHTTPError(t, err, http.StatusUnauthorized, `{"error_message":"unauthorized"}`)
}

func TestMallInscriptionValidationSkipsRequest(t *testing.T) {
	server := newMockServer(t, http.StatusOK, `{}`)
	ins, err := oneclick.NewMallInscription(testOptions(server))
	if err != nil {
		t.Fatalf("NewMallInscription: %v", err)
	}

	if _, err := ins.Start("", testEmail, testResponseURL); err == nil {
		t.Error("expected error for empty username")
	}
	if _, err := ins.Start(testUsername, "", testResponseURL); err == nil {
		t.Error("expected error for empty email")
	}
	if _, err := ins.Start(testUsername, "correo-sin-arroba", testResponseURL); err == nil {
		t.Error("expected error for email without @")
	}
	if _, err := ins.Start(testUsername, testEmail, "not-a-url"); err == nil {
		t.Error("expected error for relative response_url")
	}
	if _, err := ins.Finish("short-token"); err == nil {
		t.Error("expected error for invalid token")
	}
	if err := ins.Delete("", testUsername); err == nil {
		t.Error("expected error for empty tbk_user")
	}
	if err := ins.Delete(testTbkUser, ""); err == nil {
		t.Error("expected error for empty username")
	}

	if got := server.RequestCount(); got != 0 {
		t.Errorf("request count = %d, want 0 (validation must not hit the API)", got)
	}
}
