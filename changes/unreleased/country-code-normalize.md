## Fixed

- `l10n`: country codes are trimmed and uppercased during normalization, so `de` becomes `DE`.
- `l10n`: `ISOCountryCode` and `TaxCountryCode` are validated by the rules engine (`GOBL-L10N-ISOCOUNTRYCODE-01`, `GOBL-L10N-TAXCOUNTRYCODE-01`).
