package parser

import (
	"gitlab.com/slon/shad-go/gitfame/configs"
	"strings"
)

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
