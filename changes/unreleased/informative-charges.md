## Added

- `bill`: `Charge` accepts an `informative` flag. An informative charge is
  calculated, validated, and kept on the document, but is left out of the
  charge total, the taxable base, and every amount that follows, and cannot
  carry taxes. It models costs that must be declared on an invoice while being
  borne by the supplier, such as the Italian *Marca da Bollo* when the issuer
  absorbs the €2 instead of passing it on to the customer.
