package shell

import (
	"final/internal/dispatch"
	"final/internal/history"
	"os"
	"path/filepath"
)

// Shell struct
type Shell struct {
	cwd          string
	history      []string
	lastExitCode int
	dispatcher   *dispatch.Dispatcher
}

// Create new shell instance
func New() (*Shell, int) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, 1
	}

	// Get history
	var hist []string
	home, err := os.UserHomeDir()
	if err != nil {
		hist = []string{}
	} else {
		path := filepath.Join(home, ".gosh_history")
		hist, err = history.Load(path)
		if err != nil {
			hist = []string{}
		}
	}

	// Create shell
	shell := Shell{
		cwd:     wd,
		history: hist,
	}
	shell.dispatcher = dispatch.New(&shell)

	return &shell, 0
}

// Sets working directory
func (shell *Shell) SetCWD(path string) {
	shell.cwd = path
}

// Gets history
func (shell *Shell) GetHistory() []string {
	return shell.history
}
