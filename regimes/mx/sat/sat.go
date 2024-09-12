// Package sat provides specifications from the Mexican SAT
// (Servicio de Administración Tributaria) tax authority.
package sat

import "github.com/invopop/gobl/cbc"

// Official SAT codes to include in stamps.
const (
	StampSignature   cbc.Key = "sat-sig"          // Signature - Sello Digital del SAT (optional)
	StampSerial      cbc.Key = "sat-serial"       // Cert Serial - Número de Certificado SAT
	StampTimestamp   cbc.Key = "sat-timestamp"    // Timestamp - Fecha y hora de certificación del SAT
	StampUUID        cbc.Key = "sat-uuid"         // Folio Fiscal
	StampURL         cbc.Key = "sat-url"          // URL QR Code
	StampProviderRFC cbc.Key = "sat-provider-rfc" // Provider RFC - RFC del Proveedor de Certificación
	StampChain       cbc.Key = "sat-chain"        // Cadena original del complemento de certificación digital del SAT
)

// Custom keys used typically in meta or codes information.
const (
	KeyFormaPago    cbc.Key = "sat-forma-pago"    // for mapping to c_FormaPago’s codes
	KeyTipoRelacion cbc.Key = "sat-tipo-relacion" // for mapping to c_TipoRelacion’s codes
	KeyImpuesto     cbc.Key = "sat-impuesto"      // for mapping to c_Impuesto’s codes
)
