package config

import (
	_ "embed"
	"gopkg.in/yaml.v3"
)

//go:embed content-type.yml
var contentTypeData []byte
var (
	contentTypeList []ContentTypeItem
)

type ContentTypeItem struct {
	ContentType string `json:"content_type" yaml:"content_type"`
	Body        string `json:"body" yaml:"body"`
}

func init() {
	conf := make([]ContentTypeItem, 0, 3)
	err := yaml.Unmarshal(contentTypeData, &conf)
	if err != nil {
		panic(err)
	}
	for _, item := range conf {
		if item.ContentType == "" {
			panic("content-type can't be nil. ")
		}
		if item.Body == "" {
			panic("body can't be nil. ")
		}
	}
	contentTypeList = conf
}

func GetContentTypeList() []ContentTypeItem {
	return contentTypeList
}
