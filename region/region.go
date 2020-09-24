package region

// Code for region.
type Code string

var codes = []Code{
	"es", // Spain
}

// Codes provides a list of supported regions.
func Codes() []Code {
	return codes
}
