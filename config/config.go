package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	AppName    string `mapstructure:"APP_NAME"`
	AppVersion string `mapstructure:"APP_VERSION"`

	MysqlHost     string `mapstructure:"MYSQL_HOST"`
	MysqlPort     int    `mapstructure:"MYSQL_PORT"`
	MysqlUser     string `mapstructure:"MYSQL_USER"`
	MysqlPassword string `mapstructure:"MYSQL_PASSWORD"`
	MysqlDatabase string `mapstructure:"MYSQL_DATABASE"`

	SecretKey                string `mapstructure:"SECRET_KEY"`
	JwtAlgorithm             string `mapstructure:"JWT_ALGORITHM"`
	AccessTokenExpireMinutes int    `mapstructure:"ACCESS_TOKEN_EXPIRE_MINUTES"`

	UploadDir string `mapstructure:"UPLOAD_DIR"`

	AdminUsername string `mapstructure:"ADMIN_USERNAME"`
	AdminPassword string `mapstructure:"ADMIN_PASSWORD"`
	AdminEmail    string `mapstructure:"ADMIN_EMAIL"`

	AuditAdminUsername string `mapstructure:"AUDIT_ADMIN_USERNAME"`
	AuditAdminPassword string `mapstructure:"AUDIT_ADMIN_PASSWORD"`
	AuditAdminEmail    string `mapstructure:"AUDIT_ADMIN_EMAIL"`
}

var AppConfig Config

func LoadConfig() error {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	viper.SetDefault("APP_NAME", "安信合同管理系统")
	viper.SetDefault("APP_VERSION", "1.0.0")
	viper.SetDefault("MYSQL_HOST", "localhost")
	viper.SetDefault("MYSQL_PORT", 3306)
	viper.SetDefault("MYSQL_USER", "root")
	viper.SetDefault("MYSQL_PASSWORD", "rootroots")
	viper.SetDefault("MYSQL_DATABASE", "contract_manage")
	viper.SetDefault("JWT_ALGORITHM", "HS256")
	viper.SetDefault("ACCESS_TOKEN_EXPIRE_MINUTES", 30)
	viper.SetDefault("UPLOAD_DIR", "uploads")

	viper.SetDefault("ADMIN_USERNAME", "admin")
	viper.SetDefault("ADMIN_PASSWORD", "Admin@123456")
	viper.SetDefault("ADMIN_EMAIL", "admin@example.com")

	viper.SetDefault("AUDIT_ADMIN_USERNAME", "auditadmin")
	viper.SetDefault("AUDIT_ADMIN_PASSWORD", "Auditadmin@123456")
	viper.SetDefault("AUDIT_ADMIN_EMAIL", "audit@example.com")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}

	if err := viper.Unmarshal(&AppConfig); err != nil {
		return err
	}

	if err := validateConfig(); err != nil {
		return err
	}

	return nil
}

func validateConfig() error {
	if AppConfig.SecretKey == "" {
		return fmt.Errorf("SECRET_KEY is required")
	}

	if len(AppConfig.SecretKey) < 32 {
		return fmt.Errorf("SECRET_KEY must be at least 32 characters")
	}

	if AppConfig.SecretKey == "your-secret-key-change-in-production" {
		// 在生产环境中禁止使用默认密钥
		if os.Getenv("GO_ENV") == "production" || os.Getenv("GIN_MODE") == "release" {
			return fmt.Errorf("SECRET_KEY cannot be the default value in production environment")
		}
		fmt.Println("WARNING: Using default SECRET_KEY. Please change it in production!")
	}

	if AppConfig.MysqlPassword == "" {
		return fmt.Errorf("MYSQL_PASSWORD is required")
	}

	if AppConfig.MysqlHost == "" {
		AppConfig.MysqlHost = "localhost"
	}

	if AppConfig.MysqlPort < 1 || AppConfig.MysqlPort > 65535 {
		return fmt.Errorf("MYSQL_PORT must be between 1 and 65535")
	}

	if AppConfig.AccessTokenExpireMinutes < 5 {
		AppConfig.AccessTokenExpireMinutes = 5
	}
	if AppConfig.AccessTokenExpireMinutes > 1440 {
		AppConfig.AccessTokenExpireMinutes = 1440
	}

	AppConfig.MysqlDatabase = strings.TrimSpace(AppConfig.MysqlDatabase)
	AppConfig.MysqlUser = strings.TrimSpace(AppConfig.MysqlUser)
	AppConfig.AdminUsername = strings.TrimSpace(AppConfig.AdminUsername)

	if AppConfig.AdminPassword == "admin123" {
		// 在生产环境中禁止使用默认密码
		if os.Getenv("GO_ENV") == "production" || os.Getenv("GIN_MODE") == "release" {
			return fmt.Errorf("ADMIN_PASSWORD cannot be the default value in production environment")
		}
		fmt.Println("WARNING: Using default admin password. Please change it in production!")
	}

	if AppConfig.AuditAdminPassword == "audit123" {
		// 在生产环境中禁止使用默认密码
		if os.Getenv("GO_ENV") == "production" || os.Getenv("GIN_MODE") == "release" {
			return fmt.Errorf("AUDIT_ADMIN_PASSWORD cannot be the default value in production environment")
		}
		fmt.Println("WARNING: Using default audit admin password. Please change it in production!")
	}

	return nil
}
