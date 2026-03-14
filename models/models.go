package models

import (
	"fmt"
	"time"

	"contract-manage/config"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type UserRole string

const (
	RoleAdmin   UserRole = "admin"
	RoleManager UserRole = "manager"
	RoleUser    UserRole = "user"
)

type ContractStatus string

const (
	StatusDraft      ContractStatus = "draft"
	StatusPending    ContractStatus = "pending"
	StatusApproved   ContractStatus = "approved"
	StatusActive     ContractStatus = "active"
	StatusCompleted  ContractStatus = "completed"
	StatusTerminated ContractStatus = "terminated"
)

type ApprovalStatus string

const (
	ApprovalPending  ApprovalStatus = "pending"
	ApprovalApproved ApprovalStatus = "approved"
	ApprovalRejected ApprovalStatus = "rejected"
)

type User struct {
	ID              uint             `gorm:"primaryKey" json:"id"`
	Username        string           `gorm:"size:50;uniqueIndex;not null" json:"username"`
	Email           string           `gorm:"size:100;uniqueIndex" json:"email"`
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
	Code          string     `gorm:"size:50;uniqueIndex" json:"code"`
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
	ID              uint                `gorm:"primaryKey" json:"id"`
	ContractNo      string              `gorm:"size:50;uniqueIndex;not null" json:"contract_no"`
	Title           string              `gorm:"size:200;not null;index" json:"title"`
	CustomerID      uint                `gorm:"index" json:"customer_id"`
	ContractTypeID  uint                `gorm:"index" json:"contract_type_id"`
	Amount          float64             `json:"amount"`
	Currency        string              `gorm:"size:10;default:CNY" json:"currency"`
	Status          ContractStatus      `gorm:"size:20;default:draft" json:"status"`
	SignDate        *time.Time          `json:"sign_date"`
	StartDate       *time.Time          `json:"start_date"`
	EndDate         *time.Time          `json:"end_date"`
	PaymentTerms    string              `gorm:"type:text" json:"payment_terms"`
	Content         string              `gorm:"type:text" json:"content"`
	Notes           string              `gorm:"type:text" json:"notes"`
	CreatorID       uint                `gorm:"index" json:"creator_id"`
	CreatedAt       time.Time           `json:"created_at"`
	UpdatedAt       *time.Time          `json:"updated_at"`
	Customer        *Customer           `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
	Creator         *User               `gorm:"foreignKey:CreatorID" json:"creator,omitempty"`
	ContractType    *ContractType       `gorm:"foreignKey:ContractTypeID" json:"contract_type,omitempty"`
	Executions      []ContractExecution `gorm:"foreignKey:ContractID" json:"executions,omitempty"`
	Documents       []Document          `gorm:"foreignKey:ContractID" json:"documents,omitempty"`
	ApprovalRecords []ApprovalRecord    `gorm:"foreignKey:ContractID" json:"approval_records,omitempty"`
}

type ContractExecution struct {
	ID            uint       `gorm:"primaryKey" json:"id"`
	ContractID    uint       `gorm:"index" json:"contract_id"`
	Stage         string     `gorm:"size:100" json:"stage"`
	StageDate     *time.Time `json:"stage_date"`
	Progress      float64    `gorm:"default:0" json:"progress"`
	PaymentAmount float64    `json:"payment_amount"`
	PaymentDate   *time.Time `json:"payment_date"`
	Description   string     `gorm:"type:text" json:"description"`
	OperatorID    uint       `gorm:"index" json:"operator_id"`
	CreatedAt     time.Time  `json:"created_at"`
	Contract      *Contract  `gorm:"foreignKey:ContractID" json:"contract,omitempty"`
}

type ApprovalRecord struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	ContractID uint           `gorm:"index" json:"contract_id"`
	ApproverID uint           `gorm:"index" json:"approver_id"`
	Status     ApprovalStatus `gorm:"size:20;default:pending" json:"status"`
	Comment    string         `gorm:"type:text" json:"comment"`
	ApprovedAt *time.Time     `json:"approved_at"`
	CreatedAt  time.Time      `json:"created_at"`
	Contract   *Contract      `gorm:"foreignKey:ContractID" json:"contract,omitempty"`
	Approver   *User          `gorm:"foreignKey:ApproverID" json:"approver,omitempty"`
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
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	return AutoMigrate()
}

func AutoMigrate() error {
	return DB.AutoMigrate(
		&User{},
		&Role{},
		&Customer{},
		&ContractType{},
		&Contract{},
		&ContractExecution{},
		&ApprovalRecord{},
		&Document{},
		&Reminder{},
	)
}

func InitAdmin() error {
	var existingUser User
	err := DB.Where("username = ?", config.AppConfig.AdminUsername).First(&existingUser).Error

	if err == nil {
		fmt.Printf("管理员 %s 已存在\n", config.AppConfig.AdminUsername)
		return nil
	}

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
	return nil
}
