package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	App       AppConfig
	Database  DatabaseConfig
	JWT       JWTConfig
	Bootstrap BootstrapConfig
}

type AppConfig struct {
	Port string
	Env  string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type BootstrapConfig struct {
	Migrate               bool
	ResetDatabase         bool
	SeedMasterData        bool
	AdminPassword         string
	DefaultEntityCode     string
	DefaultEntityName     string
	DefaultDepartmentCode string
	DefaultDepartmentName string
}

type JWTConfig struct {
	Secret          string
	ExpiryHours     int
	RefreshExpHours int
}

func Load() *Config {
	jwtExpiry, _ := strconv.Atoi(getEnv("JWT_EXPIRY_HOURS", "24"))
	jwtRefreshExp, _ := strconv.Atoi(getEnv("JWT_REFRESH_EXPIRY_HOURS", "168"))

	return &Config{
		App: AppConfig{
			Port: getEnv("APP_PORT", "8080"),
			Env:  getEnv("APP_ENV", "development"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "3306"),
			User:     getEnv("DB_USER", "root"),
			Password: getEnv("DB_PASSWORD", ""),
			Name:     getEnv("DB_NAME", "e_procurement"),
		},
		Bootstrap: BootstrapConfig{
			Migrate:               getEnvBool("DB_MIGRATE", false),
			ResetDatabase:         getEnvBool("DB_RESET", false),
			SeedMasterData:        getEnvBool("DB_SEED", false),
			AdminPassword:         getEnv("SEED_ADMIN_PASSWORD", "Admin123!"),
			DefaultEntityCode:     getEnv("SEED_ENTITY_CODE", "HO"),
			DefaultEntityName:     getEnv("SEED_ENTITY_NAME", "Head Office"),
			DefaultDepartmentCode: getEnv("SEED_DEPARTMENT_CODE", "PROC"),
			DefaultDepartmentName: getEnv("SEED_DEPARTMENT_NAME", "Procurement"),
		},
		JWT: JWTConfig{
			Secret:          getEnv("JWT_SECRET", "change-me-in-production"),
			ExpiryHours:     jwtExpiry,
			RefreshExpHours: jwtRefreshExp,
		},
	}
}

func (d *DatabaseConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		d.User, d.Password, d.Host, d.Port, d.Name,
	)
}

func (d *DatabaseConfig) ServerDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/?charset=utf8mb4&parseTime=True&loc=Local",
		d.User, d.Password, d.Host, d.Port,
	)
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}

func getEnvBool(key string, fallback bool) bool {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	parsed, err := strconv.ParseBool(val)
	if err != nil {
		return fallback
	}

	return parsed
}
