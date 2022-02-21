module github.com/invopop/gobl

go 1.14

require (
	cloud.google.com/go v0.99.0
	github.com/Masterminds/semver/v3 v3.1.1
	github.com/ghodss/yaml v1.0.0
	github.com/go-ozzo/ozzo-validation/v4 v4.3.0
	github.com/google/uuid v1.1.2
	github.com/imdario/mergo v0.3.12
	github.com/invopop/jsonschema v0.1.0
	github.com/labstack/echo/v4 v4.6.3
	github.com/spf13/cobra v1.3.0
	github.com/square/go-jose/v3 v3.0.0-20200630053402-0a67ce9b0693
	github.com/stretchr/testify v1.7.0
	gitlab.com/flimzy/testy v0.9.1
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
)

// replace github.com/alecthomas/jsonschema => github.com/invopop/jsonschema v0.0.0-20211230180634-99ed368317c4

// replace github.com/invopop/jsonschema => ../jsonschema
