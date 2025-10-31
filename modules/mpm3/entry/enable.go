package entry

import "time"

type PluginEnable struct {
	Id         int       `gorm:"type:int(11);primary_key;column:id;autoIncrement:false;comment:模块插件的主键id" json:"id"`
	IsEnable   int       `gorm:"column:is_enable;not null;type:tinyint(1);comment:是否启用 1未启用 2启用" json:"is_enable"`
	Config     []byte    `gorm:"column:config;type:text;comment:启用配置;serializer:json" json:"config"`
	Operator   int       `gorm:"column:operator;type:int(11);comment:更新人/操作人" json:"operator"`
	UpdateTime time.Time `gorm:"column:update_time;type:timestamp;autoCreateTime;autoUpdateTime;comment:修改时间" json:"update_time"`
}

func (*PluginEnable) TableName() string {
	return "pm3_enable"
}

func (p *PluginEnable) IdValue() int {
	return p.Id
}
