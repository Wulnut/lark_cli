package auth

// Mask returns a sanitized version of a sensitive string, showing only the first 6 characters.
func Mask(s string) string {
	if len(s) <= 6 {
		return "***"
	}
	return s[:6] + "***"
}
