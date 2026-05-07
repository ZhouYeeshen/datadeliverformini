package config

import (
	"os"
	"time"
)

type Config struct {
	Server       ServerConfig
	DB           DBConfig
	Redis        RedisConfig
	JWT          JWTConfig
	WeChat       WeChatConfig
	TencentCloud TencentCloudConfig
	Storage      StorageConfig
}

type TencentCloudConfig struct {
	SecretID  string
	SecretKey string
	Region    string
}

type ServerConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type DBConfig struct {
	Driver   string
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

type JWTConfig struct {
	Secret     string
	ExpireTime time.Duration
}

type WeChatConfig struct {
	AppID     string
	AppSecret string
}

type StorageConfig struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	Bucket    string
}

func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Port:         getEnv("SERVER_PORT", "8080"),
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 30 * time.Second,
		},
		DB: DBConfig{
			Driver:   getEnv("DB_DRIVER", "sqlite"),
			Host:     getEnv("DB_HOST", "127.0.0.1"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", ""),
			DBName:   getEnv("DB_NAME", "business_report"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Redis: RedisConfig{
			Addr:     getEnv("REDIS_ADDR", "127.0.0.1:6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       0,
		},
		JWT: JWTConfig{
			Secret:     getEnv("JWT_SECRET", "change-me-in-production"),
			ExpireTime: 24 * time.Hour,
		},
		WeChat: WeChatConfig{
			AppID:     getEnv("WECHAT_APPID", ""),
			AppSecret: getEnv("WECHAT_APPSECRET", ""),
		},
		TencentCloud: TencentCloudConfig{
			SecretID:  getEnv("TENCENT_SECRET_ID", ""),
			SecretKey: getEnv("TENCENT_SECRET_KEY", ""),
			Region:    getEnv("TENCENT_REGION", "ap-guangzhou"),
		},
		Storage: StorageConfig{
			Endpoint:  getEnv("STORAGE_ENDPOINT", ""),
			AccessKey: getEnv("STORAGE_ACCESSKEY", ""),
			SecretKey: getEnv("STORAGE_SECRETKEY", ""),
			Bucket:    getEnv("STORAGE_BUCKET", "business-reports"),
		},
	}
}

func (c DBConfig) DSN() string {
	return "host=" + c.Host + " port=" + c.Port + " user=" + c.User +
		" password=" + c.Password + " dbname=" + c.DBName + " sslmode=" + c.SSLMode
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
