## Changed

- **breaking**: `catalogues/untdid`: the GOBL unit key to UN/ECE unit code
  mapping lives here rather than in the `eu-en16931-v2017` addon, as `UnitCode`
  and `UnitKey`, so that formats outside the EN 16931 family built on the same
  codes, such as the Mexican CFDI's `ClaveUnidad`, can share one table. The
  addon's `UnitToUNTDID` and `UnitFromUNTDID` are gone; call the catalogue
  directly.
- **breaking**: `eu-en16931-v2017`: the `untdid-unit` extension determines an
  item's unit whenever it is set, leaving the unit empty for a code GOBL has no
  key for. Normalization neither adds the extension nor removes it, so a
  document that carries a standard code keeps it alongside the unit, and one
  that does not is left alone. Item attributes are normalized the same way.
- **breaking**: `eu-en16931-v2017`: BR-23 accepts a unit code that is
  determinable, from either the unit key or the extension, rather than one
  stored in the extension. The generic `one` unit only stands in when the
  document gives neither, so a code with no GOBL equivalent no longer produces
  a unit that contradicts it.
