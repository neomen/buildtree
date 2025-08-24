package parser

import (
	"errors"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/neomen/buildtree/internal/utils"
)

var ErrEmptyInput = errors.New("input is empty")

// Node represents a file or directory in the tree
type Node struct {
	Name     string
	IsDir    bool
	Level    int
	Children []*Node
}

// ParseInput converts text input to a tree structure
func ParseInput(input string) (*Node, error) {
	input = normalizeTreeSymbols(input)
	lines := strings.Split(input, "\n")

	// Check that the input is not empty after processing
	if strings.TrimSpace(input) == "" {
		return nil, ErrEmptyInput
	}

	if len(lines) == 0 {
		return nil, ErrEmptyInput
	}

	// Parse root directory
	rootLine := strings.TrimSpace(lines[0])
	if idx := strings.Index(rootLine, "#"); idx != -1 {
		rootLine = strings.TrimSpace(rootLine[:idx])
	}
	rootLine = strings.TrimSuffix(rootLine, "/")

	root := &Node{
		Name:  rootLine,
		IsDir: true,
		Level: 0,
	}

	stack := []*Node{root}
	prevLevel := 0

	for _, line := range lines[1:] {
		line = strings.TrimRight(line, " ")
		if strings.TrimSpace(line) == "" {
			continue
		}

		// Remove comments
		if idx := strings.Index(line, "#"); idx != -1 {
			line = line[:idx]
		}

		level, name, isDir := parseLine(line)
		if name == "" {
			continue
		}

		// Adjust stack based on level
		if level <= prevLevel {
			stack = stack[:level+1]
		} else if level > len(stack)-1 {
			// Handle skipped levels
			for len(stack) <= level {
				stack = append(stack, stack[len(stack)-1])
			}
		}

		parent := stack[level]
		node := &Node{
			Name:  name,
			IsDir: isDir,
			Level: level,
		}

		parent.Children = append(parent.Children, node)

		if isDir {
			// Add to stack for children
			if level+1 < len(stack) {
				stack[level+1] = node
			} else {
				stack = append(stack, node)
			}
		}

		prevLevel = level
	}

	return root, nil
}

func parseLine(line string) (level int, name string, isDir bool) {
	// Deleting comments at the beginning
	if idx := strings.Index(line, "#"); idx != -1 {
		line = line[:idx]
	}

	line = strings.TrimRight(line, " ")
	if strings.TrimSpace(line) == "" {
		return 0, "", false
	}

	remaining := line
	for len(remaining) > 0 {
		r, size := utf8.DecodeRuneInString(remaining)
		if utils.IsTreeSymbol(r) {
			level++
			remaining = remaining[size:]
		} else {
			break
		}
	}

	realLevel := level / 4
	name = extractName(remaining)

	// Checking if this is a directory
	if strings.HasSuffix(name, "/") {
		isDir = true
		name = strings.TrimSuffix(name, "/")
	} else {
		// Checking if the name has an extension
		// If there is no dot in the name, it is possible that it is a directory without a slash
		// This is a heuristic, it can be improved
		if !strings.Contains(name, ".") && !strings.Contains(name, string(os.PathSeparator)) {
			isDir = true
		}
	}

	return realLevel, name, isDir
}

func extractName(line string) string {
	name := strings.TrimSpace(line)

	// Removing all possible prefixes of tree elements
	prefixes := []string{"── ", "-- ", "─ ", "- ", "└──", "├──", "│", "└─", "├─", "└", "├"}
	for _, prefix := range prefixes {
		name = strings.TrimPrefix(name, prefix)
	}
	name = strings.TrimSpace(name)

	// Remove any remaining tree characters from the beginning of the name
	for strings.HasPrefix(name, "├") || strings.HasPrefix(name, "└") ||
		strings.HasPrefix(name, "│") || strings.HasPrefix(name, "─") ||
		strings.HasPrefix(name, "|") || strings.HasPrefix(name, "-") {
		if len(name) > 1 {
			name = name[1:]
		} else {
			name = ""
			break
		}
		name = strings.TrimSpace(name)
	}

	// Deleting comments at the end
	if idx := strings.Index(name, "#"); idx != -1 {
		name = strings.TrimSpace(name[:idx])
	}

	return name
}

// Bringing different types of tree to a uniform condition
func normalizeTreeSymbols(input string) string {
	input = strings.ReplaceAll(input, "|--", "├──")
	input = strings.ReplaceAll(input, "'--", "└──")
	input = strings.ReplaceAll(input, "|  ", "│  ")
	input = strings.ReplaceAll(input, "|-", "├─")
	input = strings.ReplaceAll(input, "'-", "└─")
	return input
}
