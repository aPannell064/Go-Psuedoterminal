package dispatch

import (
	"final/commands/cat"
	"final/commands/cd"
	"final/commands/echo"
	"final/commands/exit"
	"final/commands/grep"
	"final/commands/history"
	"final/commands/ls"
	"final/commands/mkdir"
	"final/commands/pwd"
	"final/commands/rm"
	"final/commands/rmdir"
	"final/commands/stat"
	"final/commands/touch"
	"fmt"
	"os"
)

var cmds map[string]CommandFunc

func (d *Dispatcher) registerCommands() map[string]CommandFunc {
	cmds = make(map[string]CommandFunc)

	// Rgister cat command
	var catCommand CommandFunc = cat.CatMain
	register("cat", catCommand, true)

	// Register ls command
	var lsCommand CommandFunc = ls.LsMain
	register("ls", lsCommand, false)

	// Register cd command
	var cdCommand CommandFunc = cd.Cd(*d.state)
	register("cd", cdCommand, false)

	// Register pwd command
	var pwdCommand CommandFunc = pwd.Pwd
	register("pwd", pwdCommand, false)

	// Register echo command
	var echoCommand CommandFunc = echo.Echo
	register("echo", echoCommand, false)

	// Register touch command
	var touchCommand CommandFunc = touch.Touch
	register("touch", touchCommand, false)

	// Register mkdir command
	var mkDirCommand CommandFunc = mkdir.MKDir
	register("mkdir", mkDirCommand, false)

	// Register rm command
	var rmCommand CommandFunc = rm.RM
	register("rm", rmCommand, false)

	// Register rmdir command
	var rmdirCommand CommandFunc = rmdir.RMDir
	register("rmdir", rmdirCommand, false)

	// Register grep command
	var grepCommand CommandFunc = grep.Grep
	register("grep", grepCommand, true)

	// Register history command
	var historyCommand CommandFunc = history.History(*d.state)
	register("history", historyCommand, false)

	// Register stat command
	var statCommand CommandFunc = stat.Stat
	register("stat", statCommand, false)

	// Register exit command
	var exitCommand CommandFunc = exit.Exit(*d.state)
	register("exit", exitCommand, false)

	return cmds

}

func register(name string, cmd CommandFunc, readsInput bool) {
	if _, exists := cmds[name]; exists {
		fmt.Fprintln(os.Stderr, "Error registering command: duplicate command")
	}
	cmds[name] = cmd
}
