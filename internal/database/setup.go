package database

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"nest-api/configs"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func BuildDSN(cfg configs.DatabaseConfig) string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=Asia/Shanghai",
		cfg.Host,
		cfg.User,
		cfg.Password,
		cfg.DBName,
		cfg.Port,
		cfg.SSLMode,
	)
}

func TestConnection(cfg configs.DatabaseConfig) error {
	db, err := openSQL(cfg)
	if err != nil {
		return err
	}
	defer db.Close()

	return db.PingContext(context.Background())
}

func EnsureDatabase(cfg configs.DatabaseConfig) error {
	if cfg.DBName == "" {
		return fmt.Errorf("database name is required")
	}

	adminCfg := cfg
	adminCfg.DBName = "postgres"

	db, err := openSQL(adminCfg)
	if err != nil {
		return err
	}
	defer db.Close()

	if err := db.PingContext(context.Background()); err != nil {
		return fmt.Errorf("connect postgres failed: %w", err)
	}

	var exists bool
	err = db.QueryRowContext(
		context.Background(),
		`SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)`,
		cfg.DBName,
	).Scan(&exists)
	if err != nil {
		return err
	}

	if exists {
		return nil
	}

	_, err = db.ExecContext(
		context.Background(),
		fmt.Sprintf(`CREATE DATABASE "%s"`, escapeIdentifier(cfg.DBName)),
	)
	if err != nil {
		return fmt.Errorf("create database failed: %w", err)
	}

	return nil
}

func openSQL(cfg configs.DatabaseConfig) (*sql.DB, error) {
	db, err := sql.Open("pgx", BuildDSN(cfg))
	if err != nil {
		return nil, err
	}

	return db, nil
}

func escapeIdentifier(name string) string {
	return strings.ReplaceAll(name, `"`, `""`)
}

func EnsureAppUser(adminCfg configs.DatabaseConfig, appUser, appPassword string) error {
	if appUser == "" {
		return fmt.Errorf("application database user is required")
	}
	if appPassword == "" {
		return fmt.Errorf("application database password is required")
	}
	if adminCfg.DBName == "" {
		return fmt.Errorf("target database name is required")
	}

	targetDBName := adminCfg.DBName
	adminCfg.DBName = "postgres"

	db, err := openSQL(adminCfg)
	if err != nil {
		return err
	}
	defer db.Close()

	ctx := context.Background()

	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("connect postgres failed: %w", err)
	}

	var exists bool
	err = db.QueryRowContext(
		ctx,
		`SELECT EXISTS(SELECT 1 FROM pg_roles WHERE rolname = $1)`,
		appUser,
	).Scan(&exists)
	if err != nil {
		return err
	}

	if exists {
		_, err = db.ExecContext(
			ctx,
			fmt.Sprintf(`ALTER ROLE "%s" WITH LOGIN PASSWORD $1`, escapeIdentifier(appUser)),
			appPassword,
		)
	} else {
		_, err = db.ExecContext(
			ctx,
			fmt.Sprintf(`CREATE ROLE "%s" LOGIN PASSWORD $1`, escapeIdentifier(appUser)),
			appPassword,
		)
	}
	if err != nil {
		return fmt.Errorf("create application user failed: %w", err)
	}

	dbName := targetDBName

	_, err = db.ExecContext(
		ctx,
		fmt.Sprintf(`ALTER DATABASE "%s" OWNER TO "%s"`, escapeIdentifier(dbName), escapeIdentifier(appUser)),
	)
	if err != nil {
		return fmt.Errorf("grant database ownership failed: %w", err)
	}

	targetCfg := adminCfg
	targetCfg.DBName = dbName
	targetCfg.User = appUser
	targetCfg.Password = appPassword

	targetDB, err := openSQL(targetCfg)
	if err != nil {
		return err
	}
	defer targetDB.Close()

	if err := targetDB.PingContext(ctx); err != nil {
		return fmt.Errorf("connect target database failed: %w", err)
	}

	_, err = targetDB.ExecContext(ctx, fmt.Sprintf(`GRANT ALL ON SCHEMA public TO "%s"`, escapeIdentifier(appUser)))
	if err != nil {
		return fmt.Errorf("grant schema privileges failed: %w", err)
	}

	_, err = targetDB.ExecContext(ctx, fmt.Sprintf(`ALTER SCHEMA public OWNER TO "%s"`, escapeIdentifier(appUser)))
	if err != nil {
		return fmt.Errorf("transfer schema ownership failed: %w", err)
	}

	return nil
}
