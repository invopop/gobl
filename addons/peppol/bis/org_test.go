package bis

import (
	"testing"

	"github.com/invopop/gobl/catalogues/iso"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

// Format/checksum validators are unit-tested directly because they're pure
// functions over the identifier string.

func TestValidGLN(t *testing.T) {
	cases := map[string]bool{
		"5790000000005":  true,  // Danish GLN test value
		"4006381333931":  true,  // GS1 sample
		"5790000000006":  false, // wrong checksum
		"541058300756":   false, // 12 digits
		"57900000000051": false, // 14 digits
		"5410583ABCDEF":  false, // non-numeric
		"":               false,
	}
	for code, want := range cases {
		assert.Equal(t, want, validGLN(code), code)
	}
}

func TestValidNorwegianOrg(t *testing.T) {
	cases := map[string]bool{
		"990983666":  true,  // valid Mod 11
		"974760673":  true,  // valid Mod 11
		"912345679":  false, // wrong checksum
		"91234567":   false, // 8 digits
		"9123456789": false, // 10 digits
		"abcdefghi":  false,
		"":           false,
	}
	for code, want := range cases {
		assert.Equal(t, want, validNorwegianOrg(code), code)
	}
}

func TestValidDanishCVR(t *testing.T) {
	cases := map[string]bool{
		"13585628":  true,
		"88146328":  true,
		"13585629":  false, // wrong checksum
		"1358562":   false, // 7 digits
		"135856280": false, // 9 digits
		"abcdefgh":  false,
		"":          false,
	}
	for code, want := range cases {
		assert.Equal(t, want, validDanishCVR(code), code)
	}
}

func TestValidBelgianEnterprise(t *testing.T) {
	cases := map[string]bool{
		"0403170701":  true,  // valid (78 % 97 == 1, last 2 digits = 01)
		"0123456749":  true,  // valid Mod 97
		"0403170702":  false, // wrong checksum
		"40317070":    false, // 8 digits
		"04031707011": false, // 11 digits
		"abcdefghij":  false,
	}
	for code, want := range cases {
		assert.Equal(t, want, validBelgianEnterprise(code), code)
	}
}

func TestValidITIPA(t *testing.T) {
	cases := map[string]bool{
		"ABC123":  true,
		"XYZ987":  true,
		"123456":  true,
		"abc123":  false, // lowercase not allowed
		"ABC12":   false, // too short
		"ABC1234": false, // too long
		"ABC-12":  false, // hyphen
	}
	for code, want := range cases {
		assert.Equal(t, want, validITIPA(code), code)
	}
}

func TestValidITCodiceFiscale(t *testing.T) {
	cases := map[string]bool{
		"12345678901":       true,  // legal entity (11 digits)
		"RSSMRA80A01H501U":  true,  // person (16 alphanumerics, structured)
		"1234567890":        false, // 10 digits
		"123456789012":      false, // 12 digits
		"RSSMRA80A01H501":   false, // 15 chars
		"RSSMRA80A01H501UU": false, // 17 chars
	}
	for code, want := range cases {
		assert.Equal(t, want, validITCodiceFiscale(code), code)
	}
}

func TestValidITPartitaIVA(t *testing.T) {
	cases := map[string]bool{
		"12345678901":  true,
		"00000000000":  true,
		"1234567890":   false,
		"123456789012": false,
		"abcdefghijk":  false,
	}
	for code, want := range cases {
		assert.Equal(t, want, validITPartitaIVA(code), code)
	}
}

func TestValidSwedishOrg(t *testing.T) {
	cases := map[string]bool{
		"5560360793":  true,  // Volvo AB
		"5560743089":  true,  // Volvo Cars
		"5560747551":  true,  // IKEA
		"5560360794":  false, // wrong checksum
		"556036079":   false, // 9 digits
		"55603607930": false, // 11 digits
		"abcdefghij":  false,
	}
	for code, want := range cases {
		assert.Equal(t, want, validSwedishOrg(code), code)
	}
}

func TestLuhnValid(t *testing.T) {
	assert.True(t, luhnValid("4532015112830366"))  // valid VISA-style
	assert.True(t, luhnValid("5560360793"))        // valid Swedish org
	assert.False(t, luhnValid("4532015112830367")) // wrong checksum
	assert.True(t, luhnValid("0"))                 // sum 0 → divisible by 10
}

func TestValidAustralianABN(t *testing.T) {
	cases := map[string]bool{
		"51824753556":  true,  // valid ABN
		"83914571673":  true,  // valid ABN
		"51824753557":  false, // wrong checksum
		"5182475355":   false, // 10 digits
		"518247535561": false, // 12 digits
		"abcdefghijk":  false,
	}
	for code, want := range cases {
		assert.Equal(t, want, validAustralianABN(code), code)
	}
}

func TestValidDanishPNumber(t *testing.T) {
	assert.True(t, validDanishPNumber("1234567890"))
	assert.False(t, validDanishPNumber("123456789"))   // 9 digits
	assert.False(t, validDanishPNumber("12345678901")) // 11 digits
	assert.False(t, validDanishPNumber("12345abcde"))
}

func TestValidDanishSENumber(t *testing.T) {
	assert.True(t, validDanishSENumber("12345678"))
	assert.False(t, validDanishSENumber("1234567"))   // 7 digits
	assert.False(t, validDanishSENumber("123456789")) // 9 digits
	assert.False(t, validDanishSENumber("1234abcd"))
}

func TestOnlyDigits(t *testing.T) {
	assert.True(t, onlyDigits("12345"))
	assert.False(t, onlyDigits("12345A"))
	assert.False(t, onlyDigits(""))
}

// checkSchemeFormat dispatches to the right validator based on scheme id.

func TestCheckSchemeFormat(t *testing.T) {
	cases := []struct {
		scheme cbc.Code
		code   cbc.Code
		ok     bool
	}{
		// known schemes — valid codes
		{schemeGLN, "5790000000005", true},
		{schemeNOOrg, "990983666", true},
		{schemeDKCVR, "13585628", true},
		{schemeBEEnt, "0403170701", true},
		{schemeITIPA, "ABC123", true},
		{schemeITCF, "12345678901", true},
		{schemeITPIva, "12345678901", true},
		{schemeSEOrg, "5560360793", true},
		{schemeAUABN, "51824753556", true},
		{schemeDKPNum, "1234567890", true},
		{schemeDKSENum, "12345678", true},
		// known schemes — invalid codes (must error)
		{schemeGLN, "5790000000006", false},
		{schemeNOOrg, "912345679", false},
		{schemeDKCVR, "13585629", false},
		{schemeBEEnt, "0403170702", false},
		{schemeITIPA, "abc", false},
		// unknown scheme — always passes
		{cbc.Code("9999"), "anything", true},
		// empty code — passes
		{schemeGLN, "", true},
	}
	for _, c := range cases {
		err := checkSchemeFormat(c.scheme, c.code)
		if c.ok {
			assert.NoError(t, err, "%s/%s", c.scheme, c.code)
		} else {
			assert.Error(t, err, "%s/%s", c.scheme, c.code)
		}
	}
}

func TestIdentityFormatValid(t *testing.T) {
	t.Run("nil identity", func(t *testing.T) {
		assert.NoError(t, identityFormatValid(nil))
	})
	t.Run("non-identity value", func(t *testing.T) {
		assert.NoError(t, identityFormatValid("not an identity"))
	})
	t.Run("identity with no scheme — skipped", func(t *testing.T) {
		id := &org.Identity{Code: "anything"}
		assert.NoError(t, identityFormatValid(id))
	})
	t.Run("identity with valid GLN", func(t *testing.T) {
		id := &org.Identity{Code: "5790000000005", Ext: tax.Extensions{iso.ExtKeySchemeID: "0088"}}
		assert.NoError(t, identityFormatValid(id))
	})
	t.Run("identity with invalid GLN", func(t *testing.T) {
		id := &org.Identity{Code: "5790000000006", Ext: tax.Extensions{iso.ExtKeySchemeID: "0088"}}
		assert.Error(t, identityFormatValid(id))
	})
}

func TestInboxFormatValid(t *testing.T) {
	t.Run("nil inbox", func(t *testing.T) {
		assert.NoError(t, inboxFormatValid(nil))
	})
	t.Run("non-inbox value", func(t *testing.T) {
		assert.NoError(t, inboxFormatValid(42))
	})
	t.Run("inbox with no scheme", func(t *testing.T) {
		assert.NoError(t, inboxFormatValid(&org.Inbox{Code: "any"}))
	})
	t.Run("inbox with no code", func(t *testing.T) {
		assert.NoError(t, inboxFormatValid(&org.Inbox{Scheme: "0088"}))
	})
	t.Run("inbox with valid scheme + GLN", func(t *testing.T) {
		assert.NoError(t, inboxFormatValid(&org.Inbox{Scheme: "0088", Code: "5790000000005"}))
	})
	t.Run("inbox with valid scheme + bad GLN", func(t *testing.T) {
		assert.Error(t, inboxFormatValid(&org.Inbox{Scheme: "0088", Code: "5790000000006"}))
	})
}
