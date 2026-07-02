package install

import (
	"context"
	"fmt"
	"regexp"
	"sync"
	"time"

	"nest-api/configs"
	"nest-api/internal/database"
	"nest-api/internal/runtime"
	"nest-api/internal/utils"
	bizerr "nest-api/pkg/errors"
)

var pgRoleNamePattern = regexp.MustCompile(`^[a-z][a-z0-9_]*$`)

type Service struct{}

var installMu sync.Mutex

func (Service) Status() StatusResponse {
	return StatusResponse{
		Installed: runtime.IsInstalled(),
	}
}

func (s Service) TestDatabase(params TestDatabaseRequest) (TestDatabaseResponse, error) {
	if params.Database.Driver == "mysql" {
		return TestDatabaseResponse{}, bizerr.New("MySQL 将在后续版本中支持，请暂时使用 PostgreSQL")
	}

	cfg := toDatabaseConfig(params.Database)

	if err := database.TestConnection(cfg); err == nil {
		return TestDatabaseResponse{
			OK:      true,
			Message: "数据库连接成功",
		}, nil
	}

	adminCfg := cfg
	adminCfg.DBName = "postgres"

	if err := database.TestConnection(adminCfg); err != nil {
		return TestDatabaseResponse{}, bizerr.New("无法连接数据库，请检查主机、端口、账号和密码")
	}

	return TestDatabaseResponse{
		OK:      true,
		Message: "数据库服务连接成功，目标库不存在时将在安装过程中自动创建",
	}, nil
}

func (s Service) Install(ctx context.Context, params InstallRequest) (InstallResponse, error) {
	installMu.Lock()
	defer installMu.Unlock()

	if runtime.IsInstalled() {
		return InstallResponse{}, bizerr.New("系统已安装，请勿重复安装")
	}

	if params.Database.Driver == "mysql" {
		return InstallResponse{}, bizerr.New("MySQL 将在后续版本中支持，请暂时使用 PostgreSQL")
	}

	if params.Admin.Password != params.Admin.ConfirmPassword {
		return InstallResponse{}, bizerr.New("两次输入的管理员密码不一致")
	}

	appUser := params.AppDatabase.Username
	appPassword := params.AppDatabase.Password

	if !pgRoleNamePattern.MatchString(appUser) {
		return InstallResponse{}, bizerr.New("应用数据库用户名需以小写字母开头，仅包含小写字母、数字和下划线")
	}

	if appUser == params.Database.User {
		return InstallResponse{}, bizerr.New("应用数据库用户名不能与超级用户名相同")
	}

	cfg := toDatabaseConfig(params.Database)

	if err := database.EnsureDatabase(cfg); err != nil {
		return InstallResponse{}, bizerr.New(fmt.Sprintf("创建数据库失败: %v", err))
	}

	if err := database.EnsureAppUser(cfg, appUser, appPassword); err != nil {
		return InstallResponse{}, bizerr.New(fmt.Sprintf("创建应用数据库用户失败: %v", err))
	}

	appCfg := cfg
	appCfg.User = appUser
	appCfg.Password = appPassword

	if err := database.InitWithConfig(appCfg); err != nil {
		return InstallResponse{}, bizerr.New(fmt.Sprintf("初始化数据库失败: %v", err))
	}

	hash, err := utils.Hash(params.Admin.Password)
	if err != nil {
		return InstallResponse{}, err
	}

	_, err = database.DB.User.
		Create().
		SetName(params.Admin.Username).
		SetPassword(hash).
		SetIsAdmin(true).
		SetStatus(1).
		Save(ctx)
	if err != nil {
		return InstallResponse{}, bizerr.New(fmt.Sprintf("创建管理员失败: %v", err))
	}

	runtimeCfg := configs.RuntimeConfig{
		Installed:   true,
		InstalledAt: time.Now().Format(time.RFC3339),
		Database: configs.DatabaseRuntime{
			Driver:   params.Database.Driver,
			Host:     params.Database.Host,
			Port:     params.Database.Port,
			User:     appUser,
			Password: appPassword,
			Name:     params.Database.Name,
			SSLMode:  params.Database.SSLMode,
		},
	}

	if err := runtime.Save(runtimeCfg); err != nil {
		return InstallResponse{}, bizerr.New(fmt.Sprintf("保存配置失败: %v", err))
	}

	return InstallResponse{
		AdminUsername:    params.Admin.Username,
		AdminPassword:    params.Admin.Password,
		DatabaseUser:     appUser,
		DatabasePassword: appPassword,
	}, nil
}

func toDatabaseConfig(db DatabaseRequest) configs.DatabaseConfig {
	return configs.DatabaseConfig{
		Host:     db.Host,
		Port:     db.Port,
		User:     db.User,
		Password: db.Password,
		DBName:   db.Name,
		SSLMode:  db.SSLMode,
	}
}
