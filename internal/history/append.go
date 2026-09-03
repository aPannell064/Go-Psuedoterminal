package history

func Append(entries []string, cmd string) []string {
	if len(cmd) == 0 {
		return entries
	}
	return append(entries, cmd)
}
