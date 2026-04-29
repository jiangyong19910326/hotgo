-- ============================================================
-- HotGo PLC 模块菜单迁移脚本 (MySQL)
-- 适用版本: HotGo v2.x
--
-- 说明：
--   本脚本不硬编码菜单 ID，由 AUTO_INCREMENT 自动分配。
--   通过 ON DUPLICATE KEY UPDATE id = LAST_INSERT_ID(id) 实现幂等：
--   - 首次执行：正常插入并返回新 ID。
--   - 重复执行：触发 DUPLICATE KEY（name 唯一索引），
--     LAST_INSERT_ID(id) 返回已有行 ID，后续变量引用不受影响。
-- ============================================================
--
-- 菜单层级结构：
--   ROOT  PLC采集管理       (level=1, 顶级目录, LAYOUT)
--   ├──  矿场管理            (level=2, ParentLayout)
--   │   └──  矿场列表        (level=3, 列表页)
--   │       ├──  添加/编辑矿场
--   │       ├──  删除矿场
--   │       ├──  更新矿场状态
--   │       └──  矿场下拉选项
--   ├──  设备管理            (level=2, ParentLayout)
--   │   └──  设备列表        (level=3, 列表页)
--   │       ├──  添加/编辑设备
--   │       ├──  删除设备
--   │       └──  更新设备状态
--   ├──  数据点管理          (level=2, ParentLayout)
--   │   └──  数据点列表      (level=3, 列表页)
--   │       ├──  添加/编辑数据点
--   │       ├──  删除数据点
--   │       └──  更新数据点状态
--   ├──  报警管理            (level=2, ParentLayout)
--   │   └──  报警列表        (level=3, 列表页)
--   │       └──  处理报警
--   └──  实时看板            (level=2, ParentLayout)
--       └──  看板页          (level=3, 列表页)
-- ============================================================

-- ─────────────────────────────────────────────────────────────
-- Level 1: PLC 采集管理（根目录）
-- ─────────────────────────────────────────────────────────────
INSERT INTO `hg_admin_menu`
  (`pid`,`level`,`tree`,`title`,`name`,`path`,`icon`,`type`,`redirect`,
   `permissions`,`permission_name`,`component`,`always_show`,`active_menu`,
   `is_root`,`is_frame`,`frame_src`,`keep_alive`,`hidden`,`affix`,`sort`,
   `remark`,`status`,`updated_at`,`created_at`)
VALUES
  (0, 1, '', 'PLC采集管理', 'Plc', '/plc', 'ControlOutlined', 1, '/plc/mine',
   '', '', 'LAYOUT', 1, '',
   0, 0, '', 0, 0, 0, 220,
   '', 1, NOW(), NOW())
ON DUPLICATE KEY UPDATE `id` = LAST_INSERT_ID(`id`), `updated_at` = NOW();
SET @r = LAST_INSERT_ID();

-- ═══════════════════════════════════════
-- 矿场管理
-- ═══════════════════════════════════════

-- Level 2: 矿场目录
INSERT INTO `hg_admin_menu`
  (`pid`,`level`,`tree`,`title`,`name`,`path`,`icon`,`type`,`redirect`,
   `permissions`,`permission_name`,`component`,`always_show`,`active_menu`,
   `is_root`,`is_frame`,`frame_src`,`keep_alive`,`hidden`,`affix`,`sort`,
   `remark`,`status`,`updated_at`,`created_at`)
VALUES
  (@r, 2, CONCAT('tr_', @r, ' '), '矿场管理', 'plcMine', 'mine', 'EnvironmentOutlined', 1, '/plc/mine/index',
   '', '', 'ParentLayout', 1, '',
   0, 0, '', 0, 0, 0, 10,
   '', 1, NOW(), NOW())
ON DUPLICATE KEY UPDATE `id` = LAST_INSERT_ID(`id`), `updated_at` = NOW();
SET @m = LAST_INSERT_ID();

-- Level 3: 矿场列表页（hidden=1，侧边栏不单独显示）
INSERT INTO `hg_admin_menu`
  (`pid`,`level`,`tree`,`title`,`name`,`path`,`icon`,`type`,`redirect`,
   `permissions`,`permission_name`,`component`,`always_show`,`active_menu`,
   `is_root`,`is_frame`,`frame_src`,`keep_alive`,`hidden`,`affix`,`sort`,
   `remark`,`status`,`updated_at`,`created_at`)
VALUES
  (@m, 3, CONCAT('tr_', @r, ' tr_', @m, ' '), '矿场列表', 'plcMineIndex', 'index', '', 2, '',
   '/plc/mine/list', '', '/plc/mine/index', 1, 'plcMine',
   0, 0, '', 0, 1, 0, 10,
   '', 1, NOW(), NOW())
ON DUPLICATE KEY UPDATE `id` = LAST_INSERT_ID(`id`), `updated_at` = NOW();
SET @mi = LAST_INSERT_ID();

-- Level 4: 矿场操作按钮
INSERT INTO `hg_admin_menu`
  (`pid`,`level`,`tree`,`title`,`name`,`path`,`icon`,`type`,`redirect`,
   `permissions`,`permission_name`,`component`,`always_show`,`active_menu`,
   `is_root`,`is_frame`,`frame_src`,`keep_alive`,`hidden`,`affix`,`sort`,
   `remark`,`status`,`updated_at`,`created_at`)
VALUES
  (@mi, 4, CONCAT('tr_', @r, ' tr_', @m, ' tr_', @mi, ' '), '添加/编辑矿场', 'plcMineEdit', '', '', 3, '',
   '/plc/mine/edit', '', '', 1, '',
   0, 2, '', 0, 0, 0, 20,
   '', 1, NOW(), NOW())
ON DUPLICATE KEY UPDATE `id` = LAST_INSERT_ID(`id`), `updated_at` = NOW();
SET @me = LAST_INSERT_ID();

INSERT INTO `hg_admin_menu`
  (`pid`,`level`,`tree`,`title`,`name`,`path`,`icon`,`type`,`redirect`,
   `permissions`,`permission_name`,`component`,`always_show`,`active_menu`,
   `is_root`,`is_frame`,`frame_src`,`keep_alive`,`hidden`,`affix`,`sort`,
   `remark`,`status`,`updated_at`,`created_at`)
VALUES
  (@mi, 4, CONCAT('tr_', @r, ' tr_', @m, ' tr_', @mi, ' '), '删除矿场', 'plcMineDelete', '', '', 3, '',
   '/plc/mine/delete', '', '', 1, '',
   0, 2, '', 0, 0, 0, 30,
   '', 1, NOW(), NOW())
ON DUPLICATE KEY UPDATE `id` = LAST_INSERT_ID(`id`), `updated_at` = NOW();
SET @md = LAST_INSERT_ID();

INSERT INTO `hg_admin_menu`
  (`pid`,`level`,`tree`,`title`,`name`,`path`,`icon`,`type`,`redirect`,
   `permissions`,`permission_name`,`component`,`always_show`,`active_menu`,
   `is_root`,`is_frame`,`frame_src`,`keep_alive`,`hidden`,`affix`,`sort`,
   `remark`,`status`,`updated_at`,`created_at`)
VALUES
  (@mi, 4, CONCAT('tr_', @r, ' tr_', @m, ' tr_', @mi, ' '), '更新矿场状态', 'plcMineStatus', '', '', 3, '',
   '/plc/mine/status', '', '', 1, '',
   0, 2, '', 0, 0, 0, 40,
   '', 1, NOW(), NOW())
ON DUPLICATE KEY UPDATE `id` = LAST_INSERT_ID(`id`), `updated_at` = NOW();
SET @ms = LAST_INSERT_ID();

INSERT INTO `hg_admin_menu`
  (`pid`,`level`,`tree`,`title`,`name`,`path`,`icon`,`type`,`redirect`,
   `permissions`,`permission_name`,`component`,`always_show`,`active_menu`,
   `is_root`,`is_frame`,`frame_src`,`keep_alive`,`hidden`,`affix`,`sort`,
   `remark`,`status`,`updated_at`,`created_at`)
VALUES
  (@mi, 4, CONCAT('tr_', @r, ' tr_', @m, ' tr_', @mi, ' '), '矿场下拉选项', 'plcMineOptions', '', '', 3, '',
   '/plc/mine/options', '', '', 1, '',
   0, 2, '', 0, 0, 0, 50,
   '', 1, NOW(), NOW())
ON DUPLICATE KEY UPDATE `id` = LAST_INSERT_ID(`id`), `updated_at` = NOW();
SET @mo = LAST_INSERT_ID();

-- ═══════════════════════════════════════
-- 设备管理
-- ═══════════════════════════════════════

-- Level 2: 设备目录
INSERT INTO `hg_admin_menu`
  (`pid`,`level`,`tree`,`title`,`name`,`path`,`icon`,`type`,`redirect`,
   `permissions`,`permission_name`,`component`,`always_show`,`active_menu`,
   `is_root`,`is_frame`,`frame_src`,`keep_alive`,`hidden`,`affix`,`sort`,
   `remark`,`status`,`updated_at`,`created_at`)
VALUES
  (@r, 2, CONCAT('tr_', @r, ' '), '设备管理', 'plcDevice', 'device', 'ApiOutlined', 1, '/plc/device/index',
   '', '', 'ParentLayout', 1, '',
   0, 0, '', 0, 0, 0, 20,
   '', 1, NOW(), NOW())
ON DUPLICATE KEY UPDATE `id` = LAST_INSERT_ID(`id`), `updated_at` = NOW();
SET @dv = LAST_INSERT_ID();

-- Level 3: 设备列表页
INSERT INTO `hg_admin_menu`
  (`pid`,`level`,`tree`,`title`,`name`,`path`,`icon`,`type`,`redirect`,
   `permissions`,`permission_name`,`component`,`always_show`,`active_menu`,
   `is_root`,`is_frame`,`frame_src`,`keep_alive`,`hidden`,`affix`,`sort`,
   `remark`,`status`,`updated_at`,`created_at`)
VALUES
  (@dv, 3, CONCAT('tr_', @r, ' tr_', @dv, ' '), '设备列表', 'plcDeviceIndex', 'index', '', 2, '',
   '/plc/device/list,/plc/mine/options', '', '/plc/device/index', 1, 'plcDevice',
   0, 0, '', 0, 1, 0, 10,
   '', 1, NOW(), NOW())
ON DUPLICATE KEY UPDATE `id` = LAST_INSERT_ID(`id`), `updated_at` = NOW();
SET @dvi = LAST_INSERT_ID();

-- Level 4: 设备操作按钮
INSERT INTO `hg_admin_menu`
  (`pid`,`level`,`tree`,`title`,`name`,`path`,`icon`,`type`,`redirect`,
   `permissions`,`permission_name`,`component`,`always_show`,`active_menu`,
   `is_root`,`is_frame`,`frame_src`,`keep_alive`,`hidden`,`affix`,`sort`,
   `remark`,`status`,`updated_at`,`created_at`)
VALUES
  (@dvi, 4, CONCAT('tr_', @r, ' tr_', @dv, ' tr_', @dvi, ' '), '添加/编辑设备', 'plcDeviceEdit', '', '', 3, '',
   '/plc/device/edit', '', '', 1, '',
   0, 2, '', 0, 0, 0, 20,
   '', 1, NOW(), NOW())
ON DUPLICATE KEY UPDATE `id` = LAST_INSERT_ID(`id`), `updated_at` = NOW();
SET @dve = LAST_INSERT_ID();

INSERT INTO `hg_admin_menu`
  (`pid`,`level`,`tree`,`title`,`name`,`path`,`icon`,`type`,`redirect`,
   `permissions`,`permission_name`,`component`,`always_show`,`active_menu`,
   `is_root`,`is_frame`,`frame_src`,`keep_alive`,`hidden`,`affix`,`sort`,
   `remark`,`status`,`updated_at`,`created_at`)
VALUES
  (@dvi, 4, CONCAT('tr_', @r, ' tr_', @dv, ' tr_', @dvi, ' '), '删除设备', 'plcDeviceDelete', '', '', 3, '',
   '/plc/device/delete', '', '', 1, '',
   0, 2, '', 0, 0, 0, 30,
   '', 1, NOW(), NOW())
ON DUPLICATE KEY UPDATE `id` = LAST_INSERT_ID(`id`), `updated_at` = NOW();
SET @dvd = LAST_INSERT_ID();

INSERT INTO `hg_admin_menu`
  (`pid`,`level`,`tree`,`title`,`name`,`path`,`icon`,`type`,`redirect`,
   `permissions`,`permission_name`,`component`,`always_show`,`active_menu`,
   `is_root`,`is_frame`,`frame_src`,`keep_alive`,`hidden`,`affix`,`sort`,
   `remark`,`status`,`updated_at`,`created_at`)
VALUES
  (@dvi, 4, CONCAT('tr_', @r, ' tr_', @dv, ' tr_', @dvi, ' '), '更新设备状态', 'plcDeviceStatus', '', '', 3, '',
   '/plc/device/status', '', '', 1, '',
   0, 2, '', 0, 0, 0, 40,
   '', 1, NOW(), NOW())
ON DUPLICATE KEY UPDATE `id` = LAST_INSERT_ID(`id`), `updated_at` = NOW();
SET @dvs = LAST_INSERT_ID();

-- ═══════════════════════════════════════
-- 数据点管理
-- ═══════════════════════════════════════

-- Level 2: 数据点目录
INSERT INTO `hg_admin_menu`
  (`pid`,`level`,`tree`,`title`,`name`,`path`,`icon`,`type`,`redirect`,
   `permissions`,`permission_name`,`component`,`always_show`,`active_menu`,
   `is_root`,`is_frame`,`frame_src`,`keep_alive`,`hidden`,`affix`,`sort`,
   `remark`,`status`,`updated_at`,`created_at`)
VALUES
  (@r, 2, CONCAT('tr_', @r, ' '), '数据点管理', 'plcPoint', 'point', 'NodeIndexOutlined', 1, '/plc/point/index',
   '', '', 'ParentLayout', 1, '',
   0, 0, '', 0, 0, 0, 30,
   '', 1, NOW(), NOW())
ON DUPLICATE KEY UPDATE `id` = LAST_INSERT_ID(`id`), `updated_at` = NOW();
SET @pt = LAST_INSERT_ID();

-- Level 3: 数据点列表页
INSERT INTO `hg_admin_menu`
  (`pid`,`level`,`tree`,`title`,`name`,`path`,`icon`,`type`,`redirect`,
   `permissions`,`permission_name`,`component`,`always_show`,`active_menu`,
   `is_root`,`is_frame`,`frame_src`,`keep_alive`,`hidden`,`affix`,`sort`,
   `remark`,`status`,`updated_at`,`created_at`)
VALUES
  (@pt, 3, CONCAT('tr_', @r, ' tr_', @pt, ' '), '数据点列表', 'plcPointIndex', 'index', '', 2, '',
   '/plc/point/list', '', '/plc/point/index', 1, 'plcPoint',
   0, 0, '', 0, 1, 0, 10,
   '', 1, NOW(), NOW())
ON DUPLICATE KEY UPDATE `id` = LAST_INSERT_ID(`id`), `updated_at` = NOW();
SET @pti = LAST_INSERT_ID();

-- Level 4: 数据点操作按钮
INSERT INTO `hg_admin_menu`
  (`pid`,`level`,`tree`,`title`,`name`,`path`,`icon`,`type`,`redirect`,
   `permissions`,`permission_name`,`component`,`always_show`,`active_menu`,
   `is_root`,`is_frame`,`frame_src`,`keep_alive`,`hidden`,`affix`,`sort`,
   `remark`,`status`,`updated_at`,`created_at`)
VALUES
  (@pti, 4, CONCAT('tr_', @r, ' tr_', @pt, ' tr_', @pti, ' '), '添加/编辑数据点', 'plcPointEdit', '', '', 3, '',
   '/plc/point/edit', '', '', 1, '',
   0, 2, '', 0, 0, 0, 20,
   '', 1, NOW(), NOW())
ON DUPLICATE KEY UPDATE `id` = LAST_INSERT_ID(`id`), `updated_at` = NOW();
SET @pte = LAST_INSERT_ID();

INSERT INTO `hg_admin_menu`
  (`pid`,`level`,`tree`,`title`,`name`,`path`,`icon`,`type`,`redirect`,
   `permissions`,`permission_name`,`component`,`always_show`,`active_menu`,
   `is_root`,`is_frame`,`frame_src`,`keep_alive`,`hidden`,`affix`,`sort`,
   `remark`,`status`,`updated_at`,`created_at`)
VALUES
  (@pti, 4, CONCAT('tr_', @r, ' tr_', @pt, ' tr_', @pti, ' '), '删除数据点', 'plcPointDelete', '', '', 3, '',
   '/plc/point/delete', '', '', 1, '',
   0, 2, '', 0, 0, 0, 30,
   '', 1, NOW(), NOW())
ON DUPLICATE KEY UPDATE `id` = LAST_INSERT_ID(`id`), `updated_at` = NOW();
SET @ptd = LAST_INSERT_ID();

INSERT INTO `hg_admin_menu`
  (`pid`,`level`,`tree`,`title`,`name`,`path`,`icon`,`type`,`redirect`,
   `permissions`,`permission_name`,`component`,`always_show`,`active_menu`,
   `is_root`,`is_frame`,`frame_src`,`keep_alive`,`hidden`,`affix`,`sort`,
   `remark`,`status`,`updated_at`,`created_at`)
VALUES
  (@pti, 4, CONCAT('tr_', @r, ' tr_', @pt, ' tr_', @pti, ' '), '更新数据点状态', 'plcPointStatus', '', '', 3, '',
   '/plc/point/status', '', '', 1, '',
   0, 2, '', 0, 0, 0, 40,
   '', 1, NOW(), NOW())
ON DUPLICATE KEY UPDATE `id` = LAST_INSERT_ID(`id`), `updated_at` = NOW();
SET @pts = LAST_INSERT_ID();

-- ═══════════════════════════════════════
-- 报警管理
-- ═══════════════════════════════════════

-- Level 2: 报警目录
INSERT INTO `hg_admin_menu`
  (`pid`,`level`,`tree`,`title`,`name`,`path`,`icon`,`type`,`redirect`,
   `permissions`,`permission_name`,`component`,`always_show`,`active_menu`,
   `is_root`,`is_frame`,`frame_src`,`keep_alive`,`hidden`,`affix`,`sort`,
   `remark`,`status`,`updated_at`,`created_at`)
VALUES
  (@r, 2, CONCAT('tr_', @r, ' '), '报警管理', 'plcAlarm', 'alarm', 'AlertOutlined', 1, '/plc/alarm/index',
   '', '', 'ParentLayout', 1, '',
   0, 0, '', 0, 0, 0, 40,
   '', 1, NOW(), NOW())
ON DUPLICATE KEY UPDATE `id` = LAST_INSERT_ID(`id`), `updated_at` = NOW();
SET @al = LAST_INSERT_ID();

-- Level 3: 报警列表页
INSERT INTO `hg_admin_menu`
  (`pid`,`level`,`tree`,`title`,`name`,`path`,`icon`,`type`,`redirect`,
   `permissions`,`permission_name`,`component`,`always_show`,`active_menu`,
   `is_root`,`is_frame`,`frame_src`,`keep_alive`,`hidden`,`affix`,`sort`,
   `remark`,`status`,`updated_at`,`created_at`)
VALUES
  (@al, 3, CONCAT('tr_', @r, ' tr_', @al, ' '), '报警列表', 'plcAlarmIndex', 'index', '', 2, '',
   '/plc/alarm/list', '', '/plc/alarm/index', 1, 'plcAlarm',
   0, 0, '', 0, 1, 0, 10,
   '', 1, NOW(), NOW())
ON DUPLICATE KEY UPDATE `id` = LAST_INSERT_ID(`id`), `updated_at` = NOW();
SET @ali = LAST_INSERT_ID();

-- Level 4: 报警操作按钮
INSERT INTO `hg_admin_menu`
  (`pid`,`level`,`tree`,`title`,`name`,`path`,`icon`,`type`,`redirect`,
   `permissions`,`permission_name`,`component`,`always_show`,`active_menu`,
   `is_root`,`is_frame`,`frame_src`,`keep_alive`,`hidden`,`affix`,`sort`,
   `remark`,`status`,`updated_at`,`created_at`)
VALUES
  (@ali, 4, CONCAT('tr_', @r, ' tr_', @al, ' tr_', @ali, ' '), '处理报警', 'plcAlarmResolve', '', '', 3, '',
   '/plc/alarm/resolve', '', '', 1, '',
   0, 2, '', 0, 0, 0, 20,
   '', 1, NOW(), NOW())
ON DUPLICATE KEY UPDATE `id` = LAST_INSERT_ID(`id`), `updated_at` = NOW();
SET @alr = LAST_INSERT_ID();

-- ═══════════════════════════════════════
-- 实时看板
-- ═══════════════════════════════════════

-- Level 2: 看板目录
INSERT INTO `hg_admin_menu`
  (`pid`,`level`,`tree`,`title`,`name`,`path`,`icon`,`type`,`redirect`,
   `permissions`,`permission_name`,`component`,`always_show`,`active_menu`,
   `is_root`,`is_frame`,`frame_src`,`keep_alive`,`hidden`,`affix`,`sort`,
   `remark`,`status`,`updated_at`,`created_at`)
VALUES
  (@r, 2, CONCAT('tr_', @r, ' '), '实时看板', 'plcDashboard', 'dashboard', 'BarChartOutlined', 1, '/plc/dashboard/index',
   '', '', 'ParentLayout', 1, '',
   0, 0, '', 0, 0, 0, 50,
   '', 1, NOW(), NOW())
ON DUPLICATE KEY UPDATE `id` = LAST_INSERT_ID(`id`), `updated_at` = NOW();
SET @db = LAST_INSERT_ID();

-- Level 3: 看板页
INSERT INTO `hg_admin_menu`
  (`pid`,`level`,`tree`,`title`,`name`,`path`,`icon`,`type`,`redirect`,
   `permissions`,`permission_name`,`component`,`always_show`,`active_menu`,
   `is_root`,`is_frame`,`frame_src`,`keep_alive`,`hidden`,`affix`,`sort`,
   `remark`,`status`,`updated_at`,`created_at`)
VALUES
  (@db, 3, CONCAT('tr_', @r, ' tr_', @db, ' '), '实时看板', 'plcDashboardIndex', 'index', '', 2, '',
   '/plc/realtime,/plc/device/list', '', '/plc/dashboard/index', 1, 'plcDashboard',
   0, 0, '', 0, 1, 0, 10,
   '', 1, NOW(), NOW())
ON DUPLICATE KEY UPDATE `id` = LAST_INSERT_ID(`id`), `updated_at` = NOW();
SET @dbi = LAST_INSERT_ID();

-- ─────────────────────────────────────────────────────────────
-- 角色菜单关联（超级管理员角色 id=2）
-- 使用上方捕获的变量，INSERT IGNORE 防止重复关联
-- ─────────────────────────────────────────────────────────────
INSERT IGNORE INTO `hg_admin_role_menu` (`role_id`, `menu_id`) VALUES (2, @r);
INSERT IGNORE INTO `hg_admin_role_menu` (`role_id`, `menu_id`) VALUES (2, @m);
INSERT IGNORE INTO `hg_admin_role_menu` (`role_id`, `menu_id`) VALUES (2, @mi);
INSERT IGNORE INTO `hg_admin_role_menu` (`role_id`, `menu_id`) VALUES (2, @me);
INSERT IGNORE INTO `hg_admin_role_menu` (`role_id`, `menu_id`) VALUES (2, @md);
INSERT IGNORE INTO `hg_admin_role_menu` (`role_id`, `menu_id`) VALUES (2, @ms);
INSERT IGNORE INTO `hg_admin_role_menu` (`role_id`, `menu_id`) VALUES (2, @mo);
INSERT IGNORE INTO `hg_admin_role_menu` (`role_id`, `menu_id`) VALUES (2, @dv);
INSERT IGNORE INTO `hg_admin_role_menu` (`role_id`, `menu_id`) VALUES (2, @dvi);
INSERT IGNORE INTO `hg_admin_role_menu` (`role_id`, `menu_id`) VALUES (2, @dve);
INSERT IGNORE INTO `hg_admin_role_menu` (`role_id`, `menu_id`) VALUES (2, @dvd);
INSERT IGNORE INTO `hg_admin_role_menu` (`role_id`, `menu_id`) VALUES (2, @dvs);
INSERT IGNORE INTO `hg_admin_role_menu` (`role_id`, `menu_id`) VALUES (2, @pt);
INSERT IGNORE INTO `hg_admin_role_menu` (`role_id`, `menu_id`) VALUES (2, @pti);
INSERT IGNORE INTO `hg_admin_role_menu` (`role_id`, `menu_id`) VALUES (2, @pte);
INSERT IGNORE INTO `hg_admin_role_menu` (`role_id`, `menu_id`) VALUES (2, @ptd);
INSERT IGNORE INTO `hg_admin_role_menu` (`role_id`, `menu_id`) VALUES (2, @pts);
INSERT IGNORE INTO `hg_admin_role_menu` (`role_id`, `menu_id`) VALUES (2, @al);
INSERT IGNORE INTO `hg_admin_role_menu` (`role_id`, `menu_id`) VALUES (2, @ali);
INSERT IGNORE INTO `hg_admin_role_menu` (`role_id`, `menu_id`) VALUES (2, @alr);
INSERT IGNORE INTO `hg_admin_role_menu` (`role_id`, `menu_id`) VALUES (2, @db);
INSERT IGNORE INTO `hg_admin_role_menu` (`role_id`, `menu_id`) VALUES (2, @dbi);
