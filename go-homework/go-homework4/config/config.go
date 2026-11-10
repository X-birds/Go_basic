package config

import "os"

type Config struct {
	Server ServerConfig
	JWT    JWTConfig
	DB     DBConfig
}

type ServerConfig struct {
	Port string // 服务监听端口
}

type JWTConfig struct {
	Secret string // 签名密钥，绝对不能硬编码提交到仓库
}

type DBConfig struct {
	Path string // 数据库连接地址
}

func LoadConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
		},
		JWT: JWTConfig{
			Secret: getEnv("JWT_SECRET", "ev-secret-change-me"),
		},
		DB: DBConfig{
			Path: getEnv("DB_PATH", "blog.db"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}
