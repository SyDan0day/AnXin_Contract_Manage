package models

import (
	"time"
)

type ApprovalWorkflow struct {
	ID           uint64    `json:"id" gorm:"primaryKey"`
	ContractID   uint64    `json:"contract_id" gorm:"index;not null"`
	CurrentLevel int       `json:"current_level" gorm:"default:1"`
	MaxLevel     int       `json:"max_level" gorm:"default:2"`
	Status       string    `json:"status" gorm:"type:varchar(20);default:'pending';index"`
	CreatorRole  string    `json:"creator_role" gorm:"type:varchar(20);not null"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	
	Contract    Contract           `json:"contract,omitempty" gorm:"foreignKey:ContractID"`
	Approvals  []WorkflowApproval `json:"approvals,omitempty" gorm:"foreignKey:WorkflowID"`
}

type WorkflowApproval struct {
	ID           uint64     `json:"id" gorm:"primaryKey"`
	WorkflowID   uint64     `json:"workflow_id" gorm:"index;not null"`
	ContractID   uint64     `json:"contract_id" gorm:"index;not null"`
	ApproverID   *uint64    `json:"approver_id"`
	ApproverRole string    `json:"approver_role" gorm:"type:varchar(20);not null"`
	Level        int        `json:"level" gorm:"not null"`
	Status       string     `json:"status" gorm:"type:varchar(20);default:'pending'"`
	Comment      string     `json:"comment" gorm:"type:text"`
	ApprovedAt   *time.Time `json:"approved_at"`
	CreatedAt    time.Time `json:"created_at"`
	
	Approver User `json:"approver,omitempty" gorm:"foreignKey:ApproverID"`
}

const (
	WorkflowStatusPending   = "pending"
	WorkflowStatusApproved  = "approved"
	WorkflowStatusRejected  = "rejected"
	WorkflowStatusCompleted = "completed"
	
	WorkflowLevel1 = 1
	WorkflowLevel2 = 2
)

var ApprovalRoles = map[string]int{
	"sales":         0,
	"admin":         1,
	"director":      2,
	"super_admin":  3,
}
