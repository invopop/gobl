package internal

import (
	"context"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/flimzy/testy"
)

func TestVerify(t *testing.T) {
	type tt struct {
		in  io.Reader
		err string
	}

	tests := testy.NewTable()
	tests.Add("validation pass", func(t *testing.T) interface{} {
		f, err := os.Open("testdata/success.json")
		if err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() { _ = f.Close() })

		return tt{
			in: f,
		}
	})

	tests.Run(t, func(t *testing.T, tt tt) {
		t.Parallel()
		err := Verify(context.Background(), tt.in)
		if tt.err == "" {
			assert.Nil(t, err)
		} else {
			assert.EqualError(t, err, tt.err)
		}

	})
}
