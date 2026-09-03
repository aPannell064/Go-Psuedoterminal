package shell

import (
	"bufio"
	"final/internal/history"
	"final/internal/pipeline"

	"fmt"
	"os"
	"os/user"
	"strings"
)

const reset = "\x1b[0m"
const green = "\x1b[32m"
const blue = "\x1b[34m"

func (shell *Shell) Run() int {
	code := 0

	//Scanner to read lines from standard input
	lineScanner := bufio.NewScanner(os.Stdin)

	// Get user for promt or default to "unknown"
	usr, err := user.Current()
	if err != nil {
		usr = &user.User{Username: "unknown"}
	}

	// Get home directory for prompt formatting
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get home directory: %v\n", err)
		return 1
	}

	// Change username color if in terminal
	uName := usr.Username
	isTerminal := stdoutIsTerminal(os.Stdout)
	if isTerminal {
		uName = green + uName + reset
	}

	for {
		// Check for home to replace with "~"
		cwd := shell.cwd
		if strings.HasPrefix(cwd, home) {
			cwd = strings.Replace(cwd, home, "~", 1)
		}

		// Change path color if in terminal
		if isTerminal {
			cwd = blue + cwd + reset
		}

		// Display prompt
		prompt := fmt.Sprintf("%s:%s$ ", uName, cwd)
		_, err := os.Stdout.WriteString(prompt)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to write to stdout: %v\n", err)
			code = 1
			break
		}

		// Read one line from input
		if !lineScanner.Scan() {
			// If EOF go to next line (to look better)
			_, err := os.Stdout.WriteString("\n")
			if err != nil {
				fmt.Fprintln(os.Stderr, "Failed to write to stdout")
				code = 1
			}
			break
		}

		// Parse the line into tokens using scanArgs
		tokenScanner := bufio.NewScanner(strings.NewReader(lineScanner.Text()))
		tokenScanner.Split(ScanArgs)

		// Array for arguments
		var args []string

		// Add tokens to args
		for tokenScanner.Scan() {
			args = append(args, string(tokenScanner.Bytes()))
		}

		// Add the command to history
		if len(args) > 0 {
			shell.history = history.Append(shell.history, lineScanner.Text())

			// Execute the pipeline and set last exit code
			err = pipeline.StagePipeline(shell.dispatcher, args)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				shell.lastExitCode = 1
			} else {
				shell.lastExitCode = 0
			}
		}

	}

	// Exit shell
	code = shell.exitShell(code)
	return code
}
