package is_test

import (
	"regexp"
	"testing"

	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// patternCode is a named string type used to exercise Matches via For.
type patternCode string

func matchSet(pattern string) *rules.Set {
	return rules.For(patternCode(""),
		rules.Assert("001", "must match pattern",
			is.Matches(pattern),
		),
	)
}

func TestMatches(t *testing.T) {
	t.Run("passes when string matches pattern", func(t *testing.T) {
		set := matchSet(`^\d+$`)
		assert.Nil(t, set.Validate(patternCode("123")))
	})

	t.Run("fails when string does not match pattern", func(t *testing.T) {
		set := matchSet(`^\d+$`)
		faults := set.Validate(patternCode("abc"))
		require.NotNil(t, faults)
		assert.Equal(t, 1, faults.Len())
	})

	t.Run("passes when string is empty", func(t *testing.T) {
		set := matchSet(`^\d+$`)
		assert.Nil(t, set.Validate(patternCode("")))
	})

	t.Run("passes when pointer is nil", func(t *testing.T) {
		type Thing struct {
			Code *string `json:"code"`
		}
		set := rules.For(new(Thing),
			rules.Field("code",
				rules.Assert("001", "must match",
					is.Matches(`^\d+$`),
				),
			),
		)
		faults := set.Validate(&Thing{Code: nil})
		assert.Nil(t, faults)
	})

	t.Run("works with byte slices", func(t *testing.T) {
		type Blob struct {
			Data []byte `json:"data"`
		}
		set := rules.For(new(Blob),
			rules.Field("data",
				rules.Assert("001", "must match",
					is.Matches(`^\d+$`),
				),
			),
		)
		assert.Nil(t, set.Validate(&Blob{Data: []byte("123")}))
		assert.NotNil(t, set.Validate(&Blob{Data: []byte("abc")}))
		assert.Nil(t, set.Validate(&Blob{Data: []byte{}}))
	})

	t.Run("String returns pattern description", func(t *testing.T) {
		assert.Equal(t, `matches ^\d+$`, is.Matches(`^\d+$`).String())
	})

	t.Run("panics on invalid regex pattern", func(t *testing.T) {
		assert.Panics(t, func() {
			rules.For(patternCode(""),
				rules.Assert("001", "bad pattern",
					is.Matches(`[invalid`),
				),
			)
		})
	})
}

func TestMatchesRegexp(t *testing.T) {
	re := regexp.MustCompile(`^\d+$`)

	t.Run("passes when string matches", func(t *testing.T) {
		assert.True(t, is.MatchesRegexp(re).Check(patternCode("123")))
	})

	t.Run("fails when string does not match", func(t *testing.T) {
		assert.False(t, is.MatchesRegexp(re).Check(patternCode("abc")))
	})

	t.Run("passes when string is empty", func(t *testing.T) {
		assert.True(t, is.MatchesRegexp(re).Check(patternCode("")))
	})

	t.Run("String returns pattern description", func(t *testing.T) {
		assert.Equal(t, `matches ^\d+$`, is.MatchesRegexp(re).String())
	})
}
