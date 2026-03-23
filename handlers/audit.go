package handlers

import (
	"contract-manage/middleware"
	"contract-manage/models"
	"contract-manage/services"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuditHandler struct {
	auditService *services.AuditService
}

func NewAuditHandler() *AuditHandler {
	return &AuditHandler{
		auditService: services.NewAuditService(),
	}
}

func GetAuditService() *services.AuditService {
	return services.NewAuditService()
}

func (h *AuditHandler) GetAuditLogs(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if err != nil || pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	username := c.Query("username")
	if len(username) > 50 {
		username = username[:50]
	}

	action := c.Query("action")
	if len(action) > 50 {
		action = action[:50]
	}

	module := c.Query("module")
	if len(module) > 50 {
		module = module[:50]
	}

	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	// 简单的日期格式验证 (YYYY-MM-DD)
	dateRegex := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
	if startDate != "" && !dateRegex.MatchString(startDate) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "start_date格式错误，请使用YYYY-MM-DD格式"})
		return
	}
	if endDate != "" && !dateRegex.MatchString(endDate) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "end_date格式错误，请使用YYYY-MM-DD格式"})
		return
	}

	logs, total, err := h.auditService.GetAuditLogs(page, pageSize, username, action, module, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"logs":  logs,
		"total": total,
		"page":  page,
		"size":  pageSize,
	})
}

func (h *AuditHandler) DeleteAuditLog(c *gin.Context) {
	role, _ := middleware.GetCurrentUserRole(c)
	if role != "admin" && role != "audit_admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权限删除审计日志，仅审计管理员可操作"})
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := h.auditService.DeleteAuditLog(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

func (h *AuditHandler) DeleteAuditLogs(c *gin.Context) {
	role, _ := middleware.GetCurrentUserRole(c)
	if role != "admin" && role != "audit_admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权限删除审计日志，仅审计管理员可操作"})
		return
	}

	var input struct {
		IDs []uint `json:"ids"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if err := h.auditService.DeleteAuditLogs(input.IDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "批量删除成功"})
}

func (h *AuditHandler) ExportAuditLogs(c *gin.Context) {
	username := c.Query("username")
	if len(username) > 50 {
		username = username[:50]
	}

	action := c.Query("action")
	if len(action) > 50 {
		action = action[:50]
	}

	module := c.Query("module")
	if len(module) > 50 {
		module = module[:50]
	}

	startDate := c.DefaultQuery("start_date", "")
	endDate := c.DefaultQuery("end_date", "")

	// 简单的日期格式验证 (YYYY-MM-DD)
	dateRegex := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
	if startDate != "" && !dateRegex.MatchString(startDate) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "start_date格式错误，请使用YYYY-MM-DD格式"})
		return
	}
	if endDate != "" && !dateRegex.MatchString(endDate) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "end_date格式错误，请使用YYYY-MM-DD格式"})
		return
	}

	logs, _, err := h.auditService.GetAuditLogs(1, 10000, username, action, module, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, logs)
}

func AuditLogMiddleware(auditService *services.AuditService) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if c.Request.Method == "OPTIONS" {
			return
		}

		if strings.HasPrefix(c.Request.URL.Path, "/api/auth/login") ||
			strings.HasPrefix(c.Request.URL.Path, "/api/auth/register") {
			return
		}

		userID, _ := middleware.GetCurrentUserID(c)
		username, _ := c.Get("username")

		if userID == 0 {
			return
		}

		action := c.Request.Method + " " + c.Request.URL.Path
		module := getModuleFromPath(c.Request.URL.Path)

		clientIP := c.ClientIP()
		userAgent := c.Request.UserAgent()

		log := models.AuditLog{
			UserID:    userID,
			Username:  username.(string),
			Action:    action,
			Module:    module,
			Method:    c.Request.Method,
			Path:      c.Request.URL.Path,
			IPAddress: clientIP,
			UserAgent: userAgent,
		}

		statusCode := c.Writer.Status()
		log.StatusCode = statusCode

		go auditService.CreateAuditLog(log)
	}
}

func getModuleFromPath(path string) string {
	if strings.Contains(path, "/auth/") {
		return "auth"
	}
	if strings.Contains(path, "/contracts") {
		return "contract"
	}
	if strings.Contains(path, "/customers") {
		return "customer"
	}
	if strings.Contains(path, "/approvals") || strings.Contains(path, "/pending") {
		return "approval"
	}
	if strings.Contains(path, "/reminders") {
		return "reminder"
	}
	if strings.Contains(path, "/users") {
		return "user"
	}
	if strings.Contains(path, "/statistics") {
		return "statistics"
	}
	return "other"
}
