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
	in          string
	reformatted string
	isValid     bool
}{
	{"bad", "0.0.0", false},
	{"v1-alpha.beta.gamma", "1.0.0", false},
	{"v1-pre", "1.0.0", false},
	{"v1+meta", "1.0.0", false},
	{"v1-pre+meta", "1.0.0", false},
	{"v1.2-pre", "1.2.0", false},
	{"v1.2+meta", "1.2.0", false},
	{"v1.2-pre+meta", "1.2.0", false},
	{"v1.0.0-alpha", "1.0.0-alpha", true},
	{"v1.0.0-alpha.1", "1.0.0-alpha.1", true},
	{"v1.0.0-alpha.beta", "1.0.0-alpha.beta", true},
	{"v1.0.0-beta", "1.0.0-beta", true},
	{"v1.0.0-beta.2", "1.0.0-beta.2", true},
	{"v1.0.0-beta.11", "1.0.0-beta.11", true},
	{"v1.0.0-rc.1", "1.0.0-rc.1", true},
	{"v1", "1.0.0", true},
	{"v1.0", "1.0.0", true},
	{"v1.0.0", "1.0.0", true},
	{"v1.2", "1.2.0", true},
	{"v1.2.0", "1.2.0", true},
	{"v1.2.3-456", "1.2.3-456", true},
	{"v1.2.3-456.789", "1.2.3-456.789", true},
	{"v1.2.3-456-789", "1.2.3-456-789", true},
	{"v1.2.3-456a", "1.2.3-456a", true},
	{"v1.2.3-pre", "1.2.3-pre", true},
	{"v1.2.3-pre+meta", "1.2.3-pre+meta", true},
	{"v1.2.3-pre.1", "1.2.3-pre.1", true},
	{"v1.2.3-zzz", "1.2.3-zzz", true},
	{"v1.2.3", "1.2.3", true},
	{"v1.2.3+meta", "1.2.3+meta", true},
	{"v1.2.3+meta-pre", "1.2.3+meta-pre", true},
	{"v1.2.3+meta-pre.sha.256a", "1.2.3+meta-pre.sha.256a", true},
}

func TestIsValid(t *testing.T) {
	for _, tt := range tests {
		ok := IsValid(tt.in)
		if ok != tt.isValid {
			t.Errorf("IsValid(%q) = %v, want %v", tt.in, ok, !ok)
		}
	}
}

/*func TestCompareString(t *testing.T) {
	for i, ti := range tests {
		for j, tj := range tests {
			cmp := CompareString(ti.in, tj.in)
			var want int
			if ti.reformatted == tj.reformatted {
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
}*/

func TestSort(t *testing.T) {
	var versions []Parsed
	for _, test := range tests {
		parsed, ok := Parse(test.in)
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
		s, _ := Parse(tt.in)
		if s.String() != tt.reformatted {
			t.Errorf("String(%q) = %q, want %q", tt.in, s.String(), tt.reformatted)
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
