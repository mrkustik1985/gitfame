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

func (tf *TabularFormatter) Output(statsMap map[string]*AuthorStats) error {
	header := []string{"Name", "Lines", "Commits", "Files"}
	people := GetStats(statsMap, tf.SortOrder)
	colWidths := make([]int, len(header))
	for i, col := range header {
		colWidths[i] = len(col)
	}
	for _, person := range people {
		for i, field := range []string{person.Name, strconv.Itoa(person.Lines), strconv.Itoa(person.Commits), strconv.Itoa(person.Files)} {
			if len(field) > colWidths[i] {
				colWidths[i] = len(field)
			}
		}
	}
	for i, field := range header {
		if i != len(header)-1 {
			fmt.Printf("%-*s ", colWidths[i], field)
		} else {
			fmt.Println(field)
		}
	}
	for _, summary := range people {
		row := []string{summary.Name, strconv.Itoa(summary.Lines), strconv.Itoa(summary.Commits), strconv.Itoa(summary.Files)}
		for i, field := range row {
			if i != len(header)-1 {
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
	file := os.Stdout
	writer := csv.NewWriter(file)
	defer writer.Flush()
	err := writer.Write([]string{"Name", "Lines", "Commits", "Files"})
	if err != nil {
		return err
	}
	for _, stats := range stats {
		err := writer.Write([]string{stats.Name, strconv.Itoa(stats.Lines), strconv.Itoa(stats.Commits), strconv.Itoa(stats.Files)})
		if err != nil {
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
