package rtutils

// InAny of the arguments, a string "e" we expect.
func InAny(e string, args ...string) bool {
	for _, a := range args {
		if a == e {
			return true
		}
	}
	return false
}
