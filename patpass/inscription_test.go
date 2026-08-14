package patpass_test

import (
	"testing"

	"github.com/ppastene/transbank-sdk-go/patpass"
)

func TestNewInscription(t *testing.T) {
	m := newMockServer(t, 200, "")
	i, err := patpass.NewInscription(testOptions(m))
	if err != nil {
		t.Fatalf("NewInscription(integration) error: %v", err)
	}
	if i == nil {
		t.Fatal("NewInscription(integration) returned nil")
	}

	i, err = patpass.NewInscription(patpass.Options{
		CommerceCode:  testPatpassCommerceCode,
		Authorization: testPatpassAuth,
		Environment:   patpass.Production,
		HTTPClient:    m.Client(),
	})
	if err != nil {
		t.Fatalf("NewInscription(production) error: %v", err)
	}
	if i == nil {
		t.Fatal("NewInscription(production) returned nil")
	}
}

func TestNewInscriptionInvalid(t *testing.T) {
	m := newMockServer(t, 200, "")

	cases := []struct {
		name string
		opts patpass.Options
	}{
		{"invalid environment", patpass.Options{CommerceCode: testPatpassCommerceCode, Authorization: testPatpassAuth, Environment: patpass.Environment(99), HTTPClient: m.Client()}},
		{"non numeric commerce code", patpass.Options{CommerceCode: "2829925a", Authorization: testPatpassAuth, Environment: patpass.Integration, HTTPClient: m.Client()}},
		{"empty authorization", patpass.Options{CommerceCode: testPatpassCommerceCode, Authorization: "", Environment: patpass.Integration, HTTPClient: m.Client()}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, err := patpass.NewInscription(c.opts)
			wantValidationError(t, err)
			if got := m.RequestCount(); got != 0 {
				t.Errorf("request count = %d, want 0", got)
			}
		})
	}
}

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

func TestInscriptionStartValidation(t *testing.T) {
	m := newMockServer(t, 200, "")
	i, err := patpass.NewInscription(testOptions(m))
	if err != nil {
		t.Fatalf("NewInscription error: %v", err)
	}

	valid := []string{testURL, testName, testFLastname, testSLastname, testRut, testServiceId, testFinalURL, "", testPhone, testMobile, testPatpassName, testPersonEmail, testCommerceEmail, testAddress, testCity}
	cases := []struct {
		name string
		args []string
	}{
		{"empty url", []string{"", testName, testFLastname, testSLastname, testRut, testServiceId, testFinalURL, "", testPhone, testMobile, testPatpassName, testPersonEmail, testCommerceEmail, testAddress, testCity}},
		{"relative url", []string{"misitio.cl", testName, testFLastname, testSLastname, testRut, testServiceId, testFinalURL, "", testPhone, testMobile, testPatpassName, testPersonEmail, testCommerceEmail, testAddress, testCity}},
		{"empty finalUrl", []string{testURL, testName, testFLastname, testSLastname, testRut, testServiceId, "", "", testPhone, testMobile, testPatpassName, testPersonEmail, testCommerceEmail, testAddress, testCity}},
		{"empty userEmail", []string{testURL, testName, testFLastname, testSLastname, testRut, testServiceId, testFinalURL, "", testPhone, testMobile, testPatpassName, "", testCommerceEmail, testAddress, testCity}},
		{"invalid userEmail", []string{testURL, testName, testFLastname, testSLastname, testRut, testServiceId, testFinalURL, "", testPhone, testMobile, testPatpassName, "persona", testCommerceEmail, testAddress, testCity}},
		{"empty commerceEmail", []string{testURL, testName, testFLastname, testSLastname, testRut, testServiceId, testFinalURL, "", testPhone, testMobile, testPatpassName, testPersonEmail, "", testAddress, testCity}},
		{"empty name", []string{testURL, "", testFLastname, testSLastname, testRut, testServiceId, testFinalURL, "", testPhone, testMobile, testPatpassName, testPersonEmail, testCommerceEmail, testAddress, testCity}},
		{"empty fLastname", []string{testURL, testName, "", testSLastname, testRut, testServiceId, testFinalURL, "", testPhone, testMobile, testPatpassName, testPersonEmail, testCommerceEmail, testAddress, testCity}},
		{"empty sLastname", []string{testURL, testName, testFLastname, "", testRut, testServiceId, testFinalURL, "", testPhone, testMobile, testPatpassName, testPersonEmail, testCommerceEmail, testAddress, testCity}},
		{"empty rut", []string{testURL, testName, testFLastname, testSLastname, "", testServiceId, testFinalURL, "", testPhone, testMobile, testPatpassName, testPersonEmail, testCommerceEmail, testAddress, testCity}},
		{"empty serviceId", []string{testURL, testName, testFLastname, testSLastname, testRut, "", testFinalURL, "", testPhone, testMobile, testPatpassName, testPersonEmail, testCommerceEmail, testAddress, testCity}},
		{"empty phoneNumber", []string{testURL, testName, testFLastname, testSLastname, testRut, testServiceId, testFinalURL, "", "", testMobile, testPatpassName, testPersonEmail, testCommerceEmail, testAddress, testCity}},
		{"empty mobileNumber", []string{testURL, testName, testFLastname, testSLastname, testRut, testServiceId, testFinalURL, "", testPhone, "", testPatpassName, testPersonEmail, testCommerceEmail, testAddress, testCity}},
		{"empty patPassName", []string{testURL, testName, testFLastname, testSLastname, testRut, testServiceId, testFinalURL, "", testPhone, testMobile, "", testPersonEmail, testCommerceEmail, testAddress, testCity}},
		{"empty userAddress", []string{testURL, testName, testFLastname, testSLastname, testRut, testServiceId, testFinalURL, "", testPhone, testMobile, testPatpassName, testPersonEmail, testCommerceEmail, "", testCity}},
		{"empty userCity", []string{testURL, testName, testFLastname, testSLastname, testRut, testServiceId, testFinalURL, "", testPhone, testMobile, testPatpassName, testPersonEmail, testCommerceEmail, testAddress, ""}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			args := c.args
			_, err := i.Start(args[0], args[1], args[2], args[3], args[4], args[5], args[6], args[7], args[8], args[9], args[10], args[11], args[12], args[13], args[14])
			wantValidationError(t, err)
			if got := m.RequestCount(); got != 0 {
				t.Errorf("request count = %d, want 0", got)
			}
		})
	}

	t.Run("empty maxAmount is allowed", func(t *testing.T) {
		m2 := newMockServer(t, 200, `{"token":"`+testToken+`","url":"`+testFormURL+`"}`)
		i2, err := patpass.NewInscription(testOptions(m2))
		if err != nil {
			t.Fatalf("NewInscription error: %v", err)
		}
		_, err = i2.Start(valid[0], valid[1], valid[2], valid[3], valid[4], valid[5], valid[6], valid[7], valid[8], valid[9], valid[10], valid[11], valid[12], valid[13], valid[14])
		if err != nil {
			t.Fatalf("Start with empty maxAmount error: %v", err)
		}
		if got := m2.RequestCount(); got != 1 {
			t.Errorf("request count = %d, want 1", got)
		}
	})
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

func TestInscriptionStatusValidation(t *testing.T) {
	m := newMockServer(t, 200, "")
	i, err := patpass.NewInscription(testOptions(m))
	if err != nil {
		t.Fatalf("NewInscription error: %v", err)
	}

	cases := []struct {
		name  string
		token string
	}{
		{"empty token", ""},
		{"short token", "abc"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, err := i.Status(c.token)
			wantValidationError(t, err)
			if got := m.RequestCount(); got != 0 {
				t.Errorf("request count = %d, want 0", got)
			}
		})
	}
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
