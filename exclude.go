package aca

import (
	"unicode"
)

var (
	ExcludeControl = unicode.IsControl
	ExcludeSpace   = unicode.IsSpace
	ExcludePunct   = unicode.IsPunct
	ExcludeSymbol  = unicode.IsSymbol
)

func ExcludeNoneLetter(r rune) bool {
	return !unicode.IsLetter(r)
}

func ExcludeNoneDigit(r rune) bool {
	return !unicode.IsDigit(r)
}

func ExcludeNoneLetterOrDigit(r rune) bool {
	return !(unicode.IsLetter(r) || unicode.IsDigit(r))
}
