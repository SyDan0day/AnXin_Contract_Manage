package handlers

import (
	"archive/zip"
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
	skip, _ := strconv.Atoi(c.DefaultQuery("skip", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))
	customerID, _ := strconv.ParseUint(c.Query("customer_id"), 10, 32)
	contractTypeID, _ := strconv.ParseUint(c.Query("contract_type_id"), 10, 32)
	status := c.Query("status")

	contracts, err := h.contractService.GetContracts(skip, limit, uint(customerID), uint(contractTypeID), status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, contracts)
}

func (h *ContractHandler) GetContractByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("contract_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid contract ID"})
		return
	}

	contract, err := h.contractService.GetContractByID(uint(id))
	if err != nil {
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

func (h *ContractHandler) GetContractExecutions(c *gin.Context) {
	contractID, err := strconv.ParseUint(c.Param("contract_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid contract ID"})
		return
	}

	executions, err := h.contractService.GetContractExecutions(uint(contractID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, executions)
}

func (h *ContractHandler) CreateContractExecution(c *gin.Context) {
	contractID, err := strconv.ParseUint(c.Param("contract_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid contract ID"})
		return
	}

	var input services.ContractExecutionCreateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input.ContractID = uint(contractID)

	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	execution, err := h.contractService.CreateContractExecution(input, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, execution)
}

func (h *ContractHandler) DeleteExecution(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("execution_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.contractService.DeleteExecution(uint(id)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

	filename := file.Filename
	uploadDir := config.AppConfig.UploadDir
	if uploadDir == "" {
		uploadDir = "uploads"
	}

	filePath := fmt.Sprintf("%s/%d/%s", uploadDir, contractID, filename)

	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建目录失败"})
		return
	}

	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存文件失败"})
		return
	}

	input := services.DocumentCreateInput{
		ContractID: uint(contractID),
		Name:       filename,
		FilePath:   "/" + filePath,
		FileSize:   int(file.Size),
		FileType:   filepath.Ext(filename)[1:],
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
		"contract_no":   `合同编号[：:]\s*([A-Z0-9\-]+)`,
		"title":         `合同名称[：:]\s*(.{2,50})`,
		"customer_name": `甲方[（(]客户[）)][：:]\s*(.{2,50})`,
		"amount":        `合同金额[：:]\s*([\d,]+\.?\d*)\s*(元|万)?`,
		"sign_date":     `签订日期[：:]\s*(\d{4}[-/年]\d{1,2}[-/月]\d{1,2}[日]?)`,
		"start_date":    `开始日期[：:]\s*(\d{4}[-/年]\d{1,2}[-/月]\d{1,2}[日]?)`,
		"end_date":      `结束日期[：:]\s*(\d{4}[-/年]\d{1,2}[-/月]\d{1,2}[日]?)`,
		"contract_type": `合同类型[：:]\s*(.{2,20})`,
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
