package parser

import (
	"strings"
	"testing"
)

func TestParseInput_SimpleStructure(t *testing.T) {
	input := `project/
├── src/
│   └── main.go
└── README.md`

	root, err := ParseInput(input)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Check root node
	if root.Name != "project" {
		t.Errorf("Expected root name 'project', got '%s'", root.Name)
	}
	if !root.IsDir {
		t.Error("Root should be a directory")
	}
	if root.Level != 0 {
		t.Errorf("Root level should be 0, got %d", root.Level)
	}

	// Check children
	if len(root.Children) != 2 {
		t.Fatalf("Expected 2 children, got %d", len(root.Children))
	}

	// Check src directory
	src := root.Children[0]
	if src.Name != "src" {
		t.Errorf("Expected first child 'src', got '%s'", src.Name)
	}
	if !src.IsDir {
		t.Error("src should be a directory")
	}
	if src.Level != 1 {
		t.Errorf("src level should be 1, got %d", src.Level)
	}

	// Check main.go file
	if len(src.Children) != 1 {
		t.Fatalf("Expected 1 child in src, got %d", len(src.Children))
	}
	mainGo := src.Children[0]
	if mainGo.Name != "main.go" {
		t.Errorf("Expected file 'main.go', got '%s'", mainGo.Name)
	}
	if mainGo.IsDir {
		t.Error("main.go should be a file")
	}
	if mainGo.Level != 2 {
		t.Errorf("main.go level should be 2, got %d", mainGo.Level)
	}

	// Check README.md file
	readme := root.Children[1]
	if readme.Name != "README.md" {
		t.Errorf("Expected file 'README.md', got '%s'", readme.Name)
	}
	if readme.IsDir {
		t.Error("README.md should be a file")
	}
	if readme.Level != 1 {
		t.Errorf("README.md level should be 1, got %d", readme.Level)
	}
}

func TestParseInput_ComplexStructure(t *testing.T) {
	input := `app/
├── src/
│   ├── components/
│   │   ├── Button.js
│   │   └── Header.js
│   ├── utils/
│   │   └── helpers.js
│   └── index.js
├── public/
│   └── index.html
└── package.json`

	root, err := ParseInput(input)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Verify the structure
	if root.Name != "app" {
		t.Errorf("Expected root name 'app', got '%s'", root.Name)
	}
	if len(root.Children) != 3 {
		t.Fatalf("Expected 3 children, got %d", len(root.Children))
	}

	// Check src directory and its children
	src := findChild(root, "src")
	if src == nil {
		t.Fatal("src directory not found")
	}
	if len(src.Children) != 3 {
		t.Fatalf("Expected 3 children in src, got %d", len(src.Children))
	}

	// Check components directory
	components := findChild(src, "components")
	if components == nil {
		t.Fatal("components directory not found")
	}
	if len(components.Children) != 2 {
		t.Fatalf("Expected 2 children in components, got %d", len(components.Children))
	}

	// Check utils directory
	utils := findChild(src, "utils")
	if utils == nil {
		t.Fatal("utils directory not found")
	}
	if len(utils.Children) != 1 {
		t.Fatalf("Expected 1 child in utils, got %d", len(utils.Children))
	}
}

func TestParseInput_EmptyInput(t *testing.T) {
	_, err := ParseInput("")
	if err != ErrEmptyInput {
		t.Errorf("Expected ErrEmptyInput, got %v", err)
	}
}

func TestParseInput_Comments(t *testing.T) {
	input := `project/ # This is a comment
├── src/ # Source code
│   └── main.go # Main file
└── README.md # Documentation`

	root, err := ParseInput(input)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Check that comments are removed
	if root.Name != "project" {
		t.Errorf("Expected root name 'project', got '%s'", root.Name)
	}
	if strings.Contains(root.Name, "#") {
		t.Error("Comment should be removed from root name")
	}

	// Check children
	if len(root.Children) != 2 {
		t.Fatalf("Expected 2 children, got %d", len(root.Children))
	}

	src := root.Children[0]
	if strings.Contains(src.Name, "#") {
		t.Error("Comment should be removed from src name")
	}

	mainGo := src.Children[0]
	if strings.Contains(mainGo.Name, "#") {
		t.Error("Comment should be removed from main.go name")
	}
}

func TestParseInput_AlternativeSymbols(t *testing.T) {
	input := `project/
|-- src/
|   |-- main.go
|-- README.md`

	root, err := ParseInput(input)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Check root node
	if root.Name != "project" {
		t.Errorf("Expected root name 'project', got '%s'", root.Name)
	}

	// Check children
	if len(root.Children) != 2 {
		t.Fatalf("Expected 2 children, got %d", len(root.Children))
	}

	// Check src directory
	src := root.Children[0]
	if src.Name != "src" {
		t.Errorf("Expected first child 'src', got '%s'", src.Name)
	}

	// Check main.go file
	if len(src.Children) != 1 {
		t.Fatalf("Expected 1 child in src, got %d", len(src.Children))
	}
	mainGo := src.Children[0]
	if mainGo.Name != "main.go" {
		t.Errorf("Expected file 'main.go', got '%s'", mainGo.Name)
	}
}

func TestParseLine(t *testing.T) {
	tests := []struct {
		name          string
		line          string
		expectedLevel int
		expectedName  string
		expectedIsDir bool
	}{
		{
			name:          "Simple file",
			line:          "    └── file.txt",
			expectedLevel: 2, // 4 пробела + 3 символа + пробел └── = 8 символов / 4 = 2
			expectedName:  "file.txt",
			expectedIsDir: false,
		},
		{
			name:          "Directory with slash",
			line:          "    ├── src/",
			expectedLevel: 2, // 4 пробела + 3 символа + пробел ├── = 8 символов / 4 = 2
			expectedName:  "src",
			expectedIsDir: true,
		},
		{
			name:          "Nested structure",
			line:          "        └── deepfile.go",
			expectedLevel: 3, // 8 пробелов + 3 символа + пробел └── = 12 символов / 4 = 2
			expectedName:  "deepfile.go",
			expectedIsDir: false,
		},
		{
			name:          "With comment",
			line:          "    ├── config.yml # Configuration file",
			expectedLevel: 2, // 4 пробела + 3 символа + пробел ├── = 8 символов / 4 = 1
			expectedName:  "config.yml",
			expectedIsDir: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			level, name, isDir := parseLine(tt.line)

			if level != tt.expectedLevel {
				t.Errorf("Expected level %d, got %d", tt.expectedLevel, level)
			}

			if name != tt.expectedName {
				t.Errorf("Expected name '%s', got '%s'", tt.expectedName, name)
			}

			if isDir != tt.expectedIsDir {
				t.Errorf("Expected isDir %t, got %t", tt.expectedIsDir, isDir)
			}
		})
	}
}

func TestExtractName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Standard tree symbol",
			input:    "    └── filename.txt",
			expected: "filename.txt",
		},
		{
			name:     "Alternative tree symbol",
			input:    "    |-- filename.txt",
			expected: "filename.txt",
		},
		{
			name:     "Directory with slash",
			input:    "    ├── dirname/",
			expected: "dirname/",
		},
		{
			name:     "With extra spaces",
			input:    "    └──   filename with spaces.txt  ",
			expected: "filename with spaces.txt",
		},
		{
			name:     "With comment",
			input:    "    └── file.txt # This is a comment",
			expected: "file.txt",
		},
		{
			name:     "With comment and no space",
			input:    "    └── file.txt#This is a comment",
			expected: "file.txt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractName(tt.input)
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestNormalizeTreeSymbols(t *testing.T) {
	input := `project/
|-- src/
|   |-- main.go
'-- README.md`

	expected := `project/
├── src/
│   ├── main.go
└── README.md`

	result := normalizeTreeSymbols(input)
	if result != expected {
		t.Errorf("Normalization failed.\nExpected:\n%s\nGot:\n%s", expected, result)
	}
}

// Helper function to find a child node by name
func findChild(parent *Node, name string) *Node {
	for _, child := range parent.Children {
		if child.Name == name {
			return child
		}
	}
	return nil
}
