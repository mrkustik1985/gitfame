package parser

import (
	"fmt"
	"sort"
	"strings"
)

type StatsAuthor struct {
	Name    string `json:"name"`
	Lines   int    `json:"lines"`
	Commits int    `json:"commits"`
	Files   int    `json:"files"`
}

func getSortFunction(summaries []StatsAuthor, sortOrder []string) func(i, j int) bool {
	sortByCriteria := map[string]func(i, j int) int{
		"Lines":   compareInt(summaries, func(a StatsAuthor) int { return a.Lines }),
		"Commits": compareInt(summaries, func(a StatsAuthor) int { return a.Commits }),
		"Files":   compareInt(summaries, func(a StatsAuthor) int { return a.Files }),
		"Name":    compareStr(summaries, func(a StatsAuthor) string { return a.Name }),
	}
	var sortFunctions []func(i, j int) int
	for _, criteria := range sortOrder {
		if sortFunc, ok := sortByCriteria[criteria]; ok {
			sortFunctions = append(sortFunctions, sortFunc)
		}
	}
	return func(i, j int) bool {
		for _, sortFunc := range sortFunctions {
			result := sortFunc(i, j)
			if result != 0 {
				return result > 0
			}
		}
		return false
	}
}

func compareInt(summaries []StatsAuthor, getValue func(StatsAuthor) int) func(int, int) int {
	return func(i, j int) int {
		a, b := getValue(summaries[i]), getValue(summaries[j])
		if a == b {
			return 0
		} else if a > b {
			return 1
		}
		return -1
	}
}

func compareStr(summaries []StatsAuthor, getValue func(StatsAuthor) string) func(int, int) int {
	return func(i, j int) int {
		return -strings.Compare(getValue(summaries[i]), getValue(summaries[j]))
	}
}

func GetStats(statsMap map[string]*AuthorStats, sortOrder []string) []StatsAuthor {
	var summaries []StatsAuthor
	for author, stats := range statsMap {
		summary := StatsAuthor{
			Name:    author,
			Lines:   stats.LinesCnt,
			Commits: len(stats.Commits),
			Files:   len(stats.Files),
		}
		summaries = append(summaries, summary)
	}
	sort.Slice(summaries, getSortFunction(summaries, sortOrder))
	return summaries
}

func NewFormatter(format string, orderBy string) (Formatter, error) {
	var sortOrder []string
	switch orderBy {
	case "lines":
		sortOrder = []string{"Lines", "Commits", "Files", "Name"}
	case "commits":
		sortOrder = []string{"Commits", "Lines", "Files", "Name"}
	case "files":
		sortOrder = []string{"Files", "Lines", "Commits", "Name"}
	default:
		return nil, fmt.Errorf("invalid order")
	}
	if format == "tabular" {
		return &TabularFormatter{SortOrder: sortOrder}, nil
	}
	if format == "csv" {
		return &CSVFormatter{SortOrder: sortOrder}, nil
	}
	if format == "json" {
		return &JSONFormatter{SortOrder: sortOrder}, nil
	}
	if format == "json-lines" {
		return &JSONLinesFormatter{SortOrder: sortOrder}, nil
	}
	return nil, fmt.Errorf("invalid format")
}
