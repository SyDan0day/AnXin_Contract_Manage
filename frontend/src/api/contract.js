import request from '@/utils/request'

export const getContractList = (params) => {
  return request({
    url: '/contracts',
    method: 'get',
    params
  })
}

export const getContractDetail = (id) => {
  return request({
    url: `/contracts/${id}`,
    method: 'get'
  })
}

export const createContract = (data) => {
  return request({
    url: '/contracts',
    method: 'post',
    data
  })
}

export const updateContract = (id, data) => {
  return request({
    url: `/contracts/${id}`,
    method: 'put',
    data
  })
}

export const deleteContract = (id) => {
  return request({
    url: `/contracts/${id}`,
    method: 'delete'
  })
}

export const getContractExecutions = (contractId) => {
  return request({
    url: `/contracts/${contractId}/executions`,
    method: 'get'
  })
}

export const createContractExecution = (data) => {
  return request({
    url: `/contracts/${data.contract_id}/executions`,
    method: 'post',
    data
  })
}

export const deleteExecution = (id) => {
  return request({
    url: `/executions/${id}`,
    method: 'delete'
  })
}

export const getContractDocuments = (contractId) => {
  return request({
    url: `/contracts/${contractId}/documents`,
    method: 'get'
  })
}

export const uploadDocument = (data) => {
  return request({
    url: `/contracts/${data.contract_id}/documents`,
    method: 'post',
    data
  })
}

export const deleteDocument = (id) => {
  return request({
    url: `/documents/${id}`,
    method: 'delete'
  })
}

export const getContractLifecycle = (contractId) => {
  return request({
    url: `/contracts/${contractId}/lifecycle`,
    method: 'get'
  })
}

export const updateContractStatus = (contractId, data) => {
  return request({
    url: `/contracts/${contractId}/status`,
    method: 'put',
    data
  })
}

export const archiveContract = (contractId) => {
  return request({
    url: `/contracts/${contractId}/archive`,
    method: 'post'
  })
}

export const requestStatusChange = (contractId, data) => {
  return request({
    url: `/contracts/${contractId}/status-change`,
    method: 'post',
    data
  })
}

export const getStatusChangeRequests = (contractId) => {
  return request({
    url: `/contracts/${contractId}/status-change`,
    method: 'get'
  })
}

export const getPendingStatusChangeApprovals = () => {
  return request({
    url: '/pending-status-changes',
    method: 'get'
  })
}

export const approveStatusChangeRequest = (requestId, data) => {
  return request({
    url: `/status-change-requests/${requestId}/approve`,
    method: 'post',
    data
  })
}

export const rejectStatusChangeRequest = (requestId, data) => {
  return request({
    url: `/status-change-requests/${requestId}/reject`,
    method: 'post',
    data
  })
}