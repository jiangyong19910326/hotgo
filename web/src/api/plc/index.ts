import { http } from '@/utils/http/axios';

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
