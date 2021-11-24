package client

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadCoordinate(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("ok"))
	}))
	defer srv.Close()

	cass := Cassowary{
		BaseURL: srv.URL,
		Groups: []QueryGroup{
			{
				Name:             "default",
				ConcurrencyLevel: 1,
				Requests:         10,
			},
		},
		DisableTerminalOutput: true,
	}

	metrics, metricsGroups, err := cass.Coordinate()
	if err != nil {
		t.Error(err)
	}

	if metrics.BaseURL != srv.URL {
		t.Errorf("Wanted %s but got %s", srv.URL, metrics.BaseURL)
	}

	if metrics.TotalRequests != 10 {
		t.Errorf("Wanted %d but got %d", 10, metrics.TotalRequests)
	}

	if metrics.FailedRequests != 0 {
		t.Errorf("Wanted %d but got %d", 0, metrics.FailedRequests)
	}

	metricsGroupsWant := map[string]ResultMetrics{
		"default": {
			Name: "default", BaseURL: srv.URL, TotalRequests: metrics.TotalRequests, FailedRequests: 0,
			RespSuccess: map[string]int{"200": metrics.TotalRequests},
			RespFailed:  map[string]int{},
			RespSize:    stats{Min: 2, Max: 2, Mean: 2, Median: 2, P95: 2, P99: 2},
			// flapping values
			RequestsPerSecond: metrics.RequestsPerSecond,
			ElapsedStats:      metrics.ElapsedStats,
			TCPStats:          metrics.TCPStats,
			ProcessingStats:   metrics.ProcessingStats,
			ContentStats:      metrics.ContentStats,
		},
	}
	assert.EqualValues(t, metricsGroupsWant, metricsGroups)
}

func TestLoadCoordinateURLPaths(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("ok"))
	}))
	defer srv.Close()

	cass := Cassowary{
		BaseURL: srv.URL,
		Groups: []QueryGroup{
			{
				Name:             "default",
				ConcurrencyLevel: 1,
				Requests:         30,
				FileMode:         true,
				URLPaths:         []string{"/get_user", "/get_accounts", "/get_orders"},
			},
		},
		DisableTerminalOutput: true,
	}
	metrics, _, err := cass.Coordinate()
	if err != nil {
		t.Error(err)
	}
	if metrics.BaseURL != srv.URL {
		t.Errorf("Wanted %s baseURL, but got %s", srv.URL, metrics.BaseURL)
	}

	if metrics.TotalRequests != 30 {
		t.Errorf("Wanted %d total, but got %d", 30, metrics.TotalRequests)
	}

	if metrics.FailedRequests != 0 {
		t.Errorf("Wanted %d failed, but got %d", 0, metrics.FailedRequests)
	}
}

type URLIterator struct {
	pos  uint64
	data []string
	v    Validator
}

func (it *URLIterator) Next() *Query {
	for {
		pos := atomic.AddUint64(&it.pos, 1)
		if pos > uint64(len(it.data)) {
			if !atomic.CompareAndSwapUint64(&it.pos, pos, 0) {
				// retry
				continue
			}
			pos = 0
		} else {
			pos--
		}
		//return &Query{Method: "GET", URL: it.data[pos]}
		return &Query{Method: "POST", URL: it.data[pos], DataType: "application/json", Data: []byte("{ \"test\": \"POST\" }"), Validator: it.v}
	}
}

func NewURLIterator(data []string) *URLIterator {
	if len(data) == 0 {
		return nil
	}
	return &URLIterator{data: data, pos: 0}
}

func TestLoadCoordinateURLIterator(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("ok"))
	}))
	defer srv.Close()

	it := NewURLIterator([]string{"/test1", "/test2", "/test3"})

	cass := Cassowary{
		BaseURL: srv.URL,
		Groups: []QueryGroup{
			{
				Name:             "group_1",
				ConcurrencyLevel: 2,
				Requests:         32,
				FileMode:         true,
				URLIterator:      it,
			},
			{
				Name:             "group_2",
				ConcurrencyLevel: 1,
				Requests:         10,
				FileMode:         true,
				URLIterator:      it,
			},
		},
		DisableTerminalOutput: true,
	}
	metrics, metricsGroups, err := cass.Coordinate()
	if err != nil {
		t.Fatal(err)
	}

	if metrics.TotalRequests != cass.Groups[0].Requests+cass.Groups[1].Requests {
		t.Errorf("Wanted %d total, but got %d", cass.Groups[0].Requests+cass.Groups[1].Requests, metrics.TotalRequests)
	}

	if metrics.FailedRequests != 0 {
		t.Errorf("Wanted %d failed, but got %d", 0, metrics.FailedRequests)
	}

	if len(metrics.RespSuccess) != 1 {
		t.Errorf("Wanted 1 total failed status code, but got %d", len(metrics.RespSuccess))
	}
	if count, ok := metrics.RespSuccess["200"]; ok {
		if count != metrics.TotalRequests {
			t.Errorf("Wanted %d total success status code, but got %d", metrics.TotalRequests, count)
		}
	} else {
		t.Errorf("Wanted %d total success status code, but got %d", metrics.TotalRequests, 0)
	}
	if len(metrics.RespFailed) != 0 {
		t.Errorf("Wanted 0 total failed, but got %d", len(metrics.RespFailed))
	}

	metricsGroupsWant := map[string]ResultMetrics{
		"group_1": {
			Name: "group_1", BaseURL: srv.URL, TotalRequests: 32, FailedRequests: 0,
			BodySize:    metrics.BodySize,
			RespSuccess: map[string]int{"200": 32},
			RespFailed:  map[string]int{},
			RespSize:    stats{Min: 2, Max: 2, Mean: 2, Median: 2, P95: 2, P99: 2},
			// flapping values
			RequestsPerSecond: metricsGroups["group_1"].RequestsPerSecond,
			ElapsedStats:      metricsGroups["group_1"].ElapsedStats,
			TCPStats:          metricsGroups["group_1"].TCPStats,
			ProcessingStats:   metricsGroups["group_1"].ProcessingStats,
			ContentStats:      metricsGroups["group_1"].ContentStats,
		},
		"group_2": {
			Name: "group_2", BaseURL: srv.URL, TotalRequests: 10, FailedRequests: 0,
			BodySize:    metrics.BodySize,
			RespSuccess: map[string]int{"200": 10},
			RespFailed:  map[string]int{},
			RespSize:    stats{Min: 2, Max: 2, Mean: 2, Median: 2, P95: 2, P99: 2},
			// flapping values
			RequestsPerSecond: metricsGroups["group_2"].RequestsPerSecond,
			ElapsedStats:      metricsGroups["group_2"].ElapsedStats,
			TCPStats:          metricsGroups["group_2"].TCPStats,
			ProcessingStats:   metricsGroups["group_2"].ProcessingStats,
			ContentStats:      metricsGroups["group_2"].ContentStats,
		},
	}
	assert.EqualValues(t, metricsGroupsWant, metricsGroups)
}

func respHTTPValidator(statusCode int, respSize int64, resp []byte, err error) (bool, string) {
	if err != nil {
		return true, err.Error()
	} else if statusCode > 226 {
		return true, strconv.Itoa(statusCode)
	} else {
		return false, strconv.Itoa(statusCode)
	}
}

func NewURLIteratorWithValidator(data []string, v Validator) *URLIterator {
	if len(data) == 0 {
		return nil
	}
	return &URLIterator{data: data, pos: 0, v: v}
}

func TestLoadCoordinateURLIteratorWithValidator(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(503)
		w.Write([]byte("error"))
	}))
	defer srv.Close()

	it := NewURLIteratorWithValidator([]string{"/test1", "/test2", "/test3"}, respHTTPValidator)

	cass := Cassowary{
		BaseURL: srv.URL,
		Groups: []QueryGroup{
			{
				Name:             "default",
				ConcurrencyLevel: 1,
				Requests:         32,
				FileMode:         true,
				URLIterator:      it,
			},
		},
		DisableTerminalOutput: true,
	}
	metrics, metricsGroups, err := cass.Coordinate()
	if err != nil {
		t.Error(err)
	}

	if metrics.TotalRequests != cass.Groups[0].Requests {
		t.Errorf("Wanted %d total, but got %d", cass.Groups[0].Requests, metrics.TotalRequests)
	}

	if metrics.FailedRequests != metrics.TotalRequests {
		t.Errorf("Wanted %d failed, but got %d", 0, metrics.FailedRequests)
	}

	if len(metrics.RespSuccess) != 0 {
		t.Errorf("Wanted 0 total failed status code, but got %d", len(metrics.RespSuccess))
	}
	if len(metrics.RespFailed) != 1 {
		t.Errorf("Wanted 1 total failed status code, but got %d", len(metrics.RespFailed))
	}
	if count, ok := metrics.RespFailed["503"]; ok {
		if count != metrics.FailedRequests {
			t.Errorf("Wanted %d total failed status code, but got %d", metrics.FailedRequests, count)
		}
	} else {
		t.Errorf("Wanted %d total failed status code, but got %d", metrics.FailedRequests, 0)
	}

	metricsGroupsWant := map[string]ResultMetrics{
		"default": {
			Name: "default", BaseURL: srv.URL, TotalRequests: metrics.TotalRequests, FailedRequests: metrics.TotalRequests,
			BodySize:    metrics.BodySize,
			RespSuccess: map[string]int{},
			RespFailed:  map[string]int{"503": metrics.TotalRequests},
			// flapping values
			RequestsPerSecond: metrics.RequestsPerSecond,
			ElapsedStats:      metrics.ElapsedStats,
			TCPStats:          metrics.TCPStats,
			ProcessingStats:   metrics.ProcessingStats,
			ContentStats:      metrics.ContentStats,
		},
	}
	assert.EqualValues(t, metricsGroupsWant, metricsGroups)
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
		BaseURL: srv.URL,
		Groups: []QueryGroup{
			{
				ConcurrencyLevel: 1,
				Requests:         10,
			},
		},
		TLSConfig:             clientTLSConfig,
		DisableTerminalOutput: true,
	}

	metrics, _, err := cass.Coordinate()
	if err != nil {
		t.Error(err)
	}

	if metrics.BaseURL != srv.URL {
		t.Errorf("Wanted %s but got %s", srv.URL, metrics.BaseURL)
	}

	if metrics.TotalRequests != 10 {
		t.Errorf("Wanted %d but got %d", 10, metrics.TotalRequests)
	}

	if metrics.FailedRequests != 0 {
		t.Errorf("Wanted %d but got %d", 0, metrics.FailedRequests)
	}
}
