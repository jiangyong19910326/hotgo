-- ============================================================
-- HotGo PLC 数据采集模块数据表迁移脚本 (MySQL)
-- 适用版本: HotGo v2.x
-- 协议: MQTT 订阅 (DTU 网关推送), Topic: /dtu/{deviceCode}/data
-- ============================================================

-- --------------------------------------------------------
-- 表结构 `hg_plc_mine`  矿场
-- --------------------------------------------------------
CREATE TABLE IF NOT EXISTS `hg_plc_mine` (
  `id`          int(11)       NOT NULL AUTO_INCREMENT              COMMENT '主键',
  `name`        varchar(64)   NOT NULL DEFAULT ''                  COMMENT '矿场名称',
  `location`    varchar(128)  NOT NULL DEFAULT ''                  COMMENT '地理位置',
  `remark`      varchar(255)  NOT NULL DEFAULT ''                  COMMENT '备注',
  `status`      tinyint(1)    NOT NULL DEFAULT 1                   COMMENT '状态：1启用 2禁用',
  `created_by`  bigint(20)    NOT NULL DEFAULT 0                   COMMENT '创建者',
  `updated_by`  bigint(20)    NOT NULL DEFAULT 0                   COMMENT '更新者',
  `created_at`  datetime      DEFAULT NULL                         COMMENT '创建时间',
  `updated_at`  datetime      DEFAULT NULL                         COMMENT '修改时间',
  `deleted_at`  datetime      DEFAULT NULL                         COMMENT '删除时间',
  PRIMARY KEY (`id`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='矿场';

-- --------------------------------------------------------
-- 表结构 `hg_plc_device`  PLC 设备 (MQTT)
-- host 字段存 DTU 序列号, 与 MQTT Topic /dtu/{host}/data 对应
-- --------------------------------------------------------
CREATE TABLE IF NOT EXISTS `hg_plc_device` (
  `id`          int(11)       NOT NULL AUTO_INCREMENT              COMMENT '主键',
  `mine_id`     int(11)       NOT NULL DEFAULT 0                   COMMENT '所属矿场 ID',
  `name`        varchar(64)   NOT NULL DEFAULT ''                  COMMENT '设备名称',
  `host`        varchar(64)   NOT NULL DEFAULT ''                  COMMENT 'DTU 设备编号 (MQTT topic)',
  `remark`      varchar(255)  NOT NULL DEFAULT ''                  COMMENT '备注',
  `status`      tinyint(1)    NOT NULL DEFAULT 1                   COMMENT '状态：1启用 2禁用',
  `created_by`  bigint(20)    NOT NULL DEFAULT 0                   COMMENT '创建者',
  `updated_by`  bigint(20)    NOT NULL DEFAULT 0                   COMMENT '更新者',
  `created_at`  datetime      DEFAULT NULL                         COMMENT '创建时间',
  `updated_at`  datetime      DEFAULT NULL                         COMMENT '修改时间',
  `deleted_at`  datetime      DEFAULT NULL                         COMMENT '删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_host` (`host`),
  KEY `idx_status` (`status`),
  KEY `idx_mine_id` (`mine_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='PLC 设备 (MQTT)';

-- --------------------------------------------------------
-- 表结构 `hg_plc_point`  数据点定义
-- field 字段直接对应 MQTT payload 中的 key (UPPER_SNAKE)
-- --------------------------------------------------------
CREATE TABLE IF NOT EXISTS `hg_plc_point` (
  `id`          int(11)       NOT NULL AUTO_INCREMENT              COMMENT '主键',
  `device_id`   int(11)       NOT NULL DEFAULT 0                   COMMENT '所属设备 ID',
  `name`        varchar(64)   NOT NULL DEFAULT ''                  COMMENT '点位名称（展示用）',
  `field`       varchar(64)   NOT NULL DEFAULT ''                  COMMENT 'MQTT 字段名 (UPPER_SNAKE)',
  `data_type`   varchar(16)   NOT NULL DEFAULT 'Real'              COMMENT '数据类型：Bool/Real',
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
-- 表结构 `hg_plc_record`  采集历史记录
-- --------------------------------------------------------
CREATE TABLE IF NOT EXISTS `hg_plc_record` (
  `id`          bigint(20)    NOT NULL AUTO_INCREMENT              COMMENT '主键',
  `device_id`   int(11)       NOT NULL DEFAULT 0                   COMMENT '设备 ID',
  `point_id`    int(11)       NOT NULL DEFAULT 0                   COMMENT '数据点 ID',
  `field`       varchar(64)   NOT NULL DEFAULT ''                  COMMENT '字段标识',
  `raw_value`   varchar(64)   NOT NULL DEFAULT ''                  COMMENT '原始值',
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
  `point_name`  varchar(64)   NOT NULL DEFAULT ''                  COMMENT '点位名称（冗余）',
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
