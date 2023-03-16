package discover_dto

import "github.com/eolinker/apinto-dashboard/modules/upstream/upstream-dto"

type DiscoveryListItem struct {
	Name       string `json:"name"`
	UUID       string `json:"uuid"`
	Driver     string `json:"driver"`
	Desc       string `json:"desc"`
	UpdateTime string `json:"update_time"`
	IsDelete   bool   `json:"is_delete"`
}

type ConfigProxy []byte

func (d *ConfigProxy) String() string {
	return string(*d)
}

func (d *ConfigProxy) UnmarshalJSON(bytes []byte) error {
	*d = bytes
	return nil
}

func (d *ConfigProxy) MarshalJSON() ([]byte, error) {
	return *d, nil
}

type DiscoveryInfoProxy struct {
	Name   string      `json:"name"`
	UUID   string      `json:"uuid"`
	Driver string      `json:"driver"`
	Desc   string      `json:"desc"`
	Config ConfigProxy `json:"config"`
}

type DiscoveryEnum struct {
	Name   string              `json:"name"`
	Driver string              `json:"driver"`
	Render upstream_dto.Render `json:"render"`
}

type DriversItem struct {
	Name   string              `json:"name"`
	Render upstream_dto.Render `json:"render"`
}

type DiscoveryInfoOutput struct {
	Discovery *DiscoveryInfoProxy `json:"discovery"`
	Render    upstream_dto.Render `json:"render"`
}
