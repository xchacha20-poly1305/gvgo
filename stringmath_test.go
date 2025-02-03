package gvgo

import (
	"testing"
)

func TestPlus(t *testing.T) {
	testCases := []struct {
		source string
		extra  string
		expect string
	}{
		// Test cases with valid numbers
		{"123", "456", "579"},
		{"999", "1", "1000"},
		{"0", "0", "0"},
		{"1", "99", "100"},
		{"99", "1", "100"},
		{"", "", ""},
		{"12345", "67890", "80235"},

		// Test cases with non-numeric input
		{"abc", "123", "abc"},
		{"123", "abc", "123"},
		{"a1b2c", "456", "a1b2c"},
		{"123", "4a5b6", "123"},
		{"", "1", "1"}, // One empty string with valid number.
		{"1", "", "1"},
		{"", "abc", ""},
		{"abc", "", "abc"},

		// Test cases with leading zeros
		{"001", "002", "001"},
		{"0123", "456", "0123"},
		{"0", "123", "123"},
		{"123", "0", "123"},
	}

	for _, tc := range testCases {
		t.Run(tc.source+"+"+tc.extra, func(t *testing.T) {
			result := Plus(tc.source, tc.extra)
			if result != tc.expect {
				t.Errorf("Plus(%s, %s) = %s; want %s", tc.source, tc.extra, result, tc.expect)
			}
		})
	}
}
