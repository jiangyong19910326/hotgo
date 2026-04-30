-- ============================================================
-- HotGo 前端用户管理菜单迁移脚本 (MySQL)
-- 适用版本: HotGo v2.x
--
-- 菜单层级结构：
--   ROOT  前端用户管理   (level=1, 顶级目录, LAYOUT)
--   └──  用户列表         (level=2, 列表页)
--       ├──  添加/编辑
--       ├──  删除
--       ├──  更新状态
--       └──  重置密码
-- ============================================================

-- ─────────────────────────────────────────────────────────────
-- Level 1: 前端用户管理（根目录, LAYOUT）
-- ─────────────────────────────────────────────────────────────
INSERT INTO `hg_admin_menu`
  (`pid`,`level`,`tree`,`title`,`name`,`path`,`icon`,`type`,`redirect`,
   `permissions`,`permission_name`,`component`,`always_show`,`active_menu`,
   `is_root`,`is_frame`,`frame_src`,`keep_alive`,`hidden`,`affix`,`sort`,
   `remark`,`status`,`updated_at`,`created_at`)
VALUES
  (0, 1, '', '前端用户管理', 'FrontUser', '/frontUser', 'UserOutlined', 1, '/frontUser/index',
   '', '', 'LAYOUT', 1, '',
   0, 0, '', 0, 0, 0, 230,
   '', 1, NOW(), NOW())
ON DUPLICATE KEY UPDATE `id` = LAST_INSERT_ID(`id`), `updated_at` = NOW();
SET @fr = LAST_INSERT_ID();

-- ─────────────────────────────────────────────────────────────
-- Level 2: 用户列表
-- ─────────────────────────────────────────────────────────────
INSERT INTO `hg_admin_menu`
  (`pid`,`level`,`tree`,`title`,`name`,`path`,`icon`,`type`,`redirect`,
   `permissions`,`permission_name`,`component`,`always_show`,`active_menu`,
   `is_root`,`is_frame`,`frame_src`,`keep_alive`,`hidden`,`affix`,`sort`,
   `remark`,`status`,`updated_at`,`created_at`)
VALUES
  (@fr, 2, CONCAT('tr_', @fr, ' '), '用户列表', 'frontUserIndex', 'index', 'TeamOutlined', 2, '',
   '/frontUser/list,/frontUser/view', '', '/frontUser/index', 1, '',
   0, 0, '', 0, 0, 0, 10,
   '', 1, NOW(), NOW())
ON DUPLICATE KEY UPDATE `id` = LAST_INSERT_ID(`id`), `updated_at` = NOW();
SET @fri = LAST_INSERT_ID();

-- ─────────────────────────────────────────────────────────────
-- Level 3: 操作按钮
-- ─────────────────────────────────────────────────────────────
INSERT INTO `hg_admin_menu`
  (`pid`,`level`,`tree`,`title`,`name`,`path`,`icon`,`type`,`redirect`,
   `permissions`,`permission_name`,`component`,`always_show`,`active_menu`,
   `is_root`,`is_frame`,`frame_src`,`keep_alive`,`hidden`,`affix`,`sort`,
   `remark`,`status`,`updated_at`,`created_at`)
VALUES
  (@fri, 3, CONCAT('tr_', @fr, ' tr_', @fri, ' '), '添加/编辑前端用户', 'frontUserEdit', '', '', 3, '',
   '/frontUser/edit', '', '', 1, '',
   0, 2, '', 0, 0, 0, 20,
   '', 1, NOW(), NOW())
ON DUPLICATE KEY UPDATE `id` = LAST_INSERT_ID(`id`), `updated_at` = NOW();
SET @fre = LAST_INSERT_ID();

INSERT INTO `hg_admin_menu`
  (`pid`,`level`,`tree`,`title`,`name`,`path`,`icon`,`type`,`redirect`,
   `permissions`,`permission_name`,`component`,`always_show`,`active_menu`,
   `is_root`,`is_frame`,`frame_src`,`keep_alive`,`hidden`,`affix`,`sort`,
   `remark`,`status`,`updated_at`,`created_at`)
VALUES
  (@fri, 3, CONCAT('tr_', @fr, ' tr_', @fri, ' '), '删除前端用户', 'frontUserDelete', '', '', 3, '',
   '/frontUser/delete', '', '', 1, '',
   0, 2, '', 0, 0, 0, 30,
   '', 1, NOW(), NOW())
ON DUPLICATE KEY UPDATE `id` = LAST_INSERT_ID(`id`), `updated_at` = NOW();
SET @frd = LAST_INSERT_ID();

INSERT INTO `hg_admin_menu`
  (`pid`,`level`,`tree`,`title`,`name`,`path`,`icon`,`type`,`redirect`,
   `permissions`,`permission_name`,`component`,`always_show`,`active_menu`,
   `is_root`,`is_frame`,`frame_src`,`keep_alive`,`hidden`,`affix`,`sort`,
   `remark`,`status`,`updated_at`,`created_at`)
VALUES
  (@fri, 3, CONCAT('tr_', @fr, ' tr_', @fri, ' '), '更新前端用户状态', 'frontUserStatus', '', '', 3, '',
   '/frontUser/status', '', '', 1, '',
   0, 2, '', 0, 0, 0, 40,
   '', 1, NOW(), NOW())
ON DUPLICATE KEY UPDATE `id` = LAST_INSERT_ID(`id`), `updated_at` = NOW();
SET @frs = LAST_INSERT_ID();

INSERT INTO `hg_admin_menu`
  (`pid`,`level`,`tree`,`title`,`name`,`path`,`icon`,`type`,`redirect`,
   `permissions`,`permission_name`,`component`,`always_show`,`active_menu`,
   `is_root`,`is_frame`,`frame_src`,`keep_alive`,`hidden`,`affix`,`sort`,
   `remark`,`status`,`updated_at`,`created_at`)
VALUES
  (@fri, 3, CONCAT('tr_', @fr, ' tr_', @fri, ' '), '重置前端用户密码', 'frontUserResetPwd', '', '', 3, '',
   '/frontUser/resetPwd', '', '', 1, '',
   0, 2, '', 0, 0, 0, 50,
   '', 1, NOW(), NOW())
ON DUPLICATE KEY UPDATE `id` = LAST_INSERT_ID(`id`), `updated_at` = NOW();
SET @frr = LAST_INSERT_ID();

-- ─────────────────────────────────────────────────────────────
-- 角色菜单关联（超级管理员角色 id=2）
-- ─────────────────────────────────────────────────────────────
INSERT IGNORE INTO `hg_admin_role_menu` (`role_id`, `menu_id`) VALUES (2, @fr);
INSERT IGNORE INTO `hg_admin_role_menu` (`role_id`, `menu_id`) VALUES (2, @fri);
INSERT IGNORE INTO `hg_admin_role_menu` (`role_id`, `menu_id`) VALUES (2, @fre);
INSERT IGNORE INTO `hg_admin_role_menu` (`role_id`, `menu_id`) VALUES (2, @frd);
INSERT IGNORE INTO `hg_admin_role_menu` (`role_id`, `menu_id`) VALUES (2, @frs);
INSERT IGNORE INTO `hg_admin_role_menu` (`role_id`, `menu_id`) VALUES (2, @frr);
