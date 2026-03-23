package main

import (
	"contract-manage/config"
	"contract-manage/handlers"
	"contract-manage/middleware"
	"contract-manage/models"
	"fmt"
	"html/template"
	"time"

	"github.com/gin-gonic/gin"
)

var apiDocs = `
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>API 调试工具 - 安信合同管理系统</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        :root {
            --primary: #6366F1;
            --primary-dark: #4F46E5;
            --primary-light: #818CF8;
            --success: #10B981;
            --warning: #F59E0B;
            --danger: #EF4444;
            --info: #3B82F6;
            --bg-main: #F8FAFC;
            --bg-card: #FFFFFF;
            --bg-sidebar: #F1F5F9;
            --text-primary: #1E293B;
            --text-secondary: #64748B;
            --border: #E2E8F0;
            --shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
            --shadow-lg: 0 10px 15px -3px rgba(0, 0, 0, 0.1);
        }
        
        body {
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
            background: var(--bg-main);
            color: var(--text-primary);
            line-height: 1.5;
            min-height: 100vh;
        }
        
        /* 头部 */
        .header {
            background: var(--bg-card);
            border-bottom: 1px solid var(--border);
            padding: 12px 24px;
            position: sticky;
            top: 0;
            z-index: 100;
            box-shadow: var(--shadow);
        }
        
        .header-content {
            max-width: 1600px;
            margin: 0 auto;
            display: flex;
            align-items: center;
            justify-content: space-between;
        }
        
        .logo {
            display: flex;
            align-items: center;
            gap: 12px;
        }
        
        .logo-icon {
            width: 40px;
            height: 40px;
            background: linear-gradient(135deg, var(--primary), var(--primary-dark));
            border-radius: 10px;
            display: flex;
            align-items: center;
            justify-content: center;
            color: white;
            font-size: 20px;
        }
        
        .logo h1 {
            font-size: 20px;
            font-weight: 600;
            color: var(--text-primary);
        }
        
        .logo-subtitle {
            font-size: 12px;
            color: var(--text-secondary);
            margin-top: 2px;
        }
        
        .header-info {
            display: flex;
            align-items: center;
            gap: 20px;
        }
        
        .version-badge {
            padding: 4px 12px;
            background: linear-gradient(135deg, var(--primary), var(--primary-dark));
            color: white;
            border-radius: 20px;
            font-size: 12px;
            font-weight: 500;
        }
        
        .status-indicator {
            display: flex;
            align-items: center;
            gap: 8px;
            padding: 6px 14px;
            background: rgba(16, 185, 129, 0.1);
            border-radius: 20px;
            font-size: 12px;
            color: var(--success);
        }
        
        .status-dot {
            width: 8px;
            height: 8px;
            background: var(--success);
            border-radius: 50%;
            animation: pulse 2s infinite;
        }
        
        @keyframes pulse {
            0%, 100% { opacity: 1; transform: scale(1); }
            50% { opacity: 0.5; transform: scale(1.1); }
        }
        
        .time-display {
            font-size: 13px;
            color: var(--text-secondary);
            font-family: monospace;
        }
        
        /* 主容器 */
        .main-container {
            max-width: 1600px;
            margin: 0 auto;
            padding: 24px;
            display: grid;
            grid-template-columns: 320px 1fr;
            gap: 24px;
        }
        
        /* 侧边栏 */
        .sidebar {
            display: flex;
            flex-direction: column;
            gap: 16px;
        }
        
        .sidebar-card {
            background: var(--bg-card);
            border-radius: 16px;
            padding: 20px;
            box-shadow: var(--shadow);
            border: 1px solid var(--border);
        }
        
        .sidebar-header {
            display: flex;
            align-items: center;
            justify-content: space-between;
            margin-bottom: 16px;
        }
        
        .sidebar-title {
            font-size: 14px;
            font-weight: 600;
            color: var(--text-primary);
            display: flex;
            align-items: center;
            gap: 8px;
        }
        
        .search-box {
            width: 100%;
            padding: 10px 14px;
            background: var(--bg-main);
            border: 1px solid var(--border);
            border-radius: 10px;
            font-size: 13px;
            color: var(--text-primary);
            margin-bottom: 16px;
        }
        
        .search-box:focus {
            outline: none;
            border-color: var(--primary);
            box-shadow: 0 0 0 3px rgba(99, 102, 241, 0.1);
        }
        
        .api-group {
            margin-bottom: 20px;
        }
        
        .group-title {
            font-size: 11px;
            font-weight: 600;
            color: var(--primary);
            text-transform: uppercase;
            letter-spacing: 0.5px;
            margin-bottom: 10px;
            padding: 8px 12px;
            background: linear-gradient(135deg, rgba(99, 102, 241, 0.1), rgba(139, 92, 246, 0.1));
            border-radius: 8px;
            display: flex;
            align-items: center;
            gap: 8px;
        }
        
        .api-item {
            display: flex;
            align-items: center;
            gap: 12px;
            padding: 12px 14px;
            border-radius: 10px;
            cursor: pointer;
            transition: all 0.2s;
            margin-bottom: 6px;
            border: 1px solid transparent;
        }
        
        .api-item:hover {
            background: var(--bg-main);
            border-color: var(--border);
        }
        
        .api-item.active {
            background: linear-gradient(135deg, rgba(99, 102, 241, 0.15), rgba(139, 92, 246, 0.15));
            border-color: var(--primary);
        }
        
        .method-badge {
            font-size: 10px;
            font-weight: 700;
            padding: 4px 10px;
            border-radius: 6px;
            min-width: 55px;
            text-align: center;
            text-transform: uppercase;
        }
        
        .method-get { background: #DCFCE7; color: #16A34A; }
        .method-post { background: #DBEAFE; color: #2563EB; }
        .method-put { background: #FEF3C7; color: #D97706; }
        .method-delete { background: #FEE2E2; color: #DC2626; }
        
        .api-path {
            font-size: 13px;
            color: var(--text-primary);
            font-family: monospace;
            flex: 1;
            overflow: hidden;
            text-overflow: ellipsis;
            white-space: nowrap;
        }
        
        .api-badge {
            font-size: 10px;
            padding: 2px 6px;
            border-radius: 4px;
            background: var(--bg-main);
            color: var(--text-secondary);
        }
        
        /* 内容区 */
        .content {
            display: flex;
            flex-direction: column;
            gap: 20px;
        }
        
        .request-panel {
            background: var(--bg-card);
            border-radius: 16px;
            padding: 24px;
            box-shadow: var(--shadow);
            border: 1px solid var(--border);
        }
        
        .url-bar {
            display: flex;
            gap: 12px;
            margin-bottom: 20px;
        }
        
        .method-select {
            background: var(--bg-main);
            border: 1px solid var(--border);
            border-radius: 12px;
            padding: 14px 18px;
            font-size: 14px;
            font-weight: 600;
            color: var(--text-primary);
            cursor: pointer;
            min-width: 120px;
        }
        
        .method-select:focus {
            outline: none;
            border-color: var(--primary);
        }
        
        .url-input {
            flex: 1;
            background: var(--bg-main);
            border: 1px solid var(--border);
            border-radius: 12px;
            padding: 14px 18px;
            font-size: 14px;
            font-family: monospace;
            color: var(--text-primary);
        }
        
        .url-input:focus {
            outline: none;
            border-color: var(--primary);
            box-shadow: 0 0 0 3px rgba(99, 102, 241, 0.1);
        }
        
        .send-btn {
            background: linear-gradient(135deg, var(--primary), var(--primary-dark));
            color: white;
            border: none;
            border-radius: 12px;
            padding: 14px 32px;
            font-size: 14px;
            font-weight: 600;
            cursor: pointer;
            transition: all 0.2s;
            display: flex;
            align-items: center;
            gap: 8px;
        }
        
        .send-btn:hover {
            transform: translateY(-2px);
            box-shadow: 0 8px 20px rgba(99, 102, 241, 0.3);
        }
        
        .send-btn:active {
            transform: translateY(0);
        }
        
        .send-btn.loading {
            opacity: 0.7;
            pointer-events: none;
        }
        
        .token-section {
            background: linear-gradient(135deg, #F8FAFC, #F1F5F9);
            border-radius: 12px;
            padding: 16px;
            margin-bottom: 20px;
            border: 1px solid var(--border);
        }
        
        .token-header {
            display: flex;
            align-items: center;
            justify-content: space-between;
            margin-bottom: 10px;
        }
        
        .token-label {
            font-size: 13px;
            font-weight: 600;
            color: var(--text-primary);
            display: flex;
            align-items: center;
            gap: 8px;
        }
        
        .token-status {
            font-size: 11px;
            padding: 4px 10px;
            border-radius: 12px;
            background: #DCFCE7;
            color: #16A34A;
        }
        
        .token-status.expired {
            background: #FEE2E2;
            color: #DC2626;
        }
        
        .token-input {
            width: 100%;
            background: white;
            border: 1px solid var(--border);
            border-radius: 8px;
            padding: 12px 14px;
            font-size: 12px;
            font-family: monospace;
            color: var(--text-primary);
        }
        
        .token-input:focus {
            outline: none;
            border-color: var(--primary);
        }
        
        .quick-actions {
            display: flex;
            gap: 10px;
            margin-bottom: 20px;
        }
        
        .quick-btn {
            padding: 10px 16px;
            background: var(--bg-main);
            border: 1px solid var(--border);
            border-radius: 8px;
            font-size: 13px;
            color: var(--text-secondary);
            cursor: pointer;
            transition: all 0.2s;
            display: flex;
            align-items: center;
            gap: 6px;
        }
        
        .quick-btn:hover {
            background: var(--bg-card);
            color: var(--text-primary);
            border-color: var(--primary);
        }
        
        .tabs {
            display: flex;
            gap: 4px;
            margin-bottom: 20px;
            border-bottom: 2px solid var(--border);
            padding-bottom: 0;
        }
        
        .tab {
            padding: 12px 20px;
            font-size: 13px;
            font-weight: 500;
            background: transparent;
            color: var(--text-secondary);
            border: none;
            cursor: pointer;
            border-bottom: 2px solid transparent;
            margin-bottom: -2px;
            transition: all 0.2s;
        }
        
        .tab:hover {
            color: var(--text-primary);
            background: var(--bg-main);
            border-radius: 8px 8px 0 0;
        }
        
        .tab.active {
            color: var(--primary);
            border-bottom-color: var(--primary);
            background: rgba(99, 102, 241, 0.05);
            border-radius: 8px 8px 0 0;
        }
        
        .tab-content {
            display: none;
            animation: fadeIn 0.2s;
        }
        
        .tab-content.active {
            display: block;
        }
        
        @keyframes fadeIn {
            from { opacity: 0; transform: translateY(-10px); }
            to { opacity: 1; transform: translateY(0); }
        }
        
        .input-row {
            display: flex;
            gap: 10px;
            margin-bottom: 12px;
        }
        
        .input-row input {
            flex: 1;
            background: var(--bg-main);
            border: 1px solid var(--border);
            border-radius: 8px;
            padding: 12px 14px;
            font-size: 13px;
            color: var(--text-primary);
        }
        
        .input-row input:focus {
            outline: none;
            border-color: var(--primary);
        }
        
        .remove-btn {
            background: #FEE2E2;
            color: #DC2626;
            border: none;
            border-radius: 8px;
            padding: 8px 14px;
            cursor: pointer;
            font-size: 14px;
            transition: all 0.2s;
        }
        
        .remove-btn:hover {
            background: #FECACA;
        }
        
        .add-btn {
            background: transparent;
            color: var(--text-secondary);
            border: 1px dashed var(--border);
            border-radius: 8px;
            padding: 12px;
            cursor: pointer;
            font-size: 13px;
            width: 100%;
            transition: all 0.2s;
        }
        
        .add-btn:hover {
            background: var(--bg-main);
            color: var(--text-primary);
            border-color: var(--primary);
        }
        
        .body-editor {
            width: 100%;
            min-height: 150px;
            background: var(--bg-main);
            border: 1px solid var(--border);
            border-radius: 12px;
            padding: 16px;
            font-size: 13px;
            font-family: "Monaco", "Menlo", "Ubuntu Mono", monospace;
            color: var(--text-primary);
            resize: vertical;
            line-height: 1.6;
        }
        
        .body-editor:focus {
            outline: none;
            border-color: var(--primary);
            box-shadow: 0 0 0 3px rgba(99, 102, 241, 0.1);
        }
        
        /* 响应面板 */
        .response-panel {
            background: var(--bg-card);
            border-radius: 16px;
            box-shadow: var(--shadow);
            border: 1px solid var(--border);
            overflow: hidden;
        }
        
        .response-header {
            padding: 16px 24px;
            border-bottom: 1px solid var(--border);
            display: flex;
            align-items: center;
            justify-content: space-between;
            background: var(--bg-main);
        }
        
        .response-info {
            display: flex;
            align-items: center;
            gap: 20px;
        }
        
        .status-badge {
            padding: 8px 16px;
            border-radius: 8px;
            font-size: 13px;
            font-weight: 600;
            display: flex;
            align-items: center;
            gap: 6px;
        }
        
        .status-2xx { background: #DCFCE7; color: #16A34A; }
        .status-4xx { background: #FEF3C7; color: #D97706; }
        .status-5xx { background: #FEE2E2; color: #DC2626; }
        
        .response-meta {
            font-size: 12px;
            color: var(--text-secondary);
            display: flex;
            align-items: center;
            gap: 6px;
        }
        
        .response-placeholder {
            color: var(--text-secondary);
            font-size: 13px;
        }
        
        .response-body {
            padding: 20px;
            max-height: 500px;
            overflow: auto;
            background: #0F172A;
        }
        
        .response-body pre {
            color: #E2E8F0;
            font-size: 13px;
            font-family: "Monaco", "Menlo", "Ubuntu Mono", monospace;
            white-space: pre-wrap;
            word-break: break-all;
            line-height: 1.6;
            margin: 0;
        }
        
        .json-key { color: #89B4FA; }
        .json-string { color: #A6E3A1; }
        .json-number { color: #FAB387; }
        .json-boolean { color: #CBA6F7; }
        .json-null { color: #6C7086; }
        
        .empty-state {
            padding: 80px 40px;
            text-align: center;
            color: var(--text-secondary);
        }
        
        .empty-state .icon {
            font-size: 56px;
            margin-bottom: 16px;
            opacity: 0.5;
        }
        
        .empty-state h3 {
            font-size: 16px;
            font-weight: 500;
            margin-bottom: 8px;
            color: var(--text-primary);
        }
        
        .empty-state p {
            font-size: 13px;
            max-width: 300px;
            margin: 0 auto;
        }
        
        /* 历史记录 */
        .history-panel {
            background: var(--bg-card);
            border-radius: 16px;
            padding: 20px;
            box-shadow: var(--shadow);
            border: 1px solid var(--border);
            max-height: 300px;
            overflow-y: auto;
        }
        
        .history-item {
            padding: 12px 14px;
            border-radius: 10px;
            margin-bottom: 8px;
            cursor: pointer;
            transition: all 0.2s;
            border: 1px solid var(--border);
            display: flex;
            align-items: center;
            gap: 12px;
        }
        
        .history-item:hover {
            background: var(--bg-main);
            border-color: var(--primary);
        }
        
        .history-status {
            width: 8px;
            height: 8px;
            border-radius: 50%;
            flex-shrink: 0;
        }
        
        .history-status.success { background: var(--success); }
        .history-status.error { background: var(--danger); }
        
        .history-info {
            flex: 1;
            overflow: hidden;
        }
        
        .history-url {
            font-size: 12px;
            font-family: monospace;
            color: var(--text-primary);
            overflow: hidden;
            text-overflow: ellipsis;
            white-space: nowrap;
        }
        
        .history-meta {
            font-size: 11px;
            color: var(--text-secondary);
            margin-top: 4px;
        }
        
        /* 响应式 */
        @media (max-width: 1024px) {
            .main-container {
                grid-template-columns: 1fr;
            }
            
            .sidebar {
                order: 2;
            }
        }
        
        @media (max-width: 640px) {
            .main-container {
                padding: 16px;
            }
            
            .header-content {
                flex-direction: column;
                gap: 12px;
                align-items: flex-start;
            }
            
            .url-bar {
                flex-direction: column;
            }
            
            .method-select {
                width: 100%;
            }
        }
        
        /* 滚动条 */
        ::-webkit-scrollbar {
            width: 6px;
            height: 6px;
        }
        
        ::-webkit-scrollbar-track {
            background: var(--bg-main);
            border-radius: 3px;
        }
        
        ::-webkit-scrollbar-thumb {
            background: #CBD5E1;
            border-radius: 3px;
        }
        
        ::-webkit-scrollbar-thumb:hover {
            background: #94A3B8;
        }
        
        /* 工具提示 */
        .tooltip {
            position: relative;
        }
        
        .tooltip::after {
            content: attr(data-tooltip);
            position: absolute;
            bottom: 100%;
            left: 50%;
            transform: translateX(-50%);
            background: #1E293B;
            color: white;
            padding: 6px 12px;
            border-radius: 6px;
            font-size: 12px;
            white-space: nowrap;
            opacity: 0;
            pointer-events: none;
            transition: opacity 0.2s;
            margin-bottom: 8px;
        }
        
        .tooltip:hover::after {
            opacity: 1;
        }
    </style>
</head>
<body>
    <header class="header">
        <div class="header-content">
            <div class="logo">
                <div class="logo-icon">⚡</div>
                <div>
                    <h1>API 调试工具</h1>
                    <div class="logo-subtitle">安信合同管理系统</div>
                </div>
            </div>
            <div class="header-info">
                <span class="version-badge">v{{.Version}}</span>
                <div class="status-indicator">
                    <div class="status-dot"></div>
                    <span>服务运行中</span>
                </div>
                <span class="time-display">{{.Time}}</span>
            </div>
        </div>
    </header>
    
    <main class="main-container">
        <aside class="sidebar">
            <div class="sidebar-card">
                <div class="sidebar-header">
                    <div class="sidebar-title">
                        <span>📋</span>
                        <span>接口列表</span>
                    </div>
                </div>
                
                <input type="text" class="search-box" placeholder="搜索接口..." id="searchBox">
                
                <div id="api-list">
                    <div class="api-group">
                        <div class="group-title">🔓 公共接口</div>
                        <div class="api-item" data-method="GET" data-url="/" data-body="" data-auth="false" data-desc="服务信息">
                            <span class="method-badge method-get">GET</span>
                            <span class="api-path">/</span>
                            <span class="api-badge">公开</span>
                        </div>
                        <div class="api-item" data-method="GET" data-url="/health" data-body="" data-auth="false" data-desc="健康检查">
                            <span class="method-badge method-get">GET</span>
                            <span class="api-path">/health</span>
                            <span class="api-badge">公开</span>
                        </div>
                        <div class="api-item" data-method="POST" data-url="/api/auth/login" data-body='{"username":"admin","password":"admin123"}' data-auth="false" data-desc="用户登录">
                            <span class="method-badge method-post">POST</span>
                            <span class="api-path">/api/auth/login</span>
                            <span class="api-badge">公开</span>
                        </div>
                        <div class="api-item" data-method="POST" data-url="/api/auth/register" data-body='{"username":"test","email":"test@test.com","password":"123456","full_name":"测试"}' data-auth="false" data-desc="用户注册">
                            <span class="method-badge method-post">POST</span>
                            <span class="api-path">/api/auth/register</span>
                            <span class="api-badge">公开</span>
                        </div>
                    </div>
                    
                    <div class="api-group">
                        <div class="group-title">👤 用户管理</div>
                        <div class="api-item" data-method="GET" data-url="/api/auth/users?skip=0&limit=10" data-body="" data-auth="true" data-desc="获取用户列表">
                            <span class="method-badge method-get">GET</span>
                            <span class="api-path">/api/auth/users</span>
                            <span class="api-badge">管理员</span>
                        </div>
                        <div class="api-item" data-method="GET" data-url="/api/auth/users/1" data-body="" data-auth="true" data-desc="获取用户详情">
                            <span class="method-badge method-get">GET</span>
                            <span class="api-path">/api/auth/users/:id</span>
                            <span class="api-badge">管理员</span>
                        </div>
                        <div class="api-item" data-method="PUT" data-url="/api/auth/users/1" data-body='{"full_name":"新名字"}' data-auth="true" data-desc="更新用户">
                            <span class="method-badge method-put">PUT</span>
                            <span class="api-path">/api/auth/users/:id</span>
                            <span class="api-badge">管理员</span>
                        </div>
                        <div class="api-item" data-method="DELETE" data-url="/api/auth/users/1" data-body="" data-auth="true" data-desc="删除用户">
                            <span class="method-badge method-delete">DELETE</span>
                            <span class="api-path">/api/auth/users/:id</span>
                            <span class="api-badge">管理员</span>
                        </div>
                    </div>
                    
                    <div class="api-group">
                        <div class="group-title">🏢 客户管理</div>
                        <div class="api-item" data-method="GET" data-url="/api/customers" data-body="" data-auth="true" data-desc="获取客户列表">
                            <span class="method-badge method-get">GET</span>
                            <span class="api-path">/api/customers</span>
                            <span class="api-badge">认证</span>
                        </div>
                        <div class="api-item" data-method="POST" data-url="/api/customers" data-body='{"name":"测试客户","code":"C001","type":"customer","contact_person":"张三"}' data-auth="true" data-desc="创建客户">
                            <span class="method-badge method-post">POST</span>
                            <span class="api-path">/api/customers</span>
                            <span class="api-badge">认证</span>
                        </div>
                        <div class="api-item" data-method="GET" data-url="/api/contract-types" data-body="" data-auth="true" data-desc="获取合同类型">
                            <span class="method-badge method-get">GET</span>
                            <span class="api-path">/api/contract-types</span>
                            <span class="api-badge">认证</span>
                        </div>
                    </div>
                    
                    <div class="api-group">
                        <div class="group-title">📄 合同管理</div>
                        <div class="api-item" data-method="GET" data-url="/api/contracts" data-body="" data-auth="true" data-desc="获取合同列表">
                            <span class="method-badge method-get">GET</span>
                            <span class="api-path">/api/contracts</span>
                            <span class="api-badge">认证</span>
                        </div>
                        <div class="api-item" data-method="POST" data-url="/api/contracts" data-body='{"contract_no":"CT2024001","title":"测试合同","customer_id":1,"contract_type_id":1,"amount":100000,"currency":"CNY","status":"draft"}' data-auth="true" data-desc="创建合同">
                            <span class="method-badge method-post">POST</span>
                            <span class="api-path">/api/contracts</span>
                            <span class="api-badge">认证</span>
                        </div>
                        <div class="api-item" data-method="GET" data-url="/api/contracts/1" data-body="" data-auth="true" data-desc="获取合同详情">
                            <span class="method-badge method-get">GET</span>
                            <span class="api-path">/api/contracts/:id</span>
                            <span class="api-badge">认证</span>
                        </div>
                        <div class="api-item" data-method="POST" data-url="/api/contracts/1/documents" data-body="" data-auth="true" data-desc="上传文档" data-file="true">
                            <span class="method-badge method-post">POST</span>
                            <span class="api-path">/api/contracts/:id/documents</span>
                            <span class="api-badge">文件</span>
                        </div>
                    </div>
                    
                    <div class="api-group">
                        <div class="group-title">📊 数据统计</div>
                        <div class="api-item" data-method="GET" data-url="/api/statistics" data-body="" data-auth="true" data-desc="获取统计数据">
                            <span class="method-badge method-get">GET</span>
                            <span class="api-path">/api/statistics</span>
                            <span class="api-badge">认证</span>
                        </div>
                        <div class="api-item" data-method="GET" data-url="/api/expiring-contracts?days=30" data-body="" data-auth="true" data-desc="即将到期合同">
                            <span class="method-badge method-get">GET</span>
                            <span class="api-path">/api/expiring-contracts</span>
                            <span class="api-badge">认证</span>
                        </div>
                    </div>
                    
                    <div class="api-group">
                        <div class="group-title">📋 审批管理</div>
                        <div class="api-item" data-method="GET" data-url="/api/contracts/1/approvals" data-body="" data-auth="true" data-desc="获取审批记录">
                            <span class="method-badge method-get">GET</span>
                            <span class="api-path">/api/contracts/:id/approvals</span>
                            <span class="api-badge">认证</span>
                        </div>
                        <div class="api-item" data-method="GET" data-url="/api/pending-approvals" data-body="" data-auth="true" data-desc="待审批列表">
                            <span class="method-badge method-get">GET</span>
                            <span class="api-path">/api/pending-approvals</span>
                            <span class="api-badge">认证</span>
                        </div>
                    </div>
                </div>
            </div>
            
            <div class="sidebar-card">
                <div class="sidebar-header">
                    <div class="sidebar-title">
                        <span>🕐</span>
                        <span>请求历史</span>
                    </div>
                    <button class="quick-btn" onclick="clearHistory()">清空</button>
                </div>
                <div id="history-list" class="history-panel">
                    <div class="empty-state" style="padding: 30px 20px;">
                        <div style="font-size: 32px; margin-bottom: 8px;">📭</div>
                        <div style="font-size: 13px;">暂无请求历史</div>
                    </div>
                </div>
            </div>
        </aside>
        
        <div class="content">
            <div class="request-panel">
                <div class="url-bar">
                    <select class="method-select" id="method">
                        <option value="GET">GET</option>
                        <option value="POST">POST</option>
                        <option value="PUT">PUT</option>
                        <option value="DELETE">DELETE</option>
                        <option value="PATCH">PATCH</option>
                    </select>
                    <input type="text" class="url-input" id="url" placeholder="输入请求地址..." value="/">
                    <button class="send-btn" id="sendBtn" onclick="sendRequest()">
                        <span>发送</span>
                        <span>→</span>
                    </button>
                </div>
                
                <div class="token-section">
                    <div class="token-header">
                        <div class="token-label">
                            <span>🔐</span>
                            <span>Auth Token</span>
                        </div>
                        <span class="token-status" id="tokenStatus">有效</span>
                    </div>
                    <input type="text" class="token-input" id="token" placeholder="粘贴 JWT Token... 登录后自动保存">
                </div>
                
                <div class="quick-actions">
                    <button class="quick-btn" onclick="clearAll()">
                        <span>🗑️</span>
                        <span>清除</span>
                    </button>
                    <button class="quick-btn" onclick="formatJson()">
                        <span>✨</span>
                        <span>格式化 JSON</span>
                    </button>
                    <button class="quick-btn" onclick="copyResponse()">
                        <span>📋</span>
                        <span>复制响应</span>
                    </button>
                    <button class="quick-btn" onclick="saveTemplate()">
                        <span>💾</span>
                        <span>保存模板</span>
                    </button>
                </div>
                
                <div class="tabs">
                    <button class="tab active" onclick="switchTab('params', this)" data-tab="params">Parameters</button>
                    <button class="tab" onclick="switchTab('headers', this)" data-tab="headers">Headers</button>
                    <button class="tab" onclick="switchTab('body', this)" data-tab="body">Body</button>
                    <button class="tab" onclick="switchTab('file', this)" data-tab="file">File</button>
                </div>
                
                <div class="tab-content active" id="tab-params">
                    <div id="params-container">
                        <div class="input-row">
                            <input type="text" placeholder="Key" class="param-key">
                            <input type="text" placeholder="Value" class="param-value">
                            <button class="remove-btn" onclick="this.parentElement.remove()">×</button>
                        </div>
                    </div>
                    <button class="add-btn" onclick="addParam()">+ 添加参数</button>
                </div>
                
                <div class="tab-content" id="tab-headers">
                    <div id="headers-container">
                        <div class="input-row">
                            <input type="text" placeholder="Key" value="Content-Type" class="header-key">
                            <input type="text" placeholder="Value" value="application/json" class="header-value">
                            <button class="remove-btn" onclick="this.parentElement.remove()">×</button>
                        </div>
                    </div>
                    <button class="add-btn" onclick="addHeader()">+ 添加请求头</button>
                </div>
                
                <div class="tab-content" id="tab-body">
                    <textarea class="body-editor" id="body" placeholder='{
  "key": "value"
}'></textarea>
                </div>
                
                <div class="tab-content" id="tab-file">
                    <div style="border: 2px dashed var(--border); border-radius: 12px; padding: 40px; text-align: center; cursor: pointer;" onclick="document.getElementById('fileInput').click()">
                        <div style="font-size: 48px; margin-bottom: 16px; opacity: 0.5;">📁</div>
                        <div style="font-size: 14px; color: var(--text-primary); margin-bottom: 8px;">点击选择文件或拖拽到此处</div>
                        <div style="font-size: 12px; color: var(--text-secondary);">支持 PDF, DOCX, XLSX, JPG, PNG 等格式</div>
                        <input type="file" id="fileInput" style="display: none;" onchange="handleFileSelect(this)">
                    </div>
                    <div id="file-info" style="margin-top: 16px; display: none;">
                        <div class="input-row">
                            <input type="text" id="selected-file" readonly placeholder="已选择文件">
                            <button class="remove-btn" onclick="clearFile()">×</button>
                        </div>
                    </div>
                </div>
            </div>
            
            <div class="response-panel">
                <div class="response-header">
                    <div class="response-info" id="response-info" style="display: none;">
                        <span class="status-badge" id="status-badge">
                            <span id="status-icon">●</span>
                            <span id="status-text">200 OK</span>
                        </span>
                        <span class="response-meta">
                            <span>⏱️</span>
                            <span id="response-time">0ms</span>
                        </span>
                        <span class="response-meta">
                            <span>📦</span>
                            <span id="response-size">0 B</span>
                        </span>
                    </div>
                    <span class="response-placeholder" id="response-placeholder">
                        等待发送请求...
                    </span>
                </div>
                <div class="response-body" id="response-body">
                    <div class="empty-state">
                        <div class="icon">🚀</div>
                        <h3>准备就绪</h3>
                        <p>点击左侧接口或输入 URL 发送请求，响应结果将显示在这里</p>
                    </div>
                </div>
            </div>
        </div>
    </main>
    
    <script>
        // 状态管理
        let requestHistory = JSON.parse(localStorage.getItem('apiHistory') || '[]');
        let currentFile = null;
        
        // 初始化
        document.addEventListener('DOMContentLoaded', () => {
            initEventListeners();
            renderHistory();
            checkToken();
        });
        
        // 事件监听
        function initEventListeners() {
            // 接口点击事件
            document.querySelectorAll('.api-item').forEach(item => {
                item.addEventListener('click', function() {
                    selectApi(this);
                });
            });
            
            // 搜索功能
            document.getElementById('searchBox').addEventListener('input', function(e) {
                filterApis(e.target.value);
            });
            
            // 回车发送
            document.getElementById('url').addEventListener('keypress', e => {
                if (e.key === 'Enter') sendRequest();
            });
            
            // Token 变化检测
            document.getElementById('token').addEventListener('input', checkToken);
        }
        
        // 选择API
        function selectApi(item) {
            const method = item.getAttribute('data-method');
            const url = item.getAttribute('data-url');
            const body = item.getAttribute('data-body');
            const needsAuth = item.getAttribute('data-auth') === 'true';
            const hasFile = item.getAttribute('data-file') === 'true';
            
            document.getElementById('method').value = method;
            document.getElementById('url').value = url;
            document.getElementById('body').value = body || '';
            
            // 更新活动状态
            document.querySelectorAll('.api-item').forEach(i => i.classList.remove('active'));
            item.classList.add('active');
            
            // 切换到合适的标签页
            if (hasFile) {
                switchTab('file', document.querySelector('[data-tab="file"]'));
            } else if (body) {
                switchTab('body', document.querySelector('[data-tab="body"]'));
            } else {
                switchTab('params', document.querySelector('[data-tab="params"]'));
            }
            
            // 检查是否需要认证
            if (needsAuth && !document.getElementById('token').value) {
                showToast('此接口需要认证，请先登录获取 Token', 'warning');
            }
        }
        
        // 过滤接口
        function filterApis(keyword) {
            const items = document.querySelectorAll('.api-item');
            const lowerKeyword = keyword.toLowerCase();
            
            items.forEach(item => {
                const url = item.getAttribute('data-url').toLowerCase();
                const desc = (item.getAttribute('data-desc') || '').toLowerCase();
                const method = item.getAttribute('data-method').toLowerCase();
                
                const match = url.includes(lowerKeyword) || 
                              desc.includes(lowerKeyword) || 
                              method.includes(lowerKeyword);
                
                item.style.display = match ? 'flex' : 'none';
            });
        }
        
        // 切换标签页
        function switchTab(tab, el) {
            document.querySelectorAll('.tab').forEach(t => t.classList.remove('active'));
            document.querySelectorAll('.tab-content').forEach(c => c.classList.remove('active'));
            
            if (el) el.classList.add('active');
            document.getElementById('tab-' + tab).classList.add('active');
        }
        
        // 添加参数行
        function addParam() {
            const container = document.getElementById('params-container');
            const div = document.createElement('div');
            div.className = 'input-row';
            div.innerHTML = '<input type="text" placeholder="Key" class="param-key">' +
                            '<input type="text" placeholder="Value" class="param-value">' +
                            '<button class="remove-btn" onclick="this.parentElement.remove()">&times;</button>';
            container.appendChild(div);
        }
        
        // 添加请求头行
        function addHeader() {
            const container = document.getElementById('headers-container');
            const div = document.createElement('div');
            div.className = 'input-row';
            div.innerHTML = '<input type="text" placeholder="Key" class="header-key">' +
                            '<input type="text" placeholder="Value" class="header-value">' +
                            '<button class="remove-btn" onclick="this.parentElement.remove()">&times;</button>';
            container.appendChild(div);
        }
        
        // 获取参数
        function getParams() {
            const params = new URLSearchParams();
            document.querySelectorAll('#params-container .input-row').forEach(row => {
                const key = row.querySelector('.param-key').value;
                const value = row.querySelector('.param-value').value;
                if (key) params.append(key, value);
            });
            return params.toString();
        }
        
        // 获取请求头
        function getHeaders() {
            const headers = {};
            document.querySelectorAll('#headers-container .input-row').forEach(row => {
                const key = row.querySelector('.header-key').value;
                const value = row.querySelector('.header-value').value;
                if (key) headers[key] = value;
            });
            return headers;
        }
        
        // 检查Token状态
        function checkToken() {
            const token = document.getElementById('token').value;
            const statusEl = document.getElementById('tokenStatus');
            
            if (!token) {
                statusEl.textContent = '未设置';
                statusEl.className = 'token-status expired';
                return;
            }
            
            // 简单检查token格式
            if (token.split('.').length === 3) {
                statusEl.textContent = '有效';
                statusEl.className = 'token-status';
            } else {
                statusEl.textContent = '格式错误';
                statusEl.className = 'token-status expired';
            }
        }
        
        // 发送请求
        async function sendRequest() {
            const method = document.getElementById('method').value;
            let url = document.getElementById('url').value;
            const body = document.getElementById('body').value;
            const token = document.getElementById('token').value;
            const sendBtn = document.getElementById('sendBtn');
            
            // 添加 loading 状态
            sendBtn.classList.add('loading');
            sendBtn.innerHTML = '<span>发送中...</span>';
            
            // 处理URL
            if (!url.startsWith('http')) {
                url = window.location.origin + url;
            }
            
            // 添加查询参数
            const params = getParams();
            if (params) url += (url.includes('?') ? '&' : '?') + params;
            
            // 构建请求头
            const headers = getHeaders();
            if (token) headers['Authorization'] = 'Bearer ' + token;
            
            // 构建请求选项
            const options = { method, headers };
            
            // 处理文件上传
            if (currentFile && method === 'POST') {
                const formData = new FormData();
                formData.append('file', currentFile);
                if (body) {
                    try {
                        const jsonData = JSON.parse(body);
                        Object.keys(jsonData).forEach(key => {
                            formData.append(key, jsonData[key]);
                        });
                    } catch (e) {
                        // body不是JSON，忽略
                    }
                }
                options.body = formData;
                delete headers['Content-Type']; // 让浏览器设置正确的boundary
            } else if (body && ['POST', 'PUT', 'PATCH'].includes(method)) {
                options.body = body;
            }
            
            const startTime = Date.now();
            
            try {
                const response = await fetch(url, options);
                const time = Date.now() - startTime;
                const size = response.headers.get('content-length') || '-';
                const status = response.status;
                const statusText = response.statusText;
                
                // 更新状态显示
                const statusBadge = document.getElementById('status-badge');
                document.getElementById('status-text').textContent = status + ' ' + statusText;
                document.getElementById('status-icon').textContent = status < 300 ? '✓' : status < 500 ? '⚠' : '✕';
                statusBadge.className = 'status-badge ' + (status < 300 ? 'status-2xx' : status < 500 ? 'status-4xx' : 'status-5xx');
                
                document.getElementById('response-info').style.display = 'flex';
                document.getElementById('response-time').textContent = time + 'ms';
                document.getElementById('response-size').textContent = formatSize(parseInt(size)) || '-';
                document.getElementById('response-placeholder').style.display = 'none';
                
                // 获取响应内容
                const text = await response.text();
                const bodyEl = document.getElementById('response-body');
                
                try {
                    const json = JSON.parse(text);
                    bodyEl.innerHTML = '<pre>' + syntaxHighlight(json) + '</pre>';
                    
                    // 自动保存token
                    if (json.access_token) {
                        document.getElementById('token').value = json.access_token;
                        checkToken();
                        showToast('Token 已自动保存', 'success');
                    }
                } catch (e) {
                    bodyEl.innerHTML = '<pre>' + escapeHtml(text) + '</pre>';
                }
                
                // 添加到历史记录
                addToHistory({
                    method,
                    url: url.replace(window.location.origin, ''),
                    status,
                    time,
                    timestamp: new Date().toISOString(),
                    success: status < 400
                });
                
            } catch (error) {
                document.getElementById('response-body').innerHTML = 
                    '<pre style="color: #F87171;">❌ 请求失败\n\n' + escapeHtml(error.message) + '</pre>';
                
                addToHistory({
                    method,
                    url: url.replace(window.location.origin, ''),
                    status: 0,
                    time: 0,
                    timestamp: new Date().toISOString(),
                    success: false,
                    error: error.message
                });
            } finally {
                // 恢复按钮状态
                sendBtn.classList.remove('loading');
                sendBtn.innerHTML = '<span>发送</span><span>→</span>';
            }
        }
        
        // 语法高亮
        function syntaxHighlight(json) {
            const formatted = JSON.stringify(json, null, 2);
            return escapeHtml(formatted).replace(/("(\\u[a-zA-Z0-9]{4}|\\[^u]|[^\\"])*"(\s*:)?|\b(true|false|null)\b|-?\d+(?:\.\d*)?(?:[eE][+\-]?\d+)?)/g, 
                function (match) {
                    let cls = 'json-number';
                    if (/^"/.test(match)) {
                        if (/:$/.test(match)) cls = 'json-key';
                        else cls = 'json-string';
                    } else if (/true|false/.test(match)) cls = 'json-boolean';
                    else if (/null/.test(match)) cls = 'json-null';
                    return '<span class="' + cls + '">' + match + '</span>';
                }
            );
        }
        
        // HTML转义
        function escapeHtml(text) {
            const div = document.createElement('div');
            div.textContent = text;
            return div.innerHTML;
        }
        
        // 格式化文件大小
        function formatSize(bytes) {
            if (bytes === 0) return '0 B';
            const k = 1024;
            const sizes = ['B', 'KB', 'MB', 'GB'];
            const i = Math.floor(Math.log(bytes) / Math.log(k));
            return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
        }
        
        // 清空所有
        function clearAll() {
            document.getElementById('url').value = '/';
            document.getElementById('body').value = '';
            document.getElementById('method').value = 'GET';
            document.getElementById('response-info').style.display = 'none';
            document.getElementById('response-placeholder').style.display = 'block';
            document.getElementById('response-placeholder').textContent = '等待发送请求...';
            document.getElementById('response-body').innerHTML = '<div class="empty-state">' +
                    '<div class="icon">🚀</div>' +
                    '<h3>准备就绪</h3>' +
                    '<p>点击左侧接口或输入 URL 发送请求，响应结果将显示在这里</p>' +
                '</div>';
            clearFile();
            
            // 清空参数
            document.getElementById('params-container').innerHTML = '<div class="input-row">' +
                    '<input type="text" placeholder="Key" class="param-key">' +
                    '<input type="text" placeholder="Value" class="param-value">' +
                    '<button class="remove-btn" onclick="this.parentElement.remove()">&times;</button>' +
                '</div>';
            
            // 重置请求头
            document.getElementById('headers-container').innerHTML = '<div class="input-row">' +
                    '<input type="text" placeholder="Key" value="Content-Type" class="header-key">' +
                    '<input type="text" placeholder="Value" value="application/json" class="header-value">' +
                    '<button class="remove-btn" onclick="this.parentElement.remove()">&times;</button>' +
                '</div>';
        }
        
        // 格式化JSON
        function formatJson() {
            const body = document.getElementById('body');
            try {
                const json = JSON.parse(body.value);
                body.value = JSON.stringify(json, null, 2);
                showToast('JSON 已格式化', 'success');
            } catch (e) {
                showToast('无效的 JSON 格式', 'error');
            }
        }
        
        // 复制响应
        function copyResponse() {
            const responseBody = document.getElementById('response-body');
            const text = responseBody.textContent;
            
            if (!text || text.includes('准备就绪') || text.includes('等待发送')) {
                showToast('没有可复制的内容', 'warning');
                return;
            }
            
            navigator.clipboard.writeText(text).then(() => {
                showToast('响应内容已复制到剪贴板', 'success');
            }).catch(() => {
                showToast('复制失败', 'error');
            });
        }
        
        // 文件处理
        function handleFileSelect(input) {
            const file = input.files[0];
            if (file) {
                currentFile = file;
                document.getElementById('selected-file').value = file.name + ' (' + formatSize(file.size) + ')';
                document.getElementById('file-info').style.display = 'block';
                
                // 切换到文件标签
                switchTab('file', document.querySelector('[data-tab="file"]'));
            }
        }
        
        function clearFile() {
            currentFile = null;
            document.getElementById('fileInput').value = '';
            document.getElementById('file-info').style.display = 'none';
        }
        
        // 历史记录管理
        function addToHistory(record) {
            requestHistory.unshift(record);
            if (requestHistory.length > 50) requestHistory.pop();
            localStorage.setItem('apiHistory', JSON.stringify(requestHistory));
            renderHistory();
        }
        
        function renderHistory() {
            const container = document.getElementById('history-list');
            
            if (requestHistory.length === 0) {
                container.innerHTML = '<div class="empty-state" style="padding: 30px 20px;">' +
                        '<div style="font-size: 32px; margin-bottom: 8px;">&#128237;</div>' +
                        '<div style="font-size: 13px;">暂无请求历史</div>' +
                    '</div>';
                return;
            }
            
            container.innerHTML = requestHistory.map((item, index) => 
                '<div class="history-item" onclick="loadFromHistory(' + index + ')">' +
                    '<div class="history-status ' + (item.success ? 'success' : 'error') + '"></div>' +
                    '<div class="history-info">' +
                        '<div class="history-url">' + item.method + ' ' + item.url + '</div>' +
                        '<div class="history-meta">' +
                            (item.status ? item.status + ' • ' : '') + item.time + 'ms • ' + new Date(item.timestamp).toLocaleTimeString() +
                        '</div>' +
                    '</div>' +
                '</div>'
            ).join('');
        }
        
        function loadFromHistory(index) {
            const item = requestHistory[index];
            document.getElementById('method').value = item.method;
            document.getElementById('url').value = item.url;
            showToast('已加载历史请求', 'success');
        }
        
        function clearHistory() {
            requestHistory = [];
            localStorage.removeItem('apiHistory');
            renderHistory();
            showToast('历史记录已清空', 'success');
        }
        
        // 保存模板
        function saveTemplate() {
            const template = {
                method: document.getElementById('method').value,
                url: document.getElementById('url').value,
                headers: getHeaders(),
                body: document.getElementById('body').value
            };
            
            const dataStr = JSON.stringify(template, null, 2);
            const dataUri = 'data:application/json;charset=utf-8,'+ encodeURIComponent(dataStr);
            
            const link = document.createElement('a');
            link.setAttribute('href', dataUri);
            link.setAttribute('download', 'api-template.json');
            document.body.appendChild(link);
            link.click();
            document.body.removeChild(link);
            
            showToast('模板已保存', 'success');
        }
        
        // Toast 通知
        function showToast(message, type = 'info') {
            const toast = document.createElement('div');
            toast.style.cssText = 'position: fixed; ' +
                'bottom: 24px; ' +
                'right: 24px; ' +
                'padding: 14px 24px; ' +
                'background: ' + (type === 'success' ? '#10B981' : type === 'error' ? '#EF4444' : type === 'warning' ? '#F59E0B' : '#3B82F6') + '; ' +
                'color: white; ' +
                'border-radius: 10px; ' +
                'font-size: 14px; ' +
                'font-weight: 500; ' +
                'box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.2); ' +
                'z-index: 1000; ' +
                'animation: slideIn 0.3s ease;';
            toast.textContent = message;
            document.body.appendChild(toast);
            
            setTimeout(() => {
                toast.style.animation = 'slideOut 0.3s ease';
                setTimeout(() => toast.remove(), 300);
            }, 3000);
        }
        
        // 添加动画样式
        const style = document.createElement('style');
        style.textContent = '@keyframes slideIn { ' +
            'from { transform: translateX(100%); opacity: 0; } ' +
            'to { transform: translateX(0); opacity: 1; } ' +
            '} ' +
            '@keyframes slideOut { ' +
            'from { transform: translateX(0); opacity: 1; } ' +
            'to { transform: translateX(100%); opacity: 0; } ' +
            '}';
        document.head.appendChild(style);
    </script>
</body>
</html>`

func main() {
	if err := config.LoadConfig(); err != nil {
		panic("Failed to load config: " + err.Error())
	}

	if err := models.InitDB(); err != nil {
		panic("Failed to connect database: " + err.Error())
	}

	if err := models.InitAdmin(); err != nil {
		fmt.Println("Warning: Failed to create admin user: " + err.Error())
	}

	r := gin.Default()

	r.SetHTMLTemplate(template.Must(template.New("").Parse(apiDocs)))

	r.Use(gin.Recovery())
	r.Use(gin.Logger())
	r.Use(middleware.ErrorHandlerMiddleware())

	r.Use(middleware.SecureHeadersMiddleware())
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.RateLimitMiddleware())
	r.Use(middleware.RequestSizeLimitMiddleware(10 << 20))
	r.Use(middleware.XssProtectionMiddleware())

	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "", gin.H{
			"Version": config.AppConfig.AppVersion,
			"Time":    time.Now().Format("2006-01-02 15:04:05"),
		})
	})

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "healthy",
			"time":   time.Now().Unix(),
		})
	})

	// 静态文件服务 - 上传的文件
	r.Static("/uploads", "./uploads")

	// Swagger文档
	r.Static("/docs", "./docs")
	r.GET("/api-docs", func(c *gin.Context) {
		c.HTML(200, "", gin.H{})
	})

	authHandler := handlers.NewAuthHandler()
	customerHandler := handlers.NewCustomerHandler()
	contractHandler := handlers.NewContractHandler()
	approvalHandler := handlers.NewApprovalHandler()
	workflowHandler := handlers.NewWorkflowHandler(models.DB)
	auditHandler := handlers.NewAuditHandler()

	auth := r.Group("/api/auth")
	auth.Use(middleware.RateLimitMiddleware())
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.GET("/users", middleware.AuthMiddleware(), middleware.AdminRequiredMiddleware(), authHandler.GetUsers)
		auth.GET("/users/:user_id", middleware.AuthMiddleware(), middleware.AdminRequiredMiddleware(), authHandler.GetUserByID)
		auth.PUT("/users/:user_id", middleware.AuthMiddleware(), middleware.AdminRequiredMiddleware(), authHandler.UpdateUser)
		auth.DELETE("/users/:user_id", middleware.AuthMiddleware(), middleware.AdminRequiredMiddleware(), authHandler.DeleteUser)
	}

	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware())
	api.Use(handlers.AuditLogMiddleware(handlers.GetAuditService()))
	{
		api.GET("/customers", customerHandler.GetCustomers)
		api.GET("/customers/:customer_id", customerHandler.GetCustomerByID)
		api.POST("/customers", customerHandler.CreateCustomer)
		api.PUT("/customers/:customer_id", customerHandler.UpdateCustomer)
		api.DELETE("/customers/:customer_id", customerHandler.DeleteCustomer)

		api.GET("/contract-types", customerHandler.GetContractTypes)
		api.POST("/contract-types", customerHandler.CreateContractType)
		api.PUT("/contract-types/:type_id", customerHandler.UpdateContractType)
		api.DELETE("/contract-types/:type_id", customerHandler.DeleteContractType)

		api.GET("/contracts", contractHandler.GetContracts)
		api.POST("/contracts", contractHandler.CreateContract)
		api.GET("/contracts/:contract_id", contractHandler.GetContractByID)
		api.PUT("/contracts/:contract_id", contractHandler.UpdateContract)
		api.PUT("/contracts/:contract_id/status", contractHandler.UpdateContractStatus)
		api.POST("/contracts/:contract_id/status-change", contractHandler.CreateStatusChangeRequest)
		api.GET("/contracts/:contract_id/status-change", contractHandler.GetStatusChangeRequests)
		api.POST("/contracts/:contract_id/archive", contractHandler.ArchiveContract)
		api.DELETE("/contracts/:contract_id", contractHandler.DeleteContract)
		api.GET("/contracts/:contract_id/lifecycle", contractHandler.GetContractLifecycle)

		api.GET("/contracts/:contract_id/executions", contractHandler.GetContractExecutions)
		api.POST("/contracts/:contract_id/executions", contractHandler.CreateContractExecution)
		api.DELETE("/executions/:execution_id", contractHandler.DeleteExecution)

		api.GET("/contracts/:contract_id/documents", contractHandler.GetContractDocuments)
		api.POST("/contracts/:contract_id/documents", contractHandler.CreateContractDocument)
		api.GET("/documents/:document_id/preview", contractHandler.PreviewDocument)
		api.DELETE("/documents/:document_id", contractHandler.DeleteDocument)

		api.GET("/contracts/:contract_id/approvals", approvalHandler.GetContractApprovals)
		api.POST("/contracts/:contract_id/approvals", approvalHandler.CreateApproval)
		api.PUT("/approvals/:approval_id", approvalHandler.UpdateApproval)
		api.GET("/pending-approvals", approvalHandler.GetPendingApprovals)

		// 工作流审批路由
		api.POST("/workflow/create", workflowHandler.CreateWorkflow)
		api.GET("/workflow/:contract_id", workflowHandler.GetWorkflow)
		api.GET("/workflow/:contract_id/pending", workflowHandler.GetMyPendingApproval)
		api.POST("/workflow/approve", workflowHandler.Approve)
		api.POST("/workflow/reject", workflowHandler.Reject)

		api.GET("/pending-status-changes", contractHandler.GetPendingStatusChangeApprovals)
		api.POST("/status-change-requests/:request_id/approve", contractHandler.ApproveStatusChangeRequest)
		api.POST("/status-change-requests/:request_id/reject", contractHandler.RejectStatusChangeRequest)

		api.GET("/contracts/:contract_id/reminders", approvalHandler.GetContractReminders)
		api.POST("/contracts/:contract_id/reminders", approvalHandler.CreateReminder)

		api.POST("/reminders/:reminder_id/send", approvalHandler.SendReminder)

		api.GET("/expiring-contracts", approvalHandler.GetExpiringContracts)
		api.GET("/statistics", approvalHandler.GetStatistics)
		api.GET("/notifications/count", approvalHandler.GetNotificationCounts)

		api.GET("/audit-logs", middleware.AdminRequiredMiddleware(), auditHandler.GetAuditLogs)
		api.DELETE("/audit-logs/:id", middleware.AdminRequiredMiddleware(), auditHandler.DeleteAuditLog)
		api.POST("/audit-logs/batch-delete", middleware.AdminRequiredMiddleware(), auditHandler.DeleteAuditLogs)
		api.GET("/audit-logs/export", middleware.AdminRequiredMiddleware(), auditHandler.ExportAuditLogs)
	}

	_ = approvalHandler
	_ = contractHandler

	addr := ":8000"
	fmt.Printf("API 调试页面: http://localhost%s\n", addr)
	fmt.Printf("Swagger 文档: http://localhost%s/docs/swagger.html\n", addr)
	r.Run(addr)
}
