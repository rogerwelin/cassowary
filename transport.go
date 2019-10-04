package main

import (
	"net/http"
	"net/http/httptrace"
	"time"
)

// rounTripTrace holds timings for a single HTTP roundtrip
type roundTripTrace struct {
	tls           bool
	start         time.Time
	dnsDone       time.Time
	connectDone   time.Time
	gotConn       time.Time
	responseStart time.Time
	end           time.Time
}

// transport is a custom tansport keeping traces for each HTTP roundtrip
type transport struct {
	Transport http.RoundTripper
	traces    []*roundTripTrace
	current   *roundTripTrace
}

func newTransport(rt http.RoundTripper) *transport {
	return &transport{
		Transport: rt,
		traces:    []*roundTripTrace{},
	}
}

// RoundTrip switches to a new trace, then runs embedded RoundTripper
func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	trace := &roundTripTrace{}
	if req.URL.Scheme == "https" {
		trace.tls = true
	}
	t.current = trace
	t.traces = append(t.traces, trace)
	return t.Transport.RoundTrip(req)
}

func (t *transport) DNSStart(_ httptrace.DNSStartInfo) {
	t.current.start = time.Now()
}

func (t *transport) DNSDone(_ httptrace.DNSDoneInfo) {
	t.current.dnsDone = time.Now()
}

func (ts *transport) ConnectStart(_, _ string) {
	t := ts.current
	// No DNS resolution because we connected to IP directly.
	if t.dnsDone.IsZero() {
		t.start = time.Now()
		t.dnsDone = t.start
	}
}

func (t *transport) ConnectDone(net, addr string, err error) {
	t.current.connectDone = time.Now()
}

func (t *transport) GotConn(_ httptrace.GotConnInfo) {
	t.current.gotConn = time.Now()
}

func (t *transport) GotFirstResponseByte() {
	t.current.responseStart = time.Now()
}
