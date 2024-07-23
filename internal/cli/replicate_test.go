package cli

import (
	"context"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/flimzy/testy"
)

// These tests depend on the build_test.go for some of the basics.

func TestReplicate(t *testing.T) {
	type tt struct {
		opts *ReplicateOptions
		err  string
	}

	tests := testy.NewTable()

	tests.Add("success", func(t *testing.T) interface{} {
		return tt{
			opts: &ReplicateOptions{
				ParseOptions: &ParseOptions{
					Input: testFileReader(t, "testdata/success.json"),
				},
			},
		}
	})

	tests.Add("success just invoice", func(t *testing.T) interface{} {
		return tt{
			opts: &ReplicateOptions{
				ParseOptions: &ParseOptions{
					Input: testFileReader(t, "testdata/invoice.json"),
				},
			},
		}
	})

	tests.Run(t, func(t *testing.T, tt tt) {
		t.Parallel()
		opts := tt.opts
		got, err := Replicate(context.Background(), opts)
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
			{
				Regexp:      regexp.MustCompile(`"issue_date":.?"[^\"]+"`),
				Replacement: `"issue_date":"2024-05-06"`,
			},
		}
		if d := testy.DiffAsJSON(testy.Snapshot(t), got, replacements...); d != nil {
			t.Error(d)
		}
	})
}
