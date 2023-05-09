package common

import (
	"gopkg.in/yaml.v3"
	io "io"
)

func DecodeYAML(r io.Reader, obj any) error {
	decoder := yaml.NewDecoder(r)
	return decoder.Decode(obj)
}
