package internal

import (
	"errors"
	"strings"
	"testing"

	"github.com/ppastene/transbank-sdk-go"
)

func validationErr(t *testing.T, err error) *transbank.ValidationError {
	t.Helper()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var tbErr *transbank.ValidationError
	if !errors.As(err, &tbErr) {
		t.Fatalf("error type = %T, want *transbank.ValidationError", err)
	}
	return tbErr
}

func TestValidateURLParam(t *testing.T) {
	valid64 := strings.Repeat("a", 64)
	valid26 := strings.Repeat("b", 26)
	tests := []struct {
		name    string
		param   string
		value   string
		maxLen  int
		wantErr bool
	}{
		{name: "token valid", param: "token", value: valid64, maxLen: 64, wantErr: false},
		{name: "token empty", param: "token", value: "", maxLen: 64, wantErr: true},
		{name: "token 65 chars", param: "token", value: valid64 + "a", maxLen: 64, wantErr: true},
		{name: "token with slash", param: "token", value: "abc/def", maxLen: 64, wantErr: true},
		{name: "token with question mark", param: "token", value: "abc?def", maxLen: 64, wantErr: true},
		{name: "token with ampersand", param: "token", value: "abc&def", maxLen: 64, wantErr: true},
		{name: "token with hash", param: "token", value: "abc#def", maxLen: 64, wantErr: true},
		{name: "token with equals", param: "token", value: "abc=def", maxLen: 64, wantErr: true},
		{name: "token with spaces", param: "token", value: "abc def", maxLen: 64, wantErr: true},
		{name: "buy_order 26 chars", param: "buy_order", value: valid26, maxLen: 26, wantErr: false},
		{name: "buy_order 27 chars", param: "buy_order", value: valid26 + "b", maxLen: 26, wantErr: true},
		{name: "buy_order empty", param: "buy_order", value: "", maxLen: 26, wantErr: true},
		{name: "buy_order with dot", param: "buy_order", value: "abc.def", maxLen: 26, wantErr: true},
		{name: "short param valid", param: "id", value: "abc", maxLen: 10, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateURLParam(tt.param, tt.value, tt.maxLen)
			if tt.wantErr {
				if err := validationErr(t, err); err.Message == "" {
					t.Error("expected non-empty validation message")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}
