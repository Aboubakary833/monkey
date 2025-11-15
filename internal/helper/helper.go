package helper


// IsCharAllowedInKeyOrVar check if wether a given character
// is a valid prefix for an identifier or allowed in a variable name
func IsCharAllowedInKeyOrVar(char byte) bool {
	return 'a' <= char && char <= 'z' || 'A' <= char && char <= 'Z' || char == '_'
}

func IsDigit(char byte) bool {
	return '0' <= char && char <= '9'
}
