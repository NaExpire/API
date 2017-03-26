package business

import (
	"encoding/json"
	"io"
)

func decodeJSON(src io.Reader, dst interface{}) error {
	decoder := json.NewDecoder(src)
	err := decoder.Decode(dst)
	return err
}

func encodeJSON(dst io.Writer, src interface{}) error {
	encoder := json.NewEncoder(dst)
	err := encoder.Encode(src)
	return err
}
