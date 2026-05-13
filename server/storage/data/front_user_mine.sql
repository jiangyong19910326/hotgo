-- ============================================================
-- 前端用户矿场绑定表 (MySQL)
-- ============================================================

CREATE TABLE IF NOT EXISTS `hg_front_user_mine` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `user_id` bigint NOT NULL DEFAULT 0 COMMENT '前端用户ID',
  `mine_id` int NOT NULL DEFAULT 0 COMMENT '矿场ID',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user_mine` (`user_id`, `mine_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_mine_id` (`mine_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='前端用户矿场绑定';
