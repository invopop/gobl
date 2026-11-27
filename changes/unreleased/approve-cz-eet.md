## Added

- `addons`: approved `cz-eet-v1` as an external addon implemented by [`github.com/invopop/gobl.cz.eet`](https://github.com/invopop/gobl.cz.eet), so it is now a valid `$addons` value.
- `bill`: payment lines accept `ext`, so an addon can mark what a single line of a payment is for, as invoice lines already allow. Payment lines have their own schema, `bill/payment-line`, for extensions to name in their usage.
