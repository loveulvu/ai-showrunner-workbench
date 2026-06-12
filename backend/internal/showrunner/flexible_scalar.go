package showrunner

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
)

func unmarshalFlexibleString(data []byte) (string, error) {
	data = bytes.TrimSpace(data)
	if len(data) == 0 || bytes.Equal(data, []byte("null")) {
		return "", nil
	}

	var value string
	if err := json.Unmarshal(data, &value); err == nil {
		return value, nil
	}

	var number json.Number
	if err := json.Unmarshal(data, &number); err == nil {
		return number.String(), nil
	}
	return "", fmt.Errorf("expected string, number, or null")
}

func unmarshalFlexibleInt(data []byte) (int, error) {
	data = bytes.TrimSpace(data)
	if len(data) == 0 || bytes.Equal(data, []byte("null")) {
		return 0, nil
	}

	var value int
	if err := json.Unmarshal(data, &value); err == nil {
		return value, nil
	}

	var text string
	if err := json.Unmarshal(data, &text); err == nil {
		value, err := strconv.Atoi(text)
		if err == nil {
			return value, nil
		}
	}
	return 0, fmt.Errorf("expected integer, integer string, or null")
}
