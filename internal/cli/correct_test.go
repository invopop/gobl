package cli

import (
	"context"
	"regexp"
	"testing"

	_ "github.com/invopop/gobl"
	"github.com/invopop/gobl/cal"
	"github.com/stretchr/testify/assert"
	"gitlab.com/flimzy/testy"
)

// These tests depend on the build_test.go for some of the basics.

func TestCorrect(t *testing.T) {
	type tt struct {
		opts *CorrectOptions
		err  string
	}

	tests := testy.NewTable()
	tests.Add("options schema", func(t *testing.T) interface{} {
		return tt{
			opts: &CorrectOptions{
				ParseOptions: &ParseOptions{
					Input: testFileReader(t, "testdata/success.json"),
				},
				OptionsSchema: true,
			},
		}
	})
	tests.Add("options schema empty", func(t *testing.T) interface{} {
		return tt{
			opts: &CorrectOptions{
				ParseOptions: &ParseOptions{
					Input: testFileReader(t, "testdata/message.env.yaml"),
				},
				OptionsSchema: true,
			},
		}
	})

	tests.Add("success", func(t *testing.T) interface{} {
		return tt{
			opts: &CorrectOptions{
				ParseOptions: &ParseOptions{
					Input: testFileReader(t, "testdata/success.json"),
				},
				Date: cal.MakeDate(2023, 4, 17),
				Data: []byte(`{"type":"credit-note","ext":{"es-facturae-correction":"01"}}`),
			},
		}
	})
	tests.Add("error missing data", func(t *testing.T) interface{} {
		return tt{
			opts: &CorrectOptions{
				ParseOptions: &ParseOptions{
					Input: testFileReader(t, "testdata/success.json"),
				},
				Data: []byte(`{"issue_date":"2023-04-17"}`),
			},
			err: "missing correction type",
		}
	})
	tests.Add("error invalid extension", func(t *testing.T) interface{} {
		return tt{
			opts: &CorrectOptions{
				ParseOptions: &ParseOptions{
					Input: testFileReader(t, "testdata/success.json"),
				},
				Data: []byte(`{"type":"credit-note","issue_date":"2023-04-17","ext":{"es-facturae-correct":"01"}}`),
			},
			err: "code=422, message=doc: (preceding: (0: (ext: (es-facturae-correct: undefined.).).).)",
		}
	})

	tests.Add("success just invoice", func(t *testing.T) interface{} {
		return tt{
			opts: &CorrectOptions{
				ParseOptions: &ParseOptions{
					Input: testFileReader(t, "testdata/invoice.json"),
				},
				Date: cal.MakeDate(2023, 4, 17),
				Data: []byte(`{"type":"credit-note","ext":{"es-facturae-correction":"01"}}`),
			},
		}
	})
	tests.Add("error just invoice", func(t *testing.T) interface{} {
		return tt{
			opts: &CorrectOptions{
				ParseOptions: &ParseOptions{
					Input: testFileReader(t, "testdata/invoice.json"),
				},
				Date: cal.MakeDate(2023, 4, 17),
			},
			err: "missing correction type",
		}
	})

	tests.Run(t, func(t *testing.T, tt tt) {
		t.Parallel()
		opts := tt.opts
		got, err := Correct(context.Background(), opts)
		if tt.err == "" {
			assert.Nil(t, err)
		} else {
			if assert.Error(t, err) {
				assert.Contains(t, err.Error(), tt.err)
			}
		}
		if err != nil {
			return
		}
		replacements := []testy.Replacement{
			{
				Regexp:      regexp.MustCompile(`(?s)"sigs": \[.*\]`),
				Replacement: `"sigs": ["signature data"]`,
			},
			{
				Regexp:      regexp.MustCompile(`"uuid":.?"[^\"]+"`),
				Replacement: `"uuid":"00000000-0000-0000-0000-000000000000"`,
			},
			{
				Regexp:      regexp.MustCompile(`"val":.?"[\w\d]{64}"`),
				Replacement: `"val":"74ffc799663823235951b43a1324c70555c0ba7e3b545c1f50af34bbcc57033b"`,
			},
		}
		if d := testy.DiffAsJSON(testy.Snapshot(t), got, replacements...); d != nil {
			t.Error(d)
		}
	})
}
