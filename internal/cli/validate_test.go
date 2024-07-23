package cli

import (
	"context"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/flimzy/testy"
)

func TestValidate(t *testing.T) {
	type tt struct {
		in  io.Reader
		err string
	}

	tests := testy.NewTable()
	tests.Add("valid document", func(t *testing.T) interface{} {
		return tt{
			in: testFileReader(t, "testdata/success.json"),
		}
	})
	tests.Add("invalid type", tt{
		in: strings.NewReader(`{
				"$schema": "https://gobl.org/draft-0/envelope",
				"head": {
					"uuid": "9d8eafd5-77be-11ec-b485-5405db9a3e49",
					"dig": {
						"alg": "sha256",
						"val": "dce3bc3c8bf28f3d209f783917b3082ddc0339a66e9ba3aa63849e4357db1422"
					}
				},
				doc: {
					"$schema": "https://example.com/duck",
					"walk": "like a duck",
					"talk": "like a duck",
					"look": "like a duck"
				}
			}`),
		err: `code=400, message=unknown-schema`,
	})
	tests.Add("with signature", func(t *testing.T) interface{} {
		return tt{
			in: testFileReader(t, "testdata/signed.json"),
		}
	})
	tests.Add("draft", func(t *testing.T) interface{} {
		return tt{
			in: testFileReader(t, "testdata/draft.json"),
		}
	})
	tests.Add("without envelope", func(t *testing.T) interface{} {
		return tt{
			in: testFileReader(t, "testdata/invoice.json"),
		}
	})
	tests.Add("without totals", func(t *testing.T) interface{} {
		return tt{
			in:  testFileReader(t, "testdata/nototals.json"),
			err: "code=422, message=doc: (totals: cannot be blank.).",
		}
	})

	tests.Run(t, func(t *testing.T, tt tt) {
		t.Parallel()
		err := Validate(context.Background(), tt.in)
		if tt.err == "" {
			assert.Nil(t, err)
		} else {
			assert.EqualError(t, err, tt.err)
		}
	})
}
