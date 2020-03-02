package main

import (
	"github.com/kjk/notionapi"
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

func toDashIDs(ids []string) {
	for i := range ids {
		ids[i] = notionapi.ToDashID(ids[i])
	}
}

func findInBButNotInA(A []string, B []string) []string {
	result := make([]string, 0)
	m := make(map[string]struct{}, len(A))
	for _, a := range A {
		m[a] = struct{}{}
	}
	for _, b := range B {
		if _, ok := m[b]; !ok {
			result = append(result, b)
		}
	}
	return result
}
