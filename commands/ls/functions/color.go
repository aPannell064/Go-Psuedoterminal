package functions

import (
	"errors"
	"fmt"
	"io"
	"os"
)

// List of colors
type color string

const (
	Reset color = "\x1b[0m"
	Red   color = "\x1b[31m"
	Green color = "\x1b[32m"
	Blue  color = "\x1b[34m"
	Cyan  color = "\x1b[36m"
)

// Prints appropriate entries in color, or colorless if regular file
func (c color) ColorPrint(w io.Writer, s string) {
	_, err := w.Write([]byte(string(c) + s + string(Reset)))
	if err != nil && !errors.Is(err, io.ErrClosedPipe) {
		fmt.Fprintln(os.Stderr, err)
	}
}

func (c color) ColorString(s string) string {
	return string(c) + s + string(Reset)
}
func Print(w io.Writer, s string) {
	_, err := w.Write([]byte(s))
	if err != nil && !errors.Is(err, io.ErrClosedPipe) {
		fmt.Fprintln(os.Stderr, err)
	}
}
