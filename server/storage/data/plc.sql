-- ============================================================
-- HotGo PLC 数据采集模块数据表迁移脚本 (MySQL)
-- 适用版本: HotGo v2.x
-- 使用方式: 可直接在已有数据库中执行，使用 IF NOT EXISTS 安全幂等
-- ============================================================

-- --------------------------------------------------------
-- 表结构 `hg_plc_device`  PLC 设备
-- --------------------------------------------------------
CREATE TABLE IF NOT EXISTS `hg_plc_device` (
  `id`          int(11)       NOT NULL AUTO_INCREMENT              COMMENT '主键',
  `name`        varchar(64)   NOT NULL DEFAULT ''                  COMMENT '设备名称',
  `host`        varchar(64)   NOT NULL DEFAULT ''                  COMMENT 'IP 地址',
  `port`        int(11)       NOT NULL DEFAULT 102                 COMMENT 'TCP 端口（S7 默认 102）',
  `rack`        tinyint(4)    NOT NULL DEFAULT 0                   COMMENT '机架号（Rack，通常为 0）',
  `slot`        tinyint(4)    NOT NULL DEFAULT 1                   COMMENT '槽号（Slot，S7-200 SMART 为 1）',
  `interval_ms` int(11)       NOT NULL DEFAULT 1000               COMMENT '采集间隔（毫秒）',
  `remark`      varchar(255)  NOT NULL DEFAULT ''                  COMMENT '备注',
  `status`      tinyint(1)    NOT NULL DEFAULT 1                   COMMENT '状态：1启用 2禁用',
  `created_by`  bigint(20)    NOT NULL DEFAULT 0                   COMMENT '创建者',
  `updated_by`  bigint(20)    NOT NULL DEFAULT 0                   COMMENT '更新者',
  `created_at`  datetime      DEFAULT NULL                         COMMENT '创建时间',
  `updated_at`  datetime      DEFAULT NULL                         COMMENT '修改时间',
  `deleted_at`  datetime      DEFAULT NULL                         COMMENT '删除时间',
  PRIMARY KEY (`id`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='PLC 设备';

-- --------------------------------------------------------
-- 表结构 `hg_plc_point`  数据点定义
-- --------------------------------------------------------
CREATE TABLE IF NOT EXISTS `hg_plc_point` (
  `id`          int(11)       NOT NULL AUTO_INCREMENT              COMMENT '主键',
  `device_id`   int(11)       NOT NULL DEFAULT 0                   COMMENT '所属设备 ID',
  `name`        varchar(64)   NOT NULL DEFAULT ''                  COMMENT '点位名称（展示用）',
  `field`       varchar(64)   NOT NULL DEFAULT ''                  COMMENT '字段标识（英文，前端使用）',
  `area`        varchar(8)    NOT NULL DEFAULT 'DB'                COMMENT '存储区：DB/M/I/Q/V',
  `db_number`   int(11)       NOT NULL DEFAULT 1                   COMMENT 'DB 块号（area=DB 时有效）',
  `byte_offset` int(11)       NOT NULL DEFAULT 0                   COMMENT '字节偏移量',
  `bit_offset`  tinyint(4)    NOT NULL DEFAULT 0                   COMMENT '位偏移量（数据类型为 Bool 时有效）',
  `data_type`   varchar(16)   NOT NULL DEFAULT 'REAL'              COMMENT '数据类型：Bool/Int/DInt/Real/Word/DWord/Byte/String',
  `scale`       double        NOT NULL DEFAULT 1                   COMMENT '换算系数（工程值 = 原始值 × scale + offset）',
  `offset_val`  double        NOT NULL DEFAULT 0                   COMMENT '换算偏移量',
  `unit`        varchar(16)   NOT NULL DEFAULT ''                  COMMENT '单位（如 ℃、bar、rpm）',
  `alarm_min`   double        DEFAULT NULL                         COMMENT '报警下限（NULL 表示不启用）',
  `alarm_max`   double        DEFAULT NULL                         COMMENT '报警上限（NULL 表示不启用）',
  `remark`      varchar(255)  NOT NULL DEFAULT ''                  COMMENT '备注',
  `sort`        int(11)       NOT NULL DEFAULT 0                   COMMENT '排序',
  `status`      tinyint(1)    NOT NULL DEFAULT 1                   COMMENT '状态：1启用 2禁用',
  `created_at`  datetime      DEFAULT NULL                         COMMENT '创建时间',
  `updated_at`  datetime      DEFAULT NULL                         COMMENT '修改时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_device_field` (`device_id`, `field`),
  KEY `idx_device_id` (`device_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='PLC 数据点定义';

-- --------------------------------------------------------
-- 表结构 `hg_plc_record`  采集历史记录（按时间查询，建议按月分区或定期归档）
-- --------------------------------------------------------
CREATE TABLE IF NOT EXISTS `hg_plc_record` (
  `id`          bigint(20)    NOT NULL AUTO_INCREMENT              COMMENT '主键',
  `device_id`   int(11)       NOT NULL DEFAULT 0                   COMMENT '设备 ID',
  `point_id`    int(11)       NOT NULL DEFAULT 0                   COMMENT '数据点 ID',
  `field`       varchar(64)   NOT NULL DEFAULT ''                  COMMENT '字段标识',
  `raw_value`   varchar(64)   NOT NULL DEFAULT ''                  COMMENT '原始值（字符串存储，兼容各类型）',
  `eng_value`   double        DEFAULT NULL                         COMMENT '工程值（换算后）',
  `collected_at` datetime     NOT NULL                             COMMENT '采集时间',
  PRIMARY KEY (`id`),
  KEY `idx_point_time` (`point_id`, `collected_at`),
  KEY `idx_device_time` (`device_id`, `collected_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='PLC 采集历史记录';

-- --------------------------------------------------------
-- 表结构 `hg_plc_alarm`  报警记录
-- --------------------------------------------------------
CREATE TABLE IF NOT EXISTS `hg_plc_alarm` (
  `id`          bigint(20)    NOT NULL AUTO_INCREMENT              COMMENT '主键',
  `device_id`   int(11)       NOT NULL DEFAULT 0                   COMMENT '设备 ID',
  `point_id`    int(11)       NOT NULL DEFAULT 0                   COMMENT '数据点 ID',
  `field`       varchar(64)   NOT NULL DEFAULT ''                  COMMENT '字段标识',
  `point_name`  varchar(64)   NOT NULL DEFAULT ''                  COMMENT '点位名称（冗余，方便查询）',
  `eng_value`   double        NOT NULL DEFAULT 0                   COMMENT '触发时工程值',
  `alarm_type`  tinyint(1)    NOT NULL DEFAULT 1                   COMMENT '报警类型：1超上限 2超下限',
  `alarm_min`   double        DEFAULT NULL                         COMMENT '报警下限快照',
  `alarm_max`   double        DEFAULT NULL                         COMMENT '报警上限快照',
  `unit`        varchar(16)   NOT NULL DEFAULT ''                  COMMENT '单位',
  `is_resolved` tinyint(1)    NOT NULL DEFAULT 2                   COMMENT '是否已处理：1是 2否',
  `resolved_at` datetime      DEFAULT NULL                         COMMENT '处理时间',
  `resolved_by` bigint(20)    DEFAULT NULL                         COMMENT '处理人 ID',
  `remark`      varchar(255)  NOT NULL DEFAULT ''                  COMMENT '处理备注',
  `triggered_at` datetime     NOT NULL                             COMMENT '触发时间',
  `created_at`  datetime      DEFAULT NULL                         COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_point_id` (`point_id`),
  KEY `idx_device_id` (`device_id`),
  KEY `idx_is_resolved` (`is_resolved`),
  KEY `idx_triggered_at` (`triggered_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='PLC 报警记录';
