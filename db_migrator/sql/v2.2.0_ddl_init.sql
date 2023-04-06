CREATE TABLE `middleware_group` (
    `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '分组ID',
    `uuid` varchar(32) NOT NULL COMMENT '中间件分组唯一id',
    `namespace_id` int(11) NOT NULL COMMENT '分组id',
    `prefix` varchar(255) NOT NULL COMMENT '前缀',
    `operator` int(11) NOT NULL COMMENT '更新者/操作者id',
    `cname` varchar(255) DEFAULT NULL COMMENT '分组名称',
    `create_time` datetime NOT NULL COMMENT '创建时间',
    `update_time` datetime NOT NULL COMMENT '更新时间',
    `middlewares` text NOT NULL COMMENT '中间件列表，字符串列表',
    PRIMARY KEY (`id`),
    UNIQUE KEY `unique_prefix` (`prefix`) USING BTREE COMMENT '唯一前缀索引',
    UNIQUE KEY `unique_id` (`uuid`) USING BTREE COMMENT '分组唯一id'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;