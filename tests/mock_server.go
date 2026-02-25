package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
)

type MockServer struct {
	Server      *httptest.Server
	Response    any
	StatusCode  int
	LastRequest *http.Request
}

func NewMockServer() *MockServer {
	ms := &MockServer{
		StatusCode: http.StatusOK,
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ms.LastRequest = r // Guardamos el request por si queremos validar headers luego

		w.Header().Set("Content-Type", "application/json")

		// 1. Validar que el status code sea legal para net/http
		// Si es -1 (error local), el server no debería ni responder,
		// pero para evitar pánicos, forzamos un 500 o simplemente ignoramos.
		actualStatus := ms.StatusCode
		if actualStatus < 100 || actualStatus > 599 {
			actualStatus = http.StatusInternalServerError
		}

		w.WriteHeader(actualStatus)

		// 2. Encodear la respuesta
		if ms.Response != nil {
			json.NewEncoder(w).Encode(ms.Response)
		}
	}))

	ms.Server = server
	return ms
}

func (ms *MockServer) Close() {
	ms.Server.Close()
}

func (ms *MockServer) URL() string {
	return ms.Server.URL
}
