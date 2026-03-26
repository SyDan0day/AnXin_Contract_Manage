package handlers

import (
	"archive/zip"
	"bytes"
	"contract-manage/config"
	"contract-manage/middleware"
	"contract-manage/models"
	"contract-manage/services"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type ContractHandler struct {
	contractService *services.ContractService
}

func NewContractHandler() *ContractHandler {
	return &ContractHandler{
		contractService: services.NewContractService(),
	}
}

func (h *ContractHandler) GetContracts(c *gin.Context) {
	skip, err := strconv.Atoi(c.DefaultQuery("skip", "0"))
	if err != nil || skip < 0 {
		skip = 0
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "100"))
	if err != nil || limit < 0 {
		limit = 100
	}
	if limit > 1000 {
		limit = 1000
	}

	var customerID uint64
	if customerStr := c.Query("customer_id"); customerStr != "" {
		if id, err := strconv.ParseUint(customerStr, 10, 32); err == nil {
			customerID = id
		}
	}

	var contractTypeID uint64
	if typeStr := c.Query("contract_type_id"); typeStr != "" {
		if id, err := strconv.ParseUint(typeStr, 10, 32); err == nil {
			contractTypeID = id
		}
	}

	status := c.Query("status")
	if status != "" && !isValidContractStatus(status) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status value"})
		return
	}

	keyword := c.Query("keyword")
	if len(keyword) > 200 {
		keyword = keyword[:200]
	}
	// Escape LIKE wildcards to prevent LIKE injection
	keyword = strings.ReplaceAll(keyword, "%", "\\%")
	keyword = strings.ReplaceAll(keyword, "_", "\\_")

	// 获取当前用户信息用于权限过滤
	userID, _ := middleware.GetCurrentUserID(c)
	role, _ := middleware.GetCurrentUserRole(c)

	// 使用新方法，支持角色过滤
	contracts, total, err := h.contractService.GetContracts(services.GetContractsParams{
		Skip:           skip,
		Limit:          limit,
		CustomerID:     uint(customerID),
		ContractTypeID: uint(contractTypeID),
		Status:         status,
		Keyword:        keyword,
		UserID:         userID,
		Role:           role,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 为每个合同添加最新的拒绝信息和审批信息
	contractsWithRejection := make([]gin.H, 0, len(contracts))
	for _, contract := range contracts {
		contractData := gin.H{
			"id":               contract.ID,
			"contract_no":      contract.ContractNo,
			"title":            contract.Title,
			"customer_id":      contract.CustomerID,
			"contract_type_id": contract.ContractTypeID,
			"amount":           contract.Amount,
			"currency":         contract.Currency,
			"status":           contract.Status,
			"sign_date":        contract.SignDate,
			"start_date":       contract.StartDate,
			"end_date":         contract.EndDate,
			"payment_terms":    contract.PaymentTerms,
			"content":          contract.Content,
			"notes":            contract.Notes,
			"creator_id":       contract.CreatorID,
			"created_at":       contract.CreatedAt,
			"updated_at":       contract.UpdatedAt,
			"customer":         contract.Customer,
			"creator":          contract.Creator,
			"contract_type":    contract.ContractType,
		}

		// 获取最新的拒绝信息
		var rejectionInfo map[string]interface{}
		models.DB.Table("workflow_approvals as wa").
			Select("wa.comment, wa.approved_at, u.full_name as approver_name, wa.approver_role").
			Joins("LEFT JOIN users u ON u.id = wa.approver_id").
			Where("wa.contract_id = ? AND wa.status = 'rejected'", contract.ID).
			Order("wa.approved_at DESC").
			Limit(1).
			Scan(&rejectionInfo)

		if rejectionInfo != nil && len(rejectionInfo) > 0 {
			contractData["rejection_info"] = rejectionInfo
		}

		// 如果是审批中状态，获取当前审批级别
		if contract.Status == "pending" {
			var workflowInfo struct {
				CurrentLevel int    `json:"current_level"`
				MaxLevel     int    `json:"max_level"`
				ApproverRole string `json:"approver_role"`
			}
			models.DB.Table("approval_workflows").
				Select("current_level, max_level").
				Where("contract_id = ? AND status = ?", contract.ID, "pending").
				Order("created_at DESC").
				Limit(1).
				Scan(&workflowInfo)

			if workflowInfo.MaxLevel > 0 {
				contractData["current_approval_level"] = workflowInfo.CurrentLevel
				contractData["max_approval_level"] = workflowInfo.MaxLevel

				var currentApprover string
				switch workflowInfo.CurrentLevel {
				case 1:
					currentApprover = "销售负责人"
				case 2:
					currentApprover = "技术负责人"
				case 3:
					currentApprover = "财务负责人"
				default:
					currentApprover = "待审批"
				}
				contractData["current_approver"] = currentApprover
			}
		}

		contractsWithRejection = append(contractsWithRejection, contractData)
	}

	c.JSON(http.StatusOK, gin.H{
		"items": contractsWithRejection,
		"total": total,
	})
}

func isValidContractStatus(status string) bool {
	validStatuses := []string{
		"draft",
		"pending",
		"active",
		"completed",
		"terminated",
		"archived",
	}
	for _, s := range validStatuses {
		if status == s {
			return true
		}
	}
	return false
}

func (h *ContractHandler) GetContractByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("contract_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid contract ID"})
		return
	}

	// 获取当前用户信息
	userID, _ := middleware.GetCurrentUserID(c)
	role, _ := middleware.GetCurrentUserRole(c)

	// 使用带权限检查的方法
	contract, err := h.contractService.GetContractByIDWithAuth(uint(id), userID, role)
	if err != nil {
		if strings.Contains(err.Error(), "无权限") {
			c.JSON(http.StatusForbidden, gin.H{"error": "无权限查看此合同"})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "Contract not found"})
		return
	}

	c.JSON(http.StatusOK, contract)
}

func (h *ContractHandler) CreateContract(c *gin.Context) {
	var input services.ContractCreateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	contract, err := h.contractService.CreateContract(input, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, contract)
}

func (h *ContractHandler) UpdateContract(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("contract_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid contract ID"})
		return
	}

	var input services.ContractUpdateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	contract, err := h.contractService.UpdateContract(uint(id), input)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Contract not found"})
		return
	}

	c.JSON(http.StatusOK, contract)
}

func (h *ContractHandler) DeleteContract(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("contract_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid contract ID"})
		return
	}

	if err := h.contractService.DeleteContract(uint(id)); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Contract not found"})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *ContractHandler) GetContractDocuments(c *gin.Context) {
	contractID, err := strconv.ParseUint(c.Param("contract_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid contract ID"})
		return
	}

	documents, err := h.contractService.GetDocuments(uint(contractID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, documents)
}

func (h *ContractHandler) CreateContractDocument(c *gin.Context) {
	contractID, err := strconv.ParseUint(c.Param("contract_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid contract ID"})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请选择要上传的文件"})
		return
	}

	// 文件类型白名单（移除HTML/HTM以防止XSS）
	allowedExtensions := map[string]bool{
		"pdf": true, "doc": true, "docx": true, "xls": true, "xlsx": true,
		"jpg": true, "jpeg": true, "png": true, "gif": true, "bmp": true, "webp": true,
		"txt": true,
	}

	// MIME类型白名单
	allowedMIMETypes := map[string]bool{
		"application/pdf":    true,
		"application/msword": true,
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
		"application/vnd.ms-excel": true,
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet": true,
		"image/jpeg": true,
		"image/png":  true,
		"image/gif":  true,
		"image/bmp":  true,
		"image/webp": true,
		"text/plain": true,
	}

	// 获取文件扩展名（转为小写）
	filename := filepath.Base(file.Filename)
	ext := strings.ToLower(filepath.Ext(filename))
	if ext != "" {
		ext = ext[1:] // 移除开头的点
	}

	// 验证文件类型
	if !allowedExtensions[ext] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "不支持的文件类型: " + ext})
		return
	}

	// 检查双重扩展名，防止恶意文件上传（如 test.php.jpg）
	if strings.Contains(strings.ToLower(filename), ".php.") ||
		strings.Contains(strings.ToLower(filename), ".jsp.") ||
		strings.Contains(strings.ToLower(filename), ".asp.") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "检测到恶意文件名格式"})
		return
	}

	// 清理文件名，移除路径信息，防止路径遍历（filename already set above）

	uploadDir := config.AppConfig.UploadDir
	if uploadDir == "" {
		uploadDir = "uploads"
	}

	// 构建安全的文件路径
	filePath := filepath.Join(uploadDir, strconv.FormatUint(contractID, 10), filename)

	// 规范化路径，防止路径遍历
	cleanFilePath := filepath.Clean(filePath)
	if strings.Contains(cleanFilePath, "..") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file path"})
		return
	}

	// 读取文件内容进行安全检查
	fileContent, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法读取文件"})
		return
	}
	defer fileContent.Close()

	// 读取文件前4KB进行内容检查
	buffer := make([]byte, 4096)
	bytesRead, err := fileContent.Read(buffer)
	if err != nil && err != io.EOF {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "读取文件内容失败"})
		return
	}
	fileContent.Seek(0, io.SeekStart) // 重置读取位置

	// 检查文件内容中是否包含恶意代码
	if containsMaliciousContent(buffer[:bytesRead], ext) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文件内容包含恶意代码"})
		return
	}

	// 使用http.DetectContentType检测文件真实MIME类型
	detectedMIME := http.DetectContentType(buffer[:bytesRead])
	if !allowedMIMETypes[detectedMIME] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文件内容类型不匹配: " + detectedMIME})
		return
	}

	if err := os.MkdirAll(filepath.Dir(cleanFilePath), 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建目录失败"})
		return
	}

	if err := c.SaveUploadedFile(file, cleanFilePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存文件失败"})
		return
	}

	input := services.DocumentCreateInput{
		ContractID: uint(contractID),
		Name:       filename,
		FilePath:   "/" + cleanFilePath,
		FileSize:   int(file.Size),
		FileType:   ext,
	}

	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	document, err := h.contractService.CreateDocument(input, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, document)
}

// GeneratePreviewToken 生成预览令牌
// POST /api/documents/:document_id/preview-token
func (h *ContractHandler) GeneratePreviewToken(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("document_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid document ID"})
		return
	}

	// 验证文档是否存在
	_, err = h.contractService.GetDocumentByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Document not found"})
		return
	}

	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// 生成预览令牌
	token, err := middleware.PreviewService.GenerateToken(uint(id), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate preview token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"preview_token": token,
		"expires_in":    300, // 5分钟，单位秒
	})
}

func (h *ContractHandler) PreviewDocument(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("document_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid document ID"})
		return
	}

	document, err := h.contractService.GetDocumentByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Document not found"})
		return
	}

	// 清理文件路径，防止路径遍历
	cleanPath := filepath.Clean(document.FilePath)
	if strings.Contains(cleanPath, "..") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file path"})
		return
	}

	// 构建绝对文件路径
	absFilePath := filepath.Join(".", cleanPath)

	// 检查文件是否存在
	if _, err := os.Stat(absFilePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "合同文件不存在，请联系管理员上传", "code": "FILE_NOT_FOUND"})
		return
	}

	// 根据文件类型返回不同的内容
	fileExt := strings.ToLower(filepath.Ext(document.Name))

	// Word 文档 (.docx) 返回纯文本内容
	if fileExt == ".docx" {
		text, err := extractTextFromDocx(absFilePath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "无法读取文档内容: " + err.Error()})
			return
		}

		// 返回纯文本内容
		c.JSON(http.StatusOK, gin.H{
			"document_id": document.ID,
			"file_name":   document.Name,
			"file_type":   document.FileType,
			"file_size":   document.FileSize,
			"created_at":  document.CreatedAt,
			"content":     text,
		})
		return
	}

	switch fileExt {
	case ".pdf":
		// PDF 文件直接返回，允许在iframe中显示
		c.Header("Content-Type", "application/pdf")
		c.Header("X-Frame-Options", "SAMEORIGIN")
		c.Header("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", document.Name))
		c.File(absFilePath)
	case ".doc":
		// Word 文档返回文件内容
		c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.wordprocessingml.document")
		c.Header("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", document.Name))
		c.File(absFilePath)
	case ".xls", ".xlsx":
		// Excel 文件返回文件内容
		c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		c.Header("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", document.Name))
		c.File(absFilePath)
	case ".jpg", ".jpeg", ".png", ".gif", ".bmp", ".webp":
		// 图片文件直接返回
		c.Header("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", document.Name))
		c.File(absFilePath)
	case ".txt":
		// 文本文件返回内容
		content, err := os.ReadFile(absFilePath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "无法读取文件内容"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"document_id": document.ID,
			"file_name":   document.Name,
			"file_type":   document.FileType,
			"file_size":   document.FileSize,
			"created_at":  document.CreatedAt,
			"content":     string(content),
		})
		return
	case ".html", ".htm":
		// HTML 文件返回内容
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.File(absFilePath)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "不支持的文件类型: " + fileExt})
	}
}

// convertWordToHTML 将 Word 文档转换为 HTML
func (h *ContractHandler) convertWordToHTML(docxPath, filePath string) (string, error) {
	// 使用 mammoth 库转换 Word 到 HTML
	// 这里需要调用 Python 脚本
	// 由于 Go 调用 Python 比较复杂，我们可以使用 exec 执行 mammoth 命令行工具
	// 或者使用 Go 库

	// 简单实现：返回提示信息，实际部署时需要安装 mammoth 并调用
	return fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<head>
			<meta charset="UTF-8">
			<title>文档预览</title>
			<style>
				body { font-family: Arial, sans-serif; margin: 20px; }
				.info { background: #f0f0f0; padding: 20px; border-radius: 5px; }
			</style>
		</head>
		<body>
			<div class="info">
				<h3>Word 文档预览</h3>
				<p>文件: %s</p>
				<p>Word 文档需要下载后查看完整内容。</p>
				<p><a href="%s" download>点击下载文件</a></p>
			</div>
		</body>
		</html>
	`, filepath.Base(docxPath), filePath), nil
}

func (h *ContractHandler) DeleteDocument(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("document_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid document ID"})
		return
	}

	if err := h.contractService.DeleteDocument(uint(id)); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Document not found"})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *ContractHandler) GetContractLifecycle(c *gin.Context) {
	contractID, err := strconv.ParseUint(c.Param("contract_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid contract ID"})
		return
	}

	events, err := h.contractService.GetLifecycleEvents(uint(contractID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, events)
}

func (h *ContractHandler) UpdateContractStatus(c *gin.Context) {
	contractID, err := strconv.ParseUint(c.Param("contract_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid contract ID"})
		return
	}

	var input struct {
		Status      string `json:"status" binding:"required"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	contract, err := h.contractService.UpdateContractStatus(uint(contractID), input.Status, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, contract)
}

func (h *ContractHandler) ArchiveContract(c *gin.Context) {
	contractID, err := strconv.ParseUint(c.Param("contract_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid contract ID"})
		return
	}

	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	contract, err := h.contractService.ArchiveContract(uint(contractID), userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, contract)
}

func (h *ContractHandler) UploadContractTemplate(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请选择要上传的文件"})
		return
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(header.Filename))
	if ext != ".docx" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "仅支持 .docx 格式文件"})
		return
	}

	if header.Size > 10*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文件大小不能超过 10MB"})
		return
	}

	uploadDir := "./uploads/contracts"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建上传目录失败"})
		return
	}

	filename := fmt.Sprintf("%d_%s", time.Now().Unix(), header.Filename)
	filePath := filepath.Join(uploadDir, filename)

	out, err := os.Create(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存文件失败"})
		return
	}
	defer out.Close()

	if _, err := io.Copy(out, file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存文件失败"})
		return
	}

	parsedData, parseErr := parseDocxFile(filePath)

	if parseErr != nil {
		c.JSON(200, gin.H{
			"success":  true,
			"message":  "文件上传成功，但解析失败: " + parseErr.Error(),
			"file_url": "/uploads/contracts/" + filename,
			"data":     nil,
		})
		return
	}

	c.JSON(200, gin.H{
		"success":  true,
		"message":  "文件上传并解析成功",
		"file_url": "/uploads/contracts/" + filename,
		"data":     parsedData,
	})
}

func parseDocxFile(filePath string) (map[string]interface{}, error) {
	text, err := extractTextFromDocx(filePath)
	if err != nil {
		return nil, err
	}

	if text == "" {
		return nil, fmt.Errorf("无法读取文档内容")
	}

	data := extractContractData(text)

	return data, nil
}

func extractTextFromDocx(filePath string) (string, error) {
	r, err := zip.OpenReader(filePath)
	if err != nil {
		return "", err
	}
	defer r.Close()

	var text strings.Builder

	for _, file := range r.File {
		if file.Name == "word/document.xml" {
			rc, err := file.Open()
			if err != nil {
				return "", err
			}
			defer rc.Close()

			content, err := io.ReadAll(rc)
			if err != nil {
				return "", err
			}

			re := regexp.MustCompile(`<w:t[^>]*>([^<]*)</w:t>`)
			matches := re.FindAllStringSubmatch(string(content), -1)
			for _, match := range matches {
				if len(match) > 1 {
					text.WriteString(match[1])
					text.WriteString(" ")
				}
			}
			break
		}
	}

	return text.String(), nil
}

func contentToString(content interface{}) string {
	switch v := content.(type) {
	case string:
		return v
	case []byte:
		return string(v)
	default:
		return fmt.Sprintf("%v", v)
	}
}

func extractContractData(text string) map[string]interface{} {
	data := make(map[string]interface{})

	patterns := map[string]string{
		"contract_no":   `合同编号[：:]\s*([A-Z0-9\-]+)(?:\s|$|\n)`,
		"title":         `合同名称[：:]\s*([^\n]+?)\s*(?:\n|$)`,
		"customer_name": `甲方[（(]客户[）)][：:]\s*([^\n]+?)\s*(?:\n|$)`,
		"amount":        `合同金额[：:]\s*([\d,]+\.?\d*)\s*(?:元|万)?(?:\s|$|\n)`,
		"sign_date":     `签订日期[：:]\s*(\d{4}[-/年]\d{1,2}[-/月]\d{1,2}[日]?)(?:\s|$|\n)`,
		"start_date":    `开始日期[：:]\s*(\d{4}[-/年]\d{1,2}[-/月]\d{1,2}[日]?)(?:\s|$|\n)`,
		"end_date":      `结束日期[：:]\s*(\d{4}[-/年]\d{1,2}[-/月]\d{1,2}[日]?)(?:\s|$|\n)`,
		"contract_type": `合同类型[：:]\s*([^\n]+?)\s*(?:\n|$)`,
	}

	for key, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(text)
		if len(matches) > 1 {
			value := strings.TrimSpace(matches[1])
			value = strings.ReplaceAll(value, "年", "-")
			value = strings.ReplaceAll(value, "月", "-")
			value = strings.ReplaceAll(value, "日", "")

			switch key {
			case "amount":
				value = strings.ReplaceAll(value, ",", "")
				if num, err := strconv.ParseFloat(value, 64); err == nil {
					data[key] = num
				}
			case "sign_date", "start_date", "end_date":
				if isValidDate(value) {
					data[key] = formatDate(value)
				}
			default:
				if value != "" {
					data[key] = value
				}
			}
		}
	}

	lines := strings.Split(text, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "联系人") && strings.Contains(line, "：") {
			if match := regexp.MustCompile(`联系人[：:]\s*(.{2,20})`).FindStringSubmatch(line); len(match) > 1 {
				data["contact_person"] = strings.TrimSpace(match[1])
			}
		}
		if strings.Contains(line, "电话") && strings.Contains(line, "：") {
			if match := regexp.MustCompile(`电话[：:]\s*([\d\-]+)`).FindStringSubmatch(line); len(match) > 1 {
				data["contact_phone"] = strings.TrimSpace(match[1])
			}
		}
	}

	_ = models.DB

	return data
}

func isValidDate(date string) bool {
	re := regexp.MustCompile(`^\d{4}-\d{1,2}-\d{1,2}$`)
	return re.MatchString(date)
}

func formatDate(date string) string {
	date = strings.ReplaceAll(date, "/", "-")
	parts := strings.Split(date, "-")
	if len(parts) == 3 {
		return fmt.Sprintf("%s-%02s-%02s", parts[0], parts[1], parts[2])
	}
	return date
}

func (h *ContractHandler) CreateStatusChangeRequest(c *gin.Context) {
	contractID, err := strconv.ParseUint(c.Param("contract_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid contract ID"})
		return
	}

	var input services.StatusChangeRequestInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	if !h.contractService.IsStatusChangeRequireApproval(input.ToStatus) {
		contract, err := h.contractService.UpdateContractStatus(uint(contractID), input.ToStatus, userID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"direct": true, "contract": contract})
		return
	}

	request, err := h.contractService.CreateStatusChangeRequest(uint(contractID), input, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if request == nil {
		contract, err := h.contractService.UpdateContractStatus(uint(contractID), input.ToStatus, userID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"direct": true, "contract": contract})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"direct": false, "request": request})
}

func (h *ContractHandler) GetStatusChangeRequests(c *gin.Context) {
	contractID, err := strconv.ParseUint(c.Param("contract_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid contract ID"})
		return
	}

	requests, err := h.contractService.GetStatusChangeRequests(uint(contractID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, requests)
}

func (h *ContractHandler) GetPendingStatusChangeApprovals(c *gin.Context) {
	role, _ := middleware.GetCurrentUserRole(c)
	if role == "" {
		role = "user"
	}

	requests, err := h.contractService.GetPendingStatusChangeRequests(role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, requests)
}

func (h *ContractHandler) ApproveStatusChangeRequest(c *gin.Context) {
	requestID, err := strconv.ParseUint(c.Param("request_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request ID"})
		return
	}

	role, _ := middleware.GetCurrentUserRole(c)
	if role != "sales_manager" && role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "只有销售负责人可以审批状态变更"})
		return
	}

	var input struct {
		Comment string `json:"comment"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	result, err := h.contractService.ApproveStatusChangeRequest(uint(requestID), userID, input.Comment)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *ContractHandler) RejectStatusChangeRequest(c *gin.Context) {
	requestID, err := strconv.ParseUint(c.Param("request_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request ID"})
		return
	}

	role, _ := middleware.GetCurrentUserRole(c)
	if role != "sales_manager" && role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "只有销售负责人可以审批状态变更"})
		return
	}

	var input struct {
		Comment string `json:"comment"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	result, err := h.contractService.RejectStatusChangeRequest(uint(requestID), userID, input.Comment)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// containsMaliciousContent 检查文件内容是否包含恶意代码
func containsMaliciousContent(content []byte, fileExt string) bool {
	if len(content) == 0 {
		return false
	}

	// 转换为小写进行检查
	contentLower := bytes.ToLower(content)

	// 检查PHP标签（对所有文件类型）
	if bytes.Contains(contentLower, []byte("<?php")) || bytes.Contains(contentLower, []byte("<?=")) {
		return true
	}

	// 检查JavaScript标签（对图片文件类型）
	if fileExt == "jpg" || fileExt == "jpeg" || fileExt == "png" || fileExt == "gif" || fileExt == "bmp" || fileExt == "webp" {
		if bytes.Contains(contentLower, []byte("<script")) || bytes.Contains(contentLower, []byte("javascript:")) {
			return true
		}
	}

	// 检查常见的WebShell关键字
	webshellKeywords := [][]byte{
		[]byte("eval("),
		[]byte("system("),
		[]byte("exec("),
		[]byte("shell_exec("),
		[]byte("passthru("),
		[]byte("base64_decode("),
		[]byte("assert("),
	}

	for _, keyword := range webshellKeywords {
		if bytes.Contains(contentLower, keyword) {
			return true
		}
	}

	return false
}
