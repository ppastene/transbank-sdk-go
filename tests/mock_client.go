package tests

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type mockClient struct {
	mockServerURL string
}

func (c *mockClient) Request(method string, originalUrl string, headers map[string]string, payload any) ([]byte, int, error) {
	var bodyReader io.Reader
	if payload != nil {
		jsonData, err := json.Marshal(payload)
		if err != nil {
			return nil, 0, err
		}
		bodyReader = bytes.NewReader(jsonData)
	}

	req, err := http.NewRequest(method, c.mockServerURL, bodyReader)
	if err != nil {
		return nil, 0, err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	// 5. Leemos la respuesta "fake" que configuramos en el Mock
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, err
	}

	return respBody, resp.StatusCode, nil
}
