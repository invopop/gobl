package co

import "embed"

//go:embed data

// Data contains local data specific for Colombia
var Data embed.FS
