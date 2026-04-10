import { http, jumpExport } from '@/utils/http/axios';

// 获取短链接列表
export function List(params) {
  return http.request({
    url: '/shortLink/list',
    method: 'get',
    params,
  });
}

// 删除/批量删除短链接
export function Delete(params) {
  return http.request({
    url: '/shortLink/delete',
    method: 'POST',
    params,
  });
}

// 添加/编辑短链接
export function Edit(params) {
  return http.request({
    url: '/shortLink/edit',
    method: 'POST',
    params,
  });
}

// 修改短链接状态
export function Status(params) {
  return http.request({
    url: '/shortLink/status',
    method: 'POST',
    params,
  });
}

// 获取短链接指定详情
export function View(params) {
  return http.request({
    url: '/shortLink/view',
    method: 'GET',
    params,
  });
}

// 导出短链接
export function Export(params) {
  jumpExport('/shortLink/export', params);
}