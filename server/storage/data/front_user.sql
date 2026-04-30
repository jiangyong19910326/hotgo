-- ============================================================
-- HotGo 前端登录用户表迁移脚本 (MySQL)
-- 适用版本: HotGo v2.x
-- ============================================================

CREATE TABLE IF NOT EXISTS `hg_front_user` (
  `id`              bigint(20)    unsigned NOT NULL AUTO_INCREMENT     COMMENT '主键',
  `username`        varchar(64)            NOT NULL DEFAULT ''         COMMENT '用户名',
  `nickname`        varchar(64)            NOT NULL DEFAULT ''         COMMENT '昵称',
  `password_hash`   varchar(64)            NOT NULL DEFAULT ''         COMMENT '密码哈希',
  `salt`            varchar(16)            NOT NULL DEFAULT ''         COMMENT '密码盐',
  `avatar`          varchar(255)           NOT NULL DEFAULT ''         COMMENT '头像',
  `mobile`          varchar(20)            NOT NULL DEFAULT ''         COMMENT '手机号',
  `email`           varchar(100)           NOT NULL DEFAULT ''         COMMENT '邮箱',
  `remark`          varchar(255)           NOT NULL DEFAULT ''         COMMENT '备注',
  `status`          tinyint(1)             NOT NULL DEFAULT 1          COMMENT '状态：1启用 2禁用',
  `last_login_at`   datetime                        DEFAULT NULL       COMMENT '最近登录时间',
  `last_login_ip`   varchar(45)            NOT NULL DEFAULT ''         COMMENT '最近登录IP',
  `created_by`      bigint(20)             NOT NULL DEFAULT 0          COMMENT '创建者',
  `updated_by`      bigint(20)             NOT NULL DEFAULT 0          COMMENT '更新者',
  `created_at`      datetime                        DEFAULT NULL       COMMENT '创建时间',
  `updated_at`      datetime                        DEFAULT NULL       COMMENT '修改时间',
  `deleted_at`      datetime                        DEFAULT NULL       COMMENT '删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_username` (`username`),
  KEY `idx_status` (`status`),
  KEY `idx_mobile` (`mobile`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='前端登录用户';
