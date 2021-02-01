package rx

import (
	"fmt"
	"math/rand"
)

func Rstr(r *rand.Rand, prefix string) string {
	return fmt.Sprintf("%s-%d", prefix, r.Intn(100))
}

func Rmap(r *rand.Rand, max int) map[string]string {
	values := make(map[string]string)
	for i := 0; i < rand.Intn(max); i++ {
		values[Rstr(r, "key")] = Rstr(r, "value")
	}

	return values
}

func Rslice(r *rand.Rand, max int) []string {
	var values []string
	for i := 0; i < rand.Intn(max); i++ {
		values = append(values, Rstr(r, "value"))
	}

	return values
}
