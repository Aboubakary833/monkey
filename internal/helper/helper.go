package helper


// IsCharAllowedInKeyOrVar check if wether a given character
// is a valid prefix for an identifier or allowed in a variable name
func IsCharAllowedInKeyOrVar(char byte) bool {
	return 'a' <= char && char <= 'z' || 'A' <= char && char <= 'Z' || char == '_'
}

func IsDigit(char byte) bool {
	return '0' <= char && char <= '9'
}

// FlipMap take a map of type map[K]V and return a new map of type map[V]K
func FlipMap[K comparable, V comparable](original map[K]V) map[V]K {
	var m = make(map[V]K, len(original))

	for k, v := range original {
		m[v] = k
	}

	return m
}
