CREATE TABLE `notice_channel` (
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