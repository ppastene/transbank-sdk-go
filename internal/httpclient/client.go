package httpclient

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"
)

type DefaultClient struct {
	httpClient *http.Client
}

func NewDefaultClient() *DefaultClient {
	return &DefaultClient{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c DefaultClient) Request(method string, url string, headers map[string]string, payload any) ([]byte, int, error) {
	method = strings.ToUpper(method)

	var bodyReader io.Reader

	if payload != nil {
		jsonData, err := json.Marshal(payload)
		if err != nil {
			return nil, 0, err
		}
		bodyReader = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, 0, err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, 0, err
	}

	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}

	return respBody, resp.StatusCode, nil
}
