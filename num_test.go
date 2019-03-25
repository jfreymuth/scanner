package scanner_test

import (
	"testing"

	. "github.com/jfreymuth/scanner"
)

func TestInt(t *testing.T) {
	tests := []struct {
		in    string
		out   int
		valid bool
	}{
		{"0", 0, true},
		{"0 ", 0, true},

		{" 0", 0, true},
		{"1 2", 1, true},

		{"", 0, false},
		{"a", 0, false},
		{" a", 0, false},
		{"0a", 0, false},
		{"-", 0, false},

		{"0x", 0, false},
		{"0x(", 0, false},
		{"0xG", 0, false},
		{"0xg", 0, false},

		{"1", 1, true},
		{"-1", -1, true},
		{"12345", 12345, true},
		{"-100", -100, true},

		{"1.9", 0, false},
		{".5", 0, false},

		{"0x10", 0x10, true},
		{"0XFF", 0XFF, true},
		{"0xab", 0xab, true},
		{"012", 012, true},
	}

	for _, test := range tests {
		sc := FromString(test.in)
		if sc.IsInt() != test.valid {
			t.Errorf("input %q: IsInt returned %v", test.in, !test.valid)
		}
		out := sc.Int()
		if test.valid {
			if sc.Err() != nil {
				t.Errorf("input %q produced error: %s", test.in, sc.Err())
			} else if out != test.out {
				t.Errorf("input %q produced output %d instead of %d", test.in, out, test.out)
			}
		} else if sc.Err() == nil {
			t.Errorf("input %q should produce an error", test.in)
		}
	}
}

func TestFloat(t *testing.T) {
	tests := []struct {
		in    string
		out   float64
		valid bool
	}{
		{"0", 0, true},
		{"0 ", 0, true},
		{"0.", 0, true},
		{".0", 0, true},

		{" 0", 0, true},

		{"", 0, false},
		{"a", 0, false},
		{" a", 0, false},
		{"0a", 0, false},

		{"1", 1, true},
		{"-1", -1, true},
		{"12345", 12345, true},
		{"-100", -100, true},
		{"01", 1, true},

		{"0.1", .1, true},
		{".1", .1, true},
		{"2.", 2, true},
		{"-.1", -.1, true},
		{"09.0", 09.0, true},

		{"1e9", 1e9, true},
		{"-1.5e9", -1.5e9, true},
		{"-1e-1", -1e-1, true},
		{"-.1e-1", -.1e-1, true},
		{".1e1", .1e1, true},
		{"1.e1", 1.e1, true},

		{"1e", 0, false},
		{"1e-", 0, false},
		{"1ea", 0, false},
		{"1e.1", 0, false},
		{"1e1a", 0, false},
		{".", 0, false},
		{"-.", 0, false},
		{"-", 0, false},
		{"0x", 0, false},

		{"0xF", 15, true},
		{"0xG", 0, false},
		{"012", 012, true},
		{"-07", -07, true},
		{"-0xf0", -0xf0, true},

		{"1,", 1, true},
		{"1 a", 1, true},
		{"1+x", 1, true},
		{"1-x", 1, true},
		{"-1+x", -1, true},
		{"-1-x", -1, true},
		{"1e+x", 0, false},
		{"1e1-x", 1e1, true},
		{"1e+1-x", 1e+1, true},
		{"1e-1-x", 1e-1, true},
	}

	for _, test := range tests {
		sc := FromString(test.in)
		if sc.IsFloat() != test.valid {
			t.Errorf("input %q: IsFloat returned %v", test.in, !test.valid)
		}
		out := sc.Float()
		if test.valid {
			if sc.Err() != nil {
				t.Errorf("input %q produced error: %s", test.in, sc.Err())
			} else if out != test.out {
				t.Errorf("input %q produced output %g instead of %g", test.in, out, test.out)
			}
		} else if sc.Err() == nil {
			t.Errorf("input %q should produce an error", test.in)
		}
	}
}
