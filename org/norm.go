package org

import (
	"github.com/invopop/gobl/norm"
)

// init registers the intrinsic normalizers for the org package's types. Each
// function performs only the type's own normalization; the norm engine handles
// recursion into nested values, the global cleaning of codes and extensions,
// and the application of regime/addon normalizers.
func init() {
	norm.Register(
		norm.For(normalizeParty),
		norm.For(normalizePerson),
		norm.For(normalizeItem),
		norm.For(normalizeAddress),
		norm.For(normalizeDocumentRef),
		norm.For(normalizeEmail),
		norm.For(normalizeEndpoint),
		norm.For(normalizeIdentity),
		norm.For(normalizeInbox),
		norm.For(normalizeName),
		norm.For(normalizeNote),
		norm.For(normalizeRegistration),
		norm.For(normalizeTelephone),
		norm.For(normalizeWebsite),
		norm.For(normalizeAttachment),
	)
}
