-- ============================================================
-- HP300 多缸圆锥破碎机 PLC 测点初始化数据 (MQTT版)
-- DTU 编号: 113200006045
-- Topic:    /dtu/113200006045/data
-- 字段名直接对齐 MQTT payload key (UPPER_SNAKE)
-- ============================================================

-- 插入设备 (host = DTU 编号)
INSERT INTO `hg_plc_device`
  (`mine_id`, `name`, `host`, `remark`, `status`, `created_at`, `updated_at`)
VALUES
  (0, 'HP300多缸圆锥破碎机', '113200006045', 'HP300 DTU MQTT接入', 1, NOW(), NOW())
ON DUPLICATE KEY UPDATE `updated_at` = NOW(), `id` = LAST_INSERT_ID(`id`);

SET @dev = LAST_INSERT_ID();

-- ============================================================
-- 报警位 (Bool=1 触发, alarm_max=0.5 用于上限判定)
-- ============================================================
INSERT IGNORE INTO `hg_plc_point`
  (`device_id`,`name`,`field`,`data_type`,`scale`,`offset_val`,`unit`,`alarm_min`,`alarm_max`,`remark`,`sort`,`status`,`created_at`,`updated_at`)
VALUES
  (@dev,'鼓风机故障',          'AL_BLOWER_STARTER',      'Bool',1,0,'',NULL,0.5,'',  10,1,NOW(),NOW()),
  (@dev,'冷风机故障',          'AL_COOLER_STARTER',      'Bool',1,0,'',NULL,0.5,'',  11,1,NOW(),NOW()),
  (@dev,'破碎机故障',          'AL_CRUSHER',             'Bool',1,0,'',NULL,0.5,'',  12,1,NOW(),NOW()),
  (@dev,'非法启动',            'AL_CRUSHER_ILLEGAL_START','Bool',1,0,'',NULL,0.5,'',13,1,NOW(),NOW()),
  (@dev,'破碎机过载',          'AL_CRUSHER_OVERLOAD',    'Bool',1,0,'',NULL,0.5,'',  14,1,NOW(),NOW()),
  (@dev,'排料信号异常',        'AL_DISCHARGER_SIGNAL',   'Bool',1,0,'',NULL,0.5,'',  15,1,NOW(),NOW()),
  (@dev,'紧急停机',            'AL_EMGENCY_STOP',        'Bool',1,0,'',NULL,0.5,'',  16,1,NOW(),NOW()),
  (@dev,'加热器故障',          'AL_HEATER_STARTER',      'Bool',1,0,'',NULL,0.5,'',  17,1,NOW(),NOW()),
  (@dev,'润滑压力高',          'AL_HIGH_LUBE_PRESS',     'Bool',1,0,'',NULL,0.5,'',  18,1,NOW(),NOW()),
  (@dev,'回油温度高',          'AL_HIGH_RETURN_TEMP',    'Bool',1,0,'',NULL,0.5,'',  19,1,NOW(),NOW()),
  (@dev,'液压过滤器堵塞',      'AL_HYD_CLOGGED',         'Bool',1,0,'',NULL,0.5,'',  20,1,NOW(),NOW()),
  (@dev,'液压压力异常',        'AL_HYD_EXC_PRESS_ATT',   'Bool',1,0,'',NULL,0.5,'',  21,1,NOW(),NOW()),
  (@dev,'液压液位低',          'AL_HYD_LS_LOW',          'Bool',1,0,'',NULL,0.5,'',  22,1,NOW(),NOW()),
  (@dev,'液压压力超时',        'AL_HYD_PRESS_TO',        'Bool',1,0,'',NULL,0.5,'',  23,1,NOW(),NOW()),
  (@dev,'液压泵故障',          'AL_HYD_PUMP',            'Bool',1,0,'',NULL,0.5,'',  24,1,NOW(),NOW()),
  (@dev,'液压压力低',          'AL_LOW_HYD_PRESS',       'Bool',1,0,'',NULL,0.5,'',  25,1,NOW(),NOW()),
  (@dev,'润滑压力低',          'AL_LOW_LUBE_PRESS',      'Bool',1,0,'',NULL,0.5,'',  26,1,NOW(),NOW()),
  (@dev,'回油温度低',          'AL_LOW_RETURN_TEMP',     'Bool',1,0,'',NULL,0.5,'',  27,1,NOW(),NOW()),
  (@dev,'润滑过滤器堵塞',      'AL_LUBE_CLOGGED',        'Bool',1,0,'',NULL,0.5,'',  28,1,NOW(),NOW()),
  (@dev,'润滑液位低',          'AL_LUBE_LS_LOW',         'Bool',1,0,'',NULL,0.5,'',  29,1,NOW(),NOW()),
  (@dev,'润滑压力故障',        'AL_LUBE_PRESS_FAULT',    'Bool',1,0,'',NULL,0.5,'',  30,1,NOW(),NOW()),
  (@dev,'润滑压力传感器故障',  'AL_LUBE_PRESS_SENSOR',   'Bool',1,0,'',NULL,0.5,'',  31,1,NOW(),NOW()),
  (@dev,'润滑泵故障',          'AL_LUBE_PUMP_STARTER',   'Bool',1,0,'',NULL,0.5,'',  32,1,NOW(),NOW()),
  (@dev,'润滑回油液位低',      'AL_LUBE_RETURN_LS_LOW',  'Bool',1,0,'',NULL,0.5,'',  33,1,NOW(),NOW()),
  (@dev,'PS1B压力低',          'AL_PS1B_PRESS_LOW',      'Bool',1,0,'',NULL,0.5,'',  34,1,NOW(),NOW()),
  (@dev,'PS1压力传感器故障',   'AL_PS1_PRESS_SENSOR',    'Bool',1,0,'',NULL,0.5,'',  35,1,NOW(),NOW()),
  (@dev,'PS2B压力低',          'AL_PS2B_PRESS_LOW',      'Bool',1,0,'',NULL,0.5,'',  36,1,NOW(),NOW()),
  (@dev,'PS2压力传感器故障',   'AL_PS2_PRESS_SENSOR',    'Bool',1,0,'',NULL,0.5,'',  37,1,NOW(),NOW()),
  (@dev,'回油温度传感器故障',  'AL_RETURN_TEMP_SENSOR',  'Bool',1,0,'',NULL,0.5,'',  38,1,NOW(),NOW()),
  (@dev,'油箱温度传感器故障',  'AL_TANK_TEMP_SENSOR',    'Bool',1,0,'',NULL,0.5,'',  39,1,NOW(),NOW());

-- ============================================================
-- 运行状态 (Bool, 不报警)
-- ============================================================
INSERT IGNORE INTO `hg_plc_point`
  (`device_id`,`name`,`field`,`data_type`,`scale`,`offset_val`,`unit`,`alarm_min`,`alarm_max`,`remark`,`sort`,`status`,`created_at`,`updated_at`)
VALUES
  (@dev,'冷风电机',            'COOLER_MOTOR',           'Bool',1,0,'',NULL,NULL,'',  50,1,NOW(),NOW()),
  (@dev,'破碎机电机',          'CRUSHER_MOTOR',          'Bool',1,0,'',NULL,NULL,'',  51,1,NOW(),NOW()),
  (@dev,'排料皮带',            'DISCHARGER_CONVEYOR',    'Bool',1,0,'',NULL,NULL,'',  52,1,NOW(),NOW()),
  (@dev,'给料机电机',          'FEEDER_MOTER',           'Bool',1,0,'',NULL,NULL,'',  53,1,NOW(),NOW()),
  (@dev,'液压过滤器堵塞状态',  'HYD_CLOGGED',            'Bool',1,0,'',NULL,NULL,'',  54,1,NOW(),NOW()),
  (@dev,'液压液位状态',        'HYD_LS',                 'Bool',1,0,'',NULL,NULL,'',  55,1,NOW(),NOW()),
  (@dev,'液压泵电机',          'HYD_PUMP_MOTOR',         'Bool',1,0,'',NULL,NULL,'',  56,1,NOW(),NOW()),
  (@dev,'润滑过滤器堵塞状态',  'LUBE_CLOGGED',           'Bool',1,0,'',NULL,NULL,'',  57,1,NOW(),NOW()),
  (@dev,'润滑液位状态',        'LUBE_LS',                'Bool',1,0,'',NULL,NULL,'',  58,1,NOW(),NOW()),
  (@dev,'润滑泵电机',          'LUBE_PUMP_MOTOR',        'Bool',1,0,'',NULL,NULL,'',  59,1,NOW(),NOW()),
  (@dev,'润滑回油液位状态',    'LUBE_RETURN_LS',         'Bool',1,0,'',NULL,NULL,'',  60,1,NOW(),NOW()),
  (@dev,'电磁阀1',             'SOL1',                   'Bool',1,0,'',NULL,NULL,'',  61,1,NOW(),NOW()),
  (@dev,'电磁阀2-3',           'SOL2-3',                 'Bool',1,0,'',NULL,NULL,'',  62,1,NOW(),NOW()),
  (@dev,'电磁阀4',             'SOL4',                   'Bool',1,0,'',NULL,NULL,'',  63,1,NOW(),NOW()),
  (@dev,'电磁阀5',             'SOL5',                   'Bool',1,0,'',NULL,NULL,'',  64,1,NOW(),NOW()),
  (@dev,'电磁阀6',             'SOL6',                   'Bool',1,0,'',NULL,NULL,'',  65,1,NOW(),NOW()),
  (@dev,'电磁阀7',             'SOL7',                   'Bool',1,0,'',NULL,NULL,'',  66,1,NOW(),NOW()),
  (@dev,'电磁阀8',             'SOL8',                   'Bool',1,0,'',NULL,NULL,'',  67,1,NOW(),NOW());

-- ============================================================
-- 模拟量 (Real)
-- ============================================================
INSERT IGNORE INTO `hg_plc_point`
  (`device_id`,`name`,`field`,`data_type`,`scale`,`offset_val`,`unit`,`alarm_min`,`alarm_max`,`remark`,`sort`,`status`,`created_at`,`updated_at`)
VALUES
  (@dev,'破碎机电流',          'CRUSHER_AMPERE',         'Real',1,0,'A',  NULL,NULL,'',  100,1,NOW(),NOW()),
  (@dev,'破碎机运行小时',      'CRUSHER_TIME_H',         'Real',1,0,'h',  NULL,NULL,'',  101,1,NOW(),NOW()),
  (@dev,'破碎机运行分钟',      'CRUSHER_TIME_M',         'Real',1,0,'min',NULL,NULL,'',  102,1,NOW(),NOW()),
  (@dev,'润滑滤芯运行小时',    'LUBE_FILTER_TIME_H',     'Real',1,0,'h',  NULL,NULL,'',  103,1,NOW(),NOW()),
  (@dev,'润滑滤芯运行分钟',    'LUBE_FILTER_TIME_M',     'Real',1,0,'min',NULL,NULL,'',  104,1,NOW(),NOW()),
  (@dev,'润滑泵运行小时',      'LUBE_PUMP_TIME_H',       'Real',1,0,'h',  NULL,NULL,'',  105,1,NOW(),NOW()),
  (@dev,'润滑泵运行分钟',      'LUBE_PUMP_TIME_M',       'Real',1,0,'min',NULL,NULL,'',  106,1,NOW(),NOW()),
  (@dev,'润滑滤后压力',        'LUBE_POST_FILLTER_PRESS','Real',1,0,'bar',NULL,NULL,'',  107,1,NOW(),NOW()),
  (@dev,'润滑滤前压力',        'LUBE_PRE_FILLTER_PRESS', 'Real',1,0,'bar',NULL,NULL,'',  108,1,NOW(),NOW()),
  (@dev,'润滑回油温度',        'LUBE_RETURN_TEMP',       'Real',1,0,'℃', NULL,NULL,'',  109,1,NOW(),NOW()),
  (@dev,'润滑油箱温度',        'LUBE_TANK_TEMP',         'Real',1,0,'℃', NULL,NULL,'',  110,1,NOW(),NOW()),
  (@dev,'PS1压力',             'PS1_PRESS',              'Real',1,0,'bar',NULL,NULL,'',  111,1,NOW(),NOW()),
  (@dev,'PS2压力',             'PS2_PRESS',              'Real',1,0,'bar',NULL,NULL,'',  112,1,NOW(),NOW());

-- ============================================================
-- 控制位 (BLOWER_STARTER 等控制按钮)
-- ============================================================
INSERT IGNORE INTO `hg_plc_point`
  (`device_id`,`name`,`field`,`data_type`,`scale`,`offset_val`,`unit`,`alarm_min`,`alarm_max`,`remark`,`sort`,`status`,`created_at`,`updated_at`)
VALUES
  (@dev,'鼓风机启动',          'BLOWER_STARTER',         'Bool',1,0,'',NULL,NULL,'',  200,1,NOW(),NOW());
