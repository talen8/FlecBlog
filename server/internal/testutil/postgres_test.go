package testutil

import (
	"os"
	"testing"
)

func TestIsClearlyTestDatabaseName(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{name: "flecblog_article_lock", want: true},
		{name: "flecblog_mcp_test", want: true},
		{name: "test_flecblog", want: true},
		{name: "flecblog_testing", want: true},
		{name: "postgres", want: false},
		{name: "template0", want: false},
		{name: "template1", want: false},
		{name: "flecblog", want: false},
		{name: "flecblog_prod", want: false},
		{name: "prod_article_lock", want: false},
		{name: "contest", want: false},
		{name: "production", want: false},
		{name: "", want: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := isClearlyTestDatabaseName(tc.name); got != tc.want {
				t.Fatalf("isClearlyTestDatabaseName(%q) = %v, want %v", tc.name, got, tc.want)
			}
		})
	}
}

func TestOpenMCPTestPostgresOverPostgres(t *testing.T) {
	db := OpenMCPTestPostgres(t)
	var databaseName string
	if err := db.Raw("SELECT current_database()").Scan(&databaseName).Error; err != nil {
		t.Fatalf("read current database: %v", err)
	}
	want := os.Getenv(MCPTestPostgresDatabaseEnv)
	if databaseName != want {
		t.Fatalf("current database = %q, want %q", databaseName, want)
	}
}
