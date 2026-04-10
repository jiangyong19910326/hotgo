/**
 * PLC 开放接口客户端 — HMAC-SHA256 签名版
 *
 * 适用于：需要通过签名验签调用 /open/plc/* 接口的场景
 *
 * 使用方式：
 *   import { PlcOpenClient } from '@/api/plc/plc-open-api';
 *
 *   const client = new PlcOpenClient({
 *     baseUrl: 'http://your-server:8000',
 *     appId:   'plc_demo_app',
 *     appSecret: 'change_me_32chars_secret_key_here',
 *   });
 *
 *   const realtime = await client.getRealtime(1);
 */

// ─────────────────────────────────────────────────────────────────────────────
// 类型定义
// ─────────────────────────────────────────────────────────────────────────────

export interface PlcOpenClientOptions {
  /** 后端地址，不含末尾斜杠，如 'http://127.0.0.1:8000' */
  baseUrl: string;
  /** AppID */
  appId: string;
  /** AppSecret（HMAC 签名密钥） */
  appSecret: string;
}

export interface RealtimePoint {
  pointId: number;
  field: string;
  name: string;
  engValue: number;
  unit: string;
  alarmType: number;   // 0正常 1超上限 2超下限
  alarmMax?: number | null;
  alarmMin?: number | null;
}

export interface RealtimeRes {
  deviceId: number;
  points: RealtimePoint[];
}

export interface HistoryItem {
  collectedAt: string;
  engValue: number | null;
  rawValue: string;
}

export interface HistoryRes {
  totalCount: number;
  page: number;
  perPage: number;
  list: HistoryItem[];
}

export interface AlarmItem {
  id: number;
  deviceId: number;
  pointId: number;
  pointName: string;
  engValue: number;
  alarmType: number;
  alarmMin?: number | null;
  alarmMax?: number | null;
  unit: string;
  isResolved: number;  // 1已处理 2未处理
  resolvedAt?: string;
  remark: string;
  triggeredAt: string;
}

export interface AlarmListRes {
  totalCount: number;
  page: number;
  perPage: number;
  list: AlarmItem[];
}

// ─────────────────────────────────────────────────────────────────────────────
// HMAC-SHA256 签名（浏览器原生 SubtleCrypto）
// ─────────────────────────────────────────────────────────────────────────────

/**
 * 计算 HMAC-SHA256 并返回大写十六进制字符串
 *
 * 签名字符串格式（参数按字母顺序排列）：
 *   appId={appId}&nonce={nonce}&timestamp={timestamp}
 *
 * sign = HMAC-SHA256( signStr, appSecret )  → 大写 HEX
 */
async function computeSign(
  appId: string,
  nonce: string,
  timestamp: string,
  appSecret: string
): Promise<string> {
  const signStr = `appId=${appId}&nonce=${nonce}&timestamp=${timestamp}`;
  const enc = new TextEncoder();

  const key = await crypto.subtle.importKey(
    'raw',
    enc.encode(appSecret),
    { name: 'HMAC', hash: 'SHA-256' },
    false,
    ['sign']
  );

  const sigBuf = await crypto.subtle.sign('HMAC', key, enc.encode(signStr));
  const hexArr = Array.from(new Uint8Array(sigBuf));
  return hexArr.map((b) => b.toString(16).padStart(2, '0')).join('').toUpperCase();
}

/** 生成随机 nonce（12 位字母数字） */
function randomNonce(): string {
  const chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
  let s = '';
  for (let i = 0; i < 12; i++) s += chars[Math.floor(Math.random() * chars.length)];
  return s;
}

// ─────────────────────────────────────────────────────────────────────────────
// HTTP 核心
// ─────────────────────────────────────────────────────────────────────────────

interface ApiResponse<T = any> {
  code: number;
  message: string;
  data: T;
}

// ─────────────────────────────────────────────────────────────────────────────
// PlcOpenClient
// ─────────────────────────────────────────────────────────────────────────────

export class PlcOpenClient {
  private readonly baseUrl: string;
  private readonly appId: string;
  private readonly appSecret: string;

  constructor(opts: PlcOpenClientOptions) {
    this.baseUrl   = opts.baseUrl.replace(/\/$/, '');
    this.appId     = opts.appId;
    this.appSecret = opts.appSecret;
  }

  /** 构造带签名参数的 URL */
  private async signedUrl(path: string, params: Record<string, any> = {}): Promise<string> {
    const timestamp = String(Math.floor(Date.now() / 1000));
    const nonce     = randomNonce();
    const sign      = await computeSign(this.appId, nonce, timestamp, this.appSecret);

    const authParams: Record<string, string> = {
      appId: this.appId,
      timestamp,
      nonce,
      sign,
    };

    // 过滤掉 undefined / null / ''
    const allParams = { ...authParams, ...params };
    const qs = Object.entries(allParams)
      .filter(([, v]) => v !== undefined && v !== null && v !== '')
      .map(([k, v]) => `${encodeURIComponent(k)}=${encodeURIComponent(String(v))}`)
      .join('&');

    return `${this.baseUrl}/open${path}?${qs}`;
  }

  private async get<T>(path: string, params: Record<string, any> = {}): Promise<T> {
    const url = await this.signedUrl(path, params);
    const res = await fetch(url);
    if (!res.ok) throw new Error(`HTTP ${res.status}: ${res.statusText}`);
    const json: ApiResponse<T> = await res.json();
    if (json.code !== 0) throw new Error(json.message || `API error code: ${json.code}`);
    return json.data;
  }

  // ─── 实时数据 ──────────────────────────────────────────────────────────────

  /**
   * 获取设备所有启用数据点的当前实时值
   * @example
   * const data = await client.getRealtime(1);
   * data.points.forEach(p => console.log(p.name, p.engValue, p.unit));
   */
  async getRealtime(deviceId: number): Promise<RealtimeRes> {
    return this.get<RealtimeRes>('/plc/realtime', { deviceId });
  }

  /**
   * 轮询实时数据，返回停止函数
   * @example
   * const stop = client.pollRealtime(1, 3000, (d) => renderDashboard(d.points));
   * // 页面卸载时：stop();
   */
  pollRealtime(
    deviceId: number,
    intervalMs: number = 3000,
    onData: (data: RealtimeRes) => void,
    onError?: (err: Error) => void
  ): () => void {
    let running = true;
    let timer: ReturnType<typeof setInterval> | null = null;

    const tick = async () => {
      if (!running) return;
      try {
        const data = await this.getRealtime(deviceId);
        if (running) onData(data);
      } catch (err) {
        if (running && onError) onError(err as Error);
      }
    };

    tick();
    timer = setInterval(tick, intervalMs);

    return () => {
      running = false;
      if (timer !== null) clearInterval(timer);
    };
  }

  // ─── 历史数据 ──────────────────────────────────────────────────────────────

  /**
   * 获取数据点历史记录（分页）
   */
  async getHistory(opts: {
    pointId: number;
    startTime?: string;
    endTime?: string;
    page?: number;
    perPage?: number;
  }): Promise<HistoryRes> {
    return this.get<HistoryRes>('/plc/history', opts);
  }

  /**
   * 获取最近 N 条历史记录（适合绘图）
   * @example
   * const rows = await client.getRecentHistory(3, 300);
   */
  async getRecentHistory(pointId: number, count: number = 200): Promise<HistoryItem[]> {
    const res = await this.getHistory({ pointId, page: 1, perPage: count });
    return res.list;
  }

  // ─── 报警记录 ──────────────────────────────────────────────────────────────

  /**
   * 获取报警记录列表
   */
  async getAlarmList(opts: {
    deviceId?: number;
    isResolved?: number;
    startTime?: string;
    endTime?: string;
    page?: number;
    perPage?: number;
  } = {}): Promise<AlarmListRes> {
    return this.get<AlarmListRes>('/plc/alarm/list', opts);
  }

  /**
   * 获取未处理的报警（常用于大屏角标）
   * @example
   * const alarms = await client.getPendingAlarms(1);
   * badge.count = alarms.length;
   */
  async getPendingAlarms(deviceId?: number): Promise<AlarmItem[]> {
    const res = await this.getAlarmList({ deviceId, isResolved: 2, perPage: 100 });
    return res.list;
  }
}

// ─────────────────────────────────────────────────────────────────────────────
// Node.js 签名工具函数（服务端调用时使用，浏览器端无需此函数）
// ─────────────────────────────────────────────────────────────────────────────

/**
 * Node.js 版本签名（使用 crypto 模块）
 *
 * @example
 * import { createHmac } from 'crypto';
 * const sign = computeSignNode('myAppId', 'abc123xyz', '1712700000', 'my_secret');
 */
export function computeSignNode(
  appId: string,
  nonce: string,
  timestamp: string,
  appSecret: string
): string {
  // 此函数仅在 Node.js 环境下可用
  // eslint-disable-next-line @typescript-eslint/no-var-requires
  const { createHmac } = require('crypto');
  const signStr = `appId=${appId}&nonce=${nonce}&timestamp=${timestamp}`;
  return createHmac('sha256', appSecret).update(signStr).digest('hex').toUpperCase();
}
