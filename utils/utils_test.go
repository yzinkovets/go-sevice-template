package utils

import (
	"testing"

	"github.com/guregu/null/v5"
)

func TestCleanHostname(t *testing.T) {
	cases := []struct {
		in   string
		want null.String
	}{
		{"", null.String{}},
		{"*", null.String{}},
		{"**", null.StringFrom("**")},
		{" * ", null.String{}},
		{" some.domain ", null.StringFrom("some.domain")},
		{"some.domain", null.StringFrom("some.domain")},
	}

	for _, c := range cases {
		t.Logf("case: %+v", c)
		got := CleanHostname(c.in)
		if got != c.want {
			t.Errorf("CleanHostname(%s) == %s, want %s", c.in, got.String, c.want.String)
		}
	}
}
