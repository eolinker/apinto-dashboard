package common

import (
	"encoding/json"
	"io"
)

func DecodeJSON(r io.Reader, obj any) error {
	decoder := json.NewDecoder(r)
	return decoder.Decode(obj)
}
