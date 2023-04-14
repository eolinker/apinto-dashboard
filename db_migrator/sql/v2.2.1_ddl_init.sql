CREATE TABLE `module_plugin` (
    `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '模块插件主键ID',
    `uuid` varchar(36) NOT NULL COMMENT 'UUID',
    `name` varchar(255) NOT NULL COMMENT '插件名',
    `version` varchar(36) NOT NULL COMMENT 'version',
    `group` varchar(255) NOT NULL COMMENT '插件分组id',
    `navigation` varchar(255) COMMENT '导航id',
    `cname` varchar(255) NOT NULL COMMENT '昵称',
    `resume` varchar(255) NOT NULL COMMENT '简介',
    `icon` varchar(255) NOT NULL COMMENT '图标的文件名，相对路径',
    `type` tinyint(1) NOT NULL COMMENT '插件类型 0为框架模块 1为核心模块 2为内置模块 3为非内置',
    `front` varchar(255) COMMENT '前端模块路由 为空表示不需要在前端显示',
    `driver` varchar(255) NOT NULL COMMENT '插件类型',
    `details` mediumtext NOT NULL COMMENT '插件详情',
    `operator` int(11) DEFAULT NULL COMMENT '更新人/操作人',
    `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE KEY `unique_uuid` (`uuid`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC COMMENT='模块插件表';

CREATE TABLE `module_plugin_enable` (
    `id` int(11) NOT NULL COMMENT '模块插件的主键id',
    `name` varchar(255) NOT NULL COMMENT '模块名，可以改，默认为模块插件的name',
    `navigation` varchar(255) NOT NULL COMMENT '导航id',
    `is_enable` tinyint(1) NOT NULL COMMENT '是否启用 1未启用 2启用',
    `config` text NOT NULL COMMENT '启用配置',
    `operator` int(11) DEFAULT NULL COMMENT '更新人/操作人',
    `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE KEY `unique_name` (`name`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC COMMENT='模块插件启用表';

CREATE TABLE `module_plugin_package` (
    `id` int(11) NOT NULL COMMENT '模块插件表的主键ID',
    `package` mediumblob NOT NULL COMMENT '安装包',
    PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC COMMENT='模块插件安装包表';

