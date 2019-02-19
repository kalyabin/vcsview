package vcsview

import (
	"fmt"
	"io"
	"os/exec"
	"strings"
	"syscall"
)

type DebugFunc func(m string)

type Cli struct {
	Cmd string
	Debugger DebugFunc
}

// Execute command with specified params
// For example:
// c.BuildCommandStr("Log", "--limit=5") = "git Log --limit=5"
func (c Cli) BuildCommandStr(params ...string) string {
	p := strings.Join(params, " ")

	if p == "" {
		return c.Cmd
	}

	return c.Cmd+" "+p
}

// Log message if set Debugger
func (c Cli) Log(m string) {
	if c.Debugger != nil {
		c.Debugger(m)
	}
}

// Log non-zero command line status
// need provide text representation of command and exec.ExitError struct
func (c Cli) LogCmdNonZeroStatus(cmd string, err error) {
	if exitError, ok := err.(*exec.ExitError); ok {
		if status, ok := exitError.Sys().(syscall.WaitStatus); ok {
			msg := fmt.Sprintf("Command %s finished with %d status code", cmd, status)
			c.Log(msg)
			return
		}
	}

	msg := fmt.Sprintf("Command %s finished with non-zero status code", cmd)
	c.Log(msg)
}

// Executes CLI command with specified params
// Captures output to Buffer
// If command will fail, returns an Error struct
func (c Cli) Execute(dir string, out io.Writer, params ...string) error {
	cmdTxt := c.BuildCommandStr(params...)

	c.Log(fmt.Sprintf("Execute command: %s", cmdTxt))

	cmd := exec.Command(c.Cmd, params...)
	cmd.Dir = dir
	cmd.Stdout = out

	if err := cmd.Run(); err != nil {
		c.LogCmdNonZeroStatus(cmdTxt, err)
		return err
	}

	return nil
}