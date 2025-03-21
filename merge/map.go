package merge

import "maps"

// MapStringString creates a new map and loads it from map args
// this function takes a least 2 args and later map args take precedence.
func MapStringString(m1 map[string]string, mapArgs ...map[string]string) map[string]string {
	// populate initial map
	outMap := map[string]string{}
	maps.Copy(outMap, m1)

	// iterate all args
	for _, m := range mapArgs {
		maps.Copy(outMap, m)
	}

	return outMap
}
