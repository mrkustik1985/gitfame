package parser

import "strings"

func (p *Parser) LoadTree() ([]string, error) {
	out, err := CreateCmd("git", "ls-tree", "-r", p.Scaner.Revision, "--name-only", "--full-name", ".")
	if err != nil {
		return nil, err
	}
	if len(out) == 0 {
		return nil, nil
	}
	files := strings.Split(out, "\n")
	if err != nil {
		return nil, err
	}
	extensions := SplitByDot(p.Scaner.Extensions)
	langs, err := GetAllLangs(p.Scaner.Languages)
	if err != nil {
		return nil, err
	}
	exclude := SplitByDot(p.Scaner.Exclude)
	restrictTo := SplitByDot(p.Scaner.RestrictTo)
	needFiles := make([]string, 0)
	for _, file := range files {
		if exclude != nil && IsFilenameMatchPattern(file, exclude) {
			continue
		}
		if restrictTo != nil && !IsFilenameMatchPattern(file, restrictTo) {
			continue
		}
		if !isExtensionMatch(file, extensions) {
			continue
		}
		if !isExtensionMatch(file, langs) {
			continue
		}
		needFiles = append(needFiles, file)
	}
	return needFiles, nil
}
