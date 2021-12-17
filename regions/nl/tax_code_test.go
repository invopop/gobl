package nl

import (
	"testing"

	"gitlab.com/flimzy/testy"
)

func TestVerifyTaxCode(t *testing.T) {
	type tt struct {
		code, err string
	}

	tests := testy.NewTable()
	tests.Add("empty", tt{
		code: "",
		err:  "invalid VAT number",
	})
	tests.Add("too long", tt{
		code: "a really really long string that's way too long",
		err:  "invalid VAT number",
	})
	tests.Add("too short", tt{
		code: "shorty",
		err:  "invalid VAT number",
	})
	tests.Add("valid", tt{
		code: "NL000099995B57",
	})
	tests.Add("lowercase", tt{
		code: "nl000099995b57",
	})
	tests.Add("no B", tt{
		code: "NL000099998X57",
		err:  "invalid VAT number",
	})
	tests.Add("non numbers", tt{
		code: "NL000099998B5a",
		err:  "invalid VAT number",
	})
	tests.Add("invalid checksum", tt{
		code: "NL123456789B12",
		err:  "checkusum mismatch",
	})

	tests.Run(t, func(t *testing.T, tt tt) {
		err := VerifyTaxCode(tt.code)
		if !testy.ErrorMatches(tt.err, err) {
			t.Errorf("Unexpected error: %s", err)
		}
	})
}
