import { http } from '@/utils/http/axios';

export function FrontUserList(params?: any) {
  return http.request({ url: '/frontUser/list', method: 'GET', params });
}
export function FrontUserView(params: any) {
  return http.request({ url: '/frontUser/view', method: 'GET', params });
}
export function FrontUserEdit(params: any) {
  return http.request({ url: '/frontUser/edit', method: 'POST', params });
}
export function FrontUserDelete(params: any) {
  return http.request({ url: '/frontUser/delete', method: 'POST', params });
}
export function FrontUserStatus(params: any) {
  return http.request({ url: '/frontUser/status', method: 'POST', params });
}
export function FrontUserResetPwd(params: any) {
  return http.request({ url: '/frontUser/resetPwd', method: 'POST', params });
}
