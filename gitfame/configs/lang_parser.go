package configs

import "encoding/json"
import _ "embed"

//go:embed language_extensions.json
var languagesJSON string

type Language struct {
	Name       string   `json:"name"`
	Type       string   `json:"type"`
	Extensions []string `json:"extensions"`
}

func ParseLangs() ([]Language, error) {
	var lang []Language
	err := json.Unmarshal([]byte(languagesJSON), &lang)
	if err != nil {
		return nil, err
	}
	return lang, nil
}
