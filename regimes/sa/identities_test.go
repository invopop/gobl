package sa_test

import (
	"testing"

	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/sa"
	"github.com/stretchr/testify/assert"
)

func TestIdentityTypeDefinitions(t *testing.T) {
	assert.Equal(t, "CRN", string(sa.IdentityTypeCRN))
	assert.Equal(t, "MOM", string(sa.IdentityTypeMOM))
	assert.Equal(t, "MLS", string(sa.IdentityTypeMLS))
	assert.Equal(t, "700", string(sa.IdentityType700))
	assert.Equal(t, "SAG", string(sa.IdentityTypeSAG))
	assert.Equal(t, "NAT", string(sa.IdentityTypeNAT))
	assert.Equal(t, "IQA", string(sa.IdentityTypeIQA))
	assert.Equal(t, "PAS", string(sa.IdentityTypePAS))
	assert.Equal(t, "GCC", string(sa.IdentityTypeGCC))
	assert.Equal(t, "OTH", string(sa.IdentityTypeOTH))
}

func TestValidateIdentity(t *testing.T) {
	tests := []struct {
		name     string
		identity *org.Identity
		err      string
	}{
		{
			name: "valid CRN",
			identity: &org.Identity{
				Type: "CRN",
				Code: "1234567890",
			},
		},
		{
			name: "CRN too short",
			identity: &org.Identity{
				Type: "CRN",
				Code: "123456789",
			},
			err: "must be a 10-digit number",
		},
		{
			name: "CRN too long",
			identity: &org.Identity{
				Type: "CRN",
				Code: "12345678901",
			},
			err: "must be a 10-digit number",
		},
		{
			name: "CRN non-numeric",
			identity: &org.Identity{
				Type: "CRN",
				Code: "12345678AB",
			},
			err: "must be a 10-digit number",
		},
		{
			name: "CRN empty code",
			identity: &org.Identity{
				Type: "CRN",
				Code: "",
			},
		},
		// NAT validation
		{
			name: "valid NAT",
			identity: &org.Identity{
				Type: "NAT",
				Code: "1234567890",
			},
		},
		{
			name: "NAT wrong start digit",
			identity: &org.Identity{
				Type: "NAT",
				Code: "2234567890",
			},
			err: "must be a 10-digit number starting with 1",
		},
		{
			name: "NAT too short",
			identity: &org.Identity{
				Type: "NAT",
				Code: "123456789",
			},
			err: "must be a 10-digit number starting with 1",
		},
		{
			name: "NAT too long",
			identity: &org.Identity{
				Type: "NAT",
				Code: "12345678901",
			},
			err: "must be a 10-digit number starting with 1",
		},
		{
			name: "NAT empty code",
			identity: &org.Identity{
				Type: "NAT",
				Code: "",
			},
		},
		// 700 validation
		{
			name: "valid 700",
			identity: &org.Identity{
				Type: "700",
				Code: "7000012345",
			},
		},
		{
			name: "700 wrong start digit",
			identity: &org.Identity{
				Type: "700",
				Code: "1234567890",
			},
			err: "must be a 10-digit number starting with 7",
		},
		{
			name: "700 too short",
			identity: &org.Identity{
				Type: "700",
				Code: "700001234",
			},
			err: "must be a 10-digit number starting with 7",
		},
		{
			name: "700 too long",
			identity: &org.Identity{
				Type: "700",
				Code: "70000123456",
			},
			err: "must be a 10-digit number starting with 7",
		},
		{
			name: "700 empty code",
			identity: &org.Identity{
				Type: "700",
				Code: "",
			},
		},
		// IQA validation
		{
			name: "valid IQA",
			identity: &org.Identity{
				Type: "IQA",
				Code: "2345678901",
			},
		},
		{
			name: "IQA wrong start digit",
			identity: &org.Identity{
				Type: "IQA",
				Code: "1345678901",
			},
			err: "must be a 10-digit number starting with 2",
		},
		{
			name: "IQA too short",
			identity: &org.Identity{
				Type: "IQA",
				Code: "234567890",
			},
			err: "must be a 10-digit number starting with 2",
		},
		{
			name: "IQA too long",
			identity: &org.Identity{
				Type: "IQA",
				Code: "23456789012",
			},
			err: "must be a 10-digit number starting with 2",
		},
		// PAS, GCC, OTH skip validation
		{
			name: "PAS skips validation",
			identity: &org.Identity{
				Type: "PAS",
				Code: "AB1234567",
			},
		},
		{
			name: "GCC skips validation",
			identity: &org.Identity{
				Type: "GCC",
				Code: "GCC-12345",
			},
		},
		{
			name: "OTH skips validation",
			identity: &org.Identity{
				Type: "OTH",
				Code: "anything",
			},
		},
		{
			name: "non-CRN identity skips validation",
			identity: &org.Identity{
				Type: "MOM",
				Code: "anything",
			},
		},
		{
			name: "unknown identity type skips validation",
			identity: &org.Identity{
				Type: "other",
				Code: "invalid",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := sa.Validate(tt.identity)
			if tt.err == "" {
				assert.NoError(t, err)
			} else {
				if assert.Error(t, err) {
					assert.Contains(t, err.Error(), tt.err)
				}
			}
		})
	}
}

func TestNormalizeIdentity(t *testing.T) {
	r := sa.New()

	t.Run("normalize CRN with bad chars", func(t *testing.T) {
		id := &org.Identity{
			Type: sa.IdentityTypeCRN,
			Code: "12-345.678 90",
		}
		r.NormalizeObject(id)
		assert.Equal(t, "1234567890", id.Code.String())
	})

	t.Run("CRN already normalized", func(t *testing.T) {
		id := &org.Identity{
			Type: sa.IdentityTypeCRN,
			Code: "1234567890",
		}
		r.NormalizeObject(id)
		assert.Equal(t, "1234567890", id.Code.String())
	})

	t.Run("normalize NAT with bad chars", func(t *testing.T) {
		id := &org.Identity{
			Type: sa.IdentityTypeNAT,
			Code: "1-234-567 890",
		}
		r.NormalizeObject(id)
		assert.Equal(t, "1234567890", id.Code.String())
	})

	t.Run("normalize 700 with bad chars", func(t *testing.T) {
		id := &org.Identity{
			Type: sa.IdentityType700,
			Code: "7-000.012 345",
		}
		r.NormalizeObject(id)
		assert.Equal(t, "7000012345", id.Code.String())
	})

	t.Run("normalize IQA with bad chars", func(t *testing.T) {
		id := &org.Identity{
			Type: sa.IdentityTypeIQA,
			Code: "2-345.678 901",
		}
		r.NormalizeObject(id)
		assert.Equal(t, "2345678901", id.Code.String())
	})

	t.Run("non-CRN not normalized", func(t *testing.T) {
		id := &org.Identity{
			Type: sa.IdentityTypeMOM,
			Code: "12-345",
		}
		r.NormalizeObject(id)
		assert.Equal(t, "12-345", id.Code.String())
	})

	t.Run("PAS not normalized", func(t *testing.T) {
		id := &org.Identity{
			Type: sa.IdentityTypePAS,
			Code: "AB-1234",
		}
		r.NormalizeObject(id)
		assert.Equal(t, "AB-1234", id.Code.String())
	})
}
