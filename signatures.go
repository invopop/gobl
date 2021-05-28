package gobl

import "github.com/invopop/gobl/dsig"

// Signatures keeps together a list of signatures that we're used to sign the document
// head contents.
type Signatures []*dsig.Signature
