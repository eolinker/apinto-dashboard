CREATE TABLE `monitor`
(
    `id`           int(11) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `uuid`         varchar(50) NOT NULL COMMENT '分区uuid，随机生成的长度32的字符串',
    `namespace`    int(11) NOT NULL COMMENT '工作空间',
    `name`         varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '分区名',
    `source_type`  varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '数据源类型',
    `config`       text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '数据源配置',
    `env`          varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '关联集群的环境',
    `cluster_ids`  varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '关联的集群id列表,全选为空',
    `operator`     int(11) NULL DEFAULT NULL COMMENT '更新人/操作人',
    `create_time`  timestamp(0)      NOT NULL DEFAULT CURRENT_TIMESTAMP(0) COMMENT '创建时间',
    `update_time`  timestamp(0)      NOT NULL DEFAULT CURRENT_TIMESTAMP(0) ON UPDATE CURRENT_TIMESTAMP (0) COMMENT '更新时间',
    PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '监控分区表' ROW_FORMAT = Dynamic;

ALTER TABLE `cluster` ADD COLUMN `uuid` varchar(255) NULL AFTER `addr`;