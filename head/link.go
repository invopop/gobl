package head

import (
	"fmt"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/validation"
	"github.com/invopop/validation/is"
)

// Link defines a link between this document and another resource. Much like stamps,
// links must be defined with a specific key, but do allow for additional data that
// can help with presentation. It is important that a link once generated cannot be
// updated, so this is not suitable for dynamic or potentially insecure.
//
// Links have a specific advantage over stamps in that they are also allowed while
// the envelope is still a draft.
type Link struct {
	// Key is a unique identifier for the link.
	Key cbc.Key `json:"key"`
	// Title of the resource to use when presenting to users.
	Title string `json:"title,omitempty" jsonschema:"title=Title"`
	// Description of the resource to use when presenting to users.
	Description string `json:"description,omitempty" jsonschema:"title=Description"`
	// Expected MIME type of the link's content.
	MIME string `json:"mime,omitempty" jsonschema:"title=MIME Type,format=mime"`
	// URL of the resource.
	URL string `json:"url" jsonschema:"title=URL,format=uri"`
}

// Validate checks that the link contains the basic information we need to function.
func (l *Link) Validate() error {
	return validation.ValidateStruct(l,
		validation.Field(&l.Key, validation.Required),
		validation.Field(&l.Title),       // not required
		validation.Field(&l.Description), // not required
		validation.Field(&l.MIME),
		validation.Field(&l.URL, validation.Required, is.URL),
	)
}

// LinkByKey finds the link with the given key from the provided list.
func LinkByKey(list []*Link, k cbc.Key) *Link {
	for _, l := range list {
		if l.Key == k {
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
		if v.Key == l.Key {
			*v = *l // copy in place
			return list
		}
	}
	return append(list, l)
}

// DetectDuplicateLinks checks if the list of links contains duplicate
// keys.
var DetectDuplicateLinks = validation.By(detectDuplicateLinks)

func detectDuplicateLinks(list any) error {
	values, ok := list.([]*Link)
	if !ok || len(values) == 0 {
		return nil
	}
	set := []*Link{}
	// loop through and check order of Since value
	for _, v := range values {
		if l := LinkByKey(set, v.Key); l != nil {
			return fmt.Errorf("duplicate key '%v'", v.Key)
		}
		set = append(set, v)
	}
	return nil
}
