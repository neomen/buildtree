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
	err := ExportStructure(outputFile, tempDir, nil, nil, nil, 0, false, false)
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
		// Convert to platform-independent slash format for comparison
		expectedPath := filepath.ToSlash(file)
		if !strings.Contains(string(content), "# "+expectedPath) {
			t.Errorf("Expected file %s in export, but not found", expectedPath)
		}
		if !strings.Contains(string(content), "Content") {
			t.Errorf("Expected content for %s in export, but not found", expectedPath)
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
	// Exporting only .go files
	err := ExportStructure(outputFile, tempDir, []string{"go"}, nil, nil, 0, false, false)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	// Checking the contents
	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}
	// Check that only .go files are exported
	expectedPaths := []string{"file2.go", "dir1/file4.go"}
	for _, path := range expectedPaths {
		expectedPath := filepath.ToSlash(path)
		if !strings.Contains(string(content), "# "+expectedPath) {
			t.Errorf("Expected %s in export, but not found", expectedPath)
		}
	}

	// Check that other files are not exported
	excludedPaths := []string{"file1.txt", "file3.md"}
	for _, path := range excludedPaths {
		expectedPath := filepath.ToSlash(path)
		if strings.Contains(string(content), "# "+expectedPath) {
			t.Errorf("%s should not be in export", expectedPath)
		}
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
	err := ExportStructure(outputFile, tempDir, nil, []string{"dir2"}, nil, 0, false, false)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	// Checking the contents
	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	// Check that dir1 is exported
	dir1Path := filepath.ToSlash("dir1/file2.txt")
	if !strings.Contains(string(content), "# "+dir1Path) {
		t.Error("Expected " + dir1Path + " in export, but not found")
	}

	// Check that dir2 is not exported
	dir2Path := filepath.ToSlash("dir2/file3.txt")
	if strings.Contains(string(content), "# "+dir2Path) {
		t.Error(dir2Path + " should not be in export")
	}

	// Check that node_modules is not exported (by default)
	nodeModulesPath := filepath.ToSlash("node_modules/file4.txt")
	if strings.Contains(string(content), "# "+nodeModulesPath) {
		t.Error(nodeModulesPath + " should not be in export")
	}

	// Check that .git is not exported (by default)
	gitPath := filepath.ToSlash(".git/file5.txt")
	if strings.Contains(string(content), "# "+gitPath) {
		t.Error(gitPath + " should not be in export")
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
	// Exporting with a limit of 100KB
	err := ExportStructure(outputFile, tempDir, nil, nil, nil, 100*1024, false, false)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	// Checking the contents
	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	// Check that the small file has been exported
	smallPath := filepath.ToSlash("small.txt")
	if !strings.Contains(string(content), "# "+smallPath) {
		t.Error("Expected " + smallPath + " in export, but not found")
	}

	// Check that the large file has not been exported
	largePath := filepath.ToSlash("large.txt")
	if strings.Contains(string(content), "# "+largePath) {
		t.Error(largePath + " should not be in export")
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
	err := ExportStructure(outputFile, tempDir, nil, nil, nil, 0, false, false)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	// Checking the contents
	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	// Check that hidden files are not exported
	hiddenFiles := []string{".hidden.txt", "dir1/.hidden2", ".gitignore", "dir2/.env"}
	for _, file := range hiddenFiles {
		expectedPath := filepath.ToSlash(file)
		if strings.Contains(string(content), "# "+expectedPath) {
			t.Errorf("%s should not be in export", expectedPath)
		}
	}

	// Now we check with the inclusion of hidden files
	err = ExportStructure(outputFile, tempDir, nil, nil, nil, 0, true, false)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	// Checking the content
	content, err = os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	// Check that the hidden files have been exported
	for _, file := range hiddenFiles {
		expectedPath := filepath.ToSlash(file)
		if !strings.Contains(string(content), "# "+expectedPath) {
			t.Errorf("%s should be in export", expectedPath)
		}
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
	err := ExportStructure(outputFile, tempDir, nil, nil, nil, 0, false, false)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	// Checking the content
	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	// Check that the output file is not included in the export
	outputFileName := filepath.Base(outputFile)
	if strings.Contains(string(content), "# "+outputFileName) {
		t.Error(outputFileName + " should not be in its own export")
	}

	// Check that other files have been exported
	expectedFiles := []string{"file1.txt", "file2.txt"}
	for _, file := range expectedFiles {
		expectedPath := filepath.ToSlash(file)
		if !strings.Contains(string(content), "# "+expectedPath) {
			t.Errorf("Expected %s in export, but not found", expectedPath)
		}
	}
}

func TestExportStructure_EmptyDirectory(t *testing.T) {
	// Creating an empty temporary directory
	tempDir := t.TempDir()
	outputFile := filepath.Join(tempDir, "export.txt")
	// Exporting the structure
	err := ExportStructure(outputFile, tempDir, nil, nil, nil, 0, false, false)
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
	err := ExportStructure(outputFile, invalidDir, nil, nil, nil, 0, false, false)
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
	err := ExportStructure(invalidOutput, tempDir, nil, nil, nil, 0, false, false)
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
	err := ExportStructure(outputFile, tempDir, nil, nil, nil, 0, false, false)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	// Checking the content
	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}
	// Check that the file with limited rights is exported anyway
	restrictedPath := filepath.ToSlash("restricted.txt")
	if !strings.Contains(string(content), "# "+restrictedPath) {
		t.Error(restrictedPath + " should be in export despite restricted permissions")
	}
}

// Auxiliary function for creating test files
func createTestFiles(root string, files map[string]string) {
	for path, content := range files {
		// Convert to system-specific path format
		path = filepath.FromSlash(path)

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
