// Package i18n provides internationalization models.
package i18n

import "github.com/invopop/gobl/schema"

func init() {
	schema.Register(schema.GOBL.Add("i18n"), String{})
}
