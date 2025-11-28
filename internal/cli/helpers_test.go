package cli

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"testing"

	"github.com/invopop/gobl"
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

func testFileReaderForDoc(t *testing.T, filename string) io.Reader {
	t.Helper()
	r := testFileReader(t, filename)
	dec := json.NewDecoder(r)
	var env gobl.Envelope
	if err := dec.Decode(&env); err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	if err := enc.Encode(env.Document); err != nil {
		t.Fatal(err)
	}

	return &buf
}
