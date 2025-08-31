package exporter

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// ExportStructure saves the current directory structure with file contents to the specified file
// with filtering options
func ExportStructure(outputFile, rootDir string, filters []string, ignoreDirs []string, maxSize int64, includeHidden bool) error {
	// Convert relative path to absolute
	rootDirAbs, err := filepath.Abs(rootDir)
	if err != nil {
		return fmt.Errorf("failed to get absolute path for root directory: %w", err)
	}

	// Check the existence of the root directory BEFORE creating the output file
	info, err := os.Stat(rootDirAbs)
	if os.IsNotExist(err) {
		return fmt.Errorf("root directory does not exist: %s", rootDirAbs)
	}
	if err != nil {
		return fmt.Errorf("error checking root directory: %w", err)
	}
	if !info.IsDir() {
		return fmt.Errorf("root path is not a directory: %s", rootDirAbs)
	}

	// Get absolute path of output file
	outputFileAbs, err := filepath.Abs(outputFile)
	if err != nil {
		return fmt.Errorf("failed to get absolute path for output file: %w", err)
	}

	// Create or truncate the output file
	f, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer f.Close()

	// Default ignore directories
	defaultIgnore := []string{".git", "node_modules"}
	if len(ignoreDirs) == 0 {
		ignoreDirs = defaultIgnore
	} else {
		// Add default ignore directories to user-provided ones
		ignoreDirs = append(ignoreDirs, defaultIgnore...)
	}

	// Walk through the directory
	return filepath.WalkDir(rootDirAbs, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip the output file itself to avoid recursion
		pathAbs, err := filepath.Abs(path)
		if err != nil {
			return fmt.Errorf("failed to get absolute path for %s: %w", path, err)
		}
		if filepath.Clean(pathAbs) == filepath.Clean(outputFileAbs) {
			return nil
		}

		// Check if directory should be ignored
		if d.IsDir() {
			dirName := d.Name()

			// Check for hidden directories (unless includeHidden is true)
			if !includeHidden && strings.HasPrefix(dirName, ".") {
				return fs.SkipDir
			}

			for _, ignore := range ignoreDirs {
				if dirName == ignore {
					return fs.SkipDir
				}
			}
		}

		// Process only files
		if !d.IsDir() {
			fileName := d.Name()

			// Check for hidden files (unless includeHidden is true)
			if !includeHidden && strings.HasPrefix(fileName, ".") {
				return nil
			}

			// Check file size
			info, err := d.Info()
			if err != nil {
				return fmt.Errorf("failed to get file info for %s: %w", path, err)
			}

			if maxSize > 0 && info.Size() > maxSize {
				return nil
			}

			// Apply filters if specified
			if len(filters) > 0 {
				ext := strings.TrimPrefix(filepath.Ext(d.Name()), ".")
				matched := false
				for _, filter := range filters {
					if strings.EqualFold(ext, filter) {
						matched = true
						break
					}
				}
				if !matched {
					return nil
				}
			}

			// Get relative path
			relPath, err := filepath.Rel(rootDirAbs, path)
			if err != nil {
				return fmt.Errorf("failed to get relative path for %s: %w", path, err)
			}

			// Replacing backslashes with straight ones
			relPath = filepath.ToSlash(relPath)

			// Read file content
			content, err := os.ReadFile(path)
			if err != nil {
				return fmt.Errorf("failed to read file %s: %w", path, err)
			}

			// Write header and content
			if _, err := fmt.Fprintf(f, "# %s\n", relPath); err != nil {
				return err
			}
			if _, err := f.Write(content); err != nil {
				return err
			}
			if _, err := f.WriteString("\n\n"); err != nil {
				return err
			}
		}
		return nil
	})
}
