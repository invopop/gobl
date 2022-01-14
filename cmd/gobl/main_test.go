package main

import (
	"bytes"
	"io"
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
	}{}

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
			err := verify(nil, tt.args)
			assert.EqualError(t, err, tt.err)
			if t := testy.DiffText(tt.Snapshot(t), buf.String()); d != nil {
				t.Error(d)
			}
		})
	}
}
