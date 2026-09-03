package stateInterface

// State interface for closures to avoid circular dependencies and pass information between shell and commands
type State interface {
	GetHistory() []string
	SetCWD(string)
}
