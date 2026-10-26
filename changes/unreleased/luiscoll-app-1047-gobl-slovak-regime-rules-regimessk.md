## Added

- `sk`: invoice rules for the content that VAT Act 222/2004 § 74(1) requires.
  The supplier needs a tax ID code and a postal address with street, locality
  and code. The customer needs the same postal address, except on simplified
  invoices. An invoice not in EUR needs an exchange rate to EUR to state the VAT
  in euros.
- `sk`: the `sk-ico` identity key for the IČO, with a format check of eight
  digits. Spaces and separators are removed when it is normalized.
- `sk`: tags for the special schemes (`travel-agency`, `second-hand-goods`,
  `art`, `antiques`, `cash-basis`) and invoice notes with the wording of the
  law for reverse charge, self-billing, exempt supplies and these schemes.

## Changed

- `sk`: **breaking** an invoice with a supplier in Slovakia now fails
  validation when the supplier or the customer has no postal address.
- `sk`: the regime description no longer calls exports and intra-EU supplies
  "zero-rated", and explains the DIČ, the IČ DPH and the IČO.
