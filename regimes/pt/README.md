# 🇵🇹 GOBL Portugal Tax Regime

Portugal doesn't have an e-invoicing format per se. Tax information is reported electronically to the AT (Autoridade Tributária e Aduaneira) either periodically in batches via a SAF-T (PT) report or individually in real time via a web service.

Find example PT GOBL files in the [`examples`](../../examples/pt) (uncalculated documents) and [`examples/out`](../../examples/pt/out) (calculated envelopes) subdirectories.

## Public Documentation

- [Portaria n.o 302/2016 – SAF-T Data Structure & Taxonomies](https://info.portaldasfinancas.gov.pt/pt/informacao_fiscal/legislacao/diplomas_legislativos/Documents/Portaria_302_2016.pdf)
- [Portaria n.o 195/2020 – Comunicação de Séries Documentais, Aspetos Específicos](https://info.portaldasfinancas.gov.pt/pt/apoio_contribuinte/Faturacao/Comunicacao_Series_ATCUD/Documents/Comunicacao_de_Series_Documentais_Manual_de_Integracao_de_SW_Aspetos_Genericos.pdf)
- [Portaria n.o 195/2020 – Especificações Técnicas Código QR](https://info.portaldasfinancas.gov.pt/pt/apoio_contribuinte/Novas_regras_faturacao/Documents/Especificacoes_Tecnicas_Codigo_QR.pdf)
- [Comunicação dos elementos dos documentos de faturação à AT, por webservice](https://info.portaldasfinancas.gov.pt/pt/apoio_contribuinte/Faturacao/Fatcorews/Documents/Comunicacao_dos_elementos_dos_documentos_de_faturacao.pdf)

## Portugal-specific Requirements

### `InvoiceType` (Tipo de documento)

AT's `InvoiceType` (Tipo de documento) specifies the type of a Portuguese tax document. The following table lists all the supported invoice types and how GOBL will map them with a combination of invoice type and tax tags:

| Code | Name                                                                     | GOBL Type     | GOBL Tax Tag      |
| ---- | ------------------------------------------------------------------------ | ------------- | ----------------- |
| FT   | Fatura, emitida nos termos do artigo 36.o do Código do IVA              | `standard`    |                   |
| FS   | Fatura simplificada, emitida nos termos do artigo 40.o do Código do IVA | `standard`    | `simplified`      |
| FR   | Fatura-recibo                                                            | `standard`    | `invoice-receipt` |
| ND   | Nota de débito                                                          | `credit-note` |                   |
| NC   | Nota de crédito                                                         | `debit-note`  |                   |

### `TaxCountryRegion` (País ou região do imposto)

AT's `TaxCountryRegion` (País ou região do imposto) specifies the region of taxation (Portugal mainland, Açores or Madeira) in a Portuguese invoice. Each region has their own tax rates which can be determined automatically.

To set the specific a region different to Portugal mainland, the `pt-region` extension of each line's VAT tax should be set to one of the following values:

| Code  | Description                                         |
| ----- | --------------------------------------------------- |
| PT    | Mainland Portugal (default, no need to be explicit) |
| PT-AC | Açores                                              |
| PT-MA | Madeira                                             |

### VAT Tax Rates

The AT `TaxCode` (Código do imposto) is required for invoice items that apply VAT. GOBL helps determine this code using the `rate` field, which in Portuguese invoices is required. The following table lists the supported tax codes and how GOBL will map them:

| Code | Name            | GOBL Tax Rate                         |
| ---- | --------------- | ------------------------------------- |
| NOR  | Tipo Geral      | `standard`                            |
| INT  | Taxa Intermédia | `intermediate`                        |
| RED  | Taxa Reduzida   | `reduced`                             |
| ISE  | Isenta          | `exempt` + extension code (see below) |

AT's `TaxExemptionCode` (Código do motivo de isenção de imposto) is a code that specifies the reason the VAT tax is exempt in a Portuguese invoice. When the `exempt` tag is used, one of the following must be defined in the `ext` map's `pt-exemption-code` property:

| Code  | Description                                                                                              |
| ----- | -------------------------------------------------------------------------------------------------------- |
| `M01` | Artigo 16.°, n.° 6 do CIVA                                                                               |
| `M02` | Artigo 6.° do Decreto-Lei n.° 198/90, de 19 de junho                                                     |
| `M04` | Isento artigo 13.° do CIVA                                                                               |
| `M05` | Isento artigo 14.° do CIVA                                                                               |
| `M06` | Isento artigo 15.° do CIVA                                                                               |
| `M07` | Isento artigo 9.° do CIVA                                                                                |
| `M09` | IVA - não confere direito a dedução / Artigo 62.° alínea b) do CIVA                                      |
| `M10` | IVA - regime de isenção / Artigo 57.° do CIVA                                                            |
| `M11` | Regime particular do tabaco / Decreto-Lei n.° 346/85, de 23 de agosto                                    |
| `M12` | Regime da margem de lucro - Agências de viagens / Decreto-Lei n.° 221/85, de 3 de julho                  |
| `M13` | Regime da margem de lucro - Bens em segunda mão / Decreto-Lei n.° 199/96, de 18 de outubro               |
| `M14` | Regime da margem de lucro - Objetos de arte / Decreto-Lei n.° 199/96, de 18 de outubro                   |
| `M15` | Regime da margem de lucro - Objetos de coleção e antiguidades / Decreto-Lei n.° 199/96, de 18 de outubro |
| `M16` | Isento artigo 14.° do RITI                                                                               |
| `M19` | Outras isenções - Isenções temporárias determinadas em diploma próprio                                   |
| `M20` | IVA - regime forfetário / Artigo 59.°-D n.°2 do CIVA                                                     |
| `M21` | IVA - não confere direito à dedução (ou expressão similar) - Artigo 72.° n.° 4 do CIVA                   |
| `M25` | Mercadorias à consignação - Artigo 38.° n.° 1 alínea a) do CIVA                                          |
| `M30` | IVA - autoliquidação / Artigo 2.° n.° 1 alínea i) do CIVA                                                |
| `M31` | IVA - autoliquidação / Artigo 2.° n.° 1 alínea j) do CIVA                                                |
| `M32` | IVA - autoliquidação / Artigo 2.° n.° 1 alínea I) do CIVA                                                |
| `M33` | IVA - autoliquidação / Artigo 2.° n.° 1 alínea m) do CIVA                                                |
| `M40` | IVA - autoliquidação / Artigo 6.° n.° 6 alínea a) do CIVA, a contrário                                   |
| `M41` | IVA - autoliquidação / Artigo 8.° n.° 3 do RITI                                                          |
| `M42` | IVA - autoliquidação / Decreto-Lei n.° 21/2007, de 29 de janeiro                                         |
| `M43` | IVA - autoliquidação / Decreto-Lei n.° 362/99, de 16 de setembro                                         |
| `M99` | Não sujeito ou não tributado                                                                             |

For example, you could define an invoice line exempt of tax as follows:

```js
{
  "$schema": "https://gobl.org/draft-0/bill/invoice",
  // [...]
  "lines": [
    {
      // [...]
      "item": {
        "name": "Some service exempt of tax",
        "price": "25.00"
      },
      "tax": [
        {
            "cat": "VAT",
            "rate": "exempt",
            "ext": {
              "pt-exemption-code": "M19"
            }
        }
      ]
    }
  ]
}
```
