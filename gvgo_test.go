package gvgo

import (
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	want := Version{
		Major: "0",
		Minor: "0",
		Patch: "0",
	}
	if got := New(); !reflect.DeepEqual(got, want) {
		t.Errorf("New() = %v, want %v", got, want)
	}
}

func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		raw     string
		wantV   Version
		wantErr bool
	}{
		{
			name: "normal",
			raw:  "v1.2.3",
			wantV: Version{
				Major: "1",
				Minor: "2",
				Patch: "3",
			},
			wantErr: false,
		},
		{
			name: "not has v",
			raw:  "1.0.0",
			wantV: Version{
				Major: "1",
				Minor: "0",
				Patch: "0",
			},
			wantErr: false,
		},
		{
			name:    "empty",
			raw:     "",
			wantErr: true,
		},
		{
			name:    "invalid",
			raw:     "gvgo!",
			wantErr: true,
		},
		{
			name: "pre",
			raw:  "v1.23.0-rc.2",
			wantV: Version{
				Major: "1",
				Minor: "23",
				Patch: "0",
				Kind:  KindRc,
				Pre:   "2",
			},
			wantErr: false,
		},
		{
			name: "pre just kind",
			raw:  "0.8.2-rc",
			wantV: Version{
				Major: "0",
				Minor: "8",
				Patch: "2",
				Kind:  KindRc,
			},
			wantErr: false,
		},
		{
			name: "git",
			raw:  "v0.6.2-0.20240717063648-d3b0c53281a1",
			wantV: Version{
				Major:   "0",
				Minor:   "6",
				Patch:   "2",
				GitInfo: "20240717063648-d3b0c53281a1",
			},
			wantErr: false,
		},
		{
			name: "pre + git",
			raw:  "v1.10.0-alpha.26.0.20240727034746-0efc42a5ef8d",
			wantV: Version{
				Major:         "1",
				Minor:         "10",
				Patch:         "0",
				Kind:          KindAlpha,
				Pre:           "26",
				BuildMetadata: "0",
				GitInfo:       "20240727034746-0efc42a5ef8d",
			},
			wantErr: false,
		},
		{
			name: "too big",
			raw:  "1.99999999999",
			wantV: Version{
				Major: "1",
				Minor: "99999999999",
				Patch: "0",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotV, err := Parse(tt.raw)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && !reflect.DeepEqual(gotV, tt.wantV) {
				t.Errorf("Parse() gotV = %v, want %v", gotV, tt.wantV)
			}
		})
	}
}

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

func TestCompare(t *testing.T) {
	tests := []struct {
		name   string
		v1, v2 Version
		want   int
	}{
		{
			name: "Just main",
			v1: Version{
				Major: "1",
				Minor: "2",
				Patch: "3",
			},
			v2: Version{
				Major: "1",
				Minor: "2",
				Patch: "3",
			},
			want: 0,
		},
		{
			name: "two kind",
			v1: Version{
				Major: "1",
				Minor: "2",
				Patch: "3",
				Kind:  KindAlpha,
			},
			v2: Version{
				Major: "1",
				Minor: "2",
				Patch: "3",
				Kind:  KindBeta,
			},
			want: -1,
		},
		{
			name: "kind with normal",
			v1: Version{
				Major: "1",
				Minor: "2",
				Patch: "3",
				Kind:  KindAlpha,
			},
			v2: Version{
				Major: "1",
				Minor: "2",
				Patch: "3",
			},
			want: -1,
		},
		{
			name: "with git",
			v1: Version{
				Major:   "1",
				Minor:   "2",
				Patch:   "3",
				GitInfo: "20240719175910-8a7402abbf56",
			},
			v2: Version{
				Major:   "1",
				Minor:   "2",
				Patch:   "3",
				GitInfo: "20240717063648-d3b0c53281a1",
			},
			want: 1,
		},
		{
			name: "pre + git",
			v1: Version{
				Major:         "1",
				Minor:         "10",
				Patch:         "0",
				Kind:          KindAlpha,
				Pre:           "0",
				BuildMetadata: "0",
				GitInfo:       "20240727034746-0efc42a5ef8d",
			},
			v2: Version{
				Major:         "1",
				Minor:         "10",
				Patch:         "0",
				Kind:          KindAlpha,
				Pre:           "0",
				BuildMetadata: "0",
				GitInfo:       "20240727034745-bf12e2370b4c",
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Compare(tt.v1, tt.v2); got != tt.want {
				t.Errorf("Compare() = %v, want %v", got, tt.want)
			}
		})
	}
}
