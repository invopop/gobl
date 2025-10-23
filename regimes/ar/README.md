# 🇦🇷 GOBL Argentina Tax Regime

Argentina's tax regime implementation for GOBL, supporting the Argentine tax system administered by AFIP (Administración Federal de Ingresos Públicos).

## Overview

The Argentine tax system includes several key components:

- **IVA (Impuesto al Valor Agregado)**: Value Added Tax with multiple rates
- **CUIT/CUIL/CDI**: Tax identification numbers with check digit validation
- **Retention Taxes**: IVA Retenido, Ganancias, and Ingresos Brutos
- **Electronic Invoicing**: AFIP-compliant electronic invoice requirements

## Public Documentation

Official documentation and resources:

- [AFIP - Administración Federal de Ingresos Públicos](https://www.afip.gob.ar/)
- [Avalara Argentina VAT Guide](https://www.avalara.com/vatlive/en/country-guides/south-america/argentina/argentina-vat-compliance-and-rates.html)
- [Santander Trade - Argentina Tax System](https://santandertrade.com/en/portal/establish-overseas/argentina/tax-system)
- [Argentina Tax Guide for Expats](https://myargentinepassport.com/en/taxes-in-argentina-for-foreigners/)
- [AFIP Tax Retention System (SIRE)](https://www.afip.gob.ar/sire/percepciones-retenciones/)
- [AFIP Income Tax Withholding](https://www.afip.gob.ar/genericos/guiavirtual/directorio_subcategoria_nivel3.aspx?id_nivel1=563id_nivel2%3D607&id_nivel3=686)
- [CUIT/CUIL Validation Algorithm](https://whiz.tools/en/legal-business/argentinian-cuit-cuil-generator-validator)
- [Argentina Tax ID Guide](https://lookuptax.com/docs/tax-identification-number/Argentina-tax-id-guide)

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
  - 20, 23 (males)
  - 27, 28 (females)

#### CDI (Clave de Identificación)
- **Used by**: Foreign residents without CUIT/CUIL

#### Validation Algorithm

CUIT/CUIL validation uses modulo 11 with specific multipliers:

```
Multipliers: [5, 4, 3, 2, 7, 6, 5, 4, 3, 2]

1. Multiply each of the first 10 digits by the corresponding multiplier
2. Sum all results
3. Calculate: checkDigit = 11 - (sum % 11)
4. Special cases:
   - If checkDigit = 11, then checkDigit = 0
   - If checkDigit = 10, prefix must be adjusted:
     * Males: 20 → 23
     * Females: 27 → 28
     * Companies: 30 → 33
     And checkDigit = 9
```

**References**:
- [CUIT/CUIL Validation Tool](https://whiz.tools/en/legal-business/argentinian-cuit-cuil-generator-validator)
- [Tax ID Guide](https://lookuptax.com/docs/tax-identification-number/Argentina-tax-id-guide)
- [Algorithm Details](https://lib.rs/crates/ar_cuil_cuit_validator)

### Tax Categories

#### IVA (Impuesto al Valor Agregado)

VAT rates in Argentina:

| Rate | Percentage | Application |
|------|-----------|-------------|
| General | 21% | Most goods and services |
| Reduced | 10.5% | Essential goods (construction, medicine, transportation, food products) |
| Super-Reduced | 2.5% | Specific categories (capital goods) |
| Zero | 0% | Exports and specific categories |
| Exempt | 0% | Books, interest on loans, insurance |

**References**:
- Law 23.349 and modifications (effective April 1, 1995)
- [Avalara Argentina VAT Rates](https://www.avalara.com/vatlive/en/country-guides/south-america/argentina/argentina-vat-compliance-and-rates.html)
- [Santander Trade Tax System](https://santandertrade.com/en/portal/establish-overseas/argentina/tax-system)

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

Argentina requires electronic invoicing through AFIP's system for most transactions:

- **Factura Electrónica**: AFIP-approved electronic invoice format
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

## GOBL Code Examples

### Simple Invoice with IVA

```json
{
  "$schema": "https://gobl.org/draft-0/envelope",
  "head": {
    "uuid": "9d8eafd5-77be-11ec-b485-5405db9a3e49",
    "dig": {
      "alg": "sha256",
      "val": "..."
    }
  },
  "doc": {
    "$schema": "https://gobl.org/draft-0/bill/invoice",
    "regime": "AR",
    "currency": "ARS",
    "issue_date": "2024-01-15",
    "supplier": {
      "name": "Empresa Argentina S.A.",
      "tax_id": {
        "country": "AR",
        "code": "30714589840"
      }
    },
    "customer": {
      "name": "Cliente Ejemplo S.R.L.",
      "tax_id": {
        "country": "AR",
        "code": "30123456789"
      }
    },
    "lines": [
      {
        "i": 1,
        "quantity": "10",
        "item": {
          "name": "Producto de Ejemplo",
          "price": "1000.00"
        },
        "taxes": [
          {
            "cat": "VAT",
            "rate": "standard"
          }
        ]
      }
    ]
  }
}
```

### Invoice with IVA Retention

```json
{
  "lines": [
    {
      "i": 1,
      "quantity": "1",
      "item": {
        "name": "Servicio de Consultoría",
        "price": "50000.00"
      },
      "taxes": [
        {
          "cat": "VAT",
          "rate": "standard"
        }
      ]
    }
  ],
  "charges": [
    {
      "key": "RIVA",
      "reason": "IVA Retenido",
      "amount": "1050.00",
      "taxes": [
        {
          "cat": "RIVA",
          "percent": "10%"
        }
      ]
    }
  ]
}
```

## Extension Keys

Argentina-specific extensions can be added using GOBL's extension mechanism for:

- AFIP authorization codes (CAE/CAI)
- Point of sale identification
- Invoice type classification
- Special tax regime indicators

## Notes

- CUIT/CUIL numbers are stored without hyphens but are typically displayed as XX-XXXXXXXX-X
- IVA rates have been stable since April 1, 1995
- Electronic invoicing requirements vary by business size and tax classification
- Provincial Ingresos Brutos rates must be determined based on the specific province and activity

## Contributing

To update or improve the Argentina tax regime:

1. Ensure all changes reference official AFIP documentation
2. Include effective dates for any tax rate changes
3. Test validation rules with real CUIT/CUIL examples
4. Update documentation with clear explanations

## References

### Official Sources
- [AFIP - Administración Federal de Ingresos Públicos](https://www.afip.gob.ar/)
- [AFIP SIRE - Tax Retention System](https://www.afip.gob.ar/sire/percepciones-retenciones/)
- [AFIP Withholding Guide](https://www.afip.gob.ar/genericos/guiavirtual/directorio_subcategoria_nivel3.aspx?id_nivel1=563id_nivel2%3D607&id_nivel3=686)
- [AFIP Withholding Calculator](https://servicioscf.afip.gob.ar/calc-rg830/)

### Tax Information
- [Avalara Argentina VAT Guide](https://www.avalara.com/vatlive/en/country-guides/south-america/argentina/argentina-vat-compliance-and-rates.html)
- [Santander Trade - Argentina Tax System](https://santandertrade.com/en/portal/establish-overseas/argentina/tax-system)
- [Argentina Tax Guide for Expats](https://myargentinepassport.com/en/taxes-in-argentina-for-foreigners/)

### CUIT/CUIL Validation
- [CUIT/CUIL Validation Tool](https://whiz.tools/en/legal-business/argentinian-cuit-cuil-generator-validator)
- [Argentina Tax ID Guide](https://lookuptax.com/docs/tax-identification-number/Argentina-tax-id-guide)
- [Tax ID Validator (Rust)](https://lib.rs/crates/ar_cuil_cuit_validator)
- [CUIT Validation Reference](https://www.lawebdelprogramador.com/codigo/Visual-Basic/160-Verificar-el-CUIT-CUIL-Argentina.html)

### Regulations
- IVA Law 23.349 and modifications
- RG 2854/2010 (IVA Retention)
- RG 830/2000 (Income Tax Withholding)
- RG 4003/2017 (Withholding modifications)
