// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package gvgo implements comparison of semantic version strings.
// In this package, semantic version strings must begin with a leading "v",
// as in "v1.0.0".
//
// The general form of a semantic version string accepted by this package is
//
//	vMAJOR[.MINOR[.PATCH[-PRERELEASE][+BUILD]]]
//
// where square brackets indicate optional parts of the syntax;
// MAJOR, MINOR, and PATCH are decimal integers without extra leading zeros;
// PRERELEASE and BUILD are each a series of non-empty dot-separated identifiers
// using only alphanumeric characters and hyphens; and
// all-numeric PRERELEASE identifiers must not have leading zeros.
//
// This package follows Semantic Versioning 2.0.0 (see semver.org)
// with two exceptions. First, it requires the "v" prefix. Second, it recognizes
// vMAJOR and vMAJOR.MINOR (with no prerelease or build suffixes)
// as shorthands for vMAJOR.0.0 and vMAJOR.MINOR.0.
package gvgo

import (
	"strings"
)

// Parsed returns the parsed form of a semantic version string.
type Parsed struct {
	Major      string
	Minor      string
	Patch      string
	Short      string // Starts with "."
	Prerelease string // Starts with "-"
	Build      string // Starts with "+"
}

// IsValid reports whether v is a valid semantic version string.
func IsValid(v string) bool {
	_, ok := Parse(v)
	return ok
}

// CompareString is same as Compare.
//
// An invalid semantic version string is considered less than a valid one.
// All invalid semantic version strings compare equal to each other.
func CompareString(v, w string) int {
	pv, ok1 := Parse(v)
	pw, ok2 := Parse(w)
	if !ok1 && !ok2 {
		return 0
	}
	if !ok1 {
		return -1
	}
	if !ok2 {
		return +1
	}
	return Compare(pv, pw)
}

func (p Parsed) Compare(v Parsed) int {
	return Compare(p, v)
}

// Compare returns an integer comparing two versions according to
// semantic version precedence.
// The result will be 0 if v == w, -1 if v < w, or +1 if v > w.
func Compare(v, w Parsed) int {
	if c := compareInt(v.Major, w.Major); c != 0 {
		return c
	}
	if c := compareInt(v.Minor, w.Minor); c != 0 {
		return c
	}
	if c := compareInt(v.Patch, w.Patch); c != 0 {
		return c
	}
	return comparePrerelease(v.Prerelease, w.Prerelease)
}

// Parse parses a new parsed version, which starts with "v".
func Parse(v string) (p Parsed, ok bool) {
	if v == "" || v[0] != 'v' {
		return
	}
	p.Major, v, ok = parseInt(v[1:])
	if !ok {
		return
	}
	if v == "" {
		p.Minor = "0"
		p.Patch = "0"
		return
	}
	if v[0] != '.' {
		ok = false
		return
	}
	p.Minor, v, ok = parseInt(v[1:])
	if !ok {
		return
	}
	if v == "" {
		p.Patch = "0"
		return
	}
	if v[0] != '.' {
		ok = false
		return
	}
	p.Patch, v, ok = parseInt(v[1:])
	if !ok {
		return
	}
	if len(v) > 0 && v[0] == '-' {
		p.Prerelease, v, ok = parsePrerelease(v)
		if !ok {
			return
		}
	}
	if len(v) > 0 && v[0] == '+' {
		p.Build, v, ok = parseBuild(v)
		if !ok {
			return
		}
	}
	if v != "" {
		ok = false
		return
	}
	ok = true
	return
}

func parseInt(v string) (t, rest string, ok bool) {
	if v == "" {
		return
	}
	if v[0] < '0' || '9' < v[0] {
		return
	}
	i := 1
	for i < len(v) && '0' <= v[i] && v[i] <= '9' {
		i++
	}
	if v[0] == '0' && i != 1 {
		return
	}
	return v[:i], v[i:], true
}

func parsePrerelease(v string) (t, rest string, ok bool) {
	// "A pre-release version MAY be denoted by appending a hyphen and
	// a series of dot separated identifiers immediately following the Patch version.
	// Identifiers MUST comprise only ASCII alphanumerics and hyphen [0-9A-Za-z-].
	// Identifiers MUST NOT be empty. Numeric identifiers MUST NOT include leading zeroes."
	if v == "" || v[0] != '-' {
		return
	}
	i := 1
	start := 1
	for i < len(v) && v[i] != '+' {
		if !isIdentChar(v[i]) && v[i] != '.' {
			return
		}
		if v[i] == '.' {
			if start == i || isBadNum(v[start:i]) {
				return
			}
			start = i + 1
		}
		i++
	}
	if start == i || isBadNum(v[start:i]) {
		return
	}
	return v[:i], v[i:], true
}

func parseBuild(v string) (t, rest string, ok bool) {
	if v == "" || v[0] != '+' {
		return
	}
	i := 1
	start := 1
	for i < len(v) {
		if !isIdentChar(v[i]) && v[i] != '.' {
			return
		}
		if v[i] == '.' {
			if start == i {
				return
			}
			start = i + 1
		}
		i++
	}
	if start == i {
		return
	}
	return v[:i], v[i:], true
}

func isIdentChar(c byte) bool {
	return 'A' <= c && c <= 'Z' || 'a' <= c && c <= 'z' || '0' <= c && c <= '9' || c == '-'
}

func isBadNum(v string) bool {
	i := 0
	for i < len(v) && '0' <= v[i] && v[i] <= '9' {
		i++
	}
	return i == len(v) && i > 1 && v[0] == '0'
}

func isNum(v string) bool {
	i := 0
	for i < len(v) && '0' <= v[i] && v[i] <= '9' {
		i++
	}
	return i == len(v)
}

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

func comparePrerelease(x, y string) int {
	// "When major, minor, and Patch are equal, a pre-release version has
	// lower precedence than a normal version.
	// Example: 1.0.0-alpha < 1.0.0.
	// Precedence for two pre-release versions with the same major, minor,
	// and Patch version MUST be determined by comparing each dot separated
	// identifier from left to right until a difference is found as follows:
	// identifiers consisting of only digits are compared numerically and
	// identifiers with letters or hyphens are compared lexically in ASCII
	// sort order. Numeric identifiers always have lower precedence than
	// non-numeric identifiers. A larger set of pre-release fields has a
	// higher precedence than a smaller set, if all of the preceding
	// identifiers are equal.
	// Example: 1.0.0-alpha < 1.0.0-alpha.1 < 1.0.0-alpha.beta <
	// 1.0.0-beta < 1.0.0-beta.2 < 1.0.0-beta.11 < 1.0.0-rc.1 < 1.0.0."
	if x == y {
		return 0
	}
	if x == "" {
		return +1
	}
	if y == "" {
		return -1
	}
	for x != "" && y != "" {
		x = x[1:] // skip - or .
		y = y[1:] // skip - or .
		var dx, dy string
		dx, x = nextIdent(x)
		dy, y = nextIdent(y)
		if dx != dy {
			ix := isNum(dx)
			iy := isNum(dy)
			if ix != iy {
				if ix {
					return -1
				} else {
					return +1
				}
			}
			if ix {
				if len(dx) < len(dy) {
					return -1
				}
				if len(dx) > len(dy) {
					return +1
				}
			}
			if dx < dy {
				return -1
			} else {
				return +1
			}
		}
	}
	if x == "" {
		return -1
	} else {
		return +1
	}
}

func nextIdent(x string) (dx, rest string) {
	i := 0
	for i < len(x) && x[i] != '.' {
		i++
	}
	return x[:i], x[i:]
}

// String returns formated version, which starts without "v".
func (p Parsed) String() (s string) {
	appendOrDefault := func(v *string) {
		if *v == "" {
			s += "0"
		} else {
			s += *v
		}
	}
	appendOrDefault(&p.Major)
	s += "."
	appendOrDefault(&p.Minor)
	s += "."
	appendOrDefault(&p.Patch)
	if p.Short != "" {
		s += "-" + strings.TrimLeft(p.Short, ".")
	}
	if p.Prerelease != "" {
		s += p.Prerelease
	}
	if p.Build != "" {
		s += p.Build
	}
	return
}

func (p Parsed) IsPrerelease() bool {
	return strings.TrimLeft(p.Prerelease, "-") != ""
}

func (p Parsed) IsBuild() bool {
	return strings.TrimLeft(p.Build, "+") != ""
}
