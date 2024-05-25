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

// IsValidKind returns true if kind is one of kinds("alpha", "beta" and "rc").
func IsValidKind(kind string) bool {
	return kind == KindAlpha || kind == KindBeta || kind == KindRC
}

func compareKind(x, y string) int {
	x = strings.ToLower(x)
	y = strings.ToLower(y)

	xIsValidKind := IsValidKind(x)
	yIsValidKind := IsValidKind(y)

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
