package bill

import (
	"context"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

// TaxScheme allows for defining a specific or special scheme that applies to the
// billing document. Schemes are defined as needed for each region.
type TaxScheme string

// Tax defines a summary of the taxes which may be applied to an invoice.
type Tax struct {
	// Addons defines a list of keys used to identify tax addons that apply special
	// normalization, scenarios, and validation rules to a document.
	Addons []cbc.Key `json:"addons,omitempty" jsonschame:"title=Addons"`

	// Category of the tax already included in the line item prices, especially
	// useful for B2C retailers with customers who prefer final prices inclusive of
	// tax.
	PricesInclude cbc.Code `json:"prices_include,omitempty" jsonschema:"title=Prices Include"`

	// Tags are used to help identify specific tax scenarios or requirements that will
	// apply changes to the contents of the invoice. Tags by design should always be optional,
	// it should always be possible to build a valid invoice without any tags.
	Tags []cbc.Key `json:"tags,omitempty" jsonschema:"title=Tags"`

	// Additional extensions that are applied to the invoice as a whole as opposed to specific
	// sections.
	Ext tax.Extensions `json:"ext,omitempty" jsonschema:"title=Extensions"`

	// Any additional data that may be required for processing, but should never
	// be relied upon by recipients.
	Meta cbc.Meta `json:"meta,omitempty" jsonschema:"title=Meta"`
}

// Normalize performs normalization on the tax and embedded objects using the
// provided list of normalizers.
func (t *Tax) Normalize(normalizers tax.Normalizers) {
	if t == nil {
		return
	}
	t.Ext = tax.CleanExtensions(t.Ext)
	normalizers.Each(t)
}

// ContainsTag returns true if the tax contains the given tag.
func (t *Tax) ContainsTag(key cbc.Key) bool {
	if t == nil {
		return false
	}
	return key.In(t.Tags...)
}

// HasTags returns true if the tax object contains all of the
// provided tags. Can be be used against a nil tax object in
// case it has not yet been initialized.
func (t *Tax) HasTags(keys ...cbc.Key) bool {
	if t == nil {
		return false
	}
	for _, k := range keys {
		if !k.In(t.Tags...) {
			return false
		}
	}
	return true
}

// GetAddons provides the list of addon instances ready to use.
func (t *Tax) GetAddons() []*tax.Addon {
	if t == nil {
		return nil
	}
	addons := make([]*tax.Addon, 0, len(t.Addons))
	for _, ak := range t.Addons {
		if a := tax.AddonForKey(ak); a != nil {
			addons = append(addons, a)
		}
	}
	return addons
}

// ValidateWithContext ensures the tax details look valid.
func (t *Tax) ValidateWithContext(ctx context.Context) error {
	r := tax.RegimeFromContext(ctx)
	return tax.ValidateStructWithContext(ctx, t,
		validation.Field(&t.PricesInclude),
		validation.Field(&t.Addons, validation.Each(tax.AddonRegistered)),
		validation.Field(&t.Tags, validation.Each(r.InTags())),
		validation.Field(&t.Ext),
		validation.Field(&t.Meta),
	)
}
