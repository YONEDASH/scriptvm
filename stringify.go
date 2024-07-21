package script

import (
	"encoding/json"
	"fmt"
)

func Stringify(v any) string {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return stringifyErrorMessage(err.Error())
	}
	return string(data)
}

func stringifyErrorMessage(msg string) string {
	return fmt.Sprintf(`{"error": "%s"}`, msg)
}
