package gvgo

import (
	"cmp"
	"strings"
)

func before(s string, b string) string {
	index := strings.Index(s, b)
	if index < 0 {
		return s
	}
	return s[:index]
}

func after(s string, a string) string {
	index := strings.Index(s, a)
	if index < 0 {
		return s
	}
	return s[index+len(a):]
}

// compareInt returns cmp.Compare(x, y) interpreting x and y as decimal numbers.
// (Copied from golang.org/x/mod/semver's cmpInt.)
func compareInt(x, y string) int {
	if x == y {
		return 0
	}
	if len(x) < len(y) {
		return -1
	}
	if len(x) > len(y) {
		return +1
	}
	if x < y {
		return -1
	} else {
		return +1
	}
}

const (
	KindAlpha = "alpha"
	KindBeta  = "beta"
	KindRC    = "rc"
)

func compareKind(x, y string) int {
	x = strings.ToLower(x)
	y = strings.ToLower(y)

	xIsValidKind := x == KindAlpha || x == KindBeta || x == KindRC
	yIsValidKind := y == KindAlpha || y == KindBeta || y == KindRC

	if xIsValidKind && yIsValidKind {
		// "" < alpha < beta < rc
		return cmp.Compare(x, y)
	}
	if xIsValidKind && !yIsValidKind {
		return 1
	}
	if !xIsValidKind && yIsValidKind {
		return -1
	}

	// timestamp
	return compareInt(x, y)
}
