package vcsview

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"sync"
	"syscall"
)

// Function for debug messages
type DebugFunc func(m string)

// Function for read stdout of the command
type readerFunc func(s *bufio.Scanner)

// Function for read stderr of the command
type errorFunc func(err error)

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


// Create a command to execute in specified path with command line params
func (c Cli) createCommand(dir string, params ...string) *exec.Cmd {
	cmdTxt := c.buildCommandStr(params...)

	c.log(fmt.Sprintf("execute command: %s", cmdTxt))

	cmd := exec.Command(c.cmd, params...)
	cmd.Dir = dir

	return cmd
}

// Async execution of CLI command with specified params
// Handle stdout with reader function
// Handle error with error handler function
func (Cli) executePipe(cmd *exec.Cmd, reader readerFunc, errHandler errorFunc) func() {
	out, _ := cmd.StdoutPipe()
	s := bufio.NewScanner(out)

	return func() {
		if reader != nil {
			go reader(s)
		}

		handleErr := func(err error) {
			if errHandler != nil {
				errHandler(err)
			}
		}

		if err := cmd.Start(); err != nil {
			handleErr(err)
		}

		if err := cmd.Wait(); err != nil {
			handleErr(err)
		}
	}
}

// Create error handler to handle async errors to errors channel
func (Cli) chanErrHandler(ch chan error) errorFunc {
	return errorFunc(func(err error) {
		ch <- err
	})
}

// Sync execution of CLI command with specified params
// Captures output to Buffer
// If command will fail, returns an Error struct
func (c Cli) execute(dir string, out io.Writer, params ...string) (err error) {
	cmd := c.createCommand(dir, params...)

	wg := &sync.WaitGroup{}

	wg.Add(2)

	reader := readerFunc(func(s *bufio.Scanner) {
		defer wg.Done()

		for s.Scan() {
			out.Write(append(s.Bytes(), []byte("\n")...))
		}
	})

	errHandler := errorFunc(func(e error) {
		if e != nil {
			err = e
		}
	})

	executor := c.executePipe(cmd, reader, errHandler)

	go func() {
		defer wg.Done()

		executor()
	}()

	wg.Wait()

	if err != nil {
		c.logCmdNonZeroStatus(c.buildCommandStr(params...), err)
	}

	return err
}