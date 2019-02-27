package vcsview

import (
	"bufio"
	"os/exec"
)

// Function for read stdout of the command
type cmdReaderFunc func(s *bufio.Scanner)

// Command line executor
type Executor struct {
	// Command line interface
	cli *Cli

	// Already created command line
	cmd *exec.Cmd

	// Reader of stdout
	reader cmdReaderFunc
}

// Create a command stdout pipe and reader function using base reader
// To start the command run Start method
// For the async reading run this function in goroutine
func (e *Executor) read() func() {
	out, _ := e.cmd.StdoutPipe()
	s := bufio.NewScanner(out)

	return func() {
		e.reader(s)
	}
}

// Start command execution
// This method run async stdout reader and start the command
// To run command async start this method in goroutine
// If command cannot by started or if command fails - returns error
func (e *Executor) Start() error {
	go e.read()()

	err := e.cmd.Run()

	if err != nil && e.cli != nil {
		e.cli.logCmdNonZeroStatus(e.cmd, err)
	}

	return err
}

func NewExecutor(cli *Cli, cmd *exec.Cmd, reader cmdReaderFunc) *Executor {
	e := new(Executor)
	e.cmd = cmd
	e.cli = cli
	e.reader = reader
	return e
}