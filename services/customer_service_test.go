package services

import (
	"contract-manage/models"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupCustomerTestDB(t *testing.T) func() {
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

	return func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}
}

func TestCustomerService_CreateCustomer(t *testing.T) {
	cleanup := setupCustomerTestDB(t)
	defer cleanup()

	service := NewCustomerService()

	tests := []struct {
		name    string
		input   CustomerCreateInput
		wantErr bool
	}{
		{
			name: "valid customer creation",
			input: CustomerCreateInput{
				Name:          "测试客户",
				Type:          "customer",
				Code:          "C001",
				ContactPerson: "张三",
				ContactPhone:  "13800138000",
				ContactEmail:  "zhangsan@example.com",
				CreditRating:  "A",
			},
			wantErr: false,
		},
		{
			name: "duplicate customer code",
			input: CustomerCreateInput{
				Name: "另一个客户",
				Type: "customer",
				Code: "C001",
			},
			wantErr: true,
		},
		{
			name: "supplier customer",
			input: CustomerCreateInput{
				Name: "测试供应商",
				Type: "supplier",
				Code: "S001",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			customer, err := service.CreateCustomer(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateCustomer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && customer == nil {
				t.Error("CreateCustomer() returned nil customer without error")
			}
		})
	}
}

func TestCustomerService_GetCustomerByID(t *testing.T) {
	cleanup := setupCustomerTestDB(t)
	defer cleanup()

	service := NewCustomerService()

	created, err := service.CreateCustomer(CustomerCreateInput{
		Name: "获取测试客户",
		Type: "customer",
		Code: "GET001",
	})
	if err != nil {
		t.Fatalf("failed to create test customer: %v", err)
	}

	t.Run("get existing customer", func(t *testing.T) {
		customer, err := service.GetCustomerByID(created.ID)
		if err != nil {
			t.Errorf("GetCustomerByID() error = %v", err)
			return
		}
		if customer.Name != "获取测试客户" {
			t.Errorf("GetCustomerByID() name = %v, want %v", customer.Name, "获取测试客户")
		}
	})

	t.Run("get non-existing customer", func(t *testing.T) {
		_, err := service.GetCustomerByID(9999)
		if err == nil {
			t.Error("GetCustomerByID() expected error for non-existing customer")
		}
	})
}

func TestCustomerService_GetCustomerByCode(t *testing.T) {
	cleanup := setupCustomerTestDB(t)
	defer cleanup()

	service := NewCustomerService()

	_, err := service.CreateCustomer(CustomerCreateInput{
		Name: "编码测试客户",
		Type: "customer",
		Code: "CODE001",
	})
	if err != nil {
		t.Fatalf("failed to create test customer: %v", err)
	}

	t.Run("find by existing code", func(t *testing.T) {
		customer, err := service.GetCustomerByCode("CODE001")
		if err != nil {
			t.Errorf("GetCustomerByCode() error = %v", err)
			return
		}
		if customer.Name != "编码测试客户" {
			t.Errorf("GetCustomerByCode() name = %v, want %v", customer.Name, "编码测试客户")
		}
	})

	t.Run("find by non-existing code", func(t *testing.T) {
		_, err := service.GetCustomerByCode("NONEXISTENT")
		if err == nil {
			t.Error("GetCustomerByCode() expected error for non-existing code")
		}
	})
}

func TestCustomerService_GetCustomers(t *testing.T) {
	cleanup := setupCustomerTestDB(t)
	defer cleanup()

	service := NewCustomerService()

	service.CreateCustomer(CustomerCreateInput{Name: "客户A", Type: "customer", Code: "CA001"})
	service.CreateCustomer(CustomerCreateInput{Name: "客户B", Type: "customer", Code: "CA002"})
	service.CreateCustomer(CustomerCreateInput{Name: "供应商A", Type: "supplier", Code: "SA001"})

	t.Run("get all customers", func(t *testing.T) {
		customers, err := service.GetCustomers(0, 10, "")
		if err != nil {
			t.Errorf("GetCustomers() error = %v", err)
			return
		}
		if len(customers) != 3 {
			t.Errorf("GetCustomers() returned %v customers, want %v", len(customers), 3)
		}
	})

	t.Run("filter by customer type", func(t *testing.T) {
		customers, err := service.GetCustomers(0, 10, "customer")
		if err != nil {
			t.Errorf("GetCustomers() error = %v", err)
			return
		}
		if len(customers) != 2 {
			t.Errorf("GetCustomers() returned %v customers, want %v", len(customers), 2)
		}
	})

	t.Run("filter by supplier type", func(t *testing.T) {
		customers, err := service.GetCustomers(0, 10, "supplier")
		if err != nil {
			t.Errorf("GetCustomers() error = %v", err)
			return
		}
		if len(customers) != 1 {
			t.Errorf("GetCustomers() returned %v customers, want %v", len(customers), 1)
		}
	})

	t.Run("pagination", func(t *testing.T) {
		customers, err := service.GetCustomers(0, 2, "")
		if err != nil {
			t.Errorf("GetCustomers() error = %v", err)
			return
		}
		if len(customers) != 2 {
			t.Errorf("GetCustomers() returned %v customers, want %v", len(customers), 2)
		}
	})
}

func TestCustomerService_UpdateCustomer(t *testing.T) {
	cleanup := setupCustomerTestDB(t)
	defer cleanup()

	service := NewCustomerService()

	created, _ := service.CreateCustomer(CustomerCreateInput{
		Name:         "原始名称",
		Type:         "customer",
		Code:         "UP001",
		ContactPhone: "13800000000",
	})

	t.Run("update customer fields", func(t *testing.T) {
		updated, err := service.UpdateCustomer(created.ID, CustomerUpdateInput{
			Name:         "新名称",
			ContactPhone: "13900000000",
		})
		if err != nil {
			t.Errorf("UpdateCustomer() error = %v", err)
			return
		}
		if updated.Name != "新名称" {
			t.Errorf("UpdateCustomer() name = %v, want %v", updated.Name, "新名称")
		}
		if updated.ContactPhone != "13900000000" {
			t.Errorf("UpdateCustomer() contact_phone = %v, want %v", updated.ContactPhone, "13900000000")
		}
	})

	t.Run("update non-existing customer", func(t *testing.T) {
		_, err := service.UpdateCustomer(9999, CustomerUpdateInput{Name: "名称"})
		if err == nil {
			t.Error("UpdateCustomer() expected error for non-existing customer")
		}
	})

	t.Run("deactivate customer", func(t *testing.T) {
		isActive := false
		updated, err := service.UpdateCustomer(created.ID, CustomerUpdateInput{
			IsActive: &isActive,
		})
		if err != nil {
			t.Errorf("UpdateCustomer() error = %v", err)
			return
		}
		if updated.IsActive != false {
			t.Error("UpdateCustomer() is_active should be false")
		}
	})
}

func TestCustomerService_DeleteCustomer(t *testing.T) {
	cleanup := setupCustomerTestDB(t)
	defer cleanup()

	service := NewCustomerService()

	created, _ := service.CreateCustomer(CustomerCreateInput{
		Name: "删除测试客户",
		Type: "customer",
		Code: "DEL001",
	})

	t.Run("delete existing customer", func(t *testing.T) {
		err := service.DeleteCustomer(created.ID)
		if err != nil {
			t.Errorf("DeleteCustomer() error = %v", err)
			return
		}
		_, err = service.GetCustomerByID(created.ID)
		if err == nil {
			t.Error("DeleteCustomer() customer still exists after deletion")
		}
	})

	t.Run("delete non-existing customer", func(t *testing.T) {
		err := service.DeleteCustomer(9999)
		if err == nil {
			t.Error("DeleteCustomer() expected error for non-existing customer")
		}
	})
}

func TestContractTypeService(t *testing.T) {
	cleanup := setupCustomerTestDB(t)
	defer cleanup()

	service := NewCustomerService()

	t.Run("create contract type", func(t *testing.T) {
		ct, err := service.CreateContractType(ContractTypeCreateInput{
			Name:        "采购合同",
			Code:        "PO001",
			Description: "采购相关合同",
		})
		if err != nil {
			t.Errorf("CreateContractType() error = %v", err)
			return
		}
		if ct.Name != "采购合同" {
			t.Errorf("CreateContractType() name = %v, want %v", ct.Name, "采购合同")
		}
	})

	t.Run("get contract types", func(t *testing.T) {
		types, err := service.GetContractTypes(0, 10)
		if err != nil {
			t.Errorf("GetContractTypes() error = %v", err)
			return
		}
		if len(types) != 1 {
			t.Errorf("GetContractTypes() returned %v types, want %v", len(types), 1)
		}
	})

	t.Run("update contract type", func(t *testing.T) {
		ct, _ := service.CreateContractType(ContractTypeCreateInput{
			Name: "销售合同",
			Code: "SO001",
		})
		updated, err := service.UpdateContractType(ct.ID, ContractTypeCreateInput{
			Name:        "销售合同V2",
			Code:        "SO002",
			Description: "更新后的销售合同",
		})
		if err != nil {
			t.Errorf("UpdateContractType() error = %v", err)
			return
		}
		if updated.Name != "销售合同V2" {
			t.Errorf("UpdateContractType() name = %v, want %v", updated.Name, "销售合同V2")
		}
	})

	t.Run("get contract type by ID", func(t *testing.T) {
		ct, _ := service.CreateContractType(ContractTypeCreateInput{
			Name: "服务合同",
			Code: "SV001",
		})
		found, err := service.GetContractTypeByID(ct.ID)
		if err != nil {
			t.Errorf("GetContractTypeByID() error = %v", err)
			return
		}
		if found.Name != "服务合同" {
			t.Errorf("GetContractTypeByID() name = %v, want %v", found.Name, "服务合同")
		}
	})

	t.Run("delete contract type", func(t *testing.T) {
		ct, _ := service.CreateContractType(ContractTypeCreateInput{
			Name: "删除合同类型",
			Code: "DL001",
		})
		err := service.DeleteContractType(ct.ID)
		if err != nil {
			t.Errorf("DeleteContractType() error = %v", err)
			return
		}
		_, err = service.GetContractTypeByID(ct.ID)
		if err == nil {
			t.Error("DeleteContractType() type still exists after deletion")
		}
	})
}
