package main

import (
	"bufio"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
)

func downloadPath(url string) (string, error) {
	path, err := os.Executable()
	if err != nil {
		return "", err
	}
	binDir := filepath.Dir(path)

	out, err := os.Create(binDir + "/load.txt")
	if err != nil {
		return "", err
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", err
	}

	return binDir + "/load.txt", nil
}

func readFile(file string) ([]byte, error) {
	fileContent, err := os.ReadFile(file)
	if err != nil {
		return []byte{}, err
	}
	return fileContent, nil
}

func readLocalRemoteFile(filePath string) ([]string, error) {
	var urlLines []string
	var err error

	regex, _ := regexp.Compile(`^(https?):\/\/[^\s\/$.?#].[^\s]*$`)
	isRemote := regex.MatchString(filePath)

	if isRemote {
		filePath, err = downloadPath(filePath)
		if err != nil {
			return []string{}, err
		}
	}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	reader := bufio.NewReader(file)

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		urlLines = append(urlLines, scanner.Text())
	}
	return urlLines, nil
}
