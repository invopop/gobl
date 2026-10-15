## Added

- `convert`: `Import` provides the files embedded in the data, such as a PDF
  version of an invoice, as `Attachments` in its `ImportResult`.

## Changed

- `convert`: **breaking** `Import(data, opts...)` returns an `ImportResult`
  with the detected `Format`, the `Envelope` and the `Attachments`. Detection
  is limited to some formats with the `WithFormats` option.
- `convert`: **breaking** `Converter.Import` returns an `ImportResult`.
- `convert`: **breaking** `Output` is renamed `ExportResult`.
