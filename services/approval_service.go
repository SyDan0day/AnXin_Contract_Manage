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

func (s *ApprovalService) GetApprovalRecords(contractID uint) ([]models.ApprovalRecord, error) {
	var records []models.ApprovalRecord
	if err := models.DB.Where("contract_id = ?", contractID).Preload("Approver").Order("created_at DESC").Find(&records).Error; err != nil {
		return nil, err
	}
	return records, nil
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

	var contracts []models.Contract
	query := models.DB.Preload("Customer").Order("created_at DESC")

	if role == "manager" {
		query = query.Where("status = ?", "draft")
		if err := query.Find(&contracts).Error; err != nil {
			return nil, err
		}
		for _, c := range contracts {
			results = append(results, map[string]interface{}{
				"id":               c.ID,
				"contract_no":      c.ContractNo,
				"title":            c.Title,
				"amount":           c.Amount,
				"status":           c.Status,
				"created_at":       c.CreatedAt,
				"customer":         c.Customer,
				"creator_id":       c.CreatorID,
				"contract_type_id": c.ContractTypeID,
				"approval_id":      uint(0),
			})
		}
	} else if role == "admin" {
		query = query.Where("status = ?", "pending")
		if err := query.Find(&contracts).Error; err != nil {
			return nil, err
		}
		for _, c := range contracts {
			var latestApproval models.ApprovalRecord
			models.DB.Where("contract_id = ?", c.ID).Order("created_at DESC").First(&latestApproval)
			results = append(results, map[string]interface{}{
				"id":               c.ID,
				"contract_no":      c.ContractNo,
				"title":            c.Title,
				"amount":           c.Amount,
				"status":           c.Status,
				"created_at":       c.CreatedAt,
				"customer":         c.Customer,
				"creator_id":       c.CreatorID,
				"contract_type_id": c.ContractTypeID,
				"approval_id":      latestApproval.ID,
			})
		}
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

func (s *ApprovalService) GetExpiringContracts(days int) ([]models.Contract, error) {
	today := time.Now()
	expiryDate := today.AddDate(0, 0, days)

	var contracts []models.Contract
	if err := models.DB.Where("end_date <= ? AND end_date >= ? AND status = ?",
		expiryDate, today, models.StatusActive).Find(&contracts).Error; err != nil {
		return nil, err
	}
	return contracts, nil
}

type Statistics struct {
	TotalContracts     int64   `json:"total_contracts"`
	ActiveContracts    int64   `json:"active_contracts"`
	PendingContracts   int64   `json:"pending_contracts"`
	CompletedContracts int64   `json:"completed_contracts"`
	TotalAmount        float64 `json:"total_amount"`
	ThisMonthContracts int64   `json:"this_month_contracts"`
	ThisMonthAmount    float64 `json:"this_month_amount"`
	ExpiringSoon       int     `json:"expiring_soon"`
}

func (s *ApprovalService) GetStatistics() (*Statistics, error) {
	today := time.Now()
	thisMonthStart := time.Date(today.Year(), today.Month(), 1, 0, 0, 0, 0, time.Local)

	stats := &Statistics{}

	models.DB.Model(&models.Contract{}).Count(&stats.TotalContracts)
	models.DB.Model(&models.Contract{}).Where("status = ?", models.StatusActive).Count(&stats.ActiveContracts)
	models.DB.Model(&models.Contract{}).Where("status = ?", models.StatusPending).Count(&stats.PendingContracts)
	models.DB.Model(&models.Contract{}).Where("status = ?", models.StatusCompleted).Count(&stats.CompletedContracts)

	var totalAmount *float64
	models.DB.Model(&models.Contract{}).Where("amount IS NOT NULL").Select("SUM(amount)").Scan(&totalAmount)
	if totalAmount != nil {
		stats.TotalAmount = *totalAmount
	}

	models.DB.Model(&models.Contract{}).Where("created_at >= ?", thisMonthStart).Count(&stats.ThisMonthContracts)

	var thisMonthAmount *float64
	models.DB.Model(&models.Contract{}).Where("created_at >= ? AND amount IS NOT NULL", thisMonthStart).Select("SUM(amount)").Scan(&thisMonthAmount)
	if thisMonthAmount != nil {
		stats.ThisMonthAmount = *thisMonthAmount
	}

	expiring, _ := s.GetExpiringContracts(30)
	stats.ExpiringSoon = len(expiring)

	return stats, nil
}
