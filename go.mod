module github.com/invopop/gobl

go 1.14

require (
	cloud.google.com/go v0.99.0
	github.com/Masterminds/semver/v3 v3.1.1
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/ghodss/yaml v1.0.0
	github.com/go-ozzo/ozzo-validation/v4 v4.3.0
	github.com/google/uuid v1.1.2
	github.com/invopop/jsonschema v0.4.1-0.20220509222051-9cef489f4cb7
	github.com/kr/pretty v0.2.0 // indirect
	github.com/square/go-jose/v3 v3.0.0-20200630053402-0a67ce9b0693
	github.com/stretchr/testify v1.7.0
	golang.org/x/crypto v0.0.0-20210817164053-32db794688a5 // indirect
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)

// replace github.com/invopop/jsonschema => ../jsonschema
