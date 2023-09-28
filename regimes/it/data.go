package it

import "embed"

//go:embed data

// Data contains local structured data to lazy load when needed.
var Data embed.FS
