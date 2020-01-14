package main

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
			if r.Method == http.MethodDelete {
				w.WriteHeader(http.StatusAccepted)
				return
			}
			w.WriteHeader(http.StatusOK)
		}),
	)
	defer pgwOK.Close()

	cass := &cassowary{}
	cass.promURL = pgwOK.URL

	t1, t2, t3, t4, t5, t6, t7, t8, t9, t10, t11, t12 := 1.1, 1.2, 1.3, 1.4, 1.5, 1.6, 1.7, 1.8, 1.9, 1.10, 1.11, 1.12
	err := cass.pushPrometheusMetrics(t1, t2, t3, t4, t5, t6, t7, t8, t9, t10, t11, t12)
	if err != nil {
		t.Error("Got error but wanted OK")
	}
	if lastPath != "/metrics/job/cassowary_load_test/url/" {
		t.Errorf("Wanted %s but got %s", "/metrics/job/cassowary_load_test/url/", lastPath)
	}
	if lastMethod != "PUT" {
		t.Errorf("Wanted %s but got %s", "PUT", lastMethod)
	}
}
