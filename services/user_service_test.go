package services

import (
	"contract-manage/models"
	"os"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupTestDB(t *testing.T) func() {
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

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestUserService_CreateUser(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	service := NewUserService()

	tests := []struct {
		name    string
		input   UserCreateInput
		wantErr bool
	}{
		{
			name: "valid user creation",
			input: UserCreateInput{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password123",
				FullName: "Test User",
			},
			wantErr: false,
		},
		{
			name: "duplicate username",
			input: UserCreateInput{
				Username: "testuser",
				Email:    "test2@example.com",
				Password: "password123",
			},
			wantErr: true,
		},
		{
			name: "duplicate email",
			input: UserCreateInput{
				Username: "testuser2",
				Email:    "test@example.com",
				Password: "password123",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := service.CreateUser(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && user == nil {
				t.Error("CreateUser() returned nil user without error")
			}
		})
	}
}

func TestUserService_GetUserByID(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	service := NewUserService()

	created, err := service.CreateUser(UserCreateInput{
		Username: "gettest",
		Email:    "get@example.com",
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}

	t.Run("get existing user", func(t *testing.T) {
		user, err := service.GetUserByID(created.ID)
		if err != nil {
			t.Errorf("GetUserByID() error = %v", err)
			return
		}
		if user.Username != "gettest" {
			t.Errorf("GetUserByID() username = %v, want %v", user.Username, "gettest")
		}
	})

	t.Run("get non-existing user", func(t *testing.T) {
		_, err := service.GetUserByID(9999)
		if err == nil {
			t.Error("GetUserByID() expected error for non-existing user")
		}
	})
}

func TestUserService_GetUserByUsername(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	service := NewUserService()

	_, err := service.CreateUser(UserCreateInput{
		Username: "searchuser",
		Email:    "search@example.com",
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}

	t.Run("find existing user", func(t *testing.T) {
		user, err := service.GetUserByUsername("searchuser")
		if err != nil {
			t.Errorf("GetUserByUsername() error = %v", err)
			return
		}
		if user.Email != "search@example.com" {
			t.Errorf("GetUserByUsername() email = %v, want %v", user.Email, "search@example.com")
		}
	})

	t.Run("find non-existing user", func(t *testing.T) {
		_, err := service.GetUserByUsername("nonexistent")
		if err == nil {
			t.Error("GetUserByUsername() expected error for non-existing user")
		}
	})
}

func TestUserService_AuthenticateUser(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	service := NewUserService()

	_, err := service.CreateUser(UserCreateInput{
		Username: "authtest",
		Email:    "auth@example.com",
		Password: "correctpassword",
	})
	if err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}

	t.Run("authenticate with correct password", func(t *testing.T) {
		user, err := service.AuthenticateUser("authtest", "correctpassword")
		if err != nil {
			t.Errorf("AuthenticateUser() error = %v", err)
			return
		}
		if user.Username != "authtest" {
			t.Errorf("AuthenticateUser() username = %v, want %v", user.Username, "authtest")
		}
	})

	t.Run("authenticate with wrong password", func(t *testing.T) {
		_, err := service.AuthenticateUser("authtest", "wrongpassword")
		if err == nil {
			t.Error("AuthenticateUser() expected error for wrong password")
		}
	})

	t.Run("authenticate non-existing user", func(t *testing.T) {
		_, err := service.AuthenticateUser("nonexistent", "password")
		if err == nil {
			t.Error("AuthenticateUser() expected error for non-existing user")
		}
	})
}

func TestUserService_UpdateUser(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	service := NewUserService()

	created, _ := service.CreateUser(UserCreateInput{
		Username: "updatetest",
		Email:    "update@example.com",
		Password: "password123",
		FullName: "Original Name",
	})

	t.Run("update user fields", func(t *testing.T) {
		updated, err := service.UpdateUser(created.ID, UserUpdateInput{
			FullName: "Updated Name",
			Email:    "updated@example.com",
		})
		if err != nil {
			t.Errorf("UpdateUser() error = %v", err)
			return
		}
		if updated.FullName != "Updated Name" {
			t.Errorf("UpdateUser() full_name = %v, want %v", updated.FullName, "Updated Name")
		}
		if updated.Email != "updated@example.com" {
			t.Errorf("UpdateUser() email = %v, want %v", updated.Email, "updated@example.com")
		}
	})

	t.Run("update non-existing user", func(t *testing.T) {
		_, err := service.UpdateUser(9999, UserUpdateInput{FullName: "Name"})
		if err == nil {
			t.Error("UpdateUser() expected error for non-existing user")
		}
	})
}

func TestUserService_DeleteUser(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	service := NewUserService()

	created, _ := service.CreateUser(UserCreateInput{
		Username: "deletetest",
		Email:    "delete@example.com",
		Password: "password123",
	})

	t.Run("delete existing user", func(t *testing.T) {
		err := service.DeleteUser(created.ID)
		if err != nil {
			t.Errorf("DeleteUser() error = %v", err)
			return
		}
		_, err = service.GetUserByID(created.ID)
		if err == nil {
			t.Error("DeleteUser() user still exists after deletion")
		}
	})

	t.Run("delete non-existing user", func(t *testing.T) {
		err := service.DeleteUser(9999)
		if err == nil {
			t.Error("DeleteUser() expected error for non-existing user")
		}
	})
}

func TestUserService_GetUsers(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	service := NewUserService()

	for i := 0; i < 5; i++ {
		service.CreateUser(UserCreateInput{
			Username: "listuser",
			Email:    "",
			Password: "password",
		})
	}

	t.Run("get users with pagination", func(t *testing.T) {
		users, err := service.GetUsers(0, 3)
		if err != nil {
			t.Errorf("GetUsers() error = %v", err)
			return
		}
		if len(users) != 3 {
			t.Errorf("GetUsers() returned %v users, want %v", len(users), 3)
		}
	})

	t.Run("get users with offset", func(t *testing.T) {
		users, err := service.GetUsers(3, 10)
		if err != nil {
			t.Errorf("GetUsers() error = %v", err)
			return
		}
		if len(users) != 2 {
			t.Errorf("GetUsers() returned %v users, want %v", len(users), 2)
		}
	})
}
