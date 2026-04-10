# HotGo 功能清单文档

> **框架版本**：v2.0
> **文档日期**：2026-03-19
> **技术栈**：Go（GoFrame 2.x）+ Vue 3（TypeScript + NaiveUI）

---

## 目录

1. [认证与安全](#一认证与安全)
2. [用户管理](#二用户管理)
3. [组织架构](#三组织架构)
4. [权限管理](#四权限管理)
5. [系统配置](#五系统配置)
6. [日志管理](#六日志管理)
7. [监控管理](#七监控管理)
8. [通知与公告](#八通知与公告)
9. [支付与财务](#九支付与财务)
10. [文件与附件](#十文件与附件)
11. [定时任务](#十一定时任务)
12. [消息队列](#十二消息队列)
13. [WebSocket 实时通信](#十三websocket-实时通信)
14. [代码生成器](#十四代码生成器)
15. [插件系统](#十五插件系统)
16. [三方服务集成](#十六三方服务集成)
17. [IP 黑名单](#十七ip-黑名单)
18. [省市区管理](#十八省市区管理)
19. [TCP 服务](#十九tcp-服务)
20. [国际化](#二十国际化)
21. [前端页面与组件](#二十一前端页面与组件)
22. [API 接口清单](#二十二api-接口清单)
23. [系统配置项详解](#二十三系统配置项详解)

---

## 一、认证与安全

### 1.1 登录方式

| 登录方式 | 接口 | 说明 |
|----------|------|------|
| 账号密码登录 | `POST /admin/site/accountLogin` | 支持图形验证码校验 |
| 手机号登录 | `POST /admin/site/mobileLogin` | 短信验证码校验 |
| 微信授权登录 | `/admin/wechat/...` | OAuth 网页授权回调 |

### 1.2 Token 机制

- **JWT 签发**：登录成功后签发 AccessToken
- **自动续期**：可配置刷新间隔（`refreshInterval`）与最大续期次数（`maxRefreshTimes`）
- **多端登录**：可配置允许/禁止多设备同时在线（`multiLogin`）
- **Token 过期**：可配置有效期秒数（`expires`）

### 1.3 验证码

- `GET /admin/site/captcha` — 获取图形验证码（防止暴力登录）

### 1.4 中间件鉴权体系

| 中间件 | 适用范围 | 说明 |
|--------|----------|------|
| Admin Auth | `/admin/*` | 验证 JWT + Casbin 路由权限 |
| API Auth | `/api/*` | 验证 JWT（轻量级） |
| Home Auth | `/home/*` | 前台用户认证 |
| WebSocket Auth | `/socket` | WebSocket 连接鉴权 |
| Develop Auth | 开发模式 | 特权中间件，跳过部分校验 |

### 1.5 其他安全接口

- `POST /admin/site/logout` — 注销登录，使 Token 失效
- `GET /admin/site/ping` — 心跳检测，验证会话有效性
- `POST /admin/site/register` — 账号注册

---

## 二、用户管理

### 2.1 用户 CRUD

| 接口 | 方法 | 说明 |
|------|------|------|
| `/admin/member/list` | GET | 用户列表（分页+筛选） |
| `/admin/member/view` | GET | 获取指定用户信息 |
| `/admin/member/edit` | POST | 新增 / 修改用户 |
| `/admin/member/delete` | POST | 删除用户 |
| `/admin/member/status` | POST | 启用 / 禁用用户 |

### 2.2 用户信息维护

| 接口 | 方法 | 说明 |
|------|------|------|
| `/admin/member/updateProfile` | POST | 更新基本资料 |
| `/admin/member/updatePwd` | POST | 修改登录密码 |
| `/admin/member/resetPwd` | POST | 重置密码（管理员操作） |
| `/admin/member/updateEmail` | POST | 换绑邮箱 |
| `/admin/member/updateMobile` | POST | 换绑手机号 |
| `/admin/member/updateCash` | POST | 修改提现信息 |

### 2.3 用户资产操作

| 接口 | 方法 | 说明 |
|------|------|------|
| `/admin/member/addBalance` | POST | 增加/扣减余额 |
| `/admin/member/addIntegral` | POST | 增加/扣减积分 |

### 2.4 用户辅助功能

| 功能 | 说明 |
|------|------|
| 用户树形关系 | `GenTree()` 生成用户上下级关系树 |
| 下级用户集合 | `GetLowerIds()` 递归获取所有下级用户 ID |
| 邀请码查用户 | `GetIdByCode()` 通过邀请码反查用户 ID |
| 用户唯一性校验 | `VerifyUnique()` 手机/邮箱/用户名唯一性验证 |
| 超管识别 | `VerifySuperId()` 验证是否为超级管理员 |
| 用户选项接口 | `GET /admin/member/option` 用于下拉选择 |
| 登录用户信息 | `GET /admin/member/info` 获取当前登录用户完整信息 |
| 登录统计 | `MemberLoginStat()` 用户登录次数/时间统计 |

---

## 三、组织架构

### 3.1 部门管理

| 接口 | 方法 | 说明 |
|------|------|------|
| `/admin/dept/list` | GET | 部门树形列表 |
| `/admin/dept/view` | GET | 获取指定部门信息 |
| `/admin/dept/edit` | POST | 新增 / 修改部门 |
| `/admin/dept/delete` | POST | 删除部门 |
| `/admin/dept/option` | GET | 部门下拉选项 |
| `/admin/dept/treeOption` | GET | 部门树形选项（级联选择器用） |

**特性：**
- 支持无限层级树形结构
- 部门唯一性校验（`VerifyUnique`）
- 排序支持（`MaxSort` 自动获取最大排序值）

### 3.2 岗位管理

| 接口 | 方法 | 说明 |
|------|------|------|
| `/admin/post/list` | GET | 岗位列表 |
| `/admin/post/view` | GET | 岗位详情 |
| `/admin/post/edit` | POST | 新增 / 修改岗位 |
| `/admin/post/delete` | POST | 删除岗位 |
| `/admin/post/option` | GET | 岗位下拉选项 |

---

## 四、权限管理

### 4.1 角色管理

| 接口 | 方法 | 说明 |
|------|------|------|
| `/admin/role/list` | GET | 角色列表 |
| `/admin/role/edit` | POST | 新增 / 修改角色 |
| `/admin/role/delete` | POST | 删除角色 |
| `/admin/role/dynamic` | GET | 获取动态路由（前端菜单渲染） |
| `/admin/role/getPermissions` | GET | 获取指定角色的菜单权限 |
| `/admin/role/updatePermissions` | POST | 编辑角色菜单权限 |
| `/admin/role/dataScope/select` | GET | 获取数据权限选项列表 |
| `/admin/role/dataScope/edit` | POST | 修改角色数据权限范围 |

**数据权限范围类型：**

| 类型 | 说明 |
|------|------|
| 全部数据 | 可查看所有数据 |
| 本人数据 | 仅能查看自己的数据 |
| 自定义部门 | 按选定部门范围限制数据 |

**其他功能：**
- 角色层级：`GetSubRoleIds()` 递归获取所有下级角色 ID
- 角色校验：`VerifyRoleId()` 验证角色 ID 合法性
- Casbin 策略刷新（`make refresh`）无需重启服务

### 4.2 菜单管理

| 接口 | 方法 | 说明 |
|------|------|------|
| `/admin/menu/list` | GET | 菜单树形列表 |
| `/admin/menu/edit` | POST | 新增 / 修改菜单 |
| `/admin/menu/delete` | POST | 删除菜单 |
| `/admin/menu/fastList` | GET | 快速菜单列表（轻量级） |

**菜单功能特性：**
- 支持树形多级菜单
- 菜单类型：目录 / 菜单 / 按钮
- 路由路径与组件路径配置
- 权限标识（`perms`）与按钮级权限
- 动态路由（`GetMenuList()`）：按用户角色动态返回可见菜单
- 细粒度权限（`LoginPermissions()`）：返回当前用户所有按钮权限标识

---

## 五、系统配置

### 5.1 参数配置

| 接口 | 方法 | 说明 |
|------|------|------|
| `/admin/config/list` | GET | 获取配置列表 |
| `/admin/config/edit` | POST | 修改配置 |
| `/admin/config/get` | GET | 获取指定配置 |
| `/admin/site/config` | GET | 获取前端所需公开配置 |
| `/admin/site/loginConfig` | GET | 获取登录页配置 |

**支持的配置组：**
- 登录配置（login）
- 微信配置（wechat）
- 支付配置（pay）
- 短信配置（sms）
- 地理配置（geo）
- 上传配置（upload）
- 邮件配置（smtp）
- Token 配置（token）

### 5.2 字典管理

**字典类型（DictType）：**

| 接口 | 方法 | 说明 |
|------|------|------|
| `/admin/dictType/list` | GET | 字典类型列表 |
| `/admin/dictType/edit` | POST | 新增 / 修改字典类型 |
| `/admin/dictType/delete` | POST | 删除字典类型 |
| `/admin/dictType/tree` | GET | 字典类型树 |

**字典数据（DictData）：**

| 接口 | 方法 | 说明 |
|------|------|------|
| `/admin/dictData/list` | GET | 字典数据列表 |
| `/admin/dictData/edit` | POST | 新增 / 修改字典数据 |
| `/admin/dictData/delete` | POST | 删除字典数据 |
| `/admin/dictData/option` | GET | 字典数据选项接口 |

**字典类型支持：**
- 枚举值字典（静态配置）
- 自定义方法字典（动态查询）
- 系统内置枚举注册（`dict_register_enums.go`）
- 函数式字典注册（`dict_register_func.go`）

---

## 六、日志管理

### 6.1 访问日志（HTTP Log）

| 接口 | 方法 | 说明 |
|------|------|------|
| `/admin/log/list` | GET | 访问日志列表 |
| `/admin/log/view` | GET | 日志详情 |
| `/admin/log/delete` | POST | 删除日志 |

**记录内容：** 请求路径、方法、状态码、请求耗时、客户端 IP、请求参数、响应内容
**配置项：** 可按模块（admin / api / default）和状态码过滤，支持队列异步写入

### 6.2 登录日志

| 接口 | 方法 | 说明 |
|------|------|------|
| `/admin/loginLog/list` | GET | 登录日志列表 |
| `/admin/loginLog/view` | GET | 日志详情 |
| `/admin/loginLog/delete` | POST | 删除日志 |

**记录内容：** 登录账号、登录 IP、IP 归属地、登录时间、登录结果、失败原因

### 6.3 服务异常日志

| 接口 | 方法 | 说明 |
|------|------|------|
| `/admin/serveLog/list` | GET | 服务日志列表 |
| `/admin/serveLog/view` | GET | 日志详情（含堆栈） |
| `/admin/serveLog/delete` | POST | 删除日志 |

**级别：** WARN / ERROR / FATAL / PANIC
**记录内容：** 错误信息、堆栈 trace、发生时间、请求上下文

### 6.4 短信发送日志

| 接口 | 方法 | 说明 |
|------|------|------|
| `/admin/smsLog/list` | GET | 短信日志列表 |
| `/admin/smsLog/view` | GET | 日志详情 |
| `/admin/smsLog/delete` | POST | 删除日志 |

**记录内容：** 接收手机号、发送模板、发送状态、发送时间、供应商响应

### 6.5 邮件发送日志

| 接口 | 方法 | 说明 |
|------|------|------|
| `/admin/emsLog/list` | GET | 邮件日志列表 |
| `/admin/emsLog/view` | GET | 日志详情 |
| `/admin/emsLog/delete` | POST | 删除日志 |

---

## 七、监控管理

### 7.1 在线用户

- `GET /admin/monitor/onlineList` — 在线用户列表
- `POST /admin/monitor/forceOffline` — 强制下线指定用户

**展示信息：** 用户名、登录 IP、登录时间、Token 过期时间、在线时长

### 7.2 服务器监控（实时）

**WebSocket 事件：**
- `admin/monitor/trends` — 实时推送：CPU、内存、磁盘、网络 IO 趋势数据
- `admin/monitor/runInfo` — 推送服务器运行信息（Go 版本、启动时间、Goroutine 数等）

**监控指标：**

| 指标 | 说明 |
|------|------|
| CPU 使用率 | 实时 CPU 占用百分比 |
| 内存使用 | 已用 / 总内存（MB） |
| 磁盘使用 | 各挂载点使用情况 |
| 网络 IO | 上行 / 下行速率（KB/s） |
| 负载均衡 | loadAvg（1m/5m/15m） |
| Goroutine 数 | 当前协程数量 |
| GC 统计 | GC 次数、暂停时间 |

### 7.3 服务许可证

- `/admin/serveLicense/...` — 服务许可证管理（查看/验证授权状态）

---

## 八、通知与公告

### 8.1 公告管理

| 接口 | 方法 | 说明 |
|------|------|------|
| `/admin/notice/list` | GET | 公告列表 |
| `/admin/notice/view` | GET | 公告详情 |
| `/admin/notice/edit` | POST | 新增 / 修改公告 |
| `/admin/notice/delete` | POST | 删除公告 |
| `/admin/notice/status` | POST | 发布 / 撤回公告 |
| `/admin/notice/pullMessages` | GET | 拉取我的消息列表 |
| `/admin/notice/readAll` | POST | 一键标记全部已读 |

### 8.2 实时推送

- 公告发布后通过 **WebSocket** 实时推送给所有在线用户
- 支持：全体通知 / 指定用户私信
- 消息状态：已读 / 未读
- 消息类型：系统公告 / 私信

---

## 九、支付与财务

### 9.1 支付网关

**支持渠道：**

| 渠道 | 场景 |
|------|------|
| 支付宝 | H5、PC 端自动适配 |
| 微信支付 | H5、JSAPI |
| QQ 支付 | 通用支付 |

**核心流程：**
- `POST /api/pay/notify/...` — 支付异步回调（供各支付平台回调）
- 支付订单创建 → 支付 → 回调通知 → 订单状态更新
- 退款：`POST /admin/pay/refund/...` — 创建退款申请、查看退款状态

### 9.2 充值订单

| 接口 | 方法 | 说明 |
|------|------|------|
| `/admin/order/list` | GET | 充值订单列表 |
| `/admin/order/view` | GET | 订单详情 |
| `/admin/order/delete` | POST | 删除订单 |

**订单状态：** 待支付 / 已完成 / 已取消 / 已退款

### 9.3 提现管理

| 接口 | 方法 | 说明 |
|------|------|------|
| `/admin/cash/list` | GET | 提现申请列表 |
| `/admin/cash/view` | GET | 提现申请详情 |
| `/admin/cash/edit` | POST | 审核 / 处理提现 |
| `/admin/cash/delete` | POST | 删除提现记录 |

### 9.4 资金变动日志（积分/余额）

| 接口 | 方法 | 说明 |
|------|------|------|
| `/admin/creditsLog/list` | GET | 资金变动列表 |
| `/admin/creditsLog/view` | GET | 变动详情 |
| `/admin/creditsLog/delete` | POST | 删除记录 |

**记录内容：** 变动类型（积分/余额）、变动数量、变动前/后余额、变动原因、关联订单

---

## 十、文件与附件

### 10.1 文件上传

| 接口 | 方法 | 说明 |
|------|------|------|
| `/admin/upload/file` | POST | 普通文件上传 |
| `/admin/upload/image` | POST | 图片上传 |
| `/admin/upload/chunk` | POST | 分片上传（大文件） |
| `/admin/upload/mergeChunk` | POST | 合并分片（断点续传） |

### 10.2 附件管理

| 接口 | 方法 | 说明 |
|------|------|------|
| `/admin/attachment/list` | GET | 附件列表 |
| `/admin/attachment/view` | GET | 附件详情 |
| `/admin/attachment/delete` | POST | 删除附件 |
| `/admin/attachment/clearKind` | POST | 批量清理指定类型附件 |
| `/admin/attachment/kindOption` | GET | 附件类型选项 |

### 10.3 存储驱动（一键切换）

| 驱动 | 标识 | 说明 |
|------|------|------|
| 本地存储 | `local` | 文件存放于服务器本地目录 |
| 阿里云 OSS | `aliyun` | Aliyun Object Storage |
| 腾讯云 COS | `tencent` | Tencent Cloud Object Storage |
| UCloud US3 | `ucloud` | UCloud 对象存储 |
| 七牛云 | `qiniu` | Qiniu Cloud Storage |
| MinIO | `minio` | 自建 S3 兼容对象存储 |

---

## 十一、定时任务

### 11.1 任务管理

| 接口 | 方法 | 说明 |
|------|------|------|
| `/admin/cron/list` | GET | 定时任务列表 |
| `/admin/cron/view` | GET | 任务详情 |
| `/admin/cron/edit` | POST | 新增 / 修改任务 |
| `/admin/cron/delete` | POST | 删除任务 |
| `/admin/cron/status` | POST | 启用 / 停止任务 |

### 11.2 任务分组管理

| 接口 | 方法 | 说明 |
|------|------|------|
| `/admin/cronGroup/list` | GET | 任务分组列表 |
| `/admin/cronGroup/edit` | POST | 新增 / 修改分组 |
| `/admin/cronGroup/delete` | POST | 删除分组 |

### 11.3 内置定时任务

| 任务名 | 说明 |
|--------|------|
| `close_order` | 关闭超期未支付订单（自动取消超过 1 天的待支付订单） |
| `test` | 测试任务（无参数，输出当前时间） |
| `test2` | 测试任务 2（带参数示例） |

**特性：**
- Cron 表达式配置（标准 5 段格式）
- 任务执行结果持久化记录
- 数据库驱动（任务配置存储在 DB 中，在线编辑后即时生效）

---

## 十二、消息队列

### 12.1 驱动支持

| 驱动 | 标识 | 适用场景 |
|------|------|----------|
| 本地磁盘 | `disk` | 单机部署、开发测试 |
| Redis | `redis` | 中小规模、简单可靠 |
| RocketMQ | `rocketmq` | 高吞吐量、金融级可靠 |
| Kafka | `kafka` | 大数据、日志流处理 |

### 12.2 内置队列 Topic

| Topic | 处理器 | 说明 |
|-------|--------|------|
| `QueueLoginLogTopic` | `LoginLog` | 异步写入登录日志 |
| `QueueServeLogTopic` | `ServeLog` | 异步写入服务异常日志 |
| `QueueSysLogTopic` | `SysLog` | 异步写入操作日志 |

### 12.3 队列配置项

```yaml
queue:
  switch: true          # 总开关
  driver: redis         # 驱动选择
  groupName: hotgo      # 消费者组名
  disk:
    path: ./data/queue  # 磁盘队列路径
    batchSize: 100       # 批处理大小
    batchTime: 3         # 批处理间隔（秒）
  redis:
    timeout: 60          # 超时时间（秒）
  rocketmq:
    nameSrvAdders: [...]
    brokerAddr: ...
    retry: 3
  kafka:
    address: ...
    version: ...
    multiConsumer: true
```

---

## 十三、WebSocket 实时通信

### 13.1 连接端点

| 端点 | 方法 | 说明 |
|------|------|------|
| `/socket` | GET（Upgrade） | WebSocket 连接入口 |
| `/socket/send` | POST | HTTP 方式发送 WS 消息 |

### 13.2 消息格式

```json
// 客户端 → 服务端（WRequest）
{
  "event": "ping",
  "data": {}
}

// 服务端 → 客户端（WResponse）
{
  "event": "ping",
  "data": {},
  "code": 0,
  "errorMsg": "",
  "timestamp": 1700000000
}
```

### 13.3 内置事件

| 事件名 | 方向 | 说明 |
|--------|------|------|
| `ping` | 双向 | 心跳保活 |
| `join` | C→S | 加入消息组 |
| `quit` | C→S | 退出消息组 |
| `admin/monitor/trends` | S→C | 推送服务器实时监控数据（CPU/内存/网络） |
| `admin/monitor/runInfo` | S→C | 推送服务器运行信息 |

### 13.4 集群支持

- 集群模式下通过 **Redis Pub/Sub** 实现多实例间消息同步
- 配置 `isCluster: true` 启用

---

## 十四、代码生成器

> 仅在开发模式（`mode: develop`）下可用，可通过 `allowedIPs` 配置白名单

### 14.1 生成模板类型

| 模板类型 | 说明 |
|----------|------|
| 普通 CRUD | 标准增删改查页面 |
| 树形表 | 带层级树形结构的 CRUD |
| 关联表 | 多表关联的 CRUD |
| Queue 消费者 | 消息队列处理器模板 |
| Cron 定时任务 | 定时任务处理器模板 |

### 14.2 生成产物一览

| 层次 | 文件位置 | 是否自动 |
|------|----------|----------|
| API 请求/响应结构体 | `server/api/admin/<feature>/` | 生成 |
| Controller 控制器 | `server/internal/controller/` | 生成 |
| Service 接口 | `server/internal/service/` | `gf gen service` |
| Logic 业务逻辑 | `server/internal/logic/` | 生成 |
| DAO 数据访问 | `server/internal/dao/` | `gf gen dao` |
| Model Input | `server/internal/model/input/` | 生成 |
| Vue 前端页面 | `web/src/views/` | 生成 |
| 前端 API 封装 | `web/src/api/` | 生成 |

### 14.3 代码生成接口

| 接口 | 方法 | 说明 |
|------|------|------|
| `/admin/genCodes/list` | GET | 生成配置列表 |
| `/admin/genCodes/view` | GET | 查看生成配置 |
| `/admin/genCodes/edit` | POST | 新增 / 修改生成配置 |
| `/admin/genCodes/delete` | POST | 删除生成配置 |
| `/admin/genCodes/preview` | GET | 预览生成代码 |
| `/admin/genCodes/generate` | POST | 执行代码生成 |

### 14.4 Addon 管理接口

| 接口 | 方法 | 说明 |
|------|------|------|
| `/admin/addons/list` | GET | 插件列表 |
| `/admin/addons/build` | POST | 构建插件脚手架 |
| `/admin/addons/install` | POST | 安装插件 |
| `/admin/addons/upgrade` | POST | 升级插件 |
| `/admin/addons/unInstall` | POST | 卸载插件 |

---

## 十五、插件系统

### 15.1 架构特性

- 插件存放于 `server/addons/<addon_name>/`，与主程序完全解耦
- 每个插件独立注册路由、菜单、定时任务、队列消费者、配置项
- 插件可独立部署，支持迁移到其他项目

### 15.2 插件目录结构

```
addons/<name>/
├── api/
│   ├── admin/          # 后台 API 结构体
│   ├── api/            # 开放 API 结构体
│   ├── home/           # 前台 API 结构体
│   └── websocket/      # WebSocket 消息结构体
├── controller/         # 控制器（按入口分类）
├── crons/              # 插件定时任务
├── global/             # 插件全局初始化
├── logic/              # 插件业务逻辑
├── model/              # 插件数据模型
├── router/             # 插件路由注册
└── manifest/config/    # 插件独立配置
```

### 15.3 内置示例插件（hgexample）

**后台功能模块：**

| 模块 | 说明 |
|------|------|
| 门户首页 | 插件首页展示 |
| 系统配置 | 插件独立配置管理 |
| 表格管理 | 普通列表表格示例 |
| 树形表格 | 树形数据表格示例 |
| 租户订单 | 多租户订单管理示例 |
| WebSocket 测试 | 实时通信测试界面 |
| 组件演示 | 丰富 UI 组件展示（见下表） |

**组件演示清单：**

| 类别 | 组件 |
|------|------|
| 基础交互 | 日历（Calendar）、拖拽（Drag）、时间线（Timeline）、水印（Watermark） |
| 图标 | Ant Design Icons、Ionicons5 图标选择器 |
| 数据图表 | ECharts、MadeaPie、PP Chart、瀑布图（Waterfall） |
| 地图 | 百度地图、高德地图（AMap） |
| 表单 | 复杂表单示例、Form Builder |
| 工具功能 | Excel 导入、打印（Print）、指纹识别（FingerprintJS） |
| 文字特效 | 渐变文字、高亮文字、拼音（Pinyin）、Mint 风格 |
| 加密 | DES 加密解密组件 |
| 弹窗通知 | Modal 弹窗、Notification 通知 |

---

## 十六、三方服务集成

### 16.1 短信服务（SMS）

| 供应商 | 标识 | 说明 |
|--------|------|------|
| 阿里云 | `aliyun` | 支持模板短信 |
| 腾讯云 | `tencent` | 支持模板短信 |

**接口：**
- `/admin/sms/send` — 发送短信验证码
- `/admin/sms/template/list` — 短信模板管理

### 16.2 邮件服务（EMS）

- 标准 SMTP 协议，支持 TLS
- `/admin/ems/send` — 发送邮件
- 支持 HTML 邮件模板

### 16.3 微信授权

- `/admin/wechat/auth` — 微信 OAuth 授权入口
- `/admin/wechat/callback` — 微信授权回调处理
- 获取用户 openid、nickname、头像等信息

### 16.4 IP 归属地解析

| 方式 | 说明 |
|------|------|
| `cz88` | 离线 IP 库解析 |
| `whois` | 在线 WHOIS 查询 |

配置项：`system.ipMethod: cz88 | whois`

---

## 十七、IP 黑名单

| 接口 | 方法 | 说明 |
|------|------|------|
| `/admin/blacklist/list` | GET | 黑名单列表 |
| `/admin/blacklist/view` | GET | 黑名单详情 |
| `/admin/blacklist/edit` | POST | 新增 / 修改黑名单 |
| `/admin/blacklist/delete` | POST | 删除黑名单 |
| `/admin/blacklist/status` | POST | 启用 / 禁用 |

**工作原理：** 预过滤中间件在请求进入路由前检查客户端 IP，命中黑名单则直接拦截返回 403。

---

## 十八、省市区管理

| 接口 | 方法 | 说明 |
|------|------|------|
| `/admin/provinces/list` | GET | 省市区列表 |
| `/admin/provinces/option` | GET | 省市区下拉选项 |
| `/admin/provinces/treeOption` | GET | 三级联动树形数据 |
| `/admin/provinces/getNameByCode` | GET | 通过 Code 查询名称 |

**数据：** 内置完整中国行政区划（省 / 市 / 区县）数据

---

## 十九、TCP 服务

- TCP Server / Client 双端支持
- 长连接管理与断线自动重连
- 服务间身份认证
- 路由消息分发机制
- RPC 请求/响应通信
- 消息拦截器（中间件链）
- 消息数据自动绑定（结构化）
- 独立日志输出（与 HTTP 日志分离）

---

## 二十、国际化

**支持语言：**

| 语言 | 标识 |
|------|------|
| 简体中文 | `zh-CN` |
| 繁体中文 | `zh-TW` |
| 英文 | `en` |

**配置：**
```yaml
system:
  i18n:
    switch: true
    defaultLanguage: zh-CN
```

**前端目录：** `web/src/locale/`（各语言 JSON 文件，运行时切换）

---

## 二十一、前端页面与组件

### 21.1 核心页面

| 页面 | 路径 | 说明 |
|------|------|------|
| 控制台 | `/dashboard/console` | 系统概览 Dashboard |
| 工作台 | `/dashboard/workplace` | 个人任务工作台 |
| 登录 | `/login` | 账号/手机/扫码登录 |
| 注册 | `/register` | 账号注册 |
| 个人中心 | `/account/...` | 资料、消息、安全设置 |
| 404 | `/exception/404` | 页面不存在 |

### 21.2 组织管理页面

| 页面 | 说明 |
|------|------|
| 用户管理 | 列表、新增、编辑、角色分配 |
| 部门管理 | 树形结构管理 |
| 岗位管理 | 岗位列表管理 |
| 角色管理 | 角色与菜单权限配置 |
| 菜单管理 | 路由与权限树 |

### 21.3 系统设置页面

| 页面 | 说明 |
|------|------|
| 参数配置 | 动态系统参数编辑 |
| 字典类型 | 字典类型管理 |
| 字典数据 | 字典值管理 |
| 附件管理 | 文件列表、预览、删除 |
| 省市区 | 行政区划数据管理 |
| 定时任务 | 在线 Cron 管理 |
| 任务分组 | Cron 分组管理 |
| 黑名单 | IP 黑名单管理 |
| 插件管理 | 安装/卸载/升级插件 |
| 代码生成 | 在线代码生成器 |

### 21.4 日志 & 监控页面

| 页面 | 说明 |
|------|------|
| 访问日志 | HTTP 请求日志查看 |
| 登录日志 | 用户登录历史 |
| 服务日志 | 异常错误日志 |
| 短信日志 | 短信发送历史 |
| 邮件日志 | 邮件发送历史 |
| 在线用户 | 当前在线用户列表 |
| 服务器监控 | 实时 CPU/内存/网络图表 |

### 21.5 财务页面

| 页面 | 说明 |
|------|------|
| 充值订单 | 订单列表与详情 |
| 提现申请 | 提现列表与审核 |
| 资金变动 | 积分/余额变动记录 |
| 交易退款 | 退款列表与状态 |

### 21.6 演示 & 示例页面

| 页面 | 说明 |
|------|------|
| CRUD 列表演示 | 标准增删改查示例 |
| 普通树表演示 | 树形表格展示 |
| 选项树表演示 | 可选树形结构 |
| 测试分类 | 层级分类测试 |

---

## 二十二、API 接口清单

### 路由前缀约定

| 类型 | 默认前缀 | 说明 |
|------|----------|------|
| 后台管理 | `/admin` | 需要 JWT + Casbin 鉴权 |
| 开放 API | `/api` | 需要 JWT 鉴权 |
| 前台页面 | `/home` | 需要前台鉴权 |
| WebSocket | `/socket` | WebSocket 专用 |

### 后台路由汇总（`/admin`）

| 模块 | 前缀 | 主要操作 |
|------|------|----------|
| 网站基础 | `/admin/site` | 登录、注销、配置、验证码 |
| 控制台 | `/admin/console` | 统计数据 |
| 用户管理 | `/admin/member` | CRUD + 资产操作 |
| 角色管理 | `/admin/role` | CRUD + 权限配置 |
| 部门管理 | `/admin/dept` | CRUD + 树形 |
| 菜单管理 | `/admin/menu` | CRUD + 动态路由 |
| 岗位管理 | `/admin/post` | CRUD |
| 公告管理 | `/admin/notice` | CRUD + 推送 |
| 监控 | `/admin/monitor` | 在线用户 + 服务器信息 |
| 系统配置 | `/admin/config` | 配置读写 |
| 字典类型 | `/admin/dictType` | CRUD |
| 字典数据 | `/admin/dictData` | CRUD |
| 附件管理 | `/admin/attachment` | CRUD + 上传 |
| 省市区 | `/admin/provinces` | 查询 + 联动 |
| 定时任务 | `/admin/cron` | CRUD + 状态 |
| 任务分组 | `/admin/cronGroup` | CRUD |
| 黑名单 | `/admin/blacklist` | CRUD |
| 访问日志 | `/admin/log` | 查询 + 删除 |
| 登录日志 | `/admin/loginLog` | 查询 + 删除 |
| 服务日志 | `/admin/serveLog` | 查询 + 删除 |
| 短信日志 | `/admin/smsLog` | 查询 + 删除 |
| 邮件日志 | `/admin/emsLog` | 查询 + 删除 |
| 充值订单 | `/admin/order` | 查询 + 删除 |
| 提现管理 | `/admin/cash` | 查询 + 审核 |
| 资金变动 | `/admin/creditsLog` | 查询 + 删除 |
| 退款管理 | `/admin/pay/refund` | 查询 + 创建 |
| 插件管理 | `/admin/addons` | 安装/卸载/升级 |
| 代码生成 | `/admin/genCodes` | 配置 + 生成 |
| 短信发送 | `/admin/sms` | 发送接口 |
| 邮件发送 | `/admin/ems` | 发送接口 |
| 文件上传 | `/admin/upload` | 普通 + 分片 |
| 微信 | `/admin/wechat` | OAuth 授权 |
| 服务许可 | `/admin/serveLicense` | 授权查询 |

### 开放 API 路由（`/api`）

| 接口 | 说明 |
|------|------|
| `POST /api/pay/notify/alipay` | 支付宝异步回调 |
| `POST /api/pay/notify/wechat` | 微信支付异步回调 |
| `GET/POST /api/member/...` | 前台用户相关接口 |

---

## 二十三、系统配置项详解

### 23.1 应用配置

```yaml
system:
  appName: HotGo            # 应用名称
  debug: true               # 调试模式
  mode: develop             # 运行模式: develop|testing|staging|product
  ipMethod: whois           # IP解析: cz88|whois
  isDemo: false             # 演示系统模式（禁止写操作）
  isCluster: false          # 集群模式（启用Redis Pub/Sub）
  addonsResourcePath: ""    # 插件资源目录
```

### 23.2 HTTP 服务器配置

```yaml
server:
  address: ":8000"
  openapiPath: "/api.json"
  swaggerPath: "/swagger"
  serverRoot: "./resource/public"
  DumpRouterMap: false
  maxHeaderBytes: "100KB"
  clientMaxBodySize: "200MB"
  ErrorStack: true
  ErrorLogEnabled: true
  accessLogEnabled: false
  pprofEnabled: false
  pprofPattern: "/debug/pprof"
```

### 23.3 日志配置

```yaml
system:
  log:
    switch: true             # 访问日志开关
    queue: false             # 是否通过队列异步写入
    module: [admin, api]     # 记录的模块
    skipCode: ["200"]        # 跳过的状态码
  serveLog:
    switch: true             # 服务日志开关
    queue: false             # 是否通过队列异步写入
    levelFormat: [WARN, ERROR, FATAL, PANIC]
```

### 23.4 Token 配置

```yaml
token:
  secretKey: "your-secret"
  expires: 86400            # 有效期（秒）
  autoRefresh: true         # 自动续期
  refreshInterval: 3600     # 每隔多久续期（秒）
  maxRefreshTimes: 30       # 最大续期次数
  multiLogin: true          # 允许多端同时登录
```

### 23.5 路由配置

```yaml
router:
  admin:
    prefix: "/admin"
    exceptLogin: ["/admin/login", "/admin/site/captcha"]
    exceptAuth: ["/admin/site/ping", "/admin/site/config"]
  api:
    prefix: "/api"
  websocket:
    prefix: "/socket"
  home:
    prefix: "/home"
```

### 23.6 缓存配置

```yaml
cache:
  adapter: memory           # memory|redis|file
  fileDir: "./data/cache"   # 文件缓存目录（adapter=file时有效）
```

### 23.7 代码生成配置

```yaml
hggen:
  allowedIPs: ["127.0.0.1"]   # IP白名单
  selectDbs: ["default"]       # 可选数据库
  disableTables:               # 禁止生成的表
    - sys_user
    - sys_role
  application:
    crud:
      templates:
        - name: default        # 模板名
          group: default       # 模板组
    queue:
      templates: [...]
    cron:
      templates: [...]
  addon:
    srcPath: "./addons"
    webApiPath: "../web/src/api"
    webViewsPath: "../web/src/views"
```

---

## 附录：命令速查

### 后端命令（`server/` 目录）

| 命令 | 说明 |
|------|------|
| `make all` | 同时启动 HTTP + Queue + Cron |
| `make http` | 仅启动 HTTP 服务 |
| `make queue` | 仅启动消息队列处理器 |
| `make cron` | 仅启动定时任务处理器 |
| `make build` | 编译前端并构建 Go 二进制 |
| `make dao` | 从 DB 重新生成 DAO / Entity |
| `make service` | 从 Logic 重新生成 Service 接口 |
| `make lint` | 运行 golangci-lint 代码检查 |
| `make refresh` | 刷新 Casbin 权限（无需重启） |

### 前端命令（`web/` 目录）

| 命令 | 说明 |
|------|------|
| `pnpm install` | 安装依赖 |
| `pnpm run dev` | 启动开发服务器（代理到 :8000） |
| `pnpm run build` | 构建生产包 → `server/resource/public/admin/` |
| `pnpm run lint:eslint` | ESLint 自动修复 |
| `pnpm run lint:prettier` | Prettier 格式化 |

---

*文档根据 HotGo v2.0 源代码自动分析生成，覆盖后端所有路由、逻辑层方法、配置项及前端页面。*
