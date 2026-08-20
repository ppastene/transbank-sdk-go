package oneclick_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/ppastene/transbank-sdk-go"
	"github.com/ppastene/transbank-sdk-go/oneclick"
)

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
	tests := []struct {
		name   string
		method func(ins *oneclick.MallInscription) error
	}{
		{
			name: "start",
			method: func(ins *oneclick.MallInscription) error {
				_, err := ins.Start(testUsername, testEmail, testResponseURL)
				return err
			},
		},
		{
			name: "finish",
			method: func(ins *oneclick.MallInscription) error {
				_, err := ins.Finish(testToken)
				return err
			},
		},
		{
			name: "delete",
			method: func(ins *oneclick.MallInscription) error {
				return ins.Delete(testTbkUser, testUsername)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := newMockServer(t, http.StatusUnauthorized, `{"error_message":"unauthorized"}`)
			ins, err := oneclick.NewMallInscription(testOptions(server))
			if err != nil {
				t.Fatalf("NewMallInscription: %v", err)
			}
			err = tt.method(ins)
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
