package services

import (
	"contract-manage/models"
	"time"
)

type AuditService struct{}

func NewAuditService() *AuditService {
	return &AuditService{}
}

func (s *AuditService) CreateAuditLog(log models.AuditLog) error {
	return models.DB.Create(&log).Error
}

func (s *AuditService) GetAuditLogs(page, pageSize int, username, action, module, startDate, endDate string) ([]models.AuditLog, int64, error) {
	var logs []models.AuditLog
	var total int64

	query := models.DB.Model(&models.AuditLog{})

	if username != "" {
		query = query.Where("username LIKE ?", "%"+username+"%")
	}
	if action != "" {
		query = query.Where("action LIKE ?", "%"+action+"%")
	}
	if module != "" {
		query = query.Where("module = ?", module)
	}
	if startDate != "" {
		if t, err := time.Parse("2006-01-02", startDate); err == nil {
			query = query.Where("created_at >= ?", t)
		}
	}
	if endDate != "" {
		if t, err := time.Parse("2006-01-02", endDate); err == nil {
			t = t.Add(24 * time.Hour)
			query = query.Where("created_at < ?", t)
		}
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.Preload("User").Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&logs).Error

	return logs, total, err
}

func (s *AuditService) DeleteAuditLog(id uint) error {
	return models.DB.Delete(&models.AuditLog{}, id).Error
}

func (s *AuditService) DeleteAuditLogs(ids []uint) error {
	if len(ids) == 0 {
		return nil
	}
	return models.DB.Where("id IN ?", ids).Delete(&models.AuditLog{}).Error
}

func (s *AuditService) GetAuditStatistics(startDate, endDate string) (map[string]interface{}, error) {
	var total int64
	var loginCount int64
	var actionCounts []struct {
		Action string
		Count  int64
	}

	query := models.DB.Model(&models.AuditLog{})

	if startDate != "" {
		if t, err := time.Parse("2006-01-02", startDate); err == nil {
			query = query.Where("created_at >= ?", t)
		}
	}
	if endDate != "" {
		if t, err := time.Parse("2006-01-02", endDate); err == nil {
			t = t.Add(24 * time.Hour)
			query = query.Where("created_at < ?", t)
		}
	}

	query.Count(&total)

	models.DB.Model(&models.AuditLog{}).Where("action LIKE ?", "%login%").Count(&loginCount)

	models.DB.Model(&models.AuditLog{}).
		Select("action, COUNT(*) as count").
		Group("action").
		Order("count DESC").
		Limit(10).
		Scan(&actionCounts)

	stats := map[string]interface{}{
		"total":        total,
		"loginCount":   loginCount,
		"actionCounts": actionCounts,
	}

	return stats, nil
}

func (s *AuditService) RecordLogin(userID uint, username, ipAddress, userAgent string, success bool) {
	action := "POST /api/auth/login"
	if !success {
		action = "POST /api/auth/login (failed)"
	}

	log := models.AuditLog{
		UserID:    userID,
		Username:  username,
		Action:    action,
		Module:    "auth",
		Method:    "POST",
		Path:      "/api/auth/login",
		IPAddress: ipAddress,
		UserAgent: userAgent,
	}

	if success {
		log.StatusCode = 200
	} else {
		log.StatusCode = 401
	}

	models.DB.Create(&log)
}

func RecordLoginSuccess(userID uint, username, ipAddress, userAgent string) {
	s := NewAuditService()
	s.RecordLogin(userID, username, ipAddress, userAgent, true)
}

func RecordLoginFailure(username, ipAddress, userAgent string) {
	s := NewAuditService()
	s.RecordLogin(0, username, ipAddress, userAgent, false)
}
