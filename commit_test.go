package vcsview

import (
	"testing"
	"time"
)

func TestCommit_Id(t *testing.T) {
	cases := []string{
		"",
		"1",
		"60a470f",
	}

	for key, testCase := range cases {
		c := Commit{}

		c.id = testCase

		if id := c.Id(); id != testCase {
			t.Errorf("[%d] Commit.Id() = %v, want: %v", key, id, testCase)
		}
	}
}

func TestCommit_Date(t *testing.T) {
	c := Commit{}

	dateFormat := "2006-01-02 03:04:00"
	expectedDateStr := "2019-02-24 10:47:00"

	expectedDate, _ := time.Parse(dateFormat, expectedDateStr)
	c.date = expectedDate

	if date := c.Date().Format(dateFormat); date != expectedDateStr {
		t.Errorf("Commit.Date() = %v, want: %v", date, expectedDate)
	}
}

func TestCommit_Author(t *testing.T) {
	c := Commit{}

	a := Contributor{"name", "test@email.ltd"}
	expectedAuthorName := "name <test@email.ltd>"

	c.author = a

	if author := c.Author(); author.String() != expectedAuthorName {
		t.Errorf("Commit.Author() = %v, want: %v", author, expectedAuthorName)
	}
}

func TestCommit_Message(t *testing.T) {
	c := Commit{}

	expectedMessage := "testing message"
	c.message = expectedMessage

	if message := c.Message(); message != expectedMessage {
		t.Errorf("Commit.Message() = %v, want: %v", message, expectedMessage)
	}
}

func TestCommit_Parents(t *testing.T) {
	cases := [][]string{
		{},
		{"1", "2", "3"},
	}

	for key, testCase := range cases {
		c := Commit{}
		c.parents = testCase

		parents := c.Parents()

		if len(parents) != len(testCase) {
			t.Errorf("[%d] Commit.Parents() = %d commits, want: %d", key, len(parents), len(testCase))
		}
	}
}