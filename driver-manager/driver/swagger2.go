package driver

import (
	"bytes"
	"github.com/eolinker/apinto-dashboard/common"
	"github.com/eolinker/apinto-dashboard/entry"
	"github.com/eolinker/apinto-dashboard/enum"
	"github.com/eolinker/apinto-dashboard/model"
	"github.com/go-basic/uuid"
	"net/url"
	"strings"
	"time"
)

type OpenAPI2 struct {
	format string
}

func (o *OpenAPI2) FormatAPI(data []byte, namespaceID, appID int, groupID, prefix, label string) ([]*model.APIInfo, error) {
	openAPI2Config := new(OpenAPI2Config)
	reader := bytes.NewReader(data)
	if err := common.DecodeYAML(reader, openAPI2Config); err != nil {
		reader = bytes.NewReader(data)
		err = common.DecodeJSON(reader, openAPI2Config)
		if err != nil {
			return nil, err
		}
	}

	//格式化requestPrefix
	if prefix != "" {
		prefix = "/" + strings.Trim(prefix, "/")
	}

	apiList := make([]*model.APIInfo, 0)
	t := time.Now()
	for path, pathMap := range openAPI2Config.Paths {
		//对路径进行转义
		decodedPath, err := url.PathUnescape(path)
		if err != nil {
			decodedPath = path
		}
		for method, info := range pathMap {
			name := info.Summary
			if name == "" {
				name = info.OperationID
			}
			if name == "" {
				name = method + "-" + prefix + decodedPath
			}
			apiList = append(apiList, &model.APIInfo{
				API: &entry.API{
					NamespaceId:      namespaceID,
					UUID:             uuid.New(),
					GroupUUID:        groupID,
					Name:             name,
					RequestPath:      common.ReplaceRestfulPath(prefix+decodedPath, enum.RestfulLabel),
					RequestPathLabel: prefix + decodedPath,
					SourceType:       enum.SourceSync,
					SourceID:         appID,
					SourceLabel:      label,
					Desc:             info.Description,
					Operator:         0,
					CreateTime:       t,
					UpdateTime:       t,
				},
				Method: []string{strings.ToUpper(method)},
			})
		}
	}
	return apiList, nil
}

func CreateOpenAPI2(format string) IAPISyncFormatDriver {
	return &OpenAPI2{format: format}
}

type OpenAPI2Config struct {
	Swagger  string                              `yaml:"swagger" json:"swagger"`
	Info     *OpenAPI2Info                       `yaml:"info" json:"info"`
	Host     string                              `yaml:"host" json:"host"`
	BasePath string                              `yaml:"basePath" json:"basePath"`
	Tags     []OpenAPI2Tags                      `yaml:"tags" json:"tags"`
	Schemes  []string                            `yaml:"schemes" json:"schemes"`
	Paths    map[string]map[string]*OpenAPI2Path `yaml:"paths" json:"paths"`
}

type OpenAPI2Info struct {
	Description    string `yaml:"description" json:"description"`
	Version        string `yaml:"version" json:"version"`
	Title          string `yaml:"title"  json:"title"`
	TermsOfService string `yaml:"termsOfService" json:"termsOfService"`
}

type OpenAPI2Tags struct {
	Name         string                   `yaml:"name" json:"name"`
	Description  string                   `yaml:"description" json:"description"`
	ExternalDocs OpenAPI2TagsExternalDocs `yaml:"externalDocs" json:"externalDocs"`
}

type OpenAPI2TagsExternalDocs struct {
	Description string `yaml:"description" json:"description"`
	URL         string `yaml:"url" json:"url"`
}

type OpenAPI2Path struct {
	Tags        []string `yaml:"tags" yaml:"tags"`
	Summary     string   `yaml:"summary" yaml:"summary"`
	Description string   `yaml:"description" yaml:"description"`
	OperationID string   `yaml:"operationId" json:"operationId"`
	//Consumes    []string                        `yaml:"consumes" json:"consumes"`
	//Produces    []string                        `yaml:"produces" json:"produces"`
	//Parameters  []OpenAPI2Parameters            `yaml:"parameters" json:"parameters"`
	//Responses   map[string]OpenAPI2PathResponse `yaml:"responses" json:"responses"`
}

type OpenAPI2Parameters struct {
	Name        string            `json:"name" yaml:"name"`
	In          string            `json:"in" yaml:"in"`
	Description string            `json:"description" yaml:"description"`
	Required    bool              `json:"required" json:"required"`
	Example     string            `json:"example" yaml:"example"`
	Schema      map[string]string `json:"schema" yaml:"schema"`
}

type OpenAPI2PathResponse struct {
	Description string                 `json:"description" yaml:"description"`
	Schema      map[string]interface{} `json:"schema" yaml:"schema"`
}
