package ar

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
)

// Tax tags that can be applied in Argentina according to AFIP regulations.
// These tags represent different invoice types, tax regimes, and special scenarios.
//
// References:
// - AFIP Electronic Invoicing: https://www.afip.gob.ar/facturacion/regimen-general/
// - Invoice Types Guide: https://www.tiendanube.com/blog/tipos-de-factura/
// - Tax Regimes: https://www.afip.gob.ar/monotributo/
const (
	// Invoice Type Tags - Tipos de Comprobantes
	// Reference: AFIP Facturación Electrónica
	// Source: https://www.afip.gob.ar/facturacion/regimen-general/

	TagInvoiceTypeA cbc.Key = "invoice-type-a" // Factura Tipo A - Responsable Inscripto to Responsable Inscripto
	TagInvoiceTypeB cbc.Key = "invoice-type-b" // Factura Tipo B - Responsable Inscripto to Final Consumer/Monotributista
	TagInvoiceTypeC cbc.Key = "invoice-type-c" // Factura Tipo C - Monotributista or Exempt issuer
	TagInvoiceTypeE cbc.Key = "invoice-type-e" // Factura Tipo E - Export invoice
	TagInvoiceTypeM cbc.Key = "invoice-type-m" // Factura Tipo M - Monotributo issuer

	// Tax Regime Tags - Categorías Fiscales
	// Reference: AFIP Tax Categories
	// Source: https://www.afip.gob.ar/monotributo/

	TagResponsableInscripto cbc.Key = "responsable-inscripto" // Registered taxpayer (full IVA obligations)
	TagMonotributo          cbc.Key = "monotributo"           // Simplified tax regime for small businesses
	TagExento               cbc.Key = "exento"                // VAT exempt entity
	TagConsumidorFinal      cbc.Key = "consumidor-final"      // Final consumer (no tax ID required)
	TagNoResponsable        cbc.Key = "no-responsable"        // Not responsible for IVA collection

	// Special Regime Tags - Regímenes Especiales
	// Reference: AFIP Special Regimes

	TagExportServices cbc.Key = "export-services" // Export of services
	TagExportGoods    cbc.Key = "export-goods"    // Export of goods
)

// invoiceTags defines the available tags for Argentine invoices
func invoiceTags() *tax.TagSet {
	return &tax.TagSet{
		Schema: bill.ShortSchemaInvoice,
		List: []*cbc.Definition{
			// Invoice Type Definitions
			{
				Key: TagInvoiceTypeA,
				Name: i18n.String{
					i18n.EN: "Invoice Type A",
					i18n.ES: "Factura Tipo A",
				},
				Desc: i18n.String{
					i18n.EN: "Invoice issued by Responsable Inscripto to another Responsable Inscripto. VAT is detailed separately, indicating the net amount and corresponding tax.",
					i18n.ES: "Factura emitida por Responsable Inscripto a otro Responsable Inscripto. El IVA se discrimina por separado, indicando el monto neto y el impuesto correspondiente.",
				},
			},
			{
				Key: TagInvoiceTypeB,
				Name: i18n.String{
					i18n.EN: "Invoice Type B",
					i18n.ES: "Factura Tipo B",
				},
				Desc: i18n.String{
					i18n.EN: "Invoice issued by Responsable Inscripto to final consumers, monotributistas, or exempt entities. VAT is included in the total but not separately detailed.",
					i18n.ES: "Factura emitida por Responsable Inscripto a consumidores finales, monotributistas o exentos. El IVA está incluido en el total pero no discriminado.",
				},
			},
			{
				Key: TagInvoiceTypeC,
				Name: i18n.String{
					i18n.EN: "Invoice Type C",
					i18n.ES: "Factura Tipo C",
				},
				Desc: i18n.String{
					i18n.EN: "Invoice issued by monotributista or VAT-exempt entities. No VAT is charged or detailed.",
					i18n.ES: "Factura emitida por monotributista o sujeto exento de IVA. No se cobra ni discrimina IVA.",
				},
			},
			{
				Key: TagInvoiceTypeE,
				Name: i18n.String{
					i18n.EN: "Invoice Type E",
					i18n.ES: "Factura Tipo E",
				},
				Desc: i18n.String{
					i18n.EN: "Export invoice. Issued for export operations of goods or services. VAT rate is typically 0%.",
					i18n.ES: "Factura de exportación. Se emite para operaciones de exportación de bienes o servicios. La alícuota de IVA es típicamente 0%.",
				},
			},
			{
				Key: TagInvoiceTypeM,
				Name: i18n.String{
					i18n.EN: "Invoice Type M",
					i18n.ES: "Factura Tipo M",
				},
				Desc: i18n.String{
					i18n.EN: "Invoice issued by monotributista. Simplified invoice format for small taxpayers under the Monotributo regime.",
					i18n.ES: "Factura emitida por monotributista. Formato de factura simplificado para pequeños contribuyentes bajo el régimen de Monotributo.",
				},
			},

			// Tax Regime Definitions
			{
				Key: TagResponsableInscripto,
				Name: i18n.String{
					i18n.EN: "Registered Taxpayer",
					i18n.ES: "Responsable Inscripto",
				},
				Desc: i18n.String{
					i18n.EN: "Taxpayer registered for VAT and Income Tax. Must charge VAT on sales and can claim VAT credits on purchases.",
					i18n.ES: "Contribuyente inscripto en IVA e Impuesto a las Ganancias. Debe cobrar IVA en las ventas y puede computar crédito fiscal en las compras.",
				},
			},
			{
				Key: TagMonotributo,
				Name: i18n.String{
					i18n.EN: "Monotributo",
					i18n.ES: "Monotributo",
				},
				Desc: i18n.String{
					i18n.EN: "Simplified tax regime for small businesses. Unifies VAT, Income Tax, and social security contributions into a single monthly payment.",
					i18n.ES: "Régimen simplificado para pequeños contribuyentes. Unifica IVA, Impuesto a las Ganancias y aportes previsionales en un único pago mensual.",
				},
			},
			{
				Key: TagExento,
				Name: i18n.String{
					i18n.EN: "VAT Exempt",
					i18n.ES: "Exento de IVA",
				},
				Desc: i18n.String{
					i18n.EN: "Entity exempt from VAT. Does not charge VAT on sales and cannot claim VAT credits on purchases.",
					i18n.ES: "Sujeto exento de IVA. No cobra IVA en las ventas y no puede computar crédito fiscal en las compras.",
				},
			},
			{
				Key: TagConsumidorFinal,
				Name: i18n.String{
					i18n.EN: "Final Consumer",
					i18n.ES: "Consumidor Final",
				},
				Desc: i18n.String{
					i18n.EN: "End consumer who does not require tax identification. Purchases include VAT but cannot claim tax credits.",
					i18n.ES: "Consumidor final que no requiere identificación tributaria. Las compras incluyen IVA pero no puede computar crédito fiscal.",
				},
			},
			{
				Key: TagNoResponsable,
				Name: i18n.String{
					i18n.EN: "Not Responsible",
					i18n.ES: "No Responsable",
				},
				Desc: i18n.String{
					i18n.EN: "Entity not responsible for VAT collection. Typically applies to certain public entities or special cases.",
					i18n.ES: "Sujeto no responsable de la percepción de IVA. Típicamente aplica a ciertos entes públicos o casos especiales.",
				},
			},

			// Special Regime Tags
			{
				Key: TagExportServices,
				Name: i18n.String{
					i18n.EN: "Export of Services",
					i18n.ES: "Exportación de Servicios",
				},
				Desc: i18n.String{
					i18n.EN: "Services provided to foreign entities or individuals. Typically subject to 0% VAT rate.",
					i18n.ES: "Servicios prestados a entidades o personas del exterior. Típicamente sujetos a alícuota de IVA del 0%.",
				},
			},
			{
				Key: TagExportGoods,
				Name: i18n.String{
					i18n.EN: "Export of Goods",
					i18n.ES: "Exportación de Bienes",
				},
				Desc: i18n.String{
					i18n.EN: "Goods exported outside Argentina. Subject to 0% VAT rate.",
					i18n.ES: "Bienes exportados fuera de Argentina. Sujetos a alícuota de IVA del 0%.",
				},
			},

			// Standard GOBL tags commonly used
			{
				Key: tax.TagSimplified,
				Name: i18n.String{
					i18n.EN: "Simplified Invoice",
					i18n.ES: "Factura Simplificada",
				},
				Desc: i18n.String{
					i18n.EN: "Simplified invoice format. Used for small transactions or when full customer details are not available.",
					i18n.ES: "Formato de factura simplificada. Se utiliza para transacciones pequeñas o cuando no se dispone de todos los datos del cliente.",
				},
			},
			{
				Key: tax.TagReverseCharge,
				Name: i18n.String{
					i18n.EN: "Reverse Charge",
					i18n.ES: "Inversión del Sujeto Pasivo",
				},
				Desc: i18n.String{
					i18n.EN: "The customer is responsible for VAT payment instead of the supplier. Applies in specific circumstances defined by AFIP.",
					i18n.ES: "El cliente es responsable del pago del IVA en lugar del proveedor. Se aplica en circunstancias específicas definidas por AFIP.",
				},
			},
		},
	}
}

// invoiceScenarios defines tax scenarios that automatically add legal notes
// and enforce specific validation rules for Argentine invoices
func invoiceScenarios() *tax.ScenarioSet {
	return &tax.ScenarioSet{
		Schema: bill.ShortSchemaInvoice,
		List: []*tax.Scenario{
			// Invoice Type A - Responsable Inscripto to Responsable Inscripto
			{
				Tags: []cbc.Key{TagInvoiceTypeA},
				Note: &tax.ScenarioNote{
					Key:  org.NoteKeyLegal,
					Src:  TagInvoiceTypeA,
					Text: "Factura Tipo A - Emitida por Responsable Inscripto. IVA discriminado. / Invoice Type A - Issued by Registered Taxpayer. VAT itemized.",
				},
			},

			// Invoice Type B - Responsable Inscripto to Final Consumer
			{
				Tags: []cbc.Key{TagInvoiceTypeB},
				Note: &tax.ScenarioNote{
					Key:  org.NoteKeyLegal,
					Src:  TagInvoiceTypeB,
					Text: "Factura Tipo B - Emitida por Responsable Inscripto a Consumidor Final o Monotributista. / Invoice Type B - Issued by Registered Taxpayer to Final Consumer or Monotributo.",
				},
			},

			// Invoice Type C - Monotributista or Exempt
			{
				Tags: []cbc.Key{TagInvoiceTypeC},
				Note: &tax.ScenarioNote{
					Key:  org.NoteKeyLegal,
					Src:  TagInvoiceTypeC,
					Text: "Factura Tipo C - Emitida por Monotributista o Sujeto Exento. Sin discriminación de IVA. / Invoice Type C - Issued by Monotributo or Exempt Entity. VAT not itemized.",
				},
			},

			// Invoice Type E - Export
			{
				Tags: []cbc.Key{TagInvoiceTypeE},
				Note: &tax.ScenarioNote{
					Key:  org.NoteKeyLegal,
					Src:  TagInvoiceTypeE,
					Text: "Factura Tipo E - Operación de Exportación. Alícuota de IVA 0%. / Invoice Type E - Export Operation. VAT rate 0%.",
				},
			},

			// Export of Services
			{
				Tags: []cbc.Key{TagExportServices},
				Note: &tax.ScenarioNote{
					Key:  org.NoteKeyLegal,
					Src:  TagExportServices,
					Text: "Exportación de Servicios - Operación exenta de IVA conforme al Artículo 43 de la Ley 23.349. / Export of Services - VAT exempt operation pursuant to Article 43 of Law 23.349.",
				},
			},

			// Export of Goods
			{
				Tags: []cbc.Key{TagExportGoods},
				Note: &tax.ScenarioNote{
					Key:  org.NoteKeyLegal,
					Src:  TagExportGoods,
					Text: "Exportación de Bienes - Operación exenta de IVA conforme al Artículo 43 de la Ley 23.349. / Export of Goods - VAT exempt operation pursuant to Article 43 of Law 23.349.",
				},
			},

			// Monotributo regime
			{
				Tags: []cbc.Key{TagMonotributo},
				Note: &tax.ScenarioNote{
					Key:  org.NoteKeyLegal,
					Src:  TagMonotributo,
					Text: "Monotributo - Régimen Simplificado para Pequeños Contribuyentes (RS). Comprobante no válido como crédito fiscal. / Monotributo - Simplified Regime for Small Taxpayers. Document not valid as tax credit.",
				},
			},

			// VAT Exempt
			{
				Tags: []cbc.Key{TagExento},
				Note: &tax.ScenarioNote{
					Key:  org.NoteKeyLegal,
					Src:  TagExento,
					Text: "IVA Exento - Operación exenta del Impuesto al Valor Agregado. / VAT Exempt - Operation exempt from Value Added Tax.",
				},
			},

			// Reverse Charge
			{
				Tags: []cbc.Key{tax.TagReverseCharge},
				Note: &tax.ScenarioNote{
					Key:  org.NoteKeyLegal,
					Src:  tax.TagReverseCharge,
					Text: "Inversión del Sujeto Pasivo - El cliente es responsable del pago del IVA. / Reverse Charge - Customer is responsible for VAT payment.",
				},
			},

			// Simplified Invoice
			{
				Tags: []cbc.Key{tax.TagSimplified},
				Note: &tax.ScenarioNote{
					Key:  org.NoteKeyLegal,
					Src:  tax.TagSimplified,
					Text: "Factura Simplificada - Comprobante simplificado conforme a normativa AFIP. / Simplified Invoice - Simplified document according to AFIP regulations.",
				},
			},
		},
	}
}
