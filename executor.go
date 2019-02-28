package vcsview

import (
	"bufio"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"syscall"
)

// Function for read stdout of the command
type cmdReaderFunc func(s *bufio.Scanner)

// Function for debug messages
type DebugFunc func(m string)

// Command line executor
type Executor struct {
	// Already created command line
	cmd *exec.Cmd

	// Reader of stdout
	reader cmdReaderFunc

	// Debug function
	debugger DebugFunc

	// Text representation of command
	cmdTxt string

	// Command context
	ctx context.Context

	// Function to stop context
	cancel context.CancelFunc
}

// log message if set Debugger
func (e *Executor) log(msg string) {
	if e.debugger != nil {
		e.debugger(msg)
	}
}

// log non-zero command line status
// need provide a command and exec.ExitError struct
func (e *Executor) logCmdNonZeroStatus(err error) {
	if exitError, ok := err.(*exec.ExitError); ok {
		if status, ok := exitError.Sys().(syscall.WaitStatus); ok {
			msg := fmt.Sprintf("Command %s finished with %d status code", e.cmdTxt, status)
			e.log(msg)
			return
		}
	}

	msg := fmt.Sprintf("Command %s finished with non-zero status code", e.cmdTxt)
	e.log(msg)
}

// Create a command stdout pipe and reader function using base reader
// Read will started after sch channel will filled data
// After read sch channel will filled data
func (e *Executor) read(sch chan interface{}) {
	defer func() {
		sch <- struct{}{}
	}()

	out, _ := e.cmd.StdoutPipe()

	<-sch

	e.reader(bufio.NewScanner(out))
	e.cancel()
}

// Run command execution
// This method run async stdout reader and start the command
// To run command async start this method in goroutine
// If command cannot by started or if command fails - returns error
func (e *Executor) Run() error {
	e.log(fmt.Sprintf("execute command: %s", e.cmdTxt))

	sch := make(chan interface{})

	go e.read(sch)

	sch <- struct{}{}

	if err := e.cmd.Start(); err != nil {
		e.logCmdNonZeroStatus(err)
		<- sch
		return err
	}

	<-sch

	if err := e.cmd.Wait(); err != nil {
		e.logCmdNonZeroStatus(err)
		return err
	}

	return nil
}

func NewExecutor(cmd *exec.Cmd, reader cmdReaderFunc, debugger DebugFunc) *Executor {
	e := new(Executor)
	e.cmd = cmd
	e.reader = reader
	e.debugger = debugger
	e.cmdTxt = strings.Join(cmd.Args, " ")
	e.ctx, e.cancel = context.WithCancel(context.Background())
	return e
}