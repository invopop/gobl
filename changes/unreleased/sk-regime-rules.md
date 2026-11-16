## Added

- `sk`: rule `GOBL-SK-BILL-INVOICE-01`, the supplier of a Slovak invoice needs a tax ID code (DIČ or IČ DPH), § 74(1)(a) of the VAT Act 222/2004, and `GOBL-SK-BILL-INVOICE-02`, that tax ID must be Slovak.
- `sk`: identity type `ICO` for the IČO, with separators removed and validated as 8 digits with its mod 11 check digit, following Peppol participant identifier scheme 0158.

## Changed

- `sk`: **breaking**: the tax ID code must start with 1 to 7, the types of person of Peppol participant identifier scheme 0245. A code starting with 8 or 9 now fails `GOBL-SK-TAX-IDENTITY-01`.
- `sk`: the description explains the DIČ, the IČ DPH and the IČO, and exports and intra-EU supplies are described as exempt with the right to deduct.
