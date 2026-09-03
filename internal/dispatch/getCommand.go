package dispatch

func (d *Dispatcher) GetCommand(cmd string) (CommandFunc, bool) {
	command, exists := d.commands[cmd]
	return command, exists
}
