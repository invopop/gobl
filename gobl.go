package gobl

// import all the dependencies to ensure all init() methods are called.
import (
	_ "github.com/invopop/gobl/bill"
	_ "github.com/invopop/gobl/currency"
	_ "github.com/invopop/gobl/dsig"
	_ "github.com/invopop/gobl/i18n"
	_ "github.com/invopop/gobl/note"
	_ "github.com/invopop/gobl/num"
	_ "github.com/invopop/gobl/org"
	_ "github.com/invopop/gobl/regions"
	"github.com/invopop/gobl/schema"
	_ "github.com/invopop/gobl/uuid"
)

func init() {
	schema.Register(schema.GOBL, Envelope{})
}
