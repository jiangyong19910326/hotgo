-- PLC API 应用密钥表
-- 用于签名验签接口的 AppID/AppSecret 管理

CREATE TABLE `hg_plc_app` (
  `id`         int(10) unsigned NOT NULL AUTO_INCREMENT,
  `app_id`     varchar(32)  NOT NULL                  COMMENT 'AppID（唯一标识）',
  `app_secret` varchar(64)  NOT NULL                  COMMENT 'AppSecret（HMAC签名密钥）',
  `name`       varchar(64)  NOT NULL DEFAULT ''        COMMENT '应用名称',
  `remark`     varchar(255) NOT NULL DEFAULT ''        COMMENT '备注',
  `status`     tinyint(1)   NOT NULL DEFAULT 1         COMMENT '状态：1启用 2禁用',
  `created_at` datetime     DEFAULT NULL,
  `updated_at` datetime     DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_app_id` (`app_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='PLC API应用密钥';

-- 示例数据（secret 请在生产中替换）
INSERT INTO `hg_plc_app` (`app_id`, `app_secret`, `name`, `remark`, `status`)
VALUES ('plc_demo_app', 'change_me_32chars_secret_key_here', '演示应用', '默认测试用 AppID/Secret', 1);
