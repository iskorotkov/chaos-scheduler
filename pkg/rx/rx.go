// Package rx allows to easily create random strings, slices and maps.
package rx

import (
	"fmt"
	"math/rand"
)

// Rstr returns a new prefixed string in format '{prefix}-{random(0, 100)}'
func Rstr(r *rand.Rand, prefix string) string {
	return fmt.Sprintf("%s-%d", prefix, r.Intn(100))
}

// Rmap returns a random map with 0..max entries.
func Rmap(r *rand.Rand, max int) map[string]string {
	values := make(map[string]string)
	for i := 0; i < rand.Intn(max); i++ {
		values[Rstr(r, "key")] = Rstr(r, "value")
	}

	return values
}

// Rslice returns a random slice with 0..max elements.
func Rslice(r *rand.Rand, max int) []string {
	var values []string
	for i := 0; i < rand.Intn(max); i++ {
		values = append(values, Rstr(r, "value"))
	}

	return values
}
