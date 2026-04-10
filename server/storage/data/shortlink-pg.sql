-- ============================================================
-- HotGo 短链接模块数据表迁移脚本 (PostgreSQL)
-- 适用版本: HotGo v2.x
-- 使用方式: 可直接在已有数据库中执行，使用 IF NOT EXISTS 安全幂等
-- ============================================================

-- --------------------------------------------------------
-- 表结构 hg_short_link
-- --------------------------------------------------------

CREATE TABLE IF NOT EXISTS hg_short_link (
    id               BIGSERIAL       PRIMARY KEY,
    code             VARCHAR(32)     NOT NULL,
    original_url     TEXT            NOT NULL,
    title            VARCHAR(255)    NOT NULL DEFAULT '',
    total_clicks     BIGINT          NOT NULL DEFAULT 0,
    today_clicks     INT             NOT NULL DEFAULT 0,
    yesterday_clicks INT             NOT NULL DEFAULT 0,
    weekly_clicks    INT             NOT NULL DEFAULT 0,
    monthly_clicks   INT             NOT NULL DEFAULT 0,
    expire_at        TIMESTAMP       DEFAULT NULL,
    status           SMALLINT        NOT NULL DEFAULT 1,
    created_by       BIGINT          NOT NULL DEFAULT 0,
    updated_by       BIGINT          NOT NULL DEFAULT 0,
    created_at       TIMESTAMP       DEFAULT NULL,
    updated_at       TIMESTAMP       DEFAULT NULL,
    deleted_at       TIMESTAMP       DEFAULT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS uk_short_link_code      ON hg_short_link (code);
CREATE        INDEX IF NOT EXISTS idx_short_link_status   ON hg_short_link (status);
CREATE        INDEX IF NOT EXISTS idx_short_link_created  ON hg_short_link (created_by);
CREATE        INDEX IF NOT EXISTS idx_short_link_expire   ON hg_short_link (expire_at);
CREATE        INDEX IF NOT EXISTS idx_short_link_cat      ON hg_short_link (created_at);

COMMENT ON TABLE  hg_short_link                IS '短链接';
COMMENT ON COLUMN hg_short_link.id             IS '主键';
COMMENT ON COLUMN hg_short_link.code           IS '短码';
COMMENT ON COLUMN hg_short_link.original_url   IS '原始链接';
COMMENT ON COLUMN hg_short_link.title          IS '标题';
COMMENT ON COLUMN hg_short_link.total_clicks   IS '总点击量';
COMMENT ON COLUMN hg_short_link.today_clicks   IS '今日点击（每日凌晨重置）';
COMMENT ON COLUMN hg_short_link.yesterday_clicks IS '昨日点击';
COMMENT ON COLUMN hg_short_link.weekly_clicks  IS '本周点击';
COMMENT ON COLUMN hg_short_link.monthly_clicks IS '本月点击';
COMMENT ON COLUMN hg_short_link.expire_at      IS '过期时间，NULL 表示永不过期';
COMMENT ON COLUMN hg_short_link.status         IS '状态：1正常 2禁用';
COMMENT ON COLUMN hg_short_link.created_by     IS '创建者 member_id';
COMMENT ON COLUMN hg_short_link.updated_by     IS '更新者 member_id';
COMMENT ON COLUMN hg_short_link.created_at     IS '创建时间';
COMMENT ON COLUMN hg_short_link.updated_at     IS '修改时间';
COMMENT ON COLUMN hg_short_link.deleted_at     IS '删除时间（软删除）';

-- --------------------------------------------------------
-- 表结构 hg_short_link_log
-- 记录每次点击明细，用于趋势/来源/地区统计
-- --------------------------------------------------------

CREATE TABLE IF NOT EXISTS hg_short_link_log (
    id           BIGSERIAL       PRIMARY KEY,
    link_id      BIGINT          NOT NULL,
    code         VARCHAR(32)     NOT NULL,
    ip           VARCHAR(64)     NOT NULL DEFAULT '',
    region       VARCHAR(128)    NOT NULL DEFAULT '',
    country      VARCHAR(64)     NOT NULL DEFAULT '',
    referer      VARCHAR(1024)   NOT NULL DEFAULT '',
    referer_host VARCHAR(256)    NOT NULL DEFAULT '',
    user_agent   VARCHAR(512)    NOT NULL DEFAULT '',
    device       VARCHAR(32)     NOT NULL DEFAULT '',
    browser      VARCHAR(64)     NOT NULL DEFAULT '',
    os           VARCHAR(64)     NOT NULL DEFAULT '',
    created_at   TIMESTAMP       DEFAULT NULL
);

CREATE INDEX IF NOT EXISTS idx_short_link_log_link_id      ON hg_short_link_log (link_id);
CREATE INDEX IF NOT EXISTS idx_short_link_log_code         ON hg_short_link_log (code);
CREATE INDEX IF NOT EXISTS idx_short_link_log_referer_host ON hg_short_link_log (referer_host);
CREATE INDEX IF NOT EXISTS idx_short_link_log_region       ON hg_short_link_log (region);
CREATE INDEX IF NOT EXISTS idx_short_link_log_created_at   ON hg_short_link_log (created_at);

COMMENT ON TABLE  hg_short_link_log              IS '短链接_点击日志';
COMMENT ON COLUMN hg_short_link_log.id           IS '主键';
COMMENT ON COLUMN hg_short_link_log.link_id      IS '短链ID';
COMMENT ON COLUMN hg_short_link_log.code         IS '短码（冗余，提高查询效率）';
COMMENT ON COLUMN hg_short_link_log.ip           IS '访客 IP';
COMMENT ON COLUMN hg_short_link_log.region       IS '访问地区（省市）';
COMMENT ON COLUMN hg_short_link_log.country      IS '访问国家';
COMMENT ON COLUMN hg_short_link_log.referer      IS '来源完整 URL';
COMMENT ON COLUMN hg_short_link_log.referer_host IS '来源域名';
COMMENT ON COLUMN hg_short_link_log.user_agent   IS 'User-Agent';
COMMENT ON COLUMN hg_short_link_log.device       IS '设备类型：desktop / mobile / tablet';
COMMENT ON COLUMN hg_short_link_log.browser      IS '浏览器';
COMMENT ON COLUMN hg_short_link_log.os           IS '操作系统';
COMMENT ON COLUMN hg_short_link_log.created_at   IS '点击时间';
