package testutil

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

type RecordedRequest struct {
	Method string
	Path   string
	Header http.Header
	Body   string
}

type MockServer struct {
	t        *testing.T
	server   *httptest.Server
	status   int
	body     string
	mu       sync.Mutex
	requests []RecordedRequest
}

func NewMockServer(t *testing.T, status int, body string) *MockServer {
	t.Helper()
	m := &MockServer{t: t, status: status, body: body}
	m.server = httptest.NewServer(http.HandlerFunc(m.handle))
	t.Cleanup(m.server.Close)
	return m
}

func (m *MockServer) handle(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		m.t.Errorf("reading request body: %v", err)
	}
	m.mu.Lock()
	m.requests = append(m.requests, RecordedRequest{
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

func (m *MockServer) Client() *http.Client {
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

func (m *MockServer) RequestCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.requests)
}

func (m *MockServer) LastRequest() RecordedRequest {
	m.mu.Lock()
	defer m.mu.Unlock()
	if len(m.requests) == 0 {
		m.t.Fatal("mock server received no requests")
	}
	return m.requests[len(m.requests)-1]
}

func WantHTTPError(t *testing.T, err error, statusCode int, body string) {
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

func AssertRequest(t *testing.T, req RecordedRequest, method, commerceCode, apiKey, path string) {
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
	if got := req.Header.Get("Tbk-Api-Key-Secret"); got != apiKey {
		t.Errorf("Tbk-Api-Key-Secret = %q, want %q", got, apiKey)
	}
	if got := req.Header.Get("Content-Type"); got != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", got)
	}
}

func AssertPatpassRequest(t *testing.T, req RecordedRequest, method, commerceCode, authorization, path string) {
	t.Helper()
	if req.Method != method {
		t.Errorf("method = %q, want %q", req.Method, method)
	}
	if req.Path != path {
		t.Errorf("path = %q, want %q", req.Path, path)
	}
	if got := req.Header.Get("Commercecode"); got != commerceCode {
		t.Errorf("Commercecode = %q, want %q", got, commerceCode)
	}
	if got := req.Header.Get("Authorization"); got != authorization {
		t.Errorf("Authorization = %q, want %q", got, authorization)
	}
	if got := req.Header.Get("Tbk-Api-Key-Id"); got != "" {
		t.Errorf("Tbk-Api-Key-Id = %q, want unset", got)
	}
	if got := req.Header.Get("Content-Type"); got != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", got)
	}
}

func AssertBody(t *testing.T, req RecordedRequest, want string) {
	t.Helper()
	if req.Body != want {
		t.Errorf("body = %q, want %q", req.Body, want)
	}
}

func WantTransportError(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var tbErr *transbank.TransportError
	if !errors.As(err, &tbErr) {
		t.Fatalf("error type = %T, want *transbank.TransportError", err)
	}
	if tbErr.Err == nil {
		t.Error("expected non-nil Err in TransportError")
	}
}

func AssertBodyNotContains(t *testing.T, req RecordedRequest, substr string) {
	t.Helper()
	if strings.Contains(req.Body, substr) {
		t.Errorf("body = %q must not contain %q", req.Body, substr)
	}
}
