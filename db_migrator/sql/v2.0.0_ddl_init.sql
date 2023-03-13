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

CREATE TABLE `warn_strategy` (
    `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `namespace` int(11) NOT NULL COMMENT '工作空间',
    `partition_id` int(11) NOT NULL COMMENT '分区ID',
    `uuid` varchar(36) NOT NULL COMMENT 'uuid',
    `title` varchar(255)NOT NULL COMMENT '告警策略名称',
    `desc` varchar(255) default NULL COMMENT '描述',
    `is_enable`  tinyint(1) DEFAULT '1' COMMENT '是否启用',
    `dimension` varchar(20) NOT NULL COMMENT '告警维度 api/service/cluster/partition',
    `quota` varchar(50)NOT NULL COMMENT '告警指标',
    `every` int(11) NOT NULL comment '统计时间 单位分钟',
    `config` text NOT NULL COMMENT '配置信息',
    `operator` int(11) DEFAULT NULL COMMENT '更新人/操作人',
    `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uuid` (`uuid`),
    UNIQUE KEY `namespace_partition_title` (`namespace`,`partition_id`,`title`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COMMENT='告警策略表';

CREATE TABLE `warn_history` (
   `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
   `namespace` int(11) NOT NULL COMMENT '工作空间',
   `partition_id` int(11) NOT NULL COMMENT '分区ID',
   `strategy_title` varchar(255) NOT NULL COMMENT '策略名称',
   `status` tinyint(1) DEFAULT '0' COMMENT '发送状态 0未发送 1已发送 2发送失败 3部分成功',
   `err_msg` varchar(255)  NULL COMMENT '发送失败原因',
   `dimension` varchar(20) NOT NULL COMMENT '告警维度 api/service/cluster/partition',
   `quota` varchar(50)NOT NULL COMMENT '告警指标',
   `target` text NOT NULL COMMENT '告警目标',
   `content` text NOT NULL COMMENT '告警内容',
   `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`),
   KEY partition_title(`partition_id`,`strategy_title`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COMMENT='告警历史表';

ALTER TABLE `user_info` ADD COLUMN `notice_user_id` varchar(36) NULL COMMENT '通知用户ID' AFTER `nick_name`;