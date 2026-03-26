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

	Contract  Contract           `json:"contract,omitempty" gorm:"foreignKey:ContractID"`
	Approvals []WorkflowApproval `json:"approvals,omitempty" gorm:"foreignKey:WorkflowID"`
}

type WorkflowApproval struct {
	ID           uint64     `json:"id" gorm:"primaryKey"`
	WorkflowID   uint64     `json:"workflow_id" gorm:"index;not null"`
	ContractID   uint64     `json:"contract_id" gorm:"index;not null"`
	ApproverID   *uint64    `json:"approver_id"`
	ApproverRole string     `json:"approver_role" gorm:"type:varchar(20);not null"`
	Level        int        `json:"level" gorm:"not null"`
	Status       string     `json:"status" gorm:"type:varchar(20);default:'pending'"`
	Comment      string     `json:"comment" gorm:"type:text"`
	ApprovedAt   *time.Time `json:"approved_at"`
	CreatedAt    time.Time  `json:"created_at"`

	Approver User `json:"approver,omitempty" gorm:"foreignKey:ApproverID"`
}

const (
	WorkflowStatusPending   = "pending"
	WorkflowStatusApproved  = "approved"
	WorkflowStatusRejected  = "rejected"
	WorkflowStatusCompleted = "completed"
	WorkflowStatusCancelled = "cancelled"

	WorkflowLevel1 = 1
	WorkflowLevel2 = 2
	WorkflowLevel3 = 3
)

// 审批角色层级（数字越大，级别越高）
var ApprovalRoles = map[string]int{
	"user":             0, // 普通用户（销售），无审批权限
	"sales_manager":    1, // 销售负责人 - 一级审批
	"tech_leader":      2, // 技术负责人 - 二级审批
	"finance_leader":   3, // 财务负责人 - 三级审批
	"contract_manager": 0, // 合同负责人，无审批权限
	"admin":            4, // 管理员，可审批所有
	"audit_admin":      0, // 审计管理员，无审批权限
}

// GetApprovalLevel 获取审批层级
func GetApprovalLevel(role string) int {
	if level, ok := ApprovalRoles[role]; ok {
		return level
	}
	return 0
}

// CanApprove 判断角色是否可以审批
func CanApprove(role string) bool {
	level := GetApprovalLevel(role)
	return level > 0
}
