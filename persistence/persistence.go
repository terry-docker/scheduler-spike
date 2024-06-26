package persistence

import (
	"encoding/json"
	"os"
	"sync"

	"example.com/scheduler/common"
)

type IPersistenceManager interface {
	SaveTasks([]common.TaskConfig) error
	LoadTasks() ([]common.TaskConfig, error)
}

type Manager struct {
	mu         sync.Mutex
	configFile string
}

var _ IPersistenceManager = (*Manager)(nil)

func NewPersistenceManager(configFile string) *Manager {
	return &Manager{configFile: configFile}
}

func (pm *Manager) SaveTasks(tasks []common.TaskConfig) error {
	file, err := os.Create(pm.configFile)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	return encoder.Encode(tasks)
}

func (pm *Manager) LoadTasks() ([]common.TaskConfig, error) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	// Check if the config file exists
	if _, err := os.Stat(pm.configFile); os.IsNotExist(err) {
		// Create an empty JSON file if it does not exist, using SaveTasks to avoid re-locking.
		return []common.TaskConfig{}, pm.SaveTasks([]common.TaskConfig{})
	}

	file, err := os.Open(pm.configFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var tasks []common.TaskConfig
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&tasks); err != nil {
		return nil, err
	}

	return tasks, nil
}
