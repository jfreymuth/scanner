package scanner

import (
	"strconv"
	"unicode"
)

// IsInt returns true if the next token is an int.
func (s *Scanner) IsInt() bool {
	_, ok := s.PeekInt()
	return ok
}

// PeekInt returns the next token as an int, or (0, false) if the next token is not an int.
// PeekInt does not advance the scanner.
func (s *Scanner) PeekInt() (int, bool) {
	i, n := peekInt(s.line[s.pos:])
	return i, n > 0
}

// Int returns the next token as an int, or causes an error if the next token is not an int.
func (s *Scanner) Int() int {
	i, n := peekInt(s.line[s.pos:])
	if n == 0 {
		s.Fail("integer expected")
		return 0
	}
	s.pos += n
	s.space()
	return i
}

// IsFloat returns true if the next token is a float.
func (s *Scanner) IsFloat() bool {
	_, ok := s.PeekFloat()
	return ok
}

// PeekFloat returns the next token as a float, or (0, false) if the next token is not a float.
// PeekFloat does not advance the scanner.
func (s *Scanner) PeekFloat() (float64, bool) {
	i, n := peekInt(s.line[s.pos:])
	if n != 0 {
		return float64(i), true
	}
	f, n := peekFloat(s.line[s.pos:])
	if n != 0 {
		return f, true
	}
	return 0, false
}

// Float returns the next token as a float, or causes an error if the next token is not a float.
func (s *Scanner) Float() float64 {
	i, n := peekInt(s.line[s.pos:])
	if n != 0 {
		s.pos += n
		s.space()
		return float64(i)
	}
	f, n := peekFloat(s.line[s.pos:])
	if n != 0 {
		s.pos += n
		s.space()
		return f
	}
	s.Fail("float expected")
	return 0
}

func peekInt(line string) (i, n int) {
	l := len(line)
	for i, r := range line {
		if (i != 0 || r != '-') && r != '.' && !unicode.IsDigit(r) && !unicode.IsLetter(r) {
			l = i
			break
		}
	}

	i64, err := strconv.ParseInt(line[:l], 0, strconv.IntSize)
	if err != nil {
		return 0, 0
	}
	return int(i64), l
}

func peekFloat(line string) (f float64, n int) {
	l := len(line)
	e := false
	for i, r := range line {
		if !(i == 0 && r == '-') && !(e && (r == '-' || r == '+')) && r != '.' && !unicode.IsDigit(r) && !unicode.IsLetter(r) {
			l = i
			break
		}
		e = false
		if r == 'e' || r == 'E' {
			e = true
		}
	}

	f, err := strconv.ParseFloat(line[:l], 64)
	if err != nil {
		return 0, 0
	}
	return f, l
}
