CREATE TABLE IF NOT EXISTS `user` (
    `id` int(10) NOT NULL AUTO_INCREMENT COMMENT '主键id',
    `sex` int(11) NOT NULL DEFAULT '0' COMMENT '性别，0未知，1男2女',
    `username` varchar(36) NOT NULL COMMENT '用户名',
    `password` varchar(32) NOT NULL COMMENT '密码',
    `notice` varchar(36) DEFAULT NULL COMMENT '通知key',
    `nickname` varchar(255) DEFAULT NULL COMMENT '昵称',
    `email` varchar(255) DEFAULT NULL COMMENT '邮箱',
    `phone` varchar(20) DEFAULT NULL COMMENT '手机号',
    `avatar` varchar(255) DEFAULT NULL COMMENT '头像',
    `login_time` timestamp NULL DEFAULT NULL COMMENT '最后登录时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `user_pk2` (`username`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8 COMMENT='用户信息表';