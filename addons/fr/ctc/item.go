package ctc

import (
	"errors"
	"strings"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/validation"
)

// validateItem validates an item
func validateItem(item *org.Item) error {
	if item == nil {
		return nil
	}

	return validation.ValidateStruct(item,
		validation.Field(&item.Meta,
			validation.By(validateItemMeta),
			validation.Skip,
		),
	)
}

// validateItemMeta ensures meta values are not blank (e.g. empty or whitespace only) BR-FR-28.
//
// Note: We currently do not map ValueQuantity in UBL which may be parsed.
// One way to support this is in parsing to map this to the value and concat
// the unit code. Need to implement this in the converter.
func validateItemMeta(value any) error {
	meta, ok := value.(cbc.Meta)
	if !ok || len(meta) == 0 {
		return nil
	}

	for key, val := range meta {
		// Check if value is blank (empty or whitespace only)
		if strings.TrimSpace(val) == "" {
			return validation.Errors{
				key.String(): errors.New("value cannot be blank (BR-FR-28)"),
			}
		}
	}

	return nil
}
