package scanner

import (
	"bytes"
	"errors"
	"strconv"
	"strings"
	"unicode/utf8"
)

// String parses a string enclosed in double quotes.
func (s *Scanner) String() string {
	return s.Quote("\"", "\"", GoEscaper('"'))
}

// Char parses a single rune enclosed in signle quotes.
func (s *Scanner) Char() rune {
	q := s.Quote("'", "'", GoEscaper('\''))
	r, sz := utf8.DecodeRuneInString(q)
	if r == utf8.RuneError || len(q) != sz {
		s.Fail("invalid character")
		return 0
	}
	return r
}

// Quote returns all text, including whitespace and comments, between the start and end tokens.
// Quote causes an error if the current token is not start or if the current line does not contain end.
func (s *Scanner) Quote(start, end string, esc Escaper) string {
	if !s.Is(start) {
		s.Failf("'%s' expected", start)
		return ""
	}
	s.pos += len(start)
	out := ""
	l := strings.Index(s.line[s.pos:], end)
	if esc != nil {
		e := esc.EscapeIndex(s.line[s.pos:])
		for e >= 0 && e < l {
			str, n, err := esc.Unescape(s.line[s.pos+e:])
			out += s.line[s.pos:s.pos+e] + str
			s.pos += e + n
			if err != nil {
				s.Fail(err.Error())
				return ""
			}
			e = esc.EscapeIndex(s.line[s.pos:])
			l = strings.Index(s.line[s.pos:], end)
		}
	}
	if l >= 0 {
		if out == "" {
			out = s.line[s.pos : s.pos+l]
		} else {
			out += s.line[s.pos : s.pos+l]
		}
		s.pos += l + len(end)
		s.space()
		return out
	}
	s.pos = len(s.line)
	s.Failf("'%s' expected", end)
	return ""
}

// QuoteMultiline returns all text, including whitespace, line breaks, and comments, between the start and end tokens.
// QuoteMultiline causes an error if the current token is not start or if the input does not contain end.
func (s *Scanner) QuoteMultiline(start, end string, esc Escaper) string {
	if !s.Is(start) {
		s.Failf("'%s' expected", start)
		return ""
	}
	s.pos += len(start)
	var out bytes.Buffer
	l := strings.Index(s.line[s.pos:], end)
	if esc != nil {
		e := esc.EscapeIndex(s.line[s.pos:])
		for e >= 0 && e < l {
			str, n, err := esc.Unescape(s.line[s.pos+e:])
			out.WriteString(s.line[s.pos : s.pos+e])
			out.WriteString(str)
			s.pos += e + n
			if err != nil {
				s.Fail(err.Error())
				return ""
			}
			e = esc.EscapeIndex(s.line[s.pos:])
			l = strings.Index(s.line[s.pos:], end)
		}
	}

	for l < 0 {
		out.WriteString(s.line[s.pos:])
		out.WriteByte('\n')
		s.pos = len(s.line)
		s.update()
		if s.err != nil {
			s.Failf("'%s' expected", end)
			return ""
		}
		l = strings.Index(s.line[s.pos:], end)
		if esc != nil {
			e := esc.EscapeIndex(s.line[s.pos:])
			for e >= 0 && e < l {
				str, n, err := esc.Unescape(s.line[s.pos+e:])
				out.WriteString(s.line[s.pos : s.pos+e])
				out.WriteString(str)
				s.pos += e + n
				if err != nil {
					s.Fail(err.Error())
					return ""
				}
				e = esc.EscapeIndex(s.line[s.pos:])
				l = strings.Index(s.line[s.pos:], end)
			}
		}
	}
	out.WriteString(s.line[s.pos : s.pos+l])
	s.pos += l + len(end)
	s.space()
	return out.String()
}

// An Escaper can detect and replace escape sequences.
type Escaper interface {
	// EscapeIndex returns the index of the first escape sequence in the string, or -1
	EscapeIndex(string) int
	// Unescape parses the escape sequence at the start of the string.
	// It returns:
	//     a string the escape sequence should be replaced with,
	//     the length of the escape sequence,
	//     and an error, if the escape sequence is invalid.
	Unescape(string) (string, int, error)
}

// ErrEscape is a convenient error that can be returned by an implementation of Escaper
var ErrEscape = errors.New("invalid escape sequence")

// A GoEscaper supports escape sequences exactly as defined by the go specification.
type GoEscaper byte

// EscapeIndex implements the Escaper interface
func (q GoEscaper) EscapeIndex(s string) int { return strings.IndexByte(s, '\\') }

// Unescape implements the Escaper interface
func (q GoEscaper) Unescape(s string) (string, int, error) {
	r, mb, t, err := strconv.UnquoteChar(s, byte(q))
	if err != nil {
		return "", 0, ErrEscape
	}
	if mb {
		return string(r), len(s) - len(t), nil
	}
	return string([]byte{byte(r)}), len(s) - len(t), nil
}
