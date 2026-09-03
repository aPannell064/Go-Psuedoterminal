package echo

import (
	"bufio"
	"errors"
	"fmt"
	"io"
)

func Echo(args []string, stdin io.Reader, stdout io.Writer) error {
	// Use bufio writer tp help with multiple writes
	w := bufio.NewWriter(stdout)
	for _, arg := range args {
		_, err := w.WriteString(arg + " ")
		if err != nil {
			// Handles pipe errors differently because pipelines may be closed to prevent blocking
			if errors.Is(err, io.ErrClosedPipe) {
				return err
			}
			return fmt.Errorf("echo: Could not write to buffer %v", err)
		}
	}

	// Write '\n' just to make things look better
	_, err := w.WriteString("\n")
	if err != nil {
		// Handles pipe errors differently because pipelines may be closed to prevent blocking
		if errors.Is(err, io.ErrClosedPipe) {
			return err
		}
		return fmt.Errorf("echo: Could not write to buffer %v", err)
	}

	return w.Flush()
}
