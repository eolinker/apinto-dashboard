package driver

import (
	"bytes"
	"github.com/eolinker/apinto-dashboard/common"
	"github.com/eolinker/apinto-dashboard/entry"
	"github.com/eolinker/apinto-dashboard/enum"
	apimodel "github.com/eolinker/apinto-dashboard/modules/api/model"
	"github.com/go-basic/uuid"
	"net/url"
	"strings"
	"time"
)

type OpenAPI3 struct {
	format string
}

func (o *OpenAPI3) FormatAPI(data []byte, namespaceID, appID int, groupID, prefix, label string) ([]*apimodel.APIInfo, error) {
	openAPI3Config := new(OpenAPI3Config)
	reader := bytes.NewReader(data)
	if err := common.DecodeYAML(reader, openAPI3Config); err != nil {
		reader = bytes.NewReader(data)
		err = common.DecodeJSON(reader, openAPI3Config)
		if err != nil {
			return nil, err
		}
	}

	//格式化requestPrefix
	if prefix != "" {
		prefix = "/" + strings.Trim(prefix, "/")
	}

	apiList := make([]*apimodel.APIInfo, 0)
	t := time.Now()
	for path, pathMap := range openAPI3Config.Paths {
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
			apiList = append(apiList, &apimodel.APIInfo{
				API: &entry.API{
					NamespaceId: namespaceID,
					UUID:        uuid.New(),
					GroupUUID:   groupID,
					Name:        name,
					//对路径进行转义
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

func CreateOpenAPI3(format string) IAPISyncFormatDriver {
	return &OpenAPI3{format: format}
}

type OpenAPI3Config struct {
	Openapi    string                              `yaml:"openapi" json:"openapi"`
	Info       *OpenAPI3Info                       `yaml:"info" json:"info"`
	Tags       []*OpenAPI3Tags                     `yaml:"tags" json:"tags"`
	Paths      map[string]map[string]*OpenAPI3Path `yaml:"paths" json:"paths"` //第一个key为path 第二个key为method
	Components interface{}                         `yaml:"components" json:"components"`
}

type OpenAPI3Info struct {
	Title       string `yaml:"title" json:"title"`
	Description string `yaml:"description" json:"description"`
}

type OpenAPI3Tags struct {
	Name         string                   `yaml:"name" json:"name"`
	Description  string                   `yaml:"description" json:"description"`
	ExternalDocs OpenAPI3TagsExternalDocs `yaml:"externalDocs,omitempty"`
}

type OpenAPI3Path struct {
	Tags        []string `yaml:"tags" json:"tags"`
	Summary     string   `yaml:"summary" yaml:"summary"`
	Description string   `yaml:"description" yaml:"description"`
	OperationID string   `yaml:"operationId" json:"operationID"`
	//Parameters  []OpenAPI3Parameters            `yaml:"parameters" json:"parameters"`
	//Responses   map[string]OpenAPI3PathResponse `yaml:"responses" json:"responses"` //key为状态码
	//RequestBody OpenAPI3PathRequestBody         `yaml:"requestBody" json:"requestBody"`
}

type OpenAPI3Parameters struct {
	Name        string                   `json:"name" yaml:"name"`
	In          string                   `json:"in" yaml:"in"`
	Description string                   `json:"description" yaml:"description"`
	Required    bool                     `json:"required" json:"required"`
	Example     string                   `json:"example" yaml:"example"`
	Schema      OpenAPI3ParametersSchema `json:"schema" yaml:"schema"`
}
type OpenAPI3PathResponse struct {
	Description string                                 `json:"description" yaml:"description"`
	Content     map[string]OpenAPI3PathResponseContent `json:"content" yaml:"content"`
	Headers     map[string]OpenAPI3PathResponseHeader  `json:"headers" yaml:"headers"`
}

type OpenAPI3PathResponseHeader struct {
	Schema      map[string]interface{} `json:"schema" yaml:"schema"`
	Description string                 `json:"description" yaml:"description"`
}

type OpenAPI3PathRequestBody struct {
	Description string                                    `json:"description" yaml:"description"`
	Required    bool                                      `json:"required" yaml:"required"`
	Content     map[string]OpenAPI3PathRequestBodyContent `json:"content" yaml:"content"`
}
type OpenAPI3TagsExternalDocs struct {
	Description string `yaml:"description" json:"description"`
	URL         string `yaml:"url" json:"url"`
}

type OpenAPI3PathRequestBodyContent struct {
	Schema   map[string]interface{} `json:"schema" yaml:"schema"`
	Example  interface{}            `json:"example" yaml:"example"`
	Examples map[string]interface{} `json:"examples" yaml:"examples"`
	Encoding map[string]interface{} `json:"encoding" yaml:"encoding"`
}
type OpenAPI3PathResponseContent struct {
	Schema   map[string]interface{} `json:"schema" yaml:"schema"`
	Examples interface{}            `json:"examples" yaml:"examples"`
}

type OpenAPI3ParametersSchema struct {
	Type string `json:"type" yaml:"type"`
}
