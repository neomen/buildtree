package parser

import (
	"errors"
	"strings"
	"unicode/utf8"

	"github.com/neomen/buildtree/internal/utils"
)

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

	if len(lines) == 0 {
		return nil, ErrEmptyInput
	}

	// Parse root directory
	rootLine := strings.TrimSpace(lines[0])
	if idx := strings.Index(rootLine, "#"); idx != -1 {
		rootLine = strings.TrimSpace(rootLine[:idx])
	}
	if strings.HasSuffix(rootLine, "/") {
		rootLine = rootLine[:len(rootLine)-1]
	}

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

	if strings.HasSuffix(name, "/") {
		isDir = true
		name = strings.TrimSuffix(name, "/")
	}

	return realLevel, name, isDir
}

func extractName(line string) string {
	name := strings.TrimSpace(line)
	prefixes := []string{"── ", "-- ", "─ ", "- ", "└──", "├──", "│", "└─", "├─", "└", "├"}
	for _, prefix := range prefixes {
		name = strings.TrimPrefix(name, prefix)
	}
	name = strings.TrimSpace(name)

	// Remove any remaining tree symbols
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

	return name
}

func normalizeTreeSymbols(input string) string {
	input = strings.ReplaceAll(input, "|--", "├──")
	input = strings.ReplaceAll(input, "'--", "└──")
	input = strings.ReplaceAll(input, "|  ", "│  ")
	input = strings.ReplaceAll(input, "|-", "├─")
	input = strings.ReplaceAll(input, "'-", "└─")
	return input
}

var ErrEmptyInput = errors.New("input is empty")
