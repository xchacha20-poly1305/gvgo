package gvgo

import (
	"fmt"
	"testing"
)

type cutCase struct {
	input    string
	flag     string
	expected string
}

func TestAfter(t *testing.T) {
	tests := []cutCase{
		{
			input:    "a-b",
			flag:     "-",
			expected: "b",
		},
		{
			input:    "s-n-a-c-k",
			flag:     "-",
			expected: "n-a-c-k",
		},
	}

	for _, tt := range tests {
		s := after(tt.input, tt.flag)
		printed := fmt.Sprintf("%s: want %s, got %s", tt.input, s, tt.expected)
		if s != tt.expected {
			t.Error(printed)
			continue
		}
		t.Logf(printed)
	}
}
