# Contributing to GOBL

<img src="https://github.com/invopop/gobl/blob/main/gobl_logo_black_rgb.svg?raw=true" width="181" height="219" alt="GOBL Logo">

Go Business Language. Core library, schemas, and CLI.

Released under the Apache 2.0 [LICENSE](https://github.com/invopop/gobl/blob/main/LICENSE), Copyright 2021-2025 [Invopop S.L.](https://invopop.com).

## Language

The preferred language for contributions is English (American).

## Key Directories

The GOBL repository is organized into key directories, each serving a distinct role:

- **addons**: Optional modules that extend GOBL to support additional document formats, country-specific requirements, or custom validation rules.
- **bill**: Core billing logic, including structures and validation for invoices, deliveries, orders, and payments.
- **c14n** (Canonicalization): Ensures consistent JSON formatting for digital signatures and verification.
- **cal**: Utilities for date, time, and calendar calculations.
- **catalogues**: Standardized lists and code sets (e.g., country codes, tax categories) used throughout GOBL and its extensions.
- **cbc** (Common Basic Components): Shared building blocks such as keys, codes, and reusable definitions.
- **cmd**: Command-line tools for interacting with and processing GOBL documents.
- **currency**: Currency code definitions, exchange rate handling, and monetary value utilities.
- **data**: Autogenerated resources, including embedded files and static data required by GOBL.
- **dsig**: Digital signature support for signing and verifying GOBL documents.
- **examples**: Sample documents and usage examples for testing and reference.
- **head**: Envelope metadata and header structures used in GOBL document packaging.
- **i18n** (Internationalization): Multilingual support for text and labels across GOBL.
- **internal**: Private utilities and helpers not intended for public use.
- **l10n** (Localization): Country and region-specific definitions, including localization data.
- **note**: Structures for notes, comments, or messages within GOBL documents.
- **num**: Numeric types and precise arithmetic for financial calculations.
- **org**: Data structures for organizations, parties, and related business entities.
- **pay**: Payment methods, terms, and processing logic.
- **pkg**: Shared utility packages used across multiple parts of the codebase.
- **regimes** (Tax Regimes): Country-specific tax rules, rates, and validation logic. See [`regimes/README.md`](regimes/README.md) for more details.
- **schema**: JSON schema generation and validation logic.
- **tax**: Core tax structures and logic, used in documents and by regimes or addons.
- **uuid**: Utilities for generating and handling UUIDs.
- **wasm**: WebAssembly integration for running GOBL in browser environments.

Each directory is designed to encapsulate a specific aspect of GOBL, making it easier to locate, understand, and contribute to the relevant parts of the project.

## Development

GOBL uses the `go generate` command to automatically generate JSON schemas, definitions, and some Go code output. After any changes, be sure to run:

```bash
go generate .
```

### Linting and formatting

[golangci-lint](https://golangci-lint.run/) is used to check the code for errors and style issues. The configuration is in the [`.golangci.yaml`](.golangci.yaml) file.

Install locally and run:

```bash
golangci-lint run
golangci-lint run --fix # To autofix lint errors where possible
```

Note: we considered incorporating golangci-lint as a tool directly in the `go.mod`, but due the large amount of dependencies, decided not to do so.

### VS Code setup

[Official docs](https://golangci-lint.run/welcome/integrations/#visual-studio-code)

Add this to your `settings.json` file:

```json
"go.lintTool": "golangci-lint",
"go.lintFlags": [
  "--path-mode=abs",
  "--fast-only"
],
"go.formatTool": "custom",
"go.alternateTools": {
  "customFormatter": "golangci-lint"
},
"go.formatFlags": [
  "fmt",
  "--stdin"
]
```

## Where to Make Changes

GOBL is structured in multiple layers, and choosing the right place for your contribution is important. Consider these four main layers:

1. **GOBL Core**: The foundational packages (e.g., `bill`, `cbc`, `tax`, etc.). Changes here are reviewed carefully, as they affect the entire ecosystem and must align with best practices.
2. **Tax Regimes**: Subdirectories within the `regimes` package, each implementing country-specific tax rules, rates, and validation logic. See [`regimes/README.md`](regimes/README.md) for details.
3. **Addons**: Extensions that handle specialized formats or unique normalization/validation rules. Addons may define extension codes to support reliable conversion of GOBL documents.
4. **External Applications**: Business logic that doesn't fit within GOBL itself, typically because it involves relationships between documents or relies on persistent storage.

If you're unsure where your change belongs, feel free to open an issue or discussion for guidance.
