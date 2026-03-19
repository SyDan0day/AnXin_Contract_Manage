package services

import (
	"contract-manage/models"
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupContractTestDB(t *testing.T) func() {
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
		Username: "testuser",
		Email:    "test@example.com",
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
	models.DB.Create(&contractType)

	return func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}
}

func TestContractService_CreateContract(t *testing.T) {
	cleanup := setupContractTestDB(t)
	defer cleanup()

	service := NewContractService()

	t.Run("create contract with all fields", func(t *testing.T) {
		startDate := time.Now()
		endDate := startDate.AddDate(1, 0, 0)
		signDate := startDate

		contract, err := service.CreateContract(ContractCreateInput{
			Title:          "测试采购合同",
			CustomerID:     1,
			ContractTypeID: 1,
			Amount:         100000.00,
			Currency:       "CNY",
			SignDate:       &JSONTime{Time: signDate},
			StartDate:      &JSONTime{Time: startDate},
			EndDate:        &JSONTime{Time: endDate},
			PaymentTerms:   "预付30%",
			Content:        "合同内容...",
			Notes:          "备注",
		}, 1)

		if err != nil {
			t.Errorf("CreateContract() error = %v", err)
			return
		}
		if contract.Title != "测试采购合同" {
			t.Errorf("CreateContract() title = %v, want %v", contract.Title, "测试采购合同")
		}
		if contract.Amount != 100000.00 {
			t.Errorf("CreateContract() amount = %v, want %v", contract.Amount, 100000.00)
		}
		if contract.Status != models.StatusDraft {
			t.Errorf("CreateContract() status = %v, want %v", contract.Status, models.StatusDraft)
		}
		if contract.ContractNo == "" {
			t.Error("CreateContract() contract_no should not be empty")
		}
	})

	t.Run("create contract with default currency", func(t *testing.T) {
		contract, err := service.CreateContract(ContractCreateInput{
			Title:          "无货币合同",
			CustomerID:     1,
			ContractTypeID: 1,
			Amount:         50000.00,
		}, 1)

		if err != nil {
			t.Errorf("CreateContract() error = %v", err)
			return
		}
		if contract.Currency != "CNY" {
			t.Errorf("CreateContract() currency = %v, want %v", contract.Currency, "CNY")
		}
	})

	t.Run("verify lifecycle event created", func(t *testing.T) {
		events, err := service.GetLifecycleEvents(1)
		if err != nil {
			t.Errorf("GetLifecycleEvents() error = %v", err)
			return
		}
		if len(events) == 0 {
			t.Error("CreateContract() should create lifecycle event")
		}
	})
}

func TestContractService_GetContractByID(t *testing.T) {
	cleanup := setupContractTestDB(t)
	defer cleanup()

	service := NewContractService()

	contract := models.Contract{
		ContractNo:     "CT001",
		Title:          "获取测试合同",
		CustomerID:     1,
		ContractTypeID: 1,
		Amount:         10000.00,
		CreatorID:      1,
		Status:         models.StatusDraft,
	}
	models.DB.Create(&contract)

	t.Run("get existing contract", func(t *testing.T) {
		found, err := service.GetContractByID(contract.ID)
		if err != nil {
			t.Errorf("GetContractByID() error = %v", err)
			return
		}
		if found.Title != "获取测试合同" {
			t.Errorf("GetContractByID() title = %v, want %v", found.Title, "获取测试合同")
		}
	})

	t.Run("get non-existing contract", func(t *testing.T) {
		_, err := service.GetContractByID(9999)
		if err == nil {
			t.Error("GetContractByID() expected error for non-existing contract")
		}
	})
}

func TestContractService_GetContracts(t *testing.T) {
	cleanup := setupContractTestDB(t)
	defer cleanup()

	service := NewContractService()

	for i := 0; i < 5; i++ {
		models.DB.Create(&models.Contract{
			ContractNo:     "CT00" + string(rune('1'+i)),
			Title:          "合同" + string(rune('A'+i)),
			CustomerID:     1,
			ContractTypeID: 1,
			Amount:         float64(10000 * (i + 1)),
			CreatorID:      1,
			Status:         models.StatusDraft,
		})
	}

	t.Run("get all contracts", func(t *testing.T) {
		contracts, err := service.GetContracts(0, 10, 0, 0, "")
		if err != nil {
			t.Errorf("GetContracts() error = %v", err)
			return
		}
		if len(contracts) != 5 {
			t.Errorf("GetContracts() returned %v contracts, want %v", len(contracts), 5)
		}
	})

	t.Run("filter by customer", func(t *testing.T) {
		contracts, err := service.GetContracts(0, 10, 1, 0, "")
		if err != nil {
			t.Errorf("GetContracts() error = %v", err)
			return
		}
		if len(contracts) != 5 {
			t.Errorf("GetContracts() returned %v contracts, want %v", len(contracts), 5)
		}
	})

	t.Run("pagination", func(t *testing.T) {
		contracts, err := service.GetContracts(0, 2, 0, 0, "")
		if err != nil {
			t.Errorf("GetContracts() error = %v", err)
			return
		}
		if len(contracts) != 2 {
			t.Errorf("GetContracts() returned %v contracts, want %v", len(contracts), 2)
		}
	})
}

func TestContractService_UpdateContract(t *testing.T) {
	cleanup := setupContractTestDB(t)
	defer cleanup()

	service := NewContractService()

	contract := models.Contract{
		ContractNo:     "UP001",
		Title:          "原始标题",
		CustomerID:     1,
		ContractTypeID: 1,
		Amount:         50000.00,
		CreatorID:      1,
		Status:         models.StatusDraft,
	}
	models.DB.Create(&contract)

	t.Run("update contract fields", func(t *testing.T) {
		updated, err := service.UpdateContract(contract.ID, ContractUpdateInput{
			Title:  "更新后的标题",
			Amount: 60000.00,
		})
		if err != nil {
			t.Errorf("UpdateContract() error = %v", err)
			return
		}
		if updated.Title != "更新后的标题" {
			t.Errorf("UpdateContract() title = %v, want %v", updated.Title, "更新后的标题")
		}
		if updated.Amount != 60000.00 {
			t.Errorf("UpdateContract() amount = %v, want %v", updated.Amount, 60000.00)
		}
	})

	t.Run("update non-existing contract", func(t *testing.T) {
		_, err := service.UpdateContract(9999, ContractUpdateInput{Title: "标题"})
		if err == nil {
			t.Error("UpdateContract() expected error for non-existing contract")
		}
	})
}

func TestContractService_DeleteContract(t *testing.T) {
	cleanup := setupContractTestDB(t)
	defer cleanup()

	service := NewContractService()

	contract := models.Contract{
		ContractNo:     "DEL001",
		Title:          "待删除合同",
		CustomerID:     1,
		ContractTypeID: 1,
		CreatorID:      1,
	}
	models.DB.Create(&contract)

	t.Run("delete existing contract", func(t *testing.T) {
		err := service.DeleteContract(contract.ID)
		if err != nil {
			t.Errorf("DeleteContract() error = %v", err)
			return
		}
		_, err = service.GetContractByID(contract.ID)
		if err == nil {
			t.Error("DeleteContract() contract still exists after deletion")
		}
	})

	t.Run("delete non-existing contract", func(t *testing.T) {
		err := service.DeleteContract(9999)
		if err == nil {
			t.Error("DeleteContract() expected error for non-existing contract")
		}
	})
}

func TestContractService_UpdateContractStatus(t *testing.T) {
	cleanup := setupContractTestDB(t)
	defer cleanup()

	service := NewContractService()

	contract := models.Contract{
		ContractNo:     "ST001",
		Title:          "状态变更测试",
		CustomerID:     1,
		ContractTypeID: 1,
		CreatorID:      1,
		Status:         models.StatusDraft,
	}
	models.DB.Create(&contract)

	t.Run("update contract status", func(t *testing.T) {
		updated, err := service.UpdateContractStatus(contract.ID, "active", 1)
		if err != nil {
			t.Errorf("UpdateContractStatus() error = %v", err)
			return
		}
		if updated.Status != models.StatusActive {
			t.Errorf("UpdateContractStatus() status = %v, want %v", updated.Status, models.StatusActive)
		}
	})

	t.Run("verify lifecycle event", func(t *testing.T) {
		events, err := service.GetLifecycleEvents(contract.ID)
		if err != nil {
			t.Errorf("GetLifecycleEvents() error = %v", err)
			return
		}
		found := false
		for _, e := range events {
			if e.EventType == "status_changed" {
				found = true
				break
			}
		}
		if !found {
			t.Error("UpdateContractStatus() should create status_changed lifecycle event")
		}
	})
}

func TestContractService_ArchiveContract(t *testing.T) {
	cleanup := setupContractTestDB(t)
	defer cleanup()

	service := NewContractService()

	contract := models.Contract{
		ContractNo:     "ARC001",
		Title:          "归档测试合同",
		CustomerID:     1,
		ContractTypeID: 1,
		CreatorID:      1,
		Status:         models.StatusCompleted,
	}
	models.DB.Create(&contract)

	t.Run("archive contract", func(t *testing.T) {
		archived, err := service.ArchiveContract(contract.ID, 1)
		if err != nil {
			t.Errorf("ArchiveContract() error = %v", err)
			return
		}
		if archived.Status != models.StatusArchived {
			t.Errorf("ArchiveContract() status = %v, want %v", archived.Status, models.StatusArchived)
		}
	})

	t.Run("verify lifecycle event", func(t *testing.T) {
		events, _ := service.GetLifecycleEvents(contract.ID)
		found := false
		for _, e := range events {
			if e.EventType == models.LifecycleArchived {
				found = true
				break
			}
		}
		if !found {
			t.Error("ArchiveContract() should create archived lifecycle event")
		}
	})
}

func TestContractService_IsStatusChangeRequireApproval(t *testing.T) {
	service := NewContractService()

	tests := []struct {
		status   string
		expected bool
	}{
		{"archived", true},
		{"terminated", true},
		{"in_progress", true},
		{"pending_pay", true},
		{"draft", false},
		{"active", false},
		{"completed", false},
	}

	for _, tt := range tests {
		t.Run(tt.status, func(t *testing.T) {
			result := service.IsStatusChangeRequireApproval(tt.status)
			if result != tt.expected {
				t.Errorf("IsStatusChangeRequireApproval(%v) = %v, want %v", tt.status, result, tt.expected)
			}
		})
	}
}

func TestContractService_StatusChangeRequest(t *testing.T) {
	cleanup := setupContractTestDB(t)
	defer cleanup()

	service := NewContractService()

	contract := models.Contract{
		ContractNo:     "CR001",
		Title:          "状态变更申请测试",
		CustomerID:     1,
		ContractTypeID: 1,
		CreatorID:      1,
		Status:         models.StatusActive,
	}
	models.DB.Create(&contract)

	t.Run("create status change request for archived", func(t *testing.T) {
		request, err := service.CreateStatusChangeRequest(contract.ID, StatusChangeRequestInput{
			ToStatus: "archived",
			Reason:   "合同已完成",
		}, 1)

		if err != nil {
			t.Errorf("CreateStatusChangeRequest() error = %v", err)
			return
		}
		if request == nil {
			t.Error("CreateStatusChangeRequest() should return request")
			return
		}
		if request.ToStatus != "archived" {
			t.Errorf("CreateStatusChangeRequest() to_status = %v, want %v", request.ToStatus, "archived")
		}
		if request.Status != "pending" {
			t.Errorf("CreateStatusChangeRequest() status = %v, want %v", request.Status, "pending")
		}
	})

	t.Run("status change that doesn't need approval", func(t *testing.T) {
		request, err := service.CreateStatusChangeRequest(contract.ID, StatusChangeRequestInput{
			ToStatus: "draft",
			Reason:   "退回草稿",
		}, 1)

		if err != nil {
			t.Errorf("CreateStatusChangeRequest() error = %v", err)
			return
		}
		if request != nil {
			t.Error("CreateStatusChangeRequest() should return nil for non-approval status")
		}
	})

	t.Run("prevent duplicate pending request", func(t *testing.T) {
		_, err := service.CreateStatusChangeRequest(contract.ID, StatusChangeRequestInput{
			ToStatus: "terminated",
			Reason:   "测试",
		}, 1)

		if err == nil {
			t.Error("CreateStatusChangeRequest() should prevent duplicate pending requests")
		}
	})
}

func TestContractService_LifecycleEvents(t *testing.T) {
	cleanup := setupContractTestDB(t)
	defer cleanup()

	service := NewContractService()

	contract := models.Contract{
		ContractNo:     "LF001",
		Title:          "生命周期测试",
		CustomerID:     1,
		ContractTypeID: 1,
		CreatorID:      1,
	}
	models.DB.Create(&contract)

	t.Run("add lifecycle event", func(t *testing.T) {
		event, err := service.AddLifecycleEvent(contract.ID, LifecycleEventInput{
			EventType:   "progress",
			Description: "执行进度50%",
			Amount:      50000.00,
		}, 1)

		if err != nil {
			t.Errorf("AddLifecycleEvent() error = %v", err)
			return
		}
		if event.EventType != "progress" {
			t.Errorf("AddLifecycleEvent() event_type = %v, want %v", event.EventType, "progress")
		}
	})

	t.Run("get lifecycle events", func(t *testing.T) {
		events, err := service.GetLifecycleEvents(contract.ID)
		if err != nil {
			t.Errorf("GetLifecycleEvents() error = %v", err)
			return
		}
		if len(events) == 0 {
			t.Error("GetLifecycleEvents() should return events")
		}
	})
}

func TestContractService_ContractExecutions(t *testing.T) {
	cleanup := setupContractTestDB(t)
	defer cleanup()

	service := NewContractService()

	contract := models.Contract{
		ContractNo:     "EX001",
		Title:          "执行跟踪测试",
		CustomerID:     1,
		ContractTypeID: 1,
		CreatorID:      1,
	}
	models.DB.Create(&contract)

	t.Run("create contract execution", func(t *testing.T) {
		execution, err := service.CreateContractExecution(ContractExecutionCreateInput{
			ContractID:    contract.ID,
			Stage:         "阶段一",
			Progress:      50.0,
			PaymentAmount: 50000.00,
			Description:   "已完成阶段一",
		}, 1)

		if err != nil {
			t.Errorf("CreateContractExecution() error = %v", err)
			return
		}
		if execution.Stage != "阶段一" {
			t.Errorf("CreateContractExecution() stage = %v, want %v", execution.Stage, "阶段一")
		}
	})

	t.Run("get contract executions", func(t *testing.T) {
		executions, err := service.GetContractExecutions(contract.ID)
		if err != nil {
			t.Errorf("GetContractExecutions() error = %v", err)
			return
		}
		if len(executions) != 1 {
			t.Errorf("GetContractExecutions() returned %v executions, want %v", len(executions), 1)
		}
	})

	t.Run("delete execution", func(t *testing.T) {
		executions, _ := service.GetContractExecutions(contract.ID)
		err := service.DeleteExecution(executions[0].ID)
		if err != nil {
			t.Errorf("DeleteExecution() error = %v", err)
		}
	})
}

func TestContractService_Documents(t *testing.T) {
	cleanup := setupContractTestDB(t)
	defer cleanup()

	service := NewContractService()

	contract := models.Contract{
		ContractNo:     "DOC001",
		Title:          "文档管理测试",
		CustomerID:     1,
		ContractTypeID: 1,
		CreatorID:      1,
	}
	models.DB.Create(&contract)

	t.Run("create document", func(t *testing.T) {
		doc, err := service.CreateDocument(DocumentCreateInput{
			ContractID: contract.ID,
			Name:       "合同附件.pdf",
			FilePath:   "/uploads/contract.pdf",
			FileSize:   1024,
			FileType:   "application/pdf",
		}, 1)

		if err != nil {
			t.Errorf("CreateDocument() error = %v", err)
			return
		}
		if doc.Name != "合同附件.pdf" {
			t.Errorf("CreateDocument() name = %v, want %v", doc.Name, "合同附件.pdf")
		}
	})

	t.Run("get documents", func(t *testing.T) {
		docs, err := service.GetDocuments(contract.ID)
		if err != nil {
			t.Errorf("GetDocuments() error = %v", err)
			return
		}
		if len(docs) != 1 {
			t.Errorf("GetDocuments() returned %v documents, want %v", len(docs), 1)
		}
	})

	t.Run("get document by ID", func(t *testing.T) {
		docs, _ := service.GetDocuments(contract.ID)
		doc, err := service.GetDocumentByID(docs[0].ID)
		if err != nil {
			t.Errorf("GetDocumentByID() error = %v", err)
			return
		}
		if doc.Name != "合同附件.pdf" {
			t.Errorf("GetDocumentByID() name = %v, want %v", doc.Name, "合同附件.pdf")
		}
	})

	t.Run("delete document", func(t *testing.T) {
		docs, _ := service.GetDocuments(contract.ID)
		err := service.DeleteDocument(docs[0].ID)
		if err != nil {
			t.Errorf("DeleteDocument() error = %v", err)
		}
	})
}

func TestJSONTime_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected time.Time
		wantErr  bool
	}{
		{"RFC3339 format", `"2024-01-01T10:00:00Z"`, time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC), false},
		{"datetime format", `"2024-01-01 10:00:00"`, time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC), false},
		{"date format", `"2024-01-01"`, time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), false},
		{"empty string", `""`, time.Time{}, false},
		{"null", `"null"`, time.Time{}, false},
		{"invalid format", `"invalid"`, time.Time{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var jt JSONTime
			err := jt.UnmarshalJSON([]byte(tt.input))
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && tt.input != `"null"` && tt.input != `""` {
				if !jt.Time.Equal(tt.expected) && jt.Time.Year() != tt.expected.Year() {
					t.Errorf("UnmarshalJSON() time = %v, want %v", jt.Time, tt.expected)
				}
			}
		})
	}
}
