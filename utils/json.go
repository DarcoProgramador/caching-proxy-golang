package utils

import "encoding/json"

func UnmarshalJSONBinary(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func MarshalJSONBinary(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}
