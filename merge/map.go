package merge

import "maps"

// Maps merges multiple generic map objects and return a shallow copy.
func Maps[M ~map[K]V, K comparable, V any](m1 M, mapArgs ...M) M {
	out := M{}
	maps.Copy(out, m1)

	for _, m := range mapArgs {
		maps.Copy(out, m)
	}

	return out
}
