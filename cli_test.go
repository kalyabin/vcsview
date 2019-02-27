package vcsview

import (
	"testing"
)

func TestCli_Log(t *testing.T) {
	var (
		c = Cli{"git", nil}
		gotMessage = ""
		wantMessage = "testing message"
	)
	c.cmd = "git"

	c.Debugger = DebugFunc(func(msg string) {
		gotMessage = msg
	})

	c.log(wantMessage)

	if gotMessage != wantMessage {
		t.Fatalf("Cli.log(%s) doesn't call Debugger function", wantMessage)
	}
}

func TestCli_LogCmdNonZeroStatus(t *testing.T) {
	cases := []struct{
		cmd string
		args []string
		message string
		err error
	}{
		{"git", []string{"log"}, "Command git log finished with non-zero status code", nil},
	}

	// test with some exited status
	for _, v := range cases {
		gotMessage := ""
		debugger := DebugFunc(func(msg string) {
			gotMessage = msg
		})
		c := Cli{"git", debugger}

		cmd := c.createCommand(".", v.args...)

		c.logCmdNonZeroStatus(cmd, v.err)

		if gotMessage != v.message {
			t.Errorf("Cli.logCmdNonZeroStatus(%v) want: %v, got: %v", v.message, v.message, gotMessage)
		}
	}
}

