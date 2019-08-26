package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func setup(t *testing.T) (*os.File, func()) {
	f, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatalf("could not create temp file: %v", err)
	}
	teardown := func() {
		os.Remove(f.Name())
	}
	return f, teardown
}

func TestReadFile(t *testing.T) {
	file, teardown := setup(t)
	fmt.Println(file.Name())
	defer teardown()
	d1 := []byte("hello\ntest\n")
	n, err := file.Write(d1)
	if err != nil {
		t.Errorf("Could not write to file: %v", err)
	}

	fmt.Println(n)

	lines, err := readFile(file)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	time.Sleep(10000)
	fmt.Println(lines)

}
