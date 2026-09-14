## Added

- `bill`: `RoundToCurrency` recalculates an invoice with the `currency` rounding rule so that every amount fits the currency's precision, with any change in the amount payable carried in the totals' `rounding`.

## Changed

- `bill`: `RemoveIncludedTaxes` only switches to the `precise` rounding rule when the rule was inherited from the tax regime. A rule set on the document itself is now respected.

## Fixed

- `bill`: `RemoveIncludedTaxes` maintains the amount payable instead of the total with tax, adds to any `rounding` total already present instead of replacing it, and no longer panics on documents that have not been calculated.
- `bill`: fixed discount and charge amounts inside lines are rounded to the currency's precision when the `currency` rounding rule applies, instead of keeping the decimals they were given.
- `bill`: the delivery `type` is marked as calculated in the JSON schema, which was previously rendered as `"enum": "advice"`.
