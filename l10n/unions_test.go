package l10n_test

import (
	"testing"

	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/l10n"
	"github.com/stretchr/testify/assert"
)

func TestUnions(t *testing.T) {
	u := l10n.Unions().Code(l10n.EU)
	assert.Equal(t, l10n.EU, u.Code)
	assert.Equal(t, l10n.Unions().Len(), 1)
	assert.Nil(t, l10n.Unions().Code(l10n.Code("X")))
}

func TestUnionMember(t *testing.T) {
	u := l10n.Unions().Code(l10n.EU)
	assert.True(t, u.HasMember(l10n.DE))
	assert.False(t, u.HasMember(l10n.US))

	d1a := cal.MakeDate(1973, 1, 1) // happy times
	assert.True(t, u.HasMemberOn(d1a, l10n.GB))
	d1b := cal.MakeDate(2020, 1, 31) // end of happy times
	assert.True(t, u.HasMemberOn(d1b, l10n.GB))
	d2 := cal.MakeDate(2020, 2, 1) // sad times
	assert.False(t, u.HasMemberOn(d2, l10n.GB))
}
