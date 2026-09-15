## Fixed

- `c14n`: a properly encoded U+FFFD is no longer rejected as invalid UTF-8.
  `utf8.DecodeRuneInString` reports `RuneError` both for malformed bytes and
  for the replacement character itself, so the check now also requires a rune
  size of 1. Documents carrying text a sender had already mangled failed to
  digest with a misleading `json: unsupported value`.
