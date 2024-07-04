package parser

import (
	"bytes"
	"fmt"
	"gitlab.com/slon/shad-go/gitfame/pkg/scaner"
	"os/exec"
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
