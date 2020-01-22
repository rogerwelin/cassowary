package client

import (
	"net/url"
	"strconv"
	"strings"
)

// IsValidURL checks if string is valid URL
func IsValidURL(urlStr string) bool {
	u, err := url.Parse(urlStr)
	return err == nil && u.Scheme != "" && u.Host != ""
}

// SplitHeader splits string on colon and return a slice
func SplitHeader(header string) (int, []string) {
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

func generateSuffixes(src []string, length int) []string {
	if len(src) > length {
		return src
	}
	var urls []string
	srcLength := len(src)
	for i := 0; i < length; i++ {
		urls = append(urls, src[i%srcLength])
	}
	return urls
}
