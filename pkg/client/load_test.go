package client

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLoadCoordinate(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("ok"))
	}))
	defer srv.Close()

	cass := Cassowary{
		BaseURL:               srv.URL,
		ConcurrencyLevel:      1,
		Requests:              10,
		DisableTerminalOutput: true,
	}

	metrics, err := cass.Coordinate()
	if err != nil {
		t.Error(err)
	}

	if metrics.BaseURL != srv.URL {
		t.Errorf("Wanted %s but got %s", srv.URL, metrics.BaseURL)
	}

	if metrics.TotalRequests != 10 {
		t.Errorf("Wanted %d but got %d", 1, metrics.TotalRequests)
	}

	if metrics.FailedRequests != 0 {
		t.Errorf("Wanted %d but got %d", 0, metrics.FailedRequests)
	}
}

func TestLoadCoordinateURLPaths(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("ok"))
	}))
	defer srv.Close()

	cass := Cassowary{
		BaseURL:               srv.URL,
		ConcurrencyLevel:      1,
		Requests:              30,
		FileMode:              true,
		URLPaths:              []string{"/get_user", "/get_accounts", "/get_orders"},
		DisableTerminalOutput: true,
	}
	metrics, err := cass.Coordinate()
	if err != nil {
		t.Error(err)
	}
	if metrics.BaseURL != srv.URL {
		t.Errorf("Wanted %s but got %s", srv.URL, metrics.BaseURL)
	}

	if metrics.TotalRequests != 30 {
		t.Errorf("Wanted %d but got %d", 1, metrics.TotalRequests)
	}

	if metrics.FailedRequests != 0 {
		t.Errorf("Wanted %d but got %d", 0, metrics.FailedRequests)
	}

	if len(cass.URLPaths) != 30 {
		t.Errorf("Wanted %d but got %d", 30, len(cass.URLPaths))
	}
}

func TestCoordinateTLSConfig(t *testing.T) {
	srv := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("ok"))
	}))

	pemCerts, err := ioutil.ReadFile("testdata/ca.pem")
	if err != nil {
		t.Fatal("Invalid ca.pem path")
	}

	ca := x509.NewCertPool()
	if !ca.AppendCertsFromPEM(pemCerts) {
		t.Fatal("Failed to read CA from PEM")
	}

	cert, err := tls.LoadX509KeyPair("testdata/server.pem", "testdata/server-key.pem")
	if err != nil {
		t.Fatal("Invalid server.pem/server-key.pem path")
	}

	srv.TLS = &tls.Config{
		ClientCAs:    ca,
		ClientAuth:   tls.RequireAndVerifyClientCert,
		Certificates: []tls.Certificate{cert},
	}
	srv.StartTLS()

	cert, err = tls.LoadX509KeyPair("testdata/client.pem", "testdata/client-key.pem")
	if err != nil {
		t.Fatal("Invalid client.pem/client-key.pem path")
	}
	clientTLSConfig := &tls.Config{
		RootCAs:      ca,
		Certificates: []tls.Certificate{cert},
	}

	cass := Cassowary{
		BaseURL:               srv.URL,
		ConcurrencyLevel:      1,
		Requests:              10,
		TLSConfig:             clientTLSConfig,
		DisableTerminalOutput: true,
	}

	metrics, err := cass.Coordinate()
	if err != nil {
		t.Error(err)
	}

	if metrics.BaseURL != srv.URL {
		t.Errorf("Wanted %s but got %s", srv.URL, metrics.BaseURL)
	}

	if metrics.TotalRequests != 10 {
		t.Errorf("Wanted %d but got %d", 1, metrics.TotalRequests)
	}

	if metrics.FailedRequests != 0 {
		t.Errorf("Wanted %d but got %d", 0, metrics.FailedRequests)
	}
}
