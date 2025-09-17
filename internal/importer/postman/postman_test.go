package postman

import (
	"encoding/json"
	"testing"
)

func TestParseAPI(t *testing.T) {
	PostManCollection := "goarm.json"

	t.Run("parse postman collection", func(t *testing.T) {
		fileData, err := parsePostmanFile(PostManCollection)
		if err != nil {
			t.Fatalf("failed to parse file %s: %v", PostManCollection, err)
		}
		t.Logf("Successfully read Postman file: %s (size: %d bytes)", PostManCollection, len(fileData))

		var pc PostmanCollection
		if err := json.Unmarshal(fileData, &pc); err != nil {
			t.Fatalf("failed to unmarshal JSON into PostmanCollection: %v", err)
		}
		t.Log(pc)

		handlers := GenerateHandler(pc)
		t.Logf("Generated %d handlers:", len(handlers))
		for i, h := range handlers {
			t.Logf(
				"Handler %d -> Name: %s, Method: %s, Fields: %+v, PathVariables: %+v",
				i, h.GetName(), h.GetMethod(), h.GetFields(), h.GetPathVariable(),
			)
		}
	})
}
