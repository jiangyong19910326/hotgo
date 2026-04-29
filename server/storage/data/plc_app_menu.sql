-- ============================================================
-- PLC 应用密钥菜单 (挂在已有 "PLC采集管理" 根目录下)
-- 幂等: name 唯一索引 + ON DUPLICATE KEY UPDATE
-- ============================================================

-- 找到 PLC 根菜单 ID
SET @r = (SELECT id FROM hg_admin_menu WHERE name='Plc' LIMIT 1);

-- Level 2: 应用密钥目录
INSERT INTO `hg_admin_menu`
  (`pid`,`level`,`tree`,`title`,`name`,`path`,`icon`,`type`,`redirect`,
   `permissions`,`permission_name`,`component`,`always_show`,`active_menu`,
   `is_root`,`is_frame`,`frame_src`,`keep_alive`,`hidden`,`affix`,`sort`,
   `remark`,`status`,`updated_at`,`created_at`)
VALUES
  (@r, 2, CONCAT('tr_', @r, ' '), '应用密钥', 'plcApp', 'app', 'KeyOutlined', 1, '/plc/app/index',
   '', '', 'ParentLayout', 1, '',
   0, 0, '', 0, 0, 0, 60,
   '签名验签 AppID/AppSecret', 1, NOW(), NOW())
ON DUPLICATE KEY UPDATE `id` = LAST_INSERT_ID(`id`), `updated_at` = NOW();
SET @ap = LAST_INSERT_ID();

-- Level 3: 列表页
INSERT INTO `hg_admin_menu`
  (`pid`,`level`,`tree`,`title`,`name`,`path`,`icon`,`type`,`redirect`,
   `permissions`,`permission_name`,`component`,`always_show`,`active_menu`,
   `is_root`,`is_frame`,`frame_src`,`keep_alive`,`hidden`,`affix`,`sort`,
   `remark`,`status`,`updated_at`,`created_at`)
VALUES
  (@ap, 3, CONCAT('tr_', @r, ' tr_', @ap, ' '), '应用密钥列表', 'plcAppIndex', 'index', '', 2, '',
   '/plc/app/list', '', '/plc/app/index', 1, 'plcApp',
   0, 0, '', 0, 1, 0, 10,
   '', 1, NOW(), NOW())
ON DUPLICATE KEY UPDATE `id` = LAST_INSERT_ID(`id`), `updated_at` = NOW();
SET @api = LAST_INSERT_ID();

-- Level 4: 操作按钮
INSERT INTO `hg_admin_menu`
  (`pid`,`level`,`tree`,`title`,`name`,`path`,`icon`,`type`,`redirect`,
   `permissions`,`permission_name`,`component`,`always_show`,`active_menu`,
   `is_root`,`is_frame`,`frame_src`,`keep_alive`,`hidden`,`affix`,`sort`,
   `remark`,`status`,`updated_at`,`created_at`)
VALUES
  (@api, 4, CONCAT('tr_', @r, ' tr_', @ap, ' tr_', @api, ' '), '查看应用详情', 'plcAppView', '', '', 3, '',
   '/plc/app/view', '', '', 1, '',
   0, 2, '', 0, 0, 0, 10,
   '', 1, NOW(), NOW())
ON DUPLICATE KEY UPDATE `id` = LAST_INSERT_ID(`id`), `updated_at` = NOW();
SET @apv = LAST_INSERT_ID();

INSERT INTO `hg_admin_menu`
  (`pid`,`level`,`tree`,`title`,`name`,`path`,`icon`,`type`,`redirect`,
   `permissions`,`permission_name`,`component`,`always_show`,`active_menu`,
   `is_root`,`is_frame`,`frame_src`,`keep_alive`,`hidden`,`affix`,`sort`,
   `remark`,`status`,`updated_at`,`created_at`)
VALUES
  (@api, 4, CONCAT('tr_', @r, ' tr_', @ap, ' tr_', @api, ' '), '添加/编辑应用', 'plcAppEdit', '', '', 3, '',
   '/plc/app/edit,/plc/app/genSecret', '', '', 1, '',
   0, 2, '', 0, 0, 0, 20,
   '', 1, NOW(), NOW())
ON DUPLICATE KEY UPDATE `id` = LAST_INSERT_ID(`id`), `updated_at` = NOW();
SET @ape = LAST_INSERT_ID();

INSERT INTO `hg_admin_menu`
  (`pid`,`level`,`tree`,`title`,`name`,`path`,`icon`,`type`,`redirect`,
   `permissions`,`permission_name`,`component`,`always_show`,`active_menu`,
   `is_root`,`is_frame`,`frame_src`,`keep_alive`,`hidden`,`affix`,`sort`,
   `remark`,`status`,`updated_at`,`created_at`)
VALUES
  (@api, 4, CONCAT('tr_', @r, ' tr_', @ap, ' tr_', @api, ' '), '删除应用', 'plcAppDelete', '', '', 3, '',
   '/plc/app/delete', '', '', 1, '',
   0, 2, '', 0, 0, 0, 30,
   '', 1, NOW(), NOW())
ON DUPLICATE KEY UPDATE `id` = LAST_INSERT_ID(`id`), `updated_at` = NOW();
SET @apd = LAST_INSERT_ID();

INSERT INTO `hg_admin_menu`
  (`pid`,`level`,`tree`,`title`,`name`,`path`,`icon`,`type`,`redirect`,
   `permissions`,`permission_name`,`component`,`always_show`,`active_menu`,
   `is_root`,`is_frame`,`frame_src`,`keep_alive`,`hidden`,`affix`,`sort`,
   `remark`,`status`,`updated_at`,`created_at`)
VALUES
  (@api, 4, CONCAT('tr_', @r, ' tr_', @ap, ' tr_', @api, ' '), '更新应用状态', 'plcAppStatus', '', '', 3, '',
   '/plc/app/status', '', '', 1, '',
   0, 2, '', 0, 0, 0, 40,
   '', 1, NOW(), NOW())
ON DUPLICATE KEY UPDATE `id` = LAST_INSERT_ID(`id`), `updated_at` = NOW();
SET @aps = LAST_INSERT_ID();

-- 角色菜单关联 (超管 role_id=2)
INSERT IGNORE INTO `hg_admin_role_menu` (`role_id`, `menu_id`) VALUES (2, @ap);
INSERT IGNORE INTO `hg_admin_role_menu` (`role_id`, `menu_id`) VALUES (2, @api);
INSERT IGNORE INTO `hg_admin_role_menu` (`role_id`, `menu_id`) VALUES (2, @apv);
INSERT IGNORE INTO `hg_admin_role_menu` (`role_id`, `menu_id`) VALUES (2, @ape);
INSERT IGNORE INTO `hg_admin_role_menu` (`role_id`, `menu_id`) VALUES (2, @apd);
INSERT IGNORE INTO `hg_admin_role_menu` (`role_id`, `menu_id`) VALUES (2, @aps);
