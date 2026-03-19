package handlers

import (
	"contract-manage/middleware"
	"contract-manage/services"
	"net/http"
	"regexp"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	userService *services.UserService
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{
		userService: services.NewUserService(),
	}
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type TokenResponse struct {
	AccessToken string    `json:"access_token"`
	TokenType   string    `json:"token_type"`
	UserInfo    *UserInfo `json:"user_info"`
}

type UserInfo struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	FullName string `json:"full_name"`
	Role     string `json:"role"`
}

var (
	emailRegex    = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_]{3,20}$`)
)

func sanitizeInput(input string) string {
	input = regexp.MustCompile(`[<>\"'%;()&+]`).ReplaceAllString(input, "")
	return input
}

func validateUsername(username string) error {
	if username == "" {
		return usernameError("用户名不能为空")
	}
	if len(username) < 3 || len(username) > 20 {
		return usernameError("用户名长度必须在3-20个字符之间")
	}
	if !usernameRegex.MatchString(username) {
		return usernameError("用户名只能包含字母、数字、下划线和中文")
	}
	return nil
}

func validateEmail(email string) error {
	if email == "" {
		return nil
	}
	if !emailRegex.MatchString(email) {
		return usernameError("邮箱格式不正确")
	}
	return nil
}

func validatePassword(password string) error {
	if password == "" {
		return usernameError("密码不能为空")
	}
	if len(password) < 6 {
		return usernameError("密码长度至少6位")
	}
	if len(password) > 50 {
		return usernameError("密码长度不能超过50位")
	}
	return nil
}

type validationError struct {
	message string
}

func (e validationError) Error() string {
	return e.message
}

func usernameError(msg string) error {
	return validationError{message: msg}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var input services.UserCreateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求数据格式不正确"})
		return
	}

	input.Username = sanitizeInput(input.Username)
	input.Email = sanitizeInput(input.Email)
	input.FullName = sanitizeInput(input.FullName)
	input.Department = sanitizeInput(input.Department)
	input.Phone = sanitizeInput(input.Phone)

	if err := validateUsername(input.Username); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validateEmail(input.Email); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validatePassword(input.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Role == "" {
		input.Role = "user"
	}

	user, err := h.userService.CreateUser(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
		"message":  "注册成功",
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求数据格式不正确"})
		return
	}

	req.Username = sanitizeInput(req.Username)

	if err := validateUsername(req.Username); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validatePassword(req.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userService.AuthenticateUser(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	if !user.IsActive {
		c.JSON(http.StatusForbidden, gin.H{"error": "账号已被禁用"})
		return
	}

	token, err := middleware.GenerateTokenWithUserIDAndRole(user.ID, user.Username, string(user.Role))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "登录失败，请稍后重试"})
		return
	}

	c.JSON(http.StatusOK, TokenResponse{
		AccessToken: token,
		TokenType:   "bearer",
		UserInfo: &UserInfo{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			FullName: user.FullName,
			Role:     string(user.Role),
		},
	})
}

func (h *AuthHandler) GetUsers(c *gin.Context) {
	skip := 0
	limit := 20

	if s := c.Query("skip"); s != "" {
		if parsed, err := strconv.Atoi(s); err == nil && parsed >= 0 {
			skip = parsed
		}
	}

	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	users, err := h.userService.GetUsers(skip, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户列表失败"})
		return
	}

	c.JSON(http.StatusOK, users)
}

func (h *AuthHandler) GetUserByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("user_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户ID格式不正确"})
		return
	}

	if id < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户ID无效"})
		return
	}

	user, err := h.userService.GetUserByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *AuthHandler) UpdateUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("user_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户ID格式不正确"})
		return
	}

	if id < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户ID无效"})
		return
	}

	var input services.UserUpdateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求数据格式不正确"})
		return
	}

	if input.Email != "" {
		input.Email = sanitizeInput(input.Email)
		if err := validateEmail(input.Email); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	if input.FullName != "" {
		input.FullName = sanitizeInput(input.FullName)
	}

	if input.Department != "" {
		input.Department = sanitizeInput(input.Department)
	}

	if input.Phone != "" {
		input.Phone = sanitizeInput(input.Phone)
	}

	if input.Role != "" {
		currentUserID, _ := middleware.GetCurrentUserID(c)
		currentRole, _ := c.Get("role")

		if currentRole != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "无权限修改用户角色"})
			return
		}

		if uint(id) == currentUserID && input.Role != "admin" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "不能撤销自己的管理员权限"})
			return
		}
	}

	user, err := h.userService.UpdateUser(uint(id), input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *AuthHandler) DeleteUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("user_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户ID格式不正确"})
		return
	}

	if id < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户ID无效"})
		return
	}

	currentUserID, _ := middleware.GetCurrentUserID(c)
	if uint(id) == currentUserID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "不能删除自己的账号"})
		return
	}

	if err := h.userService.DeleteUser(uint(id)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
