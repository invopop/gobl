package convert

// Input wraps incoming data during detection and lets detectors share
// values they extract from it, so that data read by one detector does
// not need to be parsed again by the next.
type Input struct {
	// Data is the raw incoming document.
	Data []byte

	values map[any]any
}

// NewInput prepares a new Input for the provided data.
func NewInput(data []byte) *Input {
	return &Input{Data: data}
}

// Get returns the value stored for the key, if any.
func (in *Input) Get(key any) (any, bool) {
	v, ok := in.values[key]
	return v, ok
}

// Set stores a value for the key. As with context.Context, packages should
// use their own unexported key types to avoid collisions.
func (in *Input) Set(key, value any) {
	if in.values == nil {
		in.values = make(map[any]any)
	}
	in.values[key] = value
}
