package tax

import (
	"encoding/json"
	"path"
	"sort"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/data"
	"github.com/invopop/gobl/i18n"
)

const (
	cataloguesPath = "catalogues"
)

// A CatalogueDef contains a set of re-useable extensions, scenarios, and validators that
// can be used by addons or tax regimes. This structure is useful for serializing the
// data into JSON for use in external libraries.
type CatalogueDef struct {
	// Key defines a unique identifier for the catalogue.
	Key cbc.Key `json:"key"`
	// Name is the name of the catalogue.
	Name i18n.String `json:"name"`
	// Description is a human readable description of the catalogue.
	Description i18n.String `json:"description,omitempty"`
	// Extensions defines all the extensions offered by the catalogue.
	Extensions []*cbc.Definition `json:"extensions"`
}

// RegisterCatalogueDef will register the catalogue in the global list of catalogues
// and ensure the extensions it contains are available in GOBL.
//
// Unlike other sources, catalogues are loaded from JSON files in the `catalogues`
// directory of the embedded data filesystem.
func RegisterCatalogueDef(filename string) {
	catalogue := &CatalogueDef{}
	out, err := data.Content.ReadFile(path.Join(cataloguesPath, filename))
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(out, catalogue); err != nil {
		panic(err)
	}

	for _, ext := range catalogue.Extensions {
		RegisterExtension(ext)
	}
	catalogues.add(catalogue)
}

// AllCatalogueDefs provides a slice of all the addons defined.
func AllCatalogueDefs() []*CatalogueDef {
	all := make([]*CatalogueDef, len(catalogues.list))
	for i, ao := range catalogues.keys {
		all[i] = catalogues.list[ao]
	}
	return all
}

type catalogueCollection struct {
	keys []cbc.Key // ordered list
	list map[cbc.Key]*CatalogueDef
}

var catalogues = newCatalogueCollection()

func newCatalogueCollection() *catalogueCollection {
	return &catalogueCollection{
		list: make(map[cbc.Key]*CatalogueDef),
	}
}

// add will register the catalogye in the collection
func (c *catalogueCollection) add(cd *CatalogueDef) {
	c.keys = append(c.keys, cd.Key)
	sort.Slice(c.keys, func(i, j int) bool {
		return c.keys[i].String() < c.keys[j].String()
	})
	c.list[cd.Key] = cd
}
