package models

import (
	"fmt"
	"strings"
	"time"

	"contract-manage/config"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type UserRole string

const (
	RoleAdmin           UserRole = "admin"            // 超级管理员
	RoleSalesManager    UserRole = "sales_manager"    // 销售负责人 - 一级审批
	RoleTechLeader      UserRole = "tech_leader"      // 技术负责人 - 二级审批
	RoleFinanceLeader   UserRole = "finance_leader"   // 财务负责人 - 三级审批
	RoleContractManager UserRole = "contract_manager" // 合同负责人 - 可查看所有合同
	RoleUser            UserRole = "user"             // 销售 - 普通用户，仅可见自己创建的
	RoleAuditAdmin      UserRole = "audit_admin"      // 审计管理员
)

// 角色层级（用于审批流程）
var RoleLevel = map[UserRole]int{
	RoleUser:            0, // 普通用户，无审批权限
	RoleSalesManager:    1, // 一级审批
	RoleTechLeader:      2, // 二级审批
	RoleFinanceLeader:   3, // 三级审批
	RoleContractManager: 0, // 合同负责人，无审批权限但可查看所有
	RoleAdmin:           0, // 管理员
	RoleAuditAdmin:      0, // 审计管理员
}

// GetRoleLevel 获取角色审批层级
func GetRoleLevel(role UserRole) int {
	if level, ok := RoleLevel[role]; ok {
		return level
	}
	return 0
}

// IsManagerRole 判断是否为审批角色
func IsManagerRole(role UserRole) bool {
	return role == RoleSalesManager || role == RoleTechLeader || role == RoleFinanceLeader
}

// CanViewAllContracts 判断是否可以查看所有合同
func CanViewAllContracts(role UserRole) bool {
	return role == RoleAdmin || role == RoleContractManager || role == RoleAuditAdmin ||
		role == RoleSalesManager || role == RoleTechLeader || role == RoleFinanceLeader
}

type ContractStatus string

const (
	StatusDraft      ContractStatus = "draft"
	StatusPending    ContractStatus = "pending"
	StatusActive     ContractStatus = "active"
	StatusCompleted  ContractStatus = "completed"
	StatusTerminated ContractStatus = "terminated"
	StatusArchived   ContractStatus = "archived"
)

const (
	StatusDraftText      = "草稿"
	StatusPendingText    = "待审批"
	StatusActiveText     = "已生效"
	StatusCompletedText  = "已完成"
	StatusTerminatedText = "已终止"
	StatusArchivedText   = "已归档"
)

func GetStatusText(status ContractStatus) string {
	switch status {
	case StatusDraft:
		return StatusDraftText
	case StatusPending:
		return StatusPendingText
	case StatusActive:
		return StatusActiveText
	case StatusCompleted:
		return StatusCompletedText
	case StatusTerminated:
		return StatusTerminatedText
	case StatusArchived:
		return StatusArchivedText
	default:
		return string(status)
	}
}

func GetStatusOptions() []map[string]string {
	return []map[string]string{
		{"value": string(StatusDraft), "label": StatusDraftText},
		{"value": string(StatusPending), "label": StatusPendingText},
		{"value": string(StatusActive), "label": StatusActiveText},
		{"value": string(StatusCompleted), "label": StatusCompletedText},
		{"value": string(StatusTerminated), "label": StatusTerminatedText},
		{"value": string(StatusArchived), "label": StatusArchivedText},
	}
}

type ApprovalStatus string

const (
	ApprovalPending  ApprovalStatus = "pending"
	ApprovalApproved ApprovalStatus = "approved"
	ApprovalRejected ApprovalStatus = "rejected"
)

type User struct {
	ID              uint             `gorm:"primaryKey" json:"id"`
	Username        string           `gorm:"size:50;uniqueIndex:idx_users_username;not null" json:"username"`
	Email           string           `gorm:"size:100;uniqueIndex:idx_users_email" json:"email"`
	HashedPassword  string           `gorm:"size:200;not null" json:"-"`
	FullName        string           `gorm:"size:100" json:"full_name"`
	Role            UserRole         `gorm:"size:20;default:user" json:"role"`
	Department      string           `gorm:"size:100" json:"department"`
	Phone           string           `gorm:"size:20" json:"phone"`
	IsActive        bool             `gorm:"default:true" json:"is_active"`
	CreatedAt       time.Time        `json:"created_at"`
	UpdatedAt       *time.Time       `json:"updated_at"`
	Contracts       []Contract       `gorm:"foreignKey:CreatorID" json:"contracts,omitempty"`
	ApprovalRecords []ApprovalRecord `gorm:"foreignKey:ApproverID" json:"approval_records,omitempty"`
}

type Role struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"size:50;unique;not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	Permissions string    `gorm:"type:text" json:"permissions"`
	CreatedAt   time.Time `json:"created_at"`
}

type Customer struct {
	ID            uint       `gorm:"primaryKey" json:"id"`
	Name          string     `gorm:"size:200;not null;index" json:"name"`
	Type          string     `gorm:"size:20;default:customer" json:"type"`
	Code          string     `gorm:"size:50;uniqueIndex:idx_customers_code" json:"code"`
	ContactPerson string     `gorm:"size:100" json:"contact_person"`
	ContactPhone  string     `gorm:"size:20" json:"contact_phone"`
	ContactEmail  string     `gorm:"size:100" json:"contact_email"`
	Address       string     `gorm:"type:text" json:"address"`
	CreditRating  string     `gorm:"size:20" json:"credit_rating"`
	IsActive      bool       `gorm:"default:true" json:"is_active"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     *time.Time `json:"updated_at"`
	Contracts     []Contract `gorm:"foreignKey:CustomerID" json:"contracts,omitempty"`
}

type ContractType struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"size:100;unique;not null" json:"name"`
	Code        string    `gorm:"size:50;unique" json:"code"`
	Description string    `gorm:"type:text" json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

type Contract struct {
	ID              uint             `gorm:"primaryKey" json:"id"`
	ContractNo      string           `gorm:"size:50;uniqueIndex:idx_contracts_contract_no;not null" json:"contract_no"`
	Title           string           `gorm:"size:200;not null;index" json:"title"`
	CustomerID      uint             `gorm:"index" json:"customer_id"`
	ContractTypeID  uint             `gorm:"index" json:"contract_type_id"`
	Amount          float64          `json:"amount"`
	Currency        string           `gorm:"size:10;default:CNY" json:"currency"`
	Status          ContractStatus   `gorm:"size:20;default:draft;index" json:"status"`
	SignDate        *time.Time       `gorm:"index" json:"sign_date"`
	StartDate       *time.Time       `gorm:"index" json:"start_date"`
	EndDate         *time.Time       `gorm:"index" json:"end_date"`
	PaymentTerms    string           `gorm:"type:text" json:"payment_terms"`
	Content         string           `gorm:"type:text" json:"content"`
	Notes           string           `gorm:"type:text" json:"notes"`
	CreatorID       uint             `gorm:"index" json:"creator_id"`
	CreatedAt       time.Time        `json:"created_at"`
	UpdatedAt       *time.Time       `json:"updated_at"`
	Customer        *Customer        `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
	Creator         *User            `gorm:"foreignKey:CreatorID" json:"creator,omitempty"`
	ContractType    *ContractType    `gorm:"foreignKey:ContractTypeID" json:"contract_type,omitempty"`
	Documents       []Document       `gorm:"foreignKey:ContractID" json:"documents,omitempty"`
	ApprovalRecords []ApprovalRecord `gorm:"foreignKey:ContractID" json:"approval_records,omitempty"`
}

type ApprovalRecord struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	ContractID   uint           `gorm:"index" json:"contract_id"`
	ApproverID   uint           `gorm:"index" json:"approver_id"`
	Level        int            `gorm:"default:1" json:"level"`
	ApproverRole string         `gorm:"size:20" json:"approver_role"`
	Status       ApprovalStatus `gorm:"size:20;default:pending" json:"status"`
	Comment      string         `gorm:"type:text" json:"comment"`
	ApprovedAt   *time.Time     `json:"approved_at"`
	CreatedAt    time.Time      `json:"created_at"`
	Contract     *Contract      `gorm:"foreignKey:ContractID" json:"contract,omitempty"`
	Approver     *User          `gorm:"foreignKey:ApproverID" json:"approver,omitempty"`
}

type Document struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	ContractID uint      `gorm:"index" json:"contract_id"`
	Name       string    `gorm:"size:200" json:"name"`
	FilePath   string    `gorm:"size:500" json:"file_path"`
	FileSize   int       `json:"file_size"`
	FileType   string    `gorm:"size:50" json:"file_type"`
	Version    string    `gorm:"size:20;default:1.0" json:"version"`
	UploaderID uint      `gorm:"index" json:"uploader_id"`
	CreatedAt  time.Time `json:"created_at"`
	Contract   *Contract `gorm:"foreignKey:ContractID" json:"contract,omitempty"`
}

type LifecycleEventType string

const (
	LifecycleCreated    LifecycleEventType = "created"
	LifecycleSubmitted  LifecycleEventType = "submitted"
	LifecycleApproved   LifecycleEventType = "approved"
	LifecycleRejected   LifecycleEventType = "rejected"
	LifecycleActivated  LifecycleEventType = "activated"
	LifecycleProgress   LifecycleEventType = "progress"
	LifecyclePayment    LifecycleEventType = "payment"
	LifecycleCompleted  LifecycleEventType = "completed"
	LifecycleTerminated LifecycleEventType = "terminated"
	LifecycleArchived   LifecycleEventType = "archived"
)

type ContractLifecycleEvent struct {
	ID          uint               `gorm:"primaryKey" json:"id"`
	ContractID  uint               `gorm:"index" json:"contract_id"`
	EventType   LifecycleEventType `gorm:"size:50" json:"event_type"`
	FromStatus  string             `gorm:"size:50" json:"from_status"`
	ToStatus    string             `gorm:"size:50" json:"to_status"`
	Amount      float64            `json:"amount"`
	Description string             `gorm:"type:text" json:"description"`
	OperatorID  uint               `gorm:"index" json:"operator_id"`
	CreatedAt   time.Time          `json:"created_at"`
	Contract    *Contract          `gorm:"foreignKey:ContractID" json:"contract,omitempty"`
}

type StatusChangeRequest struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	ContractID  uint       `gorm:"index" json:"contract_id"`
	FromStatus  string     `gorm:"size:50" json:"from_status"`
	ToStatus    string     `gorm:"size:50" json:"to_status"`
	Reason      string     `gorm:"type:text" json:"reason"`
	RequesterID uint       `gorm:"index" json:"requester_id"`
	ApproverID  *uint      `gorm:"index" json:"approver_id,omitempty"`
	Status      string     `gorm:"size:20;default:pending" json:"status"`
	Comment     string     `gorm:"type:text" json:"comment"`
	ApprovedAt  *time.Time `json:"approved_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	Contract    *Contract  `gorm:"foreignKey:ContractID" json:"contract,omitempty"`
	Requester   *User      `gorm:"foreignKey:RequesterID" json:"requester,omitempty"`
	Approver    *User      `gorm:"foreignKey:ApproverID" json:"approver,omitempty"`
}

type Reminder struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	ContractID   uint       `gorm:"index" json:"contract_id"`
	Type         string     `gorm:"size:50" json:"type"`
	ReminderDate *time.Time `json:"reminder_date"`
	DaysBefore   int        `json:"days_before"`
	IsSent       bool       `gorm:"default:false" json:"is_sent"`
	SentAt       *time.Time `json:"sent_at"`
	CreatedAt    time.Time  `json:"created_at"`
}

type Notification struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     uint      `gorm:"index" json:"user_id"`
	Type       string    `gorm:"size:50" json:"type"`
	Title      string    `gorm:"size:200" json:"title"`
	Content    string    `gorm:"type:text" json:"content"`
	ContractID *uint     `gorm:"index" json:"contract_id,omitempty"`
	IsRead     bool      `gorm:"default:false" json:"is_read"`
	CreatedAt  time.Time `json:"created_at"`
	Contract   *Contract `gorm:"foreignKey:ContractID" json:"contract,omitempty"`
}

const (
	NotificationTypeRejected     = "rejected"
	NotificationTypeApproved     = "approved"
	NotificationTypeStatusChange = "status_change"
)

type AuditLog struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     uint      `gorm:"index" json:"user_id"`
	Username   string    `gorm:"size:100" json:"username"`
	Action     string    `gorm:"size:100" json:"action"`
	Module     string    `gorm:"size:50" json:"module"`
	Method     string    `gorm:"size:20" json:"method"`
	Path       string    `gorm:"size:255" json:"path"`
	IPAddress  string    `gorm:"size:50" json:"ip_address"`
	UserAgent  string    `gorm:"type:text" json:"user_agent"`
	Request    string    `gorm:"type:text" json:"request"`
	Response   string    `gorm:"type:text" json:"response"`
	StatusCode int       `json:"status_code"`
	CreatedAt  time.Time `json:"created_at"`
	User       *User     `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

var DB *gorm.DB

func InitDB() error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.AppConfig.MysqlUser,
		config.AppConfig.MysqlPassword,
		config.AppConfig.MysqlHost,
		config.AppConfig.MysqlPort,
		config.AppConfig.MysqlDatabase,
	)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		return err
	}

	return AutoMigrate()
}

func AutoMigrate() error {
	// 禁用外键检查，避免迁移时因外键约束导致的错误
	if err := DB.Exec("SET FOREIGN_KEY_CHECKS = 0").Error; err != nil {
		return err
	}
	defer func() {
		// 恢复外键检查
		DB.Exec("SET FOREIGN_KEY_CHECKS = 1")
	}()

	// 尝试删除可能存在的旧外键约束，这些约束可能由之前的GORM版本创建
	// 但现在已经不在模型中，导致迁移时出错
	constraints := []struct {
		table      string
		constraint string
	}{
		{"users", "uni_users_username"},
		{"users", "uni_users_email"},
		{"customers", "uni_customers_code"},
		{"contracts", "uni_contracts_contract_no"},
	}

	for _, c := range constraints {
		sql := fmt.Sprintf("ALTER TABLE %s DROP FOREIGN KEY %s", c.table, c.constraint)
		if err := DB.Exec(sql).Error; err != nil {
			// 忽略错误，可能约束不存在或其他原因
			// 但如果是其他严重错误，可以打印日志
			// fmt.Printf("Warning: failed to drop constraint %s on table %s: %v\n", c.constraint, c.table, err)
		}
	}

	err := DB.AutoMigrate(
		&User{},
		&Role{},
		&Customer{},
		&ContractType{},
		&Contract{},
		&ApprovalRecord{},
		&Document{},
		&ContractLifecycleEvent{},
		&Reminder{},
		&StatusChangeRequest{},
		&AuditLog{},
		&ApprovalWorkflow{},
		&WorkflowApproval{},
		&Notification{},
	)
	if err != nil {
		// 忽略特定的迁移错误，例如尝试删除不存在的外键约束
		if strings.Contains(err.Error(), "Can't DROP") || strings.Contains(err.Error(), "Error 1091") {
			// 记录警告但继续执行
			fmt.Printf("Warning: ignored migration error: %v\n", err)
			return nil
		}
		return err
	}
	return nil
}

func InitAdmin() error {
	var existingUser User
	err := DB.Where("username = ?", config.AppConfig.AdminUsername).First(&existingUser).Error

	if err == nil {
		fmt.Printf("管理员 %s 已存在\n", config.AppConfig.AdminUsername)
	} else {
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(config.AppConfig.AdminPassword), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		admin := User{
			Username:       config.AppConfig.AdminUsername,
			Email:          config.AppConfig.AdminEmail,
			HashedPassword: string(hashedPassword),
			FullName:       "系统管理员",
			Role:           RoleAdmin,
			IsActive:       true,
		}

		if err := DB.Create(&admin).Error; err != nil {
			return err
		}
		fmt.Printf("超级管理员已创建: %s\n", config.AppConfig.AdminUsername)
	}

	var existingAuditAdmin User
	err = DB.Where("username = ?", config.AppConfig.AuditAdminUsername).First(&existingAuditAdmin).Error

	if err == nil {
		fmt.Printf("审计管理员 %s 已存在\n", config.AppConfig.AuditAdminUsername)
		return nil
	}

	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(config.AppConfig.AuditAdminPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	auditAdmin := User{
		Username:       config.AppConfig.AuditAdminUsername,
		Email:          config.AppConfig.AuditAdminEmail,
		HashedPassword: string(hashedPassword),
		FullName:       "审计管理员",
		Role:           RoleAuditAdmin,
		IsActive:       true,
	}

	if err := DB.Create(&auditAdmin).Error; err != nil {
		return err
	}
	fmt.Printf("审计管理员已创建: %s\n", config.AppConfig.AuditAdminUsername)
	return nil
}
