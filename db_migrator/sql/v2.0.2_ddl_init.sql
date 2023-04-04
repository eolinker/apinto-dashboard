CREATE TABLE `plugin_template`
(
    `id`          int(11) NOT NULL AUTO_INCREMENT COMMENT '插件模板主键ID',
    `uuid`        varchar(36)  NOT NULL COMMENT 'UUID',
    `namespace`   int(11) NOT NULL COMMENT '工作空间',
    `name`        varchar(255) NOT NULL COMMENT '名称',
    `desc`        varchar(255) NOT NULL COMMENT '描述',
    `operator`    int(11) DEFAULT NULL COMMENT '更新人/操作人',
    `create_time` timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE KEY `namespace_name` (`namespace`,`name`) USING BTREE,
    UNIQUE KEY `uuid` (`uuid`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC COMMENT='插件模板表';

CREATE TABLE `plugin`
(
    `id`          int(11) NOT NULL AUTO_INCREMENT COMMENT '插件主键ID',
    `namespace`   int(11) NOT NULL COMMENT '工作空间',
    `name`        varchar(255) NOT NULL COMMENT '名称',
    `extended`    varchar(255) NOT NULL COMMENT '扩展ID',
    `type`        tinyint(1) NOT NULL COMMENT '插件类型内置1 自建2',
    `desc`        varchar(255) NOT NULL COMMENT '描述',
    `rely`        int(11) DEFAULT 0 COMMENT '依赖的插件ID',
    `schema`      text         NOT NULL COMMENT 'jsonSchema',
    `sort`        int(11) NOT NULL COMMENT '排序',
    `operator`    int(11) DEFAULT NULL COMMENT '更新人/操作人',
    `create_time` timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE KEY `namespace_name` (`namespace`,`name`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC COMMENT='插件表';

CREATE TABLE `cluster_plugin`
(
    `id`          int(11) NOT NULL AUTO_INCREMENT COMMENT '集群插件主键ID',
    `namespace`   int(11) NOT NULL COMMENT '工作空间',
    `cluster`     int(11) NOT NULL COMMENT '集群ID',
    `plugin_name` varchar(255) NOT NULL COMMENT '插件名',
    `status`      tinyint(2) NOT NULL COMMENT '1禁用 2启用 3全局启用',
    `config`      text         NOT NULL COMMENT '插件配置信息',
    `operator`    int(11) DEFAULT NULL COMMENT '更新人/操作人',
    `create_time` timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE KEY `cluster_plugin` (`cluster`,`plugin_name`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC COMMENT='插件表';