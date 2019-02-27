package vcsview

import (
	"bufio"
	"bytes"
	"os/exec"
	"regexp"
	"testing"
)

func TestNewExecutor(t *testing.T) {
	var (
		expectedReaderResult = "testing"
		expectedDebugResult = "testing"
		readerResult string
		debugResult string
	)

	cmd := exec.Command("git", "--version")
	cmd.Dir = gitRepoRealPath
	reader := cmdReaderFunc(func(s *bufio.Scanner) {
		for s.Scan() {
			readerResult += s.Text()
		}
	})
	debugger := DebugFunc(func(msg string) {
		debugResult += msg
	})

	e := NewExecutor(cmd, reader, debugger)

	buf := new(bytes.Buffer)
	s := bufio.NewScanner(buf)
	buf.Write([]byte(expectedReaderResult))
	e.reader(s)
	e.debugger(expectedDebugResult)

	if debugResult != expectedDebugResult {
		t.Errorf("debugResult = %v, want: %v", debugResult, expectedDebugResult)
	}

	if readerResult != expectedReaderResult {
		t.Errorf("readerResult = %v, want: %v", readerResult, expectedReaderResult)
	}
}

func TestExecutor_Log(t *testing.T) {
	var (
		expectedDebugResult = "testing"
		debugResult string
	)

	debugger := DebugFunc(func(msg string) {
		debugResult += msg
	})

	e := Executor{}
	e.log("non testing msg")
	e.debugger = debugger
	e.log(expectedDebugResult)

	if debugResult != expectedDebugResult {
		t.Errorf("debugResult = %v, want: %v", debugResult, expectedDebugResult)
	}
}

func TestExecutor_LogCmdNonZeroStatus(t *testing.T) {
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
		reader := cmdReaderFunc(func(s *bufio.Scanner) {

		})

		c := Cli{v.cmd, debugger}
		cmd := c.command(".", v.args...)

		e := NewExecutor(cmd, reader, debugger)

		e.logCmdNonZeroStatus(v.err)

		if gotMessage != v.message {
			t.Errorf("Executor.logCmdNonZeroStatus(%v) want: %v, got: %v", v.message, v.message, gotMessage)
		}
	}
}

func TestExecutor_Run(t *testing.T) {
	cases := []struct{
		cmd string
		dir string
		params []string
		resultPattern string
		resultLog string
		error bool
	}{
		{
			"git",
			".",
			[]string{"--version"},
			`git[\s]+version.*`,
			"execute command: git --version",
			false,
		},
		{
			"hg",
			".",
			[]string{"log", "--limit=5"},
			"",
			"execute command: hg log --limit=5\nCommand hg log --limit=5 finished with 65280 status code",
			true,
		},
	}

	for key, testCase := range cases {
		gotLog := ""
		gotResult := ""
		done := make(chan interface{}, 1)

		reader := cmdReaderFunc(func(s *bufio.Scanner) {
			for s.Scan() {
				if gotResult != "" {
					gotResult += "\n"
				}
				gotResult += s.Text()
			}

			done <- struct {}{}
		})

		debugger := DebugFunc(func(msg string) {
			if gotLog != "" {
				gotLog += "\n"
			}
			gotLog += msg
		})

		c := Cli{testCase.cmd, debugger}
		cmd := c.command(testCase.dir, testCase.params...)

		e := NewExecutor(cmd, reader, debugger)

		err := e.Run()

		<- done
		close(done)

		if err != nil && !testCase.error {
			t.Errorf("[%d] Cli.execute(%v) unexepcted error: %v, want no errors", key, testCase, err)
			continue
		} else if err == nil && testCase.error {
			t.Errorf("[%d] Cli.execute(%v) exepcted error, but no errors", key, testCase)
			continue
		}

		if testCase.resultPattern == "" && gotResult != "" {
			t.Errorf("[%d] Cli.execute(%v) unexpected result: %v, want empty result", key, testCase, gotResult)
			continue
		} else if testCase.resultPattern != "" && gotResult == "" {
			t.Errorf("[%d] Cli.execute(%v) unexpected empty result, want: %v", key, testCase, testCase.resultPattern)
			continue
		}

		if testCase.resultPattern != "" {
			matched, _ := regexp.MatchString(testCase.resultPattern, gotResult)
			if !matched {
				t.Errorf("[%d] Cli.execute(%v) result doesn't match regexp: %v, got: %v", key, testCase, testCase.resultPattern, gotResult)
			}
		}

		if gotLog != testCase.resultLog {
			t.Errorf("[%d] Cli.execute(%v) want log:\n\t%v,\ngot log:\n\t%v\n", key, testCase, testCase.resultLog, gotLog)
		}
	}
}
