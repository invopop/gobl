## Changed

- **breaking**: `catalogues/untdid`: `NormalizeUnit` resolves a unit and its
  extensions so the two never state the same thing twice. A code that GOBL has
  a unit key for becomes that key and the `untdid-unit` extension is dropped;
  the extension now carries only the codes GOBL cannot express, in which case
  the unit is left empty. When both are given the extension wins.
- **breaking**: `eu-en16931-v2017`: BR-23 accepts a unit code that is
  determinable, from either the unit key or the extension, rather than one
  stored in the extension. The generic `one` unit only stands in when the
  document gives neither, so a code with no GOBL equivalent no longer produces
  a unit that contradicts it. Item attributes are normalized the same way,
  without the fallback.

  Documents that stored both a unit and a matching extension lose the redundant
  extension the next time they are calculated.
