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

func TestUnion(t *testing.T) {
	u := l10n.Union(l10n.EU)
	assert.NotNil(t, u)
	assert.Equal(t, l10n.EU, u.Code)
	assert.Equal(t, "European Union", u.Name)
	assert.Equal(t, 28, len(u.Members))
}

func TestUnionMember(t *testing.T) {
	u := l10n.Unions().Code(l10n.EU)
	assert.True(t, u.HasMember(l10n.DE))
	assert.False(t, u.HasMember(l10n.US))

	// and the greek case with alt code
	assert.True(t, u.HasMember(l10n.GR))
	assert.True(t, u.HasMember(l10n.EL))

	d1a := cal.MakeDate(1973, 1, 1) // happy times
	assert.True(t, u.HasMemberOn(d1a, l10n.GB))
	d1b := cal.MakeDate(2020, 1, 31) // end of happy times
	assert.True(t, u.HasMemberOn(d1b, l10n.GB))
	d2 := cal.MakeDate(2020, 2, 1) // sad times
	assert.False(t, u.HasMemberOn(d2, l10n.GB))

	// empty code should not match members with empty AltCode
	assert.False(t, u.HasMember(l10n.Code("")))
	assert.False(t, u.HasMemberOn(d1a, l10n.Code("")))
}
