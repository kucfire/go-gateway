import request from '@/utils/request'

export function login(data) {
  return request({
    url: '/admin_login/login',
    method: 'post',
    data
  })
}

export function getInfo(token) {
  return request({
    url: '/admin_info/info',
    method: 'get',
    params: { token }
  })
}

export function logout() {
  return request({
    url: '/admin_login/logout',
    method: 'get'
  })
}

export function changePwd(data) {
  return request({
    url: '/admin_info/changepwd',
    method: 'post',
    data
  })
}

