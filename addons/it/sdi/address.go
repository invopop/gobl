package sdi

import (
	"fmt"
	"strconv"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
)

func normalizeAddress(addr *org.Address) {
	if addr == nil || addr.Country != l10n.IT.ISO() {
		return
	}

	// ensure the Code is always 5 digits with 0 padding
	if len(addr.Code) > 0 && len(addr.Code) < 5 {
		// convert to number
		code, err := strconv.Atoi(addr.Code.String())
		if err != nil {
			// not a number, ignore
			return
		}
		addr.Code = cbc.Code(fmt.Sprintf("%05d", code))
	}
}
