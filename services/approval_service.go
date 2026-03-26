package services

import (
	"contract-manage/models"
	"errors"
	"time"
)

type ApprovalService struct {
	contractService *ContractService
}

func NewApprovalService() *ApprovalService {
	return &ApprovalService{
		contractService: NewContractService(),
	}
}

func (s *ApprovalService) GetPendingStatusChangesCount() (int, error) {
	requests, err := s.contractService.GetPendingStatusChangeRequests("admin")
	if err != nil {
		return 0, err
	}
	return len(requests), nil
}

func (s *ApprovalService) GetApprovalRecordByID(id uint) (*models.ApprovalRecord, error) {
	var record models.ApprovalRecord
	if err := models.DB.Preload("Approver").First(&record, id).Error; err != nil {
		return nil, err
	}
	return &record, nil
}

func (s *ApprovalService) GetApprovalRecords(contractID uint) ([]map[string]interface{}, error) {
	var approvals []models.WorkflowApproval
	if err := models.DB.Where("contract_id = ?", contractID).Preload("Approver").Order("level ASC").Find(&approvals).Error; err != nil {
		return nil, err
	}

	var contract models.Contract
	if err := models.DB.First(&contract, contractID).Error; err != nil {
		return nil, err
	}

	var results []map[string]interface{}
	roleMap := map[int]string{
		1: "销售负责人审批",
		2: "技术负责人审批",
		3: "财务负责人审批",
	}

	for _, a := range approvals {
		recordMap := map[string]interface{}{
			"id":              a.ID,
			"contract_id":     a.ContractID,
			"level":           a.Level,
			"approver_role":   roleMap[a.Level],
			"status":          a.Status,
			"comment":         a.Comment,
			"approved_at":     a.ApprovedAt,
			"created_at":      a.CreatedAt,
			"contract_status": contract.Status,
		}
		if a.ApproverID != nil && *a.ApproverID > 0 {
			recordMap["approver"] = map[string]interface{}{
				"id":        a.Approver.ID,
				"full_name": a.Approver.FullName,
				"username":  a.Approver.Username,
			}
		}
		results = append(results, recordMap)
	}

	return results, nil
}

type ApprovalRecordCreateInput struct {
	ContractID   uint   `json:"contract_id"`
	Level        int    `json:"level"`
	ApproverRole string `json:"approver_role"`
	Status       string `json:"status"`
	Comment      string `json:"comment"`
}

func (s *ApprovalService) CreateApprovalRecord(input ApprovalRecordCreateInput, approverID uint, approverRole string) (*models.ApprovalRecord, error) {
	record := models.ApprovalRecord{
		ContractID:   input.ContractID,
		ApproverID:   approverID,
		Level:        input.Level,
		ApproverRole: approverRole,
		Status:       models.ApprovalPending,
		Comment:      input.Comment,
	}

	if input.Status != "" {
		record.Status = models.ApprovalStatus(input.Status)
	}

	if err := models.DB.Create(&record).Error; err != nil {
		return nil, err
	}
	return &record, nil
}

type ApprovalRecordUpdateInput struct {
	Status  string `json:"status" binding:"required"`
	Comment string `json:"comment"`
}

func (s *ApprovalService) UpdateApprovalRecord(id uint, input ApprovalRecordUpdateInput, contractStatus string, operatorID uint) (*models.ApprovalRecord, error) {
	record, err := s.GetApprovalRecordByID(id)
	if err != nil {
		return nil, err
	}

	if record.Status != models.ApprovalPending {
		return nil, errors.New("this approval has already been processed")
	}

	oldStatus := string(models.StatusPending)
	var newStatus string

	now := time.Now()
	record.Status = models.ApprovalStatus(input.Status)
	record.Comment = input.Comment
	record.ApprovedAt = &now

	if err := models.DB.Save(record).Error; err != nil {
		return nil, err
	}

	if contractStatus != "" {
		var contract models.Contract
		if err := models.DB.First(&contract, record.ContractID).Error; err == nil {
			oldStatus = string(contract.Status)
			contract.Status = models.ContractStatus(contractStatus)
			newStatus = contractStatus
			models.DB.Save(&contract)

			var eventType models.LifecycleEventType
			var description string
			if input.Status == "approved" {
				eventType = models.LifecycleApproved
				description = "审批通过"
			} else if input.Status == "rejected" {
				eventType = models.LifecycleRejected
				description = "审批拒绝"
			}

			if eventType != "" {
				s.contractService.AddLifecycleEvent(record.ContractID, LifecycleEventInput{
					EventType:   string(eventType),
					FromStatus:  oldStatus,
					ToStatus:    newStatus,
					Description: description,
				}, operatorID)
			}
		}
	}

	return record, nil
}

func (s *ApprovalService) GetPendingApprovalsByRole(role string, userID uint) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	approvalRoles := []string{"sales_manager", "tech_leader", "finance_leader", "admin"}
	isApprovalRole := false
	isAdmin := false
	for _, r := range approvalRoles {
		if role == r {
			isApprovalRole = true
			if role == "admin" {
				isAdmin = true
			}
			break
		}
	}

	if !isApprovalRole {
		return results, nil
	}

	var workflows []models.ApprovalWorkflow
	if err := models.DB.Where("status = ?", models.WorkflowStatusPending).Find(&workflows).Error; err != nil {
		return nil, err
	}

	for _, wf := range workflows {
		var approval models.WorkflowApproval
		var queryErr error

		// Admin sees all pending approvals, other roles see only their role's approvals
		if isAdmin {
			queryErr = models.DB.Where("workflow_id = ? AND level = ? AND status = ?",
				wf.ID, wf.CurrentLevel, models.WorkflowStatusPending).First(&approval).Error
		} else {
			queryErr = models.DB.Where("workflow_id = ? AND approver_role = ? AND level = ? AND status = ?",
				wf.ID, role, wf.CurrentLevel, models.WorkflowStatusPending).First(&approval).Error
		}

		if queryErr != nil {
			continue
		}

		var contract models.Contract
		if err := models.DB.Preload("Customer").First(&contract, wf.ContractID).Error; err != nil {
			continue
		}

		results = append(results, map[string]interface{}{
			"id":               contract.ID,
			"contract_no":      contract.ContractNo,
			"title":            contract.Title,
			"amount":           contract.Amount,
			"status":           contract.Status,
			"created_at":       contract.CreatedAt,
			"customer":         contract.Customer,
			"creator_id":       contract.CreatorID,
			"contract_type_id": contract.ContractTypeID,
			"approval_id":      approval.ID,
			"workflow_id":      wf.ID,
			"approval_level":   approval.Level,
			"approver_role":    approval.ApproverRole,
			"current_level":    wf.CurrentLevel,
		})
	}

	return results, nil
}

func (s *ApprovalService) SubmitForApproval(contractID uint, userID uint) error {
	var contract models.Contract
	if err := models.DB.First(&contract, contractID).Error; err != nil {
		return err
	}

	oldStatus := string(contract.Status)
	contract.Status = models.StatusPending
	if err := models.DB.Save(&contract).Error; err != nil {
		return err
	}

	s.contractService.AddLifecycleEvent(contractID, LifecycleEventInput{
		EventType:   string(models.LifecycleSubmitted),
		FromStatus:  oldStatus,
		ToStatus:    string(models.StatusPending),
		Description: "合同提交审批",
	}, userID)

	return nil
}

func (s *ApprovalService) GetReminderByID(id uint) (*models.Reminder, error) {
	var reminder models.Reminder
	if err := models.DB.First(&reminder, id).Error; err != nil {
		return nil, err
	}
	return &reminder, nil
}

func (s *ApprovalService) GetReminders(contractID uint) ([]models.Reminder, error) {
	var reminders []models.Reminder
	if err := models.DB.Where("contract_id = ?", contractID).Order("reminder_date DESC").Find(&reminders).Error; err != nil {
		return nil, err
	}
	return reminders, nil
}

type ReminderCreateInput struct {
	ContractID   uint      `json:"contract_id" binding:"required"`
	Type         string    `json:"type" binding:"required"`
	ReminderDate *JSONTime `json:"reminder_date" binding:"required"`
	DaysBefore   int       `json:"days_before" binding:"required"`
}

func (s *ApprovalService) CreateReminder(input ReminderCreateInput) (*models.Reminder, error) {
	reminder := models.Reminder{
		ContractID: input.ContractID,
		Type:       input.Type,
		DaysBefore: input.DaysBefore,
		IsSent:     false,
	}

	if input.ReminderDate != nil && !input.ReminderDate.Time.IsZero() {
		reminder.ReminderDate = &input.ReminderDate.Time
	}

	if err := models.DB.Create(&reminder).Error; err != nil {
		return nil, err
	}
	return &reminder, nil
}

func (s *ApprovalService) UpdateReminderSent(id uint) error {
	reminder, err := s.GetReminderByID(id)
	if err != nil {
		return err
	}

	now := time.Now()
	reminder.IsSent = true
	reminder.SentAt = &now

	return models.DB.Save(reminder).Error
}

// GetExpiringContractsWithAuth 获取即将到期合同（带权限检查）
func (s *ApprovalService) GetExpiringContractsWithAuth(days int, userID uint, role string) ([]models.Contract, error) {
	today := time.Now()
	expiryDate := today.AddDate(0, 0, days)

	var contracts []models.Contract
	query := models.DB.Where("end_date <= ? AND end_date >= ? AND status = ?",
		expiryDate, today, models.StatusActive)

	// 根据角色过滤
	userRole := models.UserRole(role)
	if !models.CanViewAllContracts(userRole) {
		query = query.Where("creator_id = ?", userID)
	}

	if err := query.Find(&contracts).Error; err != nil {
		return nil, err
	}
	return contracts, nil
}

// GetExpiringContracts 获取即将到期合同（兼容旧接口）
func (s *ApprovalService) GetExpiringContracts(days int) ([]models.Contract, error) {
	return s.GetExpiringContractsWithAuth(days, 0, "admin")
}

type Statistics struct {
	TotalContracts      int64   `json:"total_contracts"`
	ActiveContracts     int64   `json:"active_contracts"`
	PendingContracts    int64   `json:"pending_contracts"`
	CompletedContracts  int64   `json:"completed_contracts"`
	DraftContracts      int64   `json:"draft_contracts"`
	TerminatedContracts int64   `json:"terminated_contracts"`
	TotalAmount         float64 `json:"total_amount"`
	ThisMonthContracts  int64   `json:"this_month_contracts"`
	ThisMonthAmount     float64 `json:"this_month_amount"`
	ExpiringSoon        int     `json:"expiring_soon"`
}

// GetStatisticsParams 统计数据参数
type GetStatisticsParams struct {
	UserID uint
	Role   string
}

// GetStatistics 获取统计数据（支持角色过滤）
func (s *ApprovalService) GetStatistics(params GetStatisticsParams) (*Statistics, error) {
	today := time.Now()
	thisMonthStart := time.Date(today.Year(), today.Month(), 1, 0, 0, 0, 0, time.Local)

	stats := &Statistics{}

	// 根据角色获取条件
	role := models.UserRole(params.Role)
	canViewAll := models.CanViewAllContracts(role)

	// 统计总数
	models.DB.Model(&models.Contract{}).Count(&stats.TotalContracts)
	models.DB.Model(&models.Contract{}).Where("status = ?", models.StatusActive).Count(&stats.ActiveContracts)
	models.DB.Model(&models.Contract{}).Where("status = ?", models.StatusPending).Count(&stats.PendingContracts)
	models.DB.Model(&models.Contract{}).Where("status = ?", models.StatusCompleted).Count(&stats.CompletedContracts)
	models.DB.Model(&models.Contract{}).Where("status = ?", models.StatusDraft).Count(&stats.DraftContracts)
	models.DB.Model(&models.Contract{}).Where("status = ?", models.StatusTerminated).Count(&stats.TerminatedContracts)

	// 统计金额
	var totalAmount, thisMonthAmount float64
	models.DB.Model(&models.Contract{}).Where("amount IS NOT NULL").Select("SUM(amount)").Scan(&totalAmount)
	models.DB.Model(&models.Contract{}).Where("created_at >= ? AND amount IS NOT NULL", thisMonthStart).Select("SUM(amount)").Scan(&thisMonthAmount)
	stats.TotalAmount = totalAmount
	stats.ThisMonthAmount = thisMonthAmount
	models.DB.Model(&models.Contract{}).Where("created_at >= ?", thisMonthStart).Count(&stats.ThisMonthContracts)

	// 非管理员角色，需要过滤统计数据
	if !canViewAll {
		// 重新统计符合条件的数量
		models.DB.Model(&models.Contract{}).Where("creator_id = ?", params.UserID).Count(&stats.TotalContracts)
		models.DB.Model(&models.Contract{}).Where("creator_id = ? AND status = ?", params.UserID, models.StatusActive).Count(&stats.ActiveContracts)
		models.DB.Model(&models.Contract{}).Where("creator_id = ? AND status = ?", params.UserID, models.StatusPending).Count(&stats.PendingContracts)
		models.DB.Model(&models.Contract{}).Where("creator_id = ? AND status = ?", params.UserID, models.StatusCompleted).Count(&stats.CompletedContracts)
		models.DB.Model(&models.Contract{}).Where("creator_id = ? AND status = ?", params.UserID, models.StatusDraft).Count(&stats.DraftContracts)
		models.DB.Model(&models.Contract{}).Where("creator_id = ? AND status = ?", params.UserID, models.StatusTerminated).Count(&stats.TerminatedContracts)

		// 重新统计金额
		var userTotalAmount, userThisMonthAmount float64
		models.DB.Model(&models.Contract{}).Where("creator_id = ? AND amount IS NOT NULL", params.UserID).Select("SUM(amount)").Scan(&userTotalAmount)
		models.DB.Model(&models.Contract{}).Where("creator_id = ? AND created_at >= ? AND amount IS NOT NULL", params.UserID, thisMonthStart).Select("SUM(amount)").Scan(&userThisMonthAmount)
		stats.TotalAmount = userTotalAmount
		stats.ThisMonthAmount = userThisMonthAmount

		models.DB.Model(&models.Contract{}).Where("creator_id = ? AND created_at >= ?", params.UserID, thisMonthStart).Count(&stats.ThisMonthContracts)
	}

	expiring, _ := s.GetExpiringContractsWithAuth(30, params.UserID, params.Role)
	stats.ExpiringSoon = len(expiring)

	return stats, nil
}
