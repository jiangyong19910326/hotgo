-- ============================================================
-- HotGo 短链接模块数据表迁移脚本 (MySQL)
-- 适用版本: HotGo v2.x
-- 使用方式: 可直接在已有数据库中执行，使用 IF NOT EXISTS 安全幂等
-- ============================================================

-- --------------------------------------------------------
-- 表结构 `hg_short_link`
-- --------------------------------------------------------

CREATE TABLE IF NOT EXISTS `hg_short_link` (
  `id`               bigint(20)    NOT NULL AUTO_INCREMENT                   COMMENT '主键',
  `code`             varchar(32)   NOT NULL                                  COMMENT '短码',
  `original_url`     text          NOT NULL                                  COMMENT '原始链接',
  `title`            varchar(255)  DEFAULT ''                                COMMENT '标题',
  `total_clicks`     bigint(20)    DEFAULT 0                                 COMMENT '总点击量',
  `today_clicks`     int(11)       DEFAULT 0                                 COMMENT '今日点击（每日凌晨重置）',
  `yesterday_clicks` int(11)       DEFAULT 0                                 COMMENT '昨日点击',
  `weekly_clicks`    int(11)       DEFAULT 0                                 COMMENT '本周点击',
  `monthly_clicks`   int(11)       DEFAULT 0                                 COMMENT '本月点击',
  `expire_at`        datetime      DEFAULT NULL                              COMMENT '过期时间，NULL 表示永不过期',
  `status`           tinyint(1)    DEFAULT 1                                 COMMENT '状态：1正常 2禁用',
  `created_by`       bigint(20)    DEFAULT 0                                 COMMENT '创建者 member_id',
  `updated_by`       bigint(20)    DEFAULT 0                                 COMMENT '更新者 member_id',
  `created_at`       datetime      DEFAULT NULL                              COMMENT '创建时间',
  `updated_at`       datetime      DEFAULT NULL                              COMMENT '修改时间',
  `deleted_at`       datetime      DEFAULT NULL                              COMMENT '删除时间（软删除）',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_code` (`code`),
  KEY `idx_status`     (`status`),
  KEY `idx_created_by` (`created_by`),
  KEY `idx_expire_at`  (`expire_at`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='短链接';

-- --------------------------------------------------------
-- 表结构 `hg_short_link_log`
-- 记录每次点击明细，用于趋势/来源/地区统计
-- --------------------------------------------------------

CREATE TABLE IF NOT EXISTS `hg_short_link_log` (
  `id`           bigint(20)   NOT NULL AUTO_INCREMENT  COMMENT '主键',
  `link_id`      bigint(20)   NOT NULL                 COMMENT '短链ID',
  `code`         varchar(32)  NOT NULL                 COMMENT '短码（冗余，提高查询效率）',
  `ip`           varchar(64)  DEFAULT ''               COMMENT '访客 IP',
  `region`       varchar(128) DEFAULT ''               COMMENT '访问地区（省市）',
  `country`      varchar(64)  DEFAULT ''               COMMENT '访问国家',
  `referer`      varchar(1024) DEFAULT ''              COMMENT '来源完整 URL',
  `referer_host` varchar(256) DEFAULT ''               COMMENT '来源域名',
  `user_agent`   varchar(512) DEFAULT ''               COMMENT 'User-Agent',
  `device`       varchar(32)  DEFAULT ''               COMMENT '设备类型：desktop / mobile / tablet',
  `browser`      varchar(64)  DEFAULT ''               COMMENT '浏览器',
  `os`           varchar(64)  DEFAULT ''               COMMENT '操作系统',
  `created_at`   datetime     DEFAULT NULL             COMMENT '点击时间',
  PRIMARY KEY (`id`),
  KEY `idx_link_id`      (`link_id`),
  KEY `idx_code`         (`code`),
  KEY `idx_referer_host` (`referer_host`),
  KEY `idx_region`       (`region`),
  KEY `idx_created_at`   (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='短链接_点击日志';
