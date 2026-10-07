## Added

- `regimes/cl`: new tax regime for Chile. Defines the VAT category at the
  current 19% rate (Ley sobre Impuesto a las Ventas y Servicios, D.L. 825),
  validates and normalizes the RUT tax identity using its official mod-11
  check digit, and requires the supplier's and customer's RUT on standard
  invoices (except simplified invoices / "boletas", which may omit the
  customer's). Credit and debit notes are restricted to correcting standard,
  non-simplified invoices only.
