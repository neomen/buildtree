package main

import (
	"bytes"
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/neomen/buildtree/internal/parser"
)

// Mock implementations for testing
type mockParser struct {
	parseFunc func(input string) (*parser.Node, error)
}

func (m *mockParser) ParseInput(input string) (*parser.Node, error) {
	return m.parseFunc(input)
}

type mockBuilder struct {
	buildFunc func(root *parser.Node, maxDepth int) error
}

func (m *mockBuilder) BuildTree(root *parser.Node, maxDepth int) error {
	return m.buildFunc(root, maxDepth)
}

func TestRun_HelpFlag(t *testing.T) {
	// Mock dependencies
	p := &mockParser{}
	b := &mockBuilder{}

	// Capture stdout and stderr
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	// Test help flag
	exitCode := run([]string{"-h"}, &bytes.Buffer{}, stdout, stderr, p, b)

	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got %d", exitCode)
	}

	if !strings.Contains(stdout.String(), "Buildtree - Instant Directory Tree Builder") {
		t.Error("Help text was not printed")
	}
}

func TestRun_VersionFlag(t *testing.T) {
	// Mock dependencies
	p := &mockParser{}
	b := &mockBuilder{}

	// Capture stdout and stderr
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	// Test version flag
	exitCode := run([]string{"-v"}, &bytes.Buffer{}, stdout, stderr, p, b)

	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got %d", exitCode)
	}

	if !strings.Contains(stdout.String(), "buildtree v"+version) {
		t.Error("Version information was not printed")
	}
}

func TestRun_NoInput(t *testing.T) {
	// Mock dependencies
	p := &mockParser{}
	b := &mockBuilder{}

	// Capture stdout and stderr
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	// Test with no input
	exitCode := run([]string{}, &bytes.Buffer{}, stdout, stderr, p, b)

	if exitCode != 1 {
		t.Errorf("Expected exit code 1, got %d", exitCode)
	}

	if !strings.Contains(stderr.String(), "No input structure provided") {
		t.Error("Error message was not printed")
	}
}

func TestRun_FileInput(t *testing.T) {
	// Create a temporary file with test content
	tempFile, err := os.CreateTemp("", "test-input")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tempFile.Name())

	content := `project/
├── src/
│   └── main.go
└── README.md`
	_, err = tempFile.WriteString(content)
	if err != nil {
		t.Fatal(err)
	}
	tempFile.Close()

	// Mock dependencies
	p := &mockParser{
		parseFunc: func(input string) (*parser.Node, error) {
			if input != content {
				t.Errorf("Expected content %q, got %q", content, input)
			}
			return &parser.Node{Name: "project", IsDir: true}, nil
		},
	}

	b := &mockBuilder{
		buildFunc: func(root *parser.Node, maxDepth int) error {
			if root.Name != "project" {
				t.Errorf("Expected root name 'project', got %q", root.Name)
			}
			if maxDepth != 20 {
				t.Errorf("Expected maxDepth 20, got %d", maxDepth)
			}
			return nil
		},
	}

	// Capture stdout and stderr
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	// Test file input
	exitCode := run([]string{"-i", tempFile.Name()}, &bytes.Buffer{}, stdout, stderr, p, b)

	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got %d", exitCode)
	}
}

func TestRun_StdinInput(t *testing.T) {
	input := `project/
├── src/
│   └── main.go
└── README.md`

	stdin := strings.NewReader(input)

	// Mock dependencies
	p := &mockParser{
		parseFunc: func(actualInput string) (*parser.Node, error) {
			if actualInput != input {
				t.Errorf("Expected content %q, got %q", input, actualInput)
			}
			return &parser.Node{Name: "project", IsDir: true}, nil
		},
	}

	b := &mockBuilder{
		buildFunc: func(root *parser.Node, maxDepth int) error {
			if root.Name != "project" {
				t.Errorf("Expected root name 'project', got %q", root.Name)
			}
			return nil
		},
	}

	// Capture stdout and stderr
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	// Test stdin input
	exitCode := run([]string{input}, stdin, stdout, stderr, p, b)

	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got %d", exitCode)
	}
}

func TestRun_ParseError(t *testing.T) {
	// Mock dependencies with error
	p := &mockParser{
		parseFunc: func(input string) (*parser.Node, error) {
			return nil, errors.New("parse error")
		},
	}

	b := &mockBuilder{
		buildFunc: func(root *parser.Node, maxDepth int) error {
			return nil
		},
	}

	// Capture stdout and stderr
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	// Test parse error
	exitCode := run([]string{"project/"}, &bytes.Buffer{}, stdout, stderr, p, b)

	if exitCode != 1 {
		t.Errorf("Expected exit code 1, got %d", exitCode)
	}

	if !strings.Contains(stderr.String(), "Error parsing input") {
		t.Error("Error message was not printed")
	}
}

func TestRun_BuildError(t *testing.T) {
	// Mock dependencies with error
	p := &mockParser{
		parseFunc: func(input string) (*parser.Node, error) {
			return &parser.Node{Name: "project", IsDir: true}, nil
		},
	}

	b := &mockBuilder{
		buildFunc: func(root *parser.Node, maxDepth int) error {
			return errors.New("build error")
		},
	}

	// Capture stdout and stderr
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	// Test build error
	exitCode := run([]string{"project/"}, &bytes.Buffer{}, stdout, stderr, p, b)

	if exitCode != 1 {
		t.Errorf("Expected exit code 1, got %d", exitCode)
	}

	if !strings.Contains(stderr.String(), "Error building tree") {
		t.Error("Error message was not printed")
	}
}

func TestRun_MaxDepthFlag(t *testing.T) {
	// Mock dependencies
	p := &mockParser{
		parseFunc: func(input string) (*parser.Node, error) {
			return &parser.Node{Name: "project", IsDir: true}, nil
		},
	}

	b := &mockBuilder{
		buildFunc: func(root *parser.Node, maxDepth int) error {
			if maxDepth != 5 {
				t.Errorf("Expected maxDepth 5, got %d", maxDepth)
			}
			return nil
		},
	}

	// Capture stdout and stderr
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	// Test max depth flag
	exitCode := run([]string{"-d", "5", "project/"}, &bytes.Buffer{}, stdout, stderr, p, b)

	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got %d", exitCode)
	}
}

// TestMainFunction заменяем на тест, который не вызывает os.Exit
func TestMainFunctionWrapper(t *testing.T) {
	// Mock dependencies
	p := &mockParser{
		parseFunc: func(input string) (*parser.Node, error) {
			return &parser.Node{Name: "project", IsDir: true}, nil
		},
	}

	b := &mockBuilder{
		buildFunc: func(root *parser.Node, maxDepth int) error {
			return nil
		},
	}

	// Capture stdout and stderr
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	// Test main logic without calling os.Exit
	exitCode := run([]string{"-v"}, &bytes.Buffer{}, stdout, stderr, p, b)

	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got %d", exitCode)
	}

	if !strings.Contains(stdout.String(), "buildtree v"+version) {
		t.Error("Version information was not printed")
	}
}
