package is_test

import (
	"testing"

	"github.com/invopop/gobl/rules/is"
	"github.com/stretchr/testify/assert"
)

func TestDigit(t *testing.T) {
	assert.True(t, is.Digit.Check("0123456789"))
	assert.True(t, is.Digit.Check("42"))
	assert.False(t, is.Digit.Check("12.3"))
	assert.False(t, is.Digit.Check("12a"))
	assert.False(t, is.Digit.Check(""))
}

func TestE164(t *testing.T) {
	assert.True(t, is.E164.Check("+14155552671"))
	assert.True(t, is.E164.Check("+442071838750"))
	assert.True(t, is.E164.Check("14155552671")) // leading + optional
	assert.False(t, is.E164.Check("not-a-number"))
	assert.False(t, is.E164.Check("+0123456789")) // leading digit after + must be 1-9
	assert.False(t, is.E164.Check("+1"))          // too short
}

func TestSubdomain(t *testing.T) {
	assert.True(t, is.Subdomain.Check("example"))
	assert.True(t, is.Subdomain.Check("my-subdomain"))
	assert.True(t, is.Subdomain.Check("sub123"))
	assert.False(t, is.Subdomain.Check("-invalid"))
	assert.False(t, is.Subdomain.Check("invalid-"))
	assert.False(t, is.Subdomain.Check("has.dot"))
	assert.False(t, is.Subdomain.Check(""))
}

func TestDomain(t *testing.T) {
	assert.True(t, is.Domain.Check("example.com"))
	assert.True(t, is.Domain.Check("sub.example.com"))
	assert.True(t, is.Domain.Check("my-site.co.uk"))
	assert.False(t, is.Domain.Check("example"))           // no TLD
	assert.False(t, is.Domain.Check("-bad.com"))          // leading dash
	assert.False(t, is.Domain.Check(""))
	assert.False(t, is.Domain.Check(string(make([]byte, 256)))) // > 255 chars
}

func TestUTFNumeric(t *testing.T) {
	assert.True(t, is.UTFNumeric.Check("123"))
	assert.True(t, is.UTFNumeric.Check("²³"))  // superscript digits (unicode category N)
	assert.False(t, is.UTFNumeric.Check("12a"))
	assert.True(t, is.UTFNumeric.Check("")) // empty passes (no non-numeric chars)
}

func TestISBN(t *testing.T) {
	assert.True(t, is.ISBN.Check("0-306-40615-2"))   // valid ISBN-10
	assert.True(t, is.ISBN.Check("978-3-16-148410-0")) // valid ISBN-13
	assert.False(t, is.ISBN.Check("not-an-isbn"))
	assert.False(t, is.ISBN.Check(""))
}
