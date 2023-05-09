package driver

type commonParams struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

var commonDiscoveryEnumRender = `{
	"type": "object",
	"required": true,
	"properties": {
		"service_name": {
			"type": "text",
			"title": "服务名",
			"required": true,
			"x-component": "Input",
			"x-component-props": {
				"placeholder": "请输入"
			},
			"x-index": 1
		}
	}
}`

type ApintoDiscoveryConfig struct {
	Address []string          `json:"address"`
	Params  map[string]string `json:"params"`
}
