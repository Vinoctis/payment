package config

import (

)

type Component interface {
	Name() string 
	Init() error
}

type Config struct {
	App map[string]*AppConfig `yaml:"app"`
	DB map[string]*DBConfig `yaml:"databases"`
	Service map[string]*ServiceConfig `yaml:"services"`
	Registry map[string]*RegistryConfig `yaml:"etcd-nodes"`
	components map[string]Component
}

type AppConfig struct {
	Name string `yaml:"name"`
	Env string `yaml:"env"`
	LogLevel string `yaml:"log_level"`
}

type DBConfig struct {
	Driver string `yaml:"driver"`
	DSN string `yaml:"dsn"`
	Port string  `yaml:"port"`
	Host string `yaml:"host"`
	User string `yaml:"user"`
	Password string `yaml:"password"`
	DBName string	`yaml:"dbname"`
	MaxOpenConns int `yaml:"max_open_conns"`
	MaxIdleConns int `yaml:"max_idle_conns"`
	MaxLifeTime int `yaml:"max_life_time"`
}

type ServiceConfig struct {
	Name string `yaml:"name"`
	Port string `yaml:"port"`
	Addr string `yaml:"addr"`
}

type RegistryConfig struct {
	Addr string `yaml:"addr"`
	Port string `yaml:"port"`
}

func (c *Config) RegisterComponent(comp Component) error {
	if c.components == nil {
		c.components = make(map[string]Component)
	}
	c.components[comp.Name()] = comp
	return nil 
}

