create table if not exists `notice_channel` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `namespace` int(11) NOT NULL COMMENT '工作空间',
  `name` varchar(36) NOT NULL COMMENT 'email或uuid',
  `title` varchar(255) NOT NULL COMMENT '名称',
  `type` int(11) NOT NULL COMMENT '1.webhook 2.email',
  `operator` int(11) DEFAULT NULL COMMENT '更新人/操作人',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `namespace_name` (`namespace`,`name`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COMMENT='通知渠道表';

create table if not exists `dynamic_module` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `namespace` int(11) NOT NULL COMMENT '命名空间',
  `name` varchar(255) NOT NULL COMMENT '实例名',
  `title` varchar(255) DEFAULT NULL COMMENT '标题',
  `driver` varchar(255) NOT NULL COMMENT '驱动名称',
  `description` varchar(255) DEFAULT NULL COMMENT '描述',
  `version` varchar(32) NOT NULL COMMENT '当前版本',
  `config` text NOT NULL COMMENT '配置',
  `profession` varchar(255) NOT NULL COMMENT '插件指定profession',
  `create_time` datetime NOT NULL COMMENT '创建时间',
  `update_time` datetime NOT NULL COMMENT '更新时间',
  `updater` int(11) NOT NULL COMMENT '更新者ID',
  `skill` varchar(255) DEFAULT NULL COMMENT '模块提供能力',
  PRIMARY KEY (`id`),
  UNIQUE KEY `unique` (`namespace`,`name`,`profession`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;