package repository

import (
	"os"
	"gopkg.in/yaml.v3"
	"fmt"
	"github.com/joho/godotenv"
	"payment/internal/repository/mysql"
)

type DatabaseConfig struct {
	DBPayment *mysql.PaymentDB `yaml:"mysql"`
}

func LoadConfig() (*DatabaseConfig, error) {
	err := godotenv.Load()
	configPath := os.Getenv("DATABASE_CONFIG")
	if configPath == "" {
		configPath = "configs/database.yaml"
	}
	data, err := os.ReadFile(configPath)
    if err != nil {
        return nil, fmt.Errorf("读取文件失败: %w", err)
    }

    var config DatabaseConfig
    if err := yaml.Unmarshal(data, &config); err != nil {
        return nil, fmt.Errorf("解析YAML失败: %w", err)
    }
    return &config, nil
}
