package common

import (
	"encoding/json"
	"errors"
	"github.com/xeipuuv/gojsonschema"
	"io"
)

func DecodeJSON(r io.Reader, obj any) error {
	decoder := json.NewDecoder(r)
	return decoder.Decode(obj)
}

func JsonSchemaValid(schema, data string) error {
	// 解析JSON Schema
	loader := gojsonschema.NewStringLoader(schema)
	schemaDocument, err := gojsonschema.NewSchema(loader)
	if err != nil {
		return err
	}

	// 解析待校验的JSON数据
	documentLoader := gojsonschema.NewStringLoader(data)
	result, err := schemaDocument.Validate(documentLoader)
	if err != nil {
		return err
	}
	if result.Valid() {
		return nil
	}
	return errors.New(result.Errors()[0].String())
}
