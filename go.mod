module github.com/invopop/gobl

go 1.14

require (
	cloud.google.com/go v0.65.0
	github.com/alecthomas/jsonschema v0.0.0-20211222145519-a257ab6fe816
	github.com/go-ozzo/ozzo-validation v3.6.0+incompatible
	github.com/go-ozzo/ozzo-validation/v4 v4.3.0
	github.com/google/uuid v1.1.2
	github.com/square/go-jose/v3 v3.0.0-20200630053402-0a67ce9b0693
	github.com/stretchr/testify v1.6.1
	golang.org/x/crypto v0.0.0-20200820211705-5c72a883971a // indirect
)

replace github.com/alecthomas/jsonschema => github.com/invopop/jsonschema v0.0.0-20211222145519-a257ab6fe816

// replace github.com/alecthomas/jsonschema => ../jsonschema
