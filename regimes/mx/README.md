# 🇲🇽 GOBL Mexico Tax Regime

Mexico uses the CFDI (Comprobante Fiscal Digital por Internet) format for their e-invoicing system.

## Public Documentation

* [Formato de factura (Anexo 20)](http://omawww.sat.gob.mx/tramitesyservicios/Paginas/anexo_20.htm)


## Local Codes

### `UsoCFDI`

The CFDI’s `UsoCFDI` field specifies how the invoice's recipient will use the invoice to deduce taxes for the expenditure made. The following table lists all the supported values and how GOBL will map them from the invoice's tax tags:

| Code | Name | GOBL Tax Tag |
| --- | --- | --- |
| G01 | Adquisición de mercancías | `use+goods-acquisition` |
| G02 | Devoluciones, descuentos o bonificaciones | `use+returns` |
| G03 | Gastos en general | `use+general-expenses` |
| I01 | Construcciones | `use+construction` |
| I02 | Mobiliario y equipo de oficina por inversiones | `use+office-equipment` |
| I03 | Equipo de transporte | `use+transport-equipment` |
| I04 | Equipo de computo y accesorios | `use+computer-equipment` |
| I05 | Dados, troqueles, moldes, matrices y herramental | `use+manufacturing-tooling` |
| I06 | Comunicaciones telefónicas | `use+telephone-comms` |
| I07 | Comunicaciones satelitales | `use+satellite-comms` |
| I08 | Otra maquinaria y equipo | `use+other-machinery` |
| D01 | Honorarios médicos, dentales y gastos hospitalarios | `use+medical-expenses` |
| D02 | Gastos médicos por incapacidad o discapacidad | `use+medical-expenses+disability` |
| D03 | Gastos funerales | `use+funeral-expenses` |
| D04 | Donativos | `use+donation` |
| D05 | Intereses reales efectivamente pagados por créditos hipotecarios (casa habitación) | `use+mortgage-interest` |
| D06 | Aportaciones voluntarias al SAR | `use+sar-contribution` |
| D07 | Primas por seguros de gastos médicos | `use+medical-insurance` |
| D08 | Gastos de transportación escolar obligatoria | `use+school-transportation` |
| D09 | Depósitos en cuentas para el ahorro, primas que tengan como base planes de pensiones | `use+savings-deposit` |
| D10 | Pagos por servicios educativos (colegiaturas) | `use+school-fees` |
| S01 | Sin efectos fiscales | `use+no-tax-effects` |
| CP01 | Pagos | `use+suplementary-payment` |
| CN01 | Nómina | `use+payroll` |

#### Example

The following GOBL maps to the `G03` (Gastos en general) value of the `UsoCFDI` field:

```js
{
  "$schema": "https://gobl.org/draft-0/bill/invoice",

  // [...]

  "tax": {
    "tags": [
      "use+general-expenses"
    ]
  }
}
```

### `FormaPago`

The CFDI’s `FormaPago` field specifies an invoice's means of payment. The following table lists all the supported values and how GOBL will map them from the invoice's payment instructions key:

| Code | Name | GOBL Payment Instructions Key |
| --- | --- | --- |
| 01 | Efectivo | `cash` |
| 02 | Cheque nominativo | `cheque` |
| 03 | Transferencia electrónica de fondos | `credit-transfer` |
| 04 | Tarjeta de crédito | `card` |
| 05 | Monedero electrónico | `online+wallet` |
| 06 | Dinero electrónico | `online` |
| 08 | Vales de despensa | `other+grocery-vouchers  ` |
| 12 | Dación en pago | `other+in-kind` |
| 13 | Pago por subrogación | `other+subrogation` |
| 14 | Pago por consignación | `other+consignment` |
| 15 | Condonación | `other+debt-relief` |
| 17 | Compensación | `netting` |
| 23 | Novación | `other+novation` |
| 24 | Confusión | `other+merger` |
| 25 | Remisión de deuda | `other+remission` |
| 26 | Prescripción o caducidad | `other+expiration` |
| 27 | A satisfacción del acreedor | `other+satisfy-creditor` |
| 28 | Tarjeta de débito | `card+debit` |
| 29 | Tarjeta de servicios | `card+services` |
| 30 | Aplicación de anticipos | `other+advance` |
| 31 | Intermediario pagos | `other+intermediary` |
| 99 | Por definir | `other` |

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

### `ClaveUnidad`

The CFDI’s `ClaveUnidad` field specifies the unit in which the quantity of an invoice's line is given. These are UNECE codes that GOBL will map directly from the invoice's line item unit. See the [source code](../../org/unit.go) for the full list of supported units with its UNECE codes.

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

### `ClaveProdServ`

The CFDI’s `ClaveProdServ` field specifies the type of an invoice's line item. GOBL uses the line item identity type `SAT` to map the identity code directly (no transformation) to the `ClaveProdServ` field.

### Example

The following GOBL maps to the `10101602` (Patos vivos) value of the `ClaveProdServ` field:

```js
{
  "$schema": "https://gobl.org/draft-0/bill/invoice",

  // [...]

  "lines": [
    {
      // [...]

      "item": {
        "name": "Live ducks",
        "identities": [
          {
            "type": "SAT",
            "code": "10101602"
          }
        ]
      },
    }
  ]
}
```
