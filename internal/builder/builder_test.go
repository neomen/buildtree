package builder

import (
	"os"
	"strings"
	"testing"

	"github.com/neomen/buildtree/internal/parser"
)

func TestBuildTree_SimpleStructure(t *testing.T) {
	tempDir := t.TempDir()
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	// Restoring the original directory
	defer os.Chdir(originalDir)

	err = os.Chdir(tempDir)
	if err != nil {
		t.Fatal(err)
	}

	root := &parser.Node{
		Name:  "project",
		IsDir: true,
		Level: 0,
		Children: []*parser.Node{
			{
				Name:  "src",
				IsDir: true,
				Level: 1,
				Children: []*parser.Node{
					{
						Name:  "main.go",
						IsDir: false,
						Level: 2,
					},
				},
			},
			{
				Name:  "README.md",
				IsDir: false,
				Level: 1,
			},
		},
	}

	err = BuildTree(root, 0)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Check if all files and directories were created
	assertDirExists(t, "project")
	assertDirExists(t, "project/src")
	assertFileExists(t, "project/src/main.go")
	assertFileExists(t, "project/README.md")
}

func TestBuildTree_ComplexStructure(t *testing.T) {
	tempDir := t.TempDir()
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(originalDir)

	err = os.Chdir(tempDir)
	if err != nil {
		t.Fatal(err)
	}

	root := &parser.Node{
		Name:  "app",
		IsDir: true,
		Level: 0,
		Children: []*parser.Node{
			{
				Name:  "src",
				IsDir: true,
				Level: 1,
				Children: []*parser.Node{
					{
						Name:  "components",
						IsDir: true,
						Level: 2,
						Children: []*parser.Node{
							{
								Name:  "Button.js",
								IsDir: false,
								Level: 3,
							},
							{
								Name:  "Header.js",
								IsDir: false,
								Level: 3,
							},
						},
					},
					{
						Name:  "utils",
						IsDir: true,
						Level: 2,
						Children: []*parser.Node{
							{
								Name:  "helpers.js",
								IsDir: false,
								Level: 3,
							},
						},
					},
					{
						Name:  "index.js",
						IsDir: false,
						Level: 2,
					},
				},
			},
			{
				Name:  "public",
				IsDir: true,
				Level: 1,
				Children: []*parser.Node{
					{
						Name:  "index.html",
						IsDir: false,
						Level: 2,
					},
				},
			},
			{
				Name:  "package.json",
				IsDir: false,
				Level: 1,
			},
		},
	}

	err = BuildTree(root, 0)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Check if all files and directories were created
	assertDirExists(t, "app")
	assertDirExists(t, "app/src")
	assertDirExists(t, "app/src/components")
	assertDirExists(t, "app/src/utils")
	assertDirExists(t, "app/public")
	assertFileExists(t, "app/src/components/Button.js")
	assertFileExists(t, "app/src/components/Header.js")
	assertFileExists(t, "app/src/utils/helpers.js")
	assertFileExists(t, "app/src/index.js")
	assertFileExists(t, "app/public/index.html")
	assertFileExists(t, "app/package.json")
}

func TestBuildTree_MaxDepth(t *testing.T) {
	tempDir := t.TempDir()
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(originalDir)

	err = os.Chdir(tempDir)
	if err != nil {
		t.Fatal(err)
	}

	root := &parser.Node{
		Name:  "project",
		IsDir: true,
		Level: 0,
		Children: []*parser.Node{
			{
				Name:  "level1",
				IsDir: true,
				Level: 1,
				Children: []*parser.Node{
					{
						Name:  "level2",
						IsDir: true,
						Level: 2,
						Children: []*parser.Node{
							{
								Name:  "level3",
								IsDir: true,
								Level: 3,
								Children: []*parser.Node{
									{
										Name:  "file.txt",
										IsDir: false,
										Level: 4,
									},
								},
							},
						},
					},
				},
			},
		},
	}

	// Test with maxDepth = 2 (should only create up to level2)
	err = BuildTree(root, 2)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Check that directories up to level2 were created
	assertDirExists(t, "project")
	assertDirExists(t, "project/level1")
	assertDirExists(t, "project/level1/level2")

	// Check that level3 and file.txt were NOT created
	assertNotExists(t, "project/level1/level2/level3")
	assertNotExists(t, "project/level1/level2/level3/file.txt")
}

func TestBuildTree_ExistingDirectories(t *testing.T) {
	tempDir := t.TempDir()
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(originalDir)

	err = os.Chdir(tempDir)
	if err != nil {
		t.Fatal(err)
	}

	// Create the root directory manually first
	err = os.Mkdir("project", 0755)
	if err != nil {
		t.Fatalf("Failed to create directory: %v", err)
	}

	root := &parser.Node{
		Name:  "project",
		IsDir: true,
		Level: 0,
		Children: []*parser.Node{
			{
				Name:  "src",
				IsDir: true,
				Level: 1,
				Children: []*parser.Node{
					{
						Name:  "main.go",
						IsDir: false,
						Level: 2,
					},
				},
			},
		},
	}

	// Should not fail even though project directory already exists
	err = BuildTree(root, 0)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Check that the new files and directories were created
	assertDirExists(t, "project/src")
	assertFileExists(t, "project/src/main.go")
}

func TestBuildTree_InvalidNames(t *testing.T) {
	tempDir := t.TempDir()
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(originalDir)

	err = os.Chdir(tempDir)
	if err != nil {
		t.Fatal(err)
	}

	root := &parser.Node{
		Name:  "project",
		IsDir: true,
		Level: 0,
		Children: []*parser.Node{
			{
				Name:  "valid_dir",
				IsDir: true,
				Level: 1,
			},
			{
				Name:  "invalid:name", // Contains invalid character
				IsDir: true,
				Level: 1,
			},
			{
				Name:  "valid_file.txt",
				IsDir: false,
				Level: 1,
			},
			{
				Name:  "invalid*file", // Contains invalid character
				IsDir: false,
				Level: 1,
			},
		},
	}

	err = BuildTree(root, 0)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Check that only valid names were created
	assertDirExists(t, "project")
	assertDirExists(t, "project/valid_dir")
	assertNotExists(t, "project/invalid:name")
	assertFileExists(t, "project/valid_file.txt")
	assertNotExists(t, "project/invalid*file")
}

func TestBuildTree_DotDotPath(t *testing.T) {
	tempDir := t.TempDir()
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(originalDir)

	err = os.Chdir(tempDir)
	if err != nil {
		t.Fatal(err)
	}

	root := &parser.Node{
		Name:  "project",
		IsDir: true,
		Level: 0,
		Children: []*parser.Node{
			{
				Name:  "..", // Attempt to navigate up
				IsDir: true,
				Level: 1,
			},
			{
				Name:  "normal_dir",
				IsDir: true,
				Level: 1,
			},
		},
	}

	err = BuildTree(root, 0)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Check that the malicious path was not created
	// We need to check the actual contents of the project directory
	entries, err := os.ReadDir("project")
	if err != nil {
		t.Fatalf("Error reading project directory: %v", err)
	}

	// Count non-hidden entries (excluding . and ..)
	var visibleEntries []string
	for _, entry := range entries {
		if entry.Name() != "." && entry.Name() != ".." {
			visibleEntries = append(visibleEntries, entry.Name())
		}
	}

	// Should only have one visible entry: normal_dir
	if len(visibleEntries) != 1 {
		t.Errorf("Expected 1 visible entry, got %d: %v", len(visibleEntries), visibleEntries)
	}

	if len(visibleEntries) > 0 && visibleEntries[0] != "normal_dir" {
		t.Errorf("Expected 'normal_dir', got '%s'", visibleEntries[0])
	}
}

func TestBuildTree_WindowsReservedNames(t *testing.T) {
	tempDir := t.TempDir()
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(originalDir)

	err = os.Chdir(tempDir)
	if err != nil {
		t.Fatal(err)
	}

	root := &parser.Node{
		Name:  "project",
		IsDir: true,
		Level: 0,
		Children: []*parser.Node{
			{
				// Windows reserved name
				Name:  "CON",
				IsDir: true,
				Level: 1,
			},
			{
				Name:  "normal_dir",
				IsDir: true,
				Level: 1,
			},
			{
				// Windows reserved name
				Name:  "LPT1",
				IsDir: false,
				Level: 1,
			},
			{
				Name:  "normal_file.txt",
				IsDir: false,
				Level: 1,
			},
		},
	}

	err = BuildTree(root, 0)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Check that Windows reserved names were not created
	// On Windows, we can't use standard methods to check for reserved names
	// as they get redirected to devices. Instead, we check the directory contents.
	entries, err := os.ReadDir("project")
	if err != nil {
		t.Fatalf("Error reading project directory: %v", err)
	}

	// Check that only normal names were created
	foundNormalDir := false
	foundNormalFile := false

	for _, entry := range entries {
		name := entry.Name()
		switch name {
		case "normal_dir":
			foundNormalDir = true
		case "normal_file.txt":
			foundNormalFile = true
		case "CON", "LPT1":
			t.Errorf("Reserved name %s was created but should not be", name)
		}
	}

	if !foundNormalDir {
		t.Error("Normal directory was not created")
	}
	if !foundNormalFile {
		t.Error("Normal file was not created")
	}

	// Check that normal names were created
	assertDirExists(t, "project/normal_dir")
	assertFileExists(t, "project/normal_file.txt")
}

func TestBuildTree_EmptyNode(t *testing.T) {
	tempDir := t.TempDir()
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(originalDir)

	err = os.Chdir(tempDir)
	if err != nil {
		t.Fatal(err)
	}

	// Test with an empty node (should return error)
	root := &parser.Node{
		Name:  "",
		IsDir: true,
		Level: 0,
	}

	err = BuildTree(root, 0)
	if err == nil {
		t.Error("Expected error for empty node name, but got none")
	}
}

// Helper functions for assertions
func assertDirExists(t *testing.T, path string) {
	t.Helper()
	info, err := os.Stat(path)
	if err != nil {
		t.Errorf("Expected directory %s to exist, but got error: %v", path, err)
		return
	}
	if !info.IsDir() {
		t.Errorf("Expected %s to be a directory, but it's a file", path)
	}
}

func assertFileExists(t *testing.T, path string) {
	t.Helper()
	info, err := os.Stat(path)
	if err != nil {
		t.Errorf("Expected file %s to exist, but got error: %v", path, err)
		return
	}
	if info.IsDir() {
		t.Errorf("Expected %s to be a file, but it's a directory", path)
	}
}

func assertNotExists(t *testing.T, path string) {
	t.Helper()
	_, err := os.Stat(path)
	if err == nil {
		t.Errorf("Expected %s to not exist, but it does", path)
		return
	}

	// On Windows, we ignore errors related to invalid file names
	// as this is the expected behavior
	if !os.IsNotExist(err) {
		// Checking if the error is related to an invalid file name in Windows
		if isWindowsInvalidNameError(err) {
			// This is the expected behavior for Windows, so we don't consider it an error
			return
		}
		t.Errorf("Unexpected error checking if %s exists: %v", path, err)
	}
}

// isWindowsInvalidNameError checks if the error is related to an invalid file name in Windows
func isWindowsInvalidNameError(err error) bool {
	return strings.Contains(err.Error(), "The filename, directory name, or volume label syntax is incorrect")
}
