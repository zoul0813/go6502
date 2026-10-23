// Package monitor provides the original Apple 1 Woz Monitor, written by Steve
// Wozniak in 1976. See ../README.md for provenance and the matching hex listing.
package monitor

import _ "embed"

const Address uint16 = 0xff00

//go:embed wozmon.rom
var wozmon string

// Wozmon returns an independent copy so callers cannot alter the bundled ROM.
func Wozmon() []byte { return []byte(wozmon) }
