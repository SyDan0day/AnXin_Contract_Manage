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
    <title>API 调试 - 安信合同管理系统</title>
    <link href="https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700&display=swap" rel="stylesheet">
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body {
            font-family: 'Inter', -apple-system, BlinkMacSystemFont, sans-serif;
            background: linear-gradient(135deg, #F8FAFC 0%, #E2E8F0 100%);
            min-height: 100vh;
            color: #1E293B;
        }
        
        .header {
            background: white;
            padding: 16px 32px;
            box-shadow: 0 1px 3px rgba(0,0,0,0.04);
            display: flex;
            align-items: center;
            justify-content: space-between;
            position: sticky;
            top: 0;
            z-index: 100;
        }
        .header-left { display: flex; align-items: center; gap: 12px; }
        .logo {
            width: 36px;
            height: 36px;
            background: linear-gradient(135deg, #6366F1, #8B5CF6);
            border-radius: 10px;
            display: flex;
            align-items: center;
            justify-content: center;
            color: white;
            font-weight: bold;
        }
        .header h1 { font-size: 18px; font-weight: 600; color: #1E293B; }
        .header-right { display: flex; align-items: center; gap: 16px; }
        .status-dot {
            width: 8px; height: 8px;
            background: #10B981;
            border-radius: 50%;
            animation: pulse 2s infinite;
        }
        @keyframes pulse { 0%, 100% { opacity: 1; } 50% { opacity: 0.5; } }
        .version-tag { padding: 4px 12px; background: #F1F5F9; border-radius: 20px; font-size: 12px; color: #64748B; }
        
        .main-container { display: flex; max-width: 1600px; margin: 0 auto; padding: 24px; gap: 24px; }
        .sidebar { width: 320px; flex-shrink: 0; }
        .sidebar-card { background: white; border-radius: 16px; padding: 20px; box-shadow: 0 1px 3px rgba(0,0,0,0.04); }
        .sidebar-title { font-size: 14px; font-weight: 600; color: #1E293B; margin-bottom: 16px; display: flex; align-items: center; gap: 8px; }
        
        .api-group { margin-bottom: 20px; }
        .group-title {
            font-size: 11px; font-weight: 600; color: #6366F1;
            text-transform: uppercase; letter-spacing: 0.5px;
            margin-bottom: 10px; padding: 6px 10px;
            background: rgba(99, 102, 241, 0.1); border-radius: 6px;
        }
        
        .api-item {
            display: flex; align-items: center; gap: 10px;
            padding: 10px 12px; border-radius: 10px;
            cursor: pointer; transition: all 0.2s; margin-bottom: 4px;
        }
        .api-item:hover { background: #F8FAFC; }
        .api-item.active { background: linear-gradient(135deg, rgba(99, 102, 241, 0.1), rgba(139, 92, 246, 0.1)); }
        
        .method { font-size: 10px; font-weight: 700; padding: 3px 8px; border-radius: 4px; min-width: 50px; text-align: center; }
        .method-get { background: #DCFCE7; color: #16A34A; }
        .method-post { background: #DBEAFE; color: #2563EB; }
        .method-put { background: #FEF3C7; color: #D97706; }
        .method-delete { background: #FEE2E2; color: #DC2626; }
        
        .api-path { font-size: 13px; color: #475569; flex: 1; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
        
        .content { flex: 1; min-width: 0; }
        
        .request-card { background: white; border-radius: 16px; padding: 24px; box-shadow: 0 1px 3px rgba(0,0,0,0.04); margin-bottom: 20px; }
        
        .url-bar { display: flex; gap: 12px; margin-bottom: 20px; }
        
        .method-select {
            background: #F8FAFC; border: 1px solid #E2E8F0;
            border-radius: 10px; padding: 12px 16px;
            font-size: 14px; font-weight: 600; color: #1E293B;
            cursor: pointer; min-width: 100px;
        }
        
        .url-input {
            flex: 1; background: #F8FAFC; border: 1px solid #E2E8F0;
            border-radius: 10px; padding: 12px 16px;
            font-size: 14px; font-family: monospace; color: #1E293B;
        }
        .url-input:focus { outline: none; border-color: #6366F1; box-shadow: 0 0 0 3px rgba(99, 102, 241, 0.1); }
        
        .send-btn {
            background: linear-gradient(135deg, #6366F1, #8B5CF6);
            color: white; border: none; border-radius: 10px;
            padding: 12px 28px; font-size: 14px; font-weight: 600;
            cursor: pointer; transition: all 0.2s;
        }
        .send-btn:hover { transform: translateY(-2px); box-shadow: 0 8px 20px rgba(99, 102, 241, 0.3); }
        
        .tabs { display: flex; gap: 4px; margin-bottom: 16px; border-bottom: 1px solid #E2E8F0; padding-bottom: 12px; }
        
        .tab {
            padding: 8px 16px; font-size: 13px;
            background: transparent; color: #64748B;
            border: none; cursor: pointer; border-radius: 8px;
            transition: all 0.2s; font-weight: 500;
        }
        .tab:hover { background: #F1F5F9; color: #1E293B; }
        .tab.active { background: linear-gradient(135deg, rgba(99, 102, 241, 0.1), rgba(139, 92, 246, 0.1)); color: #6366F1; }
        
        .tab-content { display: none; }
        .tab-content.active { display: block; }
        
        .input-row { display: flex; gap: 10px; margin-bottom: 10px; }
        .input-row input {
            flex: 1; background: #F8FAFC; border: 1px solid #E2E8F0;
            border-radius: 8px; padding: 10px 14px; font-size: 13px; color: #1E293B;
        }
        .input-row input:focus { outline: none; border-color: #6366F1; }
        
        .remove-btn {
            background: #FEE2E2; color: #DC2626;
            border: none; border-radius: 6px;
            padding: 8px 12px; cursor: pointer; font-size: 14px;
        }
        
        .add-btn {
            background: #F1F5F9; color: #64748B;
            border: 1px dashed #CBD5E1; border-radius: 8px;
            padding: 10px; cursor: pointer; font-size: 13px;
            width: 100%; transition: all 0.2s;
        }
        .add-btn:hover { background: #E2E8F0; color: #1E293B; }
        
        textarea {
            width: 100%; min-height: 120px;
            background: #F8FAFC; border: 1px solid #E2E8F0;
            border-radius: 10px; padding: 14px;
            font-size: 13px; font-family: monospace; color: #1E293B;
            resize: vertical;
        }
        textarea:focus { outline: none; border-color: #6366F1; }
        
        .response-card { background: white; border-radius: 16px; box-shadow: 0 1px 3px rgba(0,0,0,0.04); overflow: hidden; }
        
        .response-header { padding: 16px 24px; border-bottom: 1px solid #E2E8F0; display: flex; align-items: center; justify-content: space-between; }
        
        .response-info { display: flex; align-items: center; gap: 16px; }
        
        .status-badge { padding: 6px 14px; border-radius: 8px; font-size: 13px; font-weight: 600; }
        .status-2xx { background: #DCFCE7; color: #16A34A; }
        .status-4xx { background: #FEF3C7; color: #D97706; }
        .status-5xx { background: #FEE2E2; color: #DC2626; }
        
        .response-meta { font-size: 12px; color: #64748B; }
        
        .response-body { padding: 20px; max-height: 500px; overflow: auto; background: #1E293B; }
        .response-body pre { color: #E2E8F0; font-size: 12px; font-family: monospace; white-space: pre-wrap; word-break: break-all; line-height: 1.6; }
        
        .json-key { color: #89B4FA; }
        .json-string { color: #A6E3A1; }
        .json-number { color: #FAB387; }
        .json-boolean { color: #CBA6F7; }
        .json-null { color: #6C7086; }
        
        .empty-state { padding: 60px; text-align: center; color: #94A3B8; }
        .empty-state .icon { font-size: 48px; margin-bottom: 16px; }
        
        .token-section { background: #F8FAFC; border-radius: 10px; padding: 14px; margin-bottom: 16px; }
        .token-label { font-size: 12px; color: #64748B; margin-bottom: 8px; }
        .token-input {
            width: 100%; background: white; border: 1px solid #E2E8F0;
            border-radius: 8px; padding: 10px 14px; font-size: 12px;
            font-family: monospace; color: #1E293B;
        }
        
        .quick-actions { display: flex; gap: 8px; margin-bottom: 16px; }
        
        .quick-btn {
            padding: 6px 12px; background: #F1F5F9;
            border: none; border-radius: 6px;
            font-size: 12px; color: #64748B;
            cursor: pointer; transition: all 0.2s;
        }
        .quick-btn:hover { background: #E2E8F0; color: #1E293B; }
        
        ::-webkit-scrollbar { width: 6px; height: 6px; }
        ::-webkit-scrollbar-track { background: #F1F5F9; }
        ::-webkit-scrollbar-thumb { background: #CBD5E1; border-radius: 3px; }
        ::-webkit-scrollbar-thumb:hover { background: #94A3B8; }
    </style>
</head>
<body>
    <div class="header">
        <div class="header-left">
            <div class="logo">🔌</div>
            <h1>API 调试工具</h1>
        </div>
        <div class="header-right">
            <span class="version-tag">{{.Version}}</span>
            <div class="status-dot"></div>
            <span style="font-size: 13px; color: #64748B;">{{.Time}}</span>
        </div>
    </div>
    
    <div class="main-container">
        <div class="sidebar">
            <div class="sidebar-card">
                <div class="sidebar-title">📋 接口列表</div>
                
                <div class="api-group">
                    <div class="group-title">🔓 公共接口</div>
                    <div class="api-item" data-method="GET" data-url="/" data-body="">
                        <span class="method method-get">GET</span>
                        <span class="api-path">/</span>
                    </div>
                    <div class="api-item" data-method="GET" data-url="/health" data-body="">
                        <span class="method method-get">GET</span>
                        <span class="api-path">/health</span>
                    </div>
                    <div class="api-item" data-method="POST" data-url="/api/auth/login" data-body='{"username":"admin","password":"admin123"}'>
                        <span class="method method-post">POST</span>
                        <span class="api-path">/api/auth/login</span>
                    </div>
                    <div class="api-item" data-method="POST" data-url="/api/auth/register" data-body='{"username":"test","email":"test@test.com","password":"123456","full_name":"测试"}'>
                        <span class="method method-post">POST</span>
                        <span class="api-path">/api/auth/register</span>
                    </div>
                </div>
                
                <div class="api-group">
                    <div class="group-title">👤 用户管理</div>
                    <div class="api-item" data-method="GET" data-url="/api/auth/users?skip=0&limit=10" data-body="">
                        <span class="method method-get">GET</span>
                        <span class="api-path">/api/auth/users</span>
                    </div>
                    <div class="api-item" data-method="GET" data-url="/api/auth/users/1" data-body="">
                        <span class="method method-get">GET</span>
                        <span class="api-path">/api/auth/users/:id</span>
                    </div>
                    <div class="api-item" data-method="PUT" data-url="/api/auth/users/1" data-body='{"full_name":"新名字"}'>
                        <span class="method method-put">PUT</span>
                        <span class="api-path">/api/auth/users/:id</span>
                    </div>
                    <div class="api-item" data-method="DELETE" data-url="/api/auth/users/1" data-body="">
                        <span class="method method-delete">DEL</span>
                        <span class="api-path">/api/auth/users/:id</span>
                    </div>
                </div>
                
                <div class="api-group">
                    <div class="group-title">🏢 客户管理</div>
                    <div class="api-item" data-method="GET" data-url="/api/customers" data-body="">
                        <span class="method method-get">GET</span>
                        <span class="api-path">/api/customers</span>
                    </div>
                    <div class="api-item" data-method="POST" data-url="/api/customers" data-body='{"name":"测试客户","code":"C001","type":"customer","contact_person":"张三"}'>
                        <span class="method method-post">POST</span>
                        <span class="api-path">/api/customers</span>
                    </div>
                    <div class="api-item" data-method="GET" data-url="/api/contract-types" data-body="">
                        <span class="method method-get">GET</span>
                        <span class="api-path">/api/contract-types</span>
                    </div>
                </div>
                
                <div class="api-group">
                    <div class="group-title">📄 合同管理</div>
                    <div class="api-item" data-method="GET" data-url="/api/contracts" data-body="">
                        <span class="method method-get">GET</span>
                        <span class="api-path">/api/contracts</span>
                    </div>
                    <div class="api-item" data-method="POST" data-url="/api/contracts" data-body='{"contract_no":"CT2024001","title":"测试合同","customer_id":1,"contract_type_id":1,"amount":100000,"currency":"CNY","status":"draft"}'>
                        <span class="method method-post">POST</span>
                        <span class="api-path">/api/contracts</span>
                    </div>
                    <div class="api-item" data-method="GET" data-url="/api/contracts/1" data-body="">
                        <span class="method method-get">GET</span>
                        <span class="api-path">/api/contracts/:id</span>
                    </div>
                </div>
                
                <div class="api-group">
                    <div class="group-title">📊 数据统计</div>
                    <div class="api-item" data-method="GET" data-url="/api/statistics" data-body="">
                        <span class="method method-get">GET</span>
                        <span class="api-path">/api/statistics</span>
                    </div>
                    <div class="api-item" data-method="GET" data-url="/api/expiring-contracts?days=30" data-body="">
                        <span class="method method-get">GET</span>
                        <span class="api-path">/api/expiring-contracts</span>
                    </div>
                </div>
            </div>
        </div>
        
        <div class="content">
            <div class="request-card">
                <div class="url-bar">
                    <select class="method-select" id="method">
                        <option value="GET">GET</option>
                        <option value="POST">POST</option>
                        <option value="PUT">PUT</option>
                        <option value="DELETE">DELETE</option>
                    </select>
                    <input type="text" class="url-input" id="url" placeholder="输入请求地址..." value="/">
                    <button class="send-btn" onclick="sendRequest()">发送请求</button>
                </div>
                
                <div class="token-section">
                    <div class="token-label">🔐 Auth Token（登录后自动保存）</div>
                    <input type="text" class="token-input" id="token" placeholder="粘贴 JWT Token...">
                </div>
                
                <div class="quick-actions">
                    <button class="quick-btn" onclick="clearAll()">🗑️ 清除</button>
                    <button class="quick-btn" onclick="formatJson()">📝 格式化</button>
                </div>
                
                <div class="tabs">
                    <button class="tab active" onclick="switchTab('params')">Params</button>
                    <button class="tab" onclick="switchTab('headers')">Headers</button>
                    <button class="tab" onclick="switchTab('body')">Body</button>
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
                    <textarea id="body" placeholder='{"key": "value"}'></textarea>
                </div>
            </div>
            
            <div class="response-card">
                <div class="response-header">
                    <div class="response-info" id="response-info" style="display: none;">
                        <span class="status-badge" id="status-badge">200 OK</span>
                        <span class="response-meta" id="response-time"></span>
                        <span class="response-meta" id="response-size"></span>
                    </div>
                    <span style="color: #64748B; font-size: 13px;" id="response-placeholder">响应结果</span>
                </div>
                <div class="response-body" id="response-body">
                    <div class="empty-state">
                        <div class="icon">🚀</div>
                        <div>点击左侧接口或输入URL发送请求</div>
                    </div>
                </div>
            </div>
        </div>
    </div>
    
    <script>
        // 初始化接口点击事件
        document.querySelectorAll('.api-item').forEach(item => {
            item.addEventListener('click', function() {
                const method = this.getAttribute('data-method');
                const url = this.getAttribute('data-url');
                const body = this.getAttribute('data-body');
                
                document.getElementById('method').value = method;
                document.getElementById('url').value = url;
                document.getElementById('body').value = body || '';
                
                document.querySelectorAll('.api-item').forEach(i => i.classList.remove('active'));
                this.classList.add('active');
                
                if (body) switchTab('body');
                else switchTab('params');
            });
        });
        
        function switchTab(tab) {
            document.querySelectorAll('.tab').forEach(t => t.classList.remove('active'));
            document.querySelectorAll('.tab-content').forEach(c => c.classList.remove('active'));
            event.target.classList.add('active');
            document.getElementById('tab-' + tab).classList.add('active');
        }
        
        function addParam() {
            const div = document.createElement('div');
            div.className = 'input-row';
            div.innerHTML = '<input type="text" placeholder="Key" class="param-key"><input type="text" placeholder="Value" class="param-value"><button class="remove-btn" onclick="this.parentElement.remove()">×</button>';
            document.getElementById('params-container').appendChild(div);
        }
        
        function addHeader() {
            const div = document.createElement('div');
            div.className = 'input-row';
            div.innerHTML = '<input type="text" placeholder="Key" class="header-key"><input type="text" placeholder="Value" class="header-value"><button class="remove-btn" onclick="this.parentElement.remove()">×</button>';
            document.getElementById('headers-container').appendChild(div);
        }
        
        function getParams() {
            const params = new URLSearchParams();
            document.querySelectorAll('#params-container .input-row').forEach(row => {
                const key = row.querySelector('.param-key').value;
                const value = row.querySelector('.param-value').value;
                if (key) params.append(key, value);
            });
            return params.toString();
        }
        
        function getHeaders() {
            const headers = {};
            document.querySelectorAll('#headers-container .input-row').forEach(row => {
                const key = row.querySelector('.header-key').value;
                const value = row.querySelector('.header-value').value;
                if (key) headers[key] = value;
            });
            return headers;
        }
        
        async function sendRequest() {
            const method = document.getElementById('method').value;
            let url = document.getElementById('url').value;
            const body = document.getElementById('body').value;
            const token = document.getElementById('token').value;
            
            if (!url.startsWith('http')) {
                url = window.location.origin + url;
            }
            
            const params = getParams();
            if (params) url += (url.includes('?') ? '&' : '?') + params;
            
            const headers = getHeaders();
            if (token) headers['Authorization'] = 'Bearer ' + token;
            
            const options = { method, headers };
            if (body && ['POST', 'PUT'].includes(method)) {
                options.body = body;
            }
            
            const startTime = Date.now();
            
            try {
                const response = await fetch(url, options);
                const time = Date.now() - startTime;
                const size = response.headers.get('content-length') || '-';
                const status = response.status;
                const statusText = response.statusText;
                
                const statusBadge = document.getElementById('status-badge');
                statusBadge.textContent = status + ' ' + statusText;
                statusBadge.className = 'status-badge ' + (status < 300 ? 'status-2xx' : status < 500 ? 'status-4xx' : 'status-5xx');
                document.getElementById('response-info').style.display = 'flex';
                document.getElementById('response-time').textContent = time + 'ms';
                document.getElementById('response-size').textContent = size + ' bytes';
                document.getElementById('response-placeholder').style.display = 'none';
                
                const text = await response.text();
                const bodyEl = document.getElementById('response-body');
                
                try {
                    const json = JSON.parse(text);
                    bodyEl.innerHTML = '<pre>' + syntaxHighlight(json) + '</pre>';
                    
                    if (json.access_token) {
                        document.getElementById('token').value = json.access_token;
                    }
                } catch {
                    bodyEl.innerHTML = '<pre>' + text + '</pre>';
                }
            } catch (error) {
                document.getElementById('response-body').innerHTML = '<pre style="color: #F87171;">Error: ' + error.message + '</pre>';
            }
        }
        
        function syntaxHighlight(json) {
            json = JSON.stringify(json, null, 2);
            return json.replace(/("(\\u[a-zA-Z0-9]{4}|\\[^u]|[^\\"])*"(\s*:)?|\b(true|false|null)\b|-?\d+(?:\.\d*)?(?:[eE][+\-]?\d+)?)/g, function (match) {
                let cls = 'json-number';
                if (/^"/.test(match)) {
                    if (/:$/.test(match)) cls = 'json-key';
                    else cls = 'json-string';
                } else if (/true|false/.test(match)) cls = 'json-boolean';
                else if (/null/.test(match)) cls = 'json-null';
                return '<span class="' + cls + '">' + match + '</span>';
            });
        }
        
        function clearAll() {
            document.getElementById('url').value = '/';
            document.getElementById('body').value = '';
            document.getElementById('response-info').style.display = 'none';
            document.getElementById('response-placeholder').style.display = 'block';
            document.getElementById('response-body').innerHTML = '<div class="empty-state"><div class="icon">🚀</div><div>点击左侧接口或输入URL发送请求</div></div>';
        }
        
        function formatJson() {
            const body = document.getElementById('body');
            try {
                const json = JSON.parse(body.value);
                body.value = JSON.stringify(json, null, 2);
            } catch {
                alert('无效的 JSON 格式');
            }
        }
        
        document.getElementById('url').addEventListener('keypress', e => {
            if (e.key === 'Enter') sendRequest();
        });
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

	r.Use(middleware.SecureHeadersMiddleware())
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.RateLimitMiddleware())
	r.Use(middleware.RequestSizeLimitMiddleware(10 << 20))

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
		auth.GET("/users", middleware.AuthMiddleware(), authHandler.GetUsers)
		auth.GET("/users/:user_id", middleware.AuthMiddleware(), authHandler.GetUserByID)
		auth.PUT("/users/:user_id", middleware.AuthMiddleware(), authHandler.UpdateUser)
		auth.DELETE("/users/:user_id", middleware.AuthMiddleware(), authHandler.DeleteUser)
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

		api.GET("/audit-logs", auditHandler.GetAuditLogs)
		api.DELETE("/audit-logs/:id", auditHandler.DeleteAuditLog)
		api.POST("/audit-logs/batch-delete", auditHandler.DeleteAuditLogs)
		api.GET("/audit-logs/export", auditHandler.ExportAuditLogs)
	}

	_ = approvalHandler
	_ = contractHandler

	addr := ":8000"
	fmt.Printf("API 调试页面: http://localhost%s\n", addr)
	fmt.Printf("Swagger 文档: http://localhost%s/docs/swagger.html\n", addr)
	r.Run(addr)
}
