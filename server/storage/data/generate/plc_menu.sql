-- ============================================================
-- HotGo PLC 模块菜单数据 (MySQL)
-- 执行前确保 hg_admin_menu 表已存在
-- ============================================================

SET @now = NOW();

-- ────────────────────────────────────────────────
-- 1. 顶级目录：PLC 管理（type=1，component=LAYOUT）
-- ────────────────────────────────────────────────
INSERT INTO `hg_admin_menu`
  (`pid`, `title`, `name`, `path`, `icon`, `type`, `redirect`,
   `permissions`, `permission_name`, `component`,
   `always_show`, `active_menu`, `is_root`, `is_frame`, `frame_src`,
   `keep_alive`, `hidden`, `affix`, `level`, `tree`, `sort`, `remark`, `status`, `created_at`, `updated_at`)
VALUES
  (0, 'PLC 管理', 'plc', '/plc', 'ControlOutlined', 1, '/plc/device',
   '', '', 'LAYOUT',
   1, '', 0, 0, '',
   0, 0, 0, 1, '', 100, 'PLC 数据采集', 1, @now, @now);

SET @dirId = LAST_INSERT_ID();

-- ────────────────────────────────────────────────
-- 2. 设备管理页（type=2，菜单项）
-- ────────────────────────────────────────────────
INSERT INTO `hg_admin_menu`
  (`pid`, `title`, `name`, `path`, `icon`, `type`, `redirect`,
   `permissions`, `permission_name`, `component`,
   `always_show`, `active_menu`, `is_root`, `is_frame`, `frame_src`,
   `keep_alive`, `hidden`, `affix`, `level`, `tree`, `sort`, `remark`, `status`, `created_at`, `updated_at`)
VALUES
  (@dirId, '设备管理', 'plcDevice', 'device', 'AppstoreOutlined', 2, '',
   '/plc/device/list', '获取PLC设备列表', '/plc/device/index',
   1, '', 0, 0, '',
   0, 0, 0, 2, CONCAT('tr_', @dirId, ' '), 10, '', 1, @now, @now);

SET @deviceMenuId = LAST_INSERT_ID();

-- 设备管理子按钮权限（type=3）
INSERT INTO `hg_admin_menu`
  (`pid`, `title`, `name`, `path`, `icon`, `type`, `redirect`,
   `permissions`, `permission_name`, `component`,
   `always_show`, `active_menu`, `is_root`, `is_frame`, `frame_src`,
   `keep_alive`, `hidden`, `affix`, `level`, `tree`, `sort`, `remark`, `status`, `created_at`, `updated_at`)
VALUES
  (@deviceMenuId, '新增/编辑设备', 'plcDeviceEdit', '', '', 3, '',
   '/plc/device/edit', '新增/修改PLC设备', '',
   1, '', 0, 0, '', 0, 0, 0, 3, CONCAT('tr_', @dirId, ' tr_', @deviceMenuId, ' '), 10, '', 1, @now, @now),
  (@deviceMenuId, '删除设备', 'plcDeviceDelete', '', '', 3, '',
   '/plc/device/delete', '删除PLC设备', '',
   1, '', 0, 0, '', 0, 0, 0, 3, CONCAT('tr_', @dirId, ' tr_', @deviceMenuId, ' '), 20, '', 1, @now, @now),
  (@deviceMenuId, '更新设备状态', 'plcDeviceStatus', '', '', 3, '',
   '/plc/device/status', '更新PLC设备状态', '',
   1, '', 0, 0, '', 0, 0, 0, 3, CONCAT('tr_', @dirId, ' tr_', @deviceMenuId, ' '), 30, '', 1, @now, @now);

-- ────────────────────────────────────────────────
-- 3. 数据点配置（type=2）
-- ────────────────────────────────────────────────
INSERT INTO `hg_admin_menu`
  (`pid`, `title`, `name`, `path`, `icon`, `type`, `redirect`,
   `permissions`, `permission_name`, `component`,
   `always_show`, `active_menu`, `is_root`, `is_frame`, `frame_src`,
   `keep_alive`, `hidden`, `affix`, `level`, `tree`, `sort`, `remark`, `status`, `created_at`, `updated_at`)
VALUES
  (@dirId, '数据点配置', 'plcPoint', 'point', 'SettingOutlined', 2, '',
   '/plc/point/list', '获取数据点列表', '/plc/point/index',
   1, '', 0, 0, '',
   0, 0, 0, 2, CONCAT('tr_', @dirId, ' '), 20, '', 1, @now, @now);

SET @pointMenuId = LAST_INSERT_ID();

INSERT INTO `hg_admin_menu`
  (`pid`, `title`, `name`, `path`, `icon`, `type`, `redirect`,
   `permissions`, `permission_name`, `component`,
   `always_show`, `active_menu`, `is_root`, `is_frame`, `frame_src`,
   `keep_alive`, `hidden`, `affix`, `level`, `tree`, `sort`, `remark`, `status`, `created_at`, `updated_at`)
VALUES
  (@pointMenuId, '新增/编辑数据点', 'plcPointEdit', '', '', 3, '',
   '/plc/point/edit', '新增/修改数据点', '',
   1, '', 0, 0, '', 0, 0, 0, 3, CONCAT('tr_', @dirId, ' tr_', @pointMenuId, ' '), 10, '', 1, @now, @now),
  (@pointMenuId, '删除数据点', 'plcPointDelete', '', '', 3, '',
   '/plc/point/delete', '删除数据点', '',
   1, '', 0, 0, '', 0, 0, 0, 3, CONCAT('tr_', @dirId, ' tr_', @pointMenuId, ' '), 20, '', 1, @now, @now);

-- ────────────────────────────────────────────────
-- 4. 实时监控（type=2）
-- ────────────────────────────────────────────────
INSERT INTO `hg_admin_menu`
  (`pid`, `title`, `name`, `path`, `icon`, `type`, `redirect`,
   `permissions`, `permission_name`, `component`,
   `always_show`, `active_menu`, `is_root`, `is_frame`, `frame_src`,
   `keep_alive`, `hidden`, `affix`, `level`, `tree`, `sort`, `remark`, `status`, `created_at`, `updated_at`)
VALUES
  (@dirId, '实时监控', 'plcDashboard', 'dashboard', 'MonitorOutlined', 2, '',
   '/plc/realtime', '获取实时数据', '/plc/dashboard/index',
   1, '', 0, 0, '',
   0, 0, 0, 2, CONCAT('tr_', @dirId, ' '), 30, '', 1, @now, @now);

-- ────────────────────────────────────────────────
-- 5. 报警记录（type=2）
-- ────────────────────────────────────────────────
INSERT INTO `hg_admin_menu`
  (`pid`, `title`, `name`, `path`, `icon`, `type`, `redirect`,
   `permissions`, `permission_name`, `component`,
   `always_show`, `active_menu`, `is_root`, `is_frame`, `frame_src`,
   `keep_alive`, `hidden`, `affix`, `level`, `tree`, `sort`, `remark`, `status`, `created_at`, `updated_at`)
VALUES
  (@dirId, '报警记录', 'plcAlarm', 'alarm', 'AlertOutlined', 2, '',
   '/plc/alarm/list', '获取报警列表', '/plc/alarm/index',
   1, '', 0, 0, '',
   0, 0, 0, 2, CONCAT('tr_', @dirId, ' '), 40, '', 1, @now, @now);

SET @alarmMenuId = LAST_INSERT_ID();

INSERT INTO `hg_admin_menu`
  (`pid`, `title`, `name`, `path`, `icon`, `type`, `redirect`,
   `permissions`, `permission_name`, `component`,
   `always_show`, `active_menu`, `is_root`, `is_frame`, `frame_src`,
   `keep_alive`, `hidden`, `affix`, `level`, `tree`, `sort`, `remark`, `status`, `created_at`, `updated_at`)
VALUES
  (@alarmMenuId, '处理报警', 'plcAlarmResolve', '', '', 3, '',
   '/plc/alarm/resolve', '处理报警', '',
   1, '', 0, 0, '', 0, 0, 0, 3, CONCAT('tr_', @dirId, ' tr_', @alarmMenuId, ' '), 10, '', 1, @now, @now);
