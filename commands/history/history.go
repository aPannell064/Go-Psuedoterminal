package history

import (
	"bufio"
	"errors"
	"final/internal/shell/stateInterface"
	"fmt"
	"io"
)

func History(state stateInterface.State) func(args []string, stdin io.Reader, stdout io.Writer) error {
	return func(args []string, stdin io.Reader, stdout io.Writer) error {
		// Bufio writer for writing large amounts of data
		w := bufio.NewWriter(stdout)

		// Get history
		hist := state.GetHistory()

		// Write each entry to buffered writer with entry number
		for i, cmd := range hist {
			entry := fmt.Sprintf("%4d %s\n", i+1, cmd)
			_, err := w.WriteString(entry)
			if err != nil {
				if errors.Is(err, io.ErrClosedPipe) {
					return err
				}
				return fmt.Errorf("history: Could not write to buffer %v", err)
			}
		}

		// Flush buffer
		err := w.Flush()
		if err != nil {
			if errors.Is(err, io.ErrClosedPipe) {
				return err
			}
			return fmt.Errorf("history: Could not flush buffer %v", err)
		}

		return nil
	}
}
