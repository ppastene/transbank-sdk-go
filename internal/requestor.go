package internal

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/ppastene/transbank-sdk-go"
)

const (
	headerContentType = "Content-Type"
	mimeJSON          = "application/json"
)

type Requestor struct {
	config *Config
}

func NewRequestor(config *Config) *Requestor {
	return &Requestor{config: config}
}

func (r *Requestor) Post(path string, body any, result any) error {
	return r.do(http.MethodPost, path, body, result)
}

func (r *Requestor) Get(path string, result any) error {
	return r.do(http.MethodGet, path, nil, result)
}

func (r *Requestor) Put(path string, body any, result any) error {
	return r.do(http.MethodPut, path, body, result)
}

func (r *Requestor) Delete(path string, body any, result any) error {
	return r.do(http.MethodDelete, path, body, result)
}

func (r *Requestor) do(method, path string, body any, result any) error {
	endpoint := r.config.BaseURL + path

	var reader io.Reader
	if body != nil {
		buf, err := json.Marshal(body)
		if err != nil {
			return &transbank.TransportError{Message: "encoding request body", Err: err}
		}
		reader = bytes.NewReader(buf)
	}

	req, err := http.NewRequest(method, endpoint, reader)
	if err != nil {
		return &transbank.TransportError{Message: "building request", Err: err}
	}

	for name, value := range r.config.Headers {
		req.Header.Set(name, value)
	}
	req.Header.Set(headerContentType, mimeJSON)

	resp, err := r.config.HTTP.Do(req)
	if err != nil {
		return &transbank.TransportError{Message: "request failed", Err: err}
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		msg, _ := io.ReadAll(resp.Body)
		return &transbank.HTTPError{
			StatusCode: resp.StatusCode,
			Body:       string(msg),
		}
	}

	if result != nil {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return &transbank.TransportError{Message: "decoding response", Err: err}
		}
	}

	return nil
}
