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
	assert.Equal(t, l10n.Unions().Len(), 2)
	assert.Nil(t, l10n.Unions().Code(l10n.Code("X")))
}

func TestSEPAUnion(t *testing.T) {
	u := l10n.Union(l10n.SEPA)
	assert.NotNil(t, u)
	assert.Equal(t, "Single Euro Payments Area", u.Name)

	// EU member states are in SEPA.
	assert.True(t, u.HasMember(l10n.ES))
	assert.True(t, u.HasMember(l10n.DE))
	// Greece resolves via both its ISO and tax country code.
	assert.True(t, u.HasMember(l10n.GR))
	assert.True(t, u.HasMember(l10n.EL))
	// Non-EU SEPA participants are members too — the case the EU union misses.
	assert.True(t, u.HasMember(l10n.NO))
	assert.True(t, u.HasMember(l10n.CH))
	assert.True(t, u.HasMember(l10n.GB))
	assert.True(t, u.HasMember(l10n.SM))
	// Recent additions (Albania, Montenegro, North Macedonia, Moldova, Serbia).
	assert.True(t, u.HasMember(l10n.AL))
	assert.True(t, u.HasMember(l10n.ME))
	assert.True(t, u.HasMember(l10n.MK))
	assert.True(t, u.HasMember(l10n.MD))
	assert.True(t, u.HasMember(l10n.RS))
	// Non-SEPA countries are not members.
	assert.False(t, u.HasMember(l10n.US))
	assert.False(t, u.HasMember(l10n.CO))
	assert.False(t, u.HasMember(l10n.Code("")))

	// Entry dates are respected: members are not in scope before their date.
	assert.False(t, u.HasMemberOn(cal.MakeDate(2016, 4, 30), l10n.JE))
	assert.True(t, u.HasMemberOn(cal.MakeDate(2016, 5, 1), l10n.JE))
	assert.False(t, u.HasMemberOn(cal.MakeDate(2019, 2, 28), l10n.AD))
	assert.True(t, u.HasMemberOn(cal.MakeDate(2019, 3, 1), l10n.AD))
	assert.False(t, u.HasMemberOn(cal.MakeDate(2025, 10, 4), l10n.AL))
	assert.True(t, u.HasMemberOn(cal.MakeDate(2025, 10, 5), l10n.AL))
	// Croatia was not in SEPA at the 2008 launch (joined the EU in 2013).
	assert.False(t, u.HasMemberOn(cal.MakeDate(2010, 1, 1), l10n.HR))
	assert.True(t, u.HasMemberOn(cal.MakeDate(2013, 7, 1), l10n.HR))
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
