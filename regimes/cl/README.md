# Chile (CL) - Tax Regime

Chile's tax regime implementation for GOBL based on the requirements of the Chilean Tax Authority (Servicio de Impuestos Internos - SII).

## Tax Categories

### IVA (Impuesto al Valor Agregado)

Chile's Value Added Tax (VAT) is called IVA. The consumption tax was introduced in December 1974 through [Decreto Ley Nº 825](https://www.sii.cl/normativa_legislacion/sobreventasyservicios.pdf).

- **Standard Rate**: 19% (effective October 1, 2003)
- **Historical Rate**: 18% (prior to October 1, 2003)

The current 19% rate was established by [Ley 19888](https://www.bcn.cl/leychile/Navegar?idNorma=213493), published August 13, 2003, which amended Article 14 of the VAT law to increase the rate from 18% to 19%. Chile applies a single VAT rate with no reduced rates.

**Source**: [Law No. 21.210](https://www.sii.cl/vat/faq1_eng.html) confirmed the 19% rate applies to digital services provided by non-resident taxpayers.

## Tax Identity

Chile uses the **RUT** (Rol Único Tributario) as the primary tax identification number for both individuals and legal entities, issued and maintained by the SII.

### RUT Format

The RUT consists of:
- **Numeric portion**: 6 to 8 digits
- **Check digit**: Single character (0-9 or K)
- **Total length**: 7 to 9 characters (after normalization)
- **Display format**: `XX.XXX.XXX-Y` where dots and hyphen are used for readability
- **Normalized format**: Removes all separators (e.g., `713254975`, `77668208K`)

Modern RUTs typically have 8-9 total digits, while older RUTs may have 7. The SII now handles RUT numbers ranging from 7 to 9 digits in total length.

### Check Digit Calculation

The check digit is calculated using the **Modulo 11 algorithm**, which is the exclusive validation method mandated by the SII for RUT verification in Chile.

**Implementation Note**: This algorithm is implemented in the `calculateRUTCheckDigit` function in `tax_identity.go`.

**Sources**:
- [Rol Único Tributario - Wikipedia](https://es.wikipedia.org/wiki/Rol_%C3%9Anico_Tributario)
- [Modificaciones para manejo de RUT de 7 a 9 dígitos](https://help.getcirrus.com/es/articles/8515715-modificaciones-para-manejo-de-rut-de-7-a-9-digitos-chile)

### Validation and Normalization

GOBL implements comprehensive RUT validation and normalization to ensure data consistency and compliance with SII standards.

#### Normalization Process

When processing a RUT, GOBL automatically:

1. **Removes formatting characters**: Strips all dots (`.`) and hyphens (`-`)
2. **Normalizes case**: Converts lowercase `k` to uppercase `K` for the check digit

**Examples**:
- `71.325.497-5` → `713254975`
- `7.766.820-8K` → `77668208K`
- `12.345.678-5` → `123456785`

## Electronic Invoicing

Chile has a mandatory electronic invoicing system managed by SII, based on **DTE** (Documentos Tributarios Electrónicos - Electronic Tax Documents).

### Implementation Timeline

- **2001**: E-invoicing implemented as voluntary scheme
- **2014**: VAT e-invoicing made mandatory for established organizations
- **March 2018**: B2B e-invoicing made mandatory
- **March 2021**: B2C e-invoicing made mandatory

### Validation Process

The Chilean system uses **prior validation** by the SII:
1. DTEs must be transmitted to SII for validation before being sent to customers
2. Once validated, SII returns the document to the issuer
3. Document is forwarded to recipient
4. Recipient has 8 days to accept or reject
5. If no action taken, document is considered tacitly accepted

### Document Types

Main DTEs required by SII:
- **Factura Electrónica** (Electronic Invoice)
- **Boleta Electrónica** (Electronic Receipt)
- **Nota de Crédito Electrónica** (Electronic Credit Note)
- **Nota de Débito Electrónica** (Electronic Debit Note)
- **Guía de Despacho** (Dispatch Guide)
- **Factura de Exportación** (Export Invoice)

### Archiving Requirements

- DTEs must be archived for **6 years**
- Must be kept in the **XML format** validated by SII

**Sources**:
- [SII Electronic Invoicing System](https://www.chileatiende.gob.cl/fichas/13505-emision-de-documentos-tributarios-electronicos-dte-sistema-de-facturacion-gratuito-del-sii)
- [Electronic Invoicing in Chile - EDICOM](https://edicomgroup.com/electronic-invoicing/chile)
- [What is Electronic Invoicing - SII](https://www.sii.cl/factura_electronica/que_es_fact_elect.htm)

## References

### Official Sources
- [Servicio de Impuestos Internos (SII)](https://www.sii.cl/) - Chilean Tax Authority
- [Decreto Ley Nº 825 - VAT Law](https://www.sii.cl/normativa_legislacion/sobreventasyservicios.pdf) - Law on VAT and Services (1974)
- [Ley 19888](https://www.bcn.cl/leychile/Navegar?idNorma=213493) - Law establishing 19% VAT rate (2003)
- [IVA Digital - SII](https://www.sii.cl/vat/faq1_eng.html) - Official VAT information (English)

### Tax Identification
- [Rol Único Tributario - Wikipedia](https://es.wikipedia.org/wiki/Rol_%C3%9Anico_Tributario)
- [RUT Format Changes (7-9 digits) - Cirrus Help Center](https://help.getcirrus.com/es/articles/8515715-modificaciones-para-manejo-de-rut-de-7-a-9-digitos-chile)

### Electronic Invoicing
- [DTE System - ChileAtiende](https://www.chileatiende.gob.cl/fichas/13505-emision-de-documentos-tributarios-electronicos-dte-sistema-de-facturacion-gratuito-del-sii)
- [Electronic Invoicing Overview - EDICOM](https://edicomgroup.com/electronic-invoicing/chile)
- [SII Electronic Invoice Information](https://www.sii.cl/factura_electronica/que_es_fact_elect.htm)
