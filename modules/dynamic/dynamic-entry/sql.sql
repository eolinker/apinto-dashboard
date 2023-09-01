CREATE TABLE `dynamic_quote` (
                                 `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
                                 `namespace` int(11) NOT NULL COMMENT '命名空间ID',
                                 `source` varchar(255) NOT NULL COMMENT '源name',
                                 `target` varchar(255) NOT NULL COMMENT '依赖name',
                                 PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8

CREATE TABLE `dynamic_module` (
                                  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
                                  `namespace` bigint(20) DEFAULT NULL,
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
) ENGINE=InnoDB AUTO_INCREMENT=86 DEFAULT CHARSET=utf8

