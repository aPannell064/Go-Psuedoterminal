package pipeline

import (
	"final/internal/dispatch"
	"final/internal/redirect"
	"fmt"
	"os"
)

func StagePipeline(d *dispatch.Dispatcher, args []string) error {
	stdout, newArgs, err := redirect.Redirect(args)
	if err != nil {
		return err
	}

	// Defer close stdout just in case
	if stdout != os.Stdout {
		defer stdout.Close()
	}

	// Split by pipeline stage in a map
	var stages []Stage

	// For checking all commands exist and valid pipe syntax
	allExist := true

	// Set cmd to empty
	cmd := ""

	// Read each arg
	i := 0
	for _, arg := range newArgs {
		// Ensure that all stages have a command
		if arg == "|" {
			if cmd == "" {
				return fmt.Errorf("Sytax error near unexpected token '|'")
			}

			// Reset command
			cmd = ""
			i++
		} else {
			if cmd == "" {
				// Create stage for command
				cmd = arg
				stage := Stage{
					Cmd:  cmd,
					Args: []string{},
				}
				stages = append(stages, stage)

				// Print error if command doesn't exist
				if _, exits := d.GetCommand(cmd); !exits {
					fmt.Fprintf(os.Stderr, "%s: command not found\n", cmd)
					allExist = false
				}
			} else {
				// Add argument to the command arguments
				stages[i].Args = append(stages[i].Args, arg)
			}
		}
	}

	// Ensure the existance of all commands, then run pipeline
	if allExist {
		runPipeline(stages, d, os.Stdin, stdout)
	}

	return nil
}
