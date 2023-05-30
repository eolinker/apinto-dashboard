CREATE TABLE IF NOT EXISTS `application`
(
    `id`          int(11) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `id_str`      varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '随机生成的16个长度字符串',
    `namespace`   int(11) NOT NULL COMMENT '工作空间',
    `name`        varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '应用名称',
    `desc`        varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '描述',
    `version`     varchar(36) NOT NULL COMMENT '应用版本',
    `operator`    int(11) NULL DEFAULT NULL COMMENT '更新人/操作人',
    `create_time` timestamp(0)                                                 NOT NULL DEFAULT CURRENT_TIMESTAMP(0) COMMENT '创建时间',
    `update_time` timestamp(0)                                                 NOT NULL DEFAULT CURRENT_TIMESTAMP(0) ON UPDATE CURRENT_TIMESTAMP (0) COMMENT '修改时间',
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE INDEX `namespace_name`(`namespace`, `name`) USING BTREE,
    UNIQUE INDEX `namespace_idstr`(`namespace`, `id_str`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '应用表' ROW_FORMAT = Dynamic;


CREATE TABLE IF NOT EXISTS `application_auth`
(
    `id`             int(11) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `uuid`           varchar(36) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci  NOT NULL COMMENT 'uuid',
    `title`          varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '鉴权名称',
    `namespace`      int(11) NOT NULL COMMENT '工作空间',
    `application`    int(11) NOT NULL COMMENT '应用ID',
    `is_transparent` tinyint(1) NOT NULL DEFAULT 0 COMMENT '是否透传至上游',
    `driver`         varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '鉴权类型,basic,apikey,aksk,jwt',
    `position`       varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci  NOT NULL COMMENT 'header,query',
    `token_name`     varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT 'tokenName',
    `expire_time`    int(11) NOT NULL DEFAULT 0 COMMENT '过期时间',
    `config`         text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '配置信息JSON',
    `operator`       int(11) NULL DEFAULT NULL COMMENT '更新人/操作人',
    `create_time`    timestamp(0)                                                  NOT NULL DEFAULT CURRENT_TIMESTAMP(0) COMMENT '创建时间',
    `update_time`    timestamp(0)                                                  NOT NULL DEFAULT CURRENT_TIMESTAMP(0) ON UPDATE CURRENT_TIMESTAMP (0) COMMENT '修改时间',
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE INDEX `uuid`(`uuid`) USING BTREE,
    INDEX            `application`(`application`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '应用鉴权表' ROW_FORMAT = Dynamic;
