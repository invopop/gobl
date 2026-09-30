## Added

- `tax`: `Filter` selects rate totals by category and, optionally, key and extensions. `Total.FilteredAmount` sums the matching rate amounts.
- `pay`: `Record.Taxes` takes tax filters as an alternative to `percent`, so the amount is the sum of the matching rate totals. Setting both fails validation, as does using taxes on `bill.Payment` methods.
- `pay`: `Record.Waiver` gives the reason an amount was waived and not collected from the customer, such as a tax refunded at the point of sale, as an alternative to the means `key`. GOBL doesn't define the values, and `bill.Payment` methods don't accept it.

## Changed

- `bill`: due date percentages are calculated from the amount due after advances, falling back to the payable amount. When nothing remains due, their amounts are removed.
- `bill`: advances with tax filters take their amount from the invoice's rounded tax totals.
- **breaking**: `pay`: `Record.Ref` and `DirectDebit.Ref` are `cbc.Code` instead of `string`, matching `Instructions.Ref`. JSON documents are unchanged, but references are limited to 128 characters and Go code assigning a `string` needs a conversion.
