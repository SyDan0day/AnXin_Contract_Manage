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
		MaxLevel:     2,
		Status:       models.WorkflowStatusPending,
		CreatorRole:  creatorRole,
	}
	
	if err := s.db.Create(workflow).Error; err != nil {
		return nil, err
	}
	
	approvers := []models.WorkflowApproval{
		{
			WorkflowID:   workflow.ID,
			ContractID:   contractID,
			ApproverRole: "admin",
			Level:        1,
			Status:       models.WorkflowStatusPending,
		},
		{
			WorkflowID:   workflow.ID,
			ContractID:   contractID,
			ApproverRole: "director",
			Level:        2,
			Status:       models.WorkflowStatusPending,
		},
	}
	
	if err := s.db.Create(&approvers).Error; err != nil {
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
		"status":     models.WorkflowStatusApproved,
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
	} else {
		s.db.Model(&workflow).Updates(map[string]interface{}{
			"current_level": level + 1,
			"status": models.WorkflowStatusPending,
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
		"status":     models.WorkflowStatusRejected,
		"approver_id": approverID,
		"comment":     comment,
		"approved_at": now,
	}).Error; err != nil {
		return err
	}
	
	s.db.Model(&models.ApprovalWorkflow{}).Where("id = ?", workflowID).Update("status", models.WorkflowStatusRejected)
	
	return nil
}

func (s *WorkflowService) GetPendingApprovals(role string) ([]models.WorkflowApproval, error) {
	var approvals []models.WorkflowApproval
	if err := s.db.Preload("Approver").Preload("Workflow").
		Joins("JOIN approval_workflows ON approval_workflows.id = workflow_approvals.workflow_id").
		Where("workflow_approvals.approver_role = ?", role).
		Where("workflow_approvals.level = approval_workflows.current_level").
		Where("approval_workflows.status = ?", models.WorkflowStatusPending).
		Where("workflow_approvals.status = ?", models.WorkflowStatusPending).
		Find(&approvals).Error; err != nil {
		return nil, err
	}
	return approvals, nil
}
