package client

import (
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

	if metrics.TotalRequests != 10 {
		t.Errorf("Wanted %d but got %d", 1, metrics.TotalRequests)
	}

	if metrics.FailedRequests != 0 {
		t.Errorf("Wanted %d but got %d", 0, metrics.FailedRequests)
	}
}
