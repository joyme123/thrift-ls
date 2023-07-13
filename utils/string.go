package utils

func Space(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r'
}
