package services

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"contract-manage/models"

	"gorm.io/gorm"
)

// 定义允许的状态转换
var allowedStatusTransitions = map[models.ContractStatus]map[models.ContractStatus]bool{
	models.StatusDraft: {
		models.StatusPending:    true,
		models.StatusTerminated: true,
	},
	models.StatusPending: {
		models.StatusActive:     true,
		models.StatusDraft:      true,
		models.StatusTerminated: true,
	},
	models.StatusActive: {
		models.StatusCompleted:  true,
		models.StatusTerminated: true,
	},
	models.StatusCompleted: {
		models.StatusTerminated: true,
	},
	models.StatusTerminated: {
		// 终止状态不能转换到其他状态
	},
	models.StatusArchived: {
		// 已归档状态不能转换到其他状态
	},
}

// isValidStatusTransition 检查状态转换是否允许
func isValidStatusTransition(from, to models.ContractStatus) bool {
	if allowedTransitions, exists := allowedStatusTransitions[from]; exists {
		if allowedTransitions[to] {
			return true
		}
	}
	return false
}

type ContractService struct{}

func NewContractService() *ContractService {
	return &ContractService{}
}

type JSONTime struct {
	time.Time
}

func (t *JSONTime) UnmarshalJSON(data []byte) error {
	str := strings.Trim(string(data), `"`)
	if str == "" || str == "null" {
		t.Time = time.Time{}
		return nil
	}

	formats := []string{
		"2006-01-02T15:04:05Z07:00",
		"2006-01-02T15:04:05",
		"2006-01-02",
	}

	for _, format := range formats {
		if parsed, err := time.Parse(format, str); err == nil {
			t.Time = parsed
			return nil
		}
	}
	return errors.New("invalid date format")
}

type ContractCreateInput struct {
	Title          string    `json:"title" binding:"required"`
	CustomerID     uint      `json:"customer_id" binding:"required"`
	ContractTypeID uint      `json:"contract_type_id" binding:"required"`
	Amount         float64   `json:"amount"`
	Currency       string    `json:"currency"`
	SignDate       *JSONTime `json:"sign_date"`
	StartDate      *JSONTime `json:"start_date"`
	EndDate        *JSONTime `json:"end_date"`
	PaymentTerms   string    `json:"payment_terms"`
	Content        string    `json:"content"`
	Notes          string    `json:"notes"`
}

func (s *ContractService) generateContractNo() string {
	today := time.Now()
	prefix := fmt.Sprintf("CT%s", today.Format("200601"))

	var lastContract models.Contract
	models.DB.Where("contract_no LIKE ?", prefix+"%").Order("contract_no DESC").First(&lastContract)

	var newNo string
	if lastContract.ID != 0 {
		lastNo := lastContract.ContractNo[len(lastContract.ContractNo)-4:]
		var num int
		fmt.Sscanf(lastNo, "%d", &num)
		newNo = fmt.Sprintf("%04d", num+1)
	} else {
		newNo = "0001"
	}

	return prefix + newNo
}

// GetContractByIDWithAuth 获取合同详情（带权限检查）
func (s *ContractService) GetContractByIDWithAuth(id uint, userID uint, role string) (*models.Contract, error) {
	var contract models.Contract
	if err := models.DB.Preload("Customer").Preload("Creator").Preload("ContractType").First(&contract, id).Error; err != nil {
		return nil, err
	}

	// 检查权限
	userRole := models.UserRole(role)
	if !models.CanViewAllContracts(userRole) && contract.CreatorID != userID {
		return nil, fmt.Errorf("无权限查看此合同")
	}

	return &contract, nil
}

// GetContractByID 获取合同详情（兼容旧接口，无权限检查）
func (s *ContractService) GetContractByID(id uint) (*models.Contract, error) {
	var contract models.Contract
	if err := models.DB.Preload("Customer").Preload("Creator").Preload("ContractType").First(&contract, id).Error; err != nil {
		return nil, err
	}
	return &contract, nil
}

func (s *ContractService) GetContractByNo(contractNo string) (*models.Contract, error) {
	var contract models.Contract
	if err := models.DB.Where("contract_no = ?", contractNo).First(&contract).Error; err != nil {
		return nil, err
	}
	return &contract, nil
}

// GetContractsParams 获取合同列表参数
type GetContractsParams struct {
	Skip           int
	Limit          int
	CustomerID     uint
	ContractTypeID uint
	Status         string
	Keyword        string
	UserID         uint   // 当前用户ID
	Role           string // 当前用户角色
}

// GetContracts 获取合同列表（根据角色过滤）
func (s *ContractService) GetContracts(params GetContractsParams) ([]models.Contract, int64, error) {
	var contracts []models.Contract
	var total int64

	query := models.DB.Model(&models.Contract{})

	// 根据角色过滤数据
	role := models.UserRole(params.Role)
	if !models.CanViewAllContracts(role) {
		// 非管理员角色，只能查看自己创建的合同
		query = query.Where("creator_id = ?", params.UserID)
	}

	// 构建查询条件
	if params.CustomerID > 0 {
		query = query.Where("customer_id = ?", params.CustomerID)
	}
	if params.ContractTypeID > 0 {
		query = query.Where("contract_type_id = ?", params.ContractTypeID)
	}
	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}
	if params.Keyword != "" {
		query = query.Where("contract_no LIKE ? OR title LIKE ?", "%"+params.Keyword+"%", "%"+params.Keyword+"%")
	}

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	if err := query.Preload("Customer").Preload("Creator").Preload("ContractType").
		Order("created_at DESC").Offset(params.Skip).Limit(params.Limit).Find(&contracts).Error; err != nil {
		return nil, 0, err
	}

	return contracts, total, nil
}

func (s *ContractService) CreateContract(input ContractCreateInput, creatorID uint) (*models.Contract, error) {
	contract := models.Contract{
		ContractNo:     s.generateContractNo(),
		Title:          input.Title,
		CustomerID:     input.CustomerID,
		ContractTypeID: input.ContractTypeID,
		Amount:         input.Amount,
		Currency:       input.Currency,
		PaymentTerms:   input.PaymentTerms,
		Content:        input.Content,
		Notes:          input.Notes,
		CreatorID:      creatorID,
		Status:         models.StatusDraft,
	}

	if input.SignDate != nil && !input.SignDate.Time.IsZero() {
		contract.SignDate = &input.SignDate.Time
	}
	if input.StartDate != nil && !input.StartDate.Time.IsZero() {
		contract.StartDate = &input.StartDate.Time
	}
	if input.EndDate != nil && !input.EndDate.Time.IsZero() {
		contract.EndDate = &input.EndDate.Time
	}

	if contract.Currency == "" {
		contract.Currency = "CNY"
	}

	if err := models.DB.Create(&contract).Error; err != nil {
		return nil, err
	}

	s.AddLifecycleEvent(contract.ID, LifecycleEventInput{
		EventType:   string(models.LifecycleCreated),
		Description: "合同创建",
	}, creatorID)

	return &contract, nil
}

type ContractUpdateInput struct {
	Title          string                `json:"title"`
	CustomerID     uint                  `json:"customer_id"`
	ContractTypeID uint                  `json:"contract_type_id"`
	Amount         float64               `json:"amount"`
	Currency       string                `json:"currency"`
	Status         models.ContractStatus `json:"status"`
	SignDate       *JSONTime             `json:"sign_date"`
	StartDate      *JSONTime             `json:"start_date"`
	EndDate        *JSONTime             `json:"end_date"`
	PaymentTerms   string                `json:"payment_terms"`
	Content        string                `json:"content"`
	Notes          string                `json:"notes"`
}

func (s *ContractService) UpdateContract(id uint, input ContractUpdateInput) (*models.Contract, error) {
	contract, err := s.GetContractByID(id)
	if err != nil {
		return nil, err
	}

	updates := map[string]interface{}{}
	if input.Title != "" {
		updates["title"] = input.Title
	}
	if input.CustomerID > 0 {
		updates["customer_id"] = input.CustomerID
	}
	if input.ContractTypeID > 0 {
		updates["contract_type_id"] = input.ContractTypeID
	}
	if input.Amount > 0 {
		updates["amount"] = input.Amount
	}
	if input.Currency != "" {
		updates["currency"] = input.Currency
	}
	if input.Status != "" {
		updates["status"] = input.Status
	}
	if input.SignDate != nil && !input.SignDate.Time.IsZero() {
		updates["sign_date"] = input.SignDate.Time
	}
	if input.StartDate != nil && !input.StartDate.Time.IsZero() {
		updates["start_date"] = input.StartDate.Time
	}
	if input.EndDate != nil && !input.EndDate.Time.IsZero() {
		updates["end_date"] = input.EndDate.Time
	}
	if input.PaymentTerms != "" {
		updates["payment_terms"] = input.PaymentTerms
	}
	if input.Content != "" {
		updates["content"] = input.Content
	}
	if input.Notes != "" {
		updates["notes"] = input.Notes
	}

	if err := models.DB.Model(contract).Updates(updates).Error; err != nil {
		return nil, err
	}
	return contract, nil
}

func (s *ContractService) DeleteContract(id uint) error {
	result := models.DB.Delete(&models.Contract{}, id)
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}

func (s *ContractService) GetDocumentByID(id uint) (*models.Document, error) {
	var document models.Document
	if err := models.DB.First(&document, id).Error; err != nil {
		return nil, err
	}
	return &document, nil
}

func (s *ContractService) GetDocuments(contractID uint) ([]models.Document, error) {
	var documents []models.Document
	if err := models.DB.Where("contract_id = ?", contractID).Order("created_at DESC").Find(&documents).Error; err != nil {
		return nil, err
	}
	return documents, nil
}

type DocumentCreateInput struct {
	ContractID uint   `json:"contract_id"`
	Name       string `json:"name" binding:"required"`
	FilePath   string `json:"file_path" binding:"required"`
	FileSize   int    `json:"file_size" binding:"required"`
	FileType   string `json:"file_type"`
}

func (s *ContractService) CreateDocument(input DocumentCreateInput, uploaderID uint) (*models.Document, error) {
	document := models.Document{
		ContractID: input.ContractID,
		Name:       input.Name,
		FilePath:   input.FilePath,
		FileSize:   input.FileSize,
		FileType:   input.FileType,
		Version:    "1.0",
		UploaderID: uploaderID,
	}

	if err := models.DB.Create(&document).Error; err != nil {
		return nil, err
	}
	return &document, nil
}

func (s *ContractService) DeleteDocument(id uint) error {
	result := models.DB.Delete(&models.Document{}, id)
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}

type LifecycleEventInput struct {
	EventType   string  `json:"event_type"`
	FromStatus  string  `json:"from_status"`
	ToStatus    string  `json:"to_status"`
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
}

func (s *ContractService) AddLifecycleEvent(contractID uint, input LifecycleEventInput, operatorID uint) (*models.ContractLifecycleEvent, error) {
	event := models.ContractLifecycleEvent{
		ContractID:  contractID,
		EventType:   models.LifecycleEventType(input.EventType),
		FromStatus:  input.FromStatus,
		ToStatus:    input.ToStatus,
		Amount:      input.Amount,
		Description: input.Description,
		OperatorID:  operatorID,
	}

	if err := models.DB.Create(&event).Error; err != nil {
		return nil, err
	}
	return &event, nil
}

func (s *ContractService) GetLifecycleEvents(contractID uint) ([]models.ContractLifecycleEvent, error) {
	var events []models.ContractLifecycleEvent
	if err := models.DB.Where("contract_id = ?", contractID).Order("created_at ASC").Find(&events).Error; err != nil {
		return nil, err
	}
	return events, nil
}

// UpdateContractStatus 更新合同状态
func (s *ContractService) UpdateContractStatus(contractID uint, newStatus string, operatorID uint) (*models.Contract, error) {
	var contract models.Contract
	if err := models.DB.First(&contract, contractID).Error; err != nil {
		return nil, err
	}

	newStatusEnum := models.ContractStatus(newStatus)
	oldStatus := contract.Status

	// 验证状态转换是否允许
	if !isValidStatusTransition(oldStatus, newStatusEnum) {
		return nil, fmt.Errorf("不允许的状态转换：%s -> %s", oldStatus, newStatusEnum)
	}

	contract.Status = newStatusEnum

	if err := models.DB.Save(&contract).Error; err != nil {
		return nil, err
	}

	// 如果合同状态变更为终止或归档，结束审批流程
	if newStatus == "terminated" || newStatus == "archived" {
		models.DB.Model(&models.ApprovalWorkflow{}).Where("contract_id = ? AND status = ?", contractID, "pending").
			Update("status", models.WorkflowStatusCancelled)
		models.DB.Model(&models.WorkflowApproval{}).Where("contract_id = ? AND status = ?", contractID, "pending").
			Update("status", models.WorkflowStatusCancelled)
	}

	s.AddLifecycleEvent(contractID, LifecycleEventInput{
		EventType:   "status_changed",
		FromStatus:  string(oldStatus),
		ToStatus:    string(newStatusEnum),
		Description: "合同状态变更",
	}, operatorID)

	return &contract, nil
}

func (s *ContractService) ArchiveContract(contractID uint, operatorID uint) (*models.Contract, error) {
	var contract models.Contract
	if err := models.DB.First(&contract, contractID).Error; err != nil {
		return nil, err
	}

	oldStatus := string(contract.Status)
	contract.Status = models.StatusArchived

	if err := models.DB.Save(&contract).Error; err != nil {
		return nil, err
	}

	// 归档时结束审批流程
	models.DB.Model(&models.ApprovalWorkflow{}).Where("contract_id = ? AND status = ?", contractID, "pending").
		Update("status", models.WorkflowStatusCancelled)
	models.DB.Model(&models.WorkflowApproval{}).Where("contract_id = ? AND status = ?", contractID, "pending").
		Update("status", models.WorkflowStatusCancelled)

	s.AddLifecycleEvent(contractID, LifecycleEventInput{
		EventType:   string(models.LifecycleArchived),
		FromStatus:  oldStatus,
		ToStatus:    string(models.StatusArchived),
		Description: "合同已归档",
	}, operatorID)

	return &contract, nil
}

var StatusChangeRequireApproval = []string{
	"archived",
	"terminated",
}

func (s *ContractService) IsStatusChangeRequireApproval(newStatus string) bool {
	for _, status := range StatusChangeRequireApproval {
		if status == newStatus {
			return true
		}
	}
	return false
}

type StatusChangeRequestInput struct {
	ToStatus string `json:"to_status" binding:"required"`
	Reason   string `json:"reason"`
}

func (s *ContractService) CreateStatusChangeRequest(contractID uint, input StatusChangeRequestInput, requesterID uint) (*models.StatusChangeRequest, error) {
	var contract models.Contract
	if err := models.DB.First(&contract, contractID).Error; err != nil {
		return nil, err
	}

	if !s.IsStatusChangeRequireApproval(input.ToStatus) {
		return nil, nil
	}

	var existingRequest models.StatusChangeRequest
	if err := models.DB.Where("contract_id = ? AND status = ?", contractID, "pending").First(&existingRequest).Error; err == nil {
		return nil, fmt.Errorf("该合同已有待审核的状态变更申请")
	}

	request := models.StatusChangeRequest{
		ContractID:  contractID,
		FromStatus:  string(contract.Status),
		ToStatus:    input.ToStatus,
		Reason:      input.Reason,
		RequesterID: requesterID,
		Status:      "pending",
	}

	if err := models.DB.Create(&request).Error; err != nil {
		return nil, err
	}

	return &request, nil
}

func (s *ContractService) GetStatusChangeRequests(contractID uint) ([]models.StatusChangeRequest, error) {
	var requests []models.StatusChangeRequest
	if err := models.DB.Preload("Requester").Preload("Approver").Where("contract_id = ?", contractID).Order("created_at DESC").Find(&requests).Error; err != nil {
		return nil, err
	}
	return requests, nil
}

func (s *ContractService) GetPendingStatusChangeRequests(role string) ([]models.StatusChangeRequest, error) {
	var requests []models.StatusChangeRequest
	query := models.DB.Preload("Contract.Customer").Preload("Requester").Order("created_at DESC")

	if role == "sales_manager" || role == "admin" {
		query = query.Where("status = ?", "pending")
	}

	if err := query.Find(&requests).Error; err != nil {
		return nil, err
	}
	return requests, nil
}

func (s *ContractService) ApproveStatusChangeRequest(requestID uint, approverID uint, comment string) (*models.StatusChangeRequest, error) {
	var request models.StatusChangeRequest
	if err := models.DB.Preload("Contract").First(&request, requestID).Error; err != nil {
		return nil, err
	}

	if request.Status != "pending" {
		return nil, fmt.Errorf("该申请已被处理")
	}

	now := time.Now()
	request.Status = "approved"
	request.ApproverID = &approverID
	request.Comment = comment
	request.ApprovedAt = &now

	if err := models.DB.Save(&request).Error; err != nil {
		return nil, err
	}

	var contract models.Contract
	if err := models.DB.First(&contract, request.ContractID).Error; err != nil {
		return nil, err
	}

	oldStatus := string(contract.Status)
	contract.Status = models.ContractStatus(request.ToStatus)
	if err := models.DB.Save(&contract).Error; err != nil {
		return nil, err
	}

	var eventType models.LifecycleEventType
	var description string
	switch request.ToStatus {
	case "archived":
		eventType = models.LifecycleArchived
		description = "合同已归档"
	case "terminated":
		eventType = models.LifecycleTerminated
		description = "合同已终止"
	default:
		eventType = "status_changed"
		description = "合同状态变更"
	}

	s.AddLifecycleEvent(request.ContractID, LifecycleEventInput{
		EventType:   string(eventType),
		FromStatus:  oldStatus,
		ToStatus:    request.ToStatus,
		Description: description,
	}, approverID)

	return &request, nil
}

func (s *ContractService) RejectStatusChangeRequest(requestID uint, approverID uint, comment string) (*models.StatusChangeRequest, error) {
	var request models.StatusChangeRequest
	if err := models.DB.First(&request, requestID).Error; err != nil {
		return nil, err
	}

	if request.Status != "pending" {
		return nil, fmt.Errorf("该申请已被处理")
	}

	now := time.Now()
	request.Status = "rejected"
	request.ApproverID = &approverID
	request.Comment = comment
	request.ApprovedAt = &now

	if err := models.DB.Save(&request).Error; err != nil {
		return nil, err
	}

	return &request, nil
}
