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

SAF-T's `InvoiceType` (Tipo de documento) specifies the type of a sales invoice. In GOBL, this type can be set using the `pt-saft-invoice-type` extension in the tax section. GOBL will set the extension for you based on the type and the tax tags you set in your invoice. The table below shows how this mapping is done:

| Code | Name                                                                     | GOBL Type     | GOBL Tax Tag      |
| ---- | ------------------------------------------------------------------------ | ------------- | ----------------- |
| FT   | Fatura, emitida nos termos do artigo 36.o do Código do IVA              | `standard`    |                   |
| FS   | Fatura simplificada, emitida nos termos do artigo 40.o do Código do IVA | `standard`    | `simplified`      |
| FR   | Fatura-recibo                                                            | `standard`    | `invoice-receipt` |
| ND   | Nota de débito                                                          | `credit-note` |                   |
| NC   | Nota de crédito                                                         | `debit-note`  |                   |

For example:

```js
{
  "$schema": "https://gobl.org/draft-0/bill/invoice",
  "$tags": [
    "invoice-receipt"
  ],
  // [...]
  "type": "standard",
  // [...]
  "tax": {
    "ext": {
      "pt-saft-invoice-type": "FR"
    }
  },
  // [...]
```

### `TaxCountryRegion` (País ou região do imposto)

SAF-T's `TaxCountryRegion` (País ou região do imposto) specifies the region of taxation (Portugal mainland, Açores or Madeira) in a Portuguese invoice. Each region has their own tax rates which can be determined automatically.

To set the specific a region different to Portugal mainland, the `pt-region` extension of each line's VAT tax should be set to one of the following values:

| Code  | Description                                         |
| ----- | --------------------------------------------------- |
| PT    | Mainland Portugal (default)                         |
| PT-AC | Açores                                              |
| PT-MA | Madeira                                             |

For example:

```js
{
  "$schema": "https://gobl.org/draft-0/bill/invoice",
  // [...]
  "lines": [
    {
      // [...]
      "item": {
        "name": "Some service",
        "price": "25.00"
      },
      "tax": [
        {
            "cat": "VAT",
            "rate": "exempt",
            "ext": {
              "pt-region": "PT-AC",
              // [...]
            }
        }
      ]
    }
  ]
}
```

### `ProductType` (Indicador de produto ou serviço)

SAF-T's `ProductType` (Indicador de produto ou serviço) indicates the type of each line item in an invoice. The `pt-saft-product-type` extension used at line item level allows to set the product type to one of the allowed values:

| Code | Description                        |
| ---- | ---------------------------------- |
| P    | Products (default)                 |
| S    | Services                           |
| O    | Other                              |
| E    | Excise Duties                      |
| I    | Taxes, fees and parafiscal charges |

For example:

```js
{
  "$schema": "https://gobl.org/draft-0/bill/invoice",
  // [...]
  "lines": [
    {
      // [...]
      "item": {
        "name": "Some service",
        "price": "25.00",
        "ext": {
          "pt-saft-product-type": "S"
        }
      },
      // [...]
    }
  ]
}
```

### VAT Tax Rates

The SAF-T's `TaxCode` (Código do imposto) is required for invoice items that apply VAT. GOBL provides the `pt-saft-tax-rate` extension to set this code at line tax level. It also determines it automatically this code using the `rate` field (when present). The following table lists the supported tax codes and how GOBL will map them:

| Code | Name            | GOBL Tax Rate                         |
| ---- | --------------- | ------------------------------------- |
| NOR  | Tipo Geral      | `standard`                            |
| INT  | Taxa Intermédia | `intermediate`                        |
| RED  | Taxa Reduzida   | `reduced`                             |
| ISE  | Isenta          | `exempt` + extension code (see below) |

AT's `TaxExemptionCode` (Código do motivo de isenção de imposto) is a code that specifies the reason the VAT tax is exempt in a Portuguese invoice. When the `exempt` tag is used, one of the following must be defined in the `ext` map's `pt-saft-exemption` property:

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
              // [...]
              "pt-saft-tax-rate": "ISE",
              "pt-saft-exemption": "M19"
            }
        }
      ]
    }
  ]
}
```

### Payment receipts (Documento de recibo emitido)

To report payment receipts to the AT, GOBL provides conversion from `bill.Payment` documents. You can find an example of a valid Portugese receipt in the [`example folder`](../../examples/pt).

In a payment, the SAF-T's `PaymentType` (Tipo de documento) field specifies its type. In GOBL, this type can be set using the `pt-saft-payment-type` extension. GOBL will set the extension automatically based on the type and the tax tags you set. The table below shows how this mapping is done:

| Code | Name                                       | GOBL Type | GOBL Tax Tag |
| ---- | ------------------------------------------ | --------- | ------------ |
| RG   | Outro Recibo                               | `payment` |              |
| RC   | Recibo no âmbito do regime de IVA de Caixa | `payment` | `vat-cash`   |

For example:

```js
{
  "$schema": "https://gobl.org/draft-0/bill/receipt",
  // [...]
  "type": "payment",
  // [...]
  "ext": {
    "pt-saft-payment-type": "RG"
  },
  // [...]
}
```

### `PaymentMechanism` (Meios de pagamento)

The SAF-T's `PaymentMechanism` (Meios de pagamento) field specifies the payment means in a sales invoice or payment. GOBL provides the `pt-saft-payment-means` extension to set this value in your `bill.Invoice` advances or in you `bill.Payment` method. GOBL maps certain payment mean keys automatically to this extension:

| Code | Name | GOBL Payment Mean |
| --- | --- | --- |
| CC  | Cartão crédito | `card` |
| CD  | Cartão débito | (*) |
| CH  | Cheque bancário | `cheque` |
| CI  | Letter of credit | (*) |
| CO  | Cheque ou cartão oferta | (*) |
| CS  | Compensação de saldos em conta corrente | `netting` |
| DE  | Dinheiro eletrónico | `online` |
| LC  | Letra comercial | `promissory-note` |
| MB  | Referências de pagamento para Multibanco | (*) |
| NU  | Numerário | `cash` |
| OU  | Outro | `other` |
| PR  | Permuta de bens | (*) |
| TB  | Transferência bancária ou débito direto autorizado | `credit-transfer`, `debit-transfer` or `direct-debit` |
| TR  | Títulos de compensação extrassalarial | (*) |

(*) For codes not mapped from a GOBL Payment Mean, use `other` and explicitly set the extension.

For example, in an GOBL invoice:

```js
{
  "$schema": "https://gobl.org/draft-0/bill/invoice",
  // [...]
  "payment": {
    "advances": [
      {
        "date": "2023-01-30",
        "key": "credit-transfer",
        "description": "Adiantamento",
        "amount": "100.00",
        "ext": {
          "pt-saft-payment-means": "TB"
        }
      }
    ]
  },
  // [...]
}
```

For example, in a GOBL receipt:

```js
{
  "$schema": "https://gobl.org/draft-0/bill/receipt",
  // [...]
  "method": {
    "key": "other",
    "detail": "Compensação extrassalarial",
    "ext": {
      "pt-saft-payment-means": "TR"
    }
  },
  // [...]
}
```

### Work documents (Documentos de conferência)

To report work documents to the AT, GOBL provides conversion from `bill.Invoice` and `bill.Order` documents. You can find examples of valid Portuguese work documents in the [`example folder`](../../examples/pt).

The SAF-T's `WorkType` field (Tipo de documento de conferência) specifies the type of a work document. In GOBL, this type can be set using the `pt-saft-work-type` extension in the tax section of the document. Certain work types are only valid for invoices while others are only valid for orders. For some document types, GOBL will automatically set the appropriate work type extension.

The table below shows the supported work type codes and their compatibility:

| Code | Name                            | GOBL Document | GOBL Type  |
| ---- | ------------------------------- | ------------- | ---------- |
| PF   | Pró-forma                       | Invoice       | `proforma` |
| FC   | Fatura de consignação           | Invoice       |            |
| CC   | Credito de consignação          | Invoice       |            |
| CM   | Consultas de mesa               | Order         |            |
| FO   | Folhas de obra                  | Order         |            |
| NE   | Nota de Encomenda               | Order         | `purchase` |
| OU   | Outros                          | Order         |            |
| OR   | Orçamentos                      | Order         | `quote`    |
| DC   | Documentos de conferência       | Order         |            |
| RP   | Prémio ou recibo de prémio      | Order         |            |
| RE   | Estorno ou recibo de estorno    | Order         |            |
| CS   | Imputação a co-seguradoras      | Order         |            |
| LD   | Imputação a co-seguradora líder | Order         |            |
| RA   | Resseguro aceite                | Order         |            |

Example for a proforma invoice:

```js
{
  "$schema": "https://gobl.org/draft-0/bill/invoice",
  "type": "proforma",
  // [...]
  "tax": {
    "ext": {
      "pt-saft-work-type": "PF"
    }
  },
  // [...]
}
```

Example for a purchase order:

```js
{
  "$schema": "https://gobl.org/draft-0/bill/order",
  "type": "purchase",
  // [...]
  "tax": {
    "ext": {
      "pt-saft-work-type": "NE"
    }
  },
  // [...]
}
```

### Stock Movements (Documentos de movimentação de mercadorias)

To report stock movements to the AT, GOBL provides conversion from `bill.Invoice` documents. You can find an example of a valid Portugese delivery note in the [`example folder`](../../examples/pt).

SAF-T's `MovementType` (Tipo de documento de movimentação de mercadorias) specifies the type of a delivery document. In GOBL, this type can be set using the `pt-saft-movement-type` extension. If not provided explicitly, GOBL will set the extension for you based on the type of your delivery document.

The table below shows how this mapping is done:

| Code | Name                | GOBL Type     |
| ---- | ------------------- | ------------- |
| GR   | Delivery note       | `note`        |
| GT   | Waybill             | `waybill`     |
| GA   | Fixed assets guide  |               |
| GC   | Consignment note    |               |
| GD   | Returns slip        |               |

For example:

```js
{
  "$schema": "https://gobl.org/draft-0/bill/delivery",
  // [...]
  "type": "note",
  // [...]
  "tax": {
    "ext": {
      "pt-saft-movement-type": "GR"
    }
  },
  // [...]
}
```

