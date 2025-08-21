package builder

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/neomen/buildtree/internal/parser"
	"github.com/neomen/buildtree/internal/validator"
)

// BuildTree creates the file structure from the parsed tree
func BuildTree(root *parser.Node, maxDepth int) error {
	if maxDepth < 0 {
		maxDepth = 0
		log.Println("Warning: Negative max-depth value corrected to 0 (no limit)")
	}

	// Validate root node name
	if !validator.IsValidPath(root.Name) {
		return fmt.Errorf("invalid root node name: '%s'", root.Name)
	}

	// Create root directory
	if err := createNode(root, "", maxDepth, 0); err != nil {
		return err
	}

	if maxDepth > 0 {
		log.Printf("Created structure with max depth %d", maxDepth)
	}
	return nil
}

func createNode(node *parser.Node, parentPath string, maxDepth, currentDepth int) error {
	fullPath := filepath.Join(parentPath, node.Name)

	// Check max depth
	if maxDepth > 0 && currentDepth > maxDepth {
		log.Printf("Skipping '%s' - exceeds max depth (%d)", fullPath, maxDepth)
		return nil
	}

	// Validate path
	if !validator.IsValidPath(node.Name) {
		log.Printf("Invalid name '%s' - skipping", node.Name)
		return nil
	}

	if node.IsDir {
		// Create directory
		if err := os.MkdirAll(fullPath, 0755); err != nil && !os.IsExist(err) {
			return err
		}

		// Process children
		for _, child := range node.Children {
			if err := createNode(child, fullPath, maxDepth, currentDepth+1); err != nil {
				return err
			}
		}
	} else {
		// Create file
		if err := os.WriteFile(fullPath, []byte{}, 0644); err != nil {
			return err
		}
	}

	return nil
}
