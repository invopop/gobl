## Added

- `convert`: `Input.Requested` holds the format keys an import was limited to, so a detector can claim a requested format it would not pick on its own, such as one whose documents only differ from another's by their container.

## Changed

- `convert`: `Converter.Import` takes the `*ImportOptions`, as `Export` takes the `*ExportOptions`, so converters receive the requested formats.
