/**
 * 客户管理模块 API
 * 封装客户和合同类型相关的所有 API 请求
 */

import request from '@/utils/request'

/**
 * 获取客户列表
 * @param {Object} [params] - 查询参数
 * @param {string} [params.keyword] - 关键词搜索
 * @param {string} [params.type] - 客户类型 (customer/supplier)
 * @param {number} [params.skip] - 跳过记录数（分页）
 * @param {number} [params.limit] - 返回记录数限制
 * @returns {Promise} 返回客户列表
 */
export const getCustomerList = (params) => {
  return request({
    url: '/customers',
    method: 'get',
    params
  })
}

/**
 * 获取客户详情
 * @param {number} id - 客户ID
 * @returns {Promise} 返回客户详细信息
 */
export const getCustomerDetail = (id) => {
  return request({
    url: `/customers/${id}`,
    method: 'get'
  })
}

/**
 * 创建客户
 * @param {Object} data - 客户数据
 * @param {string} data.name - 客户名称
 * @param {string} [data.code] - 客户编码
 * @param {string} [data.type] - 客户类型 (customer/supplier)
 * @param {string} [data.contact_person] - 联系人
 * @param {string} [data.contact_phone] - 联系电话
 * @param {string} [data.contact_email] - 联系邮箱
 * @param {string} [data.address] - 地址
 * @param {string} [data.credit_level] - 信用等级
 * @returns {Promise} 返回创建的客户信息
 */
export const createCustomer = (data) => {
  return request({
    url: '/customers',
    method: 'post',
    data
  })
}

/**
 * 更新客户信息
 * @param {number} id - 客户ID
 * @param {Object} data - 更新的客户数据
 * @returns {Promise} 返回更新结果
 */
export const updateCustomer = (id, data) => {
  return request({
    url: `/customers/${id}`,
    method: 'put',
    data
  })
}

/**
 * 删除客户
 * @param {number} id - 客户ID
 * @returns {Promise} 返回删除结果
 */
export const deleteCustomer = (id) => {
  return request({
    url: `/customers/${id}`,
    method: 'delete'
  })
}

/**
 * 获取合同类型列表
 * 用于合同分类管理
 * @param {Object} [params] - 查询参数
 * @param {number} [params.skip] - 跳过记录数（分页）
 * @param {number} [params.limit] - 返回记录数限制
 * @returns {Promise} 返回合同类型列表
 */
export const getContractTypeList = (params) => {
  return request({
    url: '/contract-types',
    method: 'get',
    params
  })
}

/**
 * 创建合同类型
 * @param {Object} data - 合同类型数据
 * @param {string} data.name - 类型名称
 * @param {string} [data.description] - 类型描述
 * @returns {Promise} 返回创建的合同类型
 */
export const createContractType = (data) => {
  return request({
    url: '/contract-types',
    method: 'post',
    data
  })
}

/**
 * 更新合同类型
 * @param {number} id - 合同类型ID
 * @param {Object} data - 更新的数据
 * @returns {Promise} 返回更新结果
 */
export const updateContractType = (id, data) => {
  return request({
    url: `/contract-types/${id}`,
    method: 'put',
    data
  })
}

/**
 * 删除合同类型
 * @param {number} id - 合同类型ID
 * @returns {Promise} 返回删除结果
 */
export const deleteContractType = (id) => {
  return request({
    url: `/contract-types/${id}`,
    method: 'delete'
  })
}
