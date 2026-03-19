import request from '@/utils/request'

export const getApprovalRecords = (contractId) => {
  return request({
    url: `/contracts/${contractId}/approvals`,
    method: 'get'
  })
}

export const getPendingApprovals = () => {
  return request({
    url: '/pending-approvals',
    method: 'get'
  })
}

export const createApproval = (data) => {
  return request({
    url: `/contracts/${data.contract_id}/approvals`,
    method: 'post',
    data
  })
}

export const updateApproval = (id, data) => {
  return request({
    url: `/approvals/${id}`,
    method: 'put',
    data
  })
}

export const getReminders = (contractId) => {
  return request({
    url: `/contracts/${contractId}/reminders`,
    method: 'get'
  })
}

export const createReminder = (data) => {
  return request({
    url: `/contracts/${data.contract_id}/reminders`,
    method: 'post',
    data
  })
}

export const sendReminder = (id) => {
  return request({
    url: `/reminders/${id}/send`,
    method: 'post'
  })
}

export const getExpiringContracts = (days = 30) => {
  return request({
    url: '/expiring-contracts',
    method: 'get',
    params: { days }
  })
}

export const getStatistics = () => {
  return request({
    url: '/statistics',
    method: 'get'
  })
}

export const getNotificationCounts = () => {
  return request({
    url: '/notifications/count',
    method: 'get'
  })
}