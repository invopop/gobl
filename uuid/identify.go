package uuid

// Identify defines a struct that can be used to identify a document by a UUID.
type Identify struct {
	UUID UUID `json:"uuid,omitempty" jsonschema:"title=UUID,description=Universally Unique Identifier."`
}

// GetUUID returns the UUID of the document.
func (d *Identify) GetUUID() UUID {
	return d.UUID
}

// SetUUID sets the UUID of the document.
func (d *Identify) SetUUID(id UUID) {
	d.UUID = id
}

// IdentifyParse will parse the provided string as a UUID in the UUID field,
// or panic. This is mainly meant to be used in tests.
func IdentifyParse(s string) Identify {
	return Identify{UUID: MustParse(s)}
}

// IdentifyV1 is a helper method to generate a V1 uuid ready to embed.
func IdentifyV1() Identify {
	return Identify{UUID: V1()}
}

// IdentifyV3 is a helper method to generate a V3 uuid ready to embed.
func IdentifyV3(ns UUID, data []byte) Identify {
	return Identify{UUID: V3(ns, data)}
}

// IdentifyV4 is a helper method to generate a V4 uuid ready to embed.
func IdentifyV4() Identify {
	return Identify{UUID: V4()}
}

// IdentifyV5 is a helper method to generate a V5 uuid ready to embed.
func IdentifyV5(ns UUID, data []byte) Identify {
	return Identify{UUID: V5(ns, data)}
}
