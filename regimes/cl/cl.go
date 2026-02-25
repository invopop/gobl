// Package cl handles tax regime data for Chile.
//
// This package implements validation and normalization for Chilean tax documents
// according to the requirements of the SII (Servicio de Impuestos Internos).
//
// Key regulations:
//   - Decreto Ley Nº 825 (1974): VAT Law - https://www.sii.cl/normativa_legislacion/sobreventasyservicios.pdf
//   - Ley 19888 (2003): Established 19% VAT rate - https://www.bcn.cl/leychile/Navegar?idNorma=213493
//   - Resolution 36/2024: Item description requirements (effective July 2024)
//
// Electronic invoicing references:
//   - SII Electronic Invoicing Portal: https://www.sii.cl/factura_electronica/
//   - DTE Format Specification: https://www.sii.cl/factura_electronica/factura_mercado/formato_dte.pdf
//   - Boleta Format Specification: https://www.sii.cl/factura_electronica/factura_mercado/formato_boletas_elec_202306.pdf
package cl

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/tax"
)

func init() {
	tax.RegisterRegimeDef(New())
}

// New provides the tax region definition
func New() *tax.RegimeDef {
	return &tax.RegimeDef{
		Country:   "CL",
		Currency:  currency.CLP,
		TaxScheme: tax.CategoryVAT,
		Name: i18n.String{
			i18n.EN: "Chile",
			i18n.ES: "Chile",
		},
		Description: i18n.String{
			i18n.EN: here.Doc(`
				Chile's tax system is administered by the SII (Servicio de Impuestos Internos), which oversees the collection of IVA (Impuesto al Valor Agregado), the country's value-added tax.

				Taxpayers are identified using the RUT (Rol Único Tributario), a 6-8 digit number with a check digit calculated using the modulo 11 algorithm. The check digit can be 0-9 or K, and RUTs are formatted as XX.XXX.XXX-Y. Both supplier and customer RUT are mandatory for Facturas Electrónicas (B2B invoices), while customer RUT is optional for Boletas Electrónicas (B2C receipts). Supplier address is always required, and customer address is required for B2B transactions.

				Chile applies a single standard IVA rate of 19%, effective since October 1, 2003, when it was increased from 18% by Ley 19888. Unlike many other countries, Chile does not have reduced or super-reduced VAT rates.

				Electronic invoicing has been mandatory in Chile since 2018 for B2B transactions and since 2021 for B2C transactions. The system is based on DTEs (Documentos Tributarios Electrónicos - Electronic Tax Documents) that must be validated by the SII before being sent to the recipient. The validation process is known as "prior validation," where documents are transmitted to the SII first, validated, returned to the issuer, and then forwarded to the customer. Recipients have 8 days to accept or reject documents; otherwise, they are considered tacitly accepted.

				Common document types include Factura Electrónica (electronic invoice for B2B), Boleta Electrónica (electronic receipt for B2C), Nota de Crédito Electrónica (electronic credit note), Nota de Débito Electrónica (electronic debit note), and Guía de Despacho (dispatch guide). Use the "simplified" tag to indicate Boletas Electrónicas, which have relaxed customer identification requirements. All DTEs must be archived for 6 years in the XML format validated by the SII.

				Additional SII requirements include Giro Comercial (business activity) for both parties, complete address details (Dirección, Ciudad, Comuna), and clear item descriptions (Resolution 36/2024, effective July 2024).
			`),
			i18n.ES: here.Doc(`
				El sistema tributario de Chile es administrado por el SII (Servicio de Impuestos Internos), que supervisa la recaudación del IVA (Impuesto al Valor Agregado).

				Los contribuyentes se identifican mediante el RUT (Rol Único Tributario), un número de 6-8 dígitos con un dígito verificador calculado mediante el algoritmo módulo 11. El dígito verificador puede ser 0-9 o K, y los RUT se formatean como XX.XXX.XXX-Y. El RUT del emisor y receptor es obligatorio para Facturas Electrónicas (B2B), mientras que el RUT del receptor es opcional para Boletas Electrónicas (B2C). La dirección del emisor es siempre requerida, y la dirección del receptor es requerida para transacciones B2B.

				Chile aplica una tasa única de IVA del 19%, vigente desde el 1 de octubre de 2003, cuando fue aumentada del 18% mediante la Ley 19888. A diferencia de muchos otros países, Chile no tiene tasas reducidas o super-reducidas de IVA.

				La facturación electrónica es obligatoria en Chile desde 2018 para transacciones B2B y desde 2021 para transacciones B2C. El sistema se basa en DTEs (Documentos Tributarios Electrónicos) que deben ser validados por el SII antes de enviarse al receptor. El proceso de validación se conoce como "validación previa", donde los documentos se transmiten primero al SII, se validan, se devuelven al emisor y luego se reenvían al cliente. Los receptores tienen 8 días para aceptar o rechazar documentos; de lo contrario, se consideran aceptados tácitamente.

				Los tipos de documentos comunes incluyen Factura Electrónica (factura para B2B), Boleta Electrónica (boleta para B2C), Nota de Crédito Electrónica, Nota de Débito Electrónica y Guía de Despacho. Use la etiqueta "simplified" para indicar Boletas Electrónicas, que tienen requisitos de identificación del cliente más flexibles. Todos los DTEs deben archivarse durante 6 años en el formato XML validado por el SII.

				Requisitos adicionales del SII incluyen Giro Comercial (actividad económica) para ambas partes, detalles completos de dirección (Dirección, Ciudad, Comuna), y descripciones claras de los ítems (Resolución 36/2024, vigente desde julio 2024).
			`),
		},
		TimeZone:   "America/Santiago",
		Validator:  Validate,
		Normalizer: Normalize,
		Categories: taxCategories,
	}
}

// Normalize will attempt to clean the object passed to it.
func Normalize(doc any) {
	switch obj := doc.(type) {
	case *tax.Identity:
		normalizeTaxIdentity(obj)
	}
}

// Validate checks the document type and determines if it can be validated.
func Validate(doc interface{}) error {
	switch obj := doc.(type) {
	case *bill.Invoice:
		return validateInvoice(obj)
	case *tax.Identity:
		return validateTaxIdentity(obj)
	}
	return nil
}
