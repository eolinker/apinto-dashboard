

ALTER TABLE `cluster_node` MODIFY COLUMN `admin_addr` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL COMMENT '管理地址' AFTER `cluster`, MODIFY COLUMN `service_addr` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '服务地址' AFTER `admin_addr`;