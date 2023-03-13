package model

type Info struct {
	Title       string `yaml:"title" json:"title"`
	Description string `yaml:"description" json:"description"`
}

type Tags struct {
	Name         string           `yaml:"name" json:"name"`
	Description  string           `yaml:"description" json:"description"`
	ExternalDocs TagsExternalDocs `yaml:"externalDocs,omitempty"`
}

type TagsExternalDocs struct {
	Description string `yaml:"description" json:"description"`
	URL         string `yaml:"url" json:"url"`
}

type Path struct {
	Tags        []string                `yaml:"tags" json:"tags"`
	Summary     string                  `yaml:"summary" yaml:"summary"`
	Description string                  `yaml:"description" yaml:"description"`
	OperationID string                  `yaml:"operationId" json:"operationID"`
	Parameters  []Parameters            `yaml:"parameters" json:"parameters"`
	Responses   map[string]PathResponse `yaml:"responses" json:"responses"` //key为状态码
	RequestBody PathRequestBody         `yaml:"requestBody" json:"requestBody"`
}

type PathRequestBody struct {
	Description string                            `json:"description" yaml:"description"`
	Required    bool                              `json:"required" yaml:"required"`
	Content     map[string]PathRequestBodyContent `json:"content" yaml:"content"`
}
type PathRequestBodyContent struct {
	Schema   map[string]interface{} `json:"schema" yaml:"schema"`
	Example  interface{}            `json:"example" yaml:"example"`
	Examples map[string]interface{} `json:"examples" yaml:"examples"`
	Encoding map[string]interface{} `json:"encoding" yaml:"encoding"`
}

type PathResponse struct {
	Description string                         `json:"description" yaml:"description"`
	Content     map[string]PathResponseContent `json:"content" yaml:"content"`
	Headers     map[string]PathResponseHeader  `json:"headers" yaml:"headers"`
}

type PathResponseContent struct {
	Schema   map[string]interface{} `json:"schema" yaml:"schema"`
	Examples interface{}            `json:"examples" yaml:"examples"`
}
type PathResponseHeader struct {
	Schema      map[string]interface{} `json:"schema" yaml:"schema"`
	Description string                 `json:"description" yaml:"description"`
}

type Parameters struct {
	Name        string           `json:"name" yaml:"name"`
	In          string           `json:"in" yaml:"in"`
	Description string           `json:"description" yaml:"description"`
	Required    bool             `json:"required" json:"required"`
	Example     string           `json:"example" yaml:"example"`
	Schema      ParametersSchema `json:"schema" yaml:"schema"`
}
type ParametersSchema struct {
	Type string `json:"type" yaml:"type"`
}

type SwaggerConfig struct {
	Openapi    string                      `yaml:"openapi" json:"openapi"`
	Info       *Info                       `yaml:"info" json:"info"`
	Tags       []*Tags                     `yaml:"tags" json:"tags"`
	Paths      map[string]map[string]*Path `yaml:"paths" json:"paths"` //第一个key为path 第二个key为method
	Components interface{}                 `yaml:"components" json:"components"`
}
