module github.com/invopop/gobl

go 1.14

require (
	cloud.google.com/go v0.110.2
	github.com/Masterminds/semver/v3 v3.2.1
	github.com/asaskevich/govalidator v0.0.0-20230301143203-a9d515a09cc2 // indirect
	github.com/google/uuid v1.3.0
	github.com/iancoleman/orderedmap v0.2.0 // indirect
	github.com/invopop/jsonschema v0.7.0
	github.com/invopop/validation v0.3.0
	github.com/invopop/yaml v0.1.0
	github.com/square/go-jose/v3 v3.0.0-20200630053402-0a67ce9b0693
	github.com/stretchr/testify v1.8.1
	golang.org/x/crypto v0.9.0 // indirect
)

// replace github.com/invopop/jsonschema => ../jsonschema
// replace github.com/invopop/validation => ../validation
