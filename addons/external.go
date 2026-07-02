package addons

import (
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/tax"
)

// This file is the curated list of approved external addons — addons whose
// implementation lives in a separate Go module but whose keys GOBL recognises
// (see [tax.ExternalAddon]). It sits alongside the in-core addon imports above
// so the full set of addons GOBL knows about — bundled and external — is
// visible in one place. Adding an entry here is the approval step and must be
// reviewed via pull request.
//
// Approval checklist for a new entry:
//
//   - The implementation lives in a public module under github.com/invopop and
//     auto-registers via init() + tax.RegisterAddonDef when imported.
//   - The key follows the "<addon>-vN" convention and will not collide with an
//     in-core addon key.
//   - Consumers that process documents declaring the key import the module, so
//     the strict "$addons must be registered" runtime check still passes.
//
// Listing a key here only makes it a valid `$addons` value in the JSON Schema;
// the module must still be imported for Validate/Calculate to succeed.
func init() {
	const frCTCModule = "github.com/invopop/gobl.fr.ctc"
	tax.RegisterApprovedAddon(&tax.ExternalAddon{
		Key: "fr-ctc-flow2-v1",
		Name: i18n.String{
			i18n.EN: "France CTC Flow 2 (B2B Clearance)",
			i18n.FR: "France CTC Flux 2 (Dédouanement B2B)",
		},
		Module: frCTCModule,
	})
	tax.RegisterApprovedAddon(&tax.ExternalAddon{
		Key: "fr-ctc-flow6-v1",
		Name: i18n.String{
			i18n.EN: "France CTC Flow 6 (Cycle de Vie)",
			i18n.FR: "France CTC Flux 6 (Cycle de Vie)",
		},
		Module: frCTCModule,
	})
	tax.RegisterApprovedAddon(&tax.ExternalAddon{
		Key: "fr-ctc-flow10-v1",
		Name: i18n.String{
			i18n.EN: "France CTC Flow 10 (E-Reporting)",
			i18n.FR: "France CTC Flux 10 (E-Reporting)",
		},
		Module: frCTCModule,
	})
	tax.RegisterApprovedAddon(&tax.ExternalAddon{
		Key: "sa-zatca-v1",
		Name: i18n.String{
			i18n.EN: "Saudi Arabia ZATCA",
			i18n.AR: "هيئة الزكاة والضريبة والجمارك",
		},
		Module: "github.com/invopop/gobl.sa.zatca",
	})

	tax.RegisterApprovedAddon(&tax.ExternalAddon{
		Key: "dk-oioubl-v2",
		Name: i18n.String{
			i18n.EN: "Danish OIOUBL 2.1",
			i18n.DA: "Dansk OIOUBL 2.1",
		},
		Module: "github.com/invopop/gobl.dk.oioubl",
	})
	tax.RegisterApprovedAddon(&tax.ExternalAddon{
		Key: "pt-saft-v1",
		Name: i18n.String{
			i18n.EN: "Portugal SAF-T",
		},
		Module: "github.com/invopop/gobl.pt.saft",
	})
	tax.RegisterApprovedAddon(&tax.ExternalAddon{
		Key: "br-nfe-v4",
		Name: i18n.String{
			i18n.EN: "Brazil NF-e 4.00",
		},
		Module: "github.com/invopop/gobl.br.nfe",
	})
	tax.RegisterApprovedAddon(&tax.ExternalAddon{
		Key: "br-nfse-v1",
		Name: i18n.String{
			i18n.EN: "Brazil NFS-e 1.X",
		},
		Module: "github.com/invopop/gobl.br.nfse",
	})
	tax.RegisterApprovedAddon(&tax.ExternalAddon{
		Key: "mx-cfdi-v4",
		Name: i18n.String{
			i18n.EN: "Mexican SAT CFDI v4.X",
		},
		Module: "github.com/invopop/gobl.mx.cfdi",
	})
}
