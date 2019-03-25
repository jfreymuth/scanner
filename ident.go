package scanner

import (
	"unicode"
)

// IsIdent returns true if the next token is an identifier.
func (s *Scanner) IsIdent() bool {
	return s.r == '_' || unicode.IsLetter(s.r)
}

// PeekIdent returns the next token, or an empty string if the next token is not an identifier.
// PeekIdent will not advance the scanner.
func (s *Scanner) PeekIdent() string {
	if !s.IsIdent() {
		return ""
	}
	for i, r := range s.line[s.pos:] {
		if r != '_' && !unicode.IsLetter(r) && !unicode.IsNumber(r) {
			return s.line[s.pos : s.pos+i]
		}
	}
	result := s.line[s.pos:]
	return result
}

// Ident returns the next token, or causes an error if the next token is not an identifier.
func (s *Scanner) Ident() string {
	result := s.PeekIdent()
	if result == "" {
		s.Fail("identifier expected")
	}
	s.pos += len(result)
	s.space()
	return result
}
