package main

import (
	"io"
	"os"
	"testing"
)

func testFileReader(t *testing.T, filename string) io.Reader {
	t.Helper()
	f, err := os.Open(filename)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = f.Close()
	})
	return f
}
