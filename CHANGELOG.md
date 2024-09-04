# Change Log

All notable changes to GOBL will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/) and this project adheres to [Semantic Versioning](http://semver.org/). See also the [GOBL versions](https://docs.gobl.org/overview/versions) documentation site for more details.

## [Unreleased]

This one is big...

Significant set of changes around Scenario handling. Scenarios defined by tax regimes can now set tax extensions at the document level automatically. The objective here is to move away from external projects using scenario summaries directly, and instead use the absolute values set in the document.

For example, the document format and type in Italy are now set inside the extensions and can be overriden if needed manually. This will be especially important when receiving and converting invoices into GOBL from external formats; its much easier to set specific values than trying to determine the appropriate tags.

Also included is support for defining the country in tax combos, making it possible for taxes from a customers country to applied directly if needed. Typical use case would be for selling digital goods into or between EU states for B2C customers.

Invoices in GOBL can now also finally produced for any country in the world, even if not explicitly defined inside the tax regimes.

### Changed

- `tax`: moved `NormalizeIdentity` method from the regimes common package so that it can be applied to all tax IDs, regardless of if they have a regime defined or not.
- `pt`: VAT rate key is now optional if `pt-saft-tax-rate` is provided.
- `gr`: simplified validation to use tax categories.
- `it`: always add `it-sdi-fiscal-regime` to Invoice suppliers.
- `it`: renamed extension `it-sdi-retained-tax` to `it-sdi-retained`, now with validation on retained taxes.
- `it`: renamed extension `it-sdi-natura` to `it-sdi-exempt`.
- `bill.Invoice`: deprecated the `ScenarioSummary` method, as tax regimes themselves should be using extensions to apply all the correct data to a document up front.
- `mx`: scenarios will now copy the document and relation types to the tax extensions.

### Added

- `tax`: `Combo` now supports a `country` field.
- `tax.Category`: added `Validation` method support for custom validation of a tax combo for a specific tax category.
- `tax.Scenario`: added "extensions" to be able to automatically update document level extensions based on the scenario detected.
- `it`: added `ExtKeySDIDocumentType` as an extension that will be automatically included according to the scenario.
- `it`: now adding `ExtKeySDIFormat` value to document instead of just referencing from scenarios.
- `cbc.Note`: now provides `SameAs` method that will compare key attributes, but not the text payload. This is now used in Schema Summaries.
- `bill.Line`: added `RequireLineTaxCategory` validation helper method.

### Removed

- `tax.Category`: removed `RateRequired` flag, regimes should instead should help users determine valid extensions (eg. PT and GR).

### Fixed

- `tax.Scenario`: potential issue around matching notes.
- `tax.Set`: improved validation embedded error handling.

## [v0.114.0]

### Changed

- `org.Name`: either given **or** surname are required, as opposed to both at the same time.

## [v0.113.0]

### Added

- `head`: validation rule to check for the presence of stamps
- GR: support for credit notes

## [v0.112.0]

Significant set of small changes related to renaming of the `l10n.CountryCode` type. The main reason for this is an attempt to reduce confusion between regular ISO country selection, and the specific country codes used for tax purposes. Normally they coincide, but exception cases like for Greece, whose ISO code is `GR` but use `EL` for tax purposes, or `XI` for companies in Northern Ireland, mean that there needs to be a clear selection.

### Changed

- CO: improved regime's documentation
- `l10n`: split "CountryCode" into "ISOCountryCode" and "TaxCountryCode", for the two explicit use-cases.
- `l10n`: renamed `CountryDefinitions` variable to `Countries()` method.

### Added

- Code coverage report (still a lot to improve there!)
- GR: support for simplified invoices
- `l10n`: ISO and Tax lists of country definitions available, e.g. `l10n.Countries().ISO()`
- `tax`: support for alternative country codes

### Upgraded

- [invopop/validation](https://github.com/invopop/validation) - upgrade to latest version with nil pointer fix.

### Fixed

- GR: fixed certain tax combos not getting calculated by the regime

## [v0.112.0] - 2024-07-26

### Added

- `tax`: Scenarios now handle extension key and value for filtering.
- PT: exemption text handling moved to scenarios.

## [v0.111.1] - 2024-07-25

### Added

- `org.Address`: recommended fields added

### Changed

- `org.Address`: `locality` no longer required.

## [v0.111.0] - 2024-07-24

### Added

- Including `recommended` array in more JSON Schema objects.
- `bill.Invoice`: validation and changes around acceptance of simplified invoices with customer name. A customer without a tax ID now implies that a name is also not required.
- `uuid`: Compact Base64 encoding and decoding of UUIDs for compact URLs.
- `head`: New `Link` model for associating Envelopes with static URLs.
- `head.Header`: Link array in addition to stamps.
- PT: support for debit notes
- PT: validations for debit and credit notes

### Changed

- CO: renamed credit and debit extension names to fit in UIs.
- `org.Party`: `name` is now optional, but recommended.

### Fixed

- `org.People`: `name` is now correctly validated.

## [v0.110.2] - 2024-07-23

### Added

- `org.Person`: now has `label` field.

### Changed

- `org`: Refining availability of `label` field.

## [v0.110.0] - 2024-07-23

Multiple version upgrade after merging the [gobl.cli](https://github.com/invopop/gobl.cli) project directly here instead.

### Added

- CLI: move the command line interface and wasm binary support directly into GOBL.

## [v0.83.0] - 2024-07-23

### Added

- CO: support debit notes with additional validations for required reason extensions.

### Changed

- CO: renaming `co-dian-correction` code to `co-dian-credit-code` while also adding `co-dian-debit-code` to extensions.
- CO: support debit notes
- CO: updated validation for simplified invoices
- GR: renamed greece country code to `EL` to reflect local naming in tax code, package still named `gr` for ease of use.
- l10n: extension countries like EL, XI, EU for special tax cases
- l10n: country definition extension flag to be able to filter ISO codes

### Fixed

- IT: Company fiscal code can be the same as the VAT code.

## [v0.82.0] - 2024-07-19

### Added

- `bill.Invoice`: experimental `ConvertInto` method to convert the invoice's amounts from one currency into another.
- DE: support for "de-tax-number" identity which can be used instead of regular tax ID code inside Germany.
- DE: "simplified" tax tag removes requirement for supplier tax identification.

## [v0.81.0] - 2024-07-17

### Added

- `tax.Regime`: added new "Identity Keys" definition.
- IT: `it-sdi-format` extension added with the two main document formats in Italy: `FPA12` and `FPR12` (default for B2B/C if none assigned).

### Changed

- `tax.Identity`: deprecated the `type` field, and directly removed the `uuid` and `meta` fields which no longer make sense here.
- `tax.Regime`: standardised naming around key definitions to always include `_keys` as suffix.
- IT: moved fiscal code (codice fiscale) from the `org.Party` Tax ID to the Identities array with the specific key `it-fiscal-code`. This implies that invoices can now be issued with **both** a VAT code (partita IVA) and a fiscal code (codice fiscale).
- IT: data will be normalized automatically to move the fiscal code from the tax ID to the identities array.
- IT: removed explicit support for Tax ID type field.
- ES: moved Tax ID `type` usage to the `identities` array.
- CO: moved Tax ID type definitions to `identities` array.

## [v0.80.1] - 2024-07-11

### Added

- GR: Invoice and address validations
- GR: Payment means key definitions

### Changed

- GB: removed requirement for suppliers to have a tax ID code (country is still required!)

## [v0.80.0] - 2024-06-27

### Added

- `num.Amount` - `RescaleDown` method, that helps reduce accuracy from a number if higher.
- `num.Amount` - `RescaleRange` method, ensures that the exponent is within a specific range.
- Greece tax regime
- `tax.Combo` - regime specific calculations now supported.

## [v0.79.3] - 2024-06-18

### Added

- `org.Registration`: added `other` field.

### Fixed

- Field descriptions of `org.Website`

## [v0.79.2] - 2024-06-18

### Added

- PT: `at-app-id` stamp for the application ID used to register a document.

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
