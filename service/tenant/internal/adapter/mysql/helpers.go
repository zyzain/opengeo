package mysql

import (
	"encoding/json"
)

// marshalJSON 将 map 序列化为 JSON 字符串
func marshalJSON(data map[string]string) (string, error) {
	if data == nil {
		return "{}", nil
	}
	bytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// unmarshalJSON 将 JSON 字符串反序列化为 map
func unmarshalJSON(jsonStr string) (map[string]string, error) {
	if jsonStr == "" || jsonStr == "{}" {
		return make(map[string]string), nil
	}
	var result map[string]string
	err := json.Unmarshal([]byte(jsonStr), &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
