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
