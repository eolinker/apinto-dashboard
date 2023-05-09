package open_app_entry

import "time"

type ExternalApplication struct {
	Id         int       `json:"id,omitempty"`
	UUID       string    `json:"uuid,omitempty"`
	Namespace  int       `json:"namespace,omitempty"`
	Name       string    `json:"name,omitempty"`
	Token      string    `json:"token,omitempty"`
	Desc       string    `json:"desc,omitempty"`
	Tags       string    `json:"tags,omitempty"`
	IsDisable  bool      `json:"is_disable,omitempty"`
	IsDelete   bool      `json:"is_delete,omitempty"`
	Operator   int       `json:"operator,omitempty"`
	CreateTime time.Time `json:"create_time,omitempty"`
	UpdateTime time.Time `json:"update_time,omitempty"`
}

func (e *ExternalApplication) TableName() string {
	return "external_app"
}

func (e *ExternalApplication) IdValue() int {
	return e.Id
}
