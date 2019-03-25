package scanner

import (
	"testing"
)

func TestNext(t *testing.T) {
	tests := []struct {
		in  string
		out []rune
	}{
		{"abc", []rune("abc\x00")},
		{"abc", []rune("abc")},
		{"ab c", []rune("ab c\x00")},
		{"ab\n c", []rune("ab\n c\x00")},
		{"abc", []rune("abc\x00\x00\x00\x00")},
	}
	for _, test := range tests {
		sc := FromString(test.in)
		for i, r := range test.out {
			if sc.r != r {
				t.Errorf("input %q produced %d. output %q instead of %q", test.in, i, sc.r, r)
			}
			sc.next()
		}
		if !sc.End() {
			t.Errorf("input %q: End returned false", test.in)
		}
	}
}

func TestRune(t *testing.T) {
	tests := []struct {
		in  string
		out []rune
	}{
		{"abc", []rune("abc")},
		{"ab c", []rune("abc")},
		{"ab\n c", []rune("abc")},
		{"abc", []rune("abc\x00\x00")},
	}
	for _, test := range tests {
		sc := FromString(test.in)
		for i, r := range test.out {
			p := sc.Peek()
			sr := sc.Rune()
			if sr != p {
				t.Errorf("input %q[%d]: Peek produced different output from Rune ('%q' != '%q')", test.in, i, p, sr)
			}
			if sr != r {
				t.Errorf("input %q produced %d. output \"%q\" instead of \"%q\"", test.in, i, sr, r)
			}
		}
	}
}

func TestDemand(t *testing.T) {
	tests := []struct {
		in      string
		demands []string
		ok      bool
	}{
		{"test", []string{"test"}, true},
		{"test", []string{"te", "st"}, true},
		{"test", []string{"t"}, true},
		{"a /* comment */ b // comment\n c", []string{"a", "b", "c"}, true},

		{">=", []string{">="}, true},
		{">=", []string{">", "="}, true},
		{"> =", []string{">="}, false},
		{"> =", []string{">", "="}, true},
	}
	for _, test := range tests {
		sc := FromString(test.in)
		for _, d := range test.demands {
			sc.Demand(d)
		}
		if sc.Err() == nil {
			if !test.ok {
				t.Errorf("input %q, %q should produce an error", test.in, test.demands)
			}
		} else {
			if test.ok {
				t.Errorf("input %q, %q produced error: %s", test.in, test.demands, sc.Err())
			}
		}
	}
}
