## Fixed

- `l10n`: country codes are uppercased and trimmed during normalization, so a
  tax identity or address country sent as `de` becomes `DE` and the matching
  regime's normalizers and rules apply. Previously the lowercase value was
  signed as is and reached converters verbatim, where formats that require an
  uppercase ISO code rejected the document.
- `l10n`: `ISOCountryCode` and `TaxCountryCode` are validated again through
  the rules engine (`GOBL-L10N-ISOCOUNTRYCODE-01` and
  `GOBL-L10N-TAXCOUNTRYCODE-01`). The checks were lost in the migration from
  the `validation` package, which called each type's `Validate` method
  implicitly.
