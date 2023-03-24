package common

import (
	"encoding/json"
	"github.com/xeipuuv/gojsonschema"
	"io"
)

func DecodeJSON(r io.Reader, obj any) error {
	decoder := json.NewDecoder(r)
	return decoder.Decode(obj)
}

func JsonSchemaValid(schema, data string) bool {
	// 解析JSON Schema
	loader := gojsonschema.NewStringLoader(schema)
	schemaDocument, err := gojsonschema.NewSchema(loader)
	if err != nil {
		return false
	}

	// 解析待校验的JSON数据
	documentLoader := gojsonschema.NewStringLoader(data)
	result, err := schemaDocument.Validate(documentLoader)
	if err != nil {
		return false
	}
	return result.Valid()
}
