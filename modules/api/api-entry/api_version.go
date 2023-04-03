package api_entry

import "time"

type APIVersion struct {
	Id          int
	ApiID       int
	NamespaceID int
	APIVersionConfig
	Operator   int
	CreateTime time.Time
}

type APIVersionConfig struct {
	Driver           string         `json:"driver"`
	RequestPath      string         `json:"request_path"`
	RequestPathLabel string         `json:"request_path_label"`
	ServiceID        int            `json:"service_id"`
	ServiceName      string         `json:"service_name"`
	TemplateID       int            `json:"template_id"`
	TemplateUUID     string         `json:"template_uuid"`
	Method           []string       `json:"method"`
	ProxyPath        string         `json:"proxy_path"`
	Timeout          int            `json:"timeout"`
	Retry            int            `json:"retry"`
	EnableWebsocket  bool           `json:"enable_websocket"`
	Match            []*MatchConf   `json:"match"`
	Header           []*ProxyHeader `json:"proxy_header"`
}

type MatchConf struct {
	Position  string `json:"position"`
	MatchType string `json:"match_type"`
	Key       string `json:"key"`
	Pattern   string `json:"pattern"`
}

type ProxyHeader struct {
	OptType string `json:"opt_type"`
	Key     string `json:"key"`
	Value   string `json:"value"`
}

func (d *APIVersion) SetVersionId(id int) {
	d.Id = id
}
