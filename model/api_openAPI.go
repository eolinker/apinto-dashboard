package model

type ApiOpenAPIGroups struct {
	Uuid     string              `json:"uuid"`
	Name     string              `json:"name"`
	Children []*ApiOpenAPIGroups `json:"children,omitempty"`
}

type ApiOpenAPIService struct {
	Name string `json:"name"`
	Desc string `json:"desc"`
}

type ServiceStaticDriverConf struct {
	UseVariable   bool                 `json:"use_variable"`
	AddrsVariable string               `json:"addrs_variable,omitempty"`
	StaticConf    []*ServiceStaticConf `json:"static_conf,omitempty"`
}

type ServiceStaticConf struct {
	Addr   string `json:"addr"`
	Weight int    `json:"weight"`
}
