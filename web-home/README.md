# HotGo 短链前台

基于 Vue 3 + TypeScript + NaiveUI 的独立短链服务前台。

## 功能

- **创建短链** — 输入长网址生成短链接，支持自定义短码、设置过期时间
- **短链跳转** — 访问 `/r/:code` 自动重定向到原始地址
- **访问统计** — 展示点击量、访问趋势、来源分布、地区分布
- **短链列表** — 公开短链接列表，支持搜索和分页

## 开发

```bash
pnpm install
pnpm dev      # 端口 8002，代理 /api → http://localhost:8000
```

## 构建

```bash
pnpm build    # 输出到 dist/
```

## 对接的后端 API

前台调用 `server/api/api/shortlink/` 下的接口（需在后端实现）：

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/shortlink/create` | 创建短链 |
| GET  | `/api/shortlink/redirect/:code` | 获取跳转目标 URL |
| GET  | `/api/shortlink/stats/:code` | 获取访问统计 |
| GET  | `/api/shortlink/list` | 获取短链列表 |

### 请求/响应格式

遵循 HotGo 标准格式：

```json
{ "code": 0, "message": "ok", "data": { ... } }
```
