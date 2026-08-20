package patpass_test

import (
	"testing"

	"github.com/ppastene/transbank-sdk-go/patpass"
)

func TestInscriptionStart(t *testing.T) {
	m := newMockServer(t, 200, `{"token":"`+testToken+`","url":"`+testFormURL+`"}`)
	i, err := patpass.NewInscription(testOptions(m))
	if err != nil {
		t.Fatalf("NewInscription error: %v", err)
	}

	resp, err := i.Start(testURL, testName, testFLastname, testSLastname, testRut, testServiceId, testFinalURL, "", testPhone, testMobile, testPatpassName, testPersonEmail, testCommerceEmail, testAddress, testCity)
	if err != nil {
		t.Fatalf("Start error: %v", err)
	}
	if resp.Token != testToken {
		t.Errorf("Token = %q, want %q", resp.Token, testToken)
	}
	if resp.Url != testFormURL {
		t.Errorf("Url = %q, want %q", resp.Url, testFormURL)
	}

	req := m.LastRequest()
	assertRequest(t, req, "POST", testInscriptionPath)
	assertBody(t, req, `{"ciudad":"Santiago","commerceCode":"28299257","correoComercio":"comercio@test.cl","correoPersona":"persona@test.cl","direccion":"Merced 156, Santiago, Chile","finalUrl":"http://misitio.cl/voucher","montoMaximo":"","nombre":"Diego","nombrePatPass":"Help - 8050014","pApellido":"Sanchez","rut":"12345678-9","sApellido":"Valdovinos","serviceId":"323123","telefonoCelular":"57508624","telefonoFijo":"57508624","url":"http://misitio.cl/finalizar_suscripcion"}`)
}

func TestInscriptionStatus(t *testing.T) {
	m := newMockServer(t, 200, `{"authorized":true,"voucherUrl":"`+testVoucherURL+`"}`)
	i, err := patpass.NewInscription(testOptions(m))
	if err != nil {
		t.Fatalf("NewInscription error: %v", err)
	}

	resp, err := i.Status(testToken)
	if err != nil {
		t.Fatalf("Status error: %v", err)
	}
	if !resp.Authorized {
		t.Error("Authorized = false, want true")
	}
	if resp.VoucherUrl != testVoucherURL {
		t.Errorf("VoucherUrl = %q, want %q", resp.VoucherUrl, testVoucherURL)
	}

	req := m.LastRequest()
	assertRequest(t, req, "POST", testStatusPath)
	assertBody(t, req, `{"token":"`+testToken+`"}`)
}

func TestInscriptionStartHTTPError(t *testing.T) {
	m := newMockServer(t, 401, `{"error_message":"invalid credentials"}`)
	i, err := patpass.NewInscription(testOptions(m))
	if err != nil {
		t.Fatalf("NewInscription error: %v", err)
	}

	_, err = i.Start(testURL, testName, testFLastname, testSLastname, testRut, testServiceId, testFinalURL, "", testPhone, testMobile, testPatpassName, testPersonEmail, testCommerceEmail, testAddress, testCity)
	wantHTTPError(t, err, 401, `{"error_message":"invalid credentials"}`)
}

func TestInscriptionStatusHTTPError(t *testing.T) {
	m := newMockServer(t, 401, `{"error_message":"invalid credentials"}`)
	i, err := patpass.NewInscription(testOptions(m))
	if err != nil {
		t.Fatalf("NewInscription error: %v", err)
	}

	_, err = i.Status(testToken)
	wantHTTPError(t, err, 401, `{"error_message":"invalid credentials"}`)
}
