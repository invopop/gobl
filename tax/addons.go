package tax

import (
	"context"
	"fmt"
	"sort"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/jsonschema"
	"github.com/invopop/validation"
)

// Addons adds functionality to the owner to be able to handle addons.
type Addons struct {
	// Addons defines a list of keys used to identify tax addons that apply special
	// normalization, scenarios, and validation rules to a document.
	List []cbc.Key `json:"$addons,omitempty" jsonschema:"title=Addons"`
}

// AddonDef is an interface that defines the methods that a tax add-on must implement.
type AddonDef struct {
	// Key that defines how to uniquely idenitfy the add-on.
	Key cbc.Key `json:"key" jsonschema:"title=Key"`

	// Name of the add-on
	Name i18n.String `json:"name" jsonschema:"title=Name"`

	// Description of the add-on
	Description i18n.String `json:"description,omitempty" jsonschema:"title=Description"`

	// Extensions defines the list of extensions that are associated with an add-on.
	Extensions []*cbc.KeyDefinition `json:"extensions" jsonschema:"title=Extensions"`

	// Tags is slice of tag sets that define what can be assigned to each document schema.
	Tags []*TagSet `json:"tags,omitempty" jsonschema:"title=Tags"`

	// Scenarios are applied to documents after normalization and before
	// validation to ensure that form specific extensions have been added
	// to the document.
	Scenarios []*ScenarioSet `json:"scenarios" jsonschema:"title=Scenarios"`

	// Inboxes is a list of keys that are used to identify where copies of
	// documents can be sent.
	Inboxes []*cbc.KeyDefinition `json:"inboxes,omitempty" jsonschema:"title=Inboxes"`

	// Normalizer performs the normalization rules for the add-on.
	Normalizer func(doc any) `json:"-"`

	// Validator performs the validation rules for the add-on.
	Validator func(doc any) error `json:"-"`

	// Corrections is used to provide a map of correction definitions that
	// are supported by the add-on.
	Corrections CorrectionSet `json:"corrections" jsonschema:"title=Corrections"`
}

// WithAddons prepares the Addons struct with the provided list of keys.
func WithAddons(addons ...cbc.Key) Addons {
	return Addons{List: addons}
}

// SetAddons is a helper method to set the list of addons
func (as *Addons) SetAddons(addons ...cbc.Key) {
	as.List = addons
}

// GetAddons provides the list of addon keys in use.
func (as *Addons) GetAddons() []cbc.Key {
	return as.List
}

// GetAddonDefs provides a slice of Addon Definition instances.
func (as Addons) GetAddonDefs() []*AddonDef {
	list := make([]*AddonDef, 0, len(as.List))
	for _, ak := range as.List {
		if a := AddonForKey(ak); a != nil {
			list = append(list, a)
		}
	}
	return list
}

// Validate ensures that the list of addons is valid. This struct is designed to be
// embedded, so we don't perform a regular validation on the struct itself.
func (as Addons) Validate() error {
	return validation.Validate(as.List, validation.Each(AddonRegistered))
}

type addonCollection struct {
	keys []cbc.Key // ordered list
	list map[cbc.Key]*AddonDef
}

var addons = newAddonCollection()

func newAddonCollection() *addonCollection {
	return &addonCollection{
		list: make(map[cbc.Key]*AddonDef),
	}
}

// add will register the addon in the collection
func (c *addonCollection) add(ad *AddonDef) {
	c.keys = append(c.keys, ad.Key)
	sort.Slice(c.keys, func(i, j int) bool {
		return c.keys[i].String() < c.keys[j].String()
	})
	c.list[ad.Key] = ad
}

// RegisterAddonDef adds a new add-on to the shared global list of tax add-on definitions.
// This is expected to be called from module init functions.
func RegisterAddonDef(addon *AddonDef) {
	for _, ext := range addon.Extensions {
		RegisterExtension(ext)
	}
	addons.add(addon)
}

// AddonForKey provides the add-on for the given key.
func AddonForKey(key cbc.Key) *AddonDef {
	return addons.list[key]
}

// AllAddonDefs provides a slice of all the addons defined.
func AllAddonDefs() []*AddonDef {
	all := make([]*AddonDef, len(addons.list))
	for i, ao := range addons.keys {
		all[i] = addons.list[ao]
	}
	return all
}

// WithContext adds this addon to the given context, alongside
// its validator.
func (ad *AddonDef) WithContext(ctx context.Context) context.Context {
	if ad == nil {
		return ctx
	}
	ctx = contextWithValidator(ctx, ad.Validator)
	return ctx
}

type addonValidation struct{}

// AddonRegistered will check that an add-on with the key to be validated
// has been registered.
var AddonRegistered = addonValidation{}

func (addonValidation) Validate(value interface{}) error {
	key, ok := value.(cbc.Key)
	if !ok {
		return nil
	}
	if AddonForKey(key) == nil {
		return fmt.Errorf("addon '%v' not registered", key.String())
	}
	return nil
}

// Validate checks that the add-on has been defined correctly.
func (ad *AddonDef) Validate() error {
	return validation.ValidateStruct(ad,
		validation.Field(&ad.Key, validation.Required, AddonRegistered),
		validation.Field(&ad.Name, validation.Required),
		validation.Field(&ad.Extensions),
		validation.Field(&ad.Inboxes),
		validation.Field(&ad.Tags),
		validation.Field(&ad.Scenarios),
		validation.Field(&ad.Corrections),
	)
}

// JSONSchemaExtend will add the addon options to the JSON list.
func (as Addons) JSONSchemaExtend(js *jsonschema.Schema) {
	props := js.Properties
	if asl, ok := props.Get("$addons"); ok {
		asl.Items.OneOf = make([]*jsonschema.Schema, len(AllAddonDefs()))
		for i, ao := range AllAddonDefs() {
			asl.Items.OneOf[i] = &jsonschema.Schema{
				Const: ao.Key.String(),
				Title: ao.Name.String(),
			}
		}
	}
}
