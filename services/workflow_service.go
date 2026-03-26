package services

import (
	"contract-manage/models"
	"time"

	"gorm.io/gorm"
)

type WorkflowService struct {
	db *gorm.DB
}

func NewWorkflowService(db *gorm.DB) *WorkflowService {
	return &WorkflowService{db: db}
}

func (s *WorkflowService) CreateWorkflow(contractID uint64, creatorRole string) (*models.ApprovalWorkflow, error) {
	workflow := &models.ApprovalWorkflow{
		ContractID:   contractID,
		CurrentLevel: 1,
		MaxLevel:     3,
		Status:       models.WorkflowStatusPending,
		CreatorRole:  creatorRole,
	}

	if err := s.db.Create(workflow).Error; err != nil {
		return nil, err
	}

	approvalList := []models.WorkflowApproval{
		{
			WorkflowID:   workflow.ID,
			ContractID:   contractID,
			ApproverRole: "sales_manager",
			Level:        1,
			Status:       models.WorkflowStatusPending,
		},
		{
			WorkflowID:   workflow.ID,
			ContractID:   contractID,
			ApproverRole: "tech_leader",
			Level:        2,
			Status:       models.WorkflowStatusPending,
		},
		{
			WorkflowID:   workflow.ID,
			ContractID:   contractID,
			ApproverRole: "finance_leader",
			Level:        3,
			Status:       models.WorkflowStatusPending,
		},
	}

	if err := s.db.Create(&approvalList).Error; err != nil {
		return nil, err
	}

	return workflow, nil
}

func (s *WorkflowService) GetWorkflowByContractID(contractID uint64) (*models.ApprovalWorkflow, error) {
	var workflow models.ApprovalWorkflow
	if err := s.db.Preload("Approvals.Approver").Where("contract_id = ?", contractID).First(&workflow).Error; err != nil {
		return nil, err
	}
	return &workflow, nil
}

func (s *WorkflowService) Approve(workflowID uint64, contractID uint64, level int, approverID uint64, comment string) error {
	var approval models.WorkflowApproval
	if err := s.db.Where("workflow_id = ? AND level = ?", workflowID, level).First(&approval).Error; err != nil {
		return err
	}

	if approval.Status != models.WorkflowStatusPending {
		return gorm.ErrRecordNotFound
	}

	now := time.Now()
	if err := s.db.Model(&approval).Updates(map[string]interface{}{
		"status":      models.WorkflowStatusApproved,
		"approver_id": approverID,
		"comment":     comment,
		"approved_at": now,
	}).Error; err != nil {
		return err
	}

	var workflow models.ApprovalWorkflow
	if err := s.db.First(&workflow, workflowID).Error; err != nil {
		return err
	}

	if level >= workflow.MaxLevel {
		s.db.Model(&workflow).Update("status", models.WorkflowStatusCompleted)

		var contract models.Contract
		if err := s.db.First(&contract, workflow.ContractID).Error; err == nil {
			s.db.Model(&contract).Updates(map[string]interface{}{
				"status": models.StatusActive,
			})
		}
	} else {
		s.db.Model(&workflow).Updates(map[string]interface{}{
			"current_level": level + 1,
			"status":        models.WorkflowStatusPending,
		})
	}

	return nil
}

func (s *WorkflowService) Reject(workflowID uint64, level int, approverID uint64, comment string) error {
	var approval models.WorkflowApproval
	if err := s.db.Where("workflow_id = ? AND level = ?", workflowID, level).First(&approval).Error; err != nil {
		return err
	}

	now := time.Now()
	if err := s.db.Model(&approval).Updates(map[string]interface{}{
		"status":      models.WorkflowStatusRejected,
		"approver_id": approverID,
		"comment":     comment,
		"approved_at": now,
	}).Error; err != nil {
		return err
	}

	s.db.Model(&models.ApprovalWorkflow{}).Where("id = ?", workflowID).Update("status", models.WorkflowStatusRejected)

	var workflow models.ApprovalWorkflow
	if err := s.db.First(&workflow, workflowID).Error; err == nil {
		s.db.Model(&models.Contract{}).Where("id = ?", workflow.ContractID).Update("status", models.StatusDraft)

		var contract models.Contract
		if err := s.db.First(&contract, workflow.ContractID).Error; err == nil {
			notificationSvc := NewNotificationService(s.db)
			reason := comment
			if reason == "" {
				reason = "无"
			}
			_ = notificationSvc.NotifyContractRejected(contract.ID, contract.Title, reason)
		}
	}

	return nil
}

func (s *WorkflowService) GetPendingApprovals(role string) ([]models.WorkflowApproval, error) {
	var approvals []models.WorkflowApproval
	if err := s.db.Preload("Approver").Preload("Workflow", "status = ?", models.WorkflowStatusPending).
		Where("approver_role = ?", role).
		Where("status = ?", models.WorkflowStatusPending).
		Find(&approvals).Error; err != nil {
		return nil, err
	}
	return approvals, nil
}

func (s *WorkflowService) GetPendingApprovalsByRoleAndLevel(role string) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	var workflows []models.ApprovalWorkflow
	if err := s.db.Where("status = ?", models.WorkflowStatusPending).Find(&workflows).Error; err != nil {
		return nil, err
	}

	for _, wf := range workflows {
		var approval models.WorkflowApproval
		if err := s.db.Where("workflow_id = ? AND approver_role = ? AND level = ? AND status = ?", wf.ID, role, wf.CurrentLevel, models.WorkflowStatusPending).First(&approval).Error; err == nil {
			results = append(results, map[string]interface{}{
				"id":            approval.ID,
				"workflow_id":   approval.WorkflowID,
				"contract_id":   approval.ContractID,
				"level":         approval.Level,
				"status":        approval.Status,
				"current_level": wf.CurrentLevel,
			})
		}
	}

	return results, nil
}

func (s *WorkflowService) GetWorkflowByWorkflowID(workflowID uint64) (*models.ApprovalWorkflow, error) {
	var workflow models.ApprovalWorkflow
	if err := s.db.Preload("Approvals").First(&workflow, workflowID).Error; err != nil {
		return nil, err
	}
	return &workflow, nil
}

func (s *WorkflowService) GetApprovalByLevel(workflowID uint64, level int) (*models.WorkflowApproval, error) {
	var approval models.WorkflowApproval
	if err := s.db.Where("workflow_id = ? AND level = ?", workflowID, level).First(&approval).Error; err != nil {
		return nil, err
	}
	return &approval, nil
}
