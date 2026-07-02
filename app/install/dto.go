package install

type DatabaseRequest struct {
	Driver   string `json:"driver" binding:"required,oneof=postgres mysql"`
	Host     string `json:"host" binding:"required"`
	Port     int    `json:"port" binding:"required,min=1,max=65535"`
	Name     string `json:"name" binding:"required"`
	User     string `json:"user" binding:"required"`
	Password string `json:"password" binding:"required"`
	SSLMode  string `json:"ssl_mode" binding:"required"`
}

type AppDatabaseRequest struct {
	Username string `json:"username" binding:"required,min=1,max=63"`
	Password string `json:"password" binding:"required,min=6"`
}

type AdminRequest struct {
	Username        string `json:"username" binding:"required,min=3,max=50"`
	Password        string `json:"password" binding:"required,min=6"`
	ConfirmPassword string `json:"confirm_password" binding:"required,min=6"`
}

type InstallRequest struct {
	SiteURL     string             `json:"site_url" binding:"required,url,max=255"`
	Database    DatabaseRequest    `json:"database" binding:"required"`
	AppDatabase AppDatabaseRequest `json:"app_database" binding:"required"`
	Admin       AdminRequest       `json:"admin" binding:"required"`
}

type TestDatabaseRequest struct {
	Database DatabaseRequest `json:"database" binding:"required"`
}

type StatusResponse struct {
	Installed bool   `json:"installed"`
	SiteURL   string `json:"site_url,omitempty"`
}

type TestDatabaseResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

type InstallResponse struct {
	Message          string `json:"message,omitempty"`
	AdminUsername    string `json:"admin_username,omitempty"`
	AdminPassword    string `json:"admin_password,omitempty"`
	DatabaseUser     string `json:"database_user,omitempty"`
	DatabasePassword string `json:"database_password,omitempty"`
}
