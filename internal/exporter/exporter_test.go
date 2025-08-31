package exporter

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestExportStructure_BasicExport(t *testing.T) {
	// Creating a temporary test structure
	tempDir := t.TempDir()
	createTestFiles(tempDir, map[string]string{
		"file1.txt":     "Content 1",
		"file2.go":      "Content 2",
		"dir1/file3.md": "Content 3",
	})

	// Creating the output file
	outputFile := filepath.Join(tempDir, "export.txt")

	// Exporting the structure
	err := ExportStructure(outputFile, tempDir, nil, nil, 0, false)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Checking the contents
	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	// Check that all files are exported
	expectedFiles := []string{"file1.txt", "file2.go", "dir1/file3.md"}
	for _, file := range expectedFiles {
		if !strings.Contains(string(content), "# "+file) {
			t.Errorf("Expected file %s in export, but not found", file)
		}
		if !strings.Contains(string(content), "Content") {
			t.Errorf("Expected content for %s in export, but not found", file)
		}
	}
}

func TestExportStructure_WithExtensionFilter(t *testing.T) {
	// Creating a temporary test structure
	tempDir := t.TempDir()
	createTestFiles(tempDir, map[string]string{
		"file1.txt":     "Text content",
		"file2.go":      "Go content",
		"file3.md":      "Markdown content",
		"dir1/file4.go": "Go content 2",
	})

	// Creating the output file
	outputFile := filepath.Join(tempDir, "export.txt")

	/// We export only .go files
	err := ExportStructure(outputFile, tempDir, []string{"go"}, nil, 0, false)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Checking the contents
	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	// Check that only .go files are exported
	if !strings.Contains(string(content), "# file2.go") {
		t.Error("Expected file2.go in export, but not found")
	}
	if !strings.Contains(string(content), "# dir1/file4.go") {
		t.Error("Expected dir1/file4.go in export, but not found")
	}

	// Check that other files are not exported
	if strings.Contains(string(content), "# file1.txt") {
		t.Error("file1.txt should not be in export")
	}
	if strings.Contains(string(content), "# file3.md") {
		t.Error("file3.md should not be in export")
	}
}

func TestExportStructure_WithIgnoreDirs(t *testing.T) {
	// Creating a temporary test structure
	tempDir := t.TempDir()
	createTestFiles(tempDir, map[string]string{
		"file1.txt":              "Content 1",
		"dir1/file2.txt":         "Content 2",
		"dir2/file3.txt":         "Content 3",
		"node_modules/file4.txt": "Content 4",
		".git/file5.txt":         "Content 5",
	})

	// Creating the output file
	outputFile := filepath.Join(tempDir, "export.txt")

	// Exporting, ignoring dir2 and node_modules
	err := ExportStructure(outputFile, tempDir, nil, []string{"dir2"}, 0, false)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Checking the contents
	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	// Check that dir1 is exported
	if !strings.Contains(string(content), "# dir1/file2.txt") {
		t.Error("Expected dir1/file2.txt in export, but not found")
	}

	// Check that dir2 is not exported
	if strings.Contains(string(content), "# dir2/file3.txt") {
		t.Error("dir2/file3.txt should not be in export")
	}

	// Check that node_modules is not exported (by default)
	if strings.Contains(string(content), "# node_modules/file4.txt") {
		t.Error("node_modules/file4.txt should not be in export")
	}

	// Check that .git is not exported (by default)
	if strings.Contains(string(content), "# .git/file5.txt") {
		t.Error(".git/file5.txt should not be in export")
	}
}

func TestExportStructure_WithMaxSize(t *testing.T) {
	// Creating a temporary test structure
	tempDir := t.TempDir()
	createTestFiles(tempDir, map[string]string{
		"small.txt": "Small content",
		"large.txt": strings.Repeat("a", 200*1024), // 200KB
	})

	// Creating the output file
	outputFile := filepath.Join(tempDir, "export.txt")

	// We export with a limit of 100KB
	err := ExportStructure(outputFile, tempDir, nil, nil, 100*1024, false)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Checking the contents
	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	// Check that the small file has been exported
	if !strings.Contains(string(content), "# small.txt") {
		t.Error("Expected small.txt in export, but not found")
	}

	// Check that the large file has not been exported
	if strings.Contains(string(content), "# large.txt") {
		t.Error("large.txt should not be in export")
	}
}

func TestExportStructure_WithHiddenFiles(t *testing.T) {
	// Creating a temporary test structure
	tempDir := t.TempDir()
	createTestFiles(tempDir, map[string]string{
		"file.txt":      "Regular file",
		".hidden.txt":   "Hidden file",
		"dir1/.hidden2": "Hidden in dir",
		".gitignore":    "Git ignore",
		"dir2/file.txt": "File in dir",
		"dir2/.env":     "Env file",
	})

	// Creating the output file
	outputFile := filepath.Join(tempDir, "export.txt")

	// First we check without including hidden files
	err := ExportStructure(outputFile, tempDir, nil, nil, 0, false)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Checking the contents
	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	// Check that hidden files are not exported
	if strings.Contains(string(content), "# .hidden.txt") {
		t.Error(".hidden.txt should not be in export")
	}
	if strings.Contains(string(content), "# dir1/.hidden2") {
		t.Error("dir1/.hidden2 should not be in export")
	}
	if strings.Contains(string(content), "# .gitignore") {
		t.Error(".gitignore should not be in export")
	}
	if strings.Contains(string(content), "# dir2/.env") {
		t.Error("dir2/.env should not be in export")
	}

	// Now we check with the inclusion of hidden files
	err = ExportStructure(outputFile, tempDir, nil, nil, 0, true)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Checking the content
	content, err = os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	// Check that the hidden files have been exported
	if !strings.Contains(string(content), "# .hidden.txt") {
		t.Error(".hidden.txt should be in export")
	}
	if !strings.Contains(string(content), "# dir1/.hidden2") {
		t.Error("dir1/.hidden2 should be in export")
	}
	if !strings.Contains(string(content), "# .gitignore") {
		t.Error(".gitignore should be in export")
	}
	if !strings.Contains(string(content), "# dir2/.env") {
		t.Error("dir2/.env should be in export")
	}
}

func TestExportStructure_SkipOutputFile(t *testing.T) {
	// Creating a temporary test structure
	tempDir := t.TempDir()
	createTestFiles(tempDir, map[string]string{
		"file1.txt": "Content 1",
		"file2.txt": "Content 2",
	})

	// Creating the output file INSIDE the exported directory
	outputFile := filepath.Join(tempDir, "export.txt")

	// Exporting the structure
	err := ExportStructure(outputFile, tempDir, nil, nil, 0, false)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Checking the content
	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	// Check that the output file is not included in the export
	if strings.Contains(string(content), "# export.txt") {
		t.Error("export.txt should not be in its own export")
	}

	// Check that other files have been exported
	if !strings.Contains(string(content), "# file1.txt") {
		t.Error("Expected file1.txt in export, but not found")
	}
	if !strings.Contains(string(content), "# file2.txt") {
		t.Error("Expected file2.txt in export, but not found")
	}
}

func TestExportStructure_EmptyDirectory(t *testing.T) {
	// Creating an empty temporary directory
	tempDir := t.TempDir()
	outputFile := filepath.Join(tempDir, "export.txt")

	// Exporting the structure
	err := ExportStructure(outputFile, tempDir, nil, nil, 0, false)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Checking the content
	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	// Check that the export is empty
	if len(content) > 0 {
		t.Error("Export from empty directory should be empty")
	}
}

func TestExportStructure_InvalidRootDir(t *testing.T) {
	// Creating a temporary directory for the output file
	tempDir := t.TempDir()
	outputFile := filepath.Join(tempDir, "export.txt")

	// Using a non-existent directory
	invalidDir := filepath.Join(tempDir, "nonexistent")

	// Trying to export
	err := ExportStructure(outputFile, invalidDir, nil, nil, 0, false)
	if err == nil {
		t.Fatal("Expected error for invalid root directory, but got none")
	}

	// Check that the output file has not been created
	if _, err := os.Stat(outputFile); !os.IsNotExist(err) {
		t.Error("Output file should not be created when root directory is invalid")
	}
}

func TestExportStructure_InvalidOutputFile(t *testing.T) {
	// Creating a temporary test structure
	tempDir := t.TempDir()
	createTestFiles(tempDir, map[string]string{
		"file.txt": "Content",
	})

	// Trying to use an invalid path for the output file
	invalidOutput := string([]rune{0}) // null character

	// Trying to export
	err := ExportStructure(invalidOutput, tempDir, nil, nil, 0, false)
	if err == nil {
		t.Fatal("Expected error for invalid output file path, but got none")
	}
}

func TestExportStructure_PermissionError(t *testing.T) {
	// Creating a temporary test structure
	tempDir := t.TempDir()

	// Creating a file with limited rights
	restrictedFile := filepath.Join(tempDir, "restricted.txt")
	if err := os.WriteFile(restrictedFile, []byte("Secret content"), 0400); err != nil {
		t.Fatalf("Failed to create restricted file: %v", err)
	}

	// Creating the output file
	outputFile := filepath.Join(tempDir, "export.txt")

	// Exporting the structure
	err := ExportStructure(outputFile, tempDir, nil, nil, 0, false)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Checking the content
	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	// Check that the file with limited rights is exported anyway
	if !strings.Contains(string(content), "# restricted.txt") {
		t.Error("restricted.txt should be in export despite restricted permissions")
	}
}

// Auxiliary function for creating test files
func createTestFiles(root string, files map[string]string) {
	for path, content := range files {
		// Creating a directory, if necessary
		dir := filepath.Dir(path)
		if dir != "." {
			fullDir := filepath.Join(root, dir)
			if err := os.MkdirAll(fullDir, 0755); err != nil {
				panic(err)
			}
		}

		// Creating a file
		fullPath := filepath.Join(root, path)
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			panic(err)
		}
	}
}
