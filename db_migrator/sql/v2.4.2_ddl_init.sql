CREATE TABLE IF NOT EXISTS `dynamic_quote` (
 `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
 `namespace` int(11) NOT NULL COMMENT '命名空间ID',
 `source` varchar(255) NOT NULL COMMENT '源name',
 `target` varchar(255) NOT NULL COMMENT '依赖name',
 PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;