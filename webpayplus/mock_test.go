package webpayplus_test

import (
	"testing"

	"github.com/ppastene/transbank-sdk-go"
	"github.com/ppastene/transbank-sdk-go/internal/testutil"
)

const (
	testCommerceCode = "597055555532"
	testAPIKey       = "579B532A7440BB0C9079DED94D31EA1615BACEB56610332264630D42D0A36B1C"
	testToken        = "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
	testBuyOrder     = "orden-de-compra-123"
	testSessionID    = "sesion-456"
	testReturnURL    = "https://www.mi-tienda.cl/retorno"
	testChildCode1   = "597055555536"
	testChildCode2   = "597055555537"
	testAuthCode     = "123456"
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

func assertRequest(t *testing.T, req recordedRequest, method, path string) {
	t.Helper()
	testutil.AssertRequest(t, req, method, testCommerceCode, testAPIKey, path)
}
