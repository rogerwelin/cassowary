package client

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPromGwPush(t *testing.T) {
	var (
		lastMethod string
		lastPath   string
	)

	// Fake a Pushgateway that responds with 200
	pgwOK := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			lastMethod = r.Method
			var err error
			_, err = ioutil.ReadAll(r.Body)
			if err != nil {
				t.Fatal(err)
			}
			lastPath = r.URL.EscapedPath()
			w.Header().Set("Content-Type", `text/plain; charset=utf-8`)
			w.WriteHeader(http.StatusOK)
		}),
	)
	defer pgwOK.Close()

	cass := &Cassowary{}
	cass.PromURL = pgwOK.URL

	metrics := ResultMetrics{
		FailedRequests:    1,
		TotalRequests:     100,
		RequestsPerSecond: 100.10,
		TCPStats: stats{
			Mean:   10.0,
			Median: 10.0,
			P95:    10.0,
		},
		ProcessingStats: stats{
			Mean:   1.0,
			Median: 1.0,
			P95:    1.0,
		},
	}

	err := cass.PushPrometheusMetrics(metrics)
	if err != nil {
		t.Error(err)
	}
	if lastPath != "/metrics/job/cassowary_load_test/url/" {
		t.Errorf("Wanted %s but got %s", "/metrics/job/cassowary_load_test/url/", lastPath)
	}
	if lastMethod != "PUT" {
		t.Errorf("Wanted %s but got %s", "PUT", lastMethod)
	}
}
