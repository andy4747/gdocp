package parser

import (
	"bufio"
	"os"
	"regexp"
	"strings"

	"github.com/andy4747/gdocp/pkg/models"
)

// ParseFile extracts structured information from the top-level comment of a Go file.
func ParseFile(filename string) (*models.FileInfo, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	fileInfo := &models.FileInfo{}
	inMultilineComment := false
	codeStart := false

	// Regular expressions to match each field in the comment
	reMap := map[string]*regexp.Regexp{
		"author":          regexp.MustCompile(`(?i)Author:\s*(.+)`),
		"date":            regexp.MustCompile(`(?i)Date:\s*(.+)`),
		"file":            regexp.MustCompile(`(?i)File:\s*(.+)`),
		"problem":         regexp.MustCompile(`(?i)Problem:\s*(.+)`),
		"solutionStart":   regexp.MustCompile(`(?i)Solution:\s*\{`),
		"solutionEnd":     regexp.MustCompile(`\}`),
		"noteStart":       regexp.MustCompile(`(?i)Note:\s*\{`),
		"noteEnd":         regexp.MustCompile(`\}`),
		"timeComplexity":  regexp.MustCompile(`(?i)Time Complexity:\s*(.+)`),
		"spaceComplexity": regexp.MustCompile(`(?i)Space Complexity:\s*(.+)`),
	}

	var currentField *string

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if strings.HasPrefix(line, "/*") {
			inMultilineComment = true
			continue
		}
		if strings.HasPrefix(line, "*/") {
			inMultilineComment = false
			codeStart = true
			continue
		}

		if inMultilineComment {
			// Match each line against the regex map
			switch {
			case reMap["author"].MatchString(line):
				fileInfo.Author = reMap["author"].FindStringSubmatch(line)[1]
			case reMap["date"].MatchString(line):
				fileInfo.Date = reMap["date"].FindStringSubmatch(line)[1]
			case reMap["file"].MatchString(line):
				fileInfo.File = reMap["file"].FindStringSubmatch(line)[1]
			case reMap["problem"].MatchString(line):
				fileInfo.Problem = reMap["problem"].FindStringSubmatch(line)[1]
			case reMap["solutionStart"].MatchString(line):
				currentField = &fileInfo.Solution
			case reMap["solutionEnd"].MatchString(line):
				currentField = nil
			case reMap["noteStart"].MatchString(line):
				currentField = &fileInfo.Note
			case reMap["noteEnd"].MatchString(line):
				currentField = nil
			case reMap["timeComplexity"].MatchString(line):
				fileInfo.TimeComplexity = reMap["timeComplexity"].FindStringSubmatch(line)[1]
			case reMap["spaceComplexity"].MatchString(line):
				fileInfo.SpaceComplexity = reMap["spaceComplexity"].FindStringSubmatch(line)[1]
			}

			// Capture multiline content for solution and note
			if currentField != nil {
				*currentField += line + "\n"
			}
		} else if codeStart {
			// Capture the remaining code after the comment block
			fileInfo.Code += line + "\n"
		}
	}

	// Clean up the Solution and Note fields to remove leading identifiers
	fileInfo.Solution = strings.TrimSpace(strings.TrimSuffix(fileInfo.Solution, "}"))
	fileInfo.Solution = strings.TrimSpace(strings.TrimPrefix(fileInfo.Solution, "Solution: {"))
	fileInfo.Solution = strings.TrimSpace(strings.TrimPrefix(fileInfo.Solution, "Solution:{"))
	fileInfo.Note = strings.TrimSpace(strings.TrimSuffix(fileInfo.Note, "}"))
	fileInfo.Note = strings.TrimSpace(strings.TrimPrefix(fileInfo.Note, "Note: {"))
	fileInfo.Note = strings.TrimSpace(strings.TrimPrefix(fileInfo.Note, "Note:{"))

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return fileInfo, nil
}
