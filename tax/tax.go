package tax

import (
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/schema"
)

func init() {
	schema.Register(schema.GOBL.Add("tax"),
		Set{},
		Total{},
		Region{},
	)
	org.SetTaxIdentityValidation(ValidateTaxIdentity)
	org.SetTaxIdentityNormalizer(NormalizeTaxIdentity)
}
