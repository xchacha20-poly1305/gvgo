package gvgo

import (
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		str     string
		wantErr bool
	}{
		{
			name:    "Alpha",
			str:     "v5.2.0-alpha.0",
			wantErr: false,
		},
		{
			name:    "beta",
			str:     "v1.3.14-beta.1",
			wantErr: false,
		},
		{
			name:    "rc",
			str:     "v8.8.6-rc.2",
			wantErr: false,
		},
		{
			name:    "No point kind",
			str:     "v4.9.8-rc0",
			wantErr: false,
		},
		{
			name:    "No \"v\"",
			str:     "2.4.3",
			wantErr: false,
		},
		{
			name:    "Short 1",
			str:     "0",
			wantErr: false,
		},
		{
			name:    "Short 2",
			str:     "1.2",
			wantErr: false,
		},
		{
			name:    "Short 3",
			str:     "3.4.5",
			wantErr: false,
		},
		{
			name:    "Git",
			str:     "v0.0.0-20240506185415-9bf2ced13842",
			wantErr: false,
		},
		{
			name:    "Too long",
			str:     "v1.1.1.1",
			wantErr: true,
		},
		{
			name:    "Empty",
			str:     "",
			wantErr: true,
		},
		{
			name:    "Big",
			str:     "v9999999999999999999999.99999999999999999999999999.9999999999999999999999",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		v, err := Parse(tt.str)
		if tt.wantErr {
			if err == nil {
				t.Errorf("%s wants error but got: %s", tt.name, v.String())
			}
			continue
		}
		if err != nil {
			t.Errorf("Failed to parse [%s]: %v", tt.name, err)
			continue
		}
		t.Logf("Success [%s]: %s", tt.name, v.String())
	}
}
