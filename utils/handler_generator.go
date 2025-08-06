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
	Name     string     `json:"name"`
	Request  *Request   `json:"request,omitempty"`
	Response []Response `json:"response,omitempty"`
	Item     []Item     `json:"item,omitempty"`
}

type Request struct {
	Method string `json:"method"`
	Header []KV   `json:"header,omitempty"`
	Body   *Body  `json:"body,omitempty"`
	URL    *URL   `json:"url,omitempty"`
}

type KV struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Type  string `json:"type,omitempty"`
}

type Body struct {
	Mode string `json:"mode"`
	Raw  string `json:"raw,omitempty"`
}

type URL struct {
	Raw      string   `json:"raw"`
	Protocol string   `json:"protocol,omitempty"`
	Host     []string `json:"host,omitempty"`
	Path     []string `json:"path,omitempty"`
	Query    []KV     `json:"query,omitempty"`
}

type Response struct {
	Name            string   `json:"name"`
	OriginalRequest *Request `json:"originalRequest,omitempty"`
	Status          string   `json:"status"`
	Code            int      `json:"code"`
	Body            string   `json:"body"`
	Header          []KV     `json:"header,omitempty"`
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
