// Package data contains both generated and embedded data.
package data

import "embed"

//go:embed regimes schemas

// Content contains the generated regimes and schemes
// ready to serve as an embed.FS.
var Content embed.FS
