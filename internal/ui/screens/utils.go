package screens

import (
	"harvester/pkg/rendering"
	"strconv"
)

func colorToHex(c rendering.Color) string {
	r := strconv.FormatInt(int64(c.R), 16)
	g := strconv.FormatInt(int64(c.G), 16)
	b := strconv.FormatInt(int64(c.B), 16)
	
	if len(r) == 1 { r = "0" + r }
	if len(g) == 1 { g = "0" + g }
	if len(b) == 1 { b = "0" + b }
	
	return r + g + b
}