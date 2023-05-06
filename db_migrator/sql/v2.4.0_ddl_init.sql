ALTER TABLE `api` ADD COLUMN `scheme` varchar(36) NOT NULL DEFAULT http AFTER `name`;
ALTER TABLE `api` ADD COLUMN `version` varchar(36) NOT NULL DEFAULT "20230506171816" AFTER `name`;
ALTER TABLE `api` ADD COLUMN `is_disable` tinyint(1) NOT NULL DEFAULT 0 AFTER `name`;