package middleware

import (
	"encoding/json"
	"html"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type RateLimiter struct {
	visitors map[string]*visitor
	mu       sync.RWMutex
	rate     int
	burst    int
	expires  time.Duration
}

type visitor struct {
	tokens    int
	lastVisit time.Time
}

func NewRateLimiter(rate int, burst int, expires time.Duration) *RateLimiter {
	rl := &RateLimiter{
		visitors: make(map[string]*visitor),
		rate:     rate,
		burst:    burst,
		expires:  expires,
	}
	go rl.cleanup()
	return rl
}

func (rl *RateLimiter) getVisitor(ip string) *visitor {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	v, exists := rl.visitors[ip]
	if !exists {
		rl.visitors[ip] = &visitor{tokens: rl.burst, lastVisit: time.Now()}
		return rl.visitors[ip]
	}

	now := time.Now()
	elapsed := now.Sub(v.lastVisit)
	v.tokens = min(rl.burst, v.tokens+int(elapsed.Seconds()*float64(rl.rate)))
	v.lastVisit = now

	return v
}

func (rl *RateLimiter) Allow(ip string) bool {
	v := rl.getVisitor(ip)
	if v.tokens < 1 {
		return false
	}
	v.tokens--
	return true
}

func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(rl.expires)
	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		for ip, v := range rl.visitors {
			if now.Sub(v.lastVisit) > rl.expires {
				delete(rl.visitors, ip)
			}
		}
		rl.mu.Unlock()
	}
}

var (
	rateLimiter     *RateLimiter
	rateLimiterOnce sync.Once
)

func GetRateLimiter() *RateLimiter {
	rateLimiterOnce.Do(func() {
		rateLimiter = NewRateLimiter(60, 100, time.Hour)
	})
	return rateLimiter
}

func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !GetRateLimiter().Allow(c.ClientIP()) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "请求过于频繁，请稍后再试",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

func SecureHeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")

		// 为API调试页面设置更宽松的CSP策略
		if c.Request.URL.Path == "/" {
			// 允许内联样式和脚本，用于API调试页面
			c.Header("Content-Security-Policy", "default-src 'self'; style-src 'self' 'unsafe-inline'; script-src 'self' 'unsafe-inline'")
		} else {
			// 其他页面使用严格策略
			c.Header("Content-Security-Policy", "default-src 'self'")
		}

		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Header("Permissions-Policy", "geolocation=(), microphone=(), camera=()")
		c.Next()
	}
}

func CORSMiddleware() gin.HandlerFunc {
	allowedOrigins := []string{
		"http://localhost:3000",
		"http://localhost:3001",
		"http://127.0.0.1:3000",
		"http://127.0.0.1:3001",
	}

	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")

		allowed := false
		for _, o := range allowedOrigins {
			if origin == o || strings.HasPrefix(origin, "http://localhost:") {
				allowed = true
				break
			}
		}

		if allowed {
			c.Header("Access-Control-Allow-Origin", origin)
		}

		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-Requested-With")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Max-Age", "3600")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func RequestSizeLimitMiddleware(maxSize int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.ContentLength > maxSize {
			c.JSON(http.StatusRequestEntityTooLarge, gin.H{
				"error": "请求体过大",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// XssSanitize 对字符串进行HTML转义，防止XSS攻击
func XssSanitize(s string) string {
	return html.EscapeString(s)
}

// SanitizeMap 递归清理map中的所有字符串值
func SanitizeMap(data map[string]interface{}) map[string]interface{} {
	for key, value := range data {
		switch v := value.(type) {
		case string:
			data[key] = XssSanitize(v)
		case map[string]interface{}:
			data[key] = SanitizeMap(v)
		case []interface{}:
			data[key] = SanitizeSlice(v)
		}
	}
	return data
}

// SanitizeSlice 递归清理slice中的所有字符串值
func SanitizeSlice(slice []interface{}) []interface{} {
	for i, value := range slice {
		switch v := value.(type) {
		case string:
			slice[i] = XssSanitize(v)
		case map[string]interface{}:
			slice[i] = SanitizeMap(v)
		case []interface{}:
			slice[i] = SanitizeSlice(v)
		}
	}
	return slice
}

// XssProtectionMiddleware XSS防护中间件，对JSON响应进行清理
func XssProtectionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 创建自定义响应写入器来捕获响应体
		writer := &XssSafeResponseWriter{ResponseWriter: c.Writer, body: &strings.Builder{}}
		c.Writer = writer

		c.Next()

		// 恢复原始的响应写入器
		c.Writer = writer.ResponseWriter

		// 检查Content-Type
		contentType := writer.Header().Get("Content-Type")
		if strings.Contains(contentType, "application/json") {
			// 解析JSON，清理字符串值，然后重新序列化
			var data interface{}
			bodyBytes := []byte(writer.body.String())
			if err := json.Unmarshal(bodyBytes, &data); err == nil {
				cleanedData := sanitizeJSONData(data)
				if cleanedBytes, err := json.Marshal(cleanedData); err == nil {
					// 设置正确的内容长度并写入清理后的数据
					writer.Header().Set("Content-Length", strconv.Itoa(len(cleanedBytes)))
					writer.ResponseWriter.Write(cleanedBytes)
					return
				}
			}
		}

		// 对于非JSON响应或JSON清理失败的情况，直接写入原始内容
		if writer.body.Len() > 0 {
			writer.Header().Set("Content-Length", strconv.Itoa(writer.body.Len()))
			writer.ResponseWriter.Write([]byte(writer.body.String()))
		}
	}
}

// XssSafeResponseWriter 自定义响应写入器，用于捕获响应体
type XssSafeResponseWriter struct {
	gin.ResponseWriter
	body *strings.Builder
}

func (w *XssSafeResponseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return len(b), nil
}

// sanitizeJSONData 递归清理JSON数据中的字符串值
func sanitizeJSONData(data interface{}) interface{} {
	switch v := data.(type) {
	case string:
		return XssSanitize(v)
	case map[string]interface{}:
		for key, value := range v {
			v[key] = sanitizeJSONData(value)
		}
		return v
	case []interface{}:
		for i, value := range v {
			v[i] = sanitizeJSONData(value)
		}
		return v
	default:
		return data
	}
}
