package tax

import (
	"context"
	"fmt"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/validation"
)

// Addons adds functionality to the owner to be able to handle addons.
type Addons struct {
	// Addons defines a list of keys used to identify tax addons that apply special
	// normalization, scenarios, and validation rules to a document.
	Addons []cbc.Key `json:"$addons,omitempty" jsonschema:"title=Addons"`
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
	return Addons{Addons: addons}
}

// SetAddons is a helper method to set the list of addons
func (as *Addons) SetAddons(addons ...cbc.Key) {
	as.Addons = addons
}

// GetAddons provides a slice of Addon instances.
func (as Addons) GetAddons() []*AddonDef {
	list := make([]*AddonDef, 0, len(as.Addons))
	for _, ak := range as.Addons {
		if a := AddonForKey(ak); a != nil {
			list = append(list, a)
		}
	}
	return list
}

// Validate ensures that the list of addons is valid. This struct is designed to be
// embedded, so we don't perform a regular validation on the struct itself.
func (as Addons) Validate() error {
	return validation.Validate(as.Addons, validation.Each(AddonRegistered))
}

type addonCollection struct {
	list map[cbc.Key]*AddonDef
}

var addons = newAddonCollection()

func newAddonCollection() *addonCollection {
	return &addonCollection{
		list: make(map[cbc.Key]*AddonDef),
	}
}

// add will register the addon in the collection
func (c *addonCollection) add(a *AddonDef) {
	c.list[a.Key] = a
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

// AllAddons provides a slice of all the addons defined.
func AllAddons() []*AddonDef {
	all := make([]*AddonDef, len(addons.list))
	i := 0
	for _, a := range addons.list {
		all[i] = a
		i++
	}
	return all
}

// WithContext adds this addon to the given context, alongside
// its validator.
func (a *AddonDef) WithContext(ctx context.Context) context.Context {
	if a == nil {
		return ctx
	}
	ctx = contextWithValidator(ctx, a.Validator)
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
func (ao *AddonDef) Validate() error {
	return validation.ValidateStruct(ao,
		validation.Field(&ao.Key, validation.Required, AddonRegistered),
		validation.Field(&ao.Name, validation.Required),
		validation.Field(&ao.Extensions),
		validation.Field(&ao.Tags),
		validation.Field(&ao.Scenarios),
		validation.Field(&ao.Corrections),
	)
}
