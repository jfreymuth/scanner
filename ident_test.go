package scanner_test

import (
	"testing"

	. "github.com/jfreymuth/scanner"
)

func TestIdent(t *testing.T) {
	tests := []struct {
		in    string
		out   string
		valid bool
	}{
		{"a", "a", true},
		{"aa", "aa", true},
		{"a1", "a1", true},
		{"_1", "_1", true},

		{"", "", false},
		{"1", "", false},
		{".", "", false},
		{"(", "", false},

		{"  a  ", "a", true},
		{"a(", "a", true},
	}

	for _, test := range tests {
		sc := FromString(test.in)
		out := sc.Ident()
		if test.valid {
			if sc.Err() != nil {
				t.Errorf("input %q produced error: %s", test.in, sc.Err())
			} else if out != test.out {
				t.Errorf("input %q produced output %q instead of %q", test.in, out, test.out)
			}
		} else if sc.Err() == nil {
			t.Errorf("input %q should produce an error", test.in)
		}
	}
}

func TestWhitespace(t *testing.T) {
	tests := []struct {
		in  string
		out []string
	}{
		{`a`,
			[]string{"a"}},
		{`a b`,
			[]string{"a", "b"}},
		{`            a    b`,
			[]string{"a", "b"}},
		{`a
		  b`,
			[]string{"a", "b"}},
		{`a // comment`,
			[]string{"a"}},
		{`a /* comment */ b`,
			[]string{"a", "b"}},
		{`a
		/* multi
		line
		comment */
		b`,
			[]string{"a", "b"}},
		{`a /**/ /*
		
		*/ /**/ b`,
			[]string{"a", "b"}},
	}

	for _, test := range tests {
		sc := FromString(test.in)
		for i, o := range test.out {
			out := sc.Ident()
			if sc.Err() != nil {
				t.Errorf("input %q produced error: %s", test.in, sc.Err())
			} else if out != o {
				t.Errorf("input %q produced %d. output %q instead of %q", test.in, i, out, test.out)
			}
		}
	}
}
