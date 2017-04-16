package util

import (
	"encoding/json"
	"io"
)

// DecodeJSON decodes json
func DecodeJSON(src io.Reader, dst interface{}) error {
	decoder := json.NewDecoder(src)
	err := decoder.Decode(dst)
	return err
}

// EncodeJSON encodes json
func EncodeJSON(dst io.Writer, src interface{}) error {
	encoder := json.NewEncoder(dst)
	err := encoder.Encode(src)
	return err
}

func WriteErrorJSON(dst io.Writer, errorMsg string) error {
	response := make(map[string]interface{})
	response["ok"] = false
	response["failureReason"] = errorMsg
	return EncodeJSON(dst, response)
}
