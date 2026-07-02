package runtime

import (
	"os"
	"path/filepath"
	"sync"

	"nest-api/configs"

	"gopkg.in/yaml.v3"
)

const configPath = "runtime/config.yaml"

var (
	state configs.RuntimeConfig
	mu    sync.RWMutex
)

func Init() {
	mu.Lock()
	defer mu.Unlock()

	data, err := os.ReadFile(configPath)
	if err != nil {
		return
	}

	var cfg configs.RuntimeConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return
	}

	state = cfg

	if cfg.Installed {
		applyDatabase(cfg.Database)
	}
}

func IsInstalled() bool {
	mu.RLock()
	defer mu.RUnlock()

	return state.Installed
}

func Save(cfg configs.RuntimeConfig) error {
	mu.Lock()
	defer mu.Unlock()

	if err := os.MkdirAll(filepath.Dir(configPath), 0o755); err != nil {
		return err
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	if err := os.WriteFile(configPath, data, 0o600); err != nil {
		return err
	}

	state = cfg

	if cfg.Installed {
		applyDatabase(cfg.Database)
	}

	return nil
}

func applyDatabase(db configs.DatabaseRuntime) {
	configs.Database.Host = db.Host
	configs.Database.Port = db.Port
	configs.Database.User = db.User
	configs.Database.Password = db.Password
	configs.Database.DBName = db.Name
	configs.Database.SSLMode = db.SSLMode
}
