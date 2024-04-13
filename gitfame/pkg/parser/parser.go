package parser

import (
	"bufio"
	"bytes"
	"fmt"
	"gitlab.com/slon/shad-go/gitfame/configs"
	"gitlab.com/slon/shad-go/gitfame/pkg/scaner"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type AuthorStats struct {
	Commits  map[string]bool
	Files    map[string]bool
	LinesCnt int
}

func NewAuthorStats() *AuthorStats {
	return &AuthorStats{
		Commits:  make(map[string]bool),
		Files:    make(map[string]bool),
		LinesCnt: 0,
	}
}

type Parser struct {
	Scaner *scaner.Scaner
	Stats  map[string]*AuthorStats // key - author name, value - author stats
}

func NewParser(scan *scaner.Scaner) *Parser {
	return &Parser{
		Scaner: scan,
		Stats:  make(map[string]*AuthorStats),
	}
}

func CreateCmd(name string, args ...string) (string, error) {
	var stdout, stderr bytes.Buffer
	cmd := exec.Command(name, args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf(stderr.String())
	}
	return strings.TrimSpace(stdout.String()), nil
}

func SplitByDot(s string) []string {
	if s == "" {
		return nil
	}
	return strings.Split(s, ",")
}

func GetAllLangs(lgs string) ([]string, error) {
	allLang, err := configs.ParseLangs()
	if err != nil {
		return nil, err
	}
	langs := SplitByDot(lgs)
	langLower := make(map[string]bool)
	for _, lang := range langs {
		langLower[strings.ToLower(lang)] = true
	}
	retLangs := make([]string, 0)
	for _, lang := range allLang {
		name := strings.ToLower(lang.Name)
		if langLower[name] {
			retLangs = append(retLangs, lang.Extensions...)
		}
	}
	return retLangs, nil
}

func IsFilenameMatchPattern(filename string, patterns []string) bool {
	for _, pattern := range patterns {
		matched, _ := regexp.MatchString(pattern, filename)
		if matched {
			return true
		}
	}
	return false
}

func isExtensionMatch(filename string, expectedExts []string) bool {
	for _, expected := range expectedExts {
		if filepath.Ext(filename) == expected {
			return true
		}
	}
	return len(expectedExts) == 0
}

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

func (p *Parser) ParseLastCommiter(file string) error {
	out, err := CreateCmd("git", "log", "-1", "--pretty=format:%H,%an", p.Scaner.Revision, "--", file)
	if err != nil {
		return err
	}
	commitInfo := strings.SplitN(out, ",", 2)

	//log.Println(commitInfo)
	if len(commitInfo) != 2 {
		return fmt.Errorf("unexpected output format: %s", out)
	}
	commitID, author := commitInfo[0], commitInfo[1]
	if _, ok := p.Stats[author]; !ok {
		p.Stats[author] = NewAuthorStats()
	}
	p.Stats[author].Files[file] = true
	p.Stats[author].Commits[commitID] = true
	return nil
}

func (p *Parser) ParseFile(file string) error {
	out, err := CreateCmd("git", "blame", p.Scaner.Revision, "--porcelain", file)
	if err != nil {
		return err
	}
	prefixStart := "author"
	if p.Scaner.UseCommitter {
		prefixStart = "committer"
	}
	scanner := bufio.NewScanner(strings.NewReader(out))
	fl := 0
	authorByCommit := make(map[string]string)
	for scanner.Scan() {
		line := strings.Split(scanner.Text(), " ")
		commit := line[0]
		linesCnt, _ := strconv.Atoi(line[len(line)-1])
		if author, ok := authorByCommit[commit]; ok {
			p.Stats[author].LinesCnt += linesCnt
		}
		i := 0
		for i < linesCnt {
			if !scanner.Scan() {
				break
			}
			line := scanner.Text()
			lineCopy := strings.Split(line, " ")
			if strings.HasPrefix(lineCopy[0], "\t") {
				i++
			} else {
				if lineCopy[0] == prefixStart {
					author := strings.TrimSpace(strings.TrimPrefix(line, prefixStart+" "))
					fl++
					if _, ok := authorByCommit[commit]; !ok {
						authorByCommit[commit] = author
					}
					if _, ok := p.Stats[author]; !ok {
						p.Stats[author] = NewAuthorStats()
					}
					p.Stats[author].Files[file] = true
					p.Stats[author].Commits[commit] = true
					p.Stats[author].LinesCnt += linesCnt
				}
			}
		}
	}
	if fl == 0 {
		err = p.ParseLastCommiter(file)
	}
	return err
}

func (p *Parser) ParseFiles(files []string) error {
	for _, file := range files {
		err := p.ParseFile(file)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Parser) DoRoutine() error {
	err := os.Chdir(p.Scaner.Repository)
	if err != nil {
		return err
	}
	files, err := p.LoadTree()
	if err != nil {
		return err
	}
	err = p.ParseFiles(files)
	if err != nil {
		return err
	}
	return nil
}
