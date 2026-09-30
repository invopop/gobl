## Fixed

- `bill`: an `issue_time` of midnight (`00:00:00`) is now respected during calculation. Previously it was indistinguishable from an empty string and both the time and the `issue_date` were replaced with the current date and time. Only an empty string, or `cal.EmptyTime()` from Go, now requests the current date and time.

## Added

- `cal`: `Time` gains `IsEmpty` and the `EmptyTime` constructor to distinguish a deliberately undefined time from midnight. Empty times marshal back to `""`, the JSON schema pattern now admits an empty string, and JSON `null` keeps the plain zero value.
