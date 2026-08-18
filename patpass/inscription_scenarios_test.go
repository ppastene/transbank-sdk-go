package patpass_test

import (
	"net/http"
	"testing"

	"github.com/ppastene/transbank-sdk-go/patpass"
)

func TestInscriptionStartNoValidation(t *testing.T) {
	m := newMockServer(t, 200, `{"token":"`+testToken+`","url":"`+testFormURL+`"}`)
	i, err := patpass.NewInscription(testOptionsNoValidation(m))
	if err != nil {
		t.Fatalf("NewInscription: %v", err)
	}
	resp, err := i.Start(testURL, testName, testFLastname, testSLastname, testRut, testServiceId, testFinalURL, "", testPhone, testMobile, testPatpassName, testPersonEmail, testCommerceEmail, testAddress, testCity)
	if err != nil {
		t.Fatalf("Start: %v", err)
	}
	assertRequest(t, m.LastRequest(), http.MethodPost, testInscriptionPath)
	assertBody(t, m.LastRequest(), `{"ciudad":"Santiago","commerceCode":"28299257","correoComercio":"comercio@test.cl","correoPersona":"persona@test.cl","direccion":"Merced 156, Santiago, Chile","finalUrl":"http://misitio.cl/voucher","montoMaximo":"","nombre":"Diego","nombrePatPass":"Help - 8050014","pApellido":"Sanchez","rut":"12345678-9","sApellido":"Valdovinos","serviceId":"323123","telefonoCelular":"57508624","telefonoFijo":"57508624","url":"http://misitio.cl/finalizar_suscripcion"}`)
	if resp.Token != testToken {
		t.Errorf("Token = %q, want %q", resp.Token, testToken)
	}
	if m.RequestCount() != 1 {
		t.Errorf("request count = %d, want 1", m.RequestCount())
	}
}

func TestInscriptionStartNoValidationAPIRejection(t *testing.T) {
	m := newMockServer(t, http.StatusBadRequest, `{"error_message":"bad request"}`)
	i, err := patpass.NewInscription(testOptionsNoValidation(m))
	if err != nil {
		t.Fatalf("NewInscription: %v", err)
	}
	_, err = i.Start("", testName, testFLastname, testSLastname, testRut, testServiceId, testFinalURL, "", testPhone, testMobile, testPatpassName, testPersonEmail, testCommerceEmail, testAddress, testCity)
	wantHTTPError(t, err, http.StatusBadRequest, `{"error_message":"bad request"}`)
	if m.RequestCount() != 1 {
		t.Errorf("request count = %d, want 1 (API must be called when validation is off)", m.RequestCount())
	}
}

func TestInscriptionStatusNoValidation(t *testing.T) {
	m := newMockServer(t, 200, `{"authorized":true,"voucherUrl":"`+testVoucherURL+`"}`)
	i, err := patpass.NewInscription(testOptionsNoValidation(m))
	if err != nil {
		t.Fatalf("NewInscription: %v", err)
	}
	resp, err := i.Status(testToken)
	if err != nil {
		t.Fatalf("Status: %v", err)
	}
	assertRequest(t, m.LastRequest(), http.MethodPost, testStatusPath)
	assertBody(t, m.LastRequest(), `{"token":"`+testToken+`"}`)
	if !resp.Authorized {
		t.Error("Authorized = false, want true")
	}
	if m.RequestCount() != 1 {
		t.Errorf("request count = %d, want 1", m.RequestCount())
	}
}

func TestInscriptionStatusNoValidationAPIRejection(t *testing.T) {
	m := newMockServer(t, http.StatusBadRequest, `{"error_message":"bad request"}`)
	i, err := patpass.NewInscription(testOptionsNoValidation(m))
	if err != nil {
		t.Fatalf("NewInscription: %v", err)
	}
	_, err = i.Status("")
	wantHTTPError(t, err, http.StatusBadRequest, `{"error_message":"bad request"}`)
	if m.RequestCount() != 1 {
		t.Errorf("request count = %d, want 1 (API must be called when validation is off)", m.RequestCount())
	}
}
