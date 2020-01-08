package main

import (
	"net/url"
	"strconv"
	"strings"
)

// Checks if string is valid URL
func isValidURL(urlStr string) bool {
	u, err := url.Parse(urlStr)
	return err == nil && u.Scheme != "" && u.Host != ""
}

// Split string on colon and return a slice
func splitHeader(header string) (int, []string) {
	splitted := strings.Split(header, ":")
	return len(splitted), splitted

}

// Determine if tls
func isTLS(baseURL string) (bool, error) {
	scheme, err := url.Parse(baseURL)

	if err != nil {
		return false, err
	}

	if scheme.Scheme == "http" {
		return false, nil
	}

	if scheme.Scheme == "" {
		return false, nil
	}

	return true, nil
}

func stringToFloat(input string) float64 {
	if s, err := strconv.ParseFloat(input, 64); err == nil {
		return s
	}
	return 0.00
}
