# 🇲🇽 GOBL Mexico Tax Regime

Mexico uses the CFDI (Comprobante Fiscal Digital por Internet) format for their e-invoicing system.

Example MX GOBL files can be found in the [`examples`](./examples) (YAML uncalculated documents) and [`examples/out`](./examples/out) (JSON calculated envelopes) subdirectories.

## Public Documentation

- [Formato de factura (Anexo 20)](http://omawww.sat.gob.mx/tramitesyservicios/Paginas/anexo_20.htm)

## Local Codes

Mexican invoices as defined in the CFDI specification must include a set of specific codes that will either need to be known in advance by the supplier or requested from the customer during their purchase process.

The following sections highlight these codes and how they can be defined inside your GOBL documents.

### `RegimenFiscal` - Fiscal Regime

Every Supplier and Customer in a Mexican invoice must be associated with a fiscal regime code. You'll need to ensure this field's value is requested from customers when they require an invoice.

In GOBL the `mx-cfdi-fiscal-regime` identity key is used alongside the value expected by the SAT.

#### Example

The following example will associate the supplier with the `601` fiscal regime code:

```js
{
  "$schema": "https://gobl.org/draft-0/bill/invoice",
  // [...]
  "supplier": {
    "name": "ESCUELA KEMPER URGATE",
    "tax_id": {
      "country": "MX",
      "zone": "26015",
      "code": "EKU9003173C9"
    },
    "identities": [
      {
        "key": "mx-cfdi-fiscal-regime",
        "code": "601"
      }
    ]
  }
  // [...]
}
```

### `UsoCFDI` - CFDI Use

The CFDI’s `UsoCFDI` field specifies how the invoice's recipient will use the invoice to deduce taxes for the expenditure made. In a GOBL Invoice, include the `mx-cfdi-use` identity in the customer.

This field will be validated for presence and will be checked against the list of codes defined as part of the CFDI specification.

#### Example

The following GOBL maps to the `G03` (Gastos en general) value of the `UsoCFDI` field:

```js
{
  "$schema": "https://gobl.org/draft-0/bill/invoice",

  // [...]

  "customer": {
    "name": "UNIVERSIDAD ROBOTICA ESPAÑOLA",
    "tax_id": {
      "country": "MX",
      "zone": "65000",
      "code": "URE180429TM6"
    },
    "identities": [
      {
        "key": "mx-cfdi-fiscal-regime",
        "code": "601"
      },
      {
        "key": "mx-cfdi-use",
        "code": "G01"
      }
    ]
  }

  // [...]
}
```

### `FormaPago` - Payment Means

The CFDI’s `FormaPago` field specifies an invoice's means of payment. The following table lists all the supported values and how GOBL will map them from the invoice's payment instructions key:

| Code | Name                                | GOBL Payment Instructions Key |
| ---- | ----------------------------------- | ----------------------------- |
| 01   | Efectivo                            | `cash`                        |
| 02   | Cheque nominativo                   | `cheque`                      |
| 03   | Transferencia electrónica de fondos | `credit-transfer`             |
| 04   | Tarjeta de crédito                  | `card`                        |
| 05   | Monedero electrónico                | `online+wallet`               |
| 06   | Dinero electrónico                  | `online`                      |
| 08   | Vales de despensa                   | `other+grocery-vouchers  `    |
| 12   | Dación en pago                      | `other+in-kind`               |
| 13   | Pago por subrogación                | `other+subrogation`           |
| 14   | Pago por consignación               | `other+consignment`           |
| 15   | Condonación                         | `other+debt-relief`           |
| 17   | Compensación                        | `netting`                     |
| 23   | Novación                            | `other+novation`              |
| 24   | Confusión                           | `other+merger`                |
| 25   | Remisión de deuda                   | `other+remission`             |
| 26   | Prescripción o caducidad            | `other+expiration`            |
| 27   | A satisfacción del acreedor         | `other+satisfy-creditor`      |
| 28   | Tarjeta de débito                   | `card+debit`                  |
| 29   | Tarjeta de servicios                | `card+services`               |
| 30   | Aplicación de anticipos             | `other+advance`               |
| 31   | Intermediario pagos                 | `other+intermediary`          |
| 99   | Por definir                         | `other`                       |

#### Example

The following GOBL maps to the `05` (Monedero electrónico) value of the `FormaPago` field:

```js
{
  "$schema": "https://gobl.org/draft-0/bill/invoice",

  // [...]

  "payment": {
    "instructions": {
      "key": "online+wallet"
    }
  }
}
```

### `ClaveUnidad` - Unit Code

The CFDI’s `ClaveUnidad` field specifies the unit in which the quantity of an invoice's line is given. These are UNECE codes that GOBL will map directly from the invoice's line item unit. See the [source code](../../org/unit.go) for the full list of supported units with their associated UNECE codes.

#### Example

The following GOBL maps to the `KGM` (Kilogram) value of the `ClaveUnidad` field:

```js
{
  "$schema": "https://gobl.org/draft-0/bill/invoice",

  // [...]

  "lines": [
    {
      // [...]

      "item": {
        "name": "Jasmine rice",
        "unit": "kg",
        "price": "1.27"
      },
    }
  ]
}
```

### `ClaveProdServ` - Product or Service Code

The CFDI’s `ClaveProdServ` field specifies the type of an invoice's line item. GOBL uses the line item identity key `mx-cfdi-prod-serv` to map the identity code directly to the `ClaveProdServ` field.

The catalogue of available Product or Service codes that form part of the CFDI standard is immense with some 50.000 entries to choose from. At this time, GOBL will not validate these fields, you'll have to check with local accountants to check which code should be used for your products or services.

### Example

The following GOBL maps to the `10101602` ("live ducks") value to the `ClaveProdServ` field:

```js
{
  "$schema": "https://gobl.org/draft-0/bill/invoice",

  // [...]

  "lines": [
    {
      // [...]

      "item": {
        "name": "Selección de patos vivos",
        "identities": [
          {
            "key": "mx-cfdi-prod-serv",
            "code": "10101602"
          }
        ]
      },
    }
  ]
}
```
