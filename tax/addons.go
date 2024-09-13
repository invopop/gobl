package tax

import (
	"fmt"

	"github.com/invopop/gobl/cbc"
)

// Addon is an interface that defines the methods that a tax add-on must implement.
type Addon struct {
	// Key that defines how to uniquely idenitfy the add-on.
	Key cbc.Key `json:"key" jsonschema:"title=Key"`

	// Extensions defines the list of extensions that are associated with an add-on.
	Extensions []*cbc.KeyDefinition `json:"extensions" jsonschema:"title=Extensions"`

	// Normalizer performs the normalization rules for the add-on.
	Normalizer func(doc any) `json:"-"`

	// Scenarios are applied to documents after normalization and before
	// validation to ensure that form specific extensions have been added
	// to the document.
	Scenarios []*ScenarioSet `json:"scenarios" jsonschema:"title=Scenarios"`

	// Validator performs the validation rules for the add-on.
	Validator func(doc any) error `json:"-"`

	// Corrections is used to provide a map of correction definitions that
	// are supported by the add-on.
	Corrections CorrectionSet `json:"corrections" jsonschema:"title=Corrections"`
}

type addonCollection struct {
	list map[cbc.Key]*Addon
}

var addons = newAddonCollection()

func newAddonCollection() *addonCollection {
	return &addonCollection{
		list: make(map[cbc.Key]*Addon),
	}
}

// add will register the addon in the collection
func (c *addonCollection) add(a *Addon) {
	c.list[a.Key] = a
}

// RegisterAddon adds a new add-on to the shared global list of tax add-ons. This is
// expected to be called from module init functions.
func RegisterAddon(addon *Addon) {
	for _, ext := range addon.Extensions {
		RegisterExtension(ext)
	}
	addons.add(addon)
}

// AddonForKey provides the add-on for the given key.
func AddonForKey(key cbc.Key) *Addon {
	return addons.list[key]
}

// Addons provides the map of keys to addons.
func Addons() map[cbc.Key]*Addon {
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
