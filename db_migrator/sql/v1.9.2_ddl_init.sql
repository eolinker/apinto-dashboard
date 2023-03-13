CREATE TABLE `external_app`
(
    `id`           int(11) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `uuid`         varchar(36) NOT NULL COMMENT '应用id，随机生成的长度32的字符串',
    `namespace`    int(11) NOT NULL COMMENT '工作空间',
    `name`         varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '应用名',
    `token`        varchar(36) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '鉴权token',
    `desc`         text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '备注',
    `tags`         text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '关联标签',
    `is_disable`   tinyint(1) NULL DEFAULT 0 COMMENT '禁用状态，0启用，1禁用',
    `is_delete`    tinyint(1) NULL DEFAULT 0 COMMENT '是否删除',
    `operator`     int(11) NULL DEFAULT NULL COMMENT '更新人/操作人',
    `create_time`  timestamp(0)      NOT NULL DEFAULT CURRENT_TIMESTAMP(0) COMMENT '创建时间',
    `update_time`  timestamp(0)      NOT NULL DEFAULT CURRENT_TIMESTAMP(0) ON UPDATE CURRENT_TIMESTAMP (0) COMMENT '更新时间',
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE KEY `unique_uuid` (`uuid`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '外部应用表' ROW_FORMAT = Dynamic;