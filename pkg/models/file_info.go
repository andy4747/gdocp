package models

// FileInfo represents the parsed information from a comments in Go file
type FileInfo struct {
	Author          string
	Date            string
	File            string
	Problem         string
	Solution        string
	Note            string
	TimeComplexity  string
	SpaceComplexity string
	Code            string
}

// ToMap converts the FileInfo struct into a map with string keys and values.
func (f *FileInfo) ToMap() map[string]string {
	return map[string]string{
		"Author":          f.Author,
		"Date":            f.Date,
		"File":            f.File,
		"Problem":         f.Problem,
		"Solution":        f.Solution,
		"Note":            f.Note,
		"TimeComplexity":  f.TimeComplexity,
		"SpaceComplexity": f.SpaceComplexity,
		"Code":            f.Code,
	}
}
