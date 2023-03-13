package entry

import "time"

type ServiceVersion struct {
	Id          int
	ServiceId   int
	NamespaceID int
	ServiceVersionConfig
	Operator   int
	CreateTime time.Time
}

type ServiceVersionConfig struct {
	DiscoveryId int    `json:"discovery_id,omitempty"`
	DriverName  string `json:"driver_name,omitempty"`
	Scheme      string `json:"scheme,omitempty"`
	Balance     string `json:"balance,omitempty"`
	Timeout     int    `json:"timeout,omitempty"`
	FormatAddr  string `json:"format_addr,omitempty"`
	Config      string `json:"config,omitempty"`
}

func (s *ServiceVersion) SetVersionId(id int) {
	s.Id = id
}
