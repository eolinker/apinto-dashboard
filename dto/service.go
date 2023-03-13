package dto

type ServiceListItem struct {
	Name        string `json:"name"`
	UUID        string `json:"uuid"`
	Scheme      string `json:"scheme"`
	ServiceType string `json:"service_type"`
	Config      string `json:"config"`
	UpdateTime  string `json:"update_time"`
	IsDelete    bool   `json:"is_delete"`
}

type ServiceConfigProxy []byte

func (c *ServiceConfigProxy) MarshalJSON() ([]byte, error) {
	return *c, nil
}

func (c *ServiceConfigProxy) UnmarshalJSON(bytes []byte) error {
	*c = bytes
	return nil
}

func (c *ServiceConfigProxy) String() string {
	return string(*c)
}

type ServiceInfoProxy struct {
	Name          string             `json:"name"`
	UUID          string             `json:"uuid"`
	Desc          string             `json:"desc"`
	Scheme        string             `json:"scheme"`
	DiscoveryName string             `json:"discovery_name"`
	Config        ServiceConfigProxy `json:"config"`
	Timeout       int                `json:"timeout"`
	Balance       string             `json:"balance"`
}

type ServiceInfo struct {
	Name        string
	UUID        string
	Desc        string
	Scheme      string
	DiscoveryID int
	DriverName  string
	FormatAddr  string //格式化地址概要，用于列表展示
	Config      string
	Timeout     int
	Balance     string
}

type ServiceInfoOutput struct {
	Service *ServiceInfoProxy `json:"service"`
	Render  Render            `json:"render"`
}

type Render []byte

func (r *Render) MarshalJSON() ([]byte, error) {
	return *r, nil
}
