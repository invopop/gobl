## Fixed

- `org`: `Inbox` normalization no longer re-interprets the `code` as a URL or email
  when a `scheme` is set or the key is `peppol`. Dotted participant IDs such as French
  `0225` routing codes were previously moved into the `url` field on every save.
