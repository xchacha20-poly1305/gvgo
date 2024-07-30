package gvgo

import (
	"testing"
)

func TestValidKind(t *testing.T) {
	tests := []struct {
		name string
		kind string
		want bool
	}{
		{
			name: "alpha",
			kind: KindAlpha,
			want: true,
		},
		{
			name: "beta",
			kind: KindBeta,
			want: true,
		},
		{
			name: "rc",
			kind: KindRc,
			want: true,
		},
		{
			name: "invalid",
			kind: "invalid",
			want: false,
		},
		{
			name: "empty",
			kind: "",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidKind(tt.kind); got != tt.want {
				t.Errorf("ValidKind() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVersion_String(t *testing.T) {
	tests := []struct {
		name    string
		version Version
		want    string
	}{
		{
			name:    "normal",
			version: New(),
			want:    "0.0.0",
		},
		{
			name: "pre just kind",
			version: Version{
				Major: "1",
				Minor: "6",
				Patch: "3",
				Kind:  KindRc,
			},
			want: "1.6.3-rc",
		},
		{
			name: "pre",
			version: Version{
				Major: "2",
				Minor: "5",
				Patch: "6",
				Kind:  KindBeta,
				Pre:   "9",
			},
			want: "2.5.6-beta.9",
		},
		{
			name: "pre without kind",
			version: Version{
				Major: "5",
				Minor: "1",
				Patch: "2",
				Pre:   "0",
			},
			want: "5.1.2",
		},
		{
			name: "totally not released",
			version: Version{
				Major:   "0",
				Minor:   "0",
				Patch:   "0",
				GitInfo: "20240719175910-8a7402abbf56",
			},
			want: "0.0.0-20240719175910-8a7402abbf56",
		},
		{
			name: "not release next version",
			version: Version{
				Major:         "0",
				Minor:         "6",
				Patch:         "2",
				BuildMetadata: "0",
				GitInfo:       "20240717063648-d3b0c53281a1",
			},
			want: "0.6.2-0.20240717063648-d3b0c53281a1",
		},
		{
			name: "pre with pseudo",
			version: Version{
				Major:         "1",
				Minor:         "10",
				Patch:         "0",
				Kind:          KindAlpha,
				Pre:           "26",
				BuildMetadata: "0",
				GitInfo:       "20240727034746-0efc42a5ef8d",
			},
			want: "1.10.0-alpha.26.0.20240727034746-0efc42a5ef8d",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotVersion := tt.version.String(); gotVersion != tt.want {
				t.Errorf("String() = %v, want %v", gotVersion, tt.want)
			}
		})
	}
}
