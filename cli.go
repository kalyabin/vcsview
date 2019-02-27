package vcsview

import (
	"os/exec"
)

// Common command line interface for each one VCS
type Cli struct {
	// Command line name
	cmd string

	// Debug function which fixes the log messages
	Debugger DebugFunc
}

// Create a command to execute in specified path with command line params
func (c Cli) command(dir string, params ...string) *exec.Cmd {
	cmd := exec.Command(c.cmd, params...)
	cmd.Dir = dir

	return cmd
}

// Create executor instance will execute the command
func (c *Cli) executor(cmd *exec.Cmd, reader cmdReaderFunc) *Executor {
	return NewExecutor(cmd, reader, c.Debugger)
}