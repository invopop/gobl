package tax

import (
	"fmt"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/validation"
)

// Tags defines the structure to use for allowing an object to be assigned tags
// for use in determining how the content should be handled.
type Tags struct {
	// Tags are used to help identify specific tax scenarios or requirements that will
	// apply changes to the contents of the invoice. Tags by design should always be optional,
	// it should always be possible to build a valid invoice without any tags.
	List []cbc.Key `json:"$tags,omitempty" jsonschema:"title=Tags"`
}

// TagSet defines a set of tags and their descriptions that can be used for a specific
// schema in the context of a Regime or Addon.
//
// Tags between multiple sets may be duplicated. In this case, the definition of the first
// tag will be used.
type TagSet struct {
	// Schema that the tags are associated with.
	Schema string `json:"schema" jsonschema:"title=Schema"`

	// List of tags for the schema
	List []*cbc.Definition `json:"list" jsonschema:"title=List"`
}

// TagSetForSchema will return the tag set for the provided schema, or nil if it does not exist.
func TagSetForSchema(sets []*TagSet, schema string) *TagSet {
	for _, ts := range sets {
		if ts.Schema == schema {
			return ts
		}
	}
	return nil
}

// WithTags prepares a tags struct
func WithTags(tags ...cbc.Key) Tags {
	return Tags{List: tags}
}

// SetTags is a helper method to set the list of tags.
func (ts *Tags) SetTags(tags ...cbc.Key) {
	ts.List = tags
}

// GetTags returns the list of tags that have been assigned to the object.
func (ts Tags) GetTags() []cbc.Key {
	return ts.List
}

// HasTags returns true if the list of tags contains all of the
// provided tags.
func (ts Tags) HasTags(keys ...cbc.Key) bool {
	if ts.List == nil {
		return false
	}
	for _, k := range keys {
		if !k.In(ts.List...) {
			return false
		}
	}
	return true
}

// Merge will combine the tags from the current set with the tags from the other set,
// ensuring that any duplicates are not overwritten from the original list.
func (ts *TagSet) Merge(other *TagSet) *TagSet {
	if ts == nil {
		return other
	}
	if other == nil || ts.Schema != other.Schema {
		return ts
	}
	nl := ts.List // shallow copy
	for _, t := range other.List {
		found := false
		for _, nlt := range nl {
			if nlt.Key == t.Key {
				// already there
				found = true
				break
			}
		}
		if !found {
			nl = append(nl, t)
		}
	}
	return &TagSet{
		Schema: ts.Schema,
		List:   nl,
	}
}

// Keys returns the list of keys that are available in the tag set.
func (ts *TagSet) Keys() []cbc.Key {
	if ts == nil {
		return []cbc.Key{} // empty
	}
	keys := make([]cbc.Key, len(ts.List))
	for i, k := range ts.List {
		keys[i] = k.Key
	}
	return keys
}

type tagValidation struct {
	keys []cbc.Key
}

// TagsIn provides a validation rule that will ensure the object's tags are contained
// inside the list.
func TagsIn(tags ...cbc.Key) validation.Rule {
	return &tagValidation{keys: tags}
}

// Validate performs the tag validation process.
func (tv *tagValidation) Validate(val interface{}) error {
	list, ok := val.([]cbc.Key)
	if !ok {
		ts, ok := val.(Tags)
		if !ok {
			return nil
		}
		list = ts.List
	}
	for _, k := range list {
		if !k.In(tv.keys...) {
			return fmt.Errorf("'%s' undefined", k)
		}
	}
	return nil
}
