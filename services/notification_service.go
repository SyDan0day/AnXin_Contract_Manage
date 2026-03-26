package services

import (
	"contract-manage/models"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type NotificationService struct {
	db *gorm.DB
}

func NewNotificationService(db *gorm.DB) *NotificationService {
	return &NotificationService{db: db}
}

func (s *NotificationService) CreateNotification(userID uint, notificationType, title, content string, contractID *uint) error {
	notification := models.Notification{
		UserID:     userID,
		Type:       notificationType,
		Title:      title,
		Content:    content,
		ContractID: contractID,
		IsRead:     false,
		CreatedAt:  time.Now(),
	}
	return s.db.Create(&notification).Error
}

func (s *NotificationService) GetNotificationsByUserID(userID uint) ([]models.Notification, error) {
	var notifications []models.Notification
	err := s.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Preload("Contract").
		Find(&notifications).Error
	return notifications, err
}

func (s *NotificationService) GetUnreadNotificationsByUserID(userID uint) ([]models.Notification, error) {
	var notifications []models.Notification
	err := s.db.Where("user_id = ? AND is_read = ?", userID, false).
		Order("created_at DESC").
		Preload("Contract").
		Find(&notifications).Error
	return notifications, err
}

func (s *NotificationService) MarkAsRead(notificationID uint) error {
	return s.db.Model(&models.Notification{}).Where("id = ?", notificationID).Update("is_read", true).Error
}

func (s *NotificationService) MarkAllAsRead(userID uint) error {
	return s.db.Model(&models.Notification{}).Where("user_id = ?", userID).Update("is_read", true).Error
}

func (s *NotificationService) NotifyContractRejected(contractID uint, contractTitle, reason string) error {
	var contract models.Contract
	if err := s.db.First(&contract, contractID).Error; err != nil {
		return err
	}

	title := "合同被退回"
	content := fmt.Sprintf("您的合同「%s」已被退回，退回原因：%s", contractTitle, reason)

	return s.CreateNotification(contract.CreatorID, models.NotificationTypeRejected, title, content, &contractID)
}

func (s *NotificationService) NotifyContractApproved(contractID uint, contractTitle string) error {
	var contract models.Contract
	if err := s.db.First(&contract, contractID).Error; err != nil {
		return err
	}

	title := "合同审批通过"
	content := fmt.Sprintf("您的合同「%s」已通过审批", contractTitle)

	return s.CreateNotification(contract.CreatorID, models.NotificationTypeApproved, title, content, &contractID)
}

func (s *NotificationService) NotifyStatusChange(contractID uint, contractTitle, fromStatus, toStatus string) error {
	var contract models.Contract
	if err := s.db.First(&contract, contractID).Error; err != nil {
		return err
	}

	title := "合同状态变更"
	content := fmt.Sprintf("您的合同「%s」状态已从「%s」变更为「%s」", contractTitle, fromStatus, toStatus)

	return s.CreateNotification(contract.CreatorID, models.NotificationTypeStatusChange, title, content, &contractID)
}
