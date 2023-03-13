package entry

import "time"

type DiscoveryVersion struct {
	Id          int
	DiscoveryID int
	NamespaceID int
	DiscoveryVersionConfig
	Operator   int
	CreateTime time.Time
}

type DiscoveryVersionConfig struct {
	Config string `json:"config"`
}

func (d *DiscoveryVersion) SetVersionId(id int) {
	d.Id = id
}
