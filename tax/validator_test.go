package tax_test

import (
	"context"
	"testing"

	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
	"github.com/stretchr/testify/assert"
)

func TestValidateStructWithContext(t *testing.T) {
	t.Run("with errors", func(t *testing.T) {
		type testStruct struct {
			Name string `json:"test"`
		}
		ts := &testStruct{
			Name: "",
		}
		ctx := tax.ContextWithValidator(context.Background(), func(doc any) error {
			if ts, ok := doc.(*testStruct); ok {
				if ts.Name == "" {
					return validation.NewError("test", "name is required")
				}
			}
			return nil
		})

		err := tax.ValidateStructWithContext(ctx, ts)
		assert.ErrorContains(t, err, "name is required")

	})
}
