package scanner_test

import (
	"testing"

	. "github.com/jfreymuth/scanner"
)

func TestQuote(t *testing.T) {
	tests := []struct {
		in   string
		s, e string
		out  string
		ml   bool
		ok   bool
	}{
		{`test("string")`, "\"", "\"", "string", false, true},
		{`test("   string  ")`, "\"", "\"", "   string  ", false, true},
		{`test("string /* not a comment */")`, "\"", "\"", "string /* not a comment */", false, true},

		{"a <<b>> c", "<<", ">>", "b", false, true},

		{`test("string`, "\"", "\"", "", false, false},
		{`test([[string
		]])`,
			"[[", "]]", "", false, false},
		{`test([[string
		]])`,
			"[[", "]]", `string
		`, true, true},
	}

	for _, test := range tests {
		sc := FromString(test.in)
		for !sc.Is(test.s) {
			sc.Rune()
		}
		var out string
		if test.ml {
			out = sc.QuoteMultiline(test.s, test.e, nil)
		} else {
			out = sc.Quote(test.s, test.e, nil)
		}
		if test.ok {
			if sc.Err() != nil {
				t.Errorf("input %q produced error: %s", test.in, sc.Err())
			} else if out != test.out {
				t.Errorf("input %q produced output %q instead of %q", test.in, out, test.out)
			}
		} else {
			if sc.Err() == nil {
				t.Errorf("input %q should produce an error", test.in)
			}
		}
	}
}

func TestString(t *testing.T) {
	tests := []struct {
		in  string
		out string
		ok  bool
	}{
		{`no quotes`, "", false},
		{`"string"`, "string", true},
		{`"string\n"`, "string\n", true},
		{`"string\x00"`, "string\x00", true},
		{`"string\xFF"`, "string\xFF", true},
		{`"string\u00FF"`, "string\u00FF", true},

		{`"string\x0"`, "", false},
		{`"string\x0g"`, "", false},
		{`"string\."`, "", false},

		{`"'"`, "'", true},
		{`"\'"`, "", false},
	}

	for _, test := range tests {
		sc := FromString(test.in)
		out := sc.String()
		if test.ok {
			if sc.Err() != nil {
				t.Errorf("input %q produced error: %s", test.in, sc.Err())
			} else if out != test.out {
				t.Errorf("input %q produced output %q instead of %q", test.in, out, test.out)
			}
		} else {
			if sc.Err() == nil {
				t.Errorf("input %q should produce an error", test.in)
			}
		}
	}
}

func TestChar(t *testing.T) {
	tests := []struct {
		in  string
		out rune
		ok  bool
	}{
		{`a`, 0, false},
		{`'a'`, 'a', true},
		{`'\x00'`, 0, true},
		{`'\xFF'`, 0, false},
		{`'ẏ'`, 'ẏ', true},
		{`'"'`, '"', true},
		{`'\"'`, 0, false},
		{`'\''`, '\'', true},
		{`'ab'`, 0, false},
	}

	for _, test := range tests {
		sc := FromString(test.in)
		out := sc.Char()
		if test.ok {
			if sc.Err() != nil {
				t.Errorf("input %q produced error: %s", test.in, sc.Err())
			} else if out != test.out {
				t.Errorf("input %q produced output %q instead of %q", test.in, out, test.out)
			}
		} else {
			if sc.Err() == nil {
				t.Errorf("input %q should produce an error", test.in)
			}
		}
	}
}

func TestMultilineEscape(t *testing.T) {
	tests := []struct {
		in  string
		out string
		ok  bool
	}{
		{`no quotes`, "", false},
		{`[[unmatched`, "", false},
		{`[[string]]`, "string", true},
		{`[[str\ning]]`, "str\ning", true},
		{`[[first line
		second line\x00]]`, `first line
		second line` + "\x00", true},
		{"[[\\x]]", "", false},
		{"[[\n\\x]]", "", false},
	}

	for _, test := range tests {
		sc := FromString(test.in)
		out := sc.QuoteMultiline("[[", "]]", GoEscaper(0))
		if test.ok {
			if sc.Err() != nil {
				t.Errorf("input %q produced error: %s", test.in, sc.Err())
			} else if out != test.out {
				t.Errorf("input %q produced output %q instead of %q", test.in, out, test.out)
			}
		} else {
			if sc.Err() == nil {
				t.Errorf("input %q should produce an error", test.in)
			}
		}
	}
}

func TestQuoteNext(t *testing.T) {
	tests := []struct {
		in string
		ml bool
	}{
		{`"" nextToken`, false},
		{`[[]] nextToken`, true},
		{`"test" nextToken`, false},
		{`[[test]] nextToken`, true},
		{`[[multi
		line
		test]] nextToken`, true},
		{`[[multi
		line
		test]]
nextToken`, true},
	}

	for _, test := range tests {
		sc := FromString(test.in)
		if test.ml {
			sc.QuoteMultiline("[[", "]]", GoEscaper(0))
		} else {
			sc.Quote(`"`, `"`, GoEscaper(0))
		}
		if next := sc.Ident(); next != "nextToken" {
			if sc.Err() != nil {
				t.Errorf("input %q produced error: %s", test.in, sc.Err())
			} else {
				t.Errorf("input %q produced %s instead of \"nextToken\"", test.in, next)
			}
		}
	}
}
