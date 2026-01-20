package http

import "encoding/json"

// toJSON converts a Go value to a JSON string
func toJSON(v interface{}) (string, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
