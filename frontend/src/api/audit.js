import request from '@/utils/request'

export const getAuditLogs = (params) => {
  return request({
    url: '/audit-logs',
    method: 'get',
    params
  })
}

export const deleteAuditLog = (id) => {
  return request({
    url: `/audit-logs/${id}`,
    method: 'delete'
  })
}

export const deleteAuditLogs = (ids) => {
  return request({
    url: '/audit-logs/batch-delete',
    method: 'post',
    data: { ids }
  })
}

export const exportAuditLogs = (params) => {
  return request({
    url: '/audit-logs/export',
    method: 'get',
    params
  })
}
