package main

import (
	"bytes"
	"errors"
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
		opts *verifyOpts
		in   io.Reader
		args []string
		err  string
	}{
		{
			name: "invalid stdin",
			opts: &verifyOpts{
				publicKeyFile: "testdata/id_es256.pub.jwk",
			},
			in:  strings.NewReader("this isn't JSON"),
			err: "code=400, message=error unmarshaling JSON: while decoding JSON: json: cannot unmarshal string into Go value of type gobl.Envelope",
		},
		{
			name: "read error",
			opts: &verifyOpts{
				publicKeyFile: "testdata/id_es256.pub.jwk",
			},
			in:  testy.ErrorReader(`{"foo":`, errors.New("read error")),
			err: "code=400, message=read error",
		},
		{
			name: "empty envelope",
			opts: &verifyOpts{
				publicKeyFile: "testdata/id_es256.pub.jwk",
			},
			in:  strings.NewReader(`{}`),
			err: "code=422, message=$schema: cannot be blank; doc: cannot be blank; head: cannot be blank.",
		},
		{
			name: "read from file",
			opts: &verifyOpts{
				publicKeyFile: "testdata/id_es256.pub.jwk",
			},
			args: []string{"testdata/invalid.json"},
			err:  "code=400, message=error unmarshaling JSON: while decoding JSON: json: cannot unmarshal string into Go value of type gobl.Envelope",
		},
		{
			name: "file missing",
			args: []string{"asdf"},
			err:  "open asdf: no such file or directory",
		},
		{
			name: "explicit stdin",
			opts: &verifyOpts{
				publicKeyFile: "testdata/id_es256.pub.jwk",
			},
			args: []string{"-"},
			in:   strings.NewReader("not really valid"),
			err:  "code=400, message=error unmarshaling JSON: while decoding JSON: json: cannot unmarshal string into Go value of type gobl.Envelope",
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
			opts := tt.opts
			if opts == nil {
				opts = &verifyOpts{}
			}
			err := opts.runE(c, tt.args)
			if tt.err != "" {
				assert.EqualError(t, err, tt.err)
			} else {
				assert.Nil(t, err)
			}
			if d := testy.DiffText(testy.Snapshot(t), buf.String()); d != nil {
				t.Error(d)
			}
		})
	}
}
