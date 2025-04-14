package main 

import (
	"gopkg.in/yaml.v3"
	"os"
	"fmt"
)

type cfg2 struct {
	Dsn string `yaml:"dsn"`
}

type cfg3 struct {
	ServiceName string `yaml:"service_name"`
}

type Config struct {
	Registry string `yaml:"registry"`
	DBConfig map[string]cfg2 `yaml:"databases"`
	Service map[string]cfg3 `yaml:"services"`
}

func main() {
	data, err := os.ReadFile("configs/database.yaml")
    if err!= nil {
        fmt.Println("读取文件失败:", err)
    }
    var config Config
    if err := yaml.Unmarshal(data, &config); err!= nil {
        fmt.Println("解析YAML失败:", err)
    }
    fmt.Println(config.DBConfig["payment"].Dsn)

	data2, err := os.ReadFile("configs/service.yaml")
    if err!= nil {
        fmt.Println("读取文件失败:", err)
    }
    if err := yaml.Unmarshal(data2, &config); err!= nil {
        fmt.Println("解析YAML失败:", err)
    }
    fmt.Println("---",config.Service["http"].ServiceName)

	fmt.Println(config)

    
}