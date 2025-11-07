# 🇦🇷 GOBL Argentina Tax Regime

Argentina's tax regime implementation for GOBL, supporting the Argentine tax system administered by ARCA (Agencia de Recaudación y Control Aduanero).

## Overview

The Argentine tax system includes several key components:

- **IVA (Impuesto al Valor Agregado)**: Value Added Tax with multiple rates
- **CUIT/CUIL/CDI**: Tax identification numbers with check digit validation
- **Retention Taxes**: IVA Retenido, Ganancias, and Ingresos Brutos
- **Electronic Invoicing**: AFIP-compliant electronic invoice requirements

## Argentina-specific Requirements

### Tax Identification

Argentina uses three main types of tax identification numbers:

#### CUIT (Clave Única de Identificación Tributaria)
- **Format**: 11 digits (XX-XXXXXXXX-X)
- **Used by**: Companies and legal entities
- **Prefixes**: 30, 33 (conflict resolution), 34 (foreign entities)

#### CUIL (Clave Única de Identificación Laboral)
- **Format**: 11 digits (XX-XXXXXXXX-X)
- **Used by**: Individuals and employees
- **Prefixes**:
  - 20 (males)
  - 27 (females)
  - 23 (conflict resolution)

#### CDI (Clave de Identificación)
- **Used by**: Foreign residents without CUIT/CUIL (Not validated)

### Tax Categories

#### IVA (Impuesto al Valor Agregado)

VAT rates in Argentina:

| Rate | Percentage | Application |
|------|-----------|-------------|
| Increased | 27% | Gas, water and telecom services |
| General | 21% | Most goods and services |
| Reduced | 10.5% | Essential goods (construction, medicine, transportation, food products) |

#### Retention Taxes

**IVA Retenido (Retained VAT)**
- Variable rates based on taxpayer category and registration status
- Reference: AFIP RG 2854/2010 and modifications
- Source: [AFIP SIRE](https://www.afip.gob.ar/sire/percepciones-retenciones/)

**Ganancias (Income Tax Withholding)**
- Applied to payments for services
- Rates range from 0.5% to 35% depending on service type
- Reference: AFIP RG 830/2000, RG 4003/2017
- Sources:
  - [AFIP Withholding Guide](https://www.afip.gob.ar/genericos/guiavirtual/directorio_subcategoria_nivel3.aspx?id_nivel1=563id_nivel2%3D607&id_nivel3=686)
  - [Withholding Calculator](https://servicioscf.afip.gob.ar/calc-rg830/)

**Ingresos Brutos (Gross Income Tax)**
- Provincial tax with rates set by each jurisdiction
- Typical rates: 1% to 5% depending on province and activity
- Reference: Provincial tax codes (each province has its own regulations)

### Invoice Requirements

#### Electronic Invoicing

Argentina requires electronic invoicing through ARCA's system for most transactions:

- **Factura Electrónica**: ARCA-approved electronic invoice format
- **CAE/CAI**: Código de Autorización Electrónico (Electronic Authorization Code)
- **Point of Sale (Punto de Venta)**: Required for invoice numbering

#### Invoice Types

Common invoice types in Argentina:

- **Tipo A**: Issued by Responsable Inscripto to another Responsable Inscripto
- **Tipo B**: Issued by Responsable Inscripto to Monotributista or final consumer
- **Tipo C**: Issued by Monotributista or exempt entities
- **Tipo E**: Export invoices
- **Credit Notes**: Notas de Crédito (corrective documents)

### Tax Regimes

Argentina has different tax regime classifications:

- **Responsable Inscripto**: Registered taxpayer (full IVA obligations)
- **Monotributo**: Simplified tax regime for small businesses
- **Exento**: Exempt from IVA
- **No Responsable**: Not responsible for IVA collection
- **Consumidor Final**: Final consumer (no tax ID required)

## References

### Official Sources
- [AFIP - Administración Federal de Ingresos Públicos](https://www.afip.gob.ar/)
- [Invoice type and mandatory information](https://www.argentina.gob.ar/normativa/recurso/54461/259-98a/htm)

### Regulations
- IVA Law 23.349 and modifications
- RG 2854/2010 (IVA Retention)
- RG 830/2000 (Income Tax Withholding)
- RG 4003/2017 (Withholding modifications)
