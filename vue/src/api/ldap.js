import request from '@/utils/request'

export function get() {
  return request({
    url: `/system/ldap`,
    method: 'get',
  })
}

export function update(data) {
  return request({
    url: `/system/ldap`,
    method: 'put',
    data
  })
}
