package vcsview

import (
	"fmt"
	"io"
	"os/exec"
	"strings"
	"syscall"
)

// Function for debug messages
type DebugFunc func(m string)

// Common command line interface for each one VCS
type Cli struct {
	// Command line name
	cmd string

	// Debug function wich fixes the log messages
	Debugger DebugFunc
}

// log message if set Debugger
func (c Cli) log(m string) {
	if c.Debugger != nil {
		c.Debugger(m)
	}
}

// log non-zero command line status
// need provide a command and exec.ExitError struct
func (c Cli) logCmdNonZeroStatus(cmd *exec.Cmd, err error) {
	cmdTxt := strings.Join(cmd.Args, " ")

	if exitError, ok := err.(*exec.ExitError); ok {
		if status, ok := exitError.Sys().(syscall.WaitStatus); ok {
			msg := fmt.Sprintf("Command %s finished with %d status code", cmdTxt, status)
			c.log(msg)
			return
		}
	}

	msg := fmt.Sprintf("Command %s finished with non-zero status code", cmdTxt)
	c.log(msg)
}


// Create a command to execute in specified path with command line params
func (c Cli) createCommand(dir string, params ...string) *exec.Cmd {
	cmd := exec.Command(c.cmd, params...)
	cmd.Dir = dir

	c.log(fmt.Sprintf("execute command: %s", strings.Join(cmd.Args, " ")))

	return cmd
}

// Sync execution of CLI command with specified params
// Captures output to Buffer
// If command will fail, returns an Error struct
func (c Cli) execute(dir string, out io.Writer, params ...string) (err error) {
	cmd := c.createCommand(dir, params...)
	cmd.Stdout = out

	if err := cmd.Run(); err != nil {
		c.logCmdNonZeroStatus(cmd, err)
		return err
	}

	return nil
}