import request from '@/utils/request'

export const getCustomerList = (params) => {
  return request({
    url: '/customers',
    method: 'get',
    params
  })
}

export const getCustomerDetail = (id) => {
  return request({
    url: `/customers/${id}`,
    method: 'get'
  })
}

export const createCustomer = (data) => {
  return request({
    url: '/customers',
    method: 'post',
    data
  })
}

export const updateCustomer = (id, data) => {
  return request({
    url: `/customers/${id}`,
    method: 'put',
    data
  })
}

export const deleteCustomer = (id) => {
  return request({
    url: `/customers/${id}`,
    method: 'delete'
  })
}

export const getContractTypeList = (params) => {
  return request({
    url: '/contract-types',
    method: 'get',
    params
  })
}

export const createContractType = (data) => {
  return request({
    url: '/contract-types',
    method: 'post',
    data
  })
}

export const updateContractType = (id, data) => {
  return request({
    url: `/contract-types/${id}`,
    method: 'put',
    data
  })
}

export const deleteContractType = (id) => {
  return request({
    url: `/contract-types/${id}`,
    method: 'delete'
  })
}