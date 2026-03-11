package rules_test

import (
	"testing"

	"github.com/invopop/gobl/rules"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// patternCode is a named string type used to exercise Matches via For.
type patternCode string

func matchSet(pattern string) *rules.Set {
	return rules.For(patternCode(""),
		rules.Assert("001", "must match pattern",
			rules.Matches(pattern),
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
		proto := new(Thing)
		set := rules.For(proto,
			rules.Field(&proto.Code,
				rules.Assert("001", "must match",
					rules.Matches(`^\d+$`),
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
		proto := new(Blob)
		set := rules.For(proto,
			rules.Field(&proto.Data,
				rules.Assert("001", "must match",
					rules.Matches(`^\d+$`),
				),
			),
		)
		assert.Nil(t, set.Validate(&Blob{Data: []byte("123")}))
		assert.NotNil(t, set.Validate(&Blob{Data: []byte("abc")}))
		assert.Nil(t, set.Validate(&Blob{Data: []byte{}}))
	})

	t.Run("String returns pattern description", func(t *testing.T) {
		assert.Equal(t, `matches ^\d+$`, rules.Matches(`^\d+$`).String())
	})

	t.Run("panics on invalid regex pattern", func(t *testing.T) {
		assert.Panics(t, func() {
			rules.For(patternCode(""),
				rules.Assert("001", "bad pattern",
					rules.Matches(`[invalid`),
				),
			)
		})
	})
}
