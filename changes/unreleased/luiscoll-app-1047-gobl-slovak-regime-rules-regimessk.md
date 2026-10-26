## Added

- `sk`: invoice rule for the supplier's tax ID code (DIČ or IČ DPH), which VAT
  Act 222/2004 § 74(1) requires.
- `sk`: the `sk-ico` identity key for the IČO, with a format check of eight
  digits. Spaces, dots, dashes and slashes are removed when it is normalized, and the
  EN 16931 addon maps it to the ISO 6523 scheme 0158.
- `sk`: tags for the special schemes (`travel-agency`, `second-hand-goods`,
  `art`, `antiques`, `cash-basis`) and invoice notes with the wording of the
  law for reverse charge, self-billing and these schemes.

## Changed

- `sk`: the regime description no longer calls exports and intra-EU supplies
  "zero-rated", and explains the DIČ, the IČ DPH and the IČO.
