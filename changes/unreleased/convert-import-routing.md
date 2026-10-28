## Added

- `convert`: `WithRouting(from, to)` import option, recorded on the envelope's `Head.From` and `Head.To`. Converters receive the import options, so they can set the routing before calculating the envelope.

## Changed

- `convert`: `Converter.Import` takes the `*ImportOptions`, as `Export` takes the `*ExportOptions`.
