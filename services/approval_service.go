package services

import (
	"contract-manage/models"
	"errors"
	"time"
)

type ApprovalService struct{}

func NewApprovalService() *ApprovalService {
	return &ApprovalService{}
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
	ContractID uint   `json:"contract_id" binding:"required"`
	Status     string `json:"status"`
	Comment    string `json:"comment"`
}

func (s *ApprovalService) CreateApprovalRecord(input ApprovalRecordCreateInput, approverID uint) (*models.ApprovalRecord, error) {
	record := models.ApprovalRecord{
		ContractID: input.ContractID,
		ApproverID: approverID,
		Status:     models.ApprovalPending,
		Comment:    input.Comment,
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

func (s *ApprovalService) UpdateApprovalRecord(id uint, input ApprovalRecordUpdateInput) (*models.ApprovalRecord, error) {
	record, err := s.GetApprovalRecordByID(id)
	if err != nil {
		return nil, err
	}

	if record.Status != models.ApprovalPending {
		return nil, errors.New("this approval has already been processed")
	}

	now := time.Now()
	record.Status = models.ApprovalStatus(input.Status)
	record.Comment = input.Comment
	record.ApprovedAt = &now

	if err := models.DB.Save(record).Error; err != nil {
		return nil, err
	}
	return record, nil
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
	ContractID   uint       `json:"contract_id" binding:"required"`
	Type         string     `json:"type" binding:"required"`
	ReminderDate *time.Time `json:"reminder_date" binding:"required"`
	DaysBefore   int        `json:"days_before" binding:"required"`
}

func (s *ApprovalService) CreateReminder(input ReminderCreateInput) (*models.Reminder, error) {
	reminder := models.Reminder{
		ContractID:   input.ContractID,
		Type:         input.Type,
		ReminderDate: input.ReminderDate,
		DaysBefore:   input.DaysBefore,
		IsSent:       false,
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
	TotalContracts      int64    `json:"total_contracts"`
	ActiveContracts     int64    `json:"active_contracts"`
	PendingContracts   int64    `json:"pending_contracts"`
	CompletedContracts int64    `json:"completed_contracts"`
	TotalAmount        float64  `json:"total_amount"`
	ThisMonthContracts int64    `json:"this_month_contracts"`
	ThisMonthAmount    float64  `json:"this_month_amount"`
	ExpiringSoon       int      `json:"expiring_soon"`
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