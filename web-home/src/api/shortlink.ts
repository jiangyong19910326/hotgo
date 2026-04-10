import { get, post } from '@/utils/http';

export interface ShortLink {
  id: number;
  code: string;
  original_url: string;
  short_url: string;
  title: string;
  total_clicks: number;
  today_clicks: number;
  expire_at: string | null;
  created_at: string;
}

export interface CreateParams {
  original_url: string;
  title?: string;
  expire_at?: string;
}

export interface CreateResult {
  code: string;
  short_url: string;
  original_url: string;
}

export interface ListParams {
  page?: number;
  limit?: number;
  keyword?: string;
}

export interface ListResult {
  list: ShortLink[];
  total: number;
  page: number;
  limit: number;
}

export interface StatsResult {
  code: string;
  original_url: string;
  short_url: string;
  title: string;
  total_clicks: number;
  today_clicks: number;
  yesterday_clicks: number;
  weekly_clicks: number;
  monthly_clicks: number;
  expire_at: string | null;
  created_at: string;
  trend: { date: string; clicks: number }[];
  referrers: { source: string; count: number }[];
  regions: { region: string; count: number }[];
}

export interface RedirectResult {
  original_url: string;
  code: string;
  title: string;
}

// 创建短链
export function createShortLink(params: CreateParams) {
  return post<{ data: CreateResult }>('/shortlink/create', params);
}

// 获取短链跳转目标（前端跳转模式）
export function getRedirectUrl(code: string) {
  return get<{ data: RedirectResult }>(`/shortlink/redirect/${code}`);
}

// 获取访问统计
export function getStats(code: string) {
  return get<{ data: StatsResult }>(`/shortlink/stats/${code}`);
}

// 获取短链列表
export function getList(params?: ListParams) {
  return get<{ data: ListResult }>('/shortlink/list', params);
}
