package mock

import (
	"context"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
)

const ACCESS_TOKEN = "abc123"

type MockServer struct {
	Server       *httptest.Server
	mux          *http.ServeMux
	testdataPath string
}

func checkAuthorization(r *http.Request) bool {
	return r.URL.Query().Get("apikey") == ACCESS_TOKEN
}

func getTestDataPath() string {
	dir, _ := os.Getwd()
	parts := strings.Split(dir, string(filepath.Separator))
	i := 0
	for i = 0; i < len(parts); i++ {
		if parts[i] == "climacell-exporter" {
			break
		}
	}

	return fmt.Sprintf("%c%s", filepath.Separator, filepath.Join(parts[0:i+1]...))
}

func NewMockServer() *MockServer {
	m := &MockServer{}
	m.mux = http.NewServeMux()
	m.mux.Handle("/v3/weather/realtime", m.checkAuthMiddleware(m.serveFile("realtime.json")))

	m.testdataPath = filepath.Join(getTestDataPath(), "testdata")

	m.Server = httptest.NewTLSServer(m.mux)
	return m
}

func (m *MockServer) Close() {
	m.Server.Close()
}

func (m *MockServer) Client() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, network, _ string) (net.Conn, error) {
				return net.Dial(network, m.Server.Listener.Addr().String())
			},
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
}

func (m *MockServer) checkAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !checkAuthorization(r) {
			w.WriteHeader(403)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (m *MockServer) serveFile(path string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, _ := ioutil.ReadFile(filepath.Join(m.testdataPath, path))
		w.Header().Set("Content-Length", fmt.Sprintf("%d", len(data)))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte(data))
	})
}
