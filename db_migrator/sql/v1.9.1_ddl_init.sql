ALTER TABLE `api` ADD `request_path_label` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT 'api请求路径Label' AFTER `request_path`;

ALTER TABLE `api` ADD `source_type` varchar(36) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '来源类型' AFTER `request_path_label`;

ALTER TABLE `api` ADD `source_id` int(11) NOT NULL COMMENT '来源id,用于关联外部应用' AFTER `source_type`;

ALTER TABLE `api` ADD `source_label` varchar(36) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '来源标签' AFTER `source_id`;