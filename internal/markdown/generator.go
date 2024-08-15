package markdown

import (
	"fmt"
	"go/format"
	"os"

	"github.com/andy4747/gdocp/pkg/models"
)

// GenerateMarkdown creates markdown from the parsed file information.
func GenerateMarkdown(fileInfo *models.FileInfo) (string, error) {
	code := []byte(fileInfo.Code)
	code, err := format.Source(code)
	if err != nil {
		return "", err
	}
	codeBlock := fmt.Sprintf("```go\n%s\n```", code)
	fmt.Println(codeBlock)
	return fmt.Sprintf(`# Notes

**Author**: %s
**Date**: %s
**File**: %s

## Problem
%s

## Solution
%s

**Time Complexity**: %s
**Space Complexity**: %s

## Note
%s

## Code
%s
`, fileInfo.Author, fileInfo.Date, fileInfo.File, fileInfo.Problem, fileInfo.Solution, fileInfo.TimeComplexity, fileInfo.SpaceComplexity, fileInfo.Note, codeBlock), nil
}

// WriteToFile writes the markdown content to a file.
func WriteToFile(content, filename string) error {
	return os.WriteFile(filename, []byte(content), 0644)
}
