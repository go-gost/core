package util

import (
	"bytes"
	"runtime/debug"
)

// From https://stackoverflow.com/a/70723335
func GetGoroutineID() string {
	return string(bytes.Fields(debug.Stack())[1])
}
