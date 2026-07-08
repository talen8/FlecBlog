package testutil

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"flec_blog/pkg/database"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	MCPTestPostgresDSNEnv      = "FLECBLOG_MCP_TEST_POSTGRES_DSN"
	MCPTestPostgresDatabaseEnv = "FLECBLOG_MCP_TEST_POSTGRES_DATABASE"

	mcpTestMigrationLockKey int64 = 0x464c45434d435054
)

type postgresIdentity struct {
	DatabaseName string `gorm:"column:database_name"`
	UserName     string `gorm:"column:user_name"`
	Host         string `gorm:"column:host"`
	Port         int    `gorm:"column:port"`
}

// OpenMCPTestPostgres opens the explicitly confirmed MCP test database,
// serializes schema migrations across packages, and fails closed on ambiguity.
func OpenMCPTestPostgres(t *testing.T) *gorm.DB {
	t.Helper()

	dsn := strings.TrimSpace(os.Getenv(MCPTestPostgresDSNEnv))
	if dsn == "" {
		t.Skip(MCPTestPostgresDSNEnv + " is not set")
	}
	expectedDatabase := strings.TrimSpace(os.Getenv(MCPTestPostgresDatabaseEnv))
	if expectedDatabase == "" {
		t.Fatalf("%s must explicitly name the dedicated test database", MCPTestPostgresDatabaseEnv)
	}
	if !isClearlyTestDatabaseName(expectedDatabase) {
		t.Fatalf("refusing unsafe PostgreSQL test database name %q", expectedDatabase)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("open postgres test database: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("get postgres test connection pool: %v", err)
	}
	t.Cleanup(func() {
		if err := sqlDB.Close(); err != nil {
			t.Errorf("close postgres test connection pool: %v", err)
		}
	})
	var identity postgresIdentity
	if err := db.Raw(`
		SELECT
			current_database() AS database_name,
			current_user AS user_name,
			COALESCE(inet_server_addr()::text, 'local') AS host,
			COALESCE(inet_server_port(), 0) AS port
	`).Scan(&identity).Error; err != nil {
		t.Fatalf("read postgres test identity: %v", err)
	}
	if identity.DatabaseName != expectedDatabase {
		t.Fatalf("postgres test database mismatch: connected=%q expected=%q", identity.DatabaseName, expectedDatabase)
	}
	if !isClearlyTestDatabaseName(identity.DatabaseName) {
		t.Fatalf("refusing connected PostgreSQL database %q", identity.DatabaseName)
	}
	t.Logf("postgres test identity database=%s user=%s host=%s port=%d", identity.DatabaseName, identity.UserName, identity.Host, identity.Port)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	conn, err := sqlDB.Conn(ctx)
	if err != nil {
		t.Fatalf("reserve postgres migration lock connection: %v", err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			t.Errorf("close postgres migration lock connection: %v", err)
		}
	}()

	if _, err := conn.ExecContext(ctx, "SELECT pg_advisory_lock($1)", mcpTestMigrationLockKey); err != nil {
		t.Fatalf("acquire postgres migration advisory lock: %v", err)
	}
	defer func() {
		unlockCtx, unlockCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer unlockCancel()
		var unlocked bool
		if err := conn.QueryRowContext(unlockCtx, "SELECT pg_advisory_unlock($1)", mcpTestMigrationLockKey).Scan(&unlocked); err != nil {
			t.Errorf("release postgres migration advisory lock: %v", err)
			return
		}
		if !unlocked {
			t.Errorf("release postgres migration advisory lock: lock was not held")
		}
	}()

	if err := database.RunMigrations(db); err != nil {
		t.Fatalf("run postgres test migrations: %v", err)
	}
	return db
}

func isClearlyTestDatabaseName(name string) bool {
	name = strings.ToLower(strings.TrimSpace(name))
	if name == "" {
		return false
	}
	switch name {
	case "postgres", "template0", "template1":
		return false
	}
	if name == "flecblog_article_lock" {
		return true
	}
	return name == "test" ||
		strings.HasPrefix(name, "test_") ||
		strings.HasSuffix(name, "_test") ||
		strings.Contains(name, "_test_") ||
		strings.HasPrefix(name, "testing_") ||
		strings.HasSuffix(name, "_testing") ||
		strings.Contains(name, "_testing_")
}
