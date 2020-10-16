// +build appengine

package log

import (
	"io"
)

func checkIfTerminal(w io.Writer) bool {
	return true
}
