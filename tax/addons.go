package tax

import (
	"fmt"

	"github.com/invopop/gobl/cbc"
)

// Addon is an interface that defines the methods that a tax add-on must implement.
type Addon interface {
	// Key that defines how to uniquely idenitfy the add-on.
	Key() cbc.Key

	// Extensions defines the list of extensions that are associated with an add-on.
	Extensions() []*cbc.KeyDefinition

	// Normalize performs the normalization rules for the add-on.
	Normalize(doc any) error

	// Scenarios are applied to documents after normalization and before
	// validation to ensure that form specific extensions have been added
	// to the document.
	Scenarios() []*ScenarioSet

	// Validate performs the validation rules for the add-on.
	Validate(doc any) error

	// Corrections is used to provide a map of correction definitions that
	// are supported by the add-on.
	Corrections() CorrectionSet
}

type addonCollection struct {
	list map[cbc.Key]Addon
}

var addons = newAddonCollection()

func newAddonCollection() *addonCollection {
	return &addonCollection{
		list: make(map[cbc.Key]Addon),
	}
}

// add will register the addon in the collection
func (c *addonCollection) add(a Addon) {
	c.list[a.Key()] = a
}

// RegisterAddon adds a new add-on to the shared global list of tax add-ons. This is
// expected to be called from module init functions.
func RegisterAddon(addon Addon) {
	for _, ext := range addon.Extensions() {
		RegisterExtension(ext)
	}
	addons.add(addon)
}

// AddonForKey provides the add-on for the given key.
func AddonForKey(key cbc.Key) Addon {
	return addons.list[key]
}

// Addons provides the map of keys to addons.
func Addons() map[cbc.Key]Addon {
	return addons.list
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

// BaseAddon provides a base implementation of the Addon interface
// that can be embedded in other add-ons to provide a default
// implemenation and avoid adding empty methods.
type BaseAddon struct{}

// Key provides a default implementation that panics.
func (BaseAddon) Key() cbc.Key {
	panic("Key() not implemented")
}

// Extensions provides a default implementation that returns nil.
func (BaseAddon) Extensions() []*cbc.KeyDefinition {
	return nil
}

// Normalize provides a default implementation that returns nil.
func (BaseAddon) Normalize(_ any) error {
	return nil
}

// Scenarios provides a default implementation that returns nil.
func (BaseAddon) Scenarios() []*ScenarioSet {
	return nil
}

// Validate provides a default implementation that returns nil.
func (BaseAddon) Validate(_ any) error {
	return nil
}

// Corrections provides a default implementation that returns nil.
func (BaseAddon) Corrections() CorrectionSet {
	return nil
}
