## Changed

- **breaking**: `catalogues/untdid`: the GOBL unit key to UN/ECE unit code
  mapping moved here from the `eu-en16931-v2017` addon as `UnitCode` and
  `UnitKey`, so that formats outside the EN 16931 family built on the same
  codes, such as the Mexican CFDI's `ClaveUnidad`, can share one table. The
  addon's `UnitToUNTDID` and `UnitFromUNTDID` are gone; call the catalogue
  directly. The addon still owns normalization of the `untdid-unit` extension
  for EN 16931 documents.
