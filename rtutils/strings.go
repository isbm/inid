package rtutils

import (
	"bytes"
	"fmt"
	"io"
)

// InAny of the arguments, a string "e" we expect.
func InAny(e string, args ...string) bool {
	for _, a := range args {
		if a == e {
			return true
		}
	}
	return false
}

func RCloser2String(stream io.ReadCloser) string {
	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(stream); err != nil {
		return fmt.Sprintf("Stream error: %s", err.Error())
	}
	return buf.String()
}
