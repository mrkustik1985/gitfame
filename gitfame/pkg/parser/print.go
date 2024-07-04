package parser

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

type Formatter interface {
	Output(map[string]*AuthorStats) error
}

type TabularFormatter struct {
	SortOrder []string
}
type Column struct {
	Header string
	Getter func(person StatsAuthor) string
}

func (tf *TabularFormatter) Output(statsMap map[string]*AuthorStats) error {
	columns := []Column{
		{"Name", func(p StatsAuthor) string { return p.Name }},
		{"Lines", func(p StatsAuthor) string { return strconv.Itoa(p.Lines) }},
		{"Commits", func(p StatsAuthor) string { return strconv.Itoa(p.Commits) }},
		{"Files", func(p StatsAuthor) string { return strconv.Itoa(p.Files) }},
	}

	people := GetStats(statsMap, tf.SortOrder)

	colWidths := make([]int, len(columns))
	for i, col := range columns {
		colWidths[i] = len(col.Header)
	}

	for _, person := range people {
		for i, col := range columns {
			field := col.Getter(person)
			if len(field) > colWidths[i] {
				colWidths[i] = len(field)
			}
		}
	}

	for i, col := range columns {
		if i != len(columns)-1 {
			fmt.Printf("%-*s ", colWidths[i], col.Header)
		} else {
			fmt.Println(col.Header)
		}
	}

	for _, person := range people {
		for i, col := range columns {
			field := col.Getter(person)
			if i != len(columns)-1 {
				fmt.Printf("%-*s ", colWidths[i], field)
			} else {
				fmt.Println(field)
			}
		}
	}

	return nil
}

type CSVFormatter struct {
	SortOrder []string
}

func (cf *CSVFormatter) Output(statsMap map[string]*AuthorStats) error {
	stats := GetStats(statsMap, cf.SortOrder)
	writer := csv.NewWriter(os.Stdout)
	defer writer.Flush()

	headers := []string{"Name", "Lines", "Commits", "Files"}
	if err := writer.Write(headers); err != nil {
		return err
	}

	for _, stat := range stats {
		row := []string{
			stat.Name,
			strconv.Itoa(stat.Lines),
			strconv.Itoa(stat.Commits),
			strconv.Itoa(stat.Files),
		}
		if err := writer.Write(row); err != nil {
			return err
		}
	}

	return nil
}

type JSONFormatter struct {
	SortOrder []string
}

func (jf *JSONFormatter) Output(statsMap map[string]*AuthorStats) error {
	summaries := GetStats(statsMap, jf.SortOrder)
	jsonData, err := json.Marshal(summaries)
	if err != nil {
		return err
	}
	fmt.Println(string(jsonData))
	return nil
}

type JSONLinesFormatter struct {
	SortOrder []string
}

func (jlf *JSONLinesFormatter) Output(statsMap map[string]*AuthorStats) error {
	people := GetStats(statsMap, jlf.SortOrder)
	for _, person := range people {
		jsonData, err := json.Marshal(person)
		if err != nil {
			return err
		}
		//log.Println(string(jsonData))
		fmt.Println(string(jsonData))
	}
	return nil
}
