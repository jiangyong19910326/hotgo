/**
 * PLC 数据采集 — 独立前端 API 模块
 *
 * 使用方式：
 *   1. 修改 BASE_URL 为实际后端地址
 *   2. 调用 setToken(token) 设置登录 token
 *   3. 直接调用各方法，均返回 Promise<T>
 *
 * 依赖：仅需浏览器原生 fetch（或自行替换为 axios）
 */

// ─────────────────────────────────────────────────────────────────────────────
// 配置
// ─────────────────────────────────────────────────────────────────────────────

let BASE_URL = '/admin'; // 后端地址前缀，生产环境改为 'http://your-server:8000/admin'
let AUTH_TOKEN = '';

/** 设置后端地址（默认 /admin，生产按需覆盖） */
export function setBaseUrl(url: string) {
  BASE_URL = url.replace(/\/$/, '');
}

/** 设置鉴权 token（登录后调用） */
export function setToken(token: string) {
  AUTH_TOKEN = token;
}

// ─────────────────────────────────────────────────────────────────────────────
// 类型定义
// ─────────────────────────────────────────────────────────────────────────────

export interface PageReq {
  page?: number;
  perPage?: number;
}

export interface PageRes {
  totalCount: number;
  page: number;
  perPage: number;
}

/** PLC 设备 */
export interface PlcDevice {
  id: number;
  name: string;
  host: string;
  port: number;
  rack: number;
  slot: number;
  intervalMs: number;
  remark: string;
  status: number; // 1启用 2禁用
  createdAt: string;
}

export interface PlcDeviceListReq extends PageReq {
  name?: string;
  status?: number;
}

export interface PlcDeviceListRes extends PageRes {
  list: PlcDevice[];
}

export interface PlcDeviceEditReq {
  id?: number;
  name: string;
  host: string;
  port: number;    // 默认 102
  rack?: number;   // 默认 0
  slot?: number;   // 默认 1（S7-200 SMART）
  intervalMs: number;
  remark?: string;
  status?: number;
}

/** 数据点 */
export interface PlcPoint {
  id: number;
  deviceId: number;
  name: string;
  field: string;       // 前端使用的字段标识，如 "furnace_temp"
  area: string;        // DB / M / I / Q / V
  dbNumber: number;
  byteOffset: number;
  bitOffset: number;
  dataType: string;    // Bool / Byte / Word / DWord / Int / DInt / Real / String
  scale: number;       // 换算系数
  offsetVal: number;   // 换算偏移
  unit: string;        // 单位，如 ℃ bar rpm
  alarmMin: number | null;
  alarmMax: number | null;
  remark: string;
  sort: number;
  status: number;
}

export interface PlcPointListReq extends PageReq {
  deviceId?: number;
  name?: string;
  status?: number;
}

export interface PlcPointListRes extends PageRes {
  list: PlcPoint[];
}

export interface PlcPointEditReq {
  id?: number;
  deviceId: number;
  name: string;
  field: string;
  area: string;
  dbNumber?: number;
  byteOffset: number;
  bitOffset?: number;
  dataType: string;
  scale?: number;
  offsetVal?: number;
  unit?: string;
  alarmMin?: number | null;
  alarmMax?: number | null;
  remark?: string;
  sort?: number;
  status?: number;
}

/** 实时数据 — 单个点位 */
export interface PlcRealtimeItem {
  pointId: number;
  field: string;
  name: string;
  engValue: number;    // 换算后工程值
  unit: string;
  alarmType: number;   // 0正常 1超上限 2超下限
  alarmMax?: number | null;
  alarmMin?: number | null;
}

/** 实时数据 — 整台设备 */
export interface PlcRealtimeRes {
  deviceId: number;
  points: PlcRealtimeItem[];
}

/** 历史记录单条 */
export interface PlcHistoryItem {
  collectedAt: string;
  engValue: number | null;
  rawValue: string;
}

export interface PlcHistoryListReq extends PageReq {
  pointId: number;
  startTime?: string; // 如 "2026-04-10 00:00:00"
  endTime?: string;
}

export interface PlcHistoryListRes extends PageRes {
  list: PlcHistoryItem[];
}

/** 报警记录 */
export interface PlcAlarm {
  id: number;
  deviceId: number;
  pointId: number;
  pointName: string;
  engValue: number;
  alarmType: number;   // 1超上限 2超下限
  alarmMin?: number | null;
  alarmMax?: number | null;
  unit: string;
  isResolved: number;  // 1已处理 2未处理
  resolvedAt?: string | null;
  remark: string;
  triggeredAt: string;
}

export interface PlcAlarmListReq extends PageReq {
  deviceId?: number;
  isResolved?: number; // 1已处理 2未处理
  startTime?: string;
  endTime?: string;
}

export interface PlcAlarmListRes extends PageRes {
  list: PlcAlarm[];
}

// ─────────────────────────────────────────────────────────────────────────────
// HTTP 核心（基于 fetch，可替换为 axios）
// ─────────────────────────────────────────────────────────────────────────────

interface ApiResponse<T = any> {
  code: number;
  message: string;
  data: T;
}

async function request<T>(
  method: 'GET' | 'POST',
  path: string,
  payload?: Record<string, any>
): Promise<T> {
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
  };
  if (AUTH_TOKEN) {
    headers['Authorization'] = `Bearer ${AUTH_TOKEN}`;
  }

  let url = `${BASE_URL}${path}`;
  let body: string | undefined;

  if (method === 'GET' && payload) {
    // 过滤掉 undefined / null 值
    const params = Object.entries(payload)
      .filter(([, v]) => v !== undefined && v !== null && v !== '')
      .map(([k, v]) => `${encodeURIComponent(k)}=${encodeURIComponent(String(v))}`)
      .join('&');
    if (params) url += `?${params}`;
  } else if (method === 'POST' && payload) {
    body = JSON.stringify(payload);
  }

  const res = await fetch(url, { method, headers, body });

  if (!res.ok) {
    throw new Error(`HTTP ${res.status}: ${res.statusText}`);
  }

  const json: ApiResponse<T> = await res.json();

  if (json.code !== 0) {
    throw new Error(json.message || `API error code: ${json.code}`);
  }

  return json.data;
}

// ─────────────────────────────────────────────────────────────────────────────
// 登录
// ─────────────────────────────────────────────────────────────────────────────

export interface LoginReq {
  username: string;
  password: string;
  code?: string;      // 验证码（如关闭验证码可不传）
}

export interface LoginRes {
  token: string;
}

/**
 * 登录，自动将返回的 token 写入模块（后续请求自动携带）
 */
export async function login(req: LoginReq): Promise<LoginRes> {
  const data = await request<LoginRes>('POST', '/passport/login', req as any);
  if (data.token) setToken(data.token);
  return data;
}

// ─────────────────────────────────────────────────────────────────────────────
// 设备管理
// ─────────────────────────────────────────────────────────────────────────────

/** 获取设备列表 */
export async function getDeviceList(req?: PlcDeviceListReq): Promise<PlcDeviceListRes> {
  return request<PlcDeviceListRes>('GET', '/plc/device/list', {
    page: 1,
    perPage: 20,
    ...req,
  });
}

/** 获取设备详情 */
export async function getDevice(id: number): Promise<PlcDevice> {
  return request<PlcDevice>('GET', '/plc/device/view', { id });
}

/** 新增或修改设备（有 id 为修改，无 id 为新增） */
export async function saveDevice(req: PlcDeviceEditReq): Promise<void> {
  return request<void>('POST', '/plc/device/edit', req as any);
}

/** 删除设备 */
export async function deleteDevice(id: number | number[]): Promise<void> {
  return request<void>('POST', '/plc/device/delete', { id });
}

/** 更新设备状态（status: 1启用 2禁用） */
export async function updateDeviceStatus(id: number, status: 1 | 2): Promise<void> {
  return request<void>('POST', '/plc/device/status', { id, status });
}

// ─────────────────────────────────────────────────────────────────────────────
// 数据点管理
// ─────────────────────────────────────────────────────────────────────────────

/** 获取数据点列表 */
export async function getPointList(req?: PlcPointListReq): Promise<PlcPointListRes> {
  return request<PlcPointListRes>('GET', '/plc/point/list', {
    page: 1,
    perPage: 100,
    ...req,
  });
}

/** 获取数据点详情 */
export async function getPoint(id: number): Promise<PlcPoint> {
  return request<PlcPoint>('GET', '/plc/point/view', { id });
}

/** 新增或修改数据点 */
export async function savePoint(req: PlcPointEditReq): Promise<void> {
  return request<void>('POST', '/plc/point/edit', req as any);
}

/** 删除数据点 */
export async function deletePoint(id: number | number[]): Promise<void> {
  return request<void>('POST', '/plc/point/delete', { id });
}

/** 更新数据点状态 */
export async function updatePointStatus(id: number, status: 1 | 2): Promise<void> {
  return request<void>('POST', '/plc/point/status', { id, status });
}

// ─────────────────────────────────────────────────────────────────────────────
// ★ 实时数据（前端展示最常用）
// ─────────────────────────────────────────────────────────────────────────────

/**
 * 获取指定设备所有启用数据点的当前实时值
 *
 * @example
 * const data = await getRealtime(1);
 * data.points.forEach(p => {
 *   console.log(`${p.name}: ${p.engValue} ${p.unit}  ${p.alarmType > 0 ? '⚠报警' : ''}`);
 * });
 */
export async function getRealtime(deviceId: number): Promise<PlcRealtimeRes> {
  return request<PlcRealtimeRes>('GET', '/plc/realtime', { deviceId });
}

/**
 * 轮询实时数据，返回停止函数
 *
 * @param deviceId  设备ID
 * @param intervalMs  轮询间隔（毫秒），默认 3000
 * @param onData   每次拿到数据时的回调
 * @param onError  出错时的回调（不传则静默忽略）
 * @returns stop 函数，调用后停止轮询
 *
 * @example
 * const stop = pollRealtime(1, 2000, (data) => {
 *   renderDashboard(data.points);
 * });
 * // 页面卸载时：
 * stop();
 */
export function pollRealtime(
  deviceId: number,
  intervalMs: number = 3000,
  onData: (data: PlcRealtimeRes) => void,
  onError?: (err: Error) => void
): () => void {
  let timer: ReturnType<typeof setInterval> | null = null;
  let running = true;

  const fetch = async () => {
    if (!running) return;
    try {
      const data = await getRealtime(deviceId);
      if (running) onData(data);
    } catch (err) {
      if (running && onError) onError(err as Error);
    }
  };

  fetch(); // 立即执行一次
  timer = setInterval(fetch, intervalMs);

  return () => {
    running = false;
    if (timer !== null) clearInterval(timer);
  };
}

// ─────────────────────────────────────────────────────────────────────────────
// 历史记录
// ─────────────────────────────────────────────────────────────────────────────

/** 获取指定数据点的历史记录（分页） */
export async function getHistoryList(req: PlcHistoryListReq): Promise<PlcHistoryListRes> {
  return request<PlcHistoryListRes>('GET', '/plc/history', {
    page: 1,
    perPage: 500,
    ...req,
  });
}

/**
 * 获取数据点最近 N 条历史记录（不分页，直接返回数组）
 *
 * @example
 * const records = await getRecentHistory(1, 200);
 * // 用于绘制折线图
 */
export async function getRecentHistory(
  pointId: number,
  count: number = 200
): Promise<PlcHistoryItem[]> {
  const res = await getHistoryList({ pointId, page: 1, perPage: count });
  return res.list;
}

// ─────────────────────────────────────────────────────────────────────────────
// 报警管理
// ─────────────────────────────────────────────────────────────────────────────

/** 获取报警列表 */
export async function getAlarmList(req?: PlcAlarmListReq): Promise<PlcAlarmListRes> {
  return request<PlcAlarmListRes>('GET', '/plc/alarm/list', {
    page: 1,
    perPage: 20,
    ...req,
  });
}

/**
 * 获取未处理的报警（常用于大屏角标/提醒）
 *
 * @example
 * const alarms = await getPendingAlarms(1);
 * badge.count = alarms.length;
 */
export async function getPendingAlarms(deviceId?: number): Promise<PlcAlarm[]> {
  const res = await getAlarmList({ deviceId, isResolved: 2, perPage: 100 });
  return res.list;
}

/** 标记报警已处理 */
export async function resolveAlarm(id: number, remark?: string): Promise<void> {
  return request<void>('POST', '/plc/alarm/resolve', { id, remark: remark ?? '' });
}
