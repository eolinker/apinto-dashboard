ALTER TABLE `module_plugin_enable` ADD COLUMN `is_can_disable` tinyint(1) NOT NULL DEFAULT 0 COMMENT '是否可停用' AFTER `is_enable`;
ALTER TABLE `module_plugin_enable` ADD COLUMN `is_can_uninstall` tinyint(1) NOT NULL DEFAULT 0 COMMENT '是否可卸载' AFTER `is_enable`;
ALTER TABLE `module_plugin_enable` ADD COLUMN `is_show_server` tinyint(1) NOT NULL DEFAULT 0 COMMENT '启用时是否显示server' AFTER `is_enable`;
ALTER TABLE `module_plugin_enable` ADD COLUMN `is_plugin_visible` tinyint(1) NOT NULL DEFAULT 0 COMMENT '插件是否在导航可见' AFTER `is_enable`;
ALTER TABLE `module_plugin_enable` ADD COLUMN `frontend` varchar(255) COMMENT '前端路由' AFTER `is_enable`;

ALTER TABLE `module_plugin` DROP `front`;