package middleware

import (
	"net/http"
	"strings"

	"contract-manage/services"

	"github.com/gin-gonic/gin"
)

var PreviewService *services.PreviewService

func init() {
	PreviewService = services.NewPreviewService()
}

// PreviewAuthMiddleware 预览令牌认证中间件
// 用于iframe文档预览，使用短效令牌而非JWT
func PreviewAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从 query 参数获取预览令牌
		token := c.Query("preview_token")

		if token == "" {
			// 也尝试从 header 获取
			authHeader := c.GetHeader("Authorization")
			if authHeader != "" {
				parts := strings.SplitN(authHeader, " ", 2)
				if len(parts) == 2 && parts[0] == "Preview" {
					token = strings.TrimSpace(parts[1])
				}
			}
		}

		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Preview token required"})
			c.Abort()
			return
		}

		// 验证预览令牌
		previewToken, valid := PreviewService.ValidateToken(token)
		if !valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired preview token"})
			c.Abort()
			return
		}

		// 设置用户信息到context
		c.Set("user_id", previewToken.UserID)
		c.Set("document_id", previewToken.DocumentID)
		c.Set("preview_token", token)

		c.Next()
	}
}
