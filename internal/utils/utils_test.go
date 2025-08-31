package utils

import "testing"

func TestParseSize_EdgeCases(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
		err      bool
	}{
		{"0", 0, false},
		{"1g", 1073741824, false},
		{"1024mb", 1073741824, false},
		{"100kb", 102400, false},
		{"-100kb", 0, true},
		{"abc", 0, true},
		{"100xyz", 0, true},
	}

	for _, tt := range tests {
		size, err := ParseSize(tt.input)
		if (err != nil) != tt.err {
			t.Errorf("ParseSize(%q) error = %v, expected error: %v", tt.input, err, tt.err)
		}
		if size != tt.expected {
			t.Errorf("ParseSize(%q) = %d, expected %d", tt.input, size, tt.expected)
		}
	}
}
