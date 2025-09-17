package postman

import (
	"encoding/json"
	"errors"
	"os"
)

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

func GenerateHandler(ps PostmanCollection) []Handler {
	var hs []Handler
	eHs := make(map[string]bool)

	for i := range ps.Item {
		for j := range ps.Item[i].Item {
			handlerName := ps.Item[i].Name + ps.Item[i].Item[j].Name
			if _, ok := eHs[handlerName]; ok {
				continue
			}

			eHs[handlerName] = true

			var handler Handler
			handler.SetName(handlerName)
			handler.SetMethod(ps.Item[i].Item[j].Request.Method)
			handler.SetPathVariable(ps.Item[i].Item[j].Request.Url.Variable)

			hs = append(hs, handler)
		}
	}

	return hs
}

func parsePostmanFile(input string) ([]byte, error) {
	return os.ReadFile(input)
}
