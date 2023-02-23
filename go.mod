module github.com/invopop/gobl

go 1.14

require (
	cloud.google.com/go v0.99.0
	github.com/Masterminds/semver/v3 v3.1.1
	github.com/google/uuid v1.1.2
	github.com/invopop/jsonschema v0.7.0
	github.com/invopop/validation v0.3.0
	github.com/invopop/yaml v0.1.0
	github.com/kr/pretty v0.2.0 // indirect
	github.com/square/go-jose/v3 v3.0.0-20200630053402-0a67ce9b0693
	github.com/stretchr/testify v1.8.1
	golang.org/x/crypto v0.6.0 // indirect
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
)

// replace github.com/invopop/jsonschema => ../jsonschema
// replace github.com/invopop/validation => ../validation
