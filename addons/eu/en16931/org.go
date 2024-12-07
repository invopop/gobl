package en16931

import (
	"github.com/invopop/gobl/org"
	"github.com/invopop/validation"
)

func validateOrgAttachment(a *org.Attachment) error {
	return validation.ValidateStruct(a,
		validation.Field(&a.Code,
			validation.Required,
			validation.Skip,
		),
	)
}
