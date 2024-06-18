# Change Log

All notable changes to GOBL will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/) and this project adheres to [Semantic Versioning](http://semver.org/). See also the [GOBL versions](https://docs.gobl.org/overview/versions) documentation site for more details.

## [v0.79.1] - 2024-06-18

### Added

- `num.AmountFromFloat64` method helps avoid rounding issues when using floats.
- PT: `at-hash` stamp to store the AT's hash of an invoice

### Changed

- MX: Fuel account balance totals calculated maintaining item price precision.

## [v0.79.0] - 2024-06-05

### Added

- ISO 3166-1 alpha-3 codes (and a function to access them) added to the country definitions (`l10n.CountryDef`)
- MX: `mx.TaxIdentityCodeGeneric` constant added with the generic RFC for final consumers

## Changed

- MX: customer extensions no longer required for foreign customers

## [v0.78.1] - 2024-05-30

### Fixed

- CH: fixing tax ID checksum validation that ends in `0`.

## [v0.78.0] - 2024-05-23

### Added

- Default values added to correction options schema based on values from previous invoice.

### Changed

- In `pay.Online`, renamed `name` property to `label`, and `addr` to `URL`. Also added `key` property. Auto-migration included.

### Fixed

- `bill.Invoice`: fixed issue with invalid currency codes that don't have a definition, will always resort to Tax Regime's currency.

## [v0.77.0] - 2024-05-16

Fixing important bugs with some tax regimes and tax identity code validation, along with a new 'series' property for correction options.

### Added

- `series` property added to bill Correction Options to allow a series for a credit note to be added.
- `label` property added to `org.Party`.

### Changed

- Invoice Discounts and Charges will no longer update the `base` property according to the document's sum.
- Exempt rate in `tax.Combo` no longer required when percent is empty.
- When correcting an invoice, if no new series is provided, the previous document's series will be maintained.

### Fixed

- Precision handling of calculated invoice discounts and charges.
- Multiple tax regime were not validating the presence of supplier identity code.

## [v0.76.0] - 2024-05-13

Finally, invoice multi-currency support! It's been a very long time coming, but we've finalized the details on how to handle currency conversion in invoices and potentially other documents.

### Added

- Invoice Line Item alternative pricing added to be able to define custom prices in different currencies: `line.AltPrices`
- Automatic conversion of invoice line item prices into invoice currency based on exchange rates defined in invoice if no alternative prices provided.
- If an invoice has a currency that is different from that of the tax regime, a validation rule ensures that an exchange rate is defined.
- `currency.Amount` - new model that combines a currency with an amount.
- BE: added Belgium regime.

### Changed

- _BREAKING_: refactor of `currency.ExchangeRate` to clearly define `from` and `to` currencies (this was never supported, so we're not expecting anything to actually break).
- Removed all regime specific currency validation, this is now performed by the invoice and depends on the available exchange rates.
- MX: invoice line totals validated to be **zero** or more, instead of positive.

### Fixed

- Removing code requirement from Tax ID validation in all regimes so that when issuing a document to another country, the customers tax ID code will be validated if present, but will **not** be required. Any local rules for the issuing country for foreign IDs will continue to be applied.

## [v0.75.1] - 2024-05-07

### Change

- MX: allow line taxes to be empty.

## [v0.75.0] - 2024-05-06

### Added

- Bill Invoice Tax objects now support tax extensions.
- MX Stamps for signatures from CFDI and SAT.
- MX: extension for Place of Issue code: `mx-cfdi-issue-place` that replaces previous post code option in the supplier. Automatic normalization added.
- `head` package now has `GetStamp` method to find a stamp by its provider from an array.
- `num.Percentage` has `Base()` method to access base amount.
- MX: FuelAccountBalance complement now supports `percent` as an alternative to `rate`.
- ES: added extra TicketBAI exemption reasons
- Envelope `Replicate()` and supporting methods to be able to clone/replicate an envelope or document without any potentially conflicting data.

### Changed

- `reverse-charge` tag will no longer have impact on tax calculations, each tax combo per line should define if taxes are exempt or not.
- Renaming `mx.StampProviderSATUUID` constant to just `mx.StampSATUUID`.
- MX: FuelAccountBalance complement renamed tax `code` to `cat` (Category) with explicit usage of regular tax codes to be more aligned with other usage of tax categories.

### Fixed

- UUID `IsNotZero` will not raise error for empty UUIDs.

## [v0.74.1] - 2024-05-23

UUID Unmarshal fix.

### Fixed

- UUID: parsing empty strings from JSON no longer causes error.

## [v0.74.0] - 2024-05-23

Refining UUID library and moving to using version 7 as the default in GOBL.

### Added

- Additional UUID validation rules: `Valid`, `IsV6`, `Isv7`, `HasTimestamp`, and `Timeless`.

### Changed

- Using Version 7 UUIDs as default in GOBL. This version enables ordering by UUID and uses random extra data instead of a node ID.

### Fixed

- Parsing empty UUIDs now returns an empty UUID instead of an error.

## [v0.73.0] - 2024-04-22

Refactoring UUID support.

**IMPORTANT:** When running `Calculate()`, a uuid will now be assigned automatically to the document embedded in an Envelope if not already set. This is important to ensure that links between documents can always be maintained, no matter the source.

### Added

- Schema Object: `Calculate()` will now inject UUIDs.
- Schema Object: `UUID()` method will provide the UUID of the underlying document.
- `schema.Identifiable` interface to be able to read and set UUIDs on documents.
- `uuid.Identify` that makes it easier to embed UUIDs into documents with helper methods.

### Changed

- UUID: refactored to use underlying string type instead of external package. This makes it easier to manage empty values, and avoids usage of pointers.
- Removed all pointers to UUIDs and many cases replaced with `uuid.Identify` embedded structure.

### Fixed

- none

## [v0.72.0] - 2024-04-18

Refactoring region handling for Portugal VAT and now supporting `-` in `cbc.Code`.

### Added

- Regimes: Extensions can now be used to match tax rates.
- Tax: Extensions helper methods: `Merge` and `Contains`.

### Changed

- `cbc.Code`: Now supports `-` symbol alongside `.` as a separator. Mixed feelings on this as we wanted to avoid normalization complications, but it became clear with the PT changes that a bit more flexibility here is useful. (Side note: the original intent of `cbc.Code` was to avoid dashes in tax IDs, but these are now normalized automatically.)
- PT: moving from tax tags `azores` and `madeira` to `pt-region` extension provided in taxes combo for each line.
- PT: auto-migrate invoice supplier tax ID zone to appropriate line tax combo extension.

### Fixed

- none

## [v0.71.0] - 2024-04-08

New number formatting support! Expect some possible breaking SDK changes with the `num` packages. No significant schema changes.

### Added

- This CHANGELOG.md file (finally!)
- Swiss (CH) tax regime.
- Austrian (AT) tax regime.
- `num` package now provides advanced number formatting.
- `currency` provides "definitions" loaded from JSON with support for formatting.
- Polish (PL) correction and preceding validation.
- Polish (PL) header stamps for QR code.

### Changed

- `num` package refactored so that `num.Percentage` is independent from `num.Amount`.

### Fixed

- Minor fixes around tax regime definitions.
- [invopop/yaml](https://github.com/invopop/yaml) upgraded.

## [v0.70.1] - 2024-03-25

- Last version before CHANGELOG.md.
