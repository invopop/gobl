// Package data contains both generated and embedded data.
package data

import "embed"

//go:embed currency regimes schemas addons catalogues

// Content contains the generated regimes and schemes
// ready to serve as an embed.FS.
var Content embed.FS
