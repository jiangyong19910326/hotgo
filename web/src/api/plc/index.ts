import { http } from '@/utils/http/axios';

// ─── 矿场 ────────────────────────────────────────────────────
export function MineList(params?: any) {
  return http.request({ url: '/plc/mine/list', method: 'GET', params });
}
export function MineView(params: any) {
  return http.request({ url: '/plc/mine/view', method: 'GET', params });
}
export function MineEdit(params: any) {
  return http.request({ url: '/plc/mine/edit', method: 'POST', params });
}
export function MineDelete(params: any) {
  return http.request({ url: '/plc/mine/delete', method: 'POST', params });
}
export function MineStatus(params: any) {
  return http.request({ url: '/plc/mine/status', method: 'POST', params });
}
export function MineOptions() {
  return http.request({ url: '/plc/mine/options', method: 'GET' });
}

// ─── 设备 ────────────────────────────────────────────────────
export function DeviceList(params?: any) {
  return http.request({ url: '/plc/device/list', method: 'GET', params });
}
export function DeviceView(params: any) {
  return http.request({ url: '/plc/device/view', method: 'GET', params });
}
export function DeviceEdit(params: any) {
  return http.request({ url: '/plc/device/edit', method: 'POST', params });
}
export function DeviceDelete(params: any) {
  return http.request({ url: '/plc/device/delete', method: 'POST', params });
}
export function DeviceStatus(params: any) {
  return http.request({ url: '/plc/device/status', method: 'POST', params });
}
export function DeviceControl(params: any) {
  return http.request({ url: '/plc/device/control', method: 'POST', params });
}

// ─── 数据点 ───────────────────────────────────────────────────
export function PointList(params?: any) {
  return http.request({ url: '/plc/point/list', method: 'GET', params });
}
export function PointView(params: any) {
  return http.request({ url: '/plc/point/view', method: 'GET', params });
}
export function PointEdit(params: any) {
  return http.request({ url: '/plc/point/edit', method: 'POST', params });
}
export function PointDelete(params: any) {
  return http.request({ url: '/plc/point/delete', method: 'POST', params });
}
export function PointStatus(params: any) {
  return http.request({ url: '/plc/point/status', method: 'POST', params });
}

// ─── 实时数据 ─────────────────────────────────────────────────
export function Realtime(params: any) {
  return http.request({ url: '/plc/realtime', method: 'GET', params });
}
// 看板汇总: 一次拉 设备 + 数据点 + 实时值
export function Overview(params: { deviceId: number }) {
  return http.request({ url: '/plc/overview', method: 'GET', params });
}

// 单点位历史时序 (画图用), 走前台 /api 路径不验签
export function PointHistory(params: {
  pointId: number;
  startTime?: string;
  endTime?: string;
  limit?: number;
}) {
  return http.request(
    { url: '/plc/history', method: 'GET', params },
    { urlPrefix: '/api' }
  );
}

// ─── 历史记录 ─────────────────────────────────────────────────
export function History(params: any) {
  return http.request({ url: '/plc/history', method: 'GET', params });
}

// ─── 报警 ─────────────────────────────────────────────────────
export function AlarmList(params?: any) {
  return http.request({ url: '/plc/alarm/list', method: 'GET', params });
}
export function AlarmResolve(params: any) {
  return http.request({ url: '/plc/alarm/resolve', method: 'POST', params });
}

// ─── 应用密钥 ─────────────────────────────────────────────────
export function AppList(params?: any) {
  return http.request({ url: '/plc/app/list', method: 'GET', params });
}
export function AppView(params: any) {
  return http.request({ url: '/plc/app/view', method: 'GET', params });
}
export function AppEdit(params: any) {
  return http.request({ url: '/plc/app/edit', method: 'POST', params });
}
export function AppDelete(params: any) {
  return http.request({ url: '/plc/app/delete', method: 'POST', params });
}
export function AppStatus(params: any) {
  return http.request({ url: '/plc/app/status', method: 'POST', params });
}
export function AppGenSecret() {
  return http.request({ url: '/plc/app/genSecret', method: 'GET' });
}
