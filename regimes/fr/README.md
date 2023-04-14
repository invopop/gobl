# GOBL France

## Tax IDs

France has three main company IDs which are all very closely related and may be included on Invoice documents:

- VAT code, an 11 digit number that starts with a two digit checksum followed by 9 numbers that form the SIREN. In standardised form, this is always presented with the country code (`FR`), e.g. `44732829320` or `FR44732829320`.
- SIREN is the national company ID register and consists of 9 digits with a final digit representing a checksum, e.g. `732829320`.
- SIRET is the SIREN with an additional 5 digits representing a department inside the company or tax agency, e.g. `73282932000015`

During the normalization process of Tax Identities, GOBL will automatically convert a SIREN into a VAT code by adding the checksum to beginning.

The SIREN and SIRET numbers may also be defined inside the `org.Party`'s `identities` property, for example a supplier inside an invoice might look like:

```json
{
  "tax_id": {
    "country": "FR",
    "code": "44732829320"
  },
  "name": "Dummy FR Inc.",
  "addresses": [
    {
      "num": "1",
      "street": "Rue Sundacsakn",
      "locality": "Saint-Germain-En-Laye",
      "code": "75050",
      "country": "FR"
    }
  ],
  "emails": [
    {
      "addr": "email@dummycom.fr"
    }
  ],
  "identities": [
    {
      "type": "SIREN",
      "code": "732829320"
    },
    {
      "type": "SIRET",
      "code": "73282932000015"
    }
  ]
}
```
