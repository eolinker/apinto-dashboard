CREATE TABLE `module_plugin` (
    `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '模块插件主键ID',
    `uuid` varchar(36) NOT NULL COMMENT 'UUID',
    `name` varchar(255) NOT NULL COMMENT '插件名',
    `version` varchar(36) NOT NULL COMMENT 'version',
    `group` int(11) NOT NULL COMMENT '分组id',
    `cname` varchar(255) NOT NULL COMMENT '昵称',
    `resume` varchar(255) NOT NULL COMMENT '简介',
    `icon` varchar(255) NOT NULL COMMENT '图标的文件名，相对路径',
    `type` tinyint(1) NOT NULL COMMENT '是否为内置插件 0为内置-内核模块 1为内置-非内核 2为非内置',
    `front` varchar(255) COMMENT '前端模块路由 为空表示不需要在前端显示',
    `driver` varchar(255) NOT NULL COMMENT '插件类型',
    `details` mediumtext NOT NULL COMMENT '插件详情',
    `operator` int(11) DEFAULT NULL COMMENT '更新人/操作人',
    `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE KEY `unique_uuid` (`uuid`) USING BTREE,
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC COMMENT='模块插件表';

CREATE TABLE `module_plugin_enable` (
    `id` int(11) NOT NULL COMMENT '模块插件的主键id',
    `name` varchar(255) NOT NULL COMMENT '模块名，可以改，默认为模块插件的name',
    `navigation` int(11) NOT NULL COMMENT '导航id',
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
    PRIMARY KEY (`id`) USING BTREE,
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC COMMENT='模块插件安装包表';

CREATE TABLE `navigation` (
    `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `uuid` varchar(36) NOT NULL COMMENT 'uuid，唯一id',
    `title` varchar(255) NOT NULL COMMENT '导航名称',
    `icon` text COMMENT 'Icon信息，dataurl格式',
    `icon_type` varchar(255) NOT NULL COMMENT '图标类型，可选值:url、css',
    `sort` tinyint(4) unsigned NOT NULL COMMENT '排序，数字越小优先级越高',
    `module` text COMMENT '模块ID列表',
    PRIMARY KEY (`id`),
    UNIQUE KEY `unique_uuid` (`uuid`) USING BTREE COMMENT '唯一uuid'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC COMMENT='导航表';
