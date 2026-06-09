package l10n

import (
	"github.com/invopop/gobl/cal"
)

// UnionDefs is a list of UnionDef objects.
type UnionDefs []*UnionDef

// UnionDef represents the definition of a group of countries with
// a common political and economic union.
type UnionDef struct {
	// Short identifier for the union.
	Code Code `json:"code" jsonschema:"title=Code"`
	// Name of the union
	Name string `json:"name" jsonschema:"title=Name"`
	// Members defines the list of members of the union by
	// their ISO or Tax country codes. Members may be duplicated
	// if they have left and rejoined the union during a specific
	// period.
	Members []*UnionMember `json:"members" jsonschema:"title=Members"`
}

// UnionMember represents a country that is a member of a union,
// including the date when the country joined and may have left.
type UnionMember struct {
	// ISO 3166-2 Country Code
	Code Code `json:"code" jsonschema:"title=Code"`
	// Alternative code that can be used for lookups
	AltCode Code `json:"alt_code,omitempty" jsonschema:"title=Alt Code"`
	// Date when the state became a member.
	Since cal.Date `json:"since" jsonschema:"title=Since"`
	// Date of departure.
	Until cal.Date `json:"until,omitempty" jsonschema:"title=Until"`
}

// Len provides the length of the union definitions
func (uds UnionDefs) Len() int {
	return len(uds)
}

// Code finds the union definition for the given country code
func (uds UnionDefs) Code(c Code) *UnionDef {
	for _, v := range uds {
		if v.Code == c {
			return v
		}
	}
	return nil
}

// HasMember checks if the given country code is a member of the union
// at the current point in time.
func (ud *UnionDef) HasMember(c Code) bool {
	return ud.HasMemberOn(cal.Today(), c)
}

// HasMemberOn checks if the given country code is a member of the union
// at a specific point in time.
func (ud *UnionDef) HasMemberOn(date cal.Date, c Code) bool {
	if c == "" {
		return false
	}
	for _, m := range ud.Members {
		if m.Code == c || m.AltCode == c {
			if m.On(date) {
				return true
			}
		}
	}
	return false
}

// On checks if the given country code is a member of the union
// at a specific point in time.
func (ud *UnionMember) On(date cal.Date) bool {
	return date.DaysSince(ud.Since.Date) >= 0 && (ud.Until.IsZero() || date.DaysSince(ud.Until.Date) <= 0)
}
