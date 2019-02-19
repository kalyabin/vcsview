package vcsview

import (
	"bytes"
	"regexp"
	"testing"
)

func TestCli_Log(t *testing.T) {
	var (
		c = Cli{"git", nil}
		gotMessage = ""
		wantMessage = "testing message"
	)
	c.Cmd = "git"

	c.Debugger = DebugFunc(func(msg string) {
		gotMessage = msg
	})

	c.Log(wantMessage)

	if gotMessage != wantMessage {
		t.Fatalf("Cli.Log(%s) doesn't call Debugger function", wantMessage)
	}
}

func TestCli_LogCmdNonZeroStatus(t *testing.T) {
	cases := []struct{
		cmd string
		message string
		err error
	}{
		{"git log", "Command git log finished with non-zero status code", nil},
	}

	// test with some exited status
	for _, v := range cases {
		gotMessage := ""
		debugger := DebugFunc(func(msg string) {
			gotMessage = msg
		})
		c := Cli{"git", debugger}
		c.LogCmdNonZeroStatus(v.cmd, v.err)
		if gotMessage != v.message {
			t.Errorf("Cli.LogCmdNonZeroStatus(%v) want: %v, got: %v", v.message, v.message, gotMessage)
		}
	}
}

func TestCli_BuildCommandStr(t *testing.T) {
	cases := []struct{
		params []string
		result string
	}{
		{[]string{""}, "git"},
		{[]string{"Log"}, "git Log"},
		{[]string{"Log --limit=5"}, "git Log --limit=5"},
		{[]string{"Log", "--limit=5"}, "git Log --limit=5"},
		{[]string{"Log", "-n", "5"}, "git Log -n 5"},
	}

	for _, v := range cases {
		c := Cli{"git", nil}

		if result := c.BuildCommandStr(v.params...); result != v.result {
			t.Errorf("Cli.BuildCommandStr(%#v) = %v, want: %v", v.params, result, v.result)
		}
	}
}

func TestCli_Execute(t *testing.T) {
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
			"Execute command: git --version",
			false,
		},
		{
			"hg",
			".",
			[]string{"log", "--limit=5"},
			"",
			"Execute command: hg log --limit=5\nCommand hg log --limit=5 finished with 65280 status code",
			true,
		},
	}

	for key, testCase := range cases {
		gotLog := ""
		debugger := DebugFunc(func(msg string) {
			if gotLog != "" {
				gotLog += "\n"
			}
			gotLog += msg
		})

		buffer := new(bytes.Buffer)
		c := Cli{testCase.cmd, debugger}

		err := c.Execute(testCase.dir, buffer, testCase.params...)

		if err != nil && !testCase.error {
			t.Errorf("[%d] Cli.Execute(%v) unexepcted error: %v, want no errors", key, testCase, err)
			continue
		} else if err == nil && testCase.error {
			t.Errorf("[%d] Cli.Execute(%v) exepcted error, but no errors", key, testCase)
			continue
		}

		gotResult := buffer.String()

		if testCase.resultPattern == "" && gotResult != "" {
			t.Errorf("[%d] Cli.Execute(%v) unexpected result: %v, want empty result", key, testCase, gotResult)
			continue
		} else if testCase.resultPattern != "" && gotResult == "" {
			t.Errorf("[%d] Cli.Execute(%v) unexpected empty result, want: %v", key, testCase, testCase.resultPattern)
			continue
		}

		if testCase.resultPattern != "" {
			matched, _ := regexp.MatchString(testCase.resultPattern, gotResult)
			if !matched {
				t.Errorf("[%d] Cli.Execute(%v) result doesn't match regexp: %v, got: %v", key, testCase, testCase.resultPattern, gotResult)
			}
		}

		if gotLog != testCase.resultLog {
			t.Errorf("[%d] Cli.Execute(%v) want log:\n\t%v,\ngot log:\n\t%v\n", key, testCase, testCase.resultLog, gotLog)
		}
	}
}

