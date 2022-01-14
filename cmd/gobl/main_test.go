package main

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"gitlab.com/flimzy/testy"
)

func Test_verify(t *testing.T) {
	tests := []struct {
		name string
		in   io.Reader
		args []string
		err  string
	}{
		{
			name: "invalid stdin",
			in:   strings.NewReader("this isn't JSON"),
			err:  "error unmarshaling JSON: json: cannot unmarshal string into Go value of type gobl.Envelope",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c := &cobra.Command{}
			if tt.in != nil {
				c.SetIn(tt.in)
			}
			buf := &bytes.Buffer{}
			c.SetOut(buf)
			err := verify(c, tt.args)
			assert.EqualError(t, err, tt.err)
			if d := testy.DiffText(testy.Snapshot(t), buf.String()); d != nil {
				t.Error(d)
			}
		})
	}
}
