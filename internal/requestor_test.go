package internal

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

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

	cfg := NewConfig(testCommerceCode, testAPIKey, server.URL, false)
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

	cfg := NewConfig(testCommerceCode, testAPIKey, server.URL, false)
	requestor := NewRequestor(&cfg)

	err := requestor.Get("/transactions", nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	tbErr, ok := err.(*transbank.HTTPError)
	if !ok {
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
	cfg := NewConfig(testCommerceCode, testAPIKey, "https://webpay3gint.transbank.cl", false)
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
