package v1

type ServiceConfig struct {
	Timeout      int                `json:"timeout"`           //请求超时时间 ms
	Retry        int                `json:"retry"`             //失败重试次数
	Scheme       string             `json:"scheme"`            //HTTP,HTTPS
	Discovery    string             `json:"discovery"`         //服务发现名称 匿名上游为空      注册中心名@discovery
	Nodes        []string           `json:"nodes,omitempty"`   //静态配置  discovery为空时填
	Balance      string             `json:"balance"`           //round-robin
	Plugins      map[string]*Plugin `json:"plugins,omitempty"` //插件列表  key插件名称
	Name         string             `json:"name"`              //名称
	Driver       string             `json:"driver"`            //http
	Description  string             `json:"description"`       //描述
	Service      string             `json:"service,omitempty"` //服务名 or 配置
	PassHost     string             `json:"pass_host"`         //转发host方式，pass,node,rewrite,rewrite则重写为upstream_host
	UpstreamHost string             `json:"upstream_host,omitempty"`
}
