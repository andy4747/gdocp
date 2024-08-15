package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/andy4747/gdocp/internal/markdown"
)

// MarkdownContent contains the Markdown content and associated metadata.
type MarkdownContent struct {
	Path    string
	Content string
}

// GenerateMarkdownRecursively processes all Go files in a directory recursively and stores generated markdown content in a map.
func GenerateMarkdownRecursively(baseDir string) (map[string]MarkdownContent, error) {
	result := make(map[string]MarkdownContent)
	err := filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories and test files
		if info.IsDir() || strings.HasSuffix(info.Name(), "_test.go") {
			return nil
		}

		// Process only Go files
		if strings.HasSuffix(info.Name(), ".go") {
			fileInfo, err := ParseFile(path)
			if err != nil {
				return fmt.Errorf("failed to parse file %s: %v", path, err)
			}

			if fileInfo.Author != "" && fileInfo.File != "" {
				mdContent, err := markdown.GenerateMarkdown(fileInfo)
				if err != nil {
					return fmt.Errorf("failed to generate markdown for file %s: %v", path, err)
				}

				// Skip files with empty content (due to empty Author or File fields)
				if mdContent == "" {
					fmt.Printf("Skipped generating markdown for %s: Author or File is empty\n", path)
					return nil
				}

				// Generate a UUID for the markdown content
				result[path] = MarkdownContent{
					Path:    path,
					Content: mdContent,
				}

				fmt.Printf("Markdown for file %s\n", path)
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}
