package entry

import "time"

type FlywaySchemaHistory struct {
	ID         int       `gorm:"column:id;primary_key;comment:主键ID"`
	Version    string    `gorm:"column:version;uniqueIndex:version;comment:版本号"`
	VersionNum int64     `gorm:"column:version_num;uniqueIndex:version_num;comment:版本号(转换数字后)"`
	Script     string    `gorm:"column:script;comment:完整脚本名"`
	Type       string    `gorm:"column:type;comment:ddl,dml"`
	Success    bool      `gorm:"column:success;comment:是否成功执行"`
	Md5        string    `gorm:"column:md5;comment:加密串,用于对比原数据"`
	CreateTime time.Time `gorm:"column:create_time;comment:创建时间"`
}

func (FlywaySchemaHistory) TableName() string {
	return "flyway_schema_history"
}

func (FlywaySchemaHistory) CreateTableSql() string {
	return "CREATE TABLE `flyway_schema_history` (\n  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键ID',\n  `version` varchar(20) NOT NULL COMMENT '版本号',\n  `type` varchar(20) NOT NULL COMMENT 'ddl,dml',\n  `version_num` int(11) NOT NULL COMMENT '版本号(转换数字后)',\n  `script` varchar(100) NOT NULL COMMENT '完整脚本名',\n  `success` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否成功执行',\n  `md5` varchar(255) NOT NULL COMMENT '加密串,用于对比原数据',\n  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',\n  PRIMARY KEY (`id`),\n  UNIQUE KEY `version_num_type` (`version_num`,`type`)\n) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;"
}
