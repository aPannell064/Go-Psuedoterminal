package shell

import (
	"final/internal/history"
	"fmt"
	"os"
	"path/filepath"
)

func (shell *Shell) exitShell(code int) int {
	// Get home directory for saving history
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to save history: %v\n", err)
		// Just exit with code 1 if there is a problem
		return 1
	}
	path := filepath.Join(home, ".gosh_history")

	// Get shell history
	hist := shell.GetHistory()

	// Save to history
	err = history.Save(path, hist)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to save history: %v\n", err)
		// Just exit with code 1 if there is a problem
		return 1
	}

	// Return the code arg
	return code
}
