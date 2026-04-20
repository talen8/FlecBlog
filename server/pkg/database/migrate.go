package database

import (
	"database/sql"
	"embed"
	"fmt"
	"sort"
	"strings"

	"gorm.io/gorm"

	"flec_blog/pkg/logger"
)

//go:embed sql/*.sql
var sqlFiles embed.FS

// RunMigrations 执行增量迁移
func RunMigrations(gormDB *gorm.DB) error {
	db, err := gormDB.DB()
	if err != nil {
		return fmt.Errorf("获取数据库连接失败: %w", err)
	}

	if err := createMigrationTable(db); err != nil {
		return err
	}

	if err := handleLegacyDatabase(db); err != nil {
		return err
	}

	migrations, err := loadMigrations()
	if err != nil {
		return err
	}

	executed, err := getExecutedMigrations(db)
	if err != nil {
		return err
	}

	for _, migration := range migrations {
		if executed[migration] {
			continue
		}

		logger.Info("执行迁移: %s", migration)

		if err := executeMigration(db, migration); err != nil {
			return fmt.Errorf("迁移 %s 执行失败: %w", migration, err)
		}

		if err := recordMigration(db, migration); err != nil {
			return fmt.Errorf("记录迁移 %s 失败: %w", migration, err)
		}

		logger.Info("迁移 %s 执行成功", migration)
	}

	logger.Info("所有迁移执行完成")
	return nil
}

// createMigrationTable 创建迁移记录表
func createMigrationTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version VARCHAR(50) PRIMARY KEY,
			applied_at TIMESTAMP DEFAULT NOW()
		)
	`)
	return err
}

// loadMigrations 加载所有迁移文件
func loadMigrations() ([]string, error) {
	entries, err := sqlFiles.ReadDir("sql")
	if err != nil {
		return nil, fmt.Errorf("读取迁移目录失败: %w", err)
	}

	var migrations []string
	for _, entry := range entries {
		name := entry.Name()
		if !entry.IsDir() && strings.HasSuffix(name, ".sql") {
			migrations = append(migrations, name)
		}
	}

	sort.Strings(migrations)
	return migrations, nil
}

// getExecutedMigrations 获取已执行的迁移
func getExecutedMigrations(db *sql.DB) (map[string]bool, error) {
	rows, err := db.Query("SELECT version FROM schema_migrations")
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()

	executed := make(map[string]bool)
	for rows.Next() {
		var version string
		if err := rows.Scan(&version); err != nil {
			return nil, err
		}
		executed[version] = true
	}

	return executed, rows.Err()
}

// handleLegacyDatabase 处理老用户数据库
func handleLegacyDatabase(db *sql.DB) error {
	// 检查迁移记录表是否为空
	var count int
	if err := db.QueryRow("SELECT COUNT(*) FROM schema_migrations").Scan(&count); err != nil {
		return err
	}

	if count > 0 {
		return nil
	}

	var exists bool
	query := `
		SELECT EXISTS (
			SELECT FROM information_schema.tables
			WHERE table_schema = 'public'
			AND table_name = 'users'
		)
	`
	if err := db.QueryRow(query).Scan(&exists); err != nil {
		return err
	}

	if exists {
		logger.Info("检测到老版本数据库，自动标记初始迁移为已执行")
		if err := recordMigration(db, "001_init_database.sql"); err != nil {
			return fmt.Errorf("标记初始迁移失败: %w", err)
		}
	}

	return nil
}

// executeMigration 执行单个迁移
func executeMigration(db *sql.DB, filename string) error {
	content, err := sqlFiles.ReadFile("sql/" + filename)
	if err != nil {
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	if _, err := tx.Exec(string(content)); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}

// recordMigration 记录已执行的迁移
func recordMigration(db *sql.DB, filename string) error {
	_, err := db.Exec("INSERT INTO schema_migrations (version) VALUES ($1)", filename)
	return err
}
