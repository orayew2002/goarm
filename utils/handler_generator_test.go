package utils

import (
	"encoding/json"
	"testing"
)

func TestParseAPI(t *testing.T) {
	t.Run("parse json file", func(t *testing.T) {
		if _, err := parsePostmanFile("../media.json"); err != nil {
			t.Fatal("error when parsing file:", err)
		}
	})

	t.Run("parse struct", func(t *testing.T) {
		file, err := parsePostmanFile("../media.json")
		if err != nil {
			t.Fatal("error when parsing file:", err)
		}

		var pc PostmanCollection
		if err := json.Unmarshal(file, &pc); err != nil {
			t.Fatal("error parse postman struct", err)
		}

		t.Log(pc)
	})
}
