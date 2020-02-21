package main

import (
	"strings"
	"unicode"
)

func trimAll(s string) string {
	var b strings.Builder
	b.Grow(len(s))

	for _, c := range s {
		if c != ' ' {
			b.WriteRune(c)
		}
	}
	return b.String()
}

func trimAndToSmall(s string) string {
	var b strings.Builder
	b.Grow(len(s))

	for _, c := range s {
		if !unicode.IsSpace(c) {
			b.WriteRune(unicode.ToLower(c))
		}
	}
	return b.String()
}

func trimAndConvertSpace(s string) string {
	var b strings.Builder
	b.Grow(len(s))

	for _, c := range s {
		if !unicode.IsSpace(c) {
			b.WriteRune(unicode.ToLower(c))
		} else {
			b.WriteByte('_')
		}
	}
	return b.String()
}
