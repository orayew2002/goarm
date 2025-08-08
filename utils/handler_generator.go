package utils

import (
	"encoding/json"
	"errors"
	"os"
)

type PostmanCollection struct {
	Item []Item `json:"item"`
}

type Item struct {
	Name    string  `json:"name"`
	Item    []Item  `json:"item,omitempty"`
	Request Request `json:"request,omitempty"`
}

type Request struct {
	Method string `json:"method"`
	Url    Url    `json:"url,omitempty"`
}

type Url struct {
	Raw string `json:"raw"`
}

func ParseApi(input string) error {
	file, err := parsePostmanFile(input)
	if err != nil {
		return errors.Join(errors.New("error when parse file"), err)
	}

	var ps PostmanCollection
	if err := json.Unmarshal(file, &ps); err != nil {
		return err
	}
	return nil
}

func parsePostmanFile(input string) ([]byte, error) {
	return os.ReadFile(input)
}
