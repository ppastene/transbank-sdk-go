package patpass_test

import (
	"testing"

	"github.com/ppastene/transbank-sdk-go/internal/testutil"
	"github.com/ppastene/transbank-sdk-go/patpass"
)

const (
	testPatpassCommerceCode = "28299257"
	testPatpassAuth         = "cxxXQgGD9vrVe4M41FIt"
	testToken               = "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
	testURL                 = "http://misitio.cl/finalizar_suscripcion"
	testName                = "Diego"
	testFLastname           = "Sanchez"
	testSLastname           = "Valdovinos"
	testRut                 = "12345678-9"
	testServiceId           = "323123"
	testFinalURL            = "http://misitio.cl/voucher"
	testPhone               = "57508624"
	testMobile              = "57508624"
	testPatpassName         = "Help - 8050014"
	testPersonEmail         = "persona@test.cl"
	testCommerceEmail       = "comercio@test.cl"
	testAddress             = "Merced 156, Santiago, Chile"
	testCity                = "Santiago"
	testFormURL             = "https://pagoautomaticocontarjetasint.transbank.cl/nuevo-ic-rest/tokenComercioLogin"
	testVoucherURL          = "https://pagoautomaticocontarjetasint.transbank.cl/nuevo-ic-rest/tokenVoucherLogin"
)

const (
	testInscriptionPath = "/restpatpass/v1/services/patInscription"
	testStatusPath      = "/restpatpass/v1/services/status"
)

type mockServer = testutil.MockServer
type recordedRequest = testutil.RecordedRequest

func newMockServer(t *testing.T, status int, body string) *mockServer {
	return testutil.NewMockServer(t, status, body)
}

func testOptions(m *mockServer) patpass.Options {
	return patpass.Options{
		CommerceCode:   testPatpassCommerceCode,
		Authorization:  testPatpassAuth,
		Environment:    patpass.Integration,
		HTTPClient:     m.Client(),
		ValidateInputs: true,
	}
}

func testOptionsNoValidation(m *mockServer) patpass.Options {
	return patpass.Options{
		CommerceCode:  testPatpassCommerceCode,
		Authorization: testPatpassAuth,
		Environment:   patpass.Integration,
		HTTPClient:    m.Client(),
	}
}

func wantValidationError(t *testing.T, err error) {
	t.Helper()
	testutil.WantValidationError(t, err)
}

func wantHTTPError(t *testing.T, err error, statusCode int, body string) {
	t.Helper()
	testutil.WantHTTPError(t, err, statusCode, body)
}

func assertRequest(t *testing.T, req recordedRequest, method, path string) {
	t.Helper()
	testutil.AssertPatpassRequest(t, req, method, testPatpassCommerceCode, testPatpassAuth, path)
}

func assertBody(t *testing.T, req recordedRequest, want string) {
	t.Helper()
	testutil.AssertBody(t, req, want)
}
