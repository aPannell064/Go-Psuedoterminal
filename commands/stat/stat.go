package stat

import (
	"bufio"
	"errors"
	"final/commands/stat/functions"
	"fmt"
	"io"
	"os"
)

// Prints detailed metadata about a file or directory
func Stat(args []string, stdin io.Reader, stdout io.Writer) error {
	// Ensure at least one argument is provided
	if len(args) < 1 {
		return fmt.Errorf("stat: no arguments provided")
	}

	// Use a buffered writer because several write system calls will be made
	w := bufio.NewWriter(stdout)

	// For each file:
	for i, file := range args {
		// Get the data
		data, err := functions.GetData(file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "stat: %v\n", err)
			continue
		}

		// Write the data
		err = functions.WriteData(w, data, i >= len(args)-1)
		if err != nil && !errors.Is(err, io.ErrClosedPipe) {
			return fmt.Errorf("stat: %v", err)
		}
	}

	return nil
}
