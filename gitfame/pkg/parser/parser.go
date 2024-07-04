package parser

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

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
