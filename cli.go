package vcsview

import (
	"fmt"
	"io"
	"os/exec"
	"strings"
	"syscall"
)

type DebugFunc func(m string)

// Common command line interface for each one VCS
type Cli struct {
	// Command line name
	cmd string

	// Debug function wich fixes the log messages
	Debugger DebugFunc
}

// Get full command line with specified params
// For example:
// c.buildCommandStr("log", "--limit=5") = "git log --limit=5"
func (c Cli) buildCommandStr(params ...string) string {
	p := strings.Join(params, " ")

	if p == "" {
		return c.cmd
	}

	return c.cmd +" "+p
}

// log message if set Debugger
func (c Cli) log(m string) {
	if c.Debugger != nil {
		c.Debugger(m)
	}
}

// log non-zero command line status
// need provide text representation of command and exec.ExitError struct
func (c Cli) logCmdNonZeroStatus(cmd string, err error) {
	if exitError, ok := err.(*exec.ExitError); ok {
		if status, ok := exitError.Sys().(syscall.WaitStatus); ok {
			msg := fmt.Sprintf("Command %s finished with %d status code", cmd, status)
			c.log(msg)
			return
		}
	}

	msg := fmt.Sprintf("Command %s finished with non-zero status code", cmd)
	c.log(msg)
}

// Executes CLI command with specified params
// Captures output to Buffer
// If command will fail, returns an Error struct
func (c Cli) execute(dir string, out io.Writer, params ...string) error {
	cmdTxt := c.buildCommandStr(params...)

	c.log(fmt.Sprintf("execute command: %s", cmdTxt))

	cmd := exec.Command(c.cmd, params...)
	cmd.Dir = dir
	cmd.Stdout = out

	if err := cmd.Run(); err != nil {
		c.logCmdNonZeroStatus(cmdTxt, err)
		return err
	}

	return nil
}