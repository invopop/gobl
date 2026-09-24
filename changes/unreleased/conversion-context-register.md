## Added

- `convert`: a central register of converters between GOBL and external
  formats. Converter packages call `convert.Register` from `init` to add their
  contexts. `convert.Import` detects the context of incoming data and converts
  it into an envelope, `convert.Export` converts an envelope into the first of
  a list of preferred contexts that accepts it, and `convert.Contexts`,
  `convert.ContextsFor`, and `convert.Conversions` list what is available.
