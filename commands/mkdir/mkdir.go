package mkdir

import (
	"errors"
	"fmt"
	"io"
	"os"
)

func MKDir(args []string, stdin io.Reader, stdout io.Writer) error {
	// For each argument
	for _, arg := range args {
		// Make new directory
		err := os.Mkdir(arg, 0755)
		if err != nil {
			if errors.Is(err, os.ErrExist) {
				fmt.Fprintf(os.Stderr, "mkdir: '%s' already exists\n", arg)
			} else if errors.Is(err, os.ErrNotExist) {
				fmt.Fprintf(os.Stderr, "mkdir: parent path for '%s' does not exist\n", arg)
			}
		}
	}
	return nil
}
