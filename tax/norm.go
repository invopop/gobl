package tax

import (
	"strings"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/norm"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
)

func init() {
	norm.Register("tax",
		// Expand any addon dependencies (Requires) into the addon list. This
		// runs on every pass; once a required addon key appears, the norm
		// engine re-collects the context and applies that addon's normalizers
		// on the following pass.
		norm.For(func(as *Addons) { as.normalizeAddons() }),
		norm.For(normalizeCombo),
		norm.For(normalizeNote),
		// Tax identities are an exception to the normal normalization rules:
		// they are normalized by their own country's regime (not the document
		// regime or any addon). The country-specific normalizers register
		// themselves with a tax.IdentityIn guard; here we only provide the
		// fallback used when the identity's country has no regime at all.
		norm.When(identityWithoutRegime,
			norm.For(func(id *Identity) { NormalizeIdentity(id) }),
		),
	)
}

// identityWithoutRegime passes for a tax identity whose country has no
// registered tax regime, mirroring the previous Identity.Normalize fallback.
var identityWithoutRegime rules.Test = is.Func("identity country has no regime",
	func(value any) bool {
		id, ok := value.(*Identity)
		return ok && id != nil && regimes.For(id.Country.Code()) == nil
	},
)

// normalizeCombo applies the intrinsic normalization of a tax combo: it maps
// legacy VAT rate keys onto the current key/rate model and cleans extensions.
// Regime and addon adjustments are applied separately by their registered
// normalizers.
func normalizeCombo(c *Combo) {
	if c == nil {
		return
	}

	switch c.Category {
	case CategoryVAT:
		switch c.Rate {
		case KeyZero:
			c.Key = KeyZero
			c.Rate = cbc.KeyEmpty
			if c.Percent == nil {
				c.Percent = num.NewPercentage(0, 2)
			}
		case KeyExempt:
			// This can cause problems with backwards compatibility as the "exempt"
			// rate was used too widely. Addons will need to try and account for this.
			c.Key = KeyExempt
			c.Rate = cbc.KeyEmpty
		case KeyExempt.With("reverse-charge"):
			c.Key = KeyReverseCharge
			c.Rate = cbc.KeyEmpty
			c.Percent = nil
		case KeyExempt.With("export"):
			c.Key = KeyExport
			c.Rate = cbc.KeyEmpty
		case KeyExempt.With("eea"), KeyExempt.With("export").With("eea"):
			c.Key = KeyIntraCommunity
			c.Rate = cbc.KeyEmpty
		default:
			// Make no further assumptions about the key, but try to replace standard
			// rate with general.
			if c.Rate == KeyStandard {
				c.Rate = RateGeneral
			} else if found, ok := strings.CutPrefix(c.Rate.String(), "standard+"); ok {
				c.Rate = cbc.Key(RateGeneral.String() + "+" + found)
			}
		}

		switch c.Key {
		case cbc.KeyEmpty:
			// Special case for zero percent which has no additional rates
			if c.Percent != nil && c.Percent.IsZero() {
				c.Key = KeyZero
			}
		case KeyZero:
			if c.Percent == nil {
				zp := num.PercentageZero
				c.Percent = &zp
			}
		}
	}

	c.Ext = c.Ext.Clean()
}

// normalizeNote applies the intrinsic normalization of a tax note.
func normalizeNote(n *Note) {
	if n == nil {
		return
	}
	n.Ext = n.Ext.Clean()
}
