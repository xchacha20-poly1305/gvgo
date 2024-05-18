// Package gvgo used to parse version info.
package gvgo

import (
	"cmp"
	"errors"
	"strconv"
	"strings"
)

// Something from: https://github.com/golang/go/blob/go1.22.3/src/internal/gover/gover.go

// A Version is a parsed version: major[.Minor[.Patch]][-kind[(.)pre]]
// The numbers are the original decimal strings to avoid integer overflows
// and since there is very little actual math. (Probably overflow doesn't matter in practice,
// but at the time this code was written, there was an existing test that used
// go1.99999999999, which does not fit in an int on 32-bit platforms.
// The "big decimal" representation avoids the problem entirely.)
type Version struct {
	Major string // decimal
	Minor string // decimal or ""
	Patch string // decimal or ""
	Kind  string // "", "alpha", "beta", "rc"
	Pre   string // decimal or ""
}

// Parse parses the version string. Whether starts with "v" or not are both OK.
func Parse(x string) (Version, error) {
	var v Version

	x = strings.TrimPrefix(x, "v")
	parts := strings.Split(x, "-")

	// "0.0.0"
	mainParts := strings.Split(parts[0], ".")
	for i, mainPart := range mainParts {
		switch i {
		case 0:
			if mainPart == "" {
				return v, errors.New("main part empty")
			}
			v.Major = mainPart
		case 1:
			v.Minor = mainPart
		case 2:
			v.Patch = mainPart
		default:
			return v, errors.New("main part too long")
		}
	}

	if len(parts) == 3 {
		// "0.0.0-[timestamp]-[sha1]"
		v.Kind = parts[1]
		v.Pre = parts[2]
	} else {
		// "rc0" or "rc.0"
		extra := after(x, parts[0]+"-")
		extraIndex := strings.IndexFunc(extra, func(r rune) bool {
			s := string(r)
			if _, err := strconv.Atoi(s); err == nil {
				return true
			}
			return s == "."
		})
		if extraIndex >= 0 {
			if extraSlices := strings.Split(extra, ""); extraSlices[extraIndex] == "." {
				v.Kind = before(extra, ".")
				v.Pre = after(extra, ".")
			} else {
				v.Kind = strings.Join(extraSlices[:extraIndex], "")
				v.Pre = strings.Join(extraSlices[extraIndex:], "")
			}
		}
	}

	return v, nil
}

// String print readable version. But it will not include "v" at first.
func (v Version) String() string {
	s := v.Major
	if v.Minor != "" {
		s += "." + v.Minor
		if v.Patch != "" {
			s += "." + v.Patch

		} else {
			s += ".0"
		}
	} else {
		s += ".0.0"
	}

	if v.Kind != "" {
		s += "-" + v.Kind
		if v.Pre != "" {
			s += "." + v.Pre
		}
	}

	return s
}

// Compare compares two string. It not care about if they are valid.
func Compare(a, b string) int {
	av, _ := Parse(a)
	bv, _ := Parse(b)
	return CompareVersion(av, bv)
}

// CompareVersion compares two Version.
func CompareVersion(a, b Version) int {
	if c := compareInt(a.Major, b.Major); c != 0 {
		return c
	}
	if c := compareInt(a.Minor, b.Minor); c != 0 {
		return c
	}
	if c := compareInt(a.Patch, b.Patch); c != 0 {
		return c
	}
	if c := cmp.Compare(a.Kind, b.Kind); c != 0 {
		return c
	}
	if c := compareInt(a.Pre, b.Pre); c != 0 {
		return c
	}
	return 0
}

func (v Version) Compare(y Version) int {
	return CompareVersion(v, y)
}

// IsValid checks if a string is a valid Version.
func IsValid(v string) bool {
	_, err := Parse(v)
	return err == nil
}
