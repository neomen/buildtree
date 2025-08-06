package main

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// ... (основной код программы без изменений) ...

func TestCreateStructure(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name: "Simple structure",
			input: `project/
├── src/
│   └── main.go
├── pkg/
│   └── utils.go
└── README.md`,
			expected: []string{
				"project",
				"project/src",
				"project/src/main.go",
				"project/pkg",
				"project/pkg/utils.go",
				"project/README.md",
			},
		},
		{
			name: "With comments",
			input: `app/              # Main application
├── cmd/          # Commands
│   └── run.sh    # Runner script
└── config/       # Configuration
    └── app.yaml  # Config file`,
			expected: []string{
				"app",
				"app/cmd",
				"app/cmd/run.sh",
				"app/config",
				"app/config/app.yaml",
			},
		},
		{
			name: "Multi-level nesting",
			input: `root/
├── a/
│   ├── b/
│   │   └── file1.txt
│   └── c/
│       └── file2.txt
└── d/
    └── e/
        └── file3.txt`,
			expected: []string{
				"root",
				"root/a",
				"root/a/b",
				"root/a/b/file1.txt",
				"root/a/c",
				"root/a/c/file2.txt",
				"root/d",
				"root/d/e",
				"root/d/e/file3.txt",
			},
		},
		{
			name: "Mixed files and directories",
			input: `test/
├── file1
├── dir1/
│   ├── file2
│   └── subdir/
│       └── file3
└── file4`,
			expected: []string{
				"test",
				"test/file1",
				"test/dir1",
				"test/dir1/file2",
				"test/dir1/subdir",
				"test/dir1/subdir/file3",
				"test/file4",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp directory
			tempDir, err := os.MkdirTemp("", "quicktree_test")
			if err != nil {
				t.Fatalf("Failed to create temp dir: %v", err)
			}
			defer os.RemoveAll(tempDir)

			// Change working directory to temp
			oldDir, err := os.Getwd()
			if err != nil {
				t.Fatalf("Failed to get current dir: %v", err)
			}
			defer os.Chdir(oldDir)
			os.Chdir(tempDir)

			// Process test input
			processInput(tt.input)

			// Verify created structure
			for _, path := range tt.expected {
				fullPath := filepath.Join(tempDir, path)
				if _, err := os.Stat(fullPath); os.IsNotExist(err) {
					t.Errorf("Expected path not created: %s", path)
				}
			}

			// Check for unexpected files
			filepath.Walk(tempDir, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}

				relPath, err := filepath.Rel(tempDir, path)
				if err != nil {
					return err
				}

				if relPath == "." {
					return nil
				}

				found := false
				for _, expected := range tt.expected {
					if expected == relPath {
						found = true
						break
					}
				}

				if !found {
					t.Errorf("Unexpected path created: %s", relPath)
				}
				return nil
			})
		})
	}
}

func TestCommandLineFlags(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		input    string
		expected []string
	}{
		{
			name: "Direct input",
			args: []string{"-h=false", "project/\n└── file.txt"},
			expected: []string{
				"project",
				"project/file.txt",
			},
		},
		{
			name: "File input",
			args: []string{"-f", "test_structure.txt"},
			input: `dir/
├── subdir/
│   └── file`,
			expected: []string{
				"dir",
				"dir/subdir",
				"dir/subdir/file",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp directory
			tempDir, err := os.MkdirTemp("", "quicktree_test")
			if err != nil {
				t.Fatalf("Failed to create temp dir: %v", err)
			}
			defer os.RemoveAll(tempDir)

			// Create test file if needed
			if tt.input != "" {
				filePath := filepath.Join(tempDir, "test_structure.txt")
				if err := os.WriteFile(filePath, []byte(tt.input), 0644); err != nil {
					t.Fatalf("Failed to create test file: %v", err)
				}
				// Update args with actual file path
				for i, arg := range tt.args {
					if arg == "test_structure.txt" {
						tt.args[i] = filePath
					}
				}
			}

			// Change working directory to temp
			oldDir, err := os.Getwd()
			if err != nil {
				t.Fatalf("Failed to get current dir: %v", err)
			}
			defer os.Chdir(oldDir)
			os.Chdir(tempDir)

			// Backup and restore command-line args
			oldArgs := os.Args
			defer func() { os.Args = oldArgs }()
			os.Args = append([]string{"quicktree"}, tt.args...)

			// Run main function
			main()

			// Verify created structure
			for _, path := range tt.expected {
				fullPath := filepath.Join(tempDir, path)
				if _, err := os.Stat(fullPath); os.IsNotExist(err) {
					t.Errorf("Expected path not created: %s", path)
				}
			}
		})
	}
}

func TestHelpOutput(t *testing.T) {
	// Backup and restore command-line args
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"quicktree", "-h"}

	// Capture stdout
	oldStdout := os.Stdout
	defer func() { os.Stdout = oldStdout }()
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Run main function
	main()
	w.Close()

	// Read captured output
	out, _ := io.ReadAll(r)

	// Check for help content
	if !strings.Contains(string(out), "Usage:") {
		t.Error("Help output doesn't contain usage information")
	}
	if !strings.Contains(string(out), "Examples:") {
		t.Error("Help output doesn't contain examples")
	}
}
