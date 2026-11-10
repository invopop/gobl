## Added

- `eu-en16931-v2017`: **breaking**: document-level charges must have taxes (BR-37), as discounts already must (BR-32). Charges without taxes now fail with `GOBL-EU-EN16931-BILL-CHARGE-02`.

## Fixed

- `eu-en16931-v2017`: the charge and line discount reason rules cite BR-38 and BR-42, not BR-36 (charge amount) and BR-41 (line allowance amount). Fault codes are unchanged.
