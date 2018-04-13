package cmd

import (
	"encoding/json"
)

func jsonize(data interface{}) (s string, err error) {
	b, err := json.Marshal(data)
	if err != nil {
		return
	}
	s = string(b)
	return
}
