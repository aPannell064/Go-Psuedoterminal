package functions

import (
	"bufio"
	"errors"
	"fmt"
	"io"
)

func WriteData(w *bufio.Writer, data []string, lastFile bool) error {
	// For each field:
	for _, entry := range data {

		// Write the entry
		_, err := w.WriteString(entry)
		if err != nil {
			// Handles pipe errors differently because pipelines may be closed to prevent blocking
			if errors.Is(err, io.ErrClosedPipe) {
				return err
			}
			return fmt.Errorf("Cannot write to buffer: %v", err)
		}
	}

	// Add '\n' if this is the last file
	if !lastFile {
		_, err := w.WriteString("\n")
		if err != nil {
			// Handles pipe errors differently because pipelines may be closed to prevent blocking
			if errors.Is(err, io.ErrClosedPipe) {
				return err
			}
			return fmt.Errorf("Cannot write to buffer: %v", err)
		}
	}

	// Return the error for flushing
	return w.Flush()
}
