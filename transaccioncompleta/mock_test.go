package transaccioncompleta_test

import (
	"testing"

	"github.com/ppastene/transbank-sdk-go"
	"github.com/ppastene/transbank-sdk-go/internal/testutil"
)

const (
	testCommerceCode   = "597055555530"
	testMallCode       = "597055555551"
	testChildCode1     = "597055555552"
	testChildCode2     = "597055555553"
	testAPIKey         = "579B532A7440BB0C9079DED94D31EA1615BACEB56610332264630D42D0A36B1C"
	testToken          = "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
	testBuyOrder       = "orden-compra-123"
	testChildBuyOrder1 = "orden-hija-1"
	testChildBuyOrder2 = "orden-hija-2"
	testSessionID      = "sesion-456"
	testAmount         = 10000
	testCardNumber     = "4051885600446623"
	testCardExpiry     = "22/10"
	testCVV            = "123"
	testAuthCode       = "123456"
)

const testTransactionsPath = "/rswebpaytransaction/api/webpay/v1.2/transactions"

type mockServer = testutil.MockServer
type recordedRequest = testutil.RecordedRequest

func newMockServer(t *testing.T, status int, body string) *mockServer {
	return testutil.NewMockServer(t, status, body)
}

func testOptions(m *mockServer) transbank.Options {
	return transbank.Options{
		CommerceCode: testCommerceCode,
		ApiKey:       testAPIKey,
		Environment:  transbank.Integration,
		HTTPClient:   m.Client(),
	}
}

func wantHTTPError(t *testing.T, err error, statusCode int, body string) {
	t.Helper()
	testutil.WantHTTPError(t, err, statusCode, body)
}

func assertRequest(t *testing.T, req recordedRequest, method, path string) {
	t.Helper()
	testutil.AssertRequest(t, req, method, testCommerceCode, testAPIKey, path)
}

func assertMallRequest(t *testing.T, req recordedRequest, method, path string) {
	t.Helper()
	testutil.AssertRequest(t, req, method, testMallCode, testAPIKey, path)
}

func assertBody(t *testing.T, req recordedRequest, want string) {
	t.Helper()
	testutil.AssertBody(t, req, want)
}

func assertBodyNotContains(t *testing.T, req recordedRequest, substr string) {
	t.Helper()
	testutil.AssertBodyNotContains(t, req, substr)
}
