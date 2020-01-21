package client

import (
	"reflect"
	"testing"
)

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

var testSuffixes = []struct {
	in       []string
	length   int
	expected []string
}{
	{[]string{"uri1", "uri2"}, 5, []string{"uri1", "uri2", "uri1", "uri2", "uri1"}},
	{[]string{"uri1", "uri2", "uri3", "uri4", "uri5"}, 3, []string{"uri1", "uri2", "uri3", "uri4", "uri5"}},
	{[]string{"ab", "ac"}, 4, []string{"ab", "ac", "ab", "ac"}},
}

func TestValidURL(t *testing.T) {
	for i, tt := range testURL {
		actual := IsValidURL(tt.in)
		if actual != tt.expected {
			t.Errorf("test: %d, isValidURL(%s): expected %t, actual %t", i+1, tt.in, tt.expected, actual)
		}
	}
}

func TestSplitHeaders(t *testing.T) {
	for i, tt := range testHeader {
		actual, _ := SplitHeader(tt.in)
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

func TestSuffixes(t *testing.T) {
	for i, tt := range testSuffixes {
		actual := generateSuffixes(tt.in, tt.length)
		if !reflect.DeepEqual(actual, tt.expected) {
			t.Errorf("test: %d, generateSuffixes(%v,%d): expected %v, actual %v", i, tt.in, tt.length, tt.expected, actual)
		}
	}
}
