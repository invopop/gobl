// Package common provides re-usable regime related structures and data.
package common

import (
	"github.com/invopop/gobl/cbc"
)

// Common inbox keys
const (
	InboxKeyPEPPOL cbc.Key = "peppol-id"
)

// Common Identity Type Codes that are not country specific.
const (
	IdentityTypeDUNS cbc.Code = "DUNS" // Dun & Bradstreet - Data Universal Numbering System
)
