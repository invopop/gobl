# 🇧🇷 GOBL Brazil Tax Regime

Brazil uses _Notas Fiscais Eletrônicas_ (like NFSe, NFe or NFCe) for reporting tax information to municipality, state or federal authorities.

Find example BR GOBL files in the [`examples`](../../examples/br) (uncalculated documents) and [`examples/out`](../../examples/br/out) (calculated envelopes) subdirectories.

## Public Documentation

* [NFS-e Technical Documentation at ABRASAF](https://abrasf.org.br/biblioteca/arquivos-publicos/nfs-e)
* [IBGE Municipality Codes](https://www.ibge.gov.br/explica/codigos-dos-municipios.php)

## Brazil-specific Requirements

### Addresses

Brazilian addresses have 3 subdivisions relevant for tax purposes: _bairro_ (neighbourhood), _municipio_ (municipality) and _estado_ (state). Specify them in GOBL addresses like this:

| GOBL Field | Maps to                  | Example  |
| ---------- | ------------------------ | -------- |
| `locality` | Bairro / Neighbourhood   | Centro   |
| `region`   | Município / Municipality | Salvador |
| `state`    | Estado / State           | BR       |

For example:

```js
"supplier": {
  //...
  "addresses": [
    {
      "num": "75",
      "street": "Avenida Sete de Setembro",
      "street_extra": "Bloco C",
      "locality": "Centro", // Bairro
      "region": "Salvador", // Municipio
      "state": "BA", // State
      "code": "40060-000",
      "country": "BR"
    }
  ],
  //...
```

### Service Notes

Services notes (NFSe, Notas fiscais de servicio eletrônicas) let service providers document and report the taxes (e.g. ISS) related to the services they provision. Municipal governments regulate them.

Please also see the [NFSe Addon](../../addons/br/nfse) package named `br-nfse-v1`, which you should include in your documents.

#### Municipality Code

Set the party's [IBGE municipality code](https://www.ibge.gov.br/explica/codigos-dos-municipios.php) using the `br-nfse-municipality` extension.

For example:

```js
"supplier": {
  //...
  "ext": {
    "br-igbe-municipality": "2927408"
  },
  //...
```

#### National and municipal registration

Specify the party's municipal and national registration numbers as identities using the `br-nfse-municipal-reg` and `br-nfse-national-reg` keys.

For example:

```js
"supplier": {
  //...
  "identities": [
    {
      "key": "br-nfse-municipal-reg",
      "code": "45678901234567"
    },
    {
      "key": "br-nfse-national-reg",
      "code": "12345012345678"
    }
  ],
  //...
```

#### “Simples Nacional”

Report whether the party opts in for the _Simples Nacional_ simplified regime using the `br-nfse-simples-nacional` extension set to any of these codes:

| Code | Description |
| ---- | ----------- |
| `1`  | Opt-in      |
| `2`  | Opt-out     |

For example:

```js
"supplier": {
  //...
  "ext": {
    "br-nfse-simples-nacional": "1", // Opt-in
  },
  //...
```

#### Special Tax Regime

Specify a special tax regime the party is subject to using the `br-nfse-special-regime` set to any of these codes:

| Code | Description                                 |
| ---- | ------------------------------------------- |
| `1`  | Municipal micro-enterprise                  |
| `2`  | Estimated                                   |
| `3`  | Professional Society                        |
| `4`  | Cooperative                                 |
| `5`  | Single micro-entrepreneur (MEI)             |
| `6`  | Micro-enterprise or Small Business (ME EPP) |

For example:

```js
"supplier": {
  //...
  "ext": {
    "br-nfse-special-regime": "4" // Cooperative
  },
  //...
```

#### Fiscal Incentive

Report whether the party benefits from a fiscal incentive using the `br-nfse-fiscal-incentive` extension set to any of these codes:

| Code | Description                       |
| ---- | --------------------------------- |
| `1`  | Has incentive                     |
| `2`  | Does not have incentive (Default) |

For example:

```js
"supplier": {
  //...
  "ext": {
    "br-nfse-fiscal-incentive": "2" // No tax incentive
  },
  //...
```
