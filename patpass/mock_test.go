package patpass_test

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync"
	"testing"

	"github.com/ppastene/transbank-sdk-go"
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
	if req.Method != method {
		t.Errorf("method = %q, want %q", req.Method, method)
	}
	if req.Path != path {
		t.Errorf("path = %q, want %q", req.Path, path)
	}
	if got := req.Header.Get("Commercecode"); got != testPatpassCommerceCode {
		t.Errorf("Commercecode = %q, want %q", got, testPatpassCommerceCode)
	}
	if got := req.Header.Get("Authorization"); got != testPatpassAuth {
		t.Errorf("Authorization = %q, want %q", got, testPatpassAuth)
	}
	if got := req.Header.Get("Tbk-Api-Key-Id"); got != "" {
		t.Errorf("Tbk-Api-Key-Id = %q, want unset", got)
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
