## Added

- `es`: the `IPSI` tax category (Ceuta and Melilla) now accepts the `standard`,
  `zero`, `reverse-charge`, `exempt`, `export`, and `outside-scope` keys, so
  exempt and non-subject IPSI lines can be expressed without a percent.
  `intra-community` is deliberately not available.
- `es`: reverse-charge invoices with IPSI now receive the same legal note as VAT
  and IGIC.
- `es-verifactu-v1`: IPSI tax combos are normalized like VAT and IGIC, setting
  `es-verifactu-op-class` (`S1`, `S2`, `N1`, `N2`) or `es-verifactu-exempt`
  (`E1`–`E4`, `E6`) from the tax key. A taxed IPSI combo with just a percent
  now builds to `S1` without the integrator adding the extension.
- `es-verifactu-v1`: IPSI combos receive an `es-verifactu-regime` code, as
  required by the AEAT validation rules since revision 1.1.6 (rejected from
  1 January 2027 if missing). `E1` exemptions default to `19` (exempt domestic
  operations) and everything else to `01`; `08`, `11`, `18`, and `20` must be
  set explicitly. New rule `GOBL-ES-VERIFACTU-TAX-COMBO-05` rejects regime
  codes the AEAT does not accept for IPSI.

## Changed

- `es-verifactu-v1`: the "regime is required", "op-class and exempt are
  mutually exclusive" and "taxed operations require op-class" validations now
  also apply to IPSI combos. The "E2/E3 not allowed with regime 01" rule
  remains limited to VAT and IGIC, matching the AEAT rules.
- `es-verifactu-v1`: regime code names for `08`, `11`, `18`, `19`, and `20`
  now include their IPSI meanings.
