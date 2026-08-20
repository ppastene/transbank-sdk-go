package transbank

import (
	"errors"
	"fmt"
	"testing"
)

func TestOptionsValidate(t *testing.T) {
	tests := []struct {
		name    string
		opts    Options
		wantErr bool
	}{
		{
			name: "valid integration",
			opts: Options{Environment: Integration, CommerceCode: "597055555532", ApiKey: "secret"},
		},
		{
			name: "valid production",
			opts: Options{Environment: Production, CommerceCode: "597055555532", ApiKey: "secret"},
		},
		{name: "invalid environment", opts: Options{Environment: Environment(99), CommerceCode: "597055555532", ApiKey: "secret"}, wantErr: true},
		{name: "empty commerce code", opts: Options{Environment: Integration, ApiKey: "secret"}, wantErr: true},
		{name: "non numeric commerce code", opts: Options{Environment: Integration, CommerceCode: "59705555553a", ApiKey: "secret"}, wantErr: true},
		{name: "empty api key", opts: Options{Environment: Integration, CommerceCode: "597055555532"}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				var tbErr *ValidationError
				if !errors.As(err, &tbErr) {
					t.Fatalf("error type = %T, want *transbank.ValidationError", err)
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestErrorString(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want string
	}{
		{
			name: "validation error",
			err:  &ValidationError{Message: "token must be 64 characters"},
			want: "token must be 64 characters",
		},
		{
			name: "transport error",
			err:  &TransportError{Message: "request failed", Err: errors.New("dial tcp: refused")},
			want: "request failed: dial tcp: refused",
		},
		{
			name: "http error with body",
			err:  &HTTPError{StatusCode: 401, Body: `{"error_message":"unauthorized"}`},
			want: `transbank: status 401: {"error_message":"unauthorized"}`,
		},
		{
			name: "http error without body",
			err:  &HTTPError{StatusCode: 500},
			want: "transbank: status 500",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.want {
				t.Errorf("Error() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestErrorUnwrap(t *testing.T) {
	cause := errors.New("dial tcp: refused")
	err := &TransportError{Message: "request failed", Err: cause}
	if got := err.Unwrap(); got != cause {
		t.Errorf("Unwrap() = %v, want %v", got, cause)
	}
	if !errors.Is(err, cause) {
		t.Error("errors.Is(err, cause) = false, want true")
	}
}

func TestTransportErrorUnwrapChain(t *testing.T) {
	cause := errors.New("connection refused")
	inner := &TransportError{Message: "request failed", Err: cause}
	outer := fmt.Errorf("outer: %w", inner)

	var tbErr *TransportError
	if !errors.As(outer, &tbErr) {
		t.Fatal("errors.As should find TransportError in chain")
	}
	if tbErr.Message != "request failed" {
		t.Errorf("Message = %q, want request failed", tbErr.Message)
	}
	if !errors.Is(outer, cause) {
		t.Error("errors.Is should find root cause in chain")
	}
}

func TestHTTPErrorBodyTrimming(t *testing.T) {
	tests := []struct {
		name string
		err  *HTTPError
		want string
	}{
		{
			name: "multiline body",
			err:  &HTTPError{StatusCode: 400, Body: "{\n  \"error\": \"bad\"\n}"},
			want: "transbank: status 400: {\n  \"error\": \"bad\"\n}",
		},
		{
			name: "whitespace only body",
			err:  &HTTPError{StatusCode: 401, Body: "  \n  "},
			want: "transbank: status 401: ",
		},
		{
			name: "body with leading trailing spaces",
			err:  &HTTPError{StatusCode: 500, Body: "  server error  "},
			want: "transbank: status 500: server error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.want {
				t.Errorf("Error() = %q, want %q", got, tt.want)
			}
		})
	}
}
