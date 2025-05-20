# Adding a new regime

All new features should come with tests. For example, `tax_identities.go` must have a `tax_identities_test.go` file.

- Duplicate the existing [`template`](./template/) directory and rename it to the 2-letter country code (regime code).
- Update the necessary files and code as needed.
  - Rename `template.go` to `<regime_code>.go`
    - `New` function should instantiate a [`*tax.RegimeDef`](tax/regime_def.go).
    - `Normalize` and `Validate` functions should take care of each possible element that is specific to the regime.
  - `tax_categories.go`
  - `tax_identity.go`

- Optionally, add the following files:
  - `scenarios.go`
  - `corrections.go`
  - `org_parties.go`
  - `org_identities.go`
- Add the new regime to the `regimes.go` file.
- Add the new regime to the `regimes_test.go` file.

Using an `inv.Validator` to validate invoices is deprecated. Instead, validate each element separately in the `Validate` function.