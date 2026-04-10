# PLC 数据开放接口文档

> 版本：v1.0  
> 基础路径：`/open`  
> 鉴权方式：HMAC-SHA256 签名验签  
> 数据格式：JSON  
> 字符编码：UTF-8  

---

## 目录

1. [接口规范](#1-接口规范)
2. [签名算法](#2-签名算法)
3. [通用响应格式](#3-通用响应格式)
4. [接口列表](#4-接口列表)
   - [4.1 获取设备实时数据](#41-获取设备实时数据)
   - [4.2 获取数据点历史记录](#42-获取数据点历史记录)
   - [4.3 获取报警记录列表](#43-获取报警记录列表)
5. [错误码说明](#5-错误码说明)
6. [签名代码示例](#6-签名代码示例)
7. [AppID 管理](#7-appid-管理)

---

## 1. 接口规范

### 1.1 基础信息

| 项目 | 说明 |
|------|------|
| 请求协议 | HTTP / HTTPS |
| 请求方式 | GET |
| 数据格式 | JSON |
| 基础路径 | `http://{host}:{port}/open` |

### 1.2 鉴权方式

所有接口均采用 **HMAC-SHA256 签名验签**，无需 Cookie / Session / Bearer Token。  
每次请求的 Query String 中必须附带以下 4 个签名参数：

| 参数 | 类型 | 必须 | 示例 | 说明 |
|------|------|:----:|------|------|
| `appId` | string | ✅ | `plc_demo_app` | 应用标识，由管理员在后台配置 |
| `timestamp` | string | ✅ | `1712920800` | Unix 秒级时间戳，与服务器偏差须 ≤ 300 秒 |
| `nonce` | string | ✅ | `xK9mQzPr` | 随机字符串，建议 8~16 位，每次请求须不同 |
| `sign` | string | ✅ | `3A8FD1C2...` | HMAC-SHA256 签名值（大写十六进制）|

---

## 2. 签名算法

### 2.1 算法描述

```
步骤 1 — 拼接待签名字符串
  将 appId、nonce、timestamp 三个参数按参数名 ASCII 字母升序排列，
  拼接为 URL Query 格式的字符串：

  signStr = "appId={appId}&nonce={nonce}&timestamp={timestamp}"

步骤 2 — 计算 HMAC-SHA256
  以 appSecret 为密钥，对 signStr 进行 HMAC-SHA256 运算：

  sign = HMAC_SHA256(signStr, appSecret)

步骤 3 — 转换为大写十六进制
  将二进制结果转为十六进制字符串并转大写：

  sign = HEX(sign).toUpperCase()
```

### 2.2 示例计算

```
appId     = plc_demo_app
nonce     = xK9mQzPr
timestamp = 1712920800
appSecret = change_me_32chars_secret_key_here

待签名字符串：
  appId=plc_demo_app&nonce=xK9mQzPr&timestamp=1712920800

计算结果（示例）：
  sign = A3F8B2D1E5C4907...（64位大写HEX）
```

### 2.3 注意事项

- `sign` 本身不参与签名计算
- `timestamp` 为 Unix 秒级整数，与服务器时间偏差超过 **300 秒** 会被拒绝
- `nonce` 建议每次请求使用不同的随机字符串，防止重放攻击
- 签名参数通过 **Query String** 传递，与业务参数并列，无需放入请求体

---

## 3. 通用响应格式

### 3.1 响应结构

```json
{
  "code":      0,
  "message":   "",
  "timestamp": 1712920800,
  "traceId":   "abc123",
  "data":      {}
}
```

| 字段 | 类型 | 说明 |
|------|------|------|
| `code` | int | 0=成功，非0=失败 |
| `message` | string | 错误描述，成功时为空 |
| `timestamp` | int | 服务器响应时间戳 |
| `traceId` | string | 链路追踪 ID |
| `data` | object | 业务数据，失败时为 null |

### 3.2 成功示例

```json
{
  "code": 0,
  "message": "",
  "data": { ... }
}
```

### 3.3 失败示例

```json
{
  "code": 401,
  "message": "签名验证失败",
  "data": null
}
```

---

## 4. 接口列表

---

### 4.1 获取设备实时数据

获取指定 PLC 设备所有启用数据点的当前采集值及报警状态。

#### 请求信息

```
GET /open/plc/realtime
```

#### 请求参数

| 参数 | 类型 | 必须 | 默认 | 说明 |
|------|------|:----:|------|------|
| `deviceId` | int | ✅ | | PLC 设备 ID |
| `appId` | string | ✅ | | 签名参数 |
| `timestamp` | string | ✅ | | 签名参数 |
| `nonce` | string | ✅ | | 签名参数 |
| `sign` | string | ✅ | | 签名参数 |

#### 请求示例

```
GET /open/plc/realtime?deviceId=1&appId=plc_demo_app&timestamp=1712920800&nonce=xK9mQzPr&sign=A3F8B2D1E5C49079...
```

#### 响应参数

**data 字段说明：**

| 字段 | 类型 | 说明 |
|------|------|------|
| `deviceId` | int | 设备 ID |
| `points` | array | 数据点列表 |

**points 子项字段：**

| 字段 | 类型 | 说明 |
|------|------|------|
| `pointId` | int | 数据点 ID |
| `field` | string | 字段标识，如 `furnace_temp` |
| `name` | string | 点位名称，如 `炉温` |
| `engValue` | float | 换算后工程值 |
| `unit` | string | 单位，如 `℃` `bar` `rpm` |
| `alarmType` | int | 0=正常 1=超上限 2=超下限 |
| `alarmMax` | float\|null | 报警上限（仅超上限时返回）|
| `alarmMin` | float\|null | 报警下限（仅超下限时返回）|

#### 响应示例

```json
{
  "code": 0,
  "message": "",
  "data": {
    "deviceId": 1,
    "points": [
      {
        "pointId":   1,
        "field":     "furnace_temp",
        "name":      "炉温",
        "engValue":  823.5,
        "unit":      "℃",
        "alarmType": 0
      },
      {
        "pointId":   2,
        "field":     "pressure",
        "name":      "炉压",
        "engValue":  1.35,
        "unit":      "bar",
        "alarmType": 1,
        "alarmMax":  1.2
      },
      {
        "pointId":   3,
        "field":     "motor_speed",
        "name":      "电机转速",
        "engValue":  1480.0,
        "unit":      "rpm",
        "alarmType": 0
      },
      {
        "pointId":   4,
        "field":     "valve_open",
        "name":      "阀门状态",
        "engValue":  1,
        "unit":      "",
        "alarmType": 0
      }
    ]
  }
}
```

---

### 4.2 获取数据点历史记录

获取指定数据点在时间段内的历史采集记录，支持分页。

#### 请求信息

```
GET /open/plc/history
```

#### 请求参数

| 参数 | 类型 | 必须 | 默认 | 说明 |
|------|------|:----:|------|------|
| `pointId` | int | ✅ | | 数据点 ID |
| `startTime` | string | | | 开始时间，格式 `2026-04-10 00:00:00` |
| `endTime` | string | | | 结束时间，格式 `2026-04-10 23:59:59` |
| `page` | int | | `1` | 页码，从 1 开始 |
| `perPage` | int | | `200` | 每页数量，最大 `500` |
| `appId` | string | ✅ | | 签名参数 |
| `timestamp` | string | ✅ | | 签名参数 |
| `nonce` | string | ✅ | | 签名参数 |
| `sign` | string | ✅ | | 签名参数 |

#### 请求示例

```
GET /open/plc/history?pointId=1&startTime=2026-04-10+00:00:00&endTime=2026-04-10+23:59:59&page=1&perPage=500&appId=plc_demo_app&timestamp=1712920800&nonce=xK9mQzPr&sign=A3F8B2D1...
```

#### 响应参数

**data 字段说明：**

| 字段 | 类型 | 说明 |
|------|------|------|
| `totalCount` | int | 符合条件的总记录数 |
| `page` | int | 当前页码 |
| `perPage` | int | 每页数量 |
| `list` | array | 历史记录列表（按采集时间倒序）|

**list 子项字段：**

| 字段 | 类型 | 说明 |
|------|------|------|
| `collectedAt` | string | 采集时间，格式 `2026-04-10 12:00:01` |
| `engValue` | float\|null | 换算后工程值，读取失败时为 null |
| `rawValue` | string | 原始字节值（十六进制字符串）|

#### 响应示例

```json
{
  "code": 0,
  "message": "",
  "data": {
    "totalCount": 1440,
    "page":       1,
    "perPage":    500,
    "list": [
      {
        "collectedAt": "2026-04-10 23:59:58",
        "engValue":    825.3,
        "rawValue":    "42CE999A"
      },
      {
        "collectedAt": "2026-04-10 23:59:57",
        "engValue":    824.8,
        "rawValue":    "42CE6666"
      },
      {
        "collectedAt": "2026-04-10 12:00:01",
        "engValue":    null,
        "rawValue":    ""
      }
    ]
  }
}
```

> **说明：** `engValue` 为 null 表示该时刻采集失败（PLC 离线或读取超时），`rawValue` 同时为空字符串。

---

### 4.3 获取报警记录列表

获取 PLC 设备的历史报警记录，支持按处理状态、时间段筛选和分页。

#### 请求信息

```
GET /open/plc/alarm/list
```

#### 请求参数

| 参数 | 类型 | 必须 | 默认 | 说明 |
|------|------|:----:|------|------|
| `deviceId` | int | | | 设备 ID，不传则查询所有设备 |
| `isResolved` | int | | | 处理状态：`1`=已处理 `2`=未处理，不传则查全部 |
| `startTime` | string | | | 报警触发时间起，格式 `2026-04-10 00:00:00` |
| `endTime` | string | | | 报警触发时间止 |
| `page` | int | | `1` | 页码 |
| `perPage` | int | | `20` | 每页数量 |
| `appId` | string | ✅ | | 签名参数 |
| `timestamp` | string | ✅ | | 签名参数 |
| `nonce` | string | ✅ | | 签名参数 |
| `sign` | string | ✅ | | 签名参数 |

#### 请求示例

```
GET /open/plc/alarm/list?deviceId=1&isResolved=2&page=1&perPage=20&appId=plc_demo_app&timestamp=1712920800&nonce=xK9mQzPr&sign=A3F8B2D1...
```

#### 响应参数

**data 字段说明：**

| 字段 | 类型 | 说明 |
|------|------|------|
| `totalCount` | int | 符合条件的总记录数 |
| `page` | int | 当前页码 |
| `perPage` | int | 每页数量 |
| `list` | array | 报警记录列表 |

**list 子项字段：**

| 字段 | 类型 | 说明 |
|------|------|------|
| `id` | int | 报警记录 ID |
| `deviceId` | int | 所属设备 ID |
| `pointId` | int | 触发报警的数据点 ID |
| `pointName` | string | 数据点名称 |
| `engValue` | float | 触发报警时的工程值 |
| `alarmType` | int | 1=超上限 2=超下限 |
| `alarmMax` | float\|null | 报警上限 |
| `alarmMin` | float\|null | 报警下限 |
| `unit` | string | 单位 |
| `isResolved` | int | 1=已处理 2=未处理 |
| `resolvedAt` | string | 处理时间（未处理时为空）|
| `remark` | string | 处理备注 |
| `triggeredAt` | string | 报警触发时间 |

#### 响应示例

```json
{
  "code": 0,
  "message": "",
  "data": {
    "totalCount": 2,
    "page":    1,
    "perPage": 20,
    "list": [
      {
        "id":          102,
        "deviceId":    1,
        "pointId":     2,
        "pointName":   "炉压",
        "engValue":    1.35,
        "alarmType":   1,
        "alarmMax":    1.2,
        "alarmMin":    null,
        "unit":        "bar",
        "isResolved":  2,
        "resolvedAt":  "",
        "remark":      "",
        "triggeredAt": "2026-04-10 14:23:11"
      },
      {
        "id":          101,
        "deviceId":    1,
        "pointId":     1,
        "pointName":   "炉温",
        "engValue":    950.2,
        "alarmType":   1,
        "alarmMax":    900.0,
        "alarmMin":    null,
        "unit":        "℃",
        "isResolved":  1,
        "resolvedAt":  "2026-04-10 10:15:30",
        "remark":      "已检查，属正常工况波动",
        "triggeredAt": "2026-04-10 09:58:44"
      }
    ]
  }
}
```

---

## 5. 错误码说明

| HTTP 状态码 | code | message | 原因及处理方式 |
|:-----------:|:----:|---------|--------------|
| 200 | `0` | _(空)_ | 成功 |
| 200 | `401` | 缺少签名参数 appId/timestamp/nonce/sign | 请求 Query 中缺少必传签名参数 |
| 200 | `401` | 时间戳已过期，请同步服务器时间 | `timestamp` 与服务器时间偏差超过 300 秒，检查系统时间 |
| 200 | `401` | 无效的 appId: appId [xxx] 不存在或已禁用 | appId 错误或已在后台禁用 |
| 200 | `401` | 签名验证失败 | sign 计算结果不匹配，检查签名算法和 appSecret |
| 200 | `400` | 设备ID不能为空 | 业务必填参数缺失 |
| 200 | `400` | 点位ID不能为空 | 业务必填参数缺失 |
| 200 | 其他 | 具体描述 | 业务错误，见 message 字段 |

---

## 6. 签名代码示例

### Go

```go
package main

import (
    "crypto/hmac"
    "crypto/sha256"
    "encoding/hex"
    "fmt"
    "math/rand"
    "strings"
    "time"
)

func sign(appId, nonce, timestamp, appSecret string) string {
    signStr := "appId=" + appId + "&nonce=" + nonce + "&timestamp=" + timestamp
    mac := hmac.New(sha256.New, []byte(appSecret))
    mac.Write([]byte(signStr))
    return strings.ToUpper(hex.EncodeToString(mac.Sum(nil)))
}

func randomNonce(n int) string {
    const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
    b := make([]byte, n)
    for i := range b { b[i] = chars[rand.Intn(len(chars))] }
    return string(b)
}

func buildUrl(base, appId, appSecret string, params map[string]string) string {
    timestamp := fmt.Sprintf("%d", time.Now().Unix())
    nonce     := randomNonce(12)
    s         := sign(appId, nonce, timestamp, appSecret)

    query := fmt.Sprintf("appId=%s&timestamp=%s&nonce=%s&sign=%s", appId, timestamp, nonce, s)
    for k, v := range params {
        query += "&" + k + "=" + v
    }
    return base + "?" + query
}
```

### TypeScript（浏览器）

```typescript
async function sign(appId: string, nonce: string, timestamp: string, appSecret: string): Promise<string> {
  const signStr = `appId=${appId}&nonce=${nonce}&timestamp=${timestamp}`;
  const enc = new TextEncoder();
  const key = await crypto.subtle.importKey(
    'raw', enc.encode(appSecret),
    { name: 'HMAC', hash: 'SHA-256' }, false, ['sign']
  );
  const buf = await crypto.subtle.sign('HMAC', key, enc.encode(signStr));
  return Array.from(new Uint8Array(buf))
    .map(b => b.toString(16).padStart(2, '0')).join('').toUpperCase();
}

async function buildUrl(base: string, appId: string, appSecret: string, params: Record<string, any> = {}): Promise<string> {
  const timestamp = String(Math.floor(Date.now() / 1000));
  const nonce     = Math.random().toString(36).substring(2, 14);
  const s         = await sign(appId, nonce, timestamp, appSecret);

  const authQs = `appId=${appId}&timestamp=${timestamp}&nonce=${nonce}&sign=${s}`;
  const bizQs  = Object.entries(params)
    .filter(([, v]) => v !== undefined && v !== null)
    .map(([k, v]) => `${k}=${encodeURIComponent(v)}`).join('&');
  return `${base}?${authQs}${bizQs ? '&' + bizQs : ''}`;
}

// 使用示例
const url = await buildUrl('http://your-server:8000/open/plc/realtime', 'plc_demo_app', 'your_secret', { deviceId: 1 });
const res = await fetch(url).then(r => r.json());
console.log(res.data.points);
```

### Python

```python
import hmac, hashlib, time, random, string, requests

def sign(app_id: str, nonce: str, timestamp: str, app_secret: str) -> str:
    sign_str = f"appId={app_id}&nonce={nonce}&timestamp={timestamp}"
    return hmac.new(app_secret.encode(), sign_str.encode(), hashlib.sha256).hexdigest().upper()

def build_params(app_id: str, app_secret: str, **kwargs) -> dict:
    timestamp = str(int(time.time()))
    nonce     = ''.join(random.choices(string.ascii_letters + string.digits, k=12))
    s         = sign(app_id, nonce, timestamp, app_secret)
    return {"appId": app_id, "timestamp": timestamp, "nonce": nonce, "sign": s, **kwargs}

# 使用示例
BASE = "http://your-server:8000/open"
APP_ID     = "plc_demo_app"
APP_SECRET = "change_me_32chars_secret_key_here"

# 获取实时数据
params = build_params(APP_ID, APP_SECRET, deviceId=1)
resp = requests.get(f"{BASE}/plc/realtime", params=params)
data = resp.json()["data"]
for p in data["points"]:
    print(f"{p['name']}: {p['engValue']} {p['unit']}")

# 获取历史记录
params = build_params(APP_ID, APP_SECRET, pointId=1, startTime="2026-04-10 00:00:00", perPage=500)
resp = requests.get(f"{BASE}/plc/history", params=params)
rows = resp.json()["data"]["list"]

# 获取未处理报警
params = build_params(APP_ID, APP_SECRET, isResolved=2, perPage=100)
resp = requests.get(f"{BASE}/plc/alarm/list", params=params)
alarms = resp.json()["data"]["list"]
```

### Java

```java
import javax.crypto.Mac;
import javax.crypto.spec.SecretKeySpec;

public class PlcOpenApi {

    public static String sign(String appId, String nonce, String timestamp, String appSecret) throws Exception {
        String signStr = "appId=" + appId + "&nonce=" + nonce + "&timestamp=" + timestamp;
        Mac mac = Mac.getInstance("HmacSHA256");
        mac.init(new SecretKeySpec(appSecret.getBytes("UTF-8"), "HmacSHA256"));
        byte[] bytes = mac.doFinal(signStr.getBytes("UTF-8"));
        StringBuilder sb = new StringBuilder();
        for (byte b : bytes) sb.append(String.format("%02X", b));
        return sb.toString();
    }

    public static String buildUrl(String base, String appId, String appSecret, String bizParams) throws Exception {
        String timestamp = String.valueOf(System.currentTimeMillis() / 1000);
        String nonce     = java.util.UUID.randomUUID().toString().replace("-", "").substring(0, 12);
        String s         = sign(appId, nonce, timestamp, appSecret);
        String auth      = "appId=" + appId + "&timestamp=" + timestamp + "&nonce=" + nonce + "&sign=" + s;
        return base + "?" + auth + (bizParams != null && !bizParams.isEmpty() ? "&" + bizParams : "");
    }
}
```

---

## 7. AppID 管理

### 7.1 数据表结构

```sql
CREATE TABLE `hg_plc_app` (
  `id`         int unsigned NOT NULL AUTO_INCREMENT,
  `app_id`     varchar(32)  NOT NULL COMMENT 'AppID（唯一标识）',
  `app_secret` varchar(64)  NOT NULL COMMENT 'AppSecret（HMAC签名密钥）',
  `name`       varchar(64)  NOT NULL DEFAULT '' COMMENT '应用名称',
  `remark`     varchar(255) NOT NULL DEFAULT '' COMMENT '备注',
  `status`     tinyint(1)   NOT NULL DEFAULT 1  COMMENT '1启用 2禁用',
  `created_at` datetime     DEFAULT NULL,
  `updated_at` datetime     DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_app_id` (`app_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='PLC API应用密钥';
```

### 7.2 管理操作

```sql
-- 新增应用
INSERT INTO hg_plc_app (app_id, app_secret, name, remark, status)
VALUES ('your_app_id', 'your_secret_at_least_32_chars!!!', '接入方名称', '用途说明', 1);

-- 禁用应用（立即拒绝所有该 AppID 的请求，Redis 缓存最长 60 秒后生效）
UPDATE hg_plc_app SET status = 2 WHERE app_id = 'your_app_id';

-- 更换密钥
UPDATE hg_plc_app SET app_secret = 'new_secret_value' WHERE app_id = 'your_app_id';

-- 查看所有应用
SELECT id, app_id, name, status, created_at FROM hg_plc_app;
```

### 7.3 安全建议

- `appSecret` 建议长度 ≥ 32 位，包含大小写字母+数字+特殊符号
- 不同接入方使用不同的 `appId` / `appSecret`，便于独立禁用和审计
- 定期轮换 `appSecret`，更换后通知接入方同步更新
- 服务端对 `nonce` 可增加 Redis 去重（TTL=5分钟）以防重放攻击

---

*文档生成时间：2026-04-10*
