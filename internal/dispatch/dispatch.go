package dispatch

import (
	"final/internal/shell/stateInterface"
	"fmt"
	"io"
)

type Dispatcher struct {
	commands map[string]CommandFunc
	state    *stateInterface.State
}

func New(state stateInterface.State) *Dispatcher {
	d := Dispatcher{
		state: &state,
	}
	d.commands = d.registerCommands()

	return &d
}

func (d *Dispatcher) Dispatch(cmd string, args []string, stdin io.Reader, stdout io.Writer) error {
	// Run command if it exists. It should exist at this point, but check anyway
	if commandFunc, exists := d.GetCommand(cmd); exists {
		return commandFunc(args, stdin, stdout)
	}
	return fmt.Errorf("Command not found: '%s'", cmd)
}
