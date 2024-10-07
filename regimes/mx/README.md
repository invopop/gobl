# 🇲🇽 GOBL Mexico Tax Regime

Mexico uses the CFDI (Comprobante Fiscal Digital por Internet) format for e-invoicing.

Find example MX GOBL files in the [`examples`](./examples) (YAML uncalculated documents) and [`examples/out`](./examples/out) (JSON calculated envelopes) subdirectories.

## Public Documentation

- [Invoice Format (Anexo 20)](http://omawww.sat.gob.mx/tramitesyservicios/Paginas/anexo_20.htm)
- [CFDI 4.0 Filling Guide](http://omawww.sat.gob.mx/tramitesyservicios/Paginas/documentos/Anexo_20_Guia_de_llenado_CFDI.pdf)
- [Global CFDI 4.0 Filling Guide](http://omawww.sat.gob.mx/tramitesyservicios/Paginas/documentos/GuiallenadoCFDIglobal311221.pdf)

## Mexico-specific Requirements

Mexican invoices require specific codes as defined in the CFDI specification. Suppliers need to know these in advance or request them from customers during purchase.

Here's how to include these codes in your GOBL documents:

### Issue Place (`LugarExpedicion`)

Specify the postal code where the invoice was issued using the `mx-cfdi-issue-place` extension under the `tax` field of the invoice.

#### Example

```js
{
  "$schema": "https://gobl.org/draft-0/bill/invoice",
  // [...]
  "tax": {
    "ext": {
      "mx-cfdi-issue-place": "26015"
    }
  },
  // [...]
}
```

### `RegimenFiscal` - Fiscal Regime

Every Supplier and Customer in a Mexican invoice must be associated with a fiscal regime code. You'll need to ensure this field's value is requested from customers when they require an invoice.

In GOBL the `mx-cfdi-fiscal-regime` extension key is used alongside the value expected by the SAT.

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
      "code": "EKU9003173C9"
    },
    "ext": {
      "mx-cfdi-fiscal-regime": "601"
    }
  }
  // [...]
}
```

### `DomicilioFiscalReceptor` - Receipt's Tax Address

In CFDI, `DomicilioFiscalReceptor` is a mandatory field that specifies the postal code of the recepient's tax address. In a GOBL Invoice, you can provide this value setting the `mx-cfdi-post-code` extension of the invoice's customer.

#### Example

```js
{
  "$schema": "https://gobl.org/draft-0/bill/invoice",

  // [...]

  "customer": {
    "name": "UNIVERSIDAD ROBOTICA ESPAÑOLA",
    "tax_id": {
      "country": "MX",
      "code": "URE180429TM6"
    },
    "ext": {
      "mx-cfdi-fiscal-regime": "601",
      "mx-cfdi-use": "G01",
      "mx-cfdi-post-code": "65000"
    }
  }

  // [...]
}
```

### `UsoCFDI` - CFDI Use

The CFDI’s `UsoCFDI` field specifies how the invoice's recipient will use the invoice to deduce taxes for the expenditure made. In a GOBL Invoice, include the `mx-cfdi-use` extension in the customer.

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
      "code": "URE180429TM6"
    },
    "ext": {
      "mx-cfdi-fiscal-regime": "601",
      "mx-cfdi-use": "G01",
      "mx-cfdi-post-code": "65000"
    }
  }

  // [...]
}
```

### `MetodoPago` – Payment Method

The CFDI’s `MetodoPago` field specifies whether the invoice has been fully paid at the moment of issuing the invoice (`PUE` - Pago en una sola exhibición) or whether it will be paid in one or several instalments after that (`PPD` – Pago en parcialidades o diferido).

In GOBL, the presence or absence of payment advances covering the invoice’s total payable amount will determine whether the `MetodoPago` will be set to `PUE` or `PPD`.

Please note that if you don't include any advances in your GOBL invoice, it will be assumed that the payment of the invoice is outstanding (`MetodoPago = PPD`). This implies that the invoice supplier will have to issue “CFDI de Complemento de Pago” (remittance advice) documents at a later stage when receiving the invoice payments.

#### Examples

The following GOBL will map to the `PUE` (Pago en una sola exhibición) value of the `MetodoPago` field:

```js
{
  "$schema": "https://gobl.org/draft-0/bill/invoice",

  // [...]

  "payment": {
    "advances": [
      {
        "key": "credit-transfer",
        "desc": "Full credit card payment",
        "amount": "232.00"
      }
    ]
  },

  "totals": {
    "payable": "232.00",
    "advance": "232.00",
    "due": "0.00"
    // [...]
  }
}
```

The following GOBL will map to the `PPD` (Pago en parcialidades o diferido) value of the `MetodoPago` field:

```js
{
  "$schema": "https://gobl.org/draft-0/bill/invoice",

  // [...]

  "payment": {
    "advances": [
      {
        "key": "credit-transfer",
        "desc": "Partial credit card payment",
        "amount": "100.00"
      }
    ]
  },

  "totals": {
    "payable": "232.00",
    "advance": "100.00",
    "due": "132.00"
    // [...]
  }
}
```

The following GOBL will map to the `PPD` (Pago en parcialidades o diferido) value of the `MetodoPago` field:

```js
{
  "$schema": "https://gobl.org/draft-0/bill/invoice",

  // No "payment" key [...]

  "totals": {
    "payable": "232.00",
    // [...]
  }
}
```

### `FormaPago` - Payment Means

The CFDI’s `FormaPago` field specifies an invoice's means of payment.

If the invoice hasn't been fully paid at the time of issuing the invoice (`MetodoPago = PPD`, see [the previous section](#metodopago-–-payment-method)) the value of `FormaPago` will always be set to `99` (Por definir).

Otherwise (`MetodoPago = PUE`), the `FormaPago` value will be mapped from the key of the largest payment advance in the GOBL invoice. The following table lists all the supported values and how GOBL will map them:

| Code | Name                                | GOBL Payment Advance Key      |
| ---- | ----------------------------------- | ----------------------------- |
| 01   | Efectivo                            | `cash`                        |
| 02   | Cheque nominativo                   | `cheque`                      |
| 03   | Transferencia electrónica de fondos | `credit-transfer`             |
| 04   | Tarjeta de crédito                  | `card`                        |
| 05   | Monedero electrónico                | `online+wallet`               |
| 06   | Dinero electrónico                  | `online`                      |
| 08   | Vales de despensa                   | `other+grocery-vouchers`      |
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
    "advances": [
      {
        "key": "online+wallet",
        "desc": "Prepago con monedero electrónico",
        "amount": "100.00"
      }
    ]
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

The CFDI’s `ClaveProdServ` field specifies the type of an invoice's line item. GOBL uses the line item extension key `mx-cfdi-prod-serv` to map the code directly to the `ClaveProdServ` field.

The catalogue of available Product or Service codes that form part of the CFDI standard is immense with some 50.000 entries to choose from. Due the huge number of codes GOBL will not validate these fields, you'll have to check with local accountants to check which code should be used for your products or services.

#### Example

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
        "ext": {
          "mx-cfdi-prod-serv": "10101602"
        }
      },
    }
  ]
}
```

### Generic RFCs (B2C sales)

In MX, CFDI invoices to end consumers need to use a generic RFC code (`XAXX010101000`) in the document's `Receptor` and some specific values for the supplier extensions. In GOBL, you don't need to enter those values explicitly. Just set the `simplified` tax tag in your invoice and exclude the customer, the GOBL to CFDI conversor will take care of the rest:

```js
{
  "$schema": "https://gobl.org/draft-0/bill/invoice",

  // [...]

  "tax": {
    "tags": [
      "simplified"
    ]
  }

  // No "customer" entry

}
```
Sellers can issue one simplified invoice for every B2C sale or aggregate all daily transactions into a single invoice to _the general public_ ("al público en general"). The rules on how to issue global B2C invoices are available in the [public documentation](http://omawww.sat.gob.mx/tramitesyservicios/Paginas/documentos/GuiallenadoCFDIglobal311221.pdf).

Similarly, CFDI invoices to foreign customers are required to use the foreign generic RFC (`XEXX010101000`) along with some specific values for the extensions. The country and the local tax code of the supplier are included in separated CFDI fields. And again, in GOBL, you just need to provide the actual country and tax code of the supplier (without extensions), and the GOBL to CFDI conversor will take care of producing the correct document:

```js
{
  "$schema": "https://gobl.org/draft-0/bill/invoice",

  // [...]

  "customer": {
    "name": "Almuerzos Colombianos SAS",
    "tax_id": {
      "country": "CO",
      "code": "9014514805"
    }
  }

  // [...]
}
```

### Complements

Complements let you include additional complementary information to your GOBL documents. The following sections describe the complements made available by the MX regime, their purpose and how you can use them.

#### Fuel Account Balance

The _Fuel Account Balance_ complement carries the data to produce a CFDI’s [“Complemento de Estado de Cuenta de Combustibles para Monederos Electrónicos” (version 1.2 revision B)](https://www.sat.gob.mx/consulta/21885/genera-tus-facturas-electronicas-con-el-complemento-para-el-estado-de-cuenta-de-combustibles-para-monederos-electronicos) providing detailed information about fuel purchases made with electronic wallets.

In Mexico, e-wallet suppliers use this complement to report this information in the invoices to their customers.

Learn more about this complement here:
* [Schema Documentation](https://docs.gobl.org/draft-0/regimes/mx/fuel_account_balance)
* [Example GOBL document](./examples/out/fuel-account-balance.json)

#### Food Vouchers

The _Food Vouchers_ complement carries the data to produce a CFDI's [“Complemento de Vales de Despensa” (version 1.0)](https://docs.gobl.org/draft-0/regimes/mx/food_vouchers) providing detailed information about food vouchers issued by an e-wallet supplier to its customer's employees.

In Mexico, e-wallet suppliers use this complement to report this information in the invoices to their customers.

Learn more about this complement here:
* [Schema Documentation](https://docs.gobl.org/draft-0/regimes/mx/food_vouchers)
* [Example GOBL document](./examples/out/food-vouchers.json)
