## Added

- `convert`: a central register of converters between GOBL and external
  formats. Converter packages call `convert.Register` from `init` to add their
  formats. `convert.Import` detects the format of incoming data and converts
  it into an envelope, `convert.Export` converts an envelope into the first of
  a list of preferred formats that accepts it, and `convert.Formats`,
  `convert.FormatsFor`, and `convert.Conversions` list what is available.
