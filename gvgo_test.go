// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gvgo

import (
	"math/rand/v2"
	"slices"
	"strings"
	"testing"
)

var tests = []struct {
	in       string
	stdout   string
	reparsed string
}{
	{"bad", "", "v0.0.0"},
	{"v1-alpha.beta.gamma", "", "v1.0.0"},
	{"v1-pre", "", "v1.0.0"},
	{"v1+meta", "", "v1.0.0"},
	{"v1-pre+meta", "", "v1.0.0"},
	{"v1.2-pre", "", "v1.2.0"},
	{"v1.2+meta", "", "v1.2.0"},
	{"v1.2-pre+meta", "", "v1.2.0"},
	{"v1.0.0-alpha", "v1.0.0-alpha", "v1.0.0-alpha"},
	{"v1.0.0-alpha.1", "v1.0.0-alpha.1", "v1.0.0-alpha.1"},
	{"v1.0.0-alpha.beta", "v1.0.0-alpha.beta", "v1.0.0-alpha.beta"},
	{"v1.0.0-beta", "v1.0.0-beta", "v1.0.0-beta"},
	{"v1.0.0-beta.2", "v1.0.0-beta.2", "v1.0.0-beta.2"},
	{"v1.0.0-beta.11", "v1.0.0-beta.11", "v1.0.0-beta.11"},
	{"v1.0.0-rc.1", "v1.0.0-rc.1", "v1.0.0-rc.1"},
	{"v1", "v1.0.0", "v1.0.0-0.0"},
	{"v1.0", "v1.0.0", "v1.0.0-0"},
	{"v1.0.0", "v1.0.0", "v1.0.0"},
	{"v1.2", "v1.2.0", "v1.2.0-0"},
	{"v1.2.0", "v1.2.0", "v1.2.0"},
	{"v1.2.3-456", "v1.2.3-456", "v1.2.3-456"},
	{"v1.2.3-456.789", "v1.2.3-456.789", "v1.2.3-456.789"},
	{"v1.2.3-456-789", "v1.2.3-456-789", "v1.2.3-456-789"},
	{"v1.2.3-456a", "v1.2.3-456a", "v1.2.3-456a"},
	{"v1.2.3-pre", "v1.2.3-pre", "v1.2.3-pre"},
	{"v1.2.3-pre+meta", "v1.2.3-pre", "v1.2.3-pre+meta"},
	{"v1.2.3-pre.1", "v1.2.3-pre.1", "v1.2.3-pre.1"},
	{"v1.2.3-zzz", "v1.2.3-zzz", "v1.2.3-zzz"},
	{"v1.2.3", "v1.2.3", "v1.2.3"},
	{"v1.2.3+meta", "v1.2.3", "v1.2.3+meta"},
	{"v1.2.3+meta-pre", "v1.2.3", "v1.2.3+meta-pre"},
	{"v1.2.3+meta-pre.sha.256a", "v1.2.3", "v1.2.3+meta-pre.sha.256a"},
}

func TestIsValid(t *testing.T) {
	for _, tt := range tests {
		ok := IsValid(tt.in)
		if ok != (tt.stdout != "") {
			t.Errorf("IsValid(%q) = %v, want %v", tt.in, ok, !ok)
		}
	}
}

func TestCanonicalString(t *testing.T) {
	for _, tt := range tests {
		out := CanonicalString(tt.in)
		if out != tt.stdout {
			t.Errorf("Canonical(%q) = %q, want %q", tt.in, out, tt.stdout)
		}
	}
}

func TestCompareString(t *testing.T) {
	for i, ti := range tests {
		for j, tj := range tests {
			cmp := CompareString(ti.in, tj.in)
			var want int
			if ti.stdout == tj.stdout {
				want = 0
			} else if i < j {
				want = -1
			} else {
				want = +1
			}
			if cmp != want {
				t.Errorf("Compare(%q, %q) = %d, want %d", ti.in, tj.in, cmp, want)
			}
		}
	}
}

func TestSort(t *testing.T) {
	var versions []Parsed
	for _, test := range tests {
		parsed, ok := New(test.in)
		if !ok {
			continue
		}
		versions = append(versions, parsed)
	}
	rand.Shuffle(len(versions), func(i, j int) { versions[i], versions[j] = versions[j], versions[i] })
	slices.SortFunc(versions, Compare)
	if !slices.IsSortedFunc(versions, Compare) {
		all := make([]string, 0, len(versions))
		for _, s := range all {
			all = append(all, s)
		}
		t.Errorf("list is not sorted:\n%s", strings.Join(all, "\n"))
	}
}

func TestString(t *testing.T) {
	for _, tt := range tests {
		s, _ := New(tt.in)
		if s.String() != tt.reparsed {
			t.Errorf("String(%q) = %q, want %q", tt.in, s.String(), tt.reparsed)
		}
	}
}

var (
	v1 = "v1.0.0+metadata-dash"
	v2 = "v1.0.0+metadata-dash1"
)

func BenchmarkCompare(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if CompareString(v1, v2) != 0 {
			b.Fatalf("bad compare")
		}
	}
}
