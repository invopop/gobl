package head

import (
	"fmt"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/dsig"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/uuid"
	"github.com/invopop/jsonschema"
	"github.com/invopop/validation"
	"github.com/invopop/validation/is"
)

// Link category keys defined for use inside link categories.
const (
	LinkCategoryKeyFormat       cbc.Key = "format"
	LinkCategoryKeyPortal       cbc.Key = "portal"
	LinkCategoryKeyRequest      cbc.Key = "request"
	LinkCategoryKeyResponse     cbc.Key = "response"
	LinkCategoryKeyAgreement    cbc.Key = "agreement"
	LinkCategoryKeyVerification cbc.Key = "verification"
	LinkCategoryKeyAttachment   cbc.Key = "attachment"
)

// LinkCategoryDefs provides the definitions for link categories.
var LinkCategoryDefs = []*cbc.Definition{
	{
		Key:  LinkCategoryKeyFormat,
		Name: i18n.NewString("Format"),
		Desc: i18n.NewString(here.Doc(`
			Alternative formats of the same document, such as PDF, HTML, or XML.
		`)),
	},
	{
		Key:  LinkCategoryKeyPortal,
		Name: i18n.NewString("Portal"),
		Desc: i18n.NewString(here.Doc(`
			Websites that provide access to alternative versions of this document or request
			modifications. May also allow access to other or previous business documents.
		`)),
	},
	{
		Key:  LinkCategoryKeyRequest,
		Name: i18n.NewString("Request"),
		Desc: i18n.NewString(here.Doc(`
			Documents related to requests submitted to other systems.
		`)),
	},
	{
		Key:  LinkCategoryKeyResponse,
		Name: i18n.NewString("Response"),
		Desc: i18n.NewString(here.Doc(`
			Response documents sent from third party systems often in reply to requests.
		`)),
	},
	{
		Key:  LinkCategoryKeyAgreement,
		Name: i18n.NewString("Agreement"),
		Desc: i18n.NewString(here.Doc(`
			Contracts or agreements related to this document.
		`)),
	},
	{
		Key:  LinkCategoryKeyVerification,
		Name: i18n.NewString("Verification"),
		Desc: i18n.NewString(here.Doc(`
			Evidence or verification documents that can be used to validate the authenticity
			of this document.
		`)),
	},
	{
		Key:  LinkCategoryKeyAttachment,
		Name: i18n.NewString("Attachment"),
		Desc: i18n.NewString(here.Doc(`
			General attachments for related information such as spread sheets, letters,
			presentations, specifications, etc.
		`)),
	},
}

// Link defines a link between this document and another resource. Much like stamps,
// links must be defined with a specific key, but do allow for additional data that
// can help with presentation. It is important that a link once generated cannot be
// updated, so this is not suitable for dynamic or potentially insecure URLs.
//
// Links have a specific advantage over stamps in that they are also allowed while
// the envelope has not yet been signed.
type Link struct {
	uuid.Identify
	// Category helps classify the link according to a fixed list. This is optional
	// but highly recommended as it helps receivers better understand the purpose
	// of the link and potentially how its should be presented.
	Category cbc.Key `json:"category,omitempty" jsonschema:"title=Category"`
	// Key is a unique identifier for the link within the header and category
	// if provided.
	Key cbc.Key `json:"key"`
	// Code used to identify the contents of the link.
	Code cbc.Code `json:"code,omitempty" jsonschema:"title=Code"`
	// Title of the resource to use when presenting to users.
	Title string `json:"title,omitempty" jsonschema:"title=Title"`
	// Description of the resource to use when presenting to users.
	Description string `json:"description,omitempty" jsonschema:"title=Description"`
	// Expected MIME type of the link's content when the content is a file. Can only
	// be one of the allowed types defined by EN 16931-1:2017 plus XML itself.
	MIME string `json:"mime,omitempty" jsonschema:"title=MIME Type,format=mime"`
	// Digest is used to verify the integrity of the destination document
	// when downloaded from the URL.
	Digest *dsig.Digest `json:"digest,omitempty" jsonschema:"title=Digest"`
	// URL of the resource.
	URL string `json:"url" jsonschema:"title=URL,format=uri"`
}

// Validate checks that the link contains the basic information we need to function.
func (l *Link) Validate() error {
	return validation.ValidateStruct(l,
		validation.Field(&l.UUID),
		validation.Field(&l.Category,
			validation.In(cbc.DefinitionKeys(LinkCategoryDefs)...),
		),
		validation.Field(&l.Key,
			validation.Required,
		),
		validation.Field(&l.Code),
		validation.Field(&l.Title),
		validation.Field(&l.Description),
		validation.Field(&l.MIME,
			validation.In(
				// Allow EN16931-1:2017 defined MIME types
				"application/pdf",
				"image/jpeg",
				"image/png",
				"text/csv",
				"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
				"application/vnd.oasis.opendocument.spreadsheet",
				// Alternative types
				"text/html",
				"application/xml",
				"text/xml",
				"application/json",
			),
		),
		validation.Field(&l.Digest,
			validation.When(
				l.MIME == "",
				validation.Nil.Error("must be nil when MIME type is not provided"),
			),
		),
		validation.Field(&l.URL, validation.Required, is.URL),
	)
}

// LinkByCategoryAndKey finds the link with the given category and key from the provided list.
// Category may be empty.
func LinkByCategoryAndKey(list []*Link, category cbc.Key, key cbc.Key) *Link {
	for _, l := range list {
		if l.Category == category && l.Key == key {
			return l
		}
	}
	return nil
}

// AppendLink will add the link to the provided list and return the new updated
// list. If the link already exists, it will be updated.
func AppendLink(list []*Link, l *Link) []*Link {
	if l == nil {
		return list
	}
	for _, v := range list {
		if v.Category == l.Category && v.Key == l.Key {
			*v = *l // copy in place
			return list
		}
	}
	return append(list, l)
}

// DetectDuplicateLinks checks if the list of links contains duplicate
// category and key pairs.
var DetectDuplicateLinks = validation.By(detectDuplicateLinks)

func detectDuplicateLinks(list any) error {
	values, ok := list.([]*Link)
	if !ok || len(values) == 0 {
		return nil
	}
	set := []*Link{}
	// loop through and check order of Since value
	for _, v := range values {
		if l := LinkByCategoryAndKey(set, v.Category, v.Key); l != nil {
			if v.Category == "" {
				return fmt.Errorf("duplicate key '%v'", v.Key)
			}
			return fmt.Errorf("duplicate category '%v' and key '%v'", v.Category, v.Key)
		}
		set = append(set, v)
	}
	return nil
}

// JSONSchemaExtend implements the jsonschema.Extender interface for Link and
// adds extra category details to the schema.
func (Link) JSONSchemaExtend(js *jsonschema.Schema) {
	prop, ok := js.Properties.Get("category")
	if ok {
		prop.OneOf = make([]*jsonschema.Schema, len(LinkCategoryDefs))
		for i, def := range LinkCategoryDefs {
			prop.OneOf[i] = &jsonschema.Schema{
				Const:       def.Key,
				Title:       def.Name.String(),
				Description: def.Desc.String(),
			}
		}
	}
}
