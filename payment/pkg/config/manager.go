package config 

import (
	"sync"
)

type Manager struct {
	Configs *Config
	loaders []ConfigLoader 
	loadOnce sync.Once
}

func NewManager() *Manager {
	return &Manager{
		Configs: &Config{},
		loaders: make([]ConfigLoader, 0),
	}
}

func (m *Manager) RegisterLoader(loader ConfigLoader) {
	m.loaders = append(m.loaders, loader)
}

func (m *Manager) Load() error {
	var err error
	m.loadOnce.Do(func(){
		for _, loader := range m.loaders {
			if err = loader.Load(m.Configs);err != nil {
				return
			}
		}
	})
	return err
}
