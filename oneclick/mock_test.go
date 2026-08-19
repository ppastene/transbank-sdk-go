package oneclick_test

import (
	"testing"

	"github.com/ppastene/transbank-sdk-go"
	"github.com/ppastene/transbank-sdk-go/internal/testutil"
)

const (
	testCommerceCode  = "597055555541"
	testAPIKey        = "579B532A7440BB0C9079DED94D31EA1615BACEB56610332264630D42D0A36B1C"
	testToken         = "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
	testBuyOrder      = "orden-de-compra-123"
	testChildBuyOrder = "orden-hija-456"
	testChildCode1    = "597055555542"
	testChildCode2    = "597055555543"
	testCaptureCode   = "597055555548"
	testUsername      = "juanperez"
	testEmail         = "juan.perez@gmail.com"
	testResponseURL   = "https://www.mi-tienda.cl/retorno"
	testTbkUser       = "b6bd6ba3-e718-4107-9386-d2b099a8dd42"
	testAuthCode      = "123456"
)

const (
	testInscriptionsPath = "/rswebpaytransaction/api/oneclick/v1.2/inscriptions"
	testTransactionsPath = "/rswebpaytransaction/api/oneclick/v1.2/transactions"
)

type mockServer = testutil.MockServer
type recordedRequest = testutil.RecordedRequest

func newMockServer(t *testing.T, status int, body string) *mockServer {
	return testutil.NewMockServer(t, status, body)
}

func testOptions(m *mockServer) transbank.Options {
	return transbank.Options{
		CommerceCode:   testCommerceCode,
		ApiKey:         testAPIKey,
		Environment:    transbank.Integration,
		HTTPClient:     m.Client(),
		ValidateInputs: true,
	}
}

func testOptionsNoValidation(m *mockServer) transbank.Options {
	return transbank.Options{
		CommerceCode: testCommerceCode,
		ApiKey:       testAPIKey,
		Environment:  transbank.Integration,
		HTTPClient:   m.Client(),
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
	testutil.AssertRequest(t, req, method, testCommerceCode, testAPIKey, path)
}
