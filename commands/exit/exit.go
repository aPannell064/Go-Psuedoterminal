package exit

import (
	"final/internal/history"
	"final/internal/shell/stateInterface"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
)

func Exit(state stateInterface.State) func(args []string, stdin io.Reader, stdout io.Writer) error {
	return func(args []string, stdin io.Reader, stdout io.Writer) error {
		var code int
		var err error

		// If no args code is 0
		if len(args) == 0 {
			code = 0

		} else if len(args) == 1 {
			// If one arg: code = 1 if not valid int
			code, err = strconv.Atoi(args[0])
			if err != nil {
				fmt.Fprintf(os.Stderr, "exit: Exit code '%s' is invalid\n", args[0])
				code = 1
			}
		} else {
			// If multiple args, code is 1
			fmt.Fprintln(os.Stderr, "exit: Too many arguments")
			code = 1
		}

		// Get home directory for saving history
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to save history: %v\n", err)
			return err
		}
		path := filepath.Join(home, ".gosh_history")

		// Get shell history
		hist := state.GetHistory()

		// Save to history
		err = history.Save(path, hist)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to save history: %v\n", err)
			return err
		}

		// Exit with code
		os.Exit(code)

		// Shouldn't actually run
		return nil
	}

}
