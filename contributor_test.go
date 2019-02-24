package vcsview

import (
	"testing"
)

func TestContributor_Name(t *testing.T) {
	cases := []string{
		"",
		"testing name",
	}

	for key, testCase := range cases {
		c := Contributor{}

		c.name = testCase

		if name := c.Name(); name != testCase {
			t.Errorf("[%d] Contributor.Name() = %v, want: %v", key, name, testCase)
		}
	}
}

func TestContributor_Email(t *testing.T) {
	cases := []string{
		"",
		"test@email.ltd",
	}

	for key, testCase := range cases {
		c := Contributor{}

		c.email = testCase

		if email := c.Email(); email != testCase {
			t.Errorf("[%d] Contributor.Email() = %v, want: %v", key, email, testCase)
		}
	}
}

func TestContributor_String(t *testing.T) {
	cases := []struct{
		name string
		email string
		want string
	}{
		{"", "", ""},
		{"", "email@email.ltd", " <email@email.ltd>"},
		{"test", "", "test"},
		{"test", "email@email.ltd", "test <email@email.ltd>"},
	}

	for key, testCase := range cases {
		c := Contributor{testCase.name, testCase.email}

		if contributor := c.String(); contributor != testCase.want {
			t.Errorf("[%d] Contributor.String() = %v, want: %v", key, contributor, testCase.want)
		}
	}
}
