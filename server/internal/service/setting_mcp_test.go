package service

import (
	"fmt"
	"sync"
	"testing"

	"flec_blog/config"
)

func TestParseMCPStaticOperatorUserID(t *testing.T) {
	tests := []struct {
		name    string
		raw     string
		want    uint
		wantErr bool
	}{
		{name: "empty means auto", raw: "", want: 0},
		{name: "whitespace means auto", raw: "   ", want: 0},
		{name: "zero means auto", raw: "0", want: 0},
		{name: "positive id", raw: "42", want: 42},
		{name: "trim positive id", raw: " 7 ", want: 7},
		{name: "negative rejected", raw: "-1", wantErr: true},
		{name: "decimal rejected", raw: "1.5", wantErr: true},
		{name: "text rejected", raw: "admin", wantErr: true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := parseMCPStaticOperatorUserID(tc.raw)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("parseMCPStaticOperatorUserID(%q) error = nil", tc.raw)
				}
				return
			}
			if err != nil {
				t.Fatalf("parseMCPStaticOperatorUserID(%q): %v", tc.raw, err)
			}
			if got != tc.want {
				t.Fatalf("parseMCPStaticOperatorUserID(%q) = %d, want %d", tc.raw, got, tc.want)
			}
		})
	}
}

func TestMCPConfigProvidersConcurrentReads(t *testing.T) {
	cfg := &config.Config{AI: config.AIConfig{
		MCPSecret:               "secret-0",
		MCPStaticOperatorUserID: 1,
	}}
	settingService := NewSettingService(nil)
	settingService.SetConfig(cfg)

	var wg sync.WaitGroup
	for reader := 0; reader < 8; reader++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 1000; i++ {
				_ = settingService.MCPSecret()
				_ = settingService.MCPStaticOperatorUserID()
			}
		}()
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 1; i <= 1000; i++ {
			settingService.mu.Lock()
			cfg.AI.MCPSecret = fmt.Sprintf("secret-%d", i)
			cfg.AI.MCPStaticOperatorUserID = uint(i)
			settingService.mu.Unlock()
		}
	}()

	wg.Wait()
	if settingService.MCPSecret() != "secret-1000" {
		t.Fatalf("final MCP secret = %q", settingService.MCPSecret())
	}
	if settingService.MCPStaticOperatorUserID() != 1000 {
		t.Fatalf("final MCP operator ID = %d", settingService.MCPStaticOperatorUserID())
	}
}

func TestParseMCPAdminToolsEnabled(t *testing.T) {
	tests := []struct {
		name    string
		raw     string
		want    bool
		wantErr bool
	}{
		{name: "empty defaults off", raw: "", want: false},
		{name: "false", raw: "false", want: false},
		{name: "true", raw: "true", want: true},
		{name: "trim and case", raw: " TRUE ", want: true},
		{name: "reject numeric", raw: "1", wantErr: true},
		{name: "reject invalid", raw: "enabled", wantErr: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := parseMCPAdminToolsEnabled(tc.raw)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("parseMCPAdminToolsEnabled(%q) error = nil", tc.raw)
				}
				return
			}
			if err != nil {
				t.Fatalf("parseMCPAdminToolsEnabled(%q): %v", tc.raw, err)
			}
			if got != tc.want {
				t.Fatalf("parseMCPAdminToolsEnabled(%q) = %v, want %v", tc.raw, got, tc.want)
			}
		})
	}
}

func TestMCPAdminToolsEnabledProviderDefaultsOff(t *testing.T) {
	settingService := &SettingService{}
	if settingService.MCPAdminToolsEnabled() {
		t.Fatal("nil config unexpectedly enabled MCP admin tools")
	}

	cfg := &config.Config{}
	settingService.SetConfig(cfg)
	if settingService.MCPAdminToolsEnabled() {
		t.Fatal("zero-value config unexpectedly enabled MCP admin tools")
	}
	cfg.AI.MCPAdminToolsEnabled = true
	if !settingService.MCPAdminToolsEnabled() {
		t.Fatal("enabled config was not reflected by provider")
	}
}
