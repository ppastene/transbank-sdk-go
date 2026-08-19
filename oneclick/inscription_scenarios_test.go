package oneclick_test

import (
	"net/http"
	"testing"

	"github.com/ppastene/transbank-sdk-go/oneclick"
)

func TestMallInscriptionStartNoValidation(t *testing.T) {
	server := newMockServer(t, http.StatusOK, `{"token":"`+testToken+`","url_webpay":"https://webpay.cl/form_inscription"}`)
	ins, err := oneclick.NewMallInscription(testOptionsNoValidation(server))
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
	if server.RequestCount() != 1 {
		t.Errorf("request count = %d, want 1", server.RequestCount())
	}
}

func TestMallInscriptionStartNoValidationAPIRejection(t *testing.T) {
	server := newMockServer(t, http.StatusBadRequest, `{"error_message":"bad request"}`)
	ins, err := oneclick.NewMallInscription(testOptionsNoValidation(server))
	if err != nil {
		t.Fatalf("NewMallInscription: %v", err)
	}

	_, err = ins.Start("", testEmail, testResponseURL)
	wantHTTPError(t, err, http.StatusBadRequest, `{"error_message":"bad request"}`)
	if server.RequestCount() != 1 {
		t.Errorf("request count = %d, want 1 (API must be called when validation is off)", server.RequestCount())
	}
}

func TestMallInscriptionFinishNoValidation(t *testing.T) {
	server := newMockServer(t, http.StatusOK, `{"response_code":0,"tbk_user":"`+testTbkUser+`","authorization_code":"`+testAuthCode+`","card_type":"Visa","card_number":"XXXXXXXXXXXX6623"}`)
	ins, err := oneclick.NewMallInscription(testOptionsNoValidation(server))
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
	if server.RequestCount() != 1 {
		t.Errorf("request count = %d, want 1", server.RequestCount())
	}
}

func TestMallInscriptionFinishTokenAlwaysValidated(t *testing.T) {
	server := newMockServer(t, http.StatusOK, `{}`)
	ins, err := oneclick.NewMallInscription(testOptionsNoValidation(server))
	if err != nil {
		t.Fatalf("NewMallInscription: %v", err)
	}

	_, err = ins.Finish("short-token")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	wantValidationError(t, err)
	if server.RequestCount() != 0 {
		t.Errorf("request count = %d, want 0 (token validation must not hit API)", server.RequestCount())
	}
}

func TestMallInscriptionDeleteNoValidation(t *testing.T) {
	server := newMockServer(t, http.StatusNoContent, ``)
	ins, err := oneclick.NewMallInscription(testOptionsNoValidation(server))
	if err != nil {
		t.Fatalf("NewMallInscription: %v", err)
	}

	if err := ins.Delete(testTbkUser, testUsername); err != nil {
		t.Fatalf("Delete returned error: %v", err)
	}
	if server.RequestCount() != 1 {
		t.Errorf("request count = %d, want 1", server.RequestCount())
	}
}

func TestMallInscriptionDeleteNoValidationAPIRejection(t *testing.T) {
	server := newMockServer(t, http.StatusBadRequest, `{"error_message":"bad request"}`)
	ins, err := oneclick.NewMallInscription(testOptionsNoValidation(server))
	if err != nil {
		t.Fatalf("NewMallInscription: %v", err)
	}

	err = ins.Delete("", testUsername)
	wantHTTPError(t, err, http.StatusBadRequest, `{"error_message":"bad request"}`)
	if server.RequestCount() != 1 {
		t.Errorf("request count = %d, want 1 (API must be called when validation is off)", server.RequestCount())
	}
}
