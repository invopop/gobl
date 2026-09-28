## Added

- `es`: the `IPSI` tax category accepts the `standard` and `exempt` keys, so exempt
  IPSI lines can be expressed without a percent.
- `es-verifactu-v1`: IPSI combos are normalized automatically: taxed lines get
  `es-verifactu-op-class` `S1` and regime `01`; `exempt` lines get
  `es-verifactu-exempt` `E1` and regime `19`, as required by the AEAT validation
  rules for `Impuesto` 02 since revision 1.1.6.
