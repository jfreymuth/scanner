// Package scanner provides an easy-to-use scanner for writing simple DSL parsers.
package scanner // import "github.com/jfreymuth/scanner"

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"unicode"
	"unicode/utf8"
)

// A Scanner wraps an io.Reader and provides convenient methods for parsing simple languages.
// The scanner will ignore whitespace and go-style comments, except to seperate tokens.
// Once the scanner encounters any error, most of it's methods will return the zero value.
type Scanner struct {
	source  *bufio.Scanner
	line    string
	pos, ln int
	r       rune
	size    int
	err     error
}

// New creates a scanner that will read from an io.Reader
func New(in io.Reader) *Scanner {
	s := &Scanner{source: bufio.NewScanner(in)}
	s.space()
	return s
}

// FromString creates a scanner that will read from a string.
func FromString(s string) *Scanner {
	return New(strings.NewReader(s))
}

func (s *Scanner) next() {
	s.pos += s.size
	s.r, s.size = 0, 0
	if s.err != nil {
		return
	}
	if s.pos < len(s.line) {
		s.r, s.size = utf8.DecodeRuneInString(s.line[s.pos:])
	} else {
		if !s.source.Scan() {
			s.err = s.source.Err()
			if s.err == nil {
				s.err = io.EOF
			}
			s.line = ""
			s.pos = 0
			return
		}
		s.line = s.source.Text()
		s.ln++
		s.pos = 0
		s.r = '\n'
	}
}

func (s *Scanner) update() {
	s.size = 0
	s.next()
}

func (s *Scanner) space() {
repeat:
	s.update()
	for unicode.IsSpace(s.r) {
		s.next()
	}
	if s.Is("//") {
		s.pos = len(s.line)
		s.next()
		goto repeat
	} else if s.Is("/*") {
		err := &Error{"unmatched '/*'", s.line, s.ln, s.pos}
		s.pos += 2
		for s.err == nil {
			end := strings.Index(s.line[s.pos:], "*/")
			if end >= 0 {
				s.pos += end + 2
				goto repeat
			}
			s.pos = len(s.line)
			s.update()
		}
		if s.err == io.EOF {
			s.err = err
		}
	}
}

// Peek returns the next rune, without advancing the scanner.
// Peek will never return a whitespace character.
// If the scanner has encountered an error or EOF, Peek returns 0.
func (s *Scanner) Peek() rune {
	return s.r
}

// Rune returns the next rune and advances the scanner.
// Rune will never return a whitespace character.
// If the scanner has encountered an error or EOF, Rune returns 0.
func (s *Scanner) Rune() rune {
	if s.err != nil {
		return 0
	}
	r := s.r
	s.pos += s.size
	s.space()
	return r
}

// Is returns true if the scanner's input starts with str, but does not advance the scanner.
// str should not contain whitespace.
func (s *Scanner) Is(str string) bool {
	return strings.HasPrefix(s.line[s.pos:], str)
}

// Eat returns true and consumes the string if the scanner's input starts with str.
// The scanner will not be advanced if Eat returns false.
// str should not contain whitespace.
func (s *Scanner) Eat(str string) bool {
	if s.Is(str) {
		s.pos += len(str)
		s.space()
		return true
	}
	return false
}

// Demand consumes a string, or causes an error if the scanner's input does not start with str.
// str should not contain whitespace.
func (s *Scanner) Demand(str string) {
	if !s.Eat(str) {
		s.Failf("'%s' expected", str)
	}
}

// Fail sets the scanner's error, unless it has already encountered another error.
// The error will contain the current line and position of the scanner.
func (s *Scanner) Fail(msg string) {
	if s.err == nil || s.err == io.EOF {
		s.err = &Error{msg, s.line, s.ln, s.pos}
		s.line = ""
		s.pos = 0
		s.r = 0
	}
}

// Failf sets the scanner's error, unless it has already encountered another error.
// The error will contain the current line and position of the scanner.
func (s *Scanner) Failf(format string, a ...interface{}) {
	if s.err == nil || s.err == io.EOF {
		s.err = &Error{fmt.Sprintf(format, a...), s.line, s.ln, s.pos}
		s.line = ""
		s.pos = 0
		s.r = 0
	}
}

// Err returns the first error encountered by the scanner, or nil if there was no error.
// The returned error will either be of the type *Error, or an error returned by the underlying io.Reader.
func (s *Scanner) Err() error {
	if s.err == io.EOF {
		return nil
	}
	return s.err
}

// End returns true if the scanner has reached the end of the input or encountered an error.
func (s *Scanner) End() bool {
	return s.err != nil
}

// Error is an error type that contains a reference to a specific position in a scanners input.
type Error struct {
	Message  string
	Line     string
	LineNum  int
	Position int
}

// Error returns the errors message.
func (e *Error) Error() string {
	return "scanner: " + e.Message
}

// PositionIndicator returns a user-friendly multi-line-string containing the message, line and position of the error.
func (e *Error) PositionIndicator() string {
	if e.Position > 0 {
		// TODO account for tabs
		return fmt.Sprint("Line ", e.LineNum, ": ", e.Message, "\n", e.Line, "\n", strings.Repeat(" ", e.Position-1), "^")
	}
	return fmt.Sprint("Line ", e.LineNum, ": ", e.Message, "\n", e.Line, "\n", "^")
}
