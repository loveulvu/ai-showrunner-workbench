package showrunner

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type FlexibleStringList []string

func (list *FlexibleStringList) UnmarshalJSON(data []byte) error {
	data = bytes.TrimSpace(data)
	if bytes.Equal(data, []byte("null")) {
		*list = FlexibleStringList{}
		return nil
	}

	var values []string
	if err := json.Unmarshal(data, &values); err == nil {
		if values == nil {
			values = []string{}
		}
		*list = values
		return nil
	}

	var value string
	if err := json.Unmarshal(data, &value); err == nil {
		*list = FlexibleStringList{value}
		return nil
	}

	return fmt.Errorf("expected string, array of strings, or null")
}

func (list FlexibleStringList) Text() string {
	if len(list) == 0 {
		return ""
	}
	return list[0]
}
