package en16931

import (
	"testing"

	"github.com/invopop/gobl/org"
	"github.com/stretchr/testify/assert"
)

func TestOrgIdentitiesScopeHelpers(t *testing.T) {
	legal := &org.Identity{Code: "1", Scope: org.IdentityScopeLegal}
	taxID := &org.Identity{Code: "2", Scope: org.IdentityScopeTax}

	// non-slice values are permitted (defensive guard).
	assert.True(t, orgIdentitiesSingleLegalScope("not a slice"))
	assert.True(t, orgIdentitiesSingleTaxScope("not a slice"))

	assert.True(t, orgIdentitiesSingleLegalScope([]*org.Identity{legal, taxID}))
	assert.False(t, orgIdentitiesSingleLegalScope([]*org.Identity{legal, legal}))

	assert.True(t, orgIdentitiesSingleTaxScope([]*org.Identity{legal, taxID}))
	assert.False(t, orgIdentitiesSingleTaxScope([]*org.Identity{taxID, taxID}))

	// nil entries are skipped.
	assert.Equal(t, 1, orgIdentitiesScopeCount([]*org.Identity{nil, legal}, org.IdentityScopeLegal))
}
