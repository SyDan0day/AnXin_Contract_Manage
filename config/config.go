package config

import (
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
}

var AppConfig Config

func LoadConfig() error {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	viper.SetDefault("APP_NAME", "合同管理系统")
	viper.SetDefault("APP_VERSION", "1.0.0")
	viper.SetDefault("MYSQL_HOST", "localhost")
	viper.SetDefault("MYSQL_PORT", 3306)
	viper.SetDefault("MYSQL_USER", "root")
	viper.SetDefault("MYSQL_PASSWORD", "password")
	viper.SetDefault("MYSQL_DATABASE", "contract_manage")
	viper.SetDefault("SECRET_KEY", "your-secret-key-change-in-production")
	viper.SetDefault("JWT_ALGORITHM", "HS256")
	viper.SetDefault("ACCESS_TOKEN_EXPIRE_MINUTES", 30)
	viper.SetDefault("UPLOAD_DIR", "uploads")

	viper.SetDefault("ADMIN_USERNAME", "admin")
	viper.SetDefault("ADMIN_PASSWORD", "admin123")
	viper.SetDefault("ADMIN_EMAIL", "admin@example.com")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}

	return viper.Unmarshal(&AppConfig)
}
