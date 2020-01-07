package main

import "testing"

var testURL = []struct {
	in       string
	expected bool
}{
	{"google.com", false},
	{"www.google.com", false},
	{"https://www.google.com", true},
	{"/foo/bar", false},
	{"http://localhost:8000", true},
}

var testHeader = []struct {
	in       string
	expected int
}{
	{"X-Forwarded-For: www.google.com", 2},
	{"X-Forwarded-For:www.google.com", 2},
	{"X-Forwarded-For www.google.com", 1},
	{"X-Forwarded-For", 1},
}

var testTLSScheme = []struct {
	in       string
	expected bool
}{
	{"http://localhost", false},
	{"https://localhost", true},
	{"https:/localhost", true},
	{"localhost", false},
	{"HTTP://localhost", false},
	{"HTTPS://localhost", true},
}

func TestValidURL(t *testing.T) {
	for i, tt := range testURL {
		actual := isValidURL(tt.in)
		if actual != tt.expected {
			t.Errorf("test: %d, isValidURL(%s): expected %t, actual %t", i+1, tt.in, tt.expected, actual)
		}
	}
}

func TestSplitHeaders(t *testing.T) {
	for i, tt := range testHeader {
		actual, _ := splitHeader(tt.in)
		if actual != tt.expected {
			t.Errorf("test: %d, splitHeader(%s): expected %d, actual %d", i+1, tt.in, tt.expected, actual)
		}
	}
}

func TestTLSScheme(t *testing.T) {
	for i, tt := range testTLSScheme {
		actual, _ := isTLS(tt.in)
		if actual != tt.expected {
			t.Errorf("test: %d, isTLS(%s): expected %t, actual %t", i+1, tt.in, tt.expected, actual)
		}
	}
}
