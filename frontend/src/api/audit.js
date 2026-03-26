/**
 * 审计日志模块 API
 * 封装审计日志相关的所有 API 请求（仅管理员可访问）
 */

import request from '@/utils/request'

/**
 * 获取审计日志列表
 * 记录所有敏感操作，用于安全审计和合规性检查
 * @param {Object} [params] - 查询参数
 * @param {number} [params.skip] - 跳过记录数（分页）
 * @param {number} [params.limit] - 返回记录数限制
 * @param {string} [params.start_date] - 开始日期筛选
 * @param {string} [params.end_date] - 结束日期筛选
 * @param {string} [params.user_id] - 用户ID筛选
 * @param {string} [params.action] - 操作类型筛选
 * @returns {Promise} 返回审计日志列表
 */
export const getAuditLogs = (params) => {
  return request({
    url: '/audit-logs',
    method: 'get',
    params
  })
}

/**
 * 删除单条审计日志
 * @param {number} id - 审计日志ID
 * @returns {Promise} 返回删除结果
 */
export const deleteAuditLog = (id) => {
  return request({
    url: `/audit-logs/${id}`,
    method: 'delete'
  })
}

/**
 * 批量删除审计日志
 * @param {number[]} ids - 要删除的审计日志ID数组
 * @returns {Promise} 返回批量删除结果
 */
export const deleteAuditLogs = (ids) => {
  return request({
    url: '/audit-logs/batch-delete',
    method: 'post',
    data: { ids }
  })
}

/**
 * 导出审计日志
 * @param {Object} [params] - 导出参数（与查询参数相同）
 * @returns {Promise} 返回导出的文件或数据
 */
export const exportAuditLogs = (params) => {
  return request({
    url: '/audit-logs/export',
    method: 'get',
    params
  })
}
