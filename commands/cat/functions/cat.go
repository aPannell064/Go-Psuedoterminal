package functions

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
)

// More complex cat for flags
func Cat(files []string, flags []bool, stdin io.Reader, stdout io.Writer) error {
	// Set linenum and newline
	linenum := 1
	newline := true

	// Source as an io.Reader
	var src io.Reader

	// For each file
	for _, file := range files {
		if file == "-" {
			// Use stdin if -
			src = stdin
		} else {
			// Otherwise, read from file
			f, err := os.Open(file)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s: No such file or directory\n")
				continue
			}
			defer f.Close()
			src = f
		}

		// Buffered writer
		w := bufio.NewWriter(stdout)

		err := reader(src, w, flags, &linenum, &newline)
		if err != nil && !errors.Is(err, io.ErrClosedPipe) {
			fmt.Fprint(os.Stderr, err)
		}
	}
	return nil
}
