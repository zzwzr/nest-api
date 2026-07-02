package configs

type RuntimeConfig struct {
	Installed   bool            `yaml:"installed"`
	InstalledAt string          `yaml:"installed_at,omitempty"`
	Database    DatabaseRuntime `yaml:"database"`
}

type DatabaseRuntime struct {
	Driver   string `yaml:"driver"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
	SSLMode  string `yaml:"ssl_mode"`
}
