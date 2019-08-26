package main

import (
	"bufio"
	"io"
)

func readFile(f io.ReadWriteCloser) ([]string, error) {
	var urlLines []string
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		urlLines = append(urlLines, scanner.Text())
	}
	return urlLines, nil
}
