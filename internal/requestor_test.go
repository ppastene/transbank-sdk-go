package internal

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/ppastene/transbank-sdk-go"
)

const (
	testCommerceCode = "597055555532"
	testAPIKey       = "579B532A7440BB0C9079DED94D31EA1615BACEB56610332264630D42D0A36B1C"
)

func TestRequestorPostSendsHeadersAndDecodesBody(t *testing.T) {
	var gotAPIKeyID, gotAPIKeySecret, gotContentType string
	var gotBody map[string]any

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAPIKeyID = r.Header.Get("Tbk-Api-Key-Id")
		gotAPIKeySecret = r.Header.Get("Tbk-Api-Key-Secret")
		gotContentType = r.Header.Get(headerContentType)

		if r.Method != http.MethodPost {
			t.Errorf("method = %q, want POST", r.Method)
		}
		if r.URL.Path != "/transactions" {
			t.Errorf("path = %q, want /transactions", r.URL.Path)
		}

		if err := json.NewDecoder(r.Body).Decode(&gotBody); err != nil {
			t.Errorf("decoding request body: %v", err)
		}

		w.Header().Set(headerContentType, mimeJSON)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"token":"tok123","status":"INITIALIZED"}`))
	}))
	defer server.Close()

	cfg := NewConfig(testCommerceCode, testAPIKey, server.URL, )
	cfg.Headers = map[string]string{
		"Tbk-Api-Key-Id":     testCommerceCode,
		"Tbk-Api-Key-Secret": testAPIKey,
	}
	requestor := NewRequestor(&cfg)

	var result struct {
		Token  string `json:"token"`
		Status string `json:"status"`
	}
	if err := requestor.Post("/transactions", map[string]any{"buy_order": "abc"}, &result); err != nil {
		t.Fatalf("Post returned error: %v", err)
	}

	if gotAPIKeyID != testCommerceCode {
		t.Errorf("Tbk-Api-Key-Id = %q, want %q", gotAPIKeyID, testCommerceCode)
	}
	if gotAPIKeySecret != testAPIKey {
		t.Errorf("Tbk-Api-Key-Secret = %q, want %q", gotAPIKeySecret, testAPIKey)
	}
	if !strings.Contains(gotContentType, "application/json") {
		t.Errorf("Content-Type = %q, want application/json", gotContentType)
	}
	if gotBody["buy_order"] != "abc" {
		t.Errorf("request body buy_order = %v, want abc", gotBody["buy_order"])
	}
	if result.Token != "tok123" || result.Status != "INITIALIZED" {
		t.Errorf("result = %+v, want token=tok123 status=INITIALIZED", result)
	}
}

func TestRequestorReturnsAPIErrorOnNonSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error_message":"unauthorized"}`))
	}))
	defer server.Close()

	cfg := NewConfig(testCommerceCode, testAPIKey, server.URL, )
	requestor := NewRequestor(&cfg)

	err := requestor.Get("/transactions", nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	var tbErr *transbank.HTTPError
	if !errors.As(err, &tbErr) {
		t.Fatalf("error type = %T, want *transbank.HTTPError", err)
	}
	if tbErr.StatusCode != http.StatusUnauthorized {
		t.Errorf("StatusCode = %d, want %d", tbErr.StatusCode, http.StatusUnauthorized)
	}
}

type stubClient struct {
	lastReq *http.Request
}

func (c *stubClient) Do(req *http.Request) (*http.Response, error) {
	c.lastReq = req
	return &http.Response{
		StatusCode: http.StatusOK,
		Status:     "200 OK",
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(strings.NewReader(`{"token":"tok456"}`)),
	}, nil
}

func TestRequestorUsesInjectedHTTPClient(t *testing.T) {
	client := &stubClient{}
	cfg := NewConfig(testCommerceCode, testAPIKey, "https://webpay3gint.transbank.cl", )
	cfg.HTTP = client
	cfg.Headers = map[string]string{
		"Tbk-Api-Key-Id":     testCommerceCode,
		"Tbk-Api-Key-Secret": testAPIKey,
	}
	requestor := NewRequestor(&cfg)

	var result struct {
		Token string `json:"token"`
	}
	if err := requestor.Post("/transactions", map[string]any{"buy_order": "abc"}, &result); err != nil {
		t.Fatalf("Post returned error: %v", err)
	}
	if client.lastReq == nil {
		t.Fatal("injected client was not called")
	}
	if got := client.lastReq.URL.String(); got != "https://webpay3gint.transbank.cl/transactions" {
		t.Errorf("URL = %q, want https://webpay3gint.transbank.cl/transactions", got)
	}
	if result.Token != "tok456" {
		t.Errorf("result token = %q, want tok456", result.Token)
	}
}

func TestRequestorPut(t *testing.T) {
	var gotMethod string
	var gotBody map[string]any

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		json.NewDecoder(r.Body).Decode(&gotBody)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"updated"}`))
	}))
	defer server.Close()

	cfg := NewConfig(testCommerceCode, testAPIKey, server.URL, )
	cfg.Headers = map[string]string{
		"Tbk-Api-Key-Id":     testCommerceCode,
		"Tbk-Api-Key-Secret": testAPIKey,
	}
	requestor := NewRequestor(&cfg)

	var result struct {
		Status string `json:"status"`
	}
	if err := requestor.Put("/resource/123", map[string]any{"field": "value"}, &result); err != nil {
		t.Fatalf("Put returned error: %v", err)
	}
	if gotMethod != http.MethodPut {
		t.Errorf("method = %q, want PUT", gotMethod)
	}
	if gotBody["field"] != "value" {
		t.Errorf("body field = %v, want value", gotBody["field"])
	}
	if result.Status != "updated" {
		t.Errorf("result status = %q, want updated", result.Status)
	}
}

func TestRequestorDelete(t *testing.T) {
	var gotMethod string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	cfg := NewConfig(testCommerceCode, testAPIKey, server.URL, )
	cfg.Headers = map[string]string{
		"Tbk-Api-Key-Id":     testCommerceCode,
		"Tbk-Api-Key-Secret": testAPIKey,
	}
	requestor := NewRequestor(&cfg)

	if err := requestor.Delete("/resource/123", map[string]any{"id": "123"}, nil); err != nil {
		t.Fatalf("Delete returned error: %v", err)
	}
	if gotMethod != http.MethodDelete {
		t.Errorf("method = %q, want DELETE", gotMethod)
	}
}

func TestRequestorJSONDecodeFailure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"broken`))
	}))
	defer server.Close()

	cfg := NewConfig(testCommerceCode, testAPIKey, server.URL, )
	requestor := NewRequestor(&cfg)

	var result struct{ Token string }
	err := requestor.Get("/resource", &result)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var tbErr *transbank.TransportError
	if !errors.As(err, &tbErr) {
		t.Fatalf("error type = %T, want *transbank.TransportError", err)
	}
	if tbErr.Message != "decoding response" {
		t.Errorf("Message = %q, want decoding response", tbErr.Message)
	}
}

func TestRequestorEmptyBodyWithResult(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	cfg := NewConfig(testCommerceCode, testAPIKey, server.URL, )
	requestor := NewRequestor(&cfg)

	var result struct{ Token string }
	err := requestor.Get("/resource", &result)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var tbErr *transbank.TransportError
	if !errors.As(err, &tbErr) {
		t.Fatalf("error type = %T, want *transbank.TransportError", err)
	}
	if tbErr.Message != "decoding response" {
		t.Errorf("Message = %q, want decoding response", tbErr.Message)
	}
}

func TestRequestorEmptyBodyWithoutResult(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	cfg := NewConfig(testCommerceCode, testAPIKey, server.URL, )
	requestor := NewRequestor(&cfg)

	if err := requestor.Get("/resource", nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRequestorHTTPStatusBoundary299(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(299)
		w.Write([]byte(`{"token":"tok299"}`))
	}))
	defer server.Close()

	cfg := NewConfig(testCommerceCode, testAPIKey, server.URL, )
	requestor := NewRequestor(&cfg)

	var result struct{ Token string }
	if err := requestor.Get("/resource", &result); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Token != "tok299" {
		t.Errorf("Token = %q, want tok299", result.Token)
	}
}

func TestRequestorTransportErrorUnwrap(t *testing.T) {
	cfg := NewConfig(testCommerceCode, testAPIKey, "http://127.0.0.1:1", )
	requestor := NewRequestor(&cfg)

	err := requestor.Get("/resource", nil)
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

func TestNewConfigDefaultTimeout(t *testing.T) {
	cfg := NewConfig(testCommerceCode, testAPIKey, "https://example.com", )
	httpClient, ok := cfg.HTTP.(*http.Client)
	if !ok {
		t.Fatalf("HTTP type = %T, want *http.Client", cfg.HTTP)
	}
	if httpClient.Timeout != 30*time.Second {
		t.Errorf("Timeout = %v, want 30s", httpClient.Timeout)
	}
}

func TestSetBaseURL(t *testing.T) {
	cfg := NewConfig(testCommerceCode, testAPIKey, "https://old.example.com", )
	cfg.SetBaseURL("https://new.example.com")
	if cfg.BaseURL != "https://new.example.com" {
		t.Errorf("BaseURL = %q, want https://new.example.com", cfg.BaseURL)
	}
}
