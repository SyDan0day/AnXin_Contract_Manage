/**
 * 认证模块 API
 * 封装用户认证相关的所有 API 请求，包括登录、注册、用户管理等
 */

import request from '@/utils/request'

/**
 * 用户登录
 * @param {Object} data - 登录数据
 * @param {string} data.username - 用户名
 * @param {string} data.password - 密码
 * @returns {Promise} 返回包含 access_token 和用户信息的响应
 */
export const login = (data) => {
  return request({
    url: '/auth/login',
    method: 'post',
    data
  })
}

/**
 * 用户注册
 * @param {Object} data - 注册数据
 * @param {string} data.username - 用户名
 * @param {string} data.password - 密码
 * @param {string} data.email - 邮箱
 * @param {string} [data.full_name] - 真实姓名
 * @returns {Promise} 返回注册结果
 */
export const register = (data) => {
  return request({
    url: '/auth/register',
    method: 'post',
    data
  })
}

/**
 * 获取用户列表
 * @param {Object} [params] - 查询参数
 * @param {number} [params.skip] - 跳过记录数（分页）
 * @param {number} [params.limit] - 返回记录数限制
 * @returns {Promise} 返回用户列表
 */
export const getUserList = (params) => {
  return request({
    url: '/auth/users',
    method: 'get',
    params
  })
}

/**
 * 获取用户详情
 * @param {number} id - 用户ID
 * @returns {Promise} 返回用户详细信息
 */
export const getUserDetail = (id) => {
  return request({
    url: `/auth/users/${id}`,
    method: 'get'
  })
}

/**
 * 更新用户信息
 * @param {number} id - 用户ID
 * @param {Object} data - 更新的用户数据
 * @returns {Promise} 返回更新结果
 */
export const updateUser = (id, data) => {
  return request({
    url: `/auth/users/${id}`,
    method: 'put',
    data
  })
}

/**
 * 删除用户
 * @param {number} id - 用户ID
 * @returns {Promise} 返回删除结果
 */
export const deleteUser = (id) => {
  return request({
    url: `/auth/users/${id}`,
    method: 'delete'
  })
}

/**
 * 重置用户密码
 * @param {number} id - 用户ID
 * @returns {Promise} 返回重置结果（包含新密码）
 */
export const resetPassword = (id) => {
  return request({
    url: `/auth/users/${id}/reset-password`,
    method: 'post'
  })
}
