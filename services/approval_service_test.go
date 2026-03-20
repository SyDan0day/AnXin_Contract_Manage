package services

import (
	"contract-manage/models"
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupApprovalTestDB(t *testing.T) func() {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("failed to connect test database: %v", err)
	}

	err = db.AutoMigrate(
		&models.User{},
		&models.Customer{},
		&models.ContractType{},
		&models.Contract{},
		&models.ContractExecution{},
		&models.ApprovalRecord{},
		&models.Document{},
		&models.ContractLifecycleEvent{},
		&models.StatusChangeRequest{},
		&models.Reminder{},
		&models.AuditLog{},
	)
	if err != nil {
		t.Fatalf("failed to migrate test database: %v", err)
	}

	models.DB = db

	userService := NewUserService()
	userService.CreateUser(UserCreateInput{
		Username: "admin",
		Email:    "admin@example.com",
		Password: "password123",
		Role:     models.RoleAdmin,
	})
	userService.CreateUser(UserCreateInput{
		Username: "manager",
		Email:    "manager@example.com",
		Password: "password123",
		Role:     models.RoleManager,
	})
	userService.CreateUser(UserCreateInput{
		Username: "user",
		Email:    "user@example.com",
		Password: "password123",
	})

	customerService := NewCustomerService()
	customerService.CreateCustomer(CustomerCreateInput{
		Name: "测试客户",
		Type: "customer",
		Code: "C001",
	})

	contractType := models.ContractType{
		Name: "采购合同",
		Code: "PO001",
	}
	db.Create(&contractType)

	return func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}
}

func TestApprovalService_CreateApprovalRecord(t *testing.T) {
	cleanup := setupApprovalTestDB(t)
	defer cleanup()

	service := NewApprovalService()

	contract := models.Contract{
		ContractNo:     "APR001",
		Title:          "审批测试合同",
		CustomerID:     1,
		ContractTypeID: 1,
		CreatorID:      1,
		Status:         models.StatusPending,
	}
	models.DB.Create(&contract)

	t.Run("create approval record", func(t *testing.T) {
		record, err := service.CreateApprovalRecord(ApprovalRecordCreateInput{
			ContractID:   contract.ID,
			Level:        1,
			ApproverRole: "manager",
			Status:       "pending",
			Comment:      "请审批",
		}, 1, "manager")

		if err != nil {
			t.Errorf("CreateApprovalRecord() error = %v", err)
			return
		}
		if record.Level != 1 {
			t.Errorf("CreateApprovalRecord() level = %v, want %v", record.Level, 1)
		}
		if record.Status != models.ApprovalPending {
			t.Errorf("CreateApprovalRecord() status = %v, want %v", record.Status, models.ApprovalPending)
		}
	})

	t.Run("get approval records", func(t *testing.T) {
		records, err := service.GetApprovalRecords(contract.ID)
		if err != nil {
			t.Errorf("GetApprovalRecords() error = %v", err)
			return
		}
		if len(records) != 1 {
			t.Errorf("GetApprovalRecords() returned %v records, want %v", len(records), 1)
		}
	})

	t.Run("get approval record by ID", func(t *testing.T) {
		records, _ := service.GetApprovalRecords(contract.ID)
		record, err := service.GetApprovalRecordByID(records[0].ID)
		if err != nil {
			t.Errorf("GetApprovalRecordByID() error = %v", err)
			return
		}
		if record.ContractID != contract.ID {
			t.Errorf("GetApprovalRecordByID() contract_id = %v, want %v", record.ContractID, contract.ID)
		}
	})
}

func TestApprovalService_UpdateApprovalRecord(t *testing.T) {
	cleanup := setupApprovalTestDB(t)
	defer cleanup()

	service := NewApprovalService()

	contract := models.Contract{
		ContractNo:     "UPD001",
		Title:          "审批更新测试",
		CustomerID:     1,
		ContractTypeID: 1,
		CreatorID:      1,
		Status:         models.StatusPending,
	}
	models.DB.Create(&contract)

	record := models.ApprovalRecord{
		ContractID:   contract.ID,
		ApproverID:   1,
		Level:        1,
		ApproverRole: "admin",
		Status:       models.ApprovalPending,
		Comment:      "",
	}
	models.DB.Create(&record)

	t.Run("approve record", func(t *testing.T) {
		updated, err := service.UpdateApprovalRecord(record.ID, ApprovalRecordUpdateInput{
			Status:  "approved",
			Comment: "同意",
		}, "approved", 1)

		if err != nil {
			t.Errorf("UpdateApprovalRecord() error = %v", err)
			return
		}
		if updated.Status != models.ApprovalApproved {
			t.Errorf("UpdateApprovalRecord() status = %v, want %v", updated.Status, models.ApprovalApproved)
		}
		if updated.Comment != "同意" {
			t.Errorf("UpdateApprovalRecord() comment = %v, want %v", updated.Comment, "同意")
		}
		if updated.ApprovedAt == nil {
			t.Error("UpdateApprovalRecord() approved_at should be set")
		}
	})

	t.Run("reject already approved record", func(t *testing.T) {
		_, err := service.UpdateApprovalRecord(record.ID, ApprovalRecordUpdateInput{
			Status:  "rejected",
			Comment: "拒绝",
		}, "rejected", 1)

		if err == nil {
			t.Error("UpdateApprovalRecord() should reject already processed record")
		}
	})
}

func TestApprovalService_SubmitForApproval(t *testing.T) {
	cleanup := setupApprovalTestDB(t)
	defer cleanup()

	service := NewApprovalService()

	contract := models.Contract{
		ContractNo:     "SUB001",
		Title:          "提交审批测试",
		CustomerID:     1,
		ContractTypeID: 1,
		CreatorID:      1,
		Status:         models.StatusDraft,
	}
	models.DB.Create(&contract)

	t.Run("submit for approval", func(t *testing.T) {
		err := service.SubmitForApproval(contract.ID, 1)
		if err != nil {
			t.Errorf("SubmitForApproval() error = %v", err)
			return
		}

		var updated models.Contract
		models.DB.First(&updated, contract.ID)
		if updated.Status != models.StatusPending {
			t.Errorf("SubmitForApproval() status = %v, want %v", updated.Status, models.StatusPending)
		}
	})

	t.Run("verify lifecycle event", func(t *testing.T) {
		contractService := NewContractService()
		events, _ := contractService.GetLifecycleEvents(contract.ID)
		found := false
		for _, e := range events {
			if e.EventType == models.LifecycleSubmitted {
				found = true
				break
			}
		}
		if !found {
			t.Error("SubmitForApproval() should create submitted lifecycle event")
		}
	})
}

func TestApprovalService_GetPendingApprovals(t *testing.T) {
	cleanup := setupApprovalTestDB(t)
	defer cleanup()

	service := NewApprovalService()

	draftContract := models.Contract{
		ContractNo:     "PEND001",
		Title:          "待审批合同",
		CustomerID:     1,
		ContractTypeID: 1,
		CreatorID:      1,
		Status:         models.StatusDraft,
	}
	models.DB.Create(&draftContract)

	pendingContract := models.Contract{
		ContractNo:     "PEND002",
		Title:          "待审批合同2",
		CustomerID:     1,
		ContractTypeID: 1,
		CreatorID:      1,
		Status:         models.StatusPending,
	}
	models.DB.Create(&pendingContract)

	t.Run("get pending approvals for manager", func(t *testing.T) {
		results, err := service.GetPendingApprovalsByRole("manager", 2)
		if err != nil {
			t.Errorf("GetPendingApprovalsByRole() error = %v", err)
			return
		}
		if len(results) != 1 {
			t.Errorf("GetPendingApprovalsByRole() returned %v results, want %v", len(results), 1)
		}
	})

	t.Run("get pending approvals for admin", func(t *testing.T) {
		results, err := service.GetPendingApprovalsByRole("admin", 1)
		if err != nil {
			t.Errorf("GetPendingApprovalsByRole() error = %v", err)
			return
		}
		if len(results) != 1 {
			t.Errorf("GetPendingApprovalsByRole() returned %v results, want %v", len(results), 1)
		}
	})
}

func TestApprovalService_Reminders(t *testing.T) {
	cleanup := setupApprovalTestDB(t)
	defer cleanup()

	service := NewApprovalService()

	contract := models.Contract{
		ContractNo:     "REM001",
		Title:          "提醒测试合同",
		CustomerID:     1,
		ContractTypeID: 1,
		CreatorID:      1,
		Status:         models.StatusActive,
	}
	models.DB.Create(&contract)

	reminderDate := time.Now().AddDate(0, 0, 30)

	t.Run("create reminder", func(t *testing.T) {
		reminder, err := service.CreateReminder(ReminderCreateInput{
			ContractID:   contract.ID,
			Type:         "expiry",
			ReminderDate: &JSONTime{Time: reminderDate},
			DaysBefore:   30,
		})

		if err != nil {
			t.Errorf("CreateReminder() error = %v", err)
			return
		}
		if reminder.Type != "expiry" {
			t.Errorf("CreateReminder() type = %v, want %v", reminder.Type, "expiry")
		}
		if reminder.IsSent != false {
			t.Error("CreateReminder() is_sent should be false")
		}
	})

	t.Run("get reminders", func(t *testing.T) {
		reminders, err := service.GetReminders(contract.ID)
		if err != nil {
			t.Errorf("GetReminders() error = %v", err)
			return
		}
		if len(reminders) != 1 {
			t.Errorf("GetReminders() returned %v reminders, want %v", len(reminders), 1)
		}
	})

	t.Run("get reminder by ID", func(t *testing.T) {
		reminders, _ := service.GetReminders(contract.ID)
		reminder, err := service.GetReminderByID(reminders[0].ID)
		if err != nil {
			t.Errorf("GetReminderByID() error = %v", err)
			return
		}
		if reminder.DaysBefore != 30 {
			t.Errorf("GetReminderByID() days_before = %v, want %v", reminder.DaysBefore, 30)
		}
	})

	t.Run("mark reminder as sent", func(t *testing.T) {
		reminders, _ := service.GetReminders(contract.ID)
		err := service.UpdateReminderSent(reminders[0].ID)
		if err != nil {
			t.Errorf("UpdateReminderSent() error = %v", err)
			return
		}

		updated, _ := service.GetReminderByID(reminders[0].ID)
		if !updated.IsSent {
			t.Error("UpdateReminderSent() is_sent should be true")
		}
		if updated.SentAt == nil {
			t.Error("UpdateReminderSent() sent_at should be set")
		}
	})
}

func TestApprovalService_GetExpiringContracts(t *testing.T) {
	cleanup := setupApprovalTestDB(t)
	defer cleanup()

	service := NewApprovalService()

	now := time.Now()

	activeContract := models.Contract{
		ContractNo:     "EXP001",
		Title:          "即将到期合同",
		CustomerID:     1,
		ContractTypeID: 1,
		CreatorID:      1,
		Status:         models.StatusActive,
		EndDate:        &now,
	}
	models.DB.Create(&activeContract)

	futureContract := models.Contract{
		ContractNo:     "EXP002",
		Title:          "远期合同",
		CustomerID:     1,
		ContractTypeID: 1,
		CreatorID:      1,
		Status:         models.StatusActive,
		EndDate:        &now,
	}
	models.DB.Create(&futureContract)

	t.Run("get expiring contracts within 30 days", func(t *testing.T) {
		contracts, err := service.GetExpiringContracts(30)
		if err != nil {
			t.Errorf("GetExpiringContracts() error = %v", err)
			return
		}
		if len(contracts) != 1 {
			t.Errorf("GetExpiringContracts() returned %v contracts, want %v", len(contracts), 1)
		}
	})

	t.Run("get no expiring contracts within 1 day", func(t *testing.T) {
		contracts, err := service.GetExpiringContracts(1)
		if err != nil {
			t.Errorf("GetExpiringContracts() error = %v", err)
			return
		}
		if len(contracts) != 1 {
			t.Errorf("GetExpiringContracts() returned %v contracts, want %v", len(contracts), 1)
		}
	})
}

func TestApprovalService_GetStatistics(t *testing.T) {
	cleanup := setupApprovalTestDB(t)
	defer cleanup()

	service := NewApprovalService()

	contracts := []models.Contract{
		{
			ContractNo:     "STA001",
			Title:         "草稿合同",
			CustomerID:    1,
			ContractTypeID: 1,
			CreatorID:     1,
			Amount:        10000,
			Status:        models.StatusDraft,
		},
		{
			ContractNo:     "STA002",
			Title:         "生效合同",
			CustomerID:    1,
			ContractTypeID: 1,
			CreatorID:     1,
			Amount:        20000,
			Status:        models.StatusActive,
		},
		{
			ContractNo:     "STA003",
			Title:         "待审批合同",
			CustomerID:    1,
			ContractTypeID: 1,
			CreatorID:     1,
			Amount:        30000,
			Status:        models.StatusPending,
		},
		{
			ContractNo:     "STA004",
			Title:         "已完成合同",
			CustomerID:    1,
			ContractTypeID: 1,
			CreatorID:     1,
			Amount:        40000,
			Status:        models.StatusCompleted,
		},
	}

	for _, c := range contracts {
		models.DB.Create(&c)
	}

	t.Run("get statistics", func(t *testing.T) {
		stats, err := service.GetStatistics()
		if err != nil {
			t.Errorf("GetStatistics() error = %v", err)
			return
		}
		if stats.TotalContracts != 4 {
			t.Errorf("GetStatistics() total_contracts = %v, want %v", stats.TotalContracts, 4)
		}
		if stats.ActiveContracts != 1 {
			t.Errorf("GetStatistics() active_contracts = %v, want %v", stats.ActiveContracts, 1)
		}
		if stats.PendingContracts != 1 {
			t.Errorf("GetStatistics() pending_contracts = %v, want %v", stats.PendingContracts, 1)
		}
		if stats.CompletedContracts != 1 {
			t.Errorf("GetStatistics() completed_contracts = %v, want %v", stats.CompletedContracts, 1)
		}
		if stats.TotalAmount != 100000 {
			t.Errorf("GetStatistics() total_amount = %v, want %v", stats.TotalAmount, 100000.0)
		}
	})
}

func TestApprovalService_GetPendingStatusChangesCount(t *testing.T) {
	cleanup := setupApprovalTestDB(t)
	defer cleanup()

	service := NewApprovalService()
	contractService := NewContractService()

	contract := models.Contract{
		ContractNo:     "CNT001",
		Title:         "状态变更计数测试",
		CustomerID:    1,
		ContractTypeID: 1,
		CreatorID:     1,
		Status:        models.StatusActive,
	}
	models.DB.Create(&contract)

	contractService.CreateStatusChangeRequest(contract.ID, StatusChangeRequestInput{
		ToStatus: "archived",
		Reason:   "测试",
	}, 1)

	t.Run("get pending status changes count", func(t *testing.T) {
		count, err := service.GetPendingStatusChangesCount()
		if err != nil {
			t.Errorf("GetPendingStatusChangesCount() error = %v", err)
			return
		}
		if count != 1 {
			t.Errorf("GetPendingStatusChangesCount() = %v, want %v", count, 1)
		}
	})
}
