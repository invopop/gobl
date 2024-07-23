If you need to regenerate testing data, be sure to use the `private.jwk` provided here.

Use a pre-built `gobl` binary and a command like the following to regenerate individual files:

```bash
./gobl sign -k ./internal/testdata/private.jwk -i ./internal/testdata/success.json
```
