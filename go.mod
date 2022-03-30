module github.com/invopop/gobl

go 1.14

require (
	cloud.google.com/go v0.99.0
	github.com/Masterminds/semver/v3 v3.1.1
	github.com/ghodss/yaml v1.0.0
	github.com/go-ozzo/ozzo-validation v3.6.0+incompatible
	github.com/go-ozzo/ozzo-validation/v4 v4.3.0
	github.com/go-playground/validator/v10 v10.10.1
	github.com/google/uuid v1.1.2
	github.com/invopop/jsonschema v0.2.0
	github.com/square/go-jose/v3 v3.0.0-20200630053402-0a67ce9b0693
	github.com/stretchr/testify v1.7.0
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

// replace github.com/invopop/jsonschema => ../jsonschema
