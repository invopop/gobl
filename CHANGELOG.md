# Change Log

All notable changes to GOBL will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/) and this project adheres to [Semantic Versioning](http://semver.org/). See also the [GOBL versions](https://docs.gobl.org/overview/versions) documentation site for more details.

## [Unreleased] - yyyy-mm-dd

Pending release notes...

### Added

- `cbc.Code`: Now supports `-` symbol alongside `.` as a separator. We have mixed feelings on this as we wanted to avoid normalization complications, but it became clear with the PT changes that a bit more flexibility here is useful. (Side note: the original intent here was to avoid dashes in tax IDs, but these are now normalized automatically.)
- Tax: Extensions can now be used to match tax rate values.

### Changed

- PT: moving from tax tags `azores` and `madeir` to `pt-actud-region` extension provided in the supplier.

### Fixed

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
