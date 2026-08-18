package transaccioncompleta_test

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"testing"

	"github.com/ppastene/transbank-sdk-go"
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

type recordedRequest struct {
	Method string
	Path   string
	Header http.Header
	Body   string
}

type mockServer struct {
	t        *testing.T
	server   *httptest.Server
	status   int
	body     string
	mu       sync.Mutex
	requests []recordedRequest
}

func newMockServer(t *testing.T, status int, body string) *mockServer {
	t.Helper()
	m := &mockServer{t: t, status: status, body: body}
	m.server = httptest.NewServer(http.HandlerFunc(m.handle))
	t.Cleanup(m.server.Close)
	return m
}

func (m *mockServer) handle(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		m.t.Errorf("reading request body: %v", err)
	}
	m.mu.Lock()
	m.requests = append(m.requests, recordedRequest{
		Method: r.Method,
		Path:   r.URL.Path,
		Header: r.Header,
		Body:   string(body),
	})
	m.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(m.status)
	if m.body != "" {
		w.Write([]byte(m.body))
	}
}

func (m *mockServer) Client() *http.Client {
	base, err := url.Parse(m.server.URL)
	if err != nil {
		m.t.Fatalf("parsing mock server URL: %v", err)
	}
	return &http.Client{
		Transport: mockTransport{base: base},
	}
}

type mockTransport struct {
	base *url.URL
}

func (t mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	target := *req.URL
	target.Scheme = t.base.Scheme
	target.Host = t.base.Host
	clone := req.Clone(req.Context())
	clone.URL = &target
	clone.RequestURI = ""
	return http.DefaultTransport.RoundTrip(clone)
}

func (m *mockServer) RequestCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.requests)
}

func (m *mockServer) LastRequest() recordedRequest {
	m.mu.Lock()
	defer m.mu.Unlock()
	if len(m.requests) == 0 {
		m.t.Fatal("mock server received no requests")
	}
	return m.requests[len(m.requests)-1]
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

func wantValidationError(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var tbErr *transbank.ValidationError
	if !errors.As(err, &tbErr) {
		t.Fatalf("error type = %T, want *transbank.ValidationError", err)
	}
}

func wantHTTPError(t *testing.T, err error, statusCode int, body string) {
	t.Helper()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var tbErr *transbank.HTTPError
	if !errors.As(err, &tbErr) {
		t.Fatalf("error type = %T, want *transbank.HTTPError", err)
	}
	if tbErr.StatusCode != statusCode {
		t.Errorf("StatusCode = %d, want %d", tbErr.StatusCode, statusCode)
	}
	if tbErr.Body != body {
		t.Errorf("Body = %q, want %q", tbErr.Body, body)
	}
}

func assertRequest(t *testing.T, req recordedRequest, method, path string) {
	t.Helper()
	assertRequestWith(t, req, method, path, testCommerceCode)
}

func assertMallRequest(t *testing.T, req recordedRequest, method, path string) {
	t.Helper()
	assertRequestWith(t, req, method, path, testMallCode)
}

func assertRequestWith(t *testing.T, req recordedRequest, method, path, commerceCode string) {
	t.Helper()
	if req.Method != method {
		t.Errorf("method = %q, want %q", req.Method, method)
	}
	if req.Path != path {
		t.Errorf("path = %q, want %q", req.Path, path)
	}
	if got := req.Header.Get("Tbk-Api-Key-Id"); got != commerceCode {
		t.Errorf("Tbk-Api-Key-Id = %q, want %q", got, commerceCode)
	}
	if got := req.Header.Get("Tbk-Api-Key-Secret"); got != testAPIKey {
		t.Errorf("Tbk-Api-Key-Secret = %q, want %q", got, testAPIKey)
	}
	if got := req.Header.Get("Content-Type"); got != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", got)
	}
}

func assertBody(t *testing.T, req recordedRequest, want string) {
	t.Helper()
	if req.Body != want {
		t.Errorf("body = %q, want %q", req.Body, want)
	}
}

func assertBodyNotContains(t *testing.T, req recordedRequest, substr string) {
	t.Helper()
	if strings.Contains(req.Body, substr) {
		t.Errorf("body = %q must not contain %q", req.Body, substr)
	}
}
