// Package gvgo implements version compare for Go module version.
//
// Rules detail: https://go.dev/doc/modules/version-numbers
//
// https://go.dev/ref/mod#pseudo-versions
package gvgo

import (
	"strings"
)

// A Version is a parsed Go version: major[.Minor[.Patch]][kind[pre]][.BuildMetadata.Timestamp-Commit]
// The numbers are the original decimal strings to avoid integer overflows
// and since there is very little actual math. (Probably overflow doesn't matter in practice,
// but at the time this code was written, there was an existing test that used
// go1.99999999999, which does not fit in an int on 32-bit platforms.
// The "big decimal" representation avoids the problem entirely.)
type Version struct {
	Major string // decimal
	Minor string // decimal
	Patch string // decimal
	Kind  string // "", "alpha", "beta", "rc"
	Pre   string // decimal or ""

	// For pseudo-version. We assume has Commit means use pseudo-version.
	BuildMetadata string // decimal or ""
	GitInfo       string // timestamp-commit_hash
}

func New() Version {
	return Version{
		Major: "0",
		Minor: "0",
		Patch: "0",
	}
}

// Parse parses version.
func Parse(raw string) (v Version, err error) {
	raw = strings.TrimPrefix(raw, "v")

	var mainPart [2]string
	mainPart[0], mainPart[1], _ = strings.Cut(raw, "-")
	v.Major, v.Minor, v.Patch, err = parseMainPart(mainPart[0])
	if err != nil {
		return Version{}, err
	}
	if mainPart[1] == "" {
		// Not pre-release or pseudo.
		return
	}

	var afterKind string
	switch mainPart[1][0] {
	case KindAlpha[0]:
		var found bool
		afterKind, found = strings.CutPrefix(mainPart[1], KindAlpha)
		if !found {
			return Version{}, ErrInvalidKind
		}
		v.Kind = KindAlpha
	case KindBeta[0]:
		var found bool
		afterKind, found = strings.CutPrefix(mainPart[1], KindBeta)
		if !found {
			return Version{}, ErrInvalidKind
		}
		v.Kind = KindBeta
	case KindRc[0]:
		var found bool
		afterKind, found = strings.CutPrefix(mainPart[1], KindRc)
		if !found {
			return Version{}, ErrInvalidKind
		}
		v.Kind = KindRc
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		// Not pre-release
		v.BuildMetadata, v.GitInfo, err = parsePseudo(mainPart[1])
		if err != nil {
			return Version{}, err
		}
		return
	default:
		return Version{}, ErrInvalidKind
	}
	v.BuildMetadata, v.GitInfo, err = parsePseudo(afterKind)
	if err != nil {
		return Version{}, err
	}
	return
}

// parseMainPart parses the main part of version like 1.0.0.
func parseMainPart(raw string) (major, minor, patch string, err error) {
	var success bool
	major, raw, success = cutInt(raw)
	if !success {
		return "", "", "", Error{"read major version"}
	}
	minor, raw, success = cutInt(strings.TrimPrefix(raw, "."))
	if !success {
		return major, "0", "0", nil
	}
	patch, raw, success = cutInt(strings.TrimPrefix(raw, "."))
	if !success {
		patch = "0"
	}
	return
}

// parsePseudo parse pseudo parts like "0.20240717063648-d3b0c53281a1" or "20240719175910-8a7402abbf56"
func parsePseudo(raw string) (buildMetadata, gitInfo string, err error) {
	num, rest, success := cutInt(raw)
	if !success {
		return "", "", ErrInvalidGit
	}
	if strings.HasPrefix(rest, "-") {
		// Pure git
		// "20240719175910-8a7402abbf56"
		return "", raw, nil
	}
	// "0.20240717063648-d3b0c53281a1"
	return num, strings.TrimPrefix(rest, "."), nil
}

// Pre-release kind.
const (
	KindAlpha = "alpha"
	KindBeta  = "beta"
	KindRc    = "rc"
)

// ValidKind returns true when kind is valid.
func ValidKind(kind string) bool {
	return kind == KindAlpha || kind == KindBeta || kind == KindRc
}

// String returns the readable string of version.
// It not starts with "v".
func (v Version) String() (version string) {
	version = v.Major + "." + v.Minor + "." + v.Patch
	isPre := v.Kind != ""
	hasGit := v.GitInfo != ""
	if !isPre && !hasGit {
		return
	}

	version += "-"
	if isPre {
		version += v.Kind
		if v.Pre != "" {
			version += "." + v.Pre
		}
	}
	if hasGit {
		if v.BuildMetadata != "" {
			if !strings.HasSuffix(version, "-") {
				version += "."
			}
			version += v.BuildMetadata
		}
		if !strings.HasSuffix(version, "-") {
			version += "."
		}
		version += v.GitInfo
	}
	return
}
