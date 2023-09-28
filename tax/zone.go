package tax

import (
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/validation"
)

// ZoneStore defines what is expected of a zone store, given that these database can
// get pretty big, it is more efficient to store them off-disk. Each region should
// decide what to do.
type ZoneStore interface {
	Get(code l10n.Code) *Zone
}

// Zone represents an area inside a country, like a province
// or a state, which shares the basic definitions of the country, but
// may vary in some validation rules.
type Zone struct {
	// Unique zone code.
	Code l10n.Code `json:"code" jsonschema:"title=Code"`
	// Name of the zone to be use if a locality or region is not applicable.
	Name i18n.String `json:"name,omitempty" jsonschema:"title=Name"`
	// Village, town, district, or city name which should coincide with
	// address data.
	Locality i18n.String `json:"locality,omitempty" jsonschema:"title=Locality"`
	// Province, county, or state which should match address data.
	Region i18n.String `json:"region,omitempty" jsonschema:"title=Region"`
	// Codes defines a set of regime specific code mappings.
	Codes cbc.CodeMap `json:"codes,omitempty" jsonschema:"title=Codes"`
	// Any additional information
	Meta cbc.Meta `json:"meta,omitempty" jsonschema:"title=Meta"`
}

// Validate ensures that the zone looks correct.
func (z *Zone) Validate() error {
	err := validation.ValidateStruct(z,
		validation.Field(&z.Code, validation.Required),
		validation.Field(&z.Name),
		validation.Field(&z.Locality),
		validation.Field(&z.Region),
		validation.Field(&z.Meta),
	)
	return err
}

type validateZoneCode struct {
	store ZoneStore
}

// Validate checks to see if the provided zone appears in the store.
func (v *validateZoneCode) Validate(value interface{}) error {
	code, ok := value.(l10n.Code)
	if !ok || code == "" {
		return nil
	}
	if z := v.store.Get(code); z == nil {
		return errors.New("must be a valid value")
	}
	return nil
}

// ZoneIn returns a validation rule that checks to see if the provided
// zone is in the store.
func ZoneIn(store ZoneStore) validation.Rule {
	return &validateZoneCode{store}
}

// ZoneStoreEmbedded implements the ZoneStore interface and provides a standard
// implementation for loading the embedded data on demand.
type ZoneStoreEmbedded struct {
	sync.Mutex
	src   embed.FS
	fn    string
	zones []*Zone
}

// NewZoneStoreEmbedded instantiates a new zone store that will use and embedded
// file system for loading the data.
func NewZoneStoreEmbedded(fs embed.FS, filename string) *ZoneStoreEmbedded {
	return &ZoneStoreEmbedded{src: fs, fn: filename}
}

func (s *ZoneStoreEmbedded) load() {
	s.Lock()
	defer s.Unlock()

	if len(s.zones) == 0 {
		data, err := s.src.ReadFile(s.fn)
		if err != nil {
			panic(fmt.Sprintf("expected to find zone data: %s", err))
		}
		s.zones = make([]*Zone, 0)
		if err := json.Unmarshal(data, &s.zones); err != nil {
			panic(fmt.Sprintf("parsing zone data: %s", err))
		}
	}
}

// Get will load the zone object from the JSON data.
func (s *ZoneStoreEmbedded) Get(code l10n.Code) *Zone {
	s.load()
	for _, z := range s.zones {
		if z.Code == code {
			return z
		}
	}
	return nil
}
