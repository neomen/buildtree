package utils

// IsTreeSymbol checks if a rune is a tree diagram symbol
func IsTreeSymbol(r rune) bool {
	switch r {
	case ' ', '│', '├', '└', '─', '|', '-', '+', '\\', '/', '>', ':', '\'':
		return true
	default:
		return false
	}
}
