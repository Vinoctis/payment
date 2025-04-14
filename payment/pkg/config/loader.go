package config

import (
	"fmt"
	"os"
	"gopkg.in/yaml.v3"
	"github.com/joho/godotenv"
)

type ConfigLoader interface {
	Load(cfg interface{}) error
}

type EnvLoader struct {
	Envfile string 
}

func (l *EnvLoader) Load(cfg interface{}) error {
	if err := godotenv.Load(l.Envfile); err != nil {
		return fmt.Errorf(".env环境配置加载失败：%w", err)
	} 
	return nil
}

type YmalLoader struct {
	Path string
}

func (l *YmalLoader) Load(cfg interface{}) error{
	data, err := os.ReadFile(l.Path)
	if err != nil {
		return fmt.Errorf("配置文件读取失败：%w", err)
	}

	if err := yaml.Unmarshal(data, cfg); err!= nil {
		return fmt.Errorf("配置文件解析失败：%w", err)
	}
	return nil
}



