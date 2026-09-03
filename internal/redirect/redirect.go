package redirect

import (
	"fmt"
	"os"
)

func Redirect(args []string) (*os.File, []string, error) {
	stdout := os.Stdout
	var newArgs []string
	var f *os.File
	var err error

	// For each argument
	i := 0
	for i < len(args) {
		switch args[i] {
		// If >:
		case ">":
			// Ensure a file follows
			if i+1 >= len(args) {
				return nil, nil, fmt.Errorf("missing file for >")
			} else if args[i+1] == ">" || args[i+1] == ">>" || args[i+1] == "|" {
				return nil, nil, fmt.Errorf("syntax error near unexpected token '>'")
			}

			// Get file
			file := args[i+1]
			switch file {
			case "1":
				f = os.Stdout
			case "2":
				f = os.Stderr
			default:
				f, err = os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
				if err != nil {
					return nil, nil, err
				}
			}

			// Close previous stdout before updating
			if stdout != os.Stdout {
				stdout.Close()
			}

			stdout = f
			i += 2

		case ">>":
			// Ensure a file follows
			if i+1 >= len(args) {
				return nil, nil, fmt.Errorf("missing file for >>")
			} else if args[i+1] == ">" || args[i+1] == ">>" || args[i+1] == "|" {
				return nil, nil, fmt.Errorf("syntax error near unexpected token '>>'")
			}

			// Get file
			file := args[i+1]
			switch file {
			case "1":
				f = os.Stdout
			case "2":
				f = os.Stderr
			default:
				f, err = os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
				if err != nil {
					return nil, nil, err
				}
			}

			// Close previous stdout before updating
			if stdout != os.Stdout {
				stdout.Close()
			}

			stdout = f
			i += 2

		default:
			// Normal argument
			newArgs = append(newArgs, args[i])
			i++
		}
	}

	return stdout, newArgs, nil
}
