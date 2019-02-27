package vcsview

import (
	"bufio"
	"os/exec"
)

// Function for read stdout of the command
type cmdReaderFunc func(s *bufio.Scanner)

// Command line executor
type Executor struct {
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

	return e.cmd.Run()
}