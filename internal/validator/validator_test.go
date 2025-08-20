package validator

import (
	"strings"
	"testing"
)

func TestIsValidPath_ValidNames(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"Simple filename", "file.txt", true},
		{"Directory with slash", "folder/", true},
		{"Directory without slash", "folder", true},
		{"Nested path", "path/to/file", true},
		{"File with spaces", "my file.txt", true},
		{"File with dots", "file.v1.2.3.txt", true},
		{"File with underscores", "file_name.txt", true},
		{"File with hyphens", "file-name.txt", true},
		{"File with numbers", "file123.txt", true},
		{"Mixed case", "FileNaMe.TxT", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidPath(tt.input)
			if result != tt.expected {
				t.Errorf("IsValidPath(%q) = %v, expected %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestIsValidPath_InvalidCharacters(t *testing.T) {
	invalidChars := []string{"\\", ":", "*", "?", "\"", "<", ">", "|"}

	for _, char := range invalidChars {
		t.Run("Invalid character "+char, func(t *testing.T) {
			testName := "file" + char + "name.txt"
			result := IsValidPath(testName)
			if result {
				t.Errorf("IsValidPath(%q) = true, expected false", testName)
			}
		})
	}
}

func TestIsValidPath_WindowsReservedNames(t *testing.T) {
	reservedNames := []string{
		"CON", "PRN", "AUX", "NUL",
		"COM1", "COM2", "COM3", "COM4", "COM5", "COM6", "COM7", "COM8", "COM9",
		"LPT1", "LPT2", "LPT3", "LPT4", "LPT5", "LPT6", "LPT7", "LPT8", "LPT9",
	}

	for _, name := range reservedNames {
		t.Run("Reserved name "+name, func(t *testing.T) {
			// Test with various extensions and cases
			testCases := []string{
				name,
				name + ".txt",
				name + "/",
				name + ".ext",
				"prefix" + name + "suffix", // Should be allowed if not exact match
				strings.ToLower(name),
				strings.ToUpper(name),
			}

			for i, testName := range testCases {
				result := IsValidPath(testName)
				// Only exact matches should be invalid
				expected := !(i < 4 && strings.EqualFold(testName, name))
				if result != expected {
					t.Errorf("IsValidPath(%q) = %v, expected %v", testName, result, expected)
				}
			}
		})
	}
}

func TestIsValidPath_DangerousPaths(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"Current directory", ".", false},
		{"Parent directory", "..", false},
		{"Path with parent ref", "../file.txt", false},
		{"Path with current ref", "./file.txt", false},
		{"Nested parent ref", "path/../../file.txt", false},
		{"Empty string", "", false},
		{"Only slash", "/", false},
		{"Only slashes", "///", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidPath(tt.input)
			if result != tt.expected {
				t.Errorf("IsValidPath(%q) = %v, expected %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestIsValidPath_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"Very long name", strings.Repeat("a", 255), true},
		{"Too long name", strings.Repeat("a", 256), false},
		{"Name with unicode", "café.txt", true},
		{"Name with emoji", "file🐛.txt", true},
		{"Name with mixed separators", "path\\to/file", false}, // Mixed separators
		{"Name with only spaces", "   ", false},
		{"Name with leading/trailing spaces", " file.txt ", true}, // Should be trimmed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidPath(tt.input)
			if result != tt.expected {
				t.Errorf("IsValidPath(%q) = %v, expected %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestIsValidPath_RealWorldExamples(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"Common config file", "config.yml", true},
		{"Common source file", "src/main.js", true},
		{"Package file", "package.json", true},
		{"Hidden file", ".gitignore", true},
		{"Directory with special chars", "node_modules/", true},
		{"Windows-style path", "C:\\file.txt", false},
		{"URL", "https://example.com", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidPath(tt.input)
			if result != tt.expected {
				t.Errorf("IsValidPath(%q) = %v, expected %v", tt.input, result, tt.expected)
			}
		})
	}
}

// Benchmark test for performance
func BenchmarkIsValidPath(b *testing.B) {
	testCases := []string{
		"valid-file.txt",
		"invalid:file.txt",
		"very/long/path/to/a/file.txt",
		"CON", // Windows reserved
		"normal_directory/",
	}

	for i := 0; i < b.N; i++ {
		for _, testCase := range testCases {
			IsValidPath(testCase)
		}
	}
}
